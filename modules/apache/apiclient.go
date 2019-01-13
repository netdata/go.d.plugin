package apache

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/web"
)

var validStatusKeys = map[string]string{
	"ServerVersion":                "",
	"ServerMPM":                    "",
	"Server Built":                 "",
	"ParentServerConfigGeneration": "",
	"ParentServerMPMGeneration":    "",
	"ServerUptimeSeconds":          "",
	"ServerUptime":                 "",
	"Load1":                        "",
	"Load5":                        "",
	"Load15":                       "",
	"Total Duration":               "",
	"CPUUser":                      "",
	"CPUSystem":                    "",
	"CPUChildrenUser":              "",
	"CPUChildrenSystem":            "",
	"CPULoad":                      "",
	"DurationPerReq":               "",
	"Processes":                    "",
	"Stopping":                     "",
	"Total Accesses":               "total_accesses",
	"Total kBytes":                 "total_kBytes",
	"Uptime":                       "uptime",
	"ReqPerSec":                    "req_per_sec",
	"BytesPerSec":                  "bytes_per_sec",
	"BytesPerReq":                  "bytes_per_req",
	"BusyWorkers":                  "busy_workers",
	"IdleWorkers":                  "idle_workers",
	"ConnsTotal":                   "conns_total",
	"ConnsAsyncWriting":            "conns_async_writing",
	"ConnsAsyncKeepAlive":          "conns_async_keep_alive",
	"ConnsAsyncClosing":            "conns_async_closing",
	"Scoreboard":                   "scoreboard",
}

type apiClient struct {
	req        web.Request
	httpClient *http.Client
}

func (a apiClient) serverStatus() (map[string]int64, error) {
	req, err := a.createRequest()

	if err != nil {
		return nil, fmt.Errorf("error on creating request : %v", err)
	}

	resp, err := a.doRequestOK(req)

	defer closeBody(resp)

	if err != nil {
		return nil, err
	}

	s := bufio.NewScanner(resp.Body)

	status := make(map[string]string)

	for s.Scan() {
		parts := strings.Split(s.Text(), ":")

		if len(parts) != 2 {
			continue
		}
		status[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	metrics, err := a.parseStatus(status)

	if err != nil {
		return nil, fmt.Errorf("error on parsing status : %v", err)
	}

	return metrics, nil
}

func (a *apiClient) parseStatus(status map[string]string) (map[string]int64, error) {
	metrics := make(map[string]int64)

	for key, value := range status {
		k, ok := validStatusKeys[key]

		if !ok {
			return nil, fmt.Errorf("unknown value : %s", key)
		}

		if k == "" {
			continue
		}

		switch k {
		case "req_per_sec", "bytes_per_sec", "bytes_per_req":
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, err
			}
			metrics[k] = int64(v * 100000)
		case "scoreboard":
			parseScoreboard(value, metrics)
		default:
			v, err := strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
			metrics[k] = int64(v)
		}
	}

	return metrics, nil
}

func (a apiClient) doRequest(req *http.Request) (*http.Response, error) {
	return a.httpClient.Do(req)
}

func (a apiClient) doRequestOK(req *http.Request) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)

	if resp, err = a.doRequest(req); err != nil {
		return resp, fmt.Errorf("error on request to %s : %v", req.URL, err)

	}

	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("%s returned HTTP status %d", req.URL, resp.StatusCode)
	}

	return resp, err
}

func (a apiClient) createRequest() (*http.Request, error) {
	var (
		req *http.Request
		err error
	)

	if req, err = web.NewHTTPRequest(a.req); err != nil {
		return nil, err
	}

	return req, nil
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}

func parseScoreboard(scoreboard string, metrics map[string]int64) {
	//  “_” Waiting for Connection
	// “S” Starting up
	// “R” Reading Request
	// “W” Sending Reply
	// “K” Keepalive (read)
	// “D” DNS Lookup
	// “C” Closing connection
	// “L” Logging
	// “G” Gracefully finishing
	// “I” Idle cleanup of worker
	// “.” Open slot with no current process

	var waiting, open, S, R, W, K, D, C, L, G, I int64

	for _, s := range strings.Split(scoreboard, "") {

		switch s {
		case "_":
			waiting++
		case "S":
			S++
		case "R":
			R++
		case "W":
			W++
		case "K":
			K++
		case "D":
			D++
		case "C":
			C++
		case "L":
			L++
		case "G":
			G++
		case "I":
			I++
		case ".":
			open++
		}
	}

	metrics["scoreboard_waiting"] = waiting
	metrics["scoreboard_starting"] = S
	metrics["scoreboard_reading"] = R
	metrics["scoreboard_sending"] = W
	metrics["scoreboard_keepalive"] = K
	metrics["scoreboard_dns_lookup"] = D
	metrics["scoreboard_closing"] = C
	metrics["scoreboard_logging"] = L
	metrics["scoreboard_finishing"] = G
	metrics["scoreboard_idle_cleanup"] = I
	metrics["scoreboard_open"] = open
}

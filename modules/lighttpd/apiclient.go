package lighttpd

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
	"Total Accesses": "total_accesses",
	"Total kBytes":   "total_kBytes",
	"Uptime":         "uptime",
	"BusyServers":    "busy_servers",
	"IdleServers":    "idle_servers",
	"Scoreboard":     "scoreboard",
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

		switch k {
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
	// Descriptions from https://blog.serverdensity.com/monitor-lighttpd/
	//
	// “.” = Opening the TCP connection (connect)
	// “C” = Closing the TCP connection if no other HTTP request will use it (close)
	// “E” = hard error
	// “k” = Keeping the TCP connection open for more HTTP requests from the same client to avoid the TCP handling overhead (keep-alive)
	// “r” = Read the content of the HTTP request (read)
	// “R” = Read the content of the HTTP request (read-POST)
	// “W” = Write the HTTP response to the socket (write)
	// “h” = Decide action to take with the request (handle-request)
	// “q” = Start of HTTP request (request-start)
	// “Q” = End of HTTP request (request-end)
	// “s” = Start of the HTTP request response (response-start)
	// “S” = End of the HTTP request response (response-end)
	// “_” Waiting for Connection (NOTE: not sure, copied the description from apache score board)

	var waiting, open, C, E, k, r, R, W, h, q, Q, s, S int64

	for _, v := range strings.Split(scoreboard, "") {

		switch v {
		case "_":
			waiting++
		case ".":
			open++
		case "C":
			C++
		case "E":
			E++
		case "k":
			k++
		case "r":
			r++
		case "R":
			R++
		case "W":
			W++
		case "h":
			h++
		case "q":
			q++
		case "Q":
			Q++
		case "s":
			s++
		case "S":
			S++
		}
	}

	metrics["scoreboard_waiting"] = waiting
	metrics["scoreboard_open"] = open
	metrics["scoreboard_close"] = C
	metrics["scoreboard_hard_error"] = E
	metrics["scoreboard_keepalive"] = k
	metrics["scoreboard_read"] = r
	metrics["scoreboard_read_post"] = R
	metrics["scoreboard_write"] = W
	metrics["scoreboard_handle_request"] = h
	metrics["scoreboard_request_start"] = q
	metrics["scoreboard_request_end"] = Q
	metrics["scoreboard_response_start"] = s
	metrics["scoreboard_response_end"] = S
}

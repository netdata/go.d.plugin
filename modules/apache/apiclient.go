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

type (
	scoreboard struct {
		Waiting     int `stm:"waiting"`
		Starting    int `stm:"starting"`
		Reading     int `stm:"reading"`
		Sending     int `stm:"sending"`
		KeepAlive   int `stm:"keepalive"`
		DNSLookup   int `stm:"dns_lookup"`
		Closing     int `stm:"closing"`
		Logging     int `stm:"logging"`
		Finishing   int `stm:"finishing"`
		IdleCleanup int `stm:"idle_cleanup"`
		Open        int `stm:"open"`
	}
	serverStatus struct {
		TotalAccesses       *int        `stm:"total_accesses"`
		TotalKBytes         *int        `stm:"total_kBytes"`
		Uptime              *int        `stm:"uptime"`
		ReqPerSec           *float64    `stm:"req_per_sec"`
		BytesPerSec         *float64    `stm:"bytes_per_sec"`
		BytesPerReq         *float64    `stm:"bytes_per_req"`
		BusyWorkers         *int        `stm:"busy_workers"`
		IdleWorkers         *int        `stm:"idle_workers"`
		ConnsTotal          *int        `stm:"conns_total"`
		ConnsAsyncWriting   *int        `stm:"conns_async_writing"`
		ConnsAsyncKeepAlive *int        `stm:"conns_async_keep_alive"`
		ConnsAsyncClosing   *int        `stm:"conns_async_closing"`
		Scoreboard          *scoreboard `stm:"scoreboard"`
	}
)

func newAPIClient(client *http.Client, request web.Request) *apiClient {
	return &apiClient{httpClient: client, request: request}
}

type apiClient struct {
	httpClient *http.Client
	request    web.Request
}

func (a apiClient) getServerStatus() (*serverStatus, error) {
	req, err := web.NewHTTPRequest(a.request)

	if err != nil {
		return nil, fmt.Errorf("error on creating request : %v", err)
	}

	resp, err := a.doRequestOK(req)

	defer closeBody(resp)

	if err != nil {
		return nil, err
	}

	status, err := parseResponse(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("error on parsing response from %s : %v", req.URL, err)
	}

	return status, nil
}

func (a apiClient) doRequestOK(req *http.Request) (*http.Response, error) {
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error on request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s returned HTTP status %d", req.URL, resp.StatusCode)
	}
	return resp, nil
}

func parseResponse(respBody io.ReadCloser) (*serverStatus, error) {
	s := bufio.NewScanner(respBody)
	status := &serverStatus{}

	for s.Scan() {
		parts := strings.Split(s.Text(), ":")
		if len(parts) != 2 {
			continue
		}
		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])

		switch key {
		default:
		case "BusyServers", "IdleServers":
			return nil, fmt.Errorf("lighttpd data")
		case "BusyWorkers":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.BusyWorkers = &v
			}
		case "IdleWorkers":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.IdleWorkers = &v
			}
		case "ConnsTotal":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.ConnsTotal = &v
			}
		case "ConnsAsyncWriting":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.ConnsAsyncWriting = &v
			}
		case "ConnsAsyncKeepAlive":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.ConnsAsyncKeepAlive = &v
			}
		case "ConnsAsyncClosing":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.ConnsAsyncClosing = &v
			}
		case "Total Accesses":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.TotalAccesses = &v
			}
		case "Total kBytes":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.TotalKBytes = &v
			}
		case "Uptime":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.Uptime = &v
			}
		case "ReqPerSec":
			if v, err := strconv.ParseFloat(value, 64); err != nil {
				return nil, err
			} else {
				v = v * 100000
				status.ReqPerSec = &v
			}
		case "BytesPerSec":
			if v, err := strconv.ParseFloat(value, 64); err != nil {
				return nil, err
			} else {
				v = v * 100000
				status.BytesPerSec = &v
			}
		case "BytesPerReq":
			if v, err := strconv.ParseFloat(value, 64); err != nil {
				return nil, err
			} else {
				v = v * 100000
				status.BytesPerReq = &v
			}
		case "Scoreboard":
			status.Scoreboard = &scoreboard{}
			parseScoreboard(status.Scoreboard, value)
		}
	}

	return status, nil
}

func parseScoreboard(sb *scoreboard, scoreboard string) {
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
	for _, s := range strings.Split(scoreboard, "") {

		switch s {
		case "_":
			sb.Waiting++
		case "S":
			sb.Starting++
		case "R":
			sb.Reading++
		case "W":
			sb.Sending++
		case "K":
			sb.KeepAlive++
		case "D":
			sb.DNSLookup++
		case "C":
			sb.Closing++
		case "L":
			sb.Logging++
		case "G":
			sb.Finishing++
		case "I":
			sb.IdleCleanup++
		case ".":
			sb.Open++
		}
	}
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}

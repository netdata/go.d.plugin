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

type (
	scoreboard struct {
		Waiting       int `stm:"waiting"`
		Open          int `stm:"open"`
		Close         int `stm:"close"`
		HardError     int `stm:"hard_error"`
		KeepAlive     int `stm:"keepalive"`
		Read          int `stm:"read"`
		ReadPost      int `stm:"read_post"`
		Write         int `stm:"write"`
		HandleRequest int `stm:"handle_request"`
		RequestStart  int `stm:"request_start"`
		RequestEnd    int `stm:"request_end"`
		ResponseStart int `stm:"response_start"`
		ResponseEnd   int `stm:"response_end"`
	}
	serverStatus struct {
		TotalAccesses *int        `stm:"total_accesses"`
		TotalKBytes   *int        `stm:"total_kBytes"`
		Uptime        *int        `stm:"uptime"`
		BusyServers   *int        `stm:"busy_servers"`
		IdleServers   *int        `stm:"idle_servers"`
		Scoreboard    *scoreboard `stm:"scoreboard"`
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

func parseResponse(respBody io.Reader) (*serverStatus, error) {
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
		case "BusyWorkers", "IdleWorkers":
			return nil, fmt.Errorf("apache data")
		case "BusyServers":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.BusyServers = &v
			}
		case "IdleServers":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.IdleServers = &v
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
		case "Scoreboard":
			status.Scoreboard = &scoreboard{}
			parseScoreboard(status.Scoreboard, value)
		}
	}

	return status, nil
}

func parseScoreboard(sb *scoreboard, scoreboard string) {
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

	for _, s := range strings.Split(scoreboard, "") {

		switch s {
		case "_":
			sb.Waiting++
		case ".":
			sb.Open++
		case "C":
			sb.Close++
		case "E":
			sb.HardError++
		case "k":
			sb.KeepAlive++
		case "r":
			sb.Read++
		case "R":
			sb.ReadPost++
		case "W":
			sb.Write++
		case "h":
			sb.HandleRequest++
		case "q":
			sb.RequestStart++
		case "Q":
			sb.RequestEnd++
		case "s":
			sb.ResponseStart++
		case "S":
			sb.ResponseEnd++
		}
	}
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}

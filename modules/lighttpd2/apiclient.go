package lighttpd2

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

type serverStatus struct {
	Uptime                       *int `stm:"uptime"`
	RequestsAbs                  *int `stm:"requests_abs"`
	Status1xx                    *int `stm:"status_1xx"`
	Status2xx                    *int `stm:"status_2xx"`
	Status3xx                    *int `stm:"status_3xx"`
	Status4xx                    *int `stm:"status_4xx"`
	Status5xx                    *int `stm:"status_5xx"`
	TrafficInAbs                 *int `stm:"traffic_in_abs"`
	TrafficOutAbs                *int `stm:"traffic_out_abs"`
	ConnectionsAbs               *int `stm:"connections_abs"`
	ConnectionStateStart         *int `stm:"connection_state_start"`
	ConnectionStateReadHeader    *int `stm:"connection_state_read_header"`
	ConnectionStateHandleRequest *int `stm:"connection_state_handle_request"`
	ConnectionStateWriteResponse *int `stm:"connection_state_write_response"`
	ConnectionStateKeepAlive     *int `stm:"connection_state_keepalive"`
	ConnectionStateUpgraded      *int `stm:"connection_state_upgraded"`
	MemoryUsage                  *int `stm:"memory_usage"`
}

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
		line := s.Text()

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])

		switch key {
		default:
		case "requests_abs":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.RequestsAbs = &v
			}
		case "memory_usage":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.MemoryUsage = &v
			}
		case "status_1xx":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.Status1xx = &v
			}
		case "status_2xx":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.Status2xx = &v
			}
		case "status_3xx":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.Status3xx = &v
			}
		case "status_4xx":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.Status4xx = &v
			}
		case "status_5xx":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.Status5xx = &v
			}
		case "traffic_in_abs":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.TrafficInAbs = &v
			}
		case "traffic_out_abs":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.TrafficOutAbs = &v
			}
		case "connections_abs":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.ConnectionsAbs = &v
			}
		case "connection_state_start":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.ConnectionStateStart = &v
			}
		case "connection_state_read_header":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.ConnectionStateReadHeader = &v
			}
		case "connection_state_handle_request":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.ConnectionStateHandleRequest = &v
			}
		case "connection_state_write_response":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.ConnectionStateWriteResponse = &v
			}
		case "connection_state_keep_alive":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.ConnectionStateKeepAlive = &v
			}
		case "connection_state_upgraded":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.ConnectionStateUpgraded = &v
			}
		case "uptime":
			if v, err := strconv.Atoi(value); err != nil {
				return nil, err
			} else {
				status.Uptime = &v
			}
		}
	}

	return status, nil
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}

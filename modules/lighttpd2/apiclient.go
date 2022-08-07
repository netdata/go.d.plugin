// SPDX-License-Identifier: GPL-3.0-or-later

package lighttpd2

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/web"
)

const (
	requestsAbs                  = "requests_abs"
	memoryUsage                  = "memory_usage"
	status1xx                    = "status_1xx"
	status2xx                    = "status_2xx"
	status3xx                    = "status_3xx"
	status4xx                    = "status_4xx"
	status5xx                    = "status_5xx"
	trafficInAbs                 = "traffic_in_abs"
	trafficOutAbs                = "traffic_out_abs"
	connectionsAbs               = "connections_abs"
	connectionStateStart         = "connection_state_start"
	connectionStateReadHeader    = "connection_state_read_header"
	connectionStateHandleRequest = "connection_state_handle_request"
	connectionStateWriteResponse = "connection_state_write_response"
	connectionStateKeepAlive     = "connection_state_keep_alive"
	connectionStateUpgraded      = "connection_state_upgraded"
	uptime                       = "uptime"
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
		return nil, fmt.Errorf("error on request : %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("%s returned HTTP status %d", req.URL, resp.StatusCode)
	}
	return resp, nil
}

func parseResponse(r io.Reader) (*serverStatus, error) {
	s := bufio.NewScanner(r)
	var status serverStatus

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
		case requestsAbs:
			status.Requests.Total = mustParseInt(value)
		case memoryUsage:
			status.Memory.Usage = mustParseInt(value)
		case status1xx:
			status.Responses.Status.Codes1xx = mustParseInt(value)
		case status2xx:
			status.Responses.Status.Codes2xx = mustParseInt(value)
		case status3xx:
			status.Responses.Status.Codes3xx = mustParseInt(value)
		case status4xx:
			status.Responses.Status.Codes4xx = mustParseInt(value)
		case status5xx:
			status.Responses.Status.Codes5xx = mustParseInt(value)
		case trafficInAbs:
			status.Traffic.In = mustParseInt(value)
		case trafficOutAbs:
			status.Traffic.Out = mustParseInt(value)
		case connectionsAbs:
			status.Connection.Total = mustParseInt(value)
		case connectionStateStart:
			status.Connection.State.Start = mustParseInt(value)
		case connectionStateReadHeader:
			status.Connection.State.ReadHeader = mustParseInt(value)
		case connectionStateHandleRequest:
			status.Connection.State.HandleRequest = mustParseInt(value)
		case connectionStateWriteResponse:
			status.Connection.State.WriteResponse = mustParseInt(value)
		case connectionStateKeepAlive:
			status.Connection.State.KeepAlive = mustParseInt(value)
		case connectionStateUpgraded:
			status.Connection.State.Upgraded = mustParseInt(value)
		case uptime:
			status.Uptime = mustParseInt(value)
		}
	}

	return &status, nil
}

func mustParseInt(value string) *int64 {
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(err)
	}
	return &v
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}

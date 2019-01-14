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

var validStatusKeys = map[string]bool{
	"uptime":                          true,
	"memory_usage":                    true,
	"requests_abs":                    true,
	"traffic_out_abs":                 true,
	"traffic_in_abs":                  true,
	"connections_abs":                 true,
	"requests_avg":                    true,
	"traffic_out_avg":                 true,
	"traffic_in_avg":                  true,
	"connections_avg":                 true,
	"requests_avg_5sec":               true,
	"traffic_out_avg_5sec":            true,
	"traffic_in_avg_5sec":             true,
	"connections_avg_5sec":            true,
	"connection_state_start":          true,
	"connection_state_read_header":    true,
	"connection_state_handle_request": true,
	"connection_state_write_response": true,
	"connection_state_keep_alive":     true,
	"connection_state_upgraded":       true,
	"status_1xx":                      true,
	"status_2xx":                      true,
	"status_3xx":                      true,
	"status_4xx":                      true,
	"status_5xx":                      true,
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
		line := s.Text()

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ":")

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

		if _, ok := validStatusKeys[key]; !ok {
			return nil, fmt.Errorf("unknown value : %s", key)
		}

		v, err := strconv.Atoi(value)

		if err != nil {
			return nil, err
		}

		metrics[key] = int64(v)
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

package nginx

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/web"
)

type stubStatus struct {
	Active   int `stm:"active"`
	Requests int `stm:"requests"`
	Reading  int `stm:"reading"`
	Writing  int `stm:"writing"`
	Waiting  int `stm:"waiting"`
	Accepts  int `stm:"accepts"`
	Handled  int `stm:"handled"`
}

type apiClient struct {
	req        web.Request
	httpClient *http.Client
}

func (a apiClient) stubStatus() (stubStatus, error) {
	var status stubStatus

	req, err := a.createRequest()

	if err != nil {
		return status, fmt.Errorf("error on creating request : %v", err)
	}

	resp, err := a.doRequestOK(req)

	defer closeBody(resp)

	if err != nil {
		return status, err
	}

	if err := a.parseResponse(resp, &status); err != nil {
		return status, fmt.Errorf("error on parse response : %v", err)
	}

	return status, nil
}

func (a *apiClient) parseResponse(resp *http.Response, status *stubStatus) error {
	// Active connections: 2
	//server accepts handled requests
	// 2 2 3
	//Reading: 0 Writing: 1 Waiting: 1

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	lines := strings.Split(string(b), "\n")

	if len(lines) < 4 {
		return fmt.Errorf("unparsable data, expected 4 rows, got %d", len(lines))
	}

	if err := a.parseActiveConnections(lines[0], status); err != nil {
		return err
	}

	if err := a.parseAcceptsHandledRequests(lines[2], status); err != nil {
		return err
	}

	if err := a.parseReadWriteWait(lines[3], status); err != nil {
		return err
	}

	return nil
}

func (a *apiClient) parseActiveConnections(line string, status *stubStatus) error {
	slice := strings.Fields(line)

	if len(slice) != 3 {
		return fmt.Errorf("not enough fields in %s", line)
	}

	v, err := strconv.Atoi(slice[2])
	if err != nil {
		return err
	}

	status.Active = v

	return nil
}

func (a *apiClient) parseAcceptsHandledRequests(line string, status *stubStatus) error {
	slice := strings.Fields(line)

	if len(slice) != 3 {
		return fmt.Errorf("not enough fields in %s", line)
	}

	v, err := strconv.Atoi(slice[0])
	if err != nil {
		return err
	}
	status.Accepts = v

	v, err = strconv.Atoi(slice[1])
	if err != nil {
		return err
	}
	status.Handled = v

	v, err = strconv.Atoi(slice[2])
	if err != nil {
		return err
	}
	status.Requests = v

	return nil
}

func (a *apiClient) parseReadWriteWait(line string, status *stubStatus) error {
	slice := strings.Fields(line)

	if len(slice) != 6 {
		return fmt.Errorf("not enough fields in %s", line)
	}

	v, err := strconv.Atoi(slice[1])
	if err != nil {
		return err
	}
	status.Reading = v

	v, err = strconv.Atoi(slice[3])
	if err != nil {
		return err
	}
	status.Writing = v

	v, err = strconv.Atoi(slice[5])
	if err != nil {
		return err
	}
	status.Waiting = v

	return nil
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

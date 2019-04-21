package nginx

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/web"
)

func newAPIClient(client *http.Client, request web.Request) *apiClient {
	return &apiClient{httpClient: client, request: request}
}

type apiClient struct {
	httpClient *http.Client
	request    web.Request
}

func (a apiClient) getStubStatus() ([]string, error) {
	req, err := web.NewHTTPRequest(a.request)

	if err != nil {
		return nil, fmt.Errorf("error on creating request : %v", err)
	}

	resp, err := a.doRequestOK(req)

	defer closeBody(resp)

	if err != nil {
		return nil, err
	}

	s := bufio.NewScanner(resp.Body)

	var lines []string

	for s.Scan() {
		lines = append(lines, strings.Trim(s.Text(), "\r\n "))
	}

	return lines, nil
}

func (a apiClient) doRequestOK(req *http.Request) (*http.Response, error) {
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return resp, fmt.Errorf("error on request : %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("%s returned HTTP status %d", req.URL, resp.StatusCode)
	}

	return resp, err
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}

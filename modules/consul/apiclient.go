package consul

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/netdata/go.d.plugin/pkg/web"
)

type agentCheck struct {
	Node        string
	CheckID     string
	Name        string
	Status      string
	ServiceID   string
	ServiceName string
	ServiceTags []string
}

type apiClient struct {
	aclToken string

	req        web.Request
	httpClient *http.Client
}

func (a *apiClient) localChecks() (map[string]*agentCheck, error) {
	req, err := a.createRequest("/v1/agent/checks")

	if err != nil {
		return nil, fmt.Errorf("error on creating request : %v", err)
	}

	resp, err := a.doRequest(req)

	defer closeBody(resp)

	if err != nil {
		return nil, err
	}

	var checks map[string]*agentCheck

	err = json.NewDecoder(resp.Body).Decode(&checks)

	if err != nil {
		return nil, fmt.Errorf("error on decoding resp from %s : %v", req.URL, err)
	}

	return checks, nil
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

func (a apiClient) createRequest(uri string) (*http.Request, error) {
	var (
		req *http.Request
		err error
	)

	a.req.URI = uri

	if req, err = web.NewHTTPRequest(a.req); err != nil {
		return nil, err
	}

	if a.aclToken != "" {
		req.Header.Set("X-Consul-Token", a.aclToken)
	}

	return req, nil
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}

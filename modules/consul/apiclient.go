package consul

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

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

func newAPIClient(client *http.Client, request web.Request, aclToken string) *apiClient {
	return &apiClient{
		httpClient: client,
		request:    request,
		aclToken:   aclToken,
	}
}

type apiClient struct {
	httpClient *http.Client
	request    web.Request
	aclToken   string
}

func (a *apiClient) localChecks() (map[string]*agentCheck, error) {
	req, err := a.createRequest("/v1/agent/checks")
	if err != nil {
		return nil, fmt.Errorf("error on creating request : %v", err)
	}

	resp, err := a.doRequestOK(req)
	defer closeBody(resp)
	if err != nil {
		return nil, err
	}

	var checks map[string]*agentCheck
	if err = json.NewDecoder(resp.Body).Decode(&checks); err != nil {
		return nil, fmt.Errorf("error on decoding resp from %s : %v", req.URL, err)
	}

	return checks, nil
}

func (a apiClient) doRequestOK(req *http.Request) (*http.Response, error) {
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return resp, fmt.Errorf("error on request to %s : %v", req.URL, err)
	}

	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("%s returned HTTP status %d", req.URL, resp.StatusCode)
	}

	return resp, err
}

func (a apiClient) createRequest(urlPath string) (*http.Request, error) {
	req := a.request.Copy()
	u, err := url.Parse(req.URL)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, urlPath)
	req.URL = u.String()
	if a.aclToken != "" {
		req.Headers["X-Consul-Token"] = a.aclToken
	}
	return web.NewHTTPRequest(req)
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}

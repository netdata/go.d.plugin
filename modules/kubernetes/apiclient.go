package kubernetes

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/mailru/easyjson"
)

const statsSummaryURI = "/stats/summary"

func newAPIClient(client *http.Client, request web.Request) *apiClient {
	return &apiClient{httpClient: client, request: request}
}

type apiClient struct {
	httpClient *http.Client
	request    web.Request
}

func (a apiClient) getStatsSummary() (*statsSummary, error) {
	req, err := a.createRequest()

	if err != nil {
		return nil, fmt.Errorf("error on creating request : %v", err)
	}

	resp, err := a.doRequestOK(req)

	defer closeBody(resp)

	if err != nil {
		return nil, err
	}

	var summary statsSummary

	if err = easyjson.UnmarshalFromReader(resp.Body, &summary); err != nil {
		return nil, fmt.Errorf("error on decoding response from %s : %v", req.URL, err)
	}

	return &summary, nil
}

func (a apiClient) doRequest(req *http.Request) (*http.Response, error) { return a.httpClient.Do(req) }

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
	a.request.URI = statsSummaryURI

	if req, err = web.NewHTTPRequest(a.request); err != nil {
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

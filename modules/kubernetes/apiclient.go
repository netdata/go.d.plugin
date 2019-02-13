package kubernetes

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/mailru/easyjson"
)

const statsSummaryURI = "/stats/summary"

const (
	acceptHeader    = `text/plain;version=0.0.4;q=1,*/*;q=0.1`
	userAgentHeader = `netdata/go.d.plugin`
)

func newAPIClient(client *http.Client, request web.Request) *apiClient {
	return &apiClient{
		httpClient: client,
		request:    request,
		buf:        bytes.NewBuffer(make([]byte, 0, 16000)),
	}
}

type apiClient struct {
	httpClient *http.Client
	request    web.Request

	buf     *bytes.Buffer
	gzipr   *gzip.Reader
	bodybuf *bufio.Reader
}

func (a *apiClient) getStatsSummary() (*statsSummary, error) {
	req, err := a.createRequest(statsSummaryURI)

	if err != nil {
		return nil, fmt.Errorf("error on creating request : %v", err)
	}

	if err = a.fetch(a.buf, req); err != nil {
		return nil, err
	}

	var summary statsSummary

	if err = easyjson.UnmarshalFromReader(a.buf, &summary); err != nil {
		return nil, fmt.Errorf("error on decoding response from %s : %v", req.URL, err)
	}

	return &summary, nil
}

func (a *apiClient) fetch(w io.Writer, req *http.Request) error {
	req, err := a.createRequest(statsSummaryURI)

	if err != nil {
		return fmt.Errorf("error on creating request : %v", err)
	}

	resp, err := a.doRequest(req, true)
	defer closeBody(resp)

	if err != nil {
		return err
	}

	if resp.Header.Get("Content-Encoding") != "gzip" {
		_, err = io.Copy(w, resp.Body)
		return err
	}

	if a.gzipr == nil {
		a.bodybuf = bufio.NewReader(resp.Body)
		a.gzipr, err = gzip.NewReader(a.bodybuf)
		if err != nil {
			return err
		}
	} else {
		a.bodybuf.Reset(resp.Body)
		_ = a.gzipr.Reset(a.bodybuf)
	}
	_, err = io.Copy(w, a.gzipr)
	_ = a.gzipr.Close()

	return err
}

func (a *apiClient) doRequest(req *http.Request, returnOK bool) (*http.Response, error) {
	req.Header.Add("Accept", acceptHeader)
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Set("User-Agent", userAgentHeader)

	resp, err := a.httpClient.Do(req)

	if err != nil {
		return resp, fmt.Errorf("error on request to %s : %v", req.URL, err)
	}

	if returnOK && resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("%s returned HTTP status %d", req.URL, resp.StatusCode)
	}

	return resp, err
}

func (a apiClient) createRequest(uri string) (*http.Request, error) {
	a.request.URI = uri
	return web.NewHTTPRequest(a.request)
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}

package activemq

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/netdata/go.d.plugin/pkg/web"
)

type topics struct {
	XMLName xml.Name `xml:"topics"`
	Items   []topic  `xml:"topic"`
}

type topic struct {
	XMLName xml.Name `xml:"topic"`
	Name    string   `xml:"name,attr"`
	Stats   stats    `xml:"stats"`
}

type queues struct {
	XMLName xml.Name `xml:"queues"`
	Items   []queue  `xml:"queue"`
}

type queue struct {
	XMLName xml.Name `xml:"queue"`
	Name    string   `xml:"name,attr"`
	Stats   stats    `xml:"stats"`
}

type stats struct {
	XMLName       xml.Name `xml:"stats"`
	Size          int64    `xml:"size,attr"`
	ConsumerCount int64    `xml:"consumerCount,attr"`
	EnqueueCount  int64    `xml:"enqueueCount,attr"`
	DequeueCount  int64    `xml:"dequeueCount,attr"`
}

type apiClient struct {
	webadmin   string
	req        web.Request
	httpClient *http.Client
}

func (a *apiClient) getQueues() (*queues, error) {
	req, err := a.createRequest(fmt.Sprintf(uriStats, a.webadmin, keyQueues))

	if err != nil {
		return nil, fmt.Errorf("error on creating request : %v", err)
	}

	resp, err := a.doRequestOK(req)

	defer closeBody(resp)

	if err != nil {
		return nil, err
	}

	var queues queues

	if err := xml.NewDecoder(resp.Body).Decode(&queues); err != nil {
		return nil, fmt.Errorf("error on decoding resp from %s : %s", req.URL, err)
	}

	return &queues, nil
}

func (a *apiClient) getTopics() (*topics, error) {
	req, err := a.createRequest(fmt.Sprintf(uriStats, a.webadmin, keyTopics))

	if err != nil {
		return nil, fmt.Errorf("error on creating request : %v", err)
	}

	resp, err := a.doRequestOK(req)

	defer closeBody(resp)

	if err != nil {
		return nil, err
	}

	var topics topics

	if err := xml.NewDecoder(resp.Body).Decode(&topics); err != nil {
		return nil, fmt.Errorf("error on decoding resp from %s : %s", req.URL, err)
	}

	return &topics, nil
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

	return req, nil
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}

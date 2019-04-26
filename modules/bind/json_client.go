package bind

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/netdata/go.d.plugin/pkg/web"
)

type serverStats = jsonServerStats

type jsonServerStats struct {
	OpCodes   map[string]int64
	QTypes    map[string]int64
	NSStats   map[string]int64
	SockStats map[string]int64
	Views     map[string]jsonView
}

type jsonView struct {
	Resolver jsonViewResolver
}

type jsonViewResolver struct {
	Stats      map[string]int64
	QTypes     map[string]int64
	CacheStats map[string]int64
}

func newJSONClient(client *http.Client, request web.Request) *jsonClient {
	return &jsonClient{httpClient: client, request: request}
}

type jsonClient struct {
	httpClient *http.Client
	request    web.Request
}

func (j jsonClient) serverStats() (*serverStats, error) {
	r := j.request.Copy()
	r.URL.Path = path.Join(r.URL.Path, "/server")
	req, _ := web.NewHTTPRequest(r)

	resp, err := j.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error on request : %v", err)
	}

	defer closeBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s returned HTTP status %d", req.URL, resp.StatusCode)
	}

	stats := &jsonServerStats{}

	if err = json.NewDecoder(resp.Body).Decode(stats); err != nil {
		return nil, fmt.Errorf("error on decoding response from %s : %v", req.URL, err)
	}

	return stats, nil
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}

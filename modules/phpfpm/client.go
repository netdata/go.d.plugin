package phpfpm

import (
	"encoding/json"
	"fmt"
	"github.com/netdata/go-orchestrator/module"
	fcgiclient "github.com/tomasen/fcgi_client"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/netdata/go.d.plugin/pkg/web"
	"time"
)

type (
	status struct {
		Active    int64  `json:"active processes" stm:"active"`
		MaxActive int64  `json:"max active processes" stm:"maxActive"`
		Idle      int64  `json:"idle processes" stm:"idle"`
		Requests  int64  `json:"accepted conn" stm:"requests"`
		Reached   int64  `json:"max children reached" stm:"reached"`
		Slow      int64  `json:"slow requests" stm:"slow"`
		Processes []proc `json:"processes"`
	}
	requestDuration int64
	proc            struct {
		PID      int64           `json:"pid"`
		State    string          `json:"state"`
		Duration requestDuration `json:"request duration"`
		CPU      float64         `json:"last request cpu"`
		Memory   int64           `json:"last request memory"`
	}
)

// UnmarshalJSON customise JSON for timestamp.
func (rd *requestDuration) UnmarshalJSON(b []byte) error {
	if rdc, err := strconv.Atoi(string(b)); err != nil {
		*rd = 0
	} else {
		*rd = requestDuration(rdc)
	}
	return nil
}

type client interface {
	getStatus() (*status, error)
}

type httpClient struct {
	client *http.Client
	req    web.Request
	dec    decoder
}

func (c *httpClient) getStatus() (*status, error) {
	return c.Status()
}

type socketClient struct {
	module.Base
	socket string
	env    map[string]string
}

func (c *socketClient) getStatus() (*status, error) {
	return c.Status()
}

func getStatus(p *Phpfpm) (*status, error) {
	st := &status{}
	if st, err := p.client.getStatus(); st != nil {
		return st, err
	}

	return st, nil
}

func newClient(c *http.Client, r web.Request) *httpClient {
	dec := decodeText
	if _, ok := r.URL.Query()["json"]; ok {
		dec = decodeJSON
	}

	return &httpClient{
		client: c,
		req:    r,
		dec:    dec,
	}
}

func (c httpClient) Status() (*status, error) {
	req, err := web.NewHTTPRequest(c.req)
	if err != nil {
		return nil, fmt.Errorf("error on creating request: %v", err)
	}
	return c.fetchStatus(req)
}

func (c httpClient) fetchStatus(req *http.Request) (*status, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error on request: %v", err)
	}
	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s returned HTTP status %d", req.URL, resp.StatusCode)
	}

	s := &status{}
	if err := c.dec(resp.Body, s); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return s, nil
}

func (c socketClient) Status() (*status, error) {
	var timeout = time.Duration(1) * 100000000
	socket, err := fcgiclient.DialTimeout("unix", c.socket, timeout)
	if err != nil {
		return nil, fmt.Errorf("error on connecting to socket: %v", err)
	}

	return c.fetchStatus(socket)
}

func (c socketClient) fetchStatus(socket *fcgiclient.FCGIClient) (*status, error) {

	resp, err := socket.Get(c.env)
	if err != nil {
		return nil, fmt.Errorf("error on getting data from socket: %v", err)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error on reading socket: %v", err)
	}

	socket.Close()
	st := &status{}

	err2 := json.Unmarshal(content, st)
	if err2 != nil {
		return nil, fmt.Errorf("error on json Unmarshal: %v", err)
	}
	return st, nil
}

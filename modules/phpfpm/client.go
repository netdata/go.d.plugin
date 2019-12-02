package phpfpm

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/netdata/go.d.plugin/pkg/web"
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

type client struct {
	client *http.Client
	req    web.Request
	dec    decoder
}

func newClient(c *http.Client, r web.Request) *client {
	dec := decodeText
	if _, ok := r.URL.Query()["json"]; ok {
		dec = decodeJSON
	}

	return &client{
		client: c,
		req:    r,
		dec:    dec,
	}
}

func (c client) Status() (*status, error) {
	req, err := web.NewHTTPRequest(c.req)
	if err != nil {
		return nil, fmt.Errorf("error on creating request: %v", err)
	}
	return c.fetchStatus(req)
}

func (c client) fetchStatus(req *http.Request) (*status, error) {
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

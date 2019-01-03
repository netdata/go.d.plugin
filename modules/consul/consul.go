package consul

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("consul", creator)
}

const (
	defURL         = "http://127.0.0.1:8500"
	defHTTPTimeout = time.Second
)

// New creates Consul with default values
func New() *Consul {
	return &Consul{
		HTTP: web.HTTP{
			Request: web.Request{URL: defURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defHTTPTimeout}},
		},
	}
}

type agentCheck struct {
	Node        string
	CheckID     string
	Name        string
	Status      string
	ServiceID   string
	ServiceName string
	ServiceTags []string
}

// Consul consul module
type Consul struct {
	modules.Base

	web.HTTP `yaml:",inline"`

	MaxChecks  int `yaml:"max_checks"`
	Token      string
	DataCentre string

	reqChecks *http.Request
	client    *http.Client
}

// Cleanup makes cleanup
func (Consul) Cleanup() {}

// Init makes initialization
func (Consul) Init() bool {
	return false
}

// Check makes check
func (Consul) Check() bool {
	return false
}

// Charts creates Charts
func (Consul) Charts() *Charts {
	return nil
}

// Collect collects metrics
func (c *Consul) Collect() map[string]int64 {
	return nil
}

func (c *Consul) doRequest(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func (c *Consul) doRequestReqOK(req *http.Request) (resp *http.Response, err error) {
	if resp, err = c.doRequest(req); err != nil {
		return resp, fmt.Errorf("error on request to %s : %s", req.URL, err)

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

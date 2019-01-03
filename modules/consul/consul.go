package consul

import (
	"encoding/json"
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

const (
	healthPassing  = "passing"
	healthWarning  = "warning"
	healthCritical = "critical"
	healthMaint    = "maintenance"
)

// New creates Consul with default values
func New() *Consul {
	return &Consul{
		HTTP: web.HTTP{
			Request: web.Request{URL: defURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defHTTPTimeout}},
		},
		activeChecks: make(map[string]bool),
		charts:       &Charts{},
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

	MaxChecks int `yaml:"max_checks"`
	Token     string

	charts       *Charts
	activeChecks map[string]bool
	client       *http.Client
}

// Cleanup makes cleanup
func (Consul) Cleanup() {}

// Init makes initialization
func (c *Consul) Init() bool {
	if c.URL == "" {
		c.Error("URL is not set")
		return false
	}

	c.client = web.NewHTTPClient(c.Client)

	return true
}

// Check makes check
func (c *Consul) Check() bool {
	return len(c.Collect()) > 0
}

// Charts creates Charts
func (c Consul) Charts() *Charts {
	return c.charts
}

// Collect collects metrics
func (c *Consul) Collect() map[string]int64 {
	metrics := make(map[string]int64)

	if err := c.collectLocalChecks(metrics); err != nil {
		c.Error(err)
		return nil
	}

	return metrics
}

func (c *Consul) collectLocalChecks(metrics map[string]int64) error {
	checks, err := c.getLocalChecks()

	if err != nil {
		return err
	}

	c.processLocalChecks(checks, metrics)

	return nil
}

func (c *Consul) getLocalChecks() (map[string]*agentCheck, error) {
	req, err := c.createRequest("/v1/agent/checks")

	if err != nil {
		return nil, fmt.Errorf("error on creating request : %v", err)
	}

	resp, err := c.doRequestReqOK(req)

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

func (c *Consul) processLocalChecks(checks map[string]*agentCheck, metrics map[string]int64) {
	count := len(c.activeChecks)

	for id, check := range checks {
		_, exist := c.activeChecks[id]

		if !exist {
			if c.MaxChecks != 0 && count > c.MaxChecks {
				continue
			}
			c.activeChecks[id] = true
			c.addCheckChart(check)
		}
		metrics[id+"_"+healthPassing] = 0
		metrics[id+"_"+healthCritical] = 0
		metrics[id+"_"+healthMaint] = 0
		metrics[id+"_"+healthWarning] = 0
		metrics[id+"_"+check.Status] = 1
	}
}

func (c *Consul) addCheckChart(check *agentCheck) {
	_ = c.charts.Add(createCheckChart(check))
}

func (c *Consul) doRequest(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func (c *Consul) doRequestReqOK(req *http.Request) (resp *http.Response, err error) {
	if resp, err = c.doRequest(req); err != nil {
		return resp, fmt.Errorf("error on request to %s : %v", req.URL, err)

	}

	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("%s returned HTTP status %d", req.URL, resp.StatusCode)
	}

	return resp, err
}

func (c *Consul) createRequest(uri string) (req *http.Request, err error) {
	c.Request.URI = uri

	if req, err = web.NewHTTPRequest(c.Request); err != nil {
		return
	}

	if c.Token != "" {
		req.Header.Set("X-Consul-Token", c.Token)
	}

	return
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}

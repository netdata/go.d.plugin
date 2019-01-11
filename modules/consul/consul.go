package consul

import (
	"fmt"
	"time"

	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/matcher/notsimplepattern"
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
	defMaxChecks   = 50
)

const (
	healthPassing  = "passing"
	healthWarning  = "warning"
	healthCritical = "critical"
	healthMaint    = "maintenance"
)

// New creates Consul with default values.
func New() *Consul {
	return &Consul{
		HTTP: web.HTTP{
			Request: web.Request{URL: defURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defHTTPTimeout}},
		},
		MaxChecks:    defMaxChecks,
		activeChecks: make(map[string]bool),
		charts:       charts.Copy(),
	}
}

// Consul consul module.
type Consul struct {
	modules.Base

	web.HTTP `yaml:",inline"`

	ACLToken     string `yaml:"acl_token"`
	MaxChecks    int    `yaml:"max_checks"`
	ChecksFilter string `yaml:"checks_filter"`

	charts       *Charts
	activeChecks map[string]bool
	checksFilter *notsimplepattern.Patterns
	apiClient    *apiClient
}

// Cleanup makes cleanup.
func (Consul) Cleanup() {}

// Init makes initialization.
func (c *Consul) Init() bool {
	if c.URL == "" {
		c.Error("URL is not set")
		return false
	}

	c.apiClient = &apiClient{
		aclToken:   c.ACLToken,
		req:        c.Request,
		httpClient: web.NewHTTPClient(c.Client),
	}

	if c.ChecksFilter != "" {
		sps, err := notsimplepattern.Create(c.ChecksFilter)
		if err != nil {
			c.Errorf("error on creating checks filter : %v", err)
			return false
		}

		sps.UseCache = true
		c.checksFilter = sps
	}

	return true
}

// Check makes check.
func (c *Consul) Check() bool {
	return len(c.Collect()) > 0
}

// Charts creates Charts.
func (c Consul) Charts() *Charts {
	return c.charts
}

// Collect collects metrics.
func (c *Consul) Collect() map[string]int64 {
	metrics := make(map[string]int64)

	if err := c.collectLocalChecks(metrics); err != nil {
		c.Error(err)
		return nil
	}

	return metrics
}

func (c *Consul) collectLocalChecks(metrics map[string]int64) error {
	checks, err := c.apiClient.localChecks()

	if err != nil {
		return err
	}

	c.processLocalChecks(checks, metrics)

	return nil
}

func (c *Consul) processLocalChecks(checks map[string]*agentCheck, metrics map[string]int64) {
	count := len(c.activeChecks)
	var unp int

	for id, check := range checks {

		if !c.filterChecks(id) {
			continue
		}

		_, exist := c.activeChecks[id]

		if !exist {
			if c.MaxChecks != 0 && count > c.MaxChecks {
				unp++
				continue
			}

			c.activeChecks[id] = true
			c.addCheckToChart(check)
		}

		var status int64

		switch check.Status {
		case healthPassing, healthMaint:
			status = 0
		case healthWarning:
			status = 1
		case healthCritical:
			status = 2
		default:
			panic(fmt.Sprintf("check %s unkown status %s", check.CheckID, check.Status))
		}
		metrics[id] = status
	}

	if unp > 0 {
		c.Warningf("%d checks were unprocessed due to max_checks limit (%d)", unp, c.MaxChecks)
	}
}

func (c *Consul) filterChecks(name string) bool {
	if c.checksFilter == nil {
		return true
	}
	return c.checksFilter.Match(name)
}

func (c *Consul) addCheckToChart(check *agentCheck) {
	var chart *Chart

	if check.ServiceID != "" {
		chart = c.charts.Get("service_checks")
	} else {
		chart = c.charts.Get("unbound_checks")
	}

	_ = chart.AddDim(&Dim{ID: check.CheckID})
}

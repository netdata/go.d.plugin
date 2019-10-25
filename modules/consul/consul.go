package consul

import (
	"fmt"
	"strings"
	"time"

	"github.com/netdata/go.d.plugin/pkg/matcher"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("consul", creator)
}

const (
	defaultURL         = "http://127.0.0.1:8500"
	defaultHTTPTimeout = time.Second
	defaultMaxChecks   = 50
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
			Request: web.Request{UserURL: defaultURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
		MaxChecks:    defaultMaxChecks,
		activeChecks: make(map[string]bool),
		charts:       charts.Copy(),
	}
}

// Consul Consul module.
type Consul struct {
	module.Base

	web.HTTP `yaml:",inline"`

	ACLToken     string `yaml:"acl_token"`
	MaxChecks    int    `yaml:"max_checks"`
	ChecksFilter string `yaml:"checks_filter"`

	charts       *Charts
	activeChecks map[string]bool
	checksFilter matcher.Matcher
	apiClient    *apiClient
}

// Cleanup makes cleanup.
func (Consul) Cleanup() {}

// Init makes initialization.
func (c *Consul) Init() bool {
	if err := c.Request.ParseUserURL(); err != nil {
		c.Errorf("error on parsing url '%s' : %v", c.UserURL, err)
		return false
	}

	if c.URL.Host == "" {
		c.Error("URL is not set")
		return false
	}

	client, err := web.NewHTTPClient(c.Client)

	if err != nil {
		c.Error(err)
		return false
	}

	c.apiClient = newAPIClient(client, c.Request, c.ACLToken)

	if c.ChecksFilter != "" {
		sps, err := matcher.NewSimplePatternsMatcher(c.ChecksFilter)
		if err != nil {
			c.Errorf("error on creating checks filter : %v", err)
			return false
		}

		c.checksFilter = matcher.WithCache(sps)
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

		_, exist := c.activeChecks[id]

		if !exist {
			if c.MaxChecks != 0 && count > c.MaxChecks {
				unp++
				continue
			}

			if !c.filterChecks(id) {
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
		c.Debugf("%d checks were unprocessed due to max_checks limit (%d)", unp, c.MaxChecks)
	}
}

func (c *Consul) filterChecks(name string) bool {
	if c.checksFilter == nil {
		return true
	}
	return c.checksFilter.MatchString(name)
}

func (c *Consul) addCheckToChart(check *agentCheck) {
	var chart *Chart

	if check.ServiceID != "" {
		chart = c.charts.Get("service_checks")
	} else {
		chart = c.charts.Get("unbound_checks")
	}

	tags := strings.Join(check.ServiceTags, ",")
	_ = chart.AddDim(&Dim{ID: check.CheckID})
	_ = chart.AddDim(&Dim{ID: tags})
	chart.MarkNotCreated()
}

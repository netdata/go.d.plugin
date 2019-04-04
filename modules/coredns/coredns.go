package coredns

import (
	"time"

	"github.com/netdata/go.d.plugin/pkg/matcher"
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

const (
	defaultURL         = "http://127.0.0.1:9153/metrics"
	defaultHTTPTimeout = time.Second * 2
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("coredns", creator)
}

// New creates CoreDNS with default values.
func New() *CoreDNS {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{URL: defaultURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
	}
	return &CoreDNS{
		Config:           config,
		charts:           summaryCharts.Copy(),
		collectedServers: make(map[string]bool),
	}
}

// Config is the CoreDNS module configuration.
type Config struct {
	web.HTTP                 `yaml:",inline"`
	PerServerStatsPermitFrom string `yaml:"per_server_stats_permit_for"`
	//PerZoneStatsPermitFrom   string `yaml:"per_zone_stats_permit_for"`
}

// CoreDNS CoreDNS module.
type CoreDNS struct {
	module.Base
	Config           `yaml:",inline"`
	perServerMatcher matcher.Matcher
	prom             prometheus.Prometheus
	charts           *Charts
	collectedServers map[string]bool
}

// Cleanup makes cleanup.
func (CoreDNS) Cleanup() {}

// Init makes initialization.
func (cd *CoreDNS) Init() bool {
	if cd.URL == "" {
		cd.Error("URL parameter is not set")
		return false
	}

	if cd.PerServerStatsPermitFrom != "" {
		m, err := matcher.Parse(cd.PerServerStatsPermitFrom)
		if err != nil {
			cd.Errorf("error on creating 'per_server_stats_permit_for' matcher from '%s' : %v",
				cd.PerServerStatsPermitFrom, err)
			return false
		}
		cd.perServerMatcher = matcher.WithCache(m)
	}

	client, err := web.NewHTTPClient(cd.Client)
	if err != nil {
		cd.Errorf("error on creating http client : %v", err)
		return false
	}

	cd.prom = prometheus.New(client, cd.Request)

	return true
}

// Check makes check.
func (cd CoreDNS) Check() bool {
	return len(cd.Collect()) > 0
}

// Charts creates Charts.
func (cd CoreDNS) Charts() *Charts {
	return cd.charts
}

// Collect collects metrics.
func (cd *CoreDNS) Collect() map[string]int64 {
	mx, err := cd.collect()

	if err != nil {
		cd.Error(err)
		return nil
	}

	return mx
}

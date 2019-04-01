package coredns

import (
	"time"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

const (
	defaultURL         = "http://127.0.0.1:9253/metrics"
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
		Config: config,
		charts: charts.Copy(),
	}
}

// Config is the CoreDNS module configuration.
type Config struct {
	web.HTTP `yaml:",inline"`
}

// CoreDNS CoreDNS module.
type CoreDNS struct {
	module.Base
	Config `yaml:",inline"`

	prom   prometheus.Prometheus
	charts *Charts
}

// Cleanup makes cleanup.
func (CoreDNS) Cleanup() {}

// Init makes initialization.
func (cd *CoreDNS) Init() bool {
	if cd.URL == "" {
		cd.Error("URL parameter is not set")
		return false
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

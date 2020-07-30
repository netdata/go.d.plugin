package prometheus

import (
	"time"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 5,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("prometheus", creator)
}

func New() *Prometheus {
	config := Config{
		HTTP: web.HTTP{
			Client: web.Client{
				Timeout: web.Duration{Duration: time.Second * 5},
			},
		},
	}
	return &Prometheus{
		Config:      config,
		cache:       make(metricsCache),
		skipMetrics: make(map[string]bool),
		charts:      statsCharts.Copy(),
	}
}

type (
	Config struct {
		web.HTTP `yaml:",inline"`
	}
	Prometheus struct {
		module.Base
		Config      `yaml:",inline"`
		prom        prometheus.Prometheus
		charts      *module.Charts
		cache       metricsCache
		skipMetrics map[string]bool
	}
)

func (Prometheus) Cleanup() {}

func (p *Prometheus) Init() bool {
	if p.UserURL == "" {
		p.Error("URL not set")
		return false
	}

	client, err := web.NewHTTPClient(p.Client)
	if err != nil {
		p.Error(err)
		return false
	}

	p.prom = prometheus.New(client, p.Request)
	return true
}

func (p *Prometheus) Check() bool {
	return len(p.Collect()) > 0
}

func (p Prometheus) Charts() *module.Charts {
	return p.charts
}

func (p *Prometheus) Collect() map[string]int64 {
	mx, err := p.collect()
	if err != nil {
		p.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

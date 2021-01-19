package prometheus

import (
	"time"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/prometheus/selector"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go.d.plugin/agent/module"
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
		MaxTS:          3000,
		MaxTSPerMetric: 200,
	}
	return &Prometheus{
		Config:      config,
		cache:       make(collectCache),
		skipMetrics: make(map[string]bool),
		charts:      statsCharts.Copy(),
	}
}

type (
	Config struct {
		web.HTTP        `yaml:",inline"`
		BearerTokenFile string        `yaml:"bearer_token_file"` // TODO: part of web.Request?
		MaxTS           int           `yaml:"max_time_series"`
		MaxTSPerMetric  int           `yaml:"max_time_series_per_metric"`
		Selector        selector.Expr `yaml:"selector"`
		Grouping        []GroupOption `yaml:"group"`
		ExpectedPrefix  string        `yaml:"expected_prefix"`
	}
	GroupOption struct {
		Selector string `yaml:"selector"`
		ByLabel  string `yaml:"by_label"`
	}

	Prometheus struct {
		module.Base
		Config `yaml:",inline"`

		prom   prometheus.Prometheus
		charts *module.Charts

		optGroupings []optionalGrouping
		cache        collectCache
		skipMetrics  map[string]bool
	}
	optionalGrouping struct {
		sr  selector.Selector
		grp grouper
	}
)

func (Prometheus) Cleanup() {}

func (p *Prometheus) Init() bool {
	if err := p.validateConfig(); err != nil {
		p.Errorf("validating config: %v", err)
		return false
	}

	prom, err := p.initPrometheusClient()
	if err != nil {
		p.Errorf("init prometheus client: %v", err)
		return false
	}
	p.prom = prom

	optGrps, err := p.initOptionalGrouping()
	if err != nil {
		p.Errorf("init grouping: %v", err)
		return false
	}
	p.optGroupings = optGrps

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

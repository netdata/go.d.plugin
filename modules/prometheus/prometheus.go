// SPDX-License-Identifier: GPL-3.0-or-later

package prometheus

import (
	"strings"
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/prometheus/selector"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	module.Register("prometheus", module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 5,
		},
		Create: func() module.Module { return New() },
	})
}

func New() *Prometheus {
	return &Prometheus{
		Config: Config{
			HTTP: web.HTTP{
				Client: web.Client{
					Timeout: web.Duration{Duration: time.Second * 5},
				},
				//Request: web.Request{
				//	URL: "https://node.demo.do.prometheus.io/metrics",
				//},
			},
			MaxTSPerMetric: 200,
		},
		sb:     &strings.Builder{},
		charts: &module.Charts{},
		cache:  newCache(),
	}
}

type Config struct {
	web.HTTP        `yaml:",inline"`
	Name            string        `yaml:"name"`
	Application     string        `yaml:"app"`
	BearerTokenFile string        `yaml:"bearer_token_file"`
	Selector        selector.Expr `yaml:"selector"`
	MaxTSPerMetric  int           `yaml:"max_time_series_per_metric"`
	ExpectedPrefix  string        `yaml:"expected_prefix"`
}

type Prometheus struct {
	module.Base
	Config `yaml:",inline"`

	charts *module.Charts

	prom prometheus.Prometheus
	sb   *strings.Builder

	cache *cache
}

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

	return true
}

func (p *Prometheus) Check() bool {
	return len(p.Collect()) > 0
}

func (p *Prometheus) Charts() *module.Charts {
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

func (p *Prometheus) Cleanup() {}

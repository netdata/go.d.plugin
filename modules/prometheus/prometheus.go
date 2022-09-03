// SPDX-License-Identifier: GPL-3.0-or-later

package prometheus

import (
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/matcher"
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/prometheus/selector"
	"github.com/netdata/go.d.plugin/pkg/web"
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
		Config:       config,
		cache:        make(collectCache),
		skipMetrics:  make(map[string]bool),
		charts:       statsCharts.Copy(),
		firstCollect: true,
	}
}

type (
	Config struct {
		web.HTTP               `yaml:",inline"`
		Name                   string        `yaml:"name"`
		Application            string        `yaml:"app"`
		BearerTokenFile        string        `yaml:"bearer_token_file"` // TODO: part of web.Request?
		MaxTS                  int           `yaml:"max_time_series"`
		MaxTSPerMetric         int           `yaml:"max_time_series_per_metric"`
		Selector               selector.Expr `yaml:"selector"`
		Grouping               []GroupOption `yaml:"group"`
		ExpectedPrefix         string        `yaml:"expected_prefix"`
		ForceAbsoluteAlgorithm []string      `yaml:"force_absolute_algorithm"`
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

		firstCollect           bool
		forceAbsoluteAlgorithm matcher.Matcher
		optGroupings           []optionalGrouping
		cache                  collectCache
		skipMetrics            map[string]bool
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

	mr, err := p.initForceAbsoluteAlgorithm()
	if err != nil {
		p.Errorf("init force_absolute_algorithm (%v): %v", p.ForceAbsoluteAlgorithm, err)
		return false
	}
	p.forceAbsoluteAlgorithm = mr

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

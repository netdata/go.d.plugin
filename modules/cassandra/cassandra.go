// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

import (
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/web"
)

type (
	Config struct {
		web.HTTP `yaml:",inline"`
	}
	Cassandra struct {
		module.Base
		Config `yaml:",inline"`

		prom   prometheus.Prometheus
		charts *module.Charts
	}
)

func init() {
	module.Register("cassandra", module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 5,
		},
		Create: func() module.Module { return New() },
	})
}

func New() *Cassandra {
	return &Cassandra{
		Base:   module.Base{},
		Config: Config{HTTP: web.HTTP{Client: web.Client{Timeout: web.Duration{Duration: time.Second * 5}}}},
		prom:   nil,
		cache:  cache{},
		charts: newCassandraCharts(),
	}
}

func (c *Cassandra) Init() bool {
	if err := c.validateConfig(); err != nil {
		c.Errorf("error on validating config: %v", err)
		return false
	}

	prom, err := c.initPrometheusClient()
	if err != nil {
		c.Errorf("error on init prometheus client: %v", err)
		return false
	}
	c.prom = prom

	return true
}

func (c *Cassandra) Check() bool {
	return len(c.Collect()) > 0
}

func (c *Cassandra) Charts() *Charts {
	return c.charts
}

func (c *Cassandra) Collect() map[string]int64 {
	ms, err := c.collect()
	if err != nil {
		c.Error(err)
	}

	if len(ms) == 0 {
		return nil
	}
	return ms
}

func (Cassandra) Cleanup() {}

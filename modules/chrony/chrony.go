// SPDX-License-Identifier: GPL-3.0-or-later

package chrony

import (
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/modules/chrony/client"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	module.Register("chrony", module.Creator{
		Create: func() module.Module { return New() },
	})
}

func New() *Chrony {
	return &Chrony{
		Config: Config{
			Address: "127.0.0.1:323",
			Timeout: web.Duration{Duration: time.Second},
		},
		charts: charts.Copy(),
		newClient: func(c *Chrony) (chronyClient, error) {
			return client.New(c.Logger, client.Config{
				Address: c.Config.Address,
				Timeout: c.Config.Timeout.Duration,
			})
		},
	}
}

type Config struct {
	Address string       `yaml:"address"`
	Timeout web.Duration `yaml:"timeout"`
}

type (
	Chrony struct {
		module.Base
		Config `yaml:",inline"`

		charts *module.Charts

		newClient func(c *Chrony) (chronyClient, error)
		client    chronyClient
	}
	chronyClient interface {
		Tracking() (*client.TrackingPayload, error)
		Activity() (*client.ActivityPayload, error)
		Close()
	}
)

func (c *Chrony) Init() bool {
	if err := c.validateConfig(); err != nil {
		c.Errorf("config validation: %v", err)
		return false
	}

	return true
}

func (c *Chrony) Check() bool {
	return len(c.Collect()) > 0
}

func (c *Chrony) Charts() *module.Charts {
	return c.charts
}

func (c *Chrony) Collect() map[string]int64 {
	mx, err := c.collect()
	if err != nil {
		c.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (c *Chrony) Cleanup() {
	if c.client != nil {
		c.client.Close()
		c.client = nil
	}
}

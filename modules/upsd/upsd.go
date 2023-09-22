// SPDX-License-Identifier: GPL-3.0-or-later

package upsd

import (
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	module.Register("upsd", module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 5,
		},
		Create: func() module.Module { return New() },
	})
}

func New() *Nut {
	return &Nut{
		Config: Config{
			//Address: "127.0.0.1:3493",
			Address: "192.168.1.200:3493",
			Timeout: web.Duration{Duration: time.Second * 2},
		},
		newNutConn: newNutConn,
		charts:     &module.Charts{},
		upsUnits:   make(map[string]bool),
	}
}

type Config struct {
	Address  string       `yaml:"address"`
	Username string       `yaml:"username"`
	Password string       `yaml:"password"`
	Timeout  web.Duration `yaml:"timeout"`
}

type (
	Nut struct {
		module.Base

		Config `yaml:",inline"`

		charts *module.Charts

		newNutConn func(Config) nutConn
		conn       nutConn

		upsUnits map[string]bool
	}

	nutConn interface {
		connect() error
		disconnect() error
		authenticate(string, string) error
		upsUnits() ([]upsUnit, error)
	}
)

func (n *Nut) Init() bool {
	if n.Address == "" {
		n.Error("config: 'address' not set")
		return false
	}

	return true
}

func (n *Nut) Check() bool {
	return len(n.Collect()) > 0
}

func (n *Nut) Charts() *module.Charts {
	return n.charts
}

func (n *Nut) Collect() map[string]int64 {
	mx, err := n.collect()
	if err != nil {
		n.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (n *Nut) Cleanup() {
	if n.conn != nil {
		return
	}
	if err := n.conn.disconnect(); err != nil {
		n.Warningf("error on disconnect: %v", err)
	}
	n.conn = nil
}

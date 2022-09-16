// SPDX-License-Identifier: GPL-3.0-or-later

package logind

import (
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	module.Register("logind", module.Creator{
		Defaults: module.Defaults{
			Priority: 59999, // copied from the python collector
		},
		Create: func() module.Module { return New() },
	})
}

func New() *Logind {
	return &Logind{
		Config: Config{
			Timeout: web.Duration{Duration: time.Second * 2},
		},
		newLogindConn: func(cfg Config) (logindConnection, error) {
			return newLogindConnection(cfg.Timeout.Duration)
		},
		charts: charts.Copy(),
	}
}

type Config struct {
	Timeout web.Duration `yaml:"timeout"`
}

type Logind struct {
	module.Base
	Config `yaml:",inline"`

	newLogindConn func(config Config) (logindConnection, error)
	conn          logindConnection
	charts        *module.Charts
}

func (l *Logind) Init() bool {
	return true
}

func (l *Logind) Check() bool {
	return len(l.Collect()) > 0
}

func (l *Logind) Charts() *module.Charts {
	return l.charts
}

func (l *Logind) Collect() map[string]int64 {
	mx, err := l.collect()
	if err != nil {
		l.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (l *Logind) Cleanup() {
	if l.conn != nil {
		l.conn.Close()
	}
}

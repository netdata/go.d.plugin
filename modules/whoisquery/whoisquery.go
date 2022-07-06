// SPDX-License-Identifier: GPL-3.0-or-later

package whoisquery

import (
	"errors"
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 60,
		},
		Create: func() module.Module { return New() },
	}
	module.Register("whoisquery", creator)
}

func New() *WhoisQuery {
	return &WhoisQuery{
		Config: Config{
			Timeout:       web.Duration{Duration: time.Second * 5},
			DaysUntilWarn: 90,
			DaysUntilCrit: 30,
		},
	}
}

type Config struct {
	Source        string
	Timeout       web.Duration `yaml:"timeout"`
	DaysUntilWarn int64        `yaml:"days_until_expiration_warning"`
	DaysUntilCrit int64        `yaml:"days_until_expiration_critical"`
}

type WhoisQuery struct {
	module.Base
	Config `yaml:",inline"`
	prov   provider
}

func (wq WhoisQuery) validateConfig() error {
	if wq.Source == "" {
		return errors.New("source is not set")
	}
	return nil
}

func (wq *WhoisQuery) initProvider() error {
	p, err := newProvider(wq.Config)
	if err != nil {
		return err
	}
	wq.prov = p
	return nil
}

func (wq *WhoisQuery) Init() bool {
	if err := wq.validateConfig(); err != nil {
		wq.Errorf("error on validating config: %v", err)
		return false
	}

	if err := wq.initProvider(); err != nil {
		wq.Errorf("error on initializing whois provider: %v", err)
		return false
	}
	return true
}

func (wq *WhoisQuery) Check() bool {
	return len(wq.Collect()) > 0
}

func (wq WhoisQuery) Charts() *Charts {
	return charts.Copy()
}

func (wq *WhoisQuery) Collect() map[string]int64 {
	mx, err := wq.collect()
	if err != nil {
		wq.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (WhoisQuery) Cleanup() {}

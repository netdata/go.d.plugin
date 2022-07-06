// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import (
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	module.Register("wmi", module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 5,
		},
		Create: func() module.Module { return New() },
	})
}

func New() *WMI {
	return &WMI{
		Config: Config{
			HTTP: web.HTTP{
				Client: web.Client{
					Timeout: web.Duration{Duration: time.Second * 5},
				},
			},
		},
		cache: cache{
			collection:   make(map[string]bool),
			collectors:   make(map[string]bool),
			cores:        make(map[string]bool),
			nics:         make(map[string]bool),
			volumes:      make(map[string]bool),
			thermalZones: make(map[string]bool),
		},
		charts: newCollectionCharts(),
	}
}

type (
	Config struct {
		web.HTTP `yaml:",inline"`
	}
	WMI struct {
		module.Base
		Config `yaml:",inline"`
		prom   prometheus.Prometheus
		cache  cache
		charts *Charts
	}
	cache struct {
		collectors   map[string]bool
		collection   map[string]bool
		cores        map[string]bool
		nics         map[string]bool
		volumes      map[string]bool
		thermalZones map[string]bool
	}
)

func (w *WMI) Init() bool {
	if err := w.validateConfig(); err != nil {
		w.Errorf("error on validating config: %v", err)
		return false
	}

	prom, err := w.initPrometheusClient()
	if err != nil {
		w.Errorf("error on init prometheus client: %v", err)
		return false
	}
	w.prom = prom

	return true
}

func (w *WMI) Check() bool {
	return len(w.Collect()) > 0
}

func (w *WMI) Charts() *Charts {
	return w.charts
}

func (w *WMI) Collect() map[string]int64 {
	ms, err := w.collect()
	if err != nil {
		w.Error(err)
	}

	if len(ms) == 0 {
		return nil
	}
	return ms
}

func (WMI) Cleanup() {}

// SPDX-License-Identifier: GPL-3.0-or-later

package windows

import (
	"net/http"
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

func New() *WINDOWS {
	return &WINDOWS{
		Config: Config{
			HTTP: web.HTTP{
				Client: web.Client{
					Timeout: web.Duration{Duration: time.Second * 5},
				},
			},
		},
		cache: cache{
			collection:                  make(map[string]bool),
			collectors:                  make(map[string]bool),
			cores:                       make(map[string]bool),
			nics:                        make(map[string]bool),
			volumes:                     make(map[string]bool),
			thermalZones:                make(map[string]bool),
			processes:                   make(map[string]bool),
			iis:                         make(map[string]bool),
			adcs:                        make(map[string]bool),
			services:                    make(map[string]bool),
			netFrameworkCLRExceptions:   make(map[string]bool),
			netFrameworkCLRInterops:     make(map[string]bool),
			netFrameworkCLRJIT:          make(map[string]bool),
			netFrameworkCLRLoading:      make(map[string]bool),
			netFrameworkCLRLocksThreads: make(map[string]bool),
			netFrameworkCLRMemory:       make(map[string]bool),
			netFrameworkCLRRemoting:     make(map[string]bool),
			netFrameworkCLRSecurity:     make(map[string]bool),
			mssqlInstances:              make(map[string]bool),
			mssqlDBs:                    make(map[string]bool),
		},
		charts: &module.Charts{},
	}
}

type Config struct {
	web.HTTP `yaml:",inline"`
}

type (
	WINDOWS struct {
		module.Base
		Config `yaml:",inline"`

		charts *module.Charts

		doCheck bool

		httpClient *http.Client
		prom       prometheus.Prometheus

		cache cache
	}
	cache struct {
		cores                       map[string]bool
		volumes                     map[string]bool
		nics                        map[string]bool
		thermalZones                map[string]bool
		processes                   map[string]bool
		iis                         map[string]bool
		adcs                        map[string]bool
		mssqlInstances              map[string]bool
		mssqlDBs                    map[string]bool
		services                    map[string]bool
		netFrameworkCLRExceptions   map[string]bool
		netFrameworkCLRInterops     map[string]bool
		netFrameworkCLRJIT          map[string]bool
		netFrameworkCLRLoading      map[string]bool
		netFrameworkCLRLocksThreads map[string]bool
		netFrameworkCLRMemory       map[string]bool
		netFrameworkCLRRemoting     map[string]bool
		netFrameworkCLRSecurity     map[string]bool
		collectors                  map[string]bool
		collection                  map[string]bool
	}
)

func (w *WINDOWS) Init() bool {
	if err := w.validateConfig(); err != nil {
		w.Errorf("config validation: %v", err)
		return false
	}

	httpClient, err := w.initHTTPClient()
	if err != nil {
		w.Errorf("init HTTP client: %v", err)
		return false
	}
	w.httpClient = httpClient

	prom, err := w.initPrometheusClient(w.httpClient)
	if err != nil {
		w.Errorf("init prometheus clients: %v", err)
		return false
	}
	w.prom = prom

	return true
}

func (w *WINDOWS) Check() bool {
	return len(w.Collect()) > 0
}

func (w *WINDOWS) Charts() *module.Charts {
	return w.charts
}

func (w *WINDOWS) Collect() map[string]int64 {
	ms, err := w.collect()
	if err != nil {
		w.Error(err)
	}

	if len(ms) == 0 {
		return nil
	}
	return ms
}

func (w *WINDOWS) Cleanup() {
	if w.httpClient != nil {
		w.httpClient.CloseIdleConnections()
	}
}

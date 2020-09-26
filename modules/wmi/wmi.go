package wmi

import (
	"errors"
	"time"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
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

	module.Register("wmi", creator)
}

// New creates WMI with default values.
func New() *WMI {
	config := Config{
		HTTP: web.HTTP{
			Client: web.Client{
				Timeout: web.Duration{Duration: time.Second * 5},
			},
		},
	}
	return &WMI{
		Config: config,
		cache: cache{
			collection: make(map[string]bool),
			collectors: make(map[string]bool),
			cores:      make(map[string]bool),
			nics:       make(map[string]bool),
			volumes:    make(map[string]bool),
		},
		charts: collectionCharts(),
	}
}

type (
	// Config is the WMI module configuration.
	Config struct {
		web.HTTP `yaml:",inline"`
	}

	// WMI WMI module.
	WMI struct {
		module.Base
		Config `yaml:",inline"`
		prom   prometheus.Prometheus
		cache  cache
		charts *Charts
	}

	cache struct {
		collectors map[string]bool
		collection map[string]bool
		cores      map[string]bool
		nics       map[string]bool
		volumes    map[string]bool
	}
)

func (w *WMI) validateConfig() error {
	if w.URL == "" {
		return errors.New("URL is not set")
	}
	return nil
}

func (w *WMI) initClient() error {
	client, err := web.NewHTTPClient(w.Client)
	if err != nil {
		return err
	}
	w.prom = prometheus.New(client, w.Request)
	return nil
}

func (w *WMI) Init() bool {
	if err := w.validateConfig(); err != nil {
		w.Errorf("error on validating config: %v", err)
		return false
	}

	if err := w.initClient(); err != nil {
		w.Errorf("error on creating prometheus client: %v", err)
		return false
	}
	return true
}

func (w WMI) Check() bool {
	return len(w.Collect()) > 0
}

func (w WMI) Charts() *Charts {
	return w.charts
}

func (w *WMI) Collect() map[string]int64 {
	mx, err := w.collect()
	if err != nil {
		w.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (WMI) Cleanup() {}

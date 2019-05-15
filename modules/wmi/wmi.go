package wmi

import (
	"time"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

const (
	// defaultURL         = "http://127.0.0.1:9182/metrics"
	defaultURL         = "http://100.127.0.251:9182/metrics"
	defaultHTTPTimeout = time.Second * 2
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("wmi", creator)
}

// New creates WMI with default values.
func New() *WMI {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{UserURL: defaultURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
	}
	return &WMI{
		Config: config,
		charts: &Charts{},
		collected: collected{
			collectors: make(map[string]bool),
			cores:      make(map[string]bool),
			nics:       make(map[string]bool),
		},
	}
}

// Config is the WMI module configuration.
type Config struct {
	web.HTTP `yaml:",inline"`
}

type collected struct {
	collectors map[string]bool
	cores      map[string]bool
	nics       map[string]bool
}

// WMI WMI module.
type WMI struct {
	module.Base
	Config `yaml:",inline"`

	charts *Charts
	prom   prometheus.Prometheus

	collected collected
}

// Cleanup makes cleanup.
func (WMI) Cleanup() {}

// Init makes initialization.
func (w *WMI) Init() bool {
	if err := w.ParseUserURL(); err != nil {
		w.Errorf("error on parsing url '%s' : %v", w.UserURL, err)
		return false
	}

	if w.URL.Host == "" {
		w.Error("URL is not set")
		return false
	}

	client, err := web.NewHTTPClient(w.Client)
	if err != nil {
		w.Errorf("error on creating http client : %v", err)
		return false
	}

	w.prom = prometheus.New(client, w.Request)

	return true
}

// Check makes check.
func (w WMI) Check() bool { return len(w.Collect()) > 0 }

// Charts creates Charts.
func (w WMI) Charts() *Charts { return w.charts }

// Collect collects metrics.
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

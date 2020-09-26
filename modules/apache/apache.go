package apache

import (
	"strings"
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go.d.plugin/agent/module"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("apache", creator)
}

const (
	defaultURL         = "http://127.0.0.1/server-status?auto"
	defaultHTTPTimeout = time.Second * 2
)

// New creates Apache with default values.
func New() *Apache {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{
				URL: defaultURL,
			},
			Client: web.Client{
				Timeout: web.Duration{Duration: defaultHTTPTimeout},
			},
		},
	}
	return &Apache{
		Config: config,
		charts: charts.Copy(),
	}
}

// Config is the Apache module configuration.
type Config struct {
	web.HTTP `yaml:",inline"`
}

// Apache Apache module.
type Apache struct {
	module.Base
	Config    `yaml:",inline"`
	apiClient *apiClient
	charts    *Charts
}

// Cleanup makes cleanup.
func (Apache) Cleanup() {}

// Init makes initialization.
func (a *Apache) Init() bool {
	if a.URL == "" {
		a.Error("URL is not set")
		return false
	}

	if !strings.HasSuffix(a.URL, "?auto") {
		a.Errorf("bad URL, should ends in '?auto'")
		return false
	}

	client, err := web.NewHTTPClient(a.Client)
	if err != nil {
		a.Errorf("error on creating http client : %v", err)
		return false
	}

	a.apiClient = newAPIClient(client, a.Request)

	a.Debugf("using URL %s", a.URL)
	a.Debugf("using timeout: %s", a.Timeout.Duration)

	return true
}

// Check makes check.
func (a *Apache) Check() bool {
	mx := a.Collect()

	if len(mx) == 0 {
		return false
	}

	if _, extendedStatus := mx["total_accesses"]; !extendedStatus {
		_ = a.charts.Remove("requests")
		_ = a.charts.Remove("net")
		_ = a.charts.Remove("reqpersec")
		_ = a.charts.Remove("bytespersec")
		_ = a.charts.Remove("bytesperreq")
		_ = a.charts.Remove("uptime")

		a.Warning("'ExtendedStatus' is not enabled, not all metrics collected")
	}

	return true
}

// Charts returns Charts.
func (a Apache) Charts() *module.Charts { return a.charts }

// Collect collects metrics.
func (a *Apache) Collect() map[string]int64 {
	mx, err := a.collect()

	if err != nil {
		a.Error(err)
		return nil
	}

	return mx
}

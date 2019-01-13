package apache

import (
	"time"

	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("apache", creator)
}

const (
	defURL         = "http://localhost/server-status?auto"
	defHTTPTimeout = time.Second
)

// New creates Apache with default values
func New() *Apache {
	return &Apache{
		HTTP: web.HTTP{
			Request: web.Request{URL: defURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defHTTPTimeout}},
		},
		charts: charts.Copy(),
	}
}

// Apache apache module
type Apache struct {
	modules.Base // should be embedded by every module

	web.HTTP `yaml:",inline"`

	apiClient *apiClient
	charts    *Charts
}

// Cleanup makes cleanup
func (Apache) Cleanup() {}

// Init makes initialization
func (a *Apache) Init() bool {
	if a.URL == "" {
		a.Error("URL is nop set")
		return false
	}

	a.apiClient = &apiClient{
		req:        a.Request,
		httpClient: web.NewHTTPClient(a.Client),
	}

	a.Debugf("using URL %s", a.Request.URL)
	a.Debugf("using timeout: %s", a.Timeout.Duration)

	return true
}

// Check makes check
func (a *Apache) Check() bool {
	m, err := a.apiClient.serverStatus()

	if err != nil {
		a.Error(err)
		return false
	}

	if _, ok := m["total_accesses"]; !ok {
		_ = a.charts.Remove("requests")
		_ = a.charts.Remove("net")
		_ = a.charts.Remove("reqpersec")
		_ = a.charts.Remove("bytespersec")
		_ = a.charts.Remove("bytesperreq")

	}

	return len(m) > 0
}

// Charts creates Charts
func (a Apache) Charts() *modules.Charts {
	return a.charts
}

// Collect collects metrics
func (a *Apache) Collect() map[string]int64 {
	var (
		metrics map[string]int64
		err     error
	)

	if metrics, err = a.apiClient.serverStatus(); err != nil {
		a.Error(err)
		return nil
	}

	return metrics
}

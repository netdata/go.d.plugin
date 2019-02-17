package lighttpd

import (
	"strings"
	"time"

	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		DisabledByDefault: true,
		Create:            func() module.Module { return New() },
	}

	module.Register("lighttpd", creator)
}

const (
	defaultURL         = "http://localhost/server-status?auto"
	defaultHTTPTimeout = time.Second * 2
)

// New creates Lighttpd with default values.
func New() *Lighttpd {
	return &Lighttpd{
		HTTP: web.HTTP{
			Request: web.Request{URL: defaultURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
		charts: charts.Copy(),
	}
}

// Lighttpd Lighttpd module.
type Lighttpd struct {
	module.Base
	web.HTTP  `yaml:",inline"`
	apiClient *apiClient
	charts    *Charts
}

// Cleanup makes cleanup.
func (Lighttpd) Cleanup() {}

// Init makes initialization.
func (l *Lighttpd) Init() bool {
	if l.URL == "" {
		l.Error("URL is not set")
		return false
	}

	if !strings.HasSuffix(l.URL, "?auto") {
		l.Errorf("bad URL, should end in '?auto'")
		return false
	}

	client, err := web.NewHTTPClient(l.Client)

	if err != nil {
		l.Errorf("error on creating http client : %v", err)
		return false
	}

	l.apiClient = newAPIClient(client, l.Request)

	l.Debugf("using URL %s", l.URL)
	l.Debugf("using timeout: %s", l.Timeout.Duration)

	return true
}

// Check makes check
func (l *Lighttpd) Check() bool { return len(l.Collect()) > 0 }

// Charts returns Charts.
func (l Lighttpd) Charts() *module.Charts { return l.charts }

// Collect collects metrics.
func (l *Lighttpd) Collect() map[string]int64 {
	status, err := l.apiClient.getServerStatus()

	if err != nil {
		l.Error(err)
		return nil
	}

	return stm.ToMap(status)
}

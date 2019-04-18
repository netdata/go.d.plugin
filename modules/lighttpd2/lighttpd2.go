package lighttpd2

import (
	"strings"
	"time"

	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("lighttpd2", creator)
}

const (
	defaultURL         = "http://localhost/server-status?format=plain"
	defaultHTTPTimeout = time.Second * 2
)

// New creates Lighttpd with default values.
func New() *Lighttpd2 {
	return &Lighttpd2{
		HTTP: web.HTTP{
			Request: web.Request{URL: defaultURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
		charts: charts.Copy(),
	}
}

// Lighttpd2 Lighttpd2 module.
type Lighttpd2 struct {
	module.Base
	web.HTTP  `yaml:",inline"`
	apiClient *apiClient
	charts    *Charts
}

// Cleanup makes cleanup.
func (Lighttpd2) Cleanup() {}

// Init makes initialization.
func (l *Lighttpd2) Init() bool {
	if l.URL == "" {
		l.Error("URL is not set")
		return false
	}

	if !strings.HasSuffix(l.URL, "?format=plain") {
		l.Errorf("bad URL, should end in '?format=plain'")
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
func (l *Lighttpd2) Check() bool { return len(l.Collect()) > 0 }

// Charts returns Charts.
func (l Lighttpd2) Charts() *module.Charts { return l.charts }

// Collect collects metrics.
func (l *Lighttpd2) Collect() map[string]int64 {
	status, err := l.apiClient.getServerStatus()

	if err != nil {
		l.Error(err)
		return nil
	}

	return stm.ToMap(status)
}

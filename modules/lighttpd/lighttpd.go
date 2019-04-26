package lighttpd

import (
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			Disabled: true,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("lighttpd", creator)
}

const (
	defaultURL         = "http://127.0.0.1/server-status?auto"
	defaultHTTPTimeout = time.Second * 2
)

// New creates Lighttpd with default values.
func New() *Lighttpd {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{UserURL: defaultURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
	}
	return &Lighttpd{Config: config}
}

// Config is the Lighttpd module configuration.
type Config struct {
	web.HTTP `yaml:",inline"`
}

// Lighttpd Lighttpd module.
type Lighttpd struct {
	module.Base
	Config    `yaml:",inline"`
	apiClient *apiClient
}

// Cleanup makes cleanup.
func (Lighttpd) Cleanup() {}

// Init makes initialization.
func (l *Lighttpd) Init() bool {
	if err := l.ParseUserURL(); err != nil {
		l.Errorf("error on parsing url '%s' : %v", l.Request.UserURL, err)
		return false
	}

	if l.URL.Host == "" {
		l.Error("URL is not set")
		return false
	}

	if l.URL.RawQuery != "auto" {
		l.Errorf("bad URL, should ends in '?auto'")
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
func (l Lighttpd) Charts() *Charts { return charts.Copy() }

// Collect collects metrics.
func (l *Lighttpd) Collect() map[string]int64 {
	mx, err := l.collect()

	if err != nil {
		l.Error(err)
		return nil
	}

	return mx
}

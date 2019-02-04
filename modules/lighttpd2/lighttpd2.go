package lighttpd2

import (
	"strings"
	"time"

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
	defaultHTTPTimeout = time.Second
)

// New creates Lighttpd2 with default values
func New() *Lighttpd2 {
	return &Lighttpd2{
		HTTP: web.HTTP{
			Request: web.Request{URL: defaultURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
	}
}

// Lighttpd2 lighttpd module
type Lighttpd2 struct {
	module.Base // should be embedded by every module

	web.HTTP `yaml:",inline"`

	apiClient *apiClient
}

// Cleanup makes cleanup
func (Lighttpd2) Cleanup() {}

// Init makes initialization
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
		l.Error(err)
		return false
	}

	l.apiClient = &apiClient{
		req:        l.Request,
		httpClient: client,
	}

	l.Debugf("using URL %s", l.URL)
	l.Debugf("using timeout: %s", l.Timeout.Duration)

	return true
}

// Check makes check
func (l *Lighttpd2) Check() bool {
	return len(l.Collect()) > 0
}

// Charts creates Charts
func (l Lighttpd2) Charts() *Charts {
	return charts.Copy()
}

// Collect collects metrics
func (l *Lighttpd2) Collect() map[string]int64 {
	var (
		metrics map[string]int64
		err     error
	)

	if metrics, err = l.apiClient.serverStatus(); err != nil {
		l.Error(err)
		return nil
	}

	return metrics
}

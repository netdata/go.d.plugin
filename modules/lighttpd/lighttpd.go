package lighttpd

import (
	"strings"
	"time"

	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("lighttpd", creator)
}

const (
	defaultURL         = "http://localhost/server-status?auto"
	defaultHTTPTimeout = time.Second
)

// New creates Lighttpd with default values
func New() *Lighttpd {

	return &Lighttpd{
		HTTP: web.HTTP{
			Request: web.Request{URL: defaultURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
	}
}

// Lighttpd lighttpd module
type Lighttpd struct {
	modules.Base // should be embedded by every module

	web.HTTP `yaml:",inline"`

	apiClient *apiClient
}

// Cleanup makes cleanup
func (Lighttpd) Cleanup() {}

// Init makes initialization
func (l *Lighttpd) Init() bool {
	if l.URL == "" {
		l.Error("URL is not set")
		return false
	}

	if !strings.HasSuffix(l.URL, "?auto") {
		l.Errorf("bad URL, should end in '?auto'")
		return false
	}

	l.apiClient = &apiClient{
		req:        l.Request,
		httpClient: web.NewHTTPClient(l.Client),
	}

	l.Debugf("using URL %s", l.URL)
	l.Debugf("using timeout: %s", l.Timeout.Duration)

	return true
}

// Check makes check
func (l *Lighttpd) Check() bool {
	return len(l.Collect()) > 0
}

// Charts creates Charts
func (l Lighttpd) Charts() *modules.Charts {
	return charts.Copy()
}

// Collect collects metrics
func (l *Lighttpd) Collect() map[string]int64 {
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

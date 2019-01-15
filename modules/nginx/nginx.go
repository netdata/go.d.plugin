package nginx

import (
	"time"

	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("nginx", creator)
}

const (
	defaultURL         = "http://localhost/stub_status"
	defaultHTTPTimeout = time.Second
)

// New creates Nginx with default values
func New() *Nginx {
	return &Nginx{
		HTTP: web.HTTP{
			Request: web.Request{URL: defaultURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
	}
}

// Nginx nginx module
type Nginx struct {
	modules.Base // should be embedded by every module

	web.HTTP `yaml:",inline"`

	apiClient *apiClient
}

// Cleanup makes cleanup
func (Nginx) Cleanup() {}

// Init makes initialization
func (n *Nginx) Init() bool {
	if n.URL == "" {
		n.Error("URL is not set")
		return false
	}

	n.apiClient = &apiClient{
		req:        n.Request,
		httpClient: web.NewHTTPClient(n.Client),
	}

	n.Debugf("using URL %s", n.URL)
	n.Debugf("using timeout: %s", n.Timeout.Duration)

	return true
}

// Check makes check
func (n *Nginx) Check() bool {
	return len(n.Collect()) > 0
}

// Charts creates Charts
func (Nginx) Charts() *Charts {
	return charts.Copy()
}

// Collect collects metrics
func (n *Nginx) Collect() map[string]int64 {
	status, err := n.apiClient.stubStatus()

	if err != nil {
		n.Error(err)
		return nil
	}

	return stm.ToMap(status)
}

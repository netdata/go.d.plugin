package nginx

import (
	"time"

	"github.com/netdata/go.d.plugin/pkg/stm"
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

	module.Register("nginx", creator)
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
	module.Base // should be embedded by every module

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

	client, err := web.NewHTTPClient(n.Client)

	if err != nil {
		n.Error(err)
		return false
	}

	n.apiClient = &apiClient{
		req:        n.Request,
		httpClient: client,
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

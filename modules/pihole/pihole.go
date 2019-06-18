package pihole

import (
	"time"

	"github.com/netdata/go.d.plugin/modules/pihole/client"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("pihole", creator)
}

const (
	defaultURL         = "http://192.168.88.228"
	defaultHTTPTimeout = time.Second
	defaultTopClients  = 5
	defaultTopItems    = 5
	//defaultSetupVarsPath = "/etc/pihole/setupVars.conf"
	defaultSetupVarsPath = "/opt/other/setupVars.conf1"
)

// New creates Pihole with default values.
func New() *Pihole {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{UserURL: defaultURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
	}

	return &Pihole{
		Config: config,
	}
}

// Config is the Pihole module configuration.
type Config struct {
	web.HTTP `yaml:",inline"`
}

// Pihole Pihole module.
type Pihole struct {
	module.Base
	Config `yaml:",inline"`

	charts *module.Charts
	client *client.Client
}

// Cleanup makes cleanup.
func (Pihole) Cleanup() {}

// Init makes initialization.
func (p *Pihole) Init() bool {
	httpClient, err := web.NewHTTPClient(p.Client)
	if err != nil {
		p.Errorf("error on creating http client : %v", err)
		return false
	}

	pass := p.webPassword()
	if pass == "" {
		p.Warning("no web password, not all metrics available")
	} else {
		p.Debugf("web password: %s", pass)
	}

	config := client.Configuration{
		Client:      httpClient,
		URL:         p.UserURL,
		WebPassword: pass,
	}
	p.client = client.New(config)

	return true
}

// Check makes check.
func (Pihole) Check() bool { return true }

// Charts returns Charts.
func (Pihole) Charts() *module.Charts { return charts.Copy() }

// Collect collects metrics.
func (p *Pihole) Collect() map[string]int64 {
	mx, err := p.collect()

	if err != nil {
		p.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}

	return mx
}

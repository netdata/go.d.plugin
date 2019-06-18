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
	defaultSetupVarsPath = "/opt/other/setupVars.conf"
)

// New creates Pihole with default values.
func New() *Pihole {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{UserURL: defaultURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
		SetupVarsPath: defaultSetupVarsPath,
	}

	return &Pihole{
		Config: config,
		charts: authCharts.Copy(),
		collected: collected{
			forwarders: make(map[string]bool),
			topClients: make(map[string]bool),
			topDomains: make(map[string]bool),
			topAds:     make(map[string]bool),
		},
	}
}

type piholeAPIClient interface {
	Version() (int, error)
	SummaryRaw() (*client.SummaryRaw, error)
	QueryTypes() (*client.QueryTypes, error)
	ForwardDestinations() (*[]client.ForwardDestination, error)
	TopClients(top int) (*[]client.TopClient, error)
	TopItems(top int) (*client.TopItems, error)
}

type collected struct {
	forwarders map[string]bool
	topClients map[string]bool
	topDomains map[string]bool
	topAds     map[string]bool
}

// Config is the Pihole module configuration.
type Config struct {
	web.HTTP      `yaml:",inline"`
	SetupVarsPath string `yaml:"setup_vars_path"`
}

// Pihole Pihole module.
type Pihole struct {
	module.Base
	Config `yaml:",inline"`

	collected collected
	charts    *module.Charts
	client    piholeAPIClient
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

	p.Password = p.webPassword()
	if p.Password == "" {
		p.Warning("no web password, not all metrics available")
	} else {
		p.Debugf("web password: %s", p.Password)
	}

	config := client.Configuration{
		Client:      httpClient,
		URL:         p.UserURL,
		WebPassword: p.Password,
	}
	p.client = client.New(config)

	return true
}

// Check makes check.
func (Pihole) Check() bool { return true }

// Charts returns Charts.
func (p Pihole) Charts() *module.Charts { return p.charts }

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

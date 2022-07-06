// SPDX-License-Identifier: GPL-3.0-or-later

package pihole

import (
	"time"

	"github.com/netdata/go.d.plugin/modules/pihole/client"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go.d.plugin/agent/module"
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 5,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("pihole", creator)
}

const supportedAPIVersion = 3

const (
	defaultURL           = "http://127.0.0.1"
	defaultHTTPTimeout   = time.Second * 5
	defaultTopClients    = 5
	defaultTopItems      = 5
	defaultSetupVarsPath = "/etc/pihole/setupVars.conf"
)

// New creates Pihole with default values.
func New() *Pihole {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{
				URL: defaultURL,
			},
			Client: web.Client{
				Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
		SetupVarsPath:     defaultSetupVarsPath,
		TopClientsEntries: defaultTopClients,
		TopItemsEntries:   defaultTopItems,
	}

	return &Pihole{
		Config: config,
		charts: charts.Copy(),
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

// Config is the Pihole module configuration.
type Config struct {
	web.HTTP          `yaml:",inline"`
	SetupVarsPath     string `yaml:"setup_vars_path"`
	TopClientsEntries int    `yaml:"top_clients_entries"`
	TopItemsEntries   int    `yaml:"top_items_entries"`
}

// Pihole Pihole module.
type Pihole struct {
	module.Base
	Config `yaml:",inline"`

	charts *module.Charts
	client piholeAPIClient
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
		URL:         p.URL,
		WebPassword: p.Password,
	}
	p.client = client.New(config)

	return true
}

// Check makes check.
func (p Pihole) Check() bool {
	ver, err := p.client.Version()
	if err != nil {
		p.Error(err)
		return false
	}

	if ver != supportedAPIVersion {
		p.Errorf("API version: %d, supported version: %d", ver, supportedAPIVersion)
		return false
	}

	if p.Password != "" {
		// TODO: remove panic
		panicIf(p.charts.Add(*authCharts.Copy()...))
	}

	return len(p.Collect()) > 0
}

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

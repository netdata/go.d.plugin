package pihole

import (
	"github.com/netdata/go.d.plugin/modules/pihole/client"
	"time"

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
)

// New creates Pihole with default values.
func New() *Pihole {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{UserURL: defaultURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
	}

	return &Pihole{Config: config}
}

// Config is the Pihole module configuration.
type Config struct {
	web.HTTP `yaml:",inline"`
}

// Pihole Pihole module.
type Pihole struct {
	module.Base
	Config `yaml:",inline"`

	client *client.Client
}

// Cleanup makes cleanup.
func (Pihole) Cleanup() {}

// Init makes initialization.
func (p *Pihole) Init() bool {
	c, err := client.New(p.Client, p.Request)
	if err != nil {
		p.Error(err)
		return false
	}

	p.client = c
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

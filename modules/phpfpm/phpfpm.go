package phpfpm

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

	module.Register("phpfpm", creator)
}

// New returns a php-fpm module with default values.
func New() *Phpfpm {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{UserURL: "http://127.0.0.1/status?full&json"},
			Client:  web.Client{Timeout: web.Duration{Duration: time.Second}},
		},
	}

	return &Phpfpm{
		Config: config,
	}
}

type (
	// Config is the php-fpm module configuration.
	Config struct {
		web.HTTP `yaml:",inline"`
	}
	// Phpfpm collets php-fpm metrics.
	Phpfpm struct {
		module.Base
		Config `yaml:",inline"`
		client *client
	}
)

// Init makes initialization.
func (p *Phpfpm) Init() bool {
	if err := p.ParseUserURL(); err != nil {
		p.Errorf("error on parsing url '%s' : %v", p.UserURL, err)
		return false
	}
	if p.URL.Host == "" {
		p.Error("URL is not set")
		return false
	}

	client, err := web.NewHTTPClient(p.Client)
	if err != nil {
		p.Error(err)
		return false
	}
	p.client = newClient(client, p.Request)

	p.Debugf("using URL %s", p.URL)
	p.Debugf("using timeout: %s", p.Timeout.Duration)
	return true
}

// Check checks the module can collect metrics.
func (p *Phpfpm) Check() bool {
	return len(p.Collect()) > 0
}

// Charts creates Charts.
func (Phpfpm) Charts() *Charts {
	return charts.Copy()
}

// Collect returns collected metrics.
func (p *Phpfpm) Collect() map[string]int64 {
	mx, err := p.collect()
	if err != nil {
		p.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

// Cleanup makes cleanup.
func (Phpfpm) Cleanup() {}

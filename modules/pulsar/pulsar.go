package pulsar

import (
	"errors"
	"time"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			// UpdateEvery: 60, // TODO: uncomment
		},
		Create: func() module.Module { return New() },
	}

	module.Register("pulsar", creator)
}

func New() *Pulsar {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{
				UserURL: "http://127.0.0.1:8080/metrics",
			},
			Client: web.Client{Timeout: web.Duration{Duration: time.Second}},
		},
	}

	return &Pulsar{
		Config: config,
		charts: nil, // TODO
		cache:  make(cache),
	}
}

type (
	Config struct {
		web.HTTP `yaml:",inline"`
	}

	Pulsar struct {
		module.Base
		Config `yaml:",inline"`

		prom   prometheus.Prometheus
		charts *Charts
		cache  cache
	}

	cache map[string]bool
)

func (c cache) hasP(v string) bool { ok := c[v]; c[v] = true; return ok }

func (p Pulsar) validateConfig() error {
	if p.UserURL == "" {
		return errors.New("URL is not set")
	}
	return nil
}

func (p *Pulsar) initClient() error {
	client, err := web.NewHTTPClient(p.Client)
	if err != nil {
		return err
	}

	p.prom = prometheus.New(client, p.Request)
	return nil
}

func (p *Pulsar) Init() bool {
	if err := p.validateConfig(); err != nil {
		p.Errorf("error on validating config: %v", err)
		return false
	}
	if err := p.initClient(); err != nil {
		p.Errorf("error on initializing client: %v", err)
		return false
	}
	return true
}

func (p *Pulsar) Check() bool {
	return len(p.Collect()) > 0
}

func (p *Pulsar) Charts() *Charts {
	return p.charts
}

func (p *Pulsar) Collect() map[string]int64 {
	mx, err := p.collect()
	if err != nil {
		p.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (Pulsar) Cleanup() {}

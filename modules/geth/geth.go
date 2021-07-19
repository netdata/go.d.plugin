package geth

import (
	"errors"
	"time"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go.d.plugin/agent/module"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("geth", creator)
}

func New() *Geth {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{
				URL: "http://127.0.0.1:6060/debug/metrics/prometheus",
			},
			Client: web.Client{
				Timeout: web.Duration{Duration: time.Second},
			},
		},
	}

	return &Geth{
		Config: config,
		charts: charts.Copy(),
	}
}

type (
	Config struct {
		web.HTTP `yaml:",inline"`
	}

	Geth struct {
		module.Base
		Config `yaml:",inline"`

		prom   prometheus.Prometheus
		charts *Charts
	}
)

func (v Geth) validateConfig() error {
	if v.URL == "" {
		return errors.New("URL is not set")
	}
	return nil
}

func (v *Geth) initClient() error {
	client, err := web.NewHTTPClient(v.Client)
	if err != nil {
		return err
	}

	v.prom = prometheus.New(client, v.Request)
	return nil
}

func (v *Geth) Init() bool {
	if err := v.validateConfig(); err != nil {
		v.Errorf("error on validating config: %v", err)
		return false
	}
	if err := v.initClient(); err != nil {
		v.Errorf("error on initializing client: %v", err)
		return false
	}
	return true
}

func (v *Geth) Check() bool {
	return len(v.Collect()) > 0
}

func (v *Geth) Charts() *Charts {
	return v.charts
}

func (v *Geth) Collect() map[string]int64 {
	mx, err := v.collect()
	if err != nil {
		v.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (Geth) Cleanup() {}

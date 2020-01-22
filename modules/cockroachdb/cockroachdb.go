package cockroachdb

import (
	"time"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("cockroachdb", creator)
}

func New() *CockroachDB {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{
				UserURL: "http://127.0.0.1:8080/_status/vars",
			},
			Client: web.Client{Timeout: web.Duration{Duration: time.Second}},
		},
	}

	return &CockroachDB{
		Config: config,
	}
}

type (
	Config struct {
		web.HTTP `yaml:",inline"`
	}

	CockroachDB struct {
		module.Base
		Config `yaml:",inline"`

		prom   prometheus.Prometheus
		charts *Charts
	}
)

func (c *CockroachDB) validateConfig() error {
	return nil
}

func (c *CockroachDB) createClient() error {
	client, err := web.NewHTTPClient(c.Client)
	if err != nil {
		return err
	}

	c.prom = prometheus.New(client, c.Request)
	return nil
}

func (c *CockroachDB) Init() bool {
	if err := c.validateConfig(); err != nil {
		c.Errorf("error on validating config: %v", err)
		return false
	}

	if err := c.createClient(); err != nil {
		c.Errorf("error on creating client: %v", err)
		return false
	}
	return true
}

func (c *CockroachDB) Check() bool {
	return len(c.Collect()) > 0
}

func (c CockroachDB) Charts() *Charts {
	return c.charts
}

func (c *CockroachDB) Collect() map[string]int64 {
	mx, err := c.collect()
	if err != nil {
		c.Error(err)
		return nil
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (CockroachDB) Cleanup() {}

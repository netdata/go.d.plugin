package traefik

import (
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	module.Register("traefik", module.Creator{
		Create: func() module.Module { return New() },
	})
}

func New() *Traefik {
	return &Traefik{
		Config: Config{
			HTTP: web.HTTP{
				Request: web.Request{
					URL: "http://127.0.0.1:9999/traefik",
				},
				Client: web.Client{
					Timeout: web.Duration{Duration: time.Second},
				},
			},
		},

		charts:       &module.Charts{},
		checkMetrics: true,
		cache: &cache{
			entrypoints: make(map[string]*cacheEntrypoint),
		},
	}
}

type Config struct {
	web.HTTP `yaml:",inline"`
}

type (
	Traefik struct {
		module.Base
		Config `yaml:",inline"`

		prom         prometheus.Prometheus
		charts       *module.Charts
		checkMetrics bool
		cache        *cache
	}
	cache struct {
		entrypoints map[string]*cacheEntrypoint
	}
	cacheEntrypoint struct {
		requests        *module.Chart
		reqDur          *module.Chart
		openConn        *module.Chart
		openConnMethods map[string]bool
	}
)

func (t *Traefik) Init() bool {
	if err := t.validateConfig(); err != nil {
		t.Errorf("config validation: %v", err)
		return false
	}

	prom, err := t.initPrometheusClient()
	if err != nil {
		t.Errorf("prometheus client initialization: %v", err)
		return false
	}
	t.prom = prom

	return true
}

func (t *Traefik) Check() bool {
	return len(t.Collect()) > 0
}

func (t *Traefik) Charts() *module.Charts {
	return t.charts
}

func (t *Traefik) Collect() map[string]int64 {
	mx, err := t.collect()
	if err != nil {
		t.Error(err)
		return nil
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (Traefik) Cleanup() {}

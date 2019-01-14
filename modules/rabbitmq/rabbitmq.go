package rabbitmq

import (
	"time"

	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("rabbitmq", creator)
}

var (
	defURL         = "http://localhost:15672"
	defUsername    = "guest"
	defPassword    = "guest"
	defHTTPTimeout = time.Second
)

// New creates Rabbitmq with default values
func New() *Rabbitmq {
	return &Rabbitmq{
		HTTP: web.HTTP{
			Request: web.Request{
				URL:      defURL,
				Username: defUsername,
				Password: defPassword,
			},
			Client: web.Client{Timeout: web.Duration{Duration: defHTTPTimeout}},
		},
	}
}

// Rabbitmq rabbitmq module.
type Rabbitmq struct {
	modules.Base

	web.HTTP `yaml:",inline"`

	apiClient *apiClient
}

// Cleanup makes cleanup.
func (Rabbitmq) Cleanup() {}

// Init makes initialization.
func (r *Rabbitmq) Init() bool {
	if r.URL == "" {
		r.Error("URL is not set")
		return false
	}

	r.apiClient = &apiClient{
		req:        r.Request,
		httpClient: web.NewHTTPClient(r.Client),
	}

	r.Debugf("using URL %s", r.URL)
	r.Debugf("using timeout: %s", r.Timeout.Duration)

	return true
}

// Check makes check.
func (r *Rabbitmq) Check() bool {
	return len(r.Collect()) > 0
}

// Charts creates Charts.
func (Rabbitmq) Charts() *Charts {
	return charts.Copy()
}

// Collect collects stats.
func (r *Rabbitmq) Collect() map[string]int64 {
	var (
		metrics  = make(map[string]int64)
		overview overview
		node     node
		err      error
	)

	if overview, err = r.apiClient.getOverview(); err != nil {
		r.Error(err)
		return nil
	}

	if node, err = r.apiClient.getNodeStats(); err != nil {
		r.Error(err)
		return nil
	}

	for k, v := range stm.ToMap(overview) {
		metrics[k] = v
	}

	for k, v := range stm.ToMap(node) {
		metrics[k] = v
	}

	return metrics
}

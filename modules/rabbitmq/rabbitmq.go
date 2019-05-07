package rabbitmq

import (
	"time"

	"github.com/netdata/go.d.plugin/pkg/stm"
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

	module.Register("rabbitmq", creator)
}

const (
	defaultURL         = "http://localhost:15672"
	defaultUsername    = "guest"
	defaultPassword    = "guest"
	defaultHTTPTimeout = time.Second
)

// New creates RabbitMQ with default values.
func New() *RabbitMQ {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{
				UserURL:  defaultURL,
				Username: defaultUsername,
				Password: defaultPassword,
			},
			Client: web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
	}

	return &RabbitMQ{Config: config}
}

// Config is the RabbitMQ module configuration.
type Config struct {
	web.HTTP `yaml:",inline"`
}

// RabbitMQ RabbitMQ module.
type RabbitMQ struct {
	module.Base
	Config    `yaml:",inline"`
	apiClient *apiClient
}

// Cleanup makes cleanup.
func (RabbitMQ) Cleanup() {}

// Init makes initialization.
func (r *RabbitMQ) Init() bool {
	if err := r.ParseUserURL(); err != nil {
		r.Errorf("error on parsing url '%s' : %v", r.UserURL, err)
		return false
	}

	if r.URL.Host == "" {
		r.Error("URL is not set")
		return false
	}

	client, err := web.NewHTTPClient(r.Client)

	if err != nil {
		r.Error(err)
		return false
	}

	r.apiClient = newAPIClient(client, r.Request)

	r.Debugf("using URL %s", r.URL)
	r.Debugf("using timeout: %s", r.Timeout.Duration)

	return true
}

// Check makes check.
func (r *RabbitMQ) Check() bool {
	return len(r.Collect()) > 0
}

// Charts creates Charts.
func (RabbitMQ) Charts() *Charts {
	return charts.Copy()
}

// Collect collects stats.
func (r *RabbitMQ) Collect() map[string]int64 {
	var (
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

	return stm.ToMap(overview, node)
}

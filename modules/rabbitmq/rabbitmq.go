package rabbitmq

import (
	"fmt"
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

	module.Register("rabbitmq", creator)
}

// New creates RabbitMQ with default values.
func New() *RabbitMQ {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{
				UserURL:  "http://localhost:15672",
				Username: "guest",
				Password: "guest",
			},
			Client: web.Client{Timeout: web.Duration{Duration: time.Second}},
		},
	}

	return &RabbitMQ{
		Config:          config,
		charts:          charts.Copy(),
		collectedVhosts: make(map[string]bool),
	}
}

// Config is the RabbitMQ module configuration.
type Config struct {
	web.HTTP `yaml:",inline"`
}

// RabbitMQ RabbitMQ module.
type RabbitMQ struct {
	module.Base
	Config `yaml:",inline"`

	client          *client
	collectedVhosts map[string]bool
	charts          *Charts
}

// Cleanup makes cleanup.
func (RabbitMQ) Cleanup() {}

func (r RabbitMQ) createClient() (*client, error) {
	if err := r.ParseUserURL(); err != nil {
		return nil, fmt.Errorf("error on parsing url '%s' : %v", r.UserURL, err)
	}

	httpClient, err := web.NewHTTPClient(r.Client)
	if err != nil {
		return nil, fmt.Errorf("error on creating http client : %v", err)
	}

	client := newClient(httpClient, r.Request)
	return client, nil
}

// Init makes initialization.
func (r *RabbitMQ) Init() bool {
	client, err := r.createClient()
	if err != nil {
		r.Error(err)
		return false
	}

	r.client = client
	r.Debugf("using URL %s", r.URL)
	r.Debugf("using timeout: %s", r.Timeout.Duration)
	return true
}

// Check makes check.
func (r *RabbitMQ) Check() bool {
	err := r.client.findNodeName()
	if err != nil {
		r.Error(err)
		return false
	}
	return len(r.Collect()) > 0
}

// Charts creates Charts.
func (r RabbitMQ) Charts() *Charts {
	return r.charts
}

// Collect collects stats.
func (r *RabbitMQ) Collect() map[string]int64 {
	mx, err := r.collect()
	if err != nil {
		r.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

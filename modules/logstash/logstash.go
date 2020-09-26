package logstash

import (
	"errors"
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go.d.plugin/agent/module"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("logstash", creator)
}

// New creates Logstash with default values.
func New() *Logstash {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{
				URL: "http://localhost:9600",
			},
			Client: web.Client{
				Timeout: web.Duration{Duration: time.Second},
			},
		},
	}
	return &Logstash{
		Config:             config,
		charts:             charts.Copy(),
		collectedPipelines: make(map[string]bool),
	}
}

type (
	// Config is the Logstash module configuration.
	Config struct {
		web.HTTP `yaml:",inline"`
	}
	// Logstash Logstash module.
	Logstash struct {
		module.Base
		Config             `yaml:",inline"`
		client             *client
		charts             *Charts
		collectedPipelines map[string]bool
	}
)

func (l *Logstash) validateConfig() error {
	if l.URL == "" {
		return errors.New("URL not set")
	}
	return nil
}

func (l *Logstash) createClient() error {
	client, err := web.NewHTTPClient(l.Client)
	if err != nil {
		return err
	}
	l.client = newClient(client, l.Request)
	return nil
}

// Init makes initialization.
func (l *Logstash) Init() bool {
	if err := l.validateConfig(); err != nil {
		l.Errorf("error on validating config: %v", err)
		return false
	}
	if err := l.createClient(); err != nil {
		l.Errorf("error on creating client: %v", err)
		return false
	}

	l.Debugf("using URL %s", l.URL)
	l.Debugf("using timeout: %s", l.Timeout.Duration)
	return true
}

// Check makes check.
func (l *Logstash) Check() bool {
	return len(l.Collect()) > 0
}

// Charts creates Charts.
func (l *Logstash) Charts() *Charts {
	return l.charts
}

// Collect collects metrics.
func (l *Logstash) Collect() map[string]int64 {
	mx, err := l.collect()
	if err != nil {
		l.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

// Cleanup makes cleanup.
func (Logstash) Cleanup() {}

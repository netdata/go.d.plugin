package docker_engine

import (
	"time"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

const (
	defaultURL         = "http://127.0.0.1:9323/metrics"
	defaultHTTPTimeout = time.Second * 2
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("docker_engine", creator)
}

// New creates DockerEngine with default values.
func New() *DockerEngine {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{UserURL: defaultURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
	}
	return &DockerEngine{
		Config: config,
	}
}

// Config is the DockerEngine module configuration.
type Config struct {
	web.HTTP `yaml:",inline"`
}

// DockerEngine DockerEngine module.
type DockerEngine struct {
	module.Base
	Config         `yaml:",inline"`
	prom           prometheus.Prometheus
	isSwarmManager bool
}

// Cleanup makes cleanup.
func (DockerEngine) Cleanup() {}

// Init makes initialization.
func (de *DockerEngine) Init() bool {
	if err := de.ParseUserURL(); err != nil {
		de.Errorf("error on parsing url '%s' : %v", de.UserURL, err)
		return false
	}

	if de.URL.Host == "" {
		de.Error("URL is not set")
		return false
	}

	client, err := web.NewHTTPClient(de.Client)
	if err != nil {
		de.Errorf("error on creating http client : %v", err)
		return false
	}

	de.prom = prometheus.New(client, de.Request)

	return true
}

// Check makes check.
func (de *DockerEngine) Check() bool {
	return len(de.Collect()) > 0
}

// Charts creates Charts.
func (de DockerEngine) Charts() *Charts {
	if !de.isSwarmManager {
		return charts.Copy()
	}

	c := charts.Copy()
	err := c.Add(*swarmManagerCharts.Copy()...)
	if err != nil {
		de.Error(err)
		return nil
	}
	return c
}

// Collect collects metrics.
func (de *DockerEngine) Collect() map[string]int64 {
	mx, err := de.collect()

	if err != nil {
		de.Error(err)
		return nil
	}

	return mx
}

package docker_engine

import (
	"time"

	mtx "github.com/netdata/go.d.plugin/pkg/metrics"
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/stm"
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

// New creates DockerEngine with default values
func New() *DockerEngine {
	return &DockerEngine{
		Config: Config{
			HTTP: web.HTTP{
				Request: web.Request{URL: defaultURL},
				Client:  web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
			},
		},
	}
}

// Config is the DockerEngine module configuration.
type Config struct {
	web.HTTP `yaml:",inline"`
}

// DockerEngine DockerEngine module.
type DockerEngine struct {
	module.Base
	Config `yaml:",inline"`
	prom   prometheus.Prometheus
}

type metrics struct {
	Container struct {
		Actions struct {
			Changes mtx.Gauge `stm:"changes"`
			Commit  mtx.Gauge `stm:"commit"`
			Create  mtx.Gauge `stm:"create"`
			Delete  mtx.Gauge `stm:"delete"`
			Start   mtx.Gauge `stm:"start"`
		} `stm:"actions"`
		States struct {
			Paused  mtx.Gauge `stm:"paused"`
			Running mtx.Gauge `stm:"running"`
			Stopped mtx.Gauge `stm:"stopped"`
		} `stm:"states"`
	} `stm:""`
}

// Cleanup makes cleanup.
func (DockerEngine) Cleanup() {}

// Init makes initialization.
func (de *DockerEngine) Init() bool {
	if de.URL == "" {
		de.Error("URL parameter is mandatory, please set")
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
func (de DockerEngine) Check() bool {
	return len(de.Collect()) > 0
}

// Charts creates Charts
func (DockerEngine) Charts() *Charts {
	return charts.Copy()
}

// Collect collects metrics
func (de *DockerEngine) Collect() map[string]int64 {
	raw, err := de.prom.Scrape()

	if err != nil {
		de.Error(err)
		return nil
	}

	var metrics metrics

	gatherContainerActions(raw, &metrics)
	gatherContainerStates(raw, &metrics)

	return stm.ToMap(metrics)
}

func gatherContainerActions(raw prometheus.Metrics, ms *metrics) {
	for _, metric := range raw.FindByName("engine_daemon_container_actions_seconds_count") {
		action := metric.Labels.Get("action")
		if action == "" {
			continue
		}
		value := metric.Value
		switch action {
		case "changes":
			ms.Container.Actions.Changes.Set(value)
		case "commit":
			ms.Container.Actions.Commit.Set(value)
		case "create":
			ms.Container.Actions.Create.Set(value)
		case "delete":
			ms.Container.Actions.Delete.Set(value)
		case "start":
			ms.Container.Actions.Start.Set(value)
		}
	}
}

func gatherContainerStates(raw prometheus.Metrics, ms *metrics) {
	for _, metric := range raw.FindByName("engine_daemon_container_states_containers") {
		action := metric.Labels.Get("state")
		if action == "" {
			continue
		}
		value := metric.Value
		switch action {
		case "paused":
			ms.Container.States.Paused.Set(value)
		case "running":
			ms.Container.States.Running.Set(value)
		case "stopped":
			ms.Container.States.Stopped.Set(value)
		}
	}
}

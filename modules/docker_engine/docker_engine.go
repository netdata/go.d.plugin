package docker_engine

import (
	"time"

	mtx "github.com/netdata/go.d.plugin/pkg/metrics"
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
	"github.com/prometheus/prometheus/pkg/labels"
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
			Changes mtx.Counter `stm:"changes"`
			Commit  mtx.Counter `stm:"commit"`
			Create  mtx.Counter `stm:"create"`
			Delete  mtx.Counter `stm:"delete"`
			Start   mtx.Counter `stm:"start"`
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

	ms := raw.FindByName("engine_daemon_container_actions_seconds_count")
	metrics.Container.Actions.Changes.Add(
		ms.Match(newEqualMatcher("action", "changes")).Max())
	metrics.Container.Actions.Commit.Add(
		ms.Match(newEqualMatcher("action", "commit")).Max())
	metrics.Container.Actions.Create.Add(
		ms.Match(newEqualMatcher("action", "create")).Max())
	metrics.Container.Actions.Delete.Add(
		ms.Match(newEqualMatcher("action", "delete")).Max())
	metrics.Container.Actions.Start.Add(
		ms.Match(newEqualMatcher("action", "start")).Max())

	ms = raw.FindByName("engine_daemon_container_states_containers")
	metrics.Container.States.Paused.Set(
		ms.Match(newEqualMatcher("state", "paused")).Max())
	metrics.Container.States.Running.Set(
		ms.Match(newEqualMatcher("state", "running")).Max())
	metrics.Container.States.Stopped.Set(
		ms.Match(newEqualMatcher("state", "stopped")).Max())

	return stm.ToMap(metrics)
}

func newEqualMatcher(name, value string) *labels.Matcher {
	return &labels.Matcher{Type: labels.MatchEqual, Name: name, Value: value}
}

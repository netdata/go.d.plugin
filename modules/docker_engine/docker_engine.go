package docker_engine

import (
	"time"

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

// New creates DockerEngine with default values.
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
			Changes int `stm:"changes"`
			Commit  int `stm:"commit"`
			Create  int `stm:"create"`
			Delete  int `stm:"delete"`
			Start   int `stm:"start"`
		} `stm:"actions"`
		States struct {
			Paused  int `stm:"paused"`
			Running int `stm:"running"`
			Stopped int `stm:"stopped"`
		} `stm:"states"`
	} `stm:"container"`
	Builder struct {
		FailsByReason struct {
			BuildCanceled                int `stm:"build_canceled"`
			BuildTargetNotReachableError int `stm:"build_target_not_reachable_error"`
			CommandNotSupportedError     int `stm:"command_not_supported_error"`
			DockerfileEmptyError         int `stm:"dockerfile_empty_error"`
			DockerfileSyntaxError        int `stm:"dockerfile_syntax_error"`
			ErrorProcessingCommandsError int `stm:"error_processing_commands_error"`
			MissingOnbuildArgumentsError int `stm:"missing_onbuild_arguments_error"`
			UnknownInstructionError      int `stm:"unknown_instruction_error"`
		} `stm:"fails"`
	} `stm:"builder"`
	HealthChecks struct {
		Failed int `stm:"failed"`
	} `stm:"health_checks"`
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

// Charts creates Charts.
func (DockerEngine) Charts() *Charts {
	return charts.Copy()
}

// Collect collects metrics.
func (de *DockerEngine) Collect() map[string]int64 {
	raw, err := de.prom.Scrape()

	if err != nil {
		de.Error(err)
		return nil
	}

	var mx metrics

	mx.HealthChecks.Failed = int(raw.FindByName("engine_daemon_health_checks_failed_total").Max())
	collectContainerActions(raw, &mx)
	collectContainerStates(raw, &mx)
	collectBuilderBuildsFails(raw, &mx)

	return stm.ToMap(mx)
}

func collectContainerActions(raw prometheus.Metrics, ms *metrics) {
	for _, metric := range raw.FindByName("engine_daemon_container_actions_seconds_count") {
		action := metric.Labels.Get("action")
		value := metric.Value
		switch action {
		case "changes":
			ms.Container.Actions.Changes = int(value)
		case "commit":
			ms.Container.Actions.Commit = int(value)
		case "create":
			ms.Container.Actions.Create = int(value)
		case "delete":
			ms.Container.Actions.Delete = int(value)
		case "start":
			ms.Container.Actions.Start = int(value)
		}
	}
}

func collectContainerStates(raw prometheus.Metrics, ms *metrics) {
	for _, metric := range raw.FindByName("engine_daemon_container_states_containers") {
		action := metric.Labels.Get("state")
		value := metric.Value
		switch action {
		case "paused":
			ms.Container.States.Paused = int(value)
		case "running":
			ms.Container.States.Running = int(value)
		case "stopped":
			ms.Container.States.Stopped = int(value)
		}
	}
}

func collectBuilderBuildsFails(raw prometheus.Metrics, ms *metrics) {
	for _, metric := range raw.FindByName("builder_builds_failed_total") {
		action := metric.Labels.Get("reason")
		value := metric.Value
		switch action {
		case "build_canceled":
			ms.Builder.FailsByReason.BuildCanceled = int(value)
		case "build_target_not_reachable_error":
			ms.Builder.FailsByReason.BuildTargetNotReachableError = int(value)
		case "command_not_supported_error":
			ms.Builder.FailsByReason.CommandNotSupportedError = int(value)
		case "dockerfile_empty_error":
			ms.Builder.FailsByReason.DockerfileEmptyError = int(value)
		case "dockerfile_syntax_error":
			ms.Builder.FailsByReason.DockerfileSyntaxError = int(value)
		case "error_processing_commands_error":
			ms.Builder.FailsByReason.ErrorProcessingCommandsError = int(value)
		case "missing_onbuild_arguments_error":
			ms.Builder.FailsByReason.MissingOnbuildArgumentsError = int(value)
		case "unknown_instruction_error":
			ms.Builder.FailsByReason.UnknownInstructionError = int(value)
		}
	}
}

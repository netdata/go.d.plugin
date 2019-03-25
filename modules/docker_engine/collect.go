package docker_engine

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/stm"
)

func (de *DockerEngine) collect() (map[string]int64, error) {
	raw, err := de.prom.Scrape()

	if err != nil {
		return nil, err
	}

	var mx metrics

	collectHealthChecks(raw, &mx)
	collectContainerActions(raw, &mx)
	collectContainerStates(raw, &mx)
	collectBuilderBuildsFails(raw, &mx)

	return stm.ToMap(mx), nil

}

func collectHealthChecks(raw prometheus.Metrics, mx *metrics) {
	m := raw.FindByName("engine_daemon_health_checks_failed_total")
	mx.HealthChecks.Failed.Set(m.Max())
}

func collectContainerActions(raw prometheus.Metrics, mx *metrics) {
	for _, metric := range raw.FindByName("engine_daemon_container_actions_seconds_count") {
		action := metric.Labels.Get("action")
		switch action {
		case "changes":
			mx.Container.Actions.Changes.Set(metric.Value)
		case "commit":
			mx.Container.Actions.Commit.Set(metric.Value)
		case "create":
			mx.Container.Actions.Create.Set(metric.Value)
		case "delete":
			mx.Container.Actions.Delete.Set(metric.Value)
		case "start":
			mx.Container.Actions.Start.Set(metric.Value)
		}
	}
}

func collectContainerStates(raw prometheus.Metrics, mx *metrics) {
	for _, metric := range raw.FindByName("engine_daemon_container_states_containers") {
		action := metric.Labels.Get("state")
		switch action {
		case "paused":
			mx.Container.States.Paused.Set(metric.Value)
		case "running":
			mx.Container.States.Running.Set(metric.Value)
		case "stopped":
			mx.Container.States.Stopped.Set(metric.Value)
		}
	}
}

func collectBuilderBuildsFails(raw prometheus.Metrics, ms *metrics) {
	for _, metric := range raw.FindByName("builder_builds_failed_total") {
		action := metric.Labels.Get("reason")
		switch action {
		case "build_canceled":
			ms.Builder.FailsByReason.BuildCanceled.Set(metric.Value)
		case "build_target_not_reachable_error":
			ms.Builder.FailsByReason.BuildTargetNotReachableError.Set(metric.Value)
		case "command_not_supported_error":
			ms.Builder.FailsByReason.CommandNotSupportedError.Set(metric.Value)
		case "dockerfile_empty_error":
			ms.Builder.FailsByReason.DockerfileEmptyError.Set(metric.Value)
		case "dockerfile_syntax_error":
			ms.Builder.FailsByReason.DockerfileSyntaxError.Set(metric.Value)
		case "error_processing_commands_error":
			ms.Builder.FailsByReason.ErrorProcessingCommandsError.Set(metric.Value)
		case "missing_onbuild_arguments_error":
			ms.Builder.FailsByReason.MissingOnbuildArgumentsError.Set(metric.Value)
		case "unknown_instruction_error":
			ms.Builder.FailsByReason.UnknownInstructionError.Set(metric.Value)
		}
	}
}

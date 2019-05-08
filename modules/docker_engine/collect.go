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

	collectHealthChecks(&mx, raw)
	collectContainerActions(&mx, raw)
	collectContainerStates(&mx, raw)
	collectBuilderBuildsFails(&mx, raw)

	if isSwarmManager(raw) {
		de.isSwarmManager = true
		mx.SwarmManager = &swarmManager{}
		collectSwarmManager(&mx, raw)
	}

	return stm.ToMap(mx), nil

}

func collectHealthChecks(mx *metrics, raw prometheus.Metrics) {
	mx.HealthChecks.Failed.Set(raw.FindByName("engine_daemon_health_checks_failed_total").Max())
}

func collectContainerActions(mx *metrics, raw prometheus.Metrics) {
	for _, metric := range raw.FindByName("engine_daemon_container_actions_seconds_count") {
		action := metric.Labels.Get("action")
		if action == "" {
			continue
		}
		v := metric.Value

		switch action {
		default:
		case "changes":
			mx.Container.Actions.Changes.Set(v)
		case "commit":
			mx.Container.Actions.Commit.Set(v)
		case "create":
			mx.Container.Actions.Create.Set(v)
		case "delete":
			mx.Container.Actions.Delete.Set(v)
		case "start":
			mx.Container.Actions.Start.Set(v)
		}
	}
}

func collectContainerStates(mx *metrics, raw prometheus.Metrics) {
	for _, metric := range raw.FindByName("engine_daemon_container_states_containers") {
		state := metric.Labels.Get("state")
		if state == "" {
			continue
		}
		v := metric.Value

		switch state {
		default:
		case "paused":
			mx.Container.States.Paused.Set(v)
		case "running":
			mx.Container.States.Running.Set(v)
		case "stopped":
			mx.Container.States.Stopped.Set(v)
		}
	}
}

func collectBuilderBuildsFails(mx *metrics, raw prometheus.Metrics) {
	for _, metric := range raw.FindByName("builder_builds_failed_total") {
		reason := metric.Labels.Get("reason")
		if reason == "" {
			continue
		}
		v := metric.Value

		switch reason {
		default:
		case "build_canceled":
			mx.Builder.FailsByReason.BuildCanceled.Set(v)
		case "build_target_not_reachable_error":
			mx.Builder.FailsByReason.BuildTargetNotReachableError.Set(v)
		case "command_not_supported_error":
			mx.Builder.FailsByReason.CommandNotSupportedError.Set(v)
		case "dockerfile_empty_error":
			mx.Builder.FailsByReason.DockerfileEmptyError.Set(v)
		case "dockerfile_syntax_error":
			mx.Builder.FailsByReason.DockerfileSyntaxError.Set(v)
		case "error_processing_commands_error":
			mx.Builder.FailsByReason.ErrorProcessingCommandsError.Set(v)
		case "missing_onbuild_arguments_error":
			mx.Builder.FailsByReason.MissingOnbuildArgumentsError.Set(v)
		case "unknown_instruction_error":
			mx.Builder.FailsByReason.UnknownInstructionError.Set(v)
		}
	}
}

func isSwarmManager(raw prometheus.Metrics) bool {
	return raw.FindByName("swarm_node_manager").Max() == 1
}

func collectSwarmManager(mx *metrics, raw prometheus.Metrics) {
	v := raw.FindByName("swarm_manager_configs_total").Max()
	mx.SwarmManager.Configs.Set(v)

	v = raw.FindByName("swarm_manager_networks_total").Max()
	mx.SwarmManager.Networks.Set(v)

	v = raw.FindByName("swarm_manager_secrets_total").Max()
	mx.SwarmManager.Secrets.Set(v)

	v = raw.FindByName("swarm_manager_services_total").Max()
	mx.SwarmManager.Services.Set(v)

	v = raw.FindByName("swarm_manager_leader").Max()
	mx.SwarmManager.IsLeader.Set(v)

	for _, metric := range raw.FindByName("swarm_manager_nodes") {
		state := metric.Labels.Get("state")
		if state == "" {
			continue
		}
		v := metric.Value

		switch state {
		default:
		case "disconnected":
			mx.SwarmManager.Nodes.PerState.Disconnected.Set(v)
		case "down":
			mx.SwarmManager.Nodes.PerState.Down.Set(v)
		case "ready":
			mx.SwarmManager.Nodes.PerState.Ready.Set(v)
		case "unknown":
			mx.SwarmManager.Nodes.PerState.Unknown.Set(v)
		}
		mx.SwarmManager.Nodes.Total.Add(v)
	}

	for _, metric := range raw.FindByName("swarm_manager_tasks_total") {
		state := metric.Labels.Get("state")
		if state == "" {
			continue
		}
		v := metric.Value

		switch state {
		default:
		case "accepted":
			mx.SwarmManager.Tasks.PerState.Accepted.Set(v)
		case "assigned":
			mx.SwarmManager.Tasks.PerState.Assigned.Set(v)
		case "complete":
			mx.SwarmManager.Tasks.PerState.Complete.Set(v)
		case "failed":
			mx.SwarmManager.Tasks.PerState.Failed.Set(v)
		case "new":
			mx.SwarmManager.Tasks.PerState.New.Set(v)
		case "orphaned":
			mx.SwarmManager.Tasks.PerState.Orphaned.Set(v)
		case "pending":
			mx.SwarmManager.Tasks.PerState.Pending.Set(v)
		case "preparing":
			mx.SwarmManager.Tasks.PerState.Preparing.Set(v)
		case "ready":
			mx.SwarmManager.Tasks.PerState.Ready.Set(v)
		case "rejected":
			mx.SwarmManager.Tasks.PerState.Rejected.Set(v)
		case "remove":
			mx.SwarmManager.Tasks.PerState.Remove.Set(v)
		case "running":
			mx.SwarmManager.Tasks.PerState.Running.Set(v)
		case "shutdown":
			mx.SwarmManager.Tasks.PerState.Shutdown.Set(v)
		case "starting":
			mx.SwarmManager.Tasks.PerState.Starting.Set(v)
		}
		mx.SwarmManager.Tasks.Total.Add(v)
	}
}

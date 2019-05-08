package docker_engine

import (
	mtx "github.com/netdata/go.d.plugin/pkg/metrics"
)

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
	} `stm:"container"`
	Builder struct {
		FailsByReason struct {
			BuildCanceled                mtx.Gauge `stm:"build_canceled"`
			BuildTargetNotReachableError mtx.Gauge `stm:"build_target_not_reachable_error"`
			CommandNotSupportedError     mtx.Gauge `stm:"command_not_supported_error"`
			DockerfileEmptyError         mtx.Gauge `stm:"dockerfile_empty_error"`
			DockerfileSyntaxError        mtx.Gauge `stm:"dockerfile_syntax_error"`
			ErrorProcessingCommandsError mtx.Gauge `stm:"error_processing_commands_error"`
			MissingOnbuildArgumentsError mtx.Gauge `stm:"missing_onbuild_arguments_error"`
			UnknownInstructionError      mtx.Gauge `stm:"unknown_instruction_error"`
		} `stm:"fails"`
	} `stm:"builder"`
	HealthChecks struct {
		Failed mtx.Gauge `stm:"failed"`
	} `stm:"health_checks"`

	SwarmManager *swarmManager `stm:"swarm_manager"`
}

type swarmManager struct {
	IsLeader mtx.Gauge `stm:"leader"`
	Configs  mtx.Gauge `stm:"configs_total"`
	Networks mtx.Gauge `stm:"networks_total"`
	Secrets  mtx.Gauge `stm:"secrets_total"`
	Services mtx.Gauge `stm:"services_total"`
	Nodes    struct {
		Total    mtx.Gauge `stm:"total"`
		PerState struct {
			Disconnected mtx.Gauge `stm:"disconnected"`
			Down         mtx.Gauge `stm:"down"`
			Ready        mtx.Gauge `stm:"ready"`
			Unknown      mtx.Gauge `stm:"unknown"`
		} `stm:"state"`
	} `stm:"nodes"`
	Tasks struct {
		Total    mtx.Gauge `stm:"total"`
		PerState struct {
			Accepted  mtx.Gauge `stm:"accepted"`
			Assigned  mtx.Gauge `stm:"assigned"`
			Complete  mtx.Gauge `stm:"complete"`
			Failed    mtx.Gauge `stm:"failed"`
			New       mtx.Gauge `stm:"new"`
			Orphaned  mtx.Gauge `stm:"orphaned"`
			Pending   mtx.Gauge `stm:"pending"`
			Preparing mtx.Gauge `stm:"preparing"`
			Ready     mtx.Gauge `stm:"ready"`
			Rejected  mtx.Gauge `stm:"rejected"`
			Remove    mtx.Gauge `stm:"remove"`
			Running   mtx.Gauge `stm:"running"`
			Shutdown  mtx.Gauge `stm:"shutdown"`
			Starting  mtx.Gauge `stm:"starting"`
		} `stm:"state"`
	} `stm:"tasks"`
}

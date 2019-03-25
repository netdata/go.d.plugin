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
}

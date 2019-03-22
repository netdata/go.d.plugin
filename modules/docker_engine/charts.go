package docker_engine

import "github.com/netdata/go-orchestrator/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

var charts = Charts{
	{
		ID:    "engine_daemon_container_actions",
		Title: "Container Actions",
		Units: "actions/s",
		Fam:   "containers",
		Dims: Dims{
			{ID: "container_actions_changes", Name: "changes", Algo: module.Incremental},
			{ID: "container_actions_commit", Name: "commit", Algo: module.Incremental},
			{ID: "container_actions_create", Name: "create", Algo: module.Incremental},
			{ID: "container_actions_delete", Name: "delete", Algo: module.Incremental},
			{ID: "container_actions_start", Name: "start", Algo: module.Incremental},
		},
	},
	{
		ID:    "engine_daemon_container_states_containers",
		Title: "Containers In Various States",
		Units: "count",
		Fam:   "containers",
		Dims: Dims{
			{ID: "container_states_running", Name: "running"},
			{ID: "container_states_paused", Name: "paused"},
			{ID: "container_states_stopped", Name: "stopped"},
		},
	},
	{
		ID:    "builder_builds_failed_total",
		Title: "Builder Builds Fails By Reason",
		Units: "fails/s",
		Fam:   "builder",
		Dims: Dims{
			{ID: "builder_fails_build_canceled", Name: "build_canceled", Algo: module.Incremental},
			{ID: "builder_fails_build_target_not_reachable_error", Name: "build_target_not_reachable_error", Algo: module.Incremental},
			{ID: "builder_fails_command_not_supported_error", Name: "command_not_supported_error", Algo: module.Incremental},
			{ID: "builder_fails_dockerfile_empty_error", Name: "dockerfile_empty_error", Algo: module.Incremental},
			{ID: "builder_fails_dockerfile_syntax_error", Name: "dockerfile_syntax_error", Algo: module.Incremental},
			{ID: "builder_fails_error_processing_commands_error", Name: "error_processing_commands_error", Algo: module.Incremental},
			{ID: "builder_fails_missing_onbuild_arguments_error", Name: "missing_onbuild_arguments_error", Algo: module.Incremental},
			{ID: "builder_fails_unknown_instruction_error", Name: "unknown_instruction_error", Algo: module.Incremental},
		},
	},
	{
		ID:    "engine_daemon_health_checks_failed_total",
		Title: "Health Checks",
		Units: "events/s",
		Fam:   "health checks",
		Dims: Dims{
			{ID: "health_checks_failed", Name: "fails", Algo: module.Incremental},
		},
	},
}

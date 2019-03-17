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
		Title: "Container Actions", Units: "actions/s", Fam: "container actions",
		Dims: Dims{
			{ID: "actions_changes", Name: "changes", Algo: module.Incremental},
			{ID: "actions_commit", Name: "commit", Algo: module.Incremental},
			{ID: "actions_create", Name: "create", Algo: module.Incremental},
			{ID: "actions_delete", Name: "delete", Algo: module.Incremental},
			{ID: "actions_start", Name: "start", Algo: module.Incremental},
		},
	},
	{
		ID:    "engine_daemon_container_states_containers",
		Title: "The Count Of Containers In Various States", Units: "count", Fam: "container states",
		Dims: Dims{
			{ID: "states_running", Name: "running"},
			{ID: "states_paused", Name: "paused"},
			{ID: "states_stopped", Name: "stopped"},
		},
	},
}

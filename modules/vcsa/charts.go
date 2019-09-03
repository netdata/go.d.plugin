package vcsa

import "github.com/netdata/go-orchestrator/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

var charts = Charts{
	{
		ID:    "overall_health",
		Title: "Overall System Health",
		Units: "status",
		Fam:   "health",
		Ctx:   "vcsa.overall_health",
		Dims: Dims{
			{ID: "system"},
		},
	},
	{
		ID:    "components_health",
		Title: "System Key Components Health",
		Units: "status",
		Fam:   "health",
		Ctx:   "vcsa.component_health",
		Dims: Dims{
			{ID: "appl_mgmt"},
			{ID: "database_storage"},
			{ID: "load"},
			{ID: "mem"},
			{ID: "software_packages"},
			{ID: "storage"},
			{ID: "swap"},
		},
	},
}

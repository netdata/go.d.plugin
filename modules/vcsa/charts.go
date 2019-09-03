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
		ID:    "vcenter_health",
		Title: "Health Status",
		Units: "status",
		Fam:   "health",
		Ctx:   "vcenter.health",
		Dims: Dims{
			{ID: "appl_mgmt"},
			{ID: "database_storage"},
			{ID: "load"},
			{ID: "mem"},
			{ID: "software_packages"},
			{ID: "storage"},
			{ID: "swap"},
			{ID: "system"},
		},
	},
}

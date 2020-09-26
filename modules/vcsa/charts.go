package vcsa

import "github.com/netdata/go.d.plugin/agent/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

var charts = Charts{
	{
		ID:    "system_health",
		Title: "Overall System Health",
		Units: "status",
		Fam:   "health",
		Ctx:   "vcsa.system_health",
		Dims: Dims{
			{ID: "system"},
		},
	},
	{
		ID:    "components_health",
		Title: "Components Health",
		Units: "status",
		Fam:   "health",
		Ctx:   "vcsa.components_health",
		Dims: Dims{
			{ID: "applmgmt"},
			{ID: "database_storage"},
			{ID: "load"},
			{ID: "mem"},
			{ID: "storage"},
			{ID: "swap"},
		},
	},
	{
		ID:    "software_updates_health",
		Title: "Software Updates Health",
		Units: "status",
		Fam:   "health",
		Ctx:   "vcsa.software_updates_health",
		Dims: Dims{
			{ID: "software_packages"},
		},
	},
}

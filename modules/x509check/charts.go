package x509check

import "github.com/netdata/go-orchestrator/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
	// Vars is an alias for module.Vars
	Vars = module.Vars
)

var charts = Charts{
	{
		ID:    "time_until_expiration",
		Title: "Time Until Certificate Expiration",
		Units: "seconds",
		Fam:   "expiration time",
		Ctx:   "x509check.time_until_expiration",
		Dims: Dims{
			{ID: "expiry"},
		},
		Vars: Vars{
			{ID: "days_until_expiration_warning"},
			{ID: "days_until_expiration_critical"},
		},
	},
}

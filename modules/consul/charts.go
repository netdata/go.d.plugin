package consul

import (
	"github.com/netdata/go.d.plugin/agent/module"
)

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Chart is an alias for module.Chart
	Chart = module.Chart
	// Dim is an alias for module.Dim
	Dim = module.Dim
)

var charts = Charts{
	{
		ID:    "service_checks",
		Title: "Service Checks",
		Fam:   "checks",
		Units: "status",
		Ctx:   "consul.checks",
	},
	{
		ID:    "unbound_checks",
		Title: "Unbound Checks",
		Fam:   "checks",
		Units: "status",
		Ctx:   "consul.checks",
	},
}

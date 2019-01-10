package consul

import (
	"github.com/netdata/go.d.plugin/modules"
)

type (
	// Charts is an alias for modules.Charts
	Charts = modules.Charts
	// Chart is an alias for modules.Chart
	Chart = modules.Chart
	// Dim is an alias for modules.Dim
	Dim = modules.Dim
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

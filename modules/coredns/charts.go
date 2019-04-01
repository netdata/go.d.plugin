package coredns

import "github.com/netdata/go-orchestrator/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
	// Dim is an alias for module.Dim
	Dim = module.Dim
)

var charts = Charts{
	{
		ID:    "panic_count_total",
		Title: "The Total Number Of Panics",
		Units: "events/s",
		Fam:   "panics",
		Ctx:   "coredns.panic_count_total",
		Dims: Dims{
			{ID: "panic_count_total", Name: "panics", Algo: module.Incremental},
		},
	},
	{
		ID:    "request_count_total",
		Title: "The Total Number Of Requests Per All Zones, Protocols And Families",
		Units: "events/s",
		Fam:   "requests",
		Ctx:   "coredns.request_count_total",
		Dims: Dims{
			{ID: "request_count_total", Name: "requests", Algo: module.Incremental},
		},
	},
}

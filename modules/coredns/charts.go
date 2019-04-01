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
		ID:    "panic_count",
		Title: "The Number Of Panics",
		Units: "events/s",
		Fam:   "panic",
		Ctx:   "coredns.panic_count",
		Dims: Dims{
			{
				ID: "panic_count", Name: "panic", Algo: module.Incremental,
			},
		},
	},
}

package example

import "github.com/netdata/go-orchestrator/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

var charts = Charts{
	{
		ID:    "random",
		Title: "A Random Number", Units: "random", Fam: "random",
		Dims: Dims{
			{ID: "random0", Name: "random 0"},
			{ID: "random1", Name: "random 1"},
		},
	},
}

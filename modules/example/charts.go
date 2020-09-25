package example

import "github.com/netdata/go.d.plugin/plugin/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

var charts = Charts{
	{
		ID:    "random",
		Title: "A Random Number",
		Units: "random",
		Fam:   "random",
		Ctx:   "example.random",
		Dims: Dims{
			{ID: "random0", Name: "random0"},
			{ID: "random1", Name: "random1"},
		},
	},
}

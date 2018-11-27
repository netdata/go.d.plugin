package example

import "github.com/netdata/go.d.plugin/modules"

type (
	Charts  = modules.Charts
	Options = modules.Opts
	Dims    = modules.Dims
)

var charts = Charts{
	{
		ID:    "chart1",
		Title: "A Random Number", Units: "random", Fam: "random",
		Dims: Dims{
			{ID: "random0", Name: "random 0"},
			{ID: "random1", Name: "random 1"},
		},
	},
}

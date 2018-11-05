package example

import "github.com/l2isbad/go.d.plugin/modules"

type (
	Charts  = modules.Charts
	Options = modules.Opts
	Dims    = modules.Dims
)

var charts = Charts{
	{
		ID:    "chart1",
		Title: "qwe", Units: "qw", Fam: "random",
		Dims: Dims{
			{ID: "random0", Name: "random"},
		},
	},
	{
		ID:    "chart2",
		Title: "qwe", Units: "qw", Fam: "random", Type: modules.Area,
		Dims: Dims{
			{ID: "random1", Name: "random"},
		},
	},
}

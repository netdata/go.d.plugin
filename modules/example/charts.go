package example

import "github.com/l2isbad/go.d.plugin/modules"

type (
	Charts  = modules.Charts
	Options = modules.Opts
	Dims    = modules.Dims
)

var uCharts = Charts{
	{
		ID:   "chart1",
		Opts: Options{Title: "Random Data 1", Units: "random", Fam: "random"},
		Dims: Dims{
			{ID: "random0", Name: "random"},
		},
	},
	{
		ID:   "chart2",
		Opts: Options{Title: "Random Data 2", Units: "random", Fam: "random", Type: modules.Area},
		Dims: Dims{
			{ID: "random1", Name: "random"},
		},
	},
}

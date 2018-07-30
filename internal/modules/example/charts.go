package example

import "github.com/l2isbad/go.d.plugin/internal/pkg/charts"

type (
	Charts  = charts.Charts
	Options = charts.Opts
	Dims    = charts.Dims
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
		Opts: Options{Title: "Random Data 2", Units: "random", Fam: "random", Type: charts.Area},
		Dims: Dims{
			{ID: "random1", Name: "random"},
		},
	},
}

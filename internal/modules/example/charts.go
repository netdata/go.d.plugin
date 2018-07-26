package example

import "github.com/l2isbad/go.d.plugin/internal/pkg/charts"

type (
	Charts     = charts.Charts
	Options    = charts.Opts
	Dimensions = charts.Dims
)

var uCharts = Charts{
	{
		ID:   "chart1",
		Opts: Options{Title: "Random Data 1", Units: "random", Family: "random", Type: charts.Line},
		Dims: Dimensions{
			{ID: "random0", Name: "random"},
		},
	},
	{
		ID:   "chart2",
		Opts: Options{Title: "Random Data 2", Units: "random", Family: "random", Type: charts.Area},
		Dims: Dimensions{
			{ID: "random1", Name: "random"},
		},
	},
}

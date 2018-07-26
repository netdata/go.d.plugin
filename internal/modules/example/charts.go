package example

import "github.com/l2isbad/go.d.plugin/internal/pkg/charts"

type (
	Charts     = charts.Charts
	Options    = charts.Options
	Dimensions = charts.Dimensions
)

var uCharts = Charts{
	{
		ID:      "chart1",
		Options: Options{Title: "Random Data 1", Units: "random", Family: "random", Type: charts.Line},
		Dimensions: Dimensions{
			{ID: "random0", Name: "random"},
		},
	},
	{
		ID:      "chart2",
		Options: Options{Title: "Random Data 2", Units: "random", Family: "random", Type: charts.Area},
		Dimensions: Dimensions{
			{ID: "random1", Name: "random"},
		},
	},
}

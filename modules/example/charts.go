package example

import (
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

var chartTemplate = module.Chart{
	ID:    "random_%d",
	Title: "A Random Number",
	Units: "random",
	Fam:   "random",
	Ctx:   "example.random",
}

var hiddenChartTemplate = module.Chart{
	ID:    "hidden_random_%d",
	Title: "A Random Number",
	Units: "random",
	Fam:   "random",
	Ctx:   "example.random",
	Opts: module.Opts{
		Hidden: true,
	},
}

func newChart(num int, typ module.ChartType) *module.Chart {
	chart := chartTemplate.Copy()
	chart.ID = fmt.Sprintf(chart.ID, num)
	chart.Type = typ
	return chart
}

func newHiddenChart(num int, typ module.ChartType) *module.Chart {
	chart := hiddenChartTemplate.Copy()
	chart.ID = fmt.Sprintf(chart.ID, num)
	chart.Type = typ
	return chart
}

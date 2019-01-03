package consul

import (
	"fmt"

	"github.com/netdata/go.d.plugin/modules"
)

type (
	// Charts is an alias for modules.Charts
	Charts = modules.Charts
	// Chart is an alias for modules.Chart
	Chart = modules.Chart
	// Dim is an alias for modules.Dim
	Dim = modules.Dim
)

var (
	boundCheckChart = Chart{
		ID:    "%s_check",
		Title: "Service %s[%s] Check %s[%s] Status",
		Fam:   "service checks",
		Units: "status",
		Ctx:   "consul.checks",
	}
	unboundCheckChart = Chart{
		ID:    "%s_check",
		Title: "Check %s[%s] Status",
		Fam:   "unbound checks",
		Units: "status",
		Ctx:   "consul.checks",
	}
)

func createCheckChart(check *agentCheck) (chart *Chart) {
	if check.ServiceID != "" {
		chart = boundCheckChart.Copy()
		chart.ID = fmt.Sprintf(chart.ID, check.CheckID)
		chart.Title = fmt.Sprintf(chart.Title, check.ServiceID, check.ServiceName, check.CheckID, check.Name)
	} else {
		chart = unboundCheckChart.Copy()
		chart.ID = fmt.Sprintf(chart.ID, check.CheckID)
		chart.Title = fmt.Sprintf(chart.Title, check.CheckID, check.Name)
	}
	_ = chart.AddDim(&Dim{ID: check.CheckID})
	return chart
}

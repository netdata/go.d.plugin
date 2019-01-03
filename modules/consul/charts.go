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
	// Dims is an alias for modules.Dims
	Dims = modules.Dims
)

var (
	boundCheckChart = Chart{
		ID:    "check_%s",
		Title: "Service %s[%s] Check %s[%s] Status",
		Fam:   "service checks",
		Units: "status",
		Ctx:   "consul.checks",
		Dims: Dims{
			{ID: "%s_passing", Name: "passing"},
			{ID: "%s_critical", Name: "critical"},
			{ID: "%s_maintenance", Name: "maintenance"},
			{ID: "%s_warning", Name: "warning"},
		},
	}
	unboundCheckChart = Chart{
		ID:    "check_%s",
		Title: "Check %s[%s] Status",
		Fam:   "unbound checks",
		Units: "status",
		Ctx:   "consul.checks",
		Dims: Dims{
			{ID: "%s_passing", Name: "passing"},
			{ID: "%s_critical", Name: "critical"},
			{ID: "%s_maintenance", Name: "maintenance"},
			{ID: "%s_warning", Name: "warning"},
		},
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

	for _, dim := range chart.Dims {
		dim.ID = fmt.Sprintf(dim.ID, check.CheckID)
	}
	return chart
}

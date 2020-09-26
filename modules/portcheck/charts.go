package portcheck

import (
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

var portCharts = Charts{
	{
		ID:    "port_%d_status",
		Title: "TCP Check Status",
		Units: "boolean",
		Fam:   "port %d",
		Ctx:   "portcheck.status",
		Dims: Dims{
			{ID: "port_%d_success", Name: "success"},
			{ID: "port_%d_failed", Name: "failed"},
			{ID: "port_%d_timeout", Name: "timeout"},
		},
	},
	{
		ID:    "port_%d_current_state_duration",
		Title: "Current State Duration",
		Units: "seconds",
		Fam:   "port %d",
		Ctx:   "portcheck.state_duration",
		Dims: Dims{
			{ID: "port_%d_current_state_duration", Name: "time"},
		},
	},
	{
		ID:    "port_%d_connection_latency",
		Title: "TCP Connection Latency",
		Units: "ms",
		Fam:   "port %d",
		Ctx:   "portcheck.latency",
		Dims: Dims{
			{ID: "port_%d_latency", Name: "time"},
		},
	},
}

func newPortCharts(port int) *Charts {
	cs := portCharts.Copy()
	for _, chart := range *cs {
		chart.ID = fmt.Sprintf(chart.ID, port)
		chart.Fam = fmt.Sprintf(chart.Fam, port)

		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, port)
		}
	}
	return cs
}

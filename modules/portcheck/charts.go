// SPDX-License-Identifier: GPL-3.0-or-later

package portcheck

import (
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

var chartsTmpl = module.Charts{
	checkStatusChartTmpl.Copy(),
	checkInStateDurationChartTmpl.Copy(),
	checkConnectionLatencyChartTmpl.Copy(),
}

var checkStatusChartTmpl = module.Chart{
	ID:    "port_%d_status",
	Title: "TCP Check Status",
	Units: "boolean",
	Fam:   "port %d",
	Ctx:   "portcheck.status",
	Dims: module.Dims{
		{ID: "port_%d_success", Name: "success"},
		{ID: "port_%d_failed", Name: "failed"},
		{ID: "port_%d_timeout", Name: "timeout"},
	},
}

var checkInStateDurationChartTmpl = module.Chart{
	ID:    "port_%d_current_state_duration",
	Title: "Current State Duration",
	Units: "seconds",
	Fam:   "port %d",
	Ctx:   "portcheck.state_duration",
	Dims: module.Dims{
		{ID: "port_%d_current_state_duration", Name: "time"},
	},
}

var checkConnectionLatencyChartTmpl = module.Chart{
	ID:    "port_%d_connection_latency",
	Title: "TCP Connection Latency",
	Units: "ms",
	Fam:   "port %d",
	Ctx:   "portcheck.latency",
	Dims: module.Dims{
		{ID: "port_%d_latency", Name: "time"},
	},
}

func newPortCharts(port int) *module.Charts {
	charts := chartsTmpl.Copy()
	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, port)
		chart.Fam = fmt.Sprintf(chart.Fam, port)

		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, port)
		}
	}
	return charts
}

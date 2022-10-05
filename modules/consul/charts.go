// SPDX-License-Identifier: GPL-3.0-or-later

package consul

import (
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

const (
	prioServiceHealthCheckStatus = module.Priority + iota
	prioUnboundHealthCheckStatus
)

var (
	chartTmplServiceHealthCheckStatus = module.Chart{
		ID:       "health_check_%s_status",
		Title:    "Service health check status",
		Units:    "status",
		Fam:      "health checks",
		Ctx:      "consul.service_health_check_status",
		Priority: prioServiceHealthCheckStatus,
		Dims: module.Dims{
			{ID: "health_check_%s_passing_status", Name: "passing"},
			{ID: "health_check_%s_critical_status", Name: "critical"},
			{ID: "health_check_%s_maintenance_status", Name: "maintenance"},
			{ID: "health_check_%s_warning_status", Name: "warning"},
		},
	}
	chartTmplUnboundHealthCheckStatus = module.Chart{
		ID:       "health_check_%s_status",
		Title:    "Unbound health check status",
		Units:    "status",
		Fam:      "health checks",
		Ctx:      "consul.unbound_health_check_status",
		Priority: prioUnboundHealthCheckStatus,
		Dims: module.Dims{
			{ID: "health_check_%s_passing_status", Name: "passing"},
			{ID: "health_check_%s_critical_status", Name: "critical"},
			{ID: "health_check_%s_maintenance_status", Name: "maintenance"},
			{ID: "health_check_%s_warning_status", Name: "warning"},
		},
	}
)

func newServiceHealthCheckChart(check *healthCheck) *module.Chart {
	chart := chartTmplServiceHealthCheckStatus.Copy()
	chart.ID = fmt.Sprintf(chart.ID, check.CheckID)
	chart.Labels = []module.Label{
		{Key: "node", Value: check.Node},
		{Key: "service", Value: check.ServiceName},
	}
	for _, d := range chart.Dims {
		d.ID = fmt.Sprintf(d.ID, check.CheckID)
	}
	return chart
}

func newUnboundHealthCheckChart(check *healthCheck) *module.Chart {
	chart := chartTmplUnboundHealthCheckStatus.Copy()
	chart.ID = fmt.Sprintf(chart.ID, check.CheckID)
	chart.Labels = []module.Label{
		{Key: "node", Value: check.Node},
	}
	for _, d := range chart.Dims {
		d.ID = fmt.Sprintf(d.ID, check.CheckID)
	}
	return chart
}

func (c *Consul) addHealthCheckCharts(check *healthCheck) {
	var chart *module.Chart

	if check.ServiceName != "" {
		chart = newServiceHealthCheckChart(check)
	} else {
		chart = newUnboundHealthCheckChart(check)
	}

	if err := c.Charts().Add(chart); err != nil {
		c.Warning(err)
	}
}

func (c *Consul) removeHealthCheckCharts(checkID string) {
	id := fmt.Sprintf("health_check_%s_status", checkID)

	chart := c.Charts().Get(id)
	if chart == nil {
		c.Warningf("failed to remove '%s' chart: the chart does not exist", id)
		return
	}

	chart.MarkRemove()
	chart.MarkNotCreated()
}

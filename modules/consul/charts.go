// SPDX-License-Identifier: GPL-3.0-or-later

package consul

import (
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

const (
	_ = module.Priority + iota
	prioServerLeadershipStatus

	prioAutopilotHealthState
	prioAutopilotFailureTolerance

	prioMemoryAllocated
	prioMemorySys

	prioGCPauseTime

	prioRPCRequests

	prioServiceHealthCheckStatus
	prioUnboundHealthCheckStatus
)

var (
	globalCharts = module.Charts{
		chartServerLeadershipStatus.Copy(),
		chartAutopilotHealthState.Copy(),
		chartAutopilotFailureTolerance.Copy(),
		chartMemoryAllocated.Copy(),
		chartMemorySys.Copy(),
		chartGCPauseTime.Copy(),
		chartClientRPCRequestsRate.Copy(),
	}

	chartServerLeadershipStatus = module.Chart{
		ID:       "server_leadership_status",
		Title:    "Server leadership status",
		Units:    "status",
		Fam:      "leadership",
		Ctx:      "consul.server_leadership_status",
		Priority: prioServerLeadershipStatus,
		Dims: module.Dims{
			{ID: "consul.server.isLeader.yes", Name: "leader"},
			{ID: "consul.server.isLeader.no", Name: "not_leader"},
		},
	}

	chartAutopilotHealthState = module.Chart{
		ID:       "autopilot_health_state",
		Title:    "Autopilot health state",
		Units:    "state",
		Fam:      "autopilot",
		Ctx:      "consul.autopilot_health_state",
		Priority: prioAutopilotHealthState,
		Dims: module.Dims{
			{ID: "consul.autopilot.healthy.yes", Name: "healthy"},
			{ID: "consul.autopilot.healthy.no", Name: "unhealthy"},
		},
	}
	chartAutopilotFailureTolerance = module.Chart{
		ID:       "autopilot_failure_tolerance",
		Title:    "Autopilot failure tolerance",
		Units:    "servers",
		Fam:      "autopilot",
		Ctx:      "consul.autopilot_failure_tolerance",
		Priority: prioAutopilotFailureTolerance,
		Dims: module.Dims{
			{ID: "consul.autopilot.failure_tolerance", Name: "tolerance"},
		},
	}

	chartMemoryAllocated = module.Chart{
		ID:       "memory_allocated",
		Title:    "Memory allocated by the Consul process",
		Units:    "bytes",
		Fam:      "memory",
		Ctx:      "consul.memory_allocated",
		Priority: prioMemoryAllocated,
		Dims: module.Dims{
			{ID: "consul.runtime.alloc_bytes", Name: "allocated"},
		},
	}
	chartMemorySys = module.Chart{
		ID:       "memory_sys",
		Title:    "Memory obtained from the OS",
		Units:    "bytes",
		Fam:      "memory",
		Ctx:      "consul.memory_sys",
		Priority: prioMemorySys,
		Dims: module.Dims{
			{ID: "consul.runtime.sys_bytes", Name: "sys"},
		},
	}

	chartGCPauseTime = module.Chart{
		ID:       "gc_pause_time",
		Title:    "Garbage collection stop-the-world pause time",
		Units:    "seconds",
		Fam:      "garbage collection",
		Ctx:      "consul.gc_pause_time",
		Priority: prioGCPauseTime,
		Dims: module.Dims{
			{ID: "consul.runtime.total_gc_pause_ns", Name: "gc_pause", Algo: module.Incremental, Div: 1e9},
		},
	}

	chartClientRPCRequestsRate = module.Chart{
		ID:       "client_rpc_requests_rate",
		Title:    "Client RPC requests",
		Units:    "requests/s",
		Fam:      "client rpc",
		Ctx:      "consul.client_rpc_requests_rate",
		Priority: prioRPCRequests,
		Dims: module.Dims{
			{ID: "consul.client.rpc", Name: "rpc", Algo: module.Incremental},
		},
	}
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

func newServiceHealthCheckChart(check *agentCheck) *module.Chart {
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

func newUnboundHealthCheckChart(check *agentCheck) *module.Chart {
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

func (c *Consul) addHealthCheckCharts(check *agentCheck) {
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

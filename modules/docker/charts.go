// SPDX-License-Identifier: GPL-3.0-or-later

package docker

import "github.com/netdata/go.d.plugin/agent/module"

const (
	prioContainersState = module.Priority + iota
	prioContainersHealthy
	prioContainersUnhealthy
)

var charts = module.Charts{
	containersStateChart.Copy(),
	containersHealthyChart.Copy(),
	containersUnhealthyChart.Copy(),
}

var containersStateChart = module.Chart{
	ID:       "containers_state",
	Title:    "Number of containers in different states",
	Units:    "containers",
	Fam:      "containers state",
	Ctx:      "docker.containers_state",
	Priority: prioContainersState,
	Type:     module.Stacked,
	Dims: module.Dims{
		{ID: "running_containers", Name: "running"},
		{ID: "paused_containers", Name: "paused"},
		{ID: "stopped_containers", Name: "stopped"},
	},
}

var containersHealthyChart = module.Chart{
	ID:       "healthy_containers",
	Title:    "Number of healthy containers",
	Units:    "containers",
	Fam:      "containers health",
	Ctx:      "docker.healthy_containers",
	Priority: prioContainersHealthy,
	Dims: module.Dims{
		{ID: "healthy_containers", Name: "healthy"},
	},
}

var containersUnhealthyChart = module.Chart{
	ID:       "unhealthy_containers",
	Title:    "Number of unhealthy containers",
	Units:    "containers",
	Fam:      "containers health",
	Ctx:      "docker.unhealthy_containers",
	Priority: prioContainersUnhealthy,
	Dims: module.Dims{
		{ID: "unhealthy_containers", Name: "unhealthy"},
	},
}

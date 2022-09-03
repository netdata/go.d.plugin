// SPDX-License-Identifier: GPL-3.0-or-later

package docker

import "github.com/netdata/go.d.plugin/agent/module"

const (
	prioContainersState = module.Priority + iota
	prioContainersHealthy
	prioContainersUnhealthy
	prioImagesCount
	prioImagesSize
)

var charts = module.Charts{
	containersStateChart.Copy(),
	containersHealthyChart.Copy(),
	containersUnhealthyChart.Copy(),

	imagesCountChart.Copy(),
	imagesSizeChart.Copy(),
}

var (
	containersStateChart = module.Chart{
		ID:       "containers_state",
		Title:    "Number of containers in different states",
		Units:    "containers",
		Fam:      "containers",
		Ctx:      "docker.containers_state",
		Priority: prioContainersState,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "running_containers", Name: "running"},
			{ID: "paused_containers", Name: "paused"},
			{ID: "exited_containers", Name: "exited"},
		},
	}
	containersHealthyChart = module.Chart{
		ID:       "healthy_containers",
		Title:    "Number of healthy containers",
		Units:    "containers",
		Fam:      "containers",
		Ctx:      "docker.healthy_containers",
		Priority: prioContainersHealthy,
		Dims: module.Dims{
			{ID: "healthy_containers", Name: "healthy"},
		},
	}
	containersUnhealthyChart = module.Chart{
		ID:       "unhealthy_containers",
		Title:    "Number of unhealthy containers",
		Units:    "containers",
		Fam:      "containers",
		Ctx:      "docker.unhealthy_containers",
		Priority: prioContainersUnhealthy,
		Dims: module.Dims{
			{ID: "unhealthy_containers", Name: "unhealthy"},
		},
	}
)

var (
	imagesCountChart = module.Chart{
		ID:       "images_count",
		Title:    "Number of images",
		Units:    "images",
		Fam:      "images",
		Ctx:      "docker.images",
		Priority: prioImagesCount,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "images_active", Name: "active"},
			{ID: "images_dangling", Name: "dangling"},
		},
	}
	imagesSizeChart = module.Chart{
		ID:       "images_size",
		Title:    "Images size",
		Units:    "B",
		Fam:      "images",
		Ctx:      "docker.images_size",
		Priority: prioImagesSize,
		Dims: module.Dims{
			{ID: "images_size", Name: "size"},
		},
	}
)

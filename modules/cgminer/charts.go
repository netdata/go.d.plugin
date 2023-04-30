// SPDX-License-Identifier: GPL-3.0-or-later

package cgminer

import (
	"fmt"
	"github.com/netdata/go.d.plugin/agent/module"
)

var chartTemplate = module.Chart{
	ID:    "cgminer_%s",
	Title: "Cgminer Metric",
	Units: "metric",
	Fam:   "cgminer",
	Ctx:   "cgminer.metric",
}

func newChart(id string) *module.Chart {
	chart := chartTemplate.Copy()
	chart.ID = fmt.Sprintf(chart.ID, id)
	return chart
}

func charts() *module.Charts {
	return &module.Charts{
		{
			ID:    "cgminer_hashrate",
			Title: "Hashrate",
			Units: "hash/s",
			Fam:   "cgminer",
			Ctx:   "cgminer.hashrate",
			Type:  module.Stacked,
			Dims: module.Dims{
				{ID: "hashrate_gpu0", Name: "GPU 0"},
				{ID: "hashrate_gpu1", Name: "GPU 1"},
				{ID: "hashrate_gpu2", Name: "GPU 2"},
				{ID: "hashrate_gpu3", Name: "GPU 3"},
			},
		},
		{
			ID:    "cgminer_accepted",
			Title: "Accepted Shares",
			Units: "shares",
			Fam:   "cgminer",
			Ctx:   "cgminer.accepted",
			Type:  module.Stacked,
			Dims: module.Dims{
				{ID: "accepted_gpu0", Name: "GPU 0"},
				{ID: "accepted_gpu1", Name: "GPU 1"},
				{ID: "accepted_gpu2", Name: "GPU 2"},
				{ID: "accepted_gpu3", Name: "GPU 3"},
			},
		},
		{
			ID:    "cgminer_rejected",
			Title: "Rejected Shares",
			Units: "shares",
			Fam:   "cgminer",
			Ctx:   "cgminer.rejected",
			Type:  module.Stacked,
			Dims: module.Dims{
				{ID: "rejected_gpu0", Name: "GPU 0"},
				{ID: "rejected_gpu1", Name: "GPU 1"},
				{ID: "rejected_gpu2", Name: "GPU 2"},
				{ID: "rejected_gpu3", Name: "GPU 3"},
			},
		},
		{
			ID:    "cgminer_temperature",
			Title: "Temperature",
			Units: "celsius",
			Fam:   "cgminer",
			Ctx:   "cgminer.temperature",
			Type:  module.Stacked,
			Dims: module.Dims{
				{ID: "temperature_gpu0", Name: "GPU 0"},
				{ID: "temperature_gpu1", Name: "GPU 1"},
				{ID: "temperature_gpu2", Name: "GPU 2"},
				{ID: "temperature_gpu3", Name: "GPU 3"},
			},
		},
		{
			ID:    "cgminer_fan_speed",
			Title: "Fan Speed",
			Units: "RPM",
			Fam:   "cgminer",
			Ctx:   "cgminer.fan_speed",
			Type:  module.Stacked,
			Dims: module.Dims{
				{ID: "fan_speed_gpu0", Name: "GPU 0"},
				{ID: "fan_speed_gpu1", Name: "GPU 1"},
				{ID: "fan_speed_gpu2", Name: "GPU 2"},
				{ID: "fan_speed_gpu3", Name: "GPU 3"},
			},
		},
	}
}

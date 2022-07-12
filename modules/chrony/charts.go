// SPDX-License-Identifier: GPL-3.0-or-later

package chrony

import (
	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/modules/chrony/client"
)

const scaleFactor = client.ScaleFactor

var charts = module.Charts{
	{
		ID:    "stratum",
		Title: "The distance from the reference clock",
		Units: "level",
		Fam:   "stratum",
		Ctx:   "chrony.stratum",
		Dims: module.Dims{
			{ID: "stratum", Name: "stratum"},
		},
	},
	{
		ID:    "ref_timestamp",
		Title: "last clock update interval",
		Units: "seconds",
		Fam:   "ref timestamp",
		Ctx:   "chrony.ref_timestamp",
		Dims: module.Dims{
			{ID: "ref_timestamp"},
		},
	},
	{
		ID:    "current_correction",
		Title: "Current correction",
		Units: "seconds",
		Fam:   "current correction",
		Ctx:   "chrony.current_correction",
		Dims: module.Dims{
			{ID: "current_correction", Div: scaleFactor},
		},
	},
	{
		ID:    "offset",
		Title: "The offset between clock update",
		Units: "seconds",
		Fam:   "offset",
		Ctx:   "chrony.offset",
		Dims: module.Dims{
			{ID: "last_offset", Name: "last", Div: scaleFactor},
			{ID: "rms_offset", Name: "rms", Div: scaleFactor},
		},
	},
	{
		ID:    "frequency",
		Title: "The rate at which the system clock would be wrong if the chronyd did not correct them",
		Units: "ppm",
		Fam:   "frequency",
		Ctx:   "chrony.frequency",
		Dims: module.Dims{
			{ID: "frequency", Div: scaleFactor},
		},
	},
	{
		ID:    "skew",
		Title: "The estimated error bound on the frequency",
		Units: "ppm",
		Fam:   "skew",
		Ctx:   "chrony.skew",
		Dims: module.Dims{
			{ID: "skew", Div: scaleFactor},
		},
	},
	{
		ID:    "root_delay",
		Title: "The network path delay to the stratum-1 computer",
		Units: "seconds",
		Fam:   "root delay",
		Ctx:   "chrony.root_delay",
		Dims: module.Dims{
			{ID: "root_delay", Div: scaleFactor},
		},
	},
	{
		ID:    "root_dispersion",
		Title: "The dispersion accumulated through all the computers back to the stratum-1 computer",
		Units: "seconds",
		Fam:   "root dispersion",
		Ctx:   "chrony.root_dispersion",
		Dims: module.Dims{
			{ID: "root_dispersion", Div: scaleFactor},
		},
	},
	{
		ID:    "update_interval",
		Title: "The interval between the last two clock updates",
		Units: "seconds",
		Fam:   "update interval",
		Ctx:   "chrony.update_interval",
		Dims: module.Dims{
			{ID: "update_interval", Div: scaleFactor},
		},
	},
	{
		ID:    "leap_status",
		Title: "Leap status",
		Units: "status",
		Fam:   "leap status",
		Ctx:   "chrony.leap_status",
		Dims: module.Dims{
			{ID: "leap_status_normal", Name: "normal"},
			{ID: "leap_status_insert_second", Name: "insert_second"},
			{ID: "leap_status_delete_second", Name: "delete_second"},
			{ID: "leap_status_unsynchronised", Name: "unsynchronised"},
		},
	},
	{
		ID:    "activity",
		Title: "Peers activity",
		Units: "sources",
		Fam:   "activity",
		Ctx:   "chrony.activity",
		Type:  module.Stacked,
		Dims: module.Dims{
			{ID: "online_sources", Name: "online"},
			{ID: "offline_sources", Name: "offline"},
			{ID: "burst_online_sources", Name: "burst_online"},
			{ID: "burst_offline_sources", Name: "burst_offline"},
			{ID: "unresolved_sources", Name: "unresolved"},
		},
	},
}

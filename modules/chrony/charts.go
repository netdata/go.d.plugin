// SPDX-License-Identifier: GPL-3.0-or-later

package chrony

import (
	"github.com/netdata/go.d.plugin/agent/module"
	"net"
)

var charts = module.Charts{
	{
		ID:    "running",
		Title: "chrony is functional and can be monitored",
		Units: "hop",
		Type:  module.Area,
		Ctx:   "chrony.running",
		Dims: module.Dims{
			{ID: "running", Name: "running", Algo: module.Absolute, Div: 1, Mul: 1},
		},
	},
	{
		ID:    "stratum",
		Title: "distance from reference clock",
		Units: "level",
		Type:  module.Area,
		Ctx:   "chrony.stratum",
		Dims: module.Dims{
			{ID: "stratum", Name: "stratum", Algo: module.Absolute, Div: 1, Mul: 1},
		},
	},
	{
		ID: "leap_status",
		//  LEAP_Normal = 0,
		//  LEAP_InsertSecond = 1,
		//  LEAP_DeleteSecond = 2,
		//  LEAP_Unsynchronised = 3
		Title: "Leap status can be Normal, Insert second, Delete second or Not synchronised.",
		Units: "hop",
		Ctx:   "chrony.leap_status",
		Dims: module.Dims{
			{ID: "leap_status", Name: "leap_status", Algo: module.Absolute, Div: 1, Mul: 1},
		},
	},
	{
		ID:    "root_delay",
		Title: "the total of the network path delays to the stratum-1 computer",
		Units: "seconds",
		Type:  module.Area,
		Ctx:   "chrony.root_delay",
		Dims: module.Dims{
			{ID: "root_delay", Name: "root_delay", Algo: module.Absolute, Div: scaleFactor, Mul: 1},
		},
	},
	{
		ID:    "root_dispersion",
		Title: "total dispersion accumulated through all the computers back to the stratum-1 computer",
		Units: "seconds",
		Type:  module.Area,
		Ctx:   "chrony.root_dispersion",
		Dims: module.Dims{
			{ID: "root_dispersion", Name: "root_dispersion", Algo: module.Absolute, Div: scaleFactor, Mul: 1},
		},
	},
	{
		ID:    "skew",
		Title: "estimated error bound on the frequency",
		Units: "ppm",
		Type:  module.Area,
		Ctx:   "chrony.skew",
		Dims: module.Dims{
			{ID: "skew", Name: "skew", Algo: module.Absolute, Div: scaleFactor, Mul: 1},
		},
	},
	{
		ID:    "frequency",
		Title: "the rate by which the system’s clock would be would be wrong",
		Units: "ppm",
		Type:  module.Area,
		Ctx:   "chrony.frequency",
		Dims: module.Dims{
			{ID: "frequency", Name: "frequency", Algo: module.Absolute, Div: scaleFactor, Mul: 1},
		},
	},
	{
		ID:    "offset",
		Title: "the offset between clock update",
		Units: "seconds",
		Type:  module.Area,
		Ctx:   "chrony.offset",
		Dims: module.Dims{
			{ID: "last_offset", Name: "last", Algo: module.Absolute, Div: scaleFactor, Mul: 1},
			{ID: "rms_offset", Name: "rms", Algo: module.Absolute, Div: scaleFactor, Mul: 1},
		},
	},
	{
		ID:    "update_interval",
		Title: "last clock update interval",
		Units: "seconds",
		Type:  module.Area,
		Ctx:   "chrony.update_interval",
		Dims: module.Dims{
			{ID: "update_interval", Name: "update_interval", Algo: module.Absolute, Div: scaleFactor, Mul: 1},
		},
	},
	{
		ID:    "current_correction",
		Title: "last clock update interval",
		Units: "seconds",
		Type:  module.Area,
		Ctx:   "chrony.current_correction",
		Dims: module.Dims{
			{ID: "current_correction", Name: "current_correction", Algo: module.Absolute, Div: scaleFactor, Mul: 1},
		},
	},
	{
		ID:    "ref_timestamp",
		Title: "last clock update interval",
		Units: "seconds",
		Type:  module.Line,
		Ctx:   "chrony.ref_timestamp",
		Dims: module.Dims{
			{ID: "ref_timestamp", Name: "ref_timestamp", Algo: module.Absolute, Div: 1, Mul: 1},
		},
	},
	{
		ID:    "activity",
		Title: "activity status",
		Units: "count",
		Ctx:   "chrony.activity",
		Type:  module.Area,
		Dims: module.Dims{
			{ID: "online_sources", Name: "online_sources", Algo: module.Absolute, Div: 1, Mul: 1},
			{ID: "offline_sources", Name: "offline_sources", Algo: module.Absolute, Div: 1, Mul: 1},
			{ID: "burst_online_sources", Name: "burst_online_sources", Algo: module.Absolute, Div: 1, Mul: 1},
			{ID: "burst_offline_sources", Name: "burst_offline_sources", Algo: module.Absolute, Div: 1, Mul: 1},
			{ID: "unresolved_sources", Name: "unresolved_sources", Algo: module.Absolute, Div: 1, Mul: 1},
		},
	},
	{
		ID:    "source",
		Title: "Activity Source Server",
		Units: "hop",
		Ctx:   "chrony.source",
		Type:  module.Area,
		Dims: module.Dims{
			{ID: net.IPv4zero.String(), Name: net.IPv4zero.String(), Algo: module.Absolute, Div: 1, Mul: 1},
		},
	},
}

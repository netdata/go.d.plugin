package chrony

import "github.com/netdata/go.d.plugin/agent/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

var charts = Charts{
	{
		ID:    "running",
		Title: "chrony is functional and can be monitored",
		Units: "Hop",
		Type:  module.Area,
		Ctx:   "chrony.running",
		Dims: Dims{
			{ID: "running", Name: "running", Algo: module.Absolute, Div: 1, Mul: 1},
		},
	},
	{
		ID:    "stratum",
		Title: "distance from reference clock",
		Units: "Hop",
		Type:  module.Area,
		Ctx:   "chrony.stratum",
		Dims: Dims{
			{ID: "stratum", Name: "stratum", Algo: module.Absolute, Div: 1, Mul: 1},
		},
	},
	{
		ID:    "leap_status",
		Title: "Leap status can be Normal, Insert second, Delete second or Not synchronised.",
		Ctx:   "chrony.leap_status",
		Dims: Dims{
			{ID: "leap_status", Name: "leap_status", Algo: module.Absolute, Div: 1, Mul: 1},
		},
	},
	{
		ID:    "root_delay",
		Title: "the total of the network path delays to the stratum-1 computer",
		Units: "seconds",
		Type:  module.Area,
		Ctx:   "chrony.root_delay",
		Dims: Dims{
			{ID: "root_delay", Name: "root_delay", Algo: module.Absolute, Div: scaleFactor, Mul: 1},
		},
	},
	{
		ID:    "root_dispersion",
		Title: "total dispersion accumulated through all the computers back to the stratum-1 computer",
		Units: "seconds",
		Type:  module.Area,
		Ctx:   "chrony.root_dispersion",
		Dims: Dims{
			{ID: "root_dispersion", Name: "root_dispersion", Algo: module.Absolute, Div: scaleFactor, Mul: 1},
		},
	},
	{
		ID:    "skew",
		Title: "estimated error bound on the frequency",
		Units: "ppm",
		Type:  module.Area,
		Ctx:   "chrony.skew",
		Dims: Dims{
			{ID: "skew", Name: "skew", Algo: module.Absolute, Div: scaleFactor, Mul: 1},
		},
	},
	{
		ID:    "frequency",
		Title: "the rate by which the system’s clock would be would be wrong",
		Units: "ppm",
		Type:  module.Area,
		Ctx:   "chrony.frequency",
		Dims: Dims{
			{ID: "frequency", Name: "frequency", Algo: module.Absolute, Div: scaleFactor, Mul: 1},
		},
	},
	{
		ID:    "offset",
		Title: "the offset between clock update",
		Units: "seconds",
		Type:  module.Area,
		Ctx:   "chrony.offset",
		Dims: Dims{
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
		Dims: Dims{
			{ID: "update_interval", Name: "update_interval", Algo: module.Absolute, Div: scaleFactor, Mul: 1},
		},
	},
	{
		ID:    "current_correction",
		Title: "last clock update interval",
		Units: "seconds",
		Type:  module.Area,
		Ctx:   "chrony.current_correction",
		Dims: Dims{
			{ID: "current_correction", Name: "current_correction", Algo: module.Absolute, Div: scaleFactor, Mul: 1},
		},
	},
	{
		ID:    "ref_timestamp",
		Title: "last clock update interval",
		Units: "seconds",
		Type:  module.Line,
		Ctx:   "chrony.ref_timestamp",
		Dims: Dims{
			{ID: "ref_timestamp", Name: "ref_timestamp", Algo: module.Absolute, Div: 1, Mul: 1},
		},
	},
	{
		ID:    "activity",
		Title: "activity status",
		Units: "Hop",
		Ctx:   "chrony.activity",
		Type:  module.Area,
		Dims: Dims{
			{ID: "online_sources", Name: "online_sources", Algo: module.Absolute, Div: 1, Mul: 1},
			{ID: "offline_sources", Name: "offline_sources", Algo: module.Absolute, Div: 1, Mul: 1},
			{ID: "burst_online_sources", Name: "burst_online_sources", Algo: module.Absolute, Div: 1, Mul: 1},
			{ID: "burst_offline_sources", Name: "burst_offline_sources", Algo: module.Absolute, Div: 1, Mul: 1},
			{ID: "unresolved_sources", Name: "unresolved_sources", Algo: module.Absolute, Div: 1, Mul: 1},
		},
	},
}

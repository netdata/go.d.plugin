package apache

import "github.com/netdata/go.d.plugin/modules"

type (
	// Charts is an alias for modules.Charts
	Charts = modules.Charts
	// Dims is an alias for modules.Dims
	Dims = modules.Dims
)

var charts = Charts{
	{
		ID:    "requests",
		Title: "Requests",
		Units: "requests/s",
		Fam:   "statistics",
		Ctx:   "apache.requests",
		Dims: Dims{
			{ID: "requests", Algo: modules.Incremental},
		},
	},
	{
		ID:    "connections",
		Title: "Connections",
		Units: "connections",
		Fam:   "connections",
		Ctx:   "apache.connections",
		Dims: Dims{
			{ID: "connections"},
		},
	},
	{
		ID:    "conns_async",
		Title: "Async Connections",
		Units: "connections",
		Fam:   "connections",
		Ctx:   "apache.conns_async",
		Type:  modules.Stacked,
		Dims: Dims{
			{ID: "keepalive"},
			{ID: "closing"},
			{ID: "writing"},
		},
	},
	{
		ID:    "net",
		Title: "Bandwidth",
		Units: "kilobits/s",
		Fam:   "bandwidth",
		Ctx:   "apache.net",
		Type:  modules.Area,
		Dims: Dims{
			{ID: "sent", Algo: modules.Incremental, Mul: 8},
		},
	},
	{
		ID:    "workers",
		Title: "Workers",
		Units: "workers",
		Fam:   "workers",
		Ctx:   "apache.workers",
		Type:  modules.Stacked,
		Dims: Dims{
			{ID: "idle"},
			{ID: "busy"},
		},
	},
	{
		ID:    "reqpersec",
		Title: "Lifetime Average Requests",
		Units: "requests/s",
		Fam:   "statistics",
		Ctx:   "apache.reqpersec",
		Type:  modules.Area,
		Dims: Dims{
			{ID: "requests_sec"},
		},
	},
	{
		ID:    "bytespersec",
		Title: "Lifetime Average Bandwidth",
		Units: "kilobits/s",
		Fam:   "statistics",
		Ctx:   "apache.bytesperreq",
		Type:  modules.Area,
		Dims: Dims{
			{ID: "size_sec", Mul: 8, Div: 1000},
		},
	},
	{
		ID:    "bytesperreq",
		Title: "Lifetime Average Response Size",
		Units: "bytes/request",
		Fam:   "statistics",
		Ctx:   "apache.bytesperreq",
		Type:  modules.Area,
		Dims: Dims{
			{ID: "size_req"},
		},
	},
}

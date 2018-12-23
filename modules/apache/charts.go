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
		Fam:   "requests",
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
		Title: "Lifetime Average Requests Per Second",
		Units: "requests/s",
		Fam:   "statistics",
		Ctx:   "apache.reqpersec",
		Type:  modules.Area,
		Dims: Dims{
			{ID: "requests_sec", Div: 100000},
		},
	},
	{
		ID:    "bytespersec",
		Title: "Lifetime Avgerage Served Per Second",
		Units: "KiB",
		Fam:   "statistics",
		Ctx:   "apache.bytespersec",
		Type:  modules.Area,
		Dims: Dims{
			{ID: "size_sec", Mul: 8, Div: 1024 * 100000},
		},
	},
	{
		ID:    "bytesperreq",
		Title: "Lifetime Average Response Size",
		Units: "KiB",
		Fam:   "statistics",
		Ctx:   "apache.bytesperreq",
		Type:  modules.Area,
		Dims: Dims{
			{ID: "size_req", Div: 1024 * 100000},
		},
	},
	{
		ID:    "scoreboard",
		Title: "Scoreboard",
		Units: "events",
		Fam:   "scoreboard",
		Ctx:   "apache.scoreboard",
		Dims: Dims{
			{ID: "scoreboard_waiting", Name: "waiting"},
			{ID: "scoreboard_starting", Name: "starting"},
			{ID: "scoreboard_reading", Name: "reading"},
			{ID: "scoreboard_sending", Name: "sending"},
			{ID: "scoreboard_keepalive", Name: "keepalive"},
			{ID: "scoreboard_dns_lookup", Name: "dns lookup"},
			{ID: "scoreboard_closing", Name: "closing"},
			{ID: "scoreboard_logging", Name: "logging"},
			{ID: "scoreboard_finishing", Name: "finishing"},
			{ID: "scoreboard_idle_cleanup", Name: "idle cleanup"},
			{ID: "scoreboard_open", Name: "open"},
		},
	},
	{
		ID:    "uptime",
		Title: "Uptime",
		Units: "seconds",
		Fam:   "uptime",
		Ctx:   "lighttpd.uptime",
		Dims: Dims{
			{ID: "uptime"},
		},
	},
}

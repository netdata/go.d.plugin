package apache

import "github.com/netdata/go.d.plugin/agent/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

var charts = Charts{
	{
		ID:    "requests",
		Title: "Requests",
		Units: "requests/s",
		Fam:   "requests",
		Ctx:   "apache.requests",
		Dims: Dims{
			{ID: "total_accesses", Name: "requests", Algo: module.Incremental},
		},
	},
	{
		ID:    "connections",
		Title: "Connections",
		Units: "connections",
		Fam:   "connections",
		Ctx:   "apache.connections",
		Dims: Dims{
			{ID: "conns_total", Name: "connections"},
		},
	},
	{
		ID:    "conns_async",
		Title: "Async Connections",
		Units: "connections",
		Fam:   "connections",
		Ctx:   "apache.conns_async",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "conns_async_keep_alive", Name: "keepalive"},
			{ID: "conns_async_closing", Name: "closing"},
			{ID: "conns_async_writing", Name: "writing"},
		},
	},
	{
		ID:    "scoreboard",
		Title: "Scoreboard",
		Units: "connections",
		Fam:   "connections",
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
		ID:    "net",
		Title: "Bandwidth",
		Units: "kilobits/s",
		Fam:   "bandwidth",
		Ctx:   "apache.net",
		Type:  module.Area,
		Dims: Dims{
			{ID: "total_kBytes", Name: "sent", Algo: module.Incremental, Mul: 8},
		},
	},
	{
		ID:    "workers",
		Title: "Workers Threads",
		Units: "workers",
		Fam:   "workers",
		Ctx:   "apache.workers",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "idle_workers", Name: "idle"},
			{ID: "busy_workers", Name: "busy"},
		},
	},
	{
		ID:    "reqpersec",
		Title: "Lifetime Average Number Of Requests Per Second",
		Units: "requests/s",
		Fam:   "statistics",
		Ctx:   "apache.reqpersec",
		Type:  module.Area,
		Dims: Dims{
			{ID: "req_per_sec", Name: "requests", Div: 100000},
		},
	},
	{
		ID:    "bytespersec",
		Title: "Lifetime Average Number Of Bytes Served Per Second",
		Units: "KiB/s",
		Fam:   "statistics",
		Ctx:   "apache.bytespersec",
		Type:  module.Area,
		Dims: Dims{
			{ID: "bytes_per_sec", Name: "served", Mul: 8, Div: 1024 * 100000},
		},
	},
	{
		ID:    "bytesperreq",
		Title: "Lifetime Average Response Size",
		Units: "KiB",
		Fam:   "statistics",
		Ctx:   "apache.bytesperreq",
		Type:  module.Area,
		Dims: Dims{
			{ID: "bytes_per_req", Name: "size", Div: 1024 * 100000},
		},
	},
	{
		ID:    "uptime",
		Title: "Uptime",
		Units: "seconds",
		Fam:   "availability",
		Ctx:   "apache.uptime",
		Dims: Dims{
			{ID: "uptime"},
		},
	},
}

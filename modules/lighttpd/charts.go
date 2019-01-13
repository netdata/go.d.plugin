package lighttpd

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
		Ctx:   "lighttpd.requests",
		Dims: Dims{
			{ID: "total_accesses", Name: "requests", Algo: modules.Incremental},
		},
	},
	{
		ID:    "net",
		Title: "Bandwidth",
		Units: "kilobits/s",
		Fam:   "bandwidth",
		Ctx:   "lighttpd.net",
		Type:  modules.Area,
		Dims: Dims{
			{ID: "total_kBytes", Name: "sent", Algo: modules.Incremental, Mul: 8},
		},
	},
	{
		ID:    "servers",
		Title: "Servers",
		Units: "servers",
		Fam:   "servers",
		Ctx:   "lighttpd.workers",
		Type:  modules.Stacked,
		Dims: Dims{
			{ID: "idle_servers", Name: "idle"},
			{ID: "busy_servers", Name: "busy"},
		},
	},
	{
		ID:    "scoreboard",
		Title: "ScoreBoard",
		Units: "events",
		Fam:   "scoreboard",
		Ctx:   "lighttpd.scoreboard",
		Dims: Dims{
			{ID: "scoreboard_waiting", Name: "waiting"},
			{ID: "scoreboard_open", Name: "open"},
			{ID: "scoreboard_close", Name: "close"},
			{ID: "scoreboard_hard_error", Name: "hard error"},
			{ID: "scoreboard_keepalive", Name: "keepalive"},
			{ID: "scoreboard_read", Name: "read"},
			{ID: "scoreboard_read_post", Name: "read post"},
			{ID: "scoreboard_write", Name: "write"},
			{ID: "scoreboard_handle_request", Name: "handle request"},
			{ID: "scoreboard_request_start", Name: "request start"},
			{ID: "scoreboard_request_end", Name: "request end"},
			{ID: "scoreboard_response_start", Name: "response start"},
			{ID: "scoreboard_response_end", Name: "response end"},
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

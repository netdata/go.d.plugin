package nginxvts

import "github.com/netdata/go.d.plugin/agent/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

var nginxVtsMainCharts = Charts{
	{
		ID:    "times",
		Title: "Nginx running time( Uptime=(nowMsec-loadMsec)/1000 )",
		Units: "milliseconds",
		Fam:   "main",
		Ctx:   "nginxvts.times",
		Dims: Dims{
			{ID: "loadmsec", Name: "load"},
			{ID: "nowmsec", Name: "up"},
		},
	},
	{
		ID:    "connections",
		Title: "Nginx Connections",
		Units: "connections",
		Fam:   "main",
		Ctx:   "nginxvts.connections",
		Dims: Dims{
			{ID: "connections_active", Name: "active"},
			{ID: "connections_reading", Name: "reading"},
			{ID: "connections_writing", Name: "writing"},
			{ID: "connections_waiting", Name: "waiting"},
			{ID: "connections_accepted", Name: "accepted"},
			{ID: "connections_handled", Name: "handled"},
			{ID: "connections_total", Name: "total"},
		},
	},
}
var nginxVtsSharedZonesChart = Charts{
	{
		ID:    "size",
		Title: "Shared memory size",
		Units: "bytes",
		Fam:   "sharedzones",
		Ctx:   "nginxvts.sharedzones_size",
		Dims: Dims{
			{ID: "sharedzones_maxsize", Name: "max size"},
			{ID: "sharedzones_usedsize", Name: "used size"},
		},
	},
	{
		ID:    "node",
		Title: "Number of node using shared memory",
		Units: "count",
		Fam:   "sharedzones",
		Ctx:   "nginxvts.sharedzones_node",
		Dims: Dims{
			{ID: "sharedzones_usednode", Name: "used node"},
		},
	},
}

var nginxVtsServerZonesCharts = Charts{
	{
		ID:    "requests",
		Title: "Number of client requests",
		Units: "requests/s",
		Fam:   "serverzones",
		Ctx:   "nginxvts.serverzones_requests",
		Dims: Dims{
			{ID: "total_requestcounter", Name: "requests", Algo: module.Incremental},
		},
	},
	{
		ID:    "responses",
		Title: "Total Response code",
		Units: "count",
		Fam:   "serverzones",
		Ctx:   "nginxvts.serverzones_responses",
		Dims: Dims{
			{ID: "total_responses_1xx", Name: "1xx"},
			{ID: "total_responses_2xx", Name: "2xx"},
			{ID: "total_responses_3xx", Name: "3xx"},
			{ID: "total_responses_4xx", Name: "4xx"},
			{ID: "total_responses_5xx", Name: "5xx"},
		},
	},
	{
		ID:    "traffic",
		Title: "Total server traffic",
		Units: "bytes/s",
		Fam:   "serverzones",
		Ctx:   "nginxvts.serverzones_traffic",
		Dims: Dims{
			{ID: "total_inbytes", Name: "inbytes", Algo: module.Incremental},
			{ID: "total_outbytes", Name: "outbytes", Algo: module.Incremental},
		},
	},
	{
		ID:    "cache",
		Title: "Total server cache",
		Units: "count",
		Fam:   "serverzones",
		Ctx:   "nginxvts.serverzones_cache",
		Dims: Dims{
			{ID: "total_cache_miss", Name: "miss"},
			{ID: "total_cache_bypass", Name: "bypass"},
			{ID: "total_cache_expired", Name: "expired"},
			{ID: "total_cache_stale", Name: "stale"},
			{ID: "total_cache_updating", Name: "updating"},
			{ID: "total_cache_revalidated", Name: "revalidated"},
			{ID: "total_cache_hit", Name: "hit"},
			{ID: "total_cache_scarce", Name: "scarce"},
		},
	},
}

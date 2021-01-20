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
		ID:    "total_requests",
		Title: "Nginx Total Requests",
		Units: "requests/s",
		Fam:   "main",
		Ctx:   "nginxvts.connections",
		Dims: Dims{
			{ID: "connections_requests", Name: "total requests", Algo: module.Incremental},
		},
	},
	{
		ID:    "active_connections",
		Title: "Nginx Active Connections",
		Units: "connections",
		Fam:   "main",
		Ctx:   "nginxvts.connections",
		Dims: Dims{
			{ID: "connections_active", Name: "active"},
		},
	},
	{
		ID:    "connections",
		Title: "Nginx Connections",
		Units: "requests/s",
		Fam:   "main",
		Ctx:   "nginxvts.connections",
		Dims: Dims{
			{ID: "connections_reading", Name: "reading", Algo: module.Incremental},
			{ID: "connections_writing", Name: "writing", Algo: module.Incremental},
			{ID: "connections_waiting", Name: "waiting", Algo: module.Incremental},
			{ID: "connections_accepted", Name: "accepted", Algo: module.Incremental},
			{ID: "connections_handled", Name: "handled", Algo: module.Incremental},
		},
	},
	{
		ID:    "uptime",
		Title: "Nginx Uptime",
		Units: "seconds",
		Fam:   "uptime",
		Ctx:   "nginxvts.uptime",
		Dims: Dims{
			{ID: "uptime", Name: "uptime"},
		},
	},
}
var nginxVtsSharedZonesChart = Charts{
	{
		ID:    "shared_memory_size",
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
		ID:    "shared_memory_used_nodes",
		Title: "Number of node using shared memory",
		Units: "nodes",
		Fam:   "sharedzones",
		Ctx:   "nginxvts.sharedzones_node",
		Dims: Dims{
			{ID: "sharedzones_usednode", Name: "used node"},
		},
	},
}

var nginxVtsServerZonesCharts = Charts{
	{
		ID:    "server_zones_requests_total",
		Title: "Number of client requests",
		Units: "requests/s",
		Fam:   "serverzones",
		Ctx:   "nginxvts.serverzones_requests_total",
		Dims: Dims{
			{ID: "total_requestcounter", Name: "requests", Algo: module.Incremental},
		},
	},
	{
		ID:    "responses",
		Title: "Total Response code",
		Units: "responses/s",
		Fam:   "serverzones",
		Ctx:   "nginxvts.serverzones_responses",
		Dims: Dims{
			{ID: "total_responses_1xx", Name: "1xx", Algo: module.Incremental},
			{ID: "total_responses_2xx", Name: "2xx", Algo: module.Incremental},
			{ID: "total_responses_3xx", Name: "3xx", Algo: module.Incremental},
			{ID: "total_responses_4xx", Name: "4xx", Algo: module.Incremental},
			{ID: "total_responses_5xx", Name: "5xx", Algo: module.Incremental},
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
		Units: "responses/s",
		Fam:   "serverzones",
		Ctx:   "nginxvts.serverzones_cache",
		Dims: Dims{
			{ID: "total_cache_miss", Name: "miss", Algo: module.Incremental},
			{ID: "total_cache_bypass", Name: "bypass", Algo: module.Incremental},
			{ID: "total_cache_expired", Name: "expired", Algo: module.Incremental},
			{ID: "total_cache_stale", Name: "stale", Algo: module.Incremental},
			{ID: "total_cache_updating", Name: "updating", Algo: module.Incremental},
			{ID: "total_cache_revalidated", Name: "revalidated", Algo: module.Incremental},
			{ID: "total_cache_hit", Name: "hit", Algo: module.Incremental},
			{ID: "total_cache_scarce", Name: "scarce", Algo: module.Incremental},
		},
	},
}

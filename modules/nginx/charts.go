package nginx

import "github.com/netdata/go-orchestrator/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

var charts = Charts{
	{
		ID:    "connections",
		Title: "Active Connections",
		Units: "connections",
		Fam:   "active connections",
		Ctx:   "nginx.connections",
		Dims: Dims{
			{ID: "active"},
		},
	},
	{
		ID:    "requests",
		Title: "Requests",
		Units: "requests/s",
		Fam:   "requests",
		Ctx:   "nginx.requests",
		Dims: Dims{
			{ID: "requests", Algo: module.Incremental},
		},
	},
	{
		ID:    "connection_statuses",
		Title: "Active Connections By Status",
		Units: "connections",
		Fam:   "status",
		Ctx:   "nginx.connection_status",
		Dims: Dims{
			{ID: "reading"},
			{ID: "writing"},
			{ID: "waiting", Name: "idle"},
		},
	},
	{
		ID:    "connect_rate",
		Title: "Connections Rate",
		Units: "connections/s",
		Fam:   "connections rate",
		Ctx:   "nginx.connect_rate",
		Dims: Dims{
			{ID: "accepts", Name: "accepted", Algo: module.Incremental},
			{ID: "handled", Algo: module.Incremental},
		},
	},
}

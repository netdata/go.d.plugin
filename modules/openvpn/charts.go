package openvpn

import "github.com/netdata/go-orchestrator/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

var charts = Charts{
	{
		ID:    "active_clients",
		Title: "Total Number Of Active Clients",
		Units: "active clients",
		Fam:   "clients",
		Ctx:   "openvpn.active_clients",
		Dims: Dims{
			{ID: "clients"},
		},
	},
	{
		ID:    "total_traffic",
		Title: "Total Traffic",
		Units: "KiB/s",
		Fam:   "traffic",
		Ctx:   "openvpn.total_traffic",
		Type:  module.Area,
		Dims: Dims{
			{ID: "bytes_in", Name: "in", Algo: module.Incremental, Div: 1 << 10},
			{ID: "bytes_out", Name: "out", Algo: module.Incremental, Div: -1 << 10},
		},
	},
}

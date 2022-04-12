package ovpn_status_log

import "github.com/netdata/go.d.plugin/agent/module"

var charts = module.Charts{
	{
		ID:    "active_clients",
		Title: "Total Number Of Active Clients",
		Units: "active clients",
		Fam:   "active_clients",
		Ctx:   "openvpn.active_clients",
		Dims: module.Dims{
			{ID: "clients"},
		},
	},
	{
		ID:    "total_traffic",
		Title: "Total Traffic",
		Units: "kilobits/s",
		Fam:   "total_traffic",
		Ctx:   "openvpn.total_traffic",
		Type:  module.Area,
		Dims: module.Dims{
			{ID: "bytes_in", Name: "in", Algo: module.Incremental, Mul: 8, Div: 1000},
			{ID: "bytes_out", Name: "out", Algo: module.Incremental, Mul: 8, Div: -1000},
		},
	},
}

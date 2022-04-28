package openvpn_status_log

import "github.com/netdata/go.d.plugin/agent/module"

var charts = module.Charts{
	{
		ID:    "active_clients",
		Title: "Total Number Of Active Clients",
		Units: "active clients",
		Fam:   "active_clients",
		Ctx:   "openvpn.active_clients",
		Dims: module.Dims{
			{
				ID: "clients",
			},
		},
	},
	{
		ID:    "total_traffic",
		Title: "Total Traffic",
		Units: "kilobits/s",
		Fam:   "traffic",
		Ctx:   "openvpn.total_traffic",
		Type:  module.Area,
		Dims: module.Dims{
			{
				ID:   "bytes_in",
				Name: "in",
				Algo: module.Incremental,
				Mul:  8,
				Div:  1024,
			},
			{
				ID:   "bytes_out",
				Name: "out",
				Algo: module.Incremental,
				Mul:  8,
				Div:  -1024,
			},
		},
	},
}

var userCharts = module.Charts{
	{
		ID:    "%s_user_traffic",
		Title: "User Traffic",
		Units: "kilobits/s",
		Fam:   "user %s",
		Ctx:   "openvpn.user_traffic",
		Type:  module.Area,
		Dims: module.Dims{
			{
				ID:   "%s_bytes_in",
				Name: "in",
				Algo: module.Incremental,
				Mul:  8,
				Div:  1024,
			},
			{
				ID:   "%s_bytes_out",
				Name: "out",
				Algo: module.Incremental,
				Mul:  8,
				Div:  -1024,
			},
		},
	},
	{
		ID:    "%s_user_connection_time",
		Title: "User Connection Time",
		Units: "seconds",
		Fam:   "user %s",
		Ctx:   "openvpn.user_connection_time",
		Dims: module.Dims{
			{
				ID:   "%s_connection_time",
				Name: "time",
			},
		},
	},
}

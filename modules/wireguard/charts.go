package wireguard

import "github.com/netdata/go.d.plugin/agent/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

var charts = Charts{
	{
		ID:    "total_data_%v",
		Title: "Total data Received/Sent (%v interface)",
		Units: "KB",
		Fam:   "total %v",
		Ctx:   "wireguard.total",
		Dims: Dims{
			{ID: "received_total", Name: "received", Div: 1024},
			{ID: "sent_total", Name: "sent", Div: 1024},
		},
	},
}

var bandwitchChart = Charts{
	{
		ID:    "bandwidth_%d",
		Title: "Peer %v Bandwidth",
		Units: "Kb/s",
		Fam:   "network %d",
		Ctx:   "wireguard.bandwidth",
		Dims: Dims{
			{ID: "received_%v", Name: "received", Div: 1000},
			{ID: "sent_%v", Name: "sent", Div: 1000},
		},
	},
}

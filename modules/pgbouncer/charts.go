package pgbouncer

import "github.com/netdata/go.d.plugin/agent/module"

const (
	prioClientConnections = module.Priority + iota
	prioServerConnections
)

var (
	globalCharts = module.Charts{
		clientConnectionsChart.Copy(),
		serverConnectionsChart.Copy(),
	}

	clientConnectionsChart = module.Chart{
		ID:       "client_connections",
		Title:    "Client connections",
		Units:    "connections",
		Fam:      "connections",
		Ctx:      "pgbouncer.client_connections",
		Type:     module.Stacked,
		Priority: prioClientConnections,
		Dims: module.Dims{
			{ID: "free_clients", Name: "free"},
			{ID: "used_clients", Name: "used"},
		},
	}
	serverConnectionsChart = module.Chart{
		ID:       "server_connections",
		Title:    "Server connections",
		Units:    "connections",
		Fam:      "connections",
		Ctx:      "pgbouncer.server_connections",
		Type:     module.Stacked,
		Priority: prioServerConnections,
		Dims: module.Dims{
			{ID: "free_servers", Name: "free"},
			{ID: "used_servers", Name: "used"},
		},
	}
)

package lighttpd2

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
		Ctx:   "lighttpd2.requests",
		Dims: Dims{
			{ID: "requests_abs", Name: "requests", Algo: module.Incremental},
		},
	},
	{
		ID:    "status_codes",
		Title: "Status Codes",
		Units: "requests/s",
		Fam:   "requests",
		Ctx:   "lighttpd2.status_codes",
		Dims: Dims{
			{ID: "status_1xx", Name: "1xx", Algo: module.Incremental},
			{ID: "status_2xx", Name: "2xx", Algo: module.Incremental},
			{ID: "status_3xx", Name: "3xx", Algo: module.Incremental},
			{ID: "status_4xx", Name: "4xx", Algo: module.Incremental},
			{ID: "status_5xx", Name: "5xx", Algo: module.Incremental},
		},
	},
	{
		ID:    "traffic",
		Title: "Traffic",
		Units: "kilobits/s",
		Fam:   "traffic",
		Ctx:   "lighttpd2.traffic",
		Type:  module.Area,
		Dims: Dims{
			{ID: "traffic_in_abs", Name: "in", Algo: module.Incremental, Mul: 8},
			{ID: "traffic_out_abs", Name: "out", Algo: module.Incremental, Mul: -8},
		},
	},
	{
		ID:    "connections",
		Title: "Connections",
		Units: "connections",
		Fam:   "connections",
		Ctx:   "lighttpd2.connections",
		Dims: Dims{
			{ID: "connection_abs", Name: "connections"},
		},
	},
	{
		ID:    "connection_states",
		Title: "Connection States",
		Units: "state",
		Fam:   "connections",
		Ctx:   "lighttpd2.connection_states",
		Dims: Dims{
			{ID: "connection_state_start", Name: "start"},
			{ID: "connection_state_read_header", Name: "read header"},
			{ID: "connection_state_handle_request", Name: "handle request"},
			{ID: "connection_state_write_response", Name: "write response"},
			{ID: "connection_state_keep_alive", Name: "keepalive"},
			{ID: "connection_state_upgraded", Name: "upgraded"},
		},
	},
	{
		ID:    "memory_usage",
		Title: "Memory Usage",
		Units: "KiB",
		Fam:   "memory",
		Ctx:   "lighttpd2.memory_usage",
		Dims: Dims{
			{ID: "memory_usage", Name: "usage", Div: 1024},
		},
	},
	{
		ID:    "uptime",
		Title: "Uptime",
		Units: "seconds",
		Fam:   "uptime",
		Ctx:   "lighttpd2.uptime",
		Dims: Dims{
			{ID: "uptime"},
		},
	},
}

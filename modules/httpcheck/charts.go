package httpcheck

import (
	"github.com/netdata/go.d.plugin/agent/module"
)

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Chart is an alias for module.Chart
	Chart = module.Chart
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

var charts = Charts{
	{
		ID:    "response_time",
		Title: "HTTP Response Time",
		Units: "ms",
		Fam:   "response",
		Ctx:   "httpcheck.response_time",
		Dims: Dims{
			{ID: "time"},
		},
	},
	{
		ID:    "response_length",
		Title: "HTTP Response Body Length",
		Units: "characters",
		Fam:   "response",
		Ctx:   "httpcheck.response_length",
		Dims: Dims{
			{ID: "length"},
		},
	},
	{
		ID:    "request_status",
		Title: "HTTP Check Status",
		Units: "boolean",
		Fam:   "status",
		Ctx:   "httpcheck.status",
		Dims: Dims{
			{ID: "success"},
			{ID: "no_connection", Name: "no connection"},
			{ID: "timeout"},
			{ID: "bad_content", Name: "bad content"},
			{ID: "bad_status", Name: "bad status"},
			//{ID: "dns_lookup_error", Name: "dns lookup error"},
			//{ID: "address_parse_error", Name: "address parse error"},
			//{ID: "redirect_error", Name: "redirect error"},
			//{ID: "body_read_error", Name: "body read error"},
		},
	},
	{
		ID:    "current_state_duration",
		Title: "HTTP Current State Duration",
		Units: "seconds",
		Fam:   "status",
		Ctx:   "httpcheck.in_state",
		Dims: Dims{
			{ID: "in_state", Name: "time"},
		},
	},
}

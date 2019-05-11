package httpcheck

import (
	"github.com/netdata/go-orchestrator/module"
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
			{ID: "response_time", Name: "time"},
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
			{ID: "no_connection"},
			{ID: "timeout"},
			{ID: "dns_lookup_error"},
			{ID: "address_parse_error"},
			{ID: "redirect_error"},
			{ID: "body_read_error"},
			{ID: "bad_content"},
			{ID: "bad_status"},
		},
	},
}

var bodyLengthChart = Chart{
	ID:    "response_length",
	Title: "HTTP Response Body Length",
	Units: "characters",
	Fam:   "response",
	Ctx:   "httpcheck.response_length",
	Dims: Dims{
		{ID: "response_length", Name: "length"},
	},
}

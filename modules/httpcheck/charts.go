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
			{ID: "response_time", Name: "time", Div: 1000000},
		},
	},
	{
		ID:    "response_length",
		Title: "HTTP Response Body Length",
		Units: "characters",
		Fam:   "response",
		Ctx:   "httpcheck.response_length",
		Dims: Dims{
			{ID: "response_length", Name: "length", Div: 1000000},
		},
	},
	{
		ID:    "request_status",
		Title: "HTTP Request Status",
		Units: "boolean",
		Fam:   "status",
		Ctx:   "httpcheck.status",
		Dims: Dims{
			{ID: "success"},
			{ID: "failed"},
			{ID: "timeout"},
		},
	},
	{
		ID:    "response_check_status",
		Title: "HTTP Response Status Code Check ",
		Units: "boolean",
		Fam:   "status",
		Ctx:   "httpcheck.check_status",
		Dims: Dims{
			{ID: "bad_status", Name: "bad"},
		},
	},
}

var respCheckContentChart = Chart{
	ID:    "response_check_content",
	Title: "HTTP Response Content Check",
	Units: "boolean",
	Fam:   "status",
	Ctx:   "httpcheck.check_content",
	Dims: Dims{
		{ID: "bad_content", Name: "bad content"},
	},
}

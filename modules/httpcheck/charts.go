package httpcheck

import (
	"github.com/netdata/go.d.plugin/modules"
)

type (
	// Charts is an alias for modules.Charts
	Charts = modules.Charts
	// Dims is an alias for modules.Dims
	Dims = modules.Dims
)

var charts = Charts{
	{
		ID:    "response_time",
		Title: "HTTP Response Time", Units: "ms", Fam: "response",
		Dims: Dims{
			{ID: "response_time", Name: "time", Div: 1000000},
		},
	},
	{
		ID:    "response_length",
		Title: "HTTP Response Body Length", Units: "characters", Fam: "response",
		Dims: Dims{
			{ID: "response_length", Name: "length", Div: 1000000},
		},
	},
	{
		ID:    "response_status",
		Title: "HTTP Response Status", Units: "boolean", Fam: "status",
		Dims: Dims{
			{ID: "success"},
			{ID: "failed"},
			{ID: "timeout"},
		},
	},
	{
		ID:    "response_check_status",
		Title: "HTTP Response Check Status", Units: "boolean", Fam: "status",
		Dims: Dims{
			{ID: "bad_status", Name: "bad status"},
		},
	},
	{
		ID:    "response_check_content",
		Title: "HTTP Response Check Content", Units: "boolean", Fam: "status",
		Dims: Dims{
			{ID: "bad_content", Name: "bad content"},
		},
	},
}

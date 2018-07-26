package httpcheck

import (
	"github.com/l2isbad/go.d.plugin/internal/pkg/charts"
)

type (
	Charts     = charts.Charts
	Options    = charts.Opts
	Dimensions = charts.Dims
	Dimension  = charts.Dim
)

var uCharts = Charts{
	{
		ID:   "response_time",
		Opts: Options{Title: "HTTP Response Time", Units: "ms", Family: "response"},
		Dims: Dimensions{
			{ID: "response_time", Name: "time", Div: 1000000},
		},
	},
	{
		ID:   "response_length",
		Opts: Options{Title: "HTTP Response Body Length", Units: "characters", Family: "response"},
		Dims: Dimensions{
			{ID: "response_length", Name: "length", Div: 1000000},
		},
	},
	{
		ID:   "response_status",
		Opts: Options{Title: "HTTP Response Status", Units: "boolean", Family: "status"},
		Dims: Dimensions{
			{ID: "success"},
			{ID: "failed"},
			{ID: "timeout"},
		},
	},
	{
		ID:   "response_check_status",
		Opts: Options{Title: "HTTP Response Check Status", Units: "boolean", Family: "status"},
		Dims: Dimensions{
			{ID: "bad_status", Name: "bad status"},
		},
	},
	{
		ID:   "response_check_content",
		Opts: Options{Title: "HTTP Response Check Content", Units: "boolean", Family: "status"},
		Dims: Dimensions{
			{ID: "bad_content", Name: "bad content"},
		},
	},
}

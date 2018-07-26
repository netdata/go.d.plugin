package httpcheck

import (
	"github.com/l2isbad/go.d.plugin/internal/pkg/charts"
)

type (
	Charts     = charts.Charts
	Options    = charts.Options
	Dimensions = charts.Dimensions
	Dimension  = charts.Dimension
)

var uCharts = Charts{
	{
		ID:      "response_time",
		Options: Options{Title: "HTTP Response Time", Units: "ms", Family: "response"},
		Dimensions: Dimensions{
			{ID: "response_time", Name: "time", Divisor: 1000000},
		},
	},
	{
		ID:      "response_length",
		Options: Options{Title: "HTTP Response Body Length", Units: "characters", Family: "response"},
		Dimensions: Dimensions{
			{ID: "response_length", Name: "length", Divisor: 1000000},
		},
	},
	{
		ID:      "response_status",
		Options: Options{Title: "HTTP Response Status", Units: "boolean", Family: "status"},
		Dimensions: Dimensions{
			{ID: "success"},
			{ID: "failed"},
			{ID: "timeout"},
		},
	},
	{
		ID:      "response_check_status",
		Options: Options{Title: "HTTP Response Check Status", Units: "boolean", Family: "status"},
		Dimensions: Dimensions{
			{ID: "bad_status", Name: "bad status"},
		},
	},
	{
		ID:      "response_check_content",
		Options: Options{Title: "HTTP Response Check Content", Units: "boolean", Family: "status"},
		Dimensions: Dimensions{
			{ID: "bad_content", Name: "bad content"},
		},
	},
}

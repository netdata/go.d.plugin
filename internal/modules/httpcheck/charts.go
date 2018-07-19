package httpcheck

import "github.com/l2isbad/go.d.plugin/internal/pkg/charts/raw"

type (
	Charts      = raw.Charts
	Order       = raw.Order
	Definitions = raw.Definitions
	Chart       = raw.Chart
	Options     = raw.Options
	Dimensions  = raw.Dimensions
	Dimension   = raw.Dimension
)

var charts = Charts{
	Order: Order{
		"response_time", "response_length", "response_status", "response_check_status", "response_check_content"},
	Definitions: Definitions{
		Chart{
			ID:      "response_time",
			Options: Options{"HTTP Response Time", "ms", "response", "httpcheck.response_time"},
			Dimensions: Dimensions{
				Dimension{"response_time", "time", "", 1, 1e6},
			},
		},
		Chart{
			ID:      "response_length",
			Options: Options{"HTTP Response Body Length", "characters", "response", "httpcheck.response_length"},
			Dimensions: Dimensions{
				Dimension{"response_length", "length"},
			},
		},
		Chart{
			ID:      "response_status",
			Options: Options{"HTTP Response Status", "boolean", "status", "httpcheck.status"},
			Dimensions: Dimensions{
				Dimension{"success"},
				Dimension{"failed"},
				Dimension{"timeout"},
			},
		},
		Chart{
			ID:      "response_check_status",
			Options: Options{"HTTP Response Check Status", "boolean", "status", "httpcheck.check_status"},
			Dimensions: Dimensions{
				Dimension{"bad_status", "bad status"},
			},
		},
		Chart{
			ID:      "response_check_content",
			Options: Options{"HTTP Response Check Content", "boolean", "status", "httpcheck.check_content"},
			Dimensions: Dimensions{
				Dimension{"bad_content", "bad content"},
			},
		},
	},
}

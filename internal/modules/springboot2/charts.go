package springboot2

import "github.com/l2isbad/go.d.plugin/internal/pkg/charts/raw"

var charts = raw.Charts{
	Order: raw.Order{"response_code", "threads", "gc_time", "gc_ope", "heap"},
	Definitions: raw.Definitions{
		{
			ID:      "response_code",
			Options: raw.Options{"Response Codes", "requests/s", "response"},
			Dimensions: raw.Dimensions{
				raw.Dimension{"response_time", "time", "", 1, 1e6},
			},
		},
		{
			ID:      "response_length",
			Options: raw.Options{"HTTP Response Body Length", "characters", "response"},
			Dimensions: raw.Dimensions{
				raw.Dimension{"response_length", "length"},
			},
		},
		raw.Chart{
			ID:      "response_status",
			Options: raw.Options{"HTTP Response Status", "boolean", "status"},
			Dimensions: raw.Dimensions{
				raw.Dimension{"success"},
				raw.Dimension{"failed"},
				raw.Dimension{"timeout"},
			},
		},
		raw.Chart{
			ID:      "response_check_status",
			Options: raw.Options{"HTTP Response Check Status", "boolean", "status"},
			Dimensions: raw.Dimensions{
				raw.Dimension{"bad_status", "bad status"},
			},
		},
		raw.Chart{
			ID:      "response_check_content",
			Options: raw.Options{"HTTP Response Check Content", "boolean", "status"},
			Dimensions: raw.Dimensions{
				raw.Dimension{"bad_content", "bad content"},
			},
		},
	},
}

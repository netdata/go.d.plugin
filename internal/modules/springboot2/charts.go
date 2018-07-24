package springboot2

import "github.com/l2isbad/go.d.plugin/internal/pkg/charts/raw"

var charts = raw.Charts{
	Order: raw.Order{"response_code", "threads", "gc_time", "gc_ope", "heap"},
	Definitions: raw.Definitions{
		{
			ID:      "response_code",
			Options: raw.Options{"Response Codes", "requests/s", "response"},
			Dimensions: raw.Dimensions{
				{"resp_1xx", "1xx", "incremental"},
				{"resp_2xx", "2xx", "incremental"},
				{"resp_3xx", "3xx", "incremental"},
				{"resp_4xx", "4xx", "incremental"},
				{"resp_5xx", "5xx", "incremental"},
				{"resp_other", "Other", "incremental"},
			},
		},
		{
			ID:      "threads",
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

package web_log

import "github.com/l2isbad/go.d.plugin/charts/raw"

type (
	Charts      = raw.Charts
	Order       = raw.Order
	Definitions = raw.Definitions
	Chart       = raw.Chart
	Options     = raw.Options
	Dimensions  = raw.Dimensions
	Dimension   = raw.Dimension
)

var uCharts = Charts{
	Order: Order{
		"response_statuses", "response_codes", "bandwidth", "response_time", "response_time_upstream",
		"requests_per_url", "requests_per_user_defined", "http_method", "http_version",
		"requests_per_ipproto", "clients", "clients_all",
	},
	Definitions: Definitions{
		Chart{
			ID:      "response_statuses",
			Options: Options{"Response Statuses", "requests/s", "responses", "", raw.Stacked},
			Dimensions: Dimensions{
				Dimension{"successful_requests", "success", raw.Incremental},
				Dimension{"server_errors", "error", raw.Incremental},
				Dimension{"redirects", "redirect", raw.Incremental},
				Dimension{"bad_requests", "bad", raw.Incremental},
				Dimension{"other_requests", "other", raw.Incremental},
			},
		},
		Chart{
			ID:      "response_codes",
			Options: Options{"Response Codes", "requests/s", "responses", "", raw.Stacked},
			Dimensions: Dimensions{
				Dimension{"2xx", "", raw.Incremental},
				Dimension{"5xx", "", raw.Incremental},
				Dimension{"3xx", "", raw.Incremental},
				Dimension{"4xx", "", raw.Incremental},
				Dimension{"1xx", "", raw.Incremental},
				Dimension{"0xx", "", raw.Incremental},
				Dimension{"unmatched", "", raw.Incremental},
			},
		},
		Chart{
			ID:      "bandwidth",
			Options: Options{"Bandwidth", "kilobits/s", "bandwidth", "", raw.Area},
			Dimensions: Dimensions{
				Dimension{"resp_length", "received", raw.Incremental},
				Dimension{"bytes_sent", "sent", raw.Incremental},
			},
		},
		Chart{
			ID:      "response_time",
			Options: Options{"Processing Time", "milliseconds", "timings", "", raw.Area},
			Dimensions: Dimensions{
				Dimension{"resp_time_min", "min", raw.Incremental, 1, 1000},
				Dimension{"resp_time_max", "max", raw.Incremental, 1, 1000},
				Dimension{"resp_time_avg", "avg", raw.Incremental, 1, 1000},
			},
		},
		Chart{
			ID:      "response_time_upstream",
			Options: Options{"Processing Time Upstream", "milliseconds", "timings", "", raw.Area},
			Dimensions: Dimensions{
				Dimension{"resp_time_upstream_min", "min", raw.Incremental, 1, 1000},
				Dimension{"resp_time_upstream_max", "max", raw.Incremental, 1, 1000},
				Dimension{"resp_time_upstream_avg", "avg", raw.Incremental, 1, 1000},
			},
		},
		Chart{
			ID:      "clients",
			Options: Options{"Current Poll Unique Client IPs", "unique ips", "clients", "", raw.Stacked},
			Dimensions: Dimensions{
				Dimension{"unique_cur_ipv4", "ipv4", raw.Incremental},
				Dimension{"unique_cur_ipv6", "ipv6", raw.Incremental},
			},
		},
		Chart{
			ID:      "clients_all",
			Options: Options{"All Time Unique Client IPs", "unique ips", "clients", "", raw.Stacked},
			Dimensions: Dimensions{
				Dimension{"unique_tot_ipv4", "ipv4"},
				Dimension{"unique_tot_ipv6", "ipv6"},
			},
		},
		Chart{
			ID:      "http_method",
			Options: Options{"Requests Per HTTP Method", "requests/s", "http methods", "", raw.Stacked},
			Dimensions: Dimensions{
				Dimension{"GET", "", raw.Incremental},
			},
		},
		Chart{
			ID:         "http_version",
			Options:    Options{"Requests Per HTTP Version", "requests/s", "http versions", "", raw.Stacked},
			Dimensions: Dimensions{},
		},
		Chart{
			ID:         "requests_per_ipproto",
			Options:    Options{"Requests Per IP Protocol", "requests/s", "ip protocols", "", raw.Stacked},
			Dimensions: Dimensions{},
		},
		Chart{
			ID:      "requests_per_url",
			Options: Options{"Requests Per Url", "requests/s", "urls", "", raw.Stacked},
			Dimensions: Dimensions{
				Dimension{"url_other", "other", raw.Incremental},
			},
		},
		Chart{
			ID:      "requests_per_user_defined",
			Options: Options{"Requests Per User Defined Pattern", "requests/s", "user defined", "", raw.Stacked},
			Dimensions: Dimensions{
				Dimension{"user_pattern_other", "other", raw.Incremental},
			},
		},
	},
}

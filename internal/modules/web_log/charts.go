package web_log

import (
	"github.com/l2isbad/go.d.plugin/internal/pkg/charts"
)

type (
	Chart      = charts.Chart
	Options    = charts.Options
	Dimensions = charts.Dimensions
)

var (
	chartRespStatuses = Chart{
		ID: "response_statuses",
		Options: Options{
			Title: "Response Statuses", Units: "requests/s", Family: "responses", Type: charts.Stacked},
		Dimensions: Dimensions{
			{ID: "successful_requests", Name: "success", Algorithm: charts.Incremental},
			{ID: "server_errors", Name: "error", Algorithm: charts.Incremental},
			{ID: "redirects", Name: "redirect", Algorithm: charts.Incremental},
			{ID: "bad_requests", Name: "bad", Algorithm: charts.Incremental},
			{ID: "other_requests", Name: "other", Algorithm: charts.Incremental},
		},
	}
	chartRespCodes = Chart{
		ID: "response_codes",
		Options: Options{
			Title: "Response Codes", Units: "requests/s", Family: "responses", Type: charts.Stacked},
		Dimensions: Dimensions{
			{ID: "2xx", Algorithm: charts.Incremental},
			{ID: "5xx", Algorithm: charts.Incremental},
			{ID: "3xx", Algorithm: charts.Incremental},
			{ID: "4xx", Algorithm: charts.Incremental},
			{ID: "1xx", Algorithm: charts.Incremental},
			{ID: "0xx", Algorithm: charts.Incremental},
			{ID: "unmatched", Algorithm: charts.Incremental},
		},
	}
	chartRespCodesDetailed = Chart{
		ID: "detailed_response_codes",
		Options: Options{
			Title: "Detailed Response Codes", Units: "requests/s", Family: "responses", Type: charts.Stacked},
	}
	chartBandwidth = Chart{
		ID: "bandwidth",
		Options: Options{
			Title: "chartBandwidth", Units: "kilobits/s", Family: "bandwidth", Type: charts.Area},
		Dimensions: Dimensions{
			{ID: "resp_length", Name: "received", Algorithm: charts.Incremental, Multiplier: 8, Divisor: 1000},
			{ID: "bytes_sent", Name: "sent", Algorithm: charts.Incremental, Multiplier: -8, Divisor: 1000},
		},
	}
	chartRespTime = Chart{
		ID: "response_time",
		Options: Options{
			Title: "Processing Time", Units: "milliseconds", Family: "timings", Type: charts.Area},
		Dimensions: Dimensions{
			{ID: "resp_time_min", Name: "min", Algorithm: charts.Incremental, Divisor: 1000},
			{ID: "resp_time_max", Name: "max", Algorithm: charts.Incremental, Divisor: 1000},
			{ID: "resp_time_avg", Name: "avg", Algorithm: charts.Incremental, Divisor: 1000},
		},
	}
	chartRespTimeHist = Chart{
		ID: "response_time_histogram",
		Options: Options{
			Title: "Processing Time Histogram", Units: "requests/s", Family: "timings"},
	}
	chartRespTimeUpstream = Chart{
		ID: "response_time_upstream",
		Options: Options{
			Title: "Processing Time Upstream", Units: "milliseconds", Family: "timings", Type: charts.Area},
		Dimensions: Dimensions{
			{ID: "resp_time_upstream_min", Name: "min", Algorithm: charts.Incremental, Divisor: 1000},
			{ID: "resp_time_upstream_max", Name: "max", Algorithm: charts.Incremental, Divisor: 1000},
			{ID: "resp_time_upstream_avg", Name: "avg", Algorithm: charts.Incremental, Divisor: 1000},
		},
	}
	chartRespTimeUpstreamHist = Chart{
		ID: "response_time_upstream_histogram",
		Options: Options{
			Title: "Processing Time Upstream Histogram", Units: "requests/s", Family: "timings"},
	}
	chartReqPerURL = Chart{
		ID: "requests_per_url",
		Options: Options{
			Title: "Requests Per Url", Units: "requests/s", Family: "urls", Type: charts.Stacked},
	}
	chartReqPerUserDef = Chart{
		ID: "requests_per_user_defined",
		Options: Options{
			Title: "Requests Per User Defined Pattern", Units: "requests/s", Family: "user defined", Type: charts.Stacked},
	}
	chartReqPerHTTPMethod = Chart{
		ID: "requests_per_http_method",
		Options: Options{
			Title: "Requests Per HTTP Method", Units: "requests/s", Family: "http methods", Type: charts.Stacked},
		Dimensions: Dimensions{
			{ID: "GET", Algorithm: charts.Incremental},
		},
	}
	chartReqPerHTTPVer = Chart{
		ID: "requests_per_http_version",
		Options: Options{
			Title: "Requests Per HTTP Version", Units: "requests/s", Family: "http versions", Type: charts.Stacked},
	}
	chartReqPerIPProto = Chart{
		ID: "requests_per_ip_proto",
		Options: Options{
			Title: "Requests Per IP Protocol", Units: "requests/s", Family: "ip protocols", Type: charts.Stacked},
		Dimensions: Dimensions{
			{ID: "req_ipv4", Name: "ipv4", Algorithm: charts.Incremental},
			{ID: "req_ipv6", Name: "ipv6", Algorithm: charts.Incremental},
		},
	}
	chartClientsCurr = Chart{
		ID: "clients_current",
		Options: Options{
			Title: "Current Poll Unique Client IPs", Units: "unique ips", Family: "clients", Type: charts.Stacked},
		Dimensions: Dimensions{
			{ID: "unique_cur_ipv4", Name: "ipv4", Algorithm: charts.Incremental},
			{ID: "unique_cur_ipv6", Name: "ipv6", Algorithm: charts.Incremental},
		},
	}
	chartClientsAll = Chart{
		ID: "clients_all_time",
		Options: Options{
			Title: "All Time Unique Client IPs", Units: "unique ips", Family: "clients", Type: charts.Stacked},
		Dimensions: Dimensions{
			{ID: "unique_all_ipv4", Name: "ipv4"},
			{ID: "unique_all_ipv6", Name: "ipv6"},
		},
	}
)

func chartRespCodesDetailedPerFam() []Chart {
	return []Chart{
		{
			ID:      chartRespCodesDetailed.ID + "_1xx",
			Options: Options{Title: "Detailed Response Codes 1xx", Units: "requests/s", Family: "responses", Type: charts.Stacked},
		},
		{
			ID:      chartRespCodesDetailed.ID + "_2xx",
			Options: Options{Title: "Detailed Response Codes 2xx", Units: "requests/s", Family: "responses", Type: charts.Stacked},
		},
		{
			ID:      chartRespCodesDetailed.ID + "_3xx",
			Options: Options{Title: "Detailed Response Codes 3xx", Units: "requests/s", Family: "responses", Type: charts.Stacked},
		},
		{
			ID:      chartRespCodesDetailed.ID + "_4xx",
			Options: Options{Title: "Detailed Response Codes 4xx", Units: "requests/s", Family: "responses", Type: charts.Stacked},
		},
		{
			ID:      chartRespCodesDetailed.ID + "_5xx",
			Options: Options{Title: "Detailed Response Codes 5xx", Units: "requests/s", Family: "responses", Type: charts.Stacked},
		},
		{
			ID:      chartRespCodesDetailed.ID + "_other",
			Options: Options{Title: "Detailed Response Codes Other", Units: "requests/s", Family: "responses", Type: charts.Stacked},
		},
	}
}

func chartPerCategoryStats(id string) []Chart {
	return []Chart{
		{
			ID: chartRespCodesDetailed.ID + "_" + id,
			Options: Options{
				Title: "Detailed Response Codes", Units: "requests/s", Family: id,
				Context: "web_log.url_detailed_response_codes", Type: charts.Stacked},
		},
		{
			ID: chartBandwidth.ID + "_" + id,
			Options: Options{
				Title: "chartBandwidth", Units: "kilobits/s", Family: id, Context: "web_log.url_bandwidth", Type: charts.Area},
			Dimensions: Dimensions{
				{ID: id + "_resp_length", Name: "received", Algorithm: charts.Incremental, Multiplier: 8, Divisor: 1000},
				{ID: id + "_bytes_sent", Name: "sent", Algorithm: charts.Incremental, Multiplier: -8, Divisor: 1000},
			},
		},
		{
			ID: chartRespTime.ID + "_" + id,
			Options: Options{
				Title: "Processing Time", Units: "milliseconds", Family: id, Context: "web_log.url_response_time", Type: charts.Area},
			Dimensions: Dimensions{
				{ID: id + "_resp_time_min", Name: "min", Algorithm: charts.Incremental, Divisor: 1000},
				{ID: id + "_resp_time_max", Name: "max", Algorithm: charts.Incremental, Divisor: 1000},
				{ID: id + "_resp_time_avg", Name: "avg", Algorithm: charts.Incremental, Divisor: 1000},
			},
		},
	}
}

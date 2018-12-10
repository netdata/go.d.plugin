package charts

import (
	"github.com/netdata/go.d.plugin/modules"
)

type (
	// Chart is an alias for modules.Chart
	Chart = modules.Chart
	// Dims is an alias for modules.Dims
	Dims = modules.Dims
)

var (
	stacked     = modules.Stacked
	area        = modules.Area
	incremental = modules.Incremental
)

var (
	ResponseStatuses = Chart{
		ID:    "response_statuses",
		Title: "Response Statuses", Units: "requests/s", Fam: "responses", Type: stacked,
		Dims: Dims{
			{ID: "successful_requests", Name: "success", Algo: incremental},
			{ID: "server_errors", Name: "error", Algo: incremental},
			{ID: "redirects", Name: "redirect", Algo: incremental},
			{ID: "bad_requests", Name: "bad", Algo: incremental},
			{ID: "other_requests", Name: "other", Algo: incremental},
		},
	}
	ResponseCodes = Chart{
		ID:    "response_codes",
		Title: "Response Codes", Units: "requests/s", Fam: "responses", Type: stacked,
		Dims: Dims{
			{ID: "2xx", Algo: incremental},
			{ID: "5xx", Algo: incremental},
			{ID: "3xx", Algo: incremental},
			{ID: "4xx", Algo: incremental},
			{ID: "1xx", Algo: incremental},
			{ID: "0xx", Algo: incremental},
			{ID: "unmatched", Algo: incremental},
		},
	}
	ResponseCodesDetailed = Chart{
		ID:    "detailed_response_codes",
		Title: "Detailed Response Codes", Units: "requests/s", Fam: "responses", Type: stacked,
	}
	Bandwidth = Chart{
		ID:    "bandwidth",
		Title: "Bandwidth", Units: "kilobits/s", Fam: "bandwidth", Type: area,
		Dims: Dims{
			{ID: "resp_length", Name: "received", Algo: incremental, Mul: 8, Div: 1000},
			{ID: "bytes_sent", Name: "sent", Algo: incremental, Mul: -8, Div: 1000},
		},
	}
	ResponseTime = Chart{
		ID:    "response_time",
		Title: "Processing Time", Units: "milliseconds", Fam: "timings", Type: area,
		Dims: Dims{
			{ID: "resp_time_min", Name: "min", Algo: incremental, Div: 1000},
			{ID: "resp_time_max", Name: "max", Algo: incremental, Div: 1000},
			{ID: "resp_time_avg", Name: "avg", Algo: incremental, Div: 1000},
		},
	}
	ResponseTimeHistogram = Chart{
		ID:    "response_time_histogram",
		Title: "Processing Time Histogram", Units: "requests/s", Fam: "timings",
	}
	ResponseTimeUpstream = Chart{
		ID:    "response_time_upstream",
		Title: "Processing Time Upstream", Units: "milliseconds", Fam: "timings", Type: area,
		Dims: Dims{
			{ID: "resp_time_upstream_min", Name: "min", Algo: incremental, Div: 1000},
			{ID: "resp_time_upstream_max", Name: "max", Algo: incremental, Div: 1000},
			{ID: "resp_time_upstream_avg", Name: "avg", Algo: incremental, Div: 1000},
		},
	}
	ResponseTimeUpstreamHistogram = Chart{
		ID:    "response_time_upstream_histogram",
		Title: "Processing Time Upstream Histogram", Units: "requests/s", Fam: "timings",
	}
	RequestsPerURL = Chart{
		ID:    "requests_per_url",
		Title: "Requests Per Url", Units: "requests/s", Fam: "urls", Type: stacked,
	}
	RequestsPerUserDefined = Chart{
		ID:    "requests_per_user_defined",
		Title: "Requests Per User Defined Pattern", Units: "requests/s", Fam: "user defined", Type: stacked,
	}
	RequestsPerHTTPMethod = Chart{
		ID:    "requests_per_http_method",
		Title: "Requests Per HTTP Method", Units: "requests/s", Fam: "http methods", Type: stacked,
		Dims: Dims{
			{ID: "GET", Algo: incremental},
		},
	}
	RequestsPerHTTPVersion = Chart{
		ID:    "requests_per_http_version",
		Title: "Requests Per HTTP Version", Units: "requests/s", Fam: "http versions", Type: stacked,
	}
	RequestsPerIPProto = Chart{
		ID:    "requests_per_ip_proto",
		Title: "Requests Per IP Protocol", Units: "requests/s", Fam: "ip protocols", Type: stacked,
		Dims: Dims{
			{ID: "req_ipv4", Name: "ipv4", Algo: incremental},
			{ID: "req_ipv6", Name: "ipv6", Algo: incremental},
		},
	}
	CurrentPollIPs = Chart{
		ID:    "clients_current",
		Title: "Current Poll Unique Client IPs", Units: "unique ips", Fam: "clients", Type: stacked,
		Dims: Dims{
			{ID: "unique_cur_ipv4", Name: "ipv4", Algo: incremental},
			{ID: "unique_cur_ipv6", Name: "ipv6", Algo: incremental},
		},
	}
	AllTimeIPs = Chart{
		ID:    "clients_all_time",
		Title: "All Time Unique Client IPs", Units: "unique ips", Fam: "clients", Type: stacked,
		Dims: Dims{
			{ID: "unique_all_ipv4", Name: "ipv4"},
			{ID: "unique_all_ipv6", Name: "ipv6"},
		},
	}
)

func ResponseCodesDetailedPerFamily() []Chart {
	return []Chart{
		{
			ID:    ResponseCodesDetailed.ID + "_1xx",
			Title: "Detailed Response Codes 1xx", Units: "requests/s", Fam: "responses", Type: stacked,
		},
		{
			ID:    ResponseCodesDetailed.ID + "_2xx",
			Title: "Detailed Response Codes 2xx", Units: "requests/s", Fam: "responses", Type: stacked,
		},
		{
			ID:    ResponseCodesDetailed.ID + "_3xx",
			Title: "Detailed Response Codes 3xx", Units: "requests/s", Fam: "responses", Type: stacked,
		},
		{
			ID:    ResponseCodesDetailed.ID + "_4xx",
			Title: "Detailed Response Codes 4xx", Units: "requests/s", Fam: "responses", Type: stacked,
		},
		{
			ID:    ResponseCodesDetailed.ID + "_5xx",
			Title: "Detailed Response Codes 5xx", Units: "requests/s", Fam: "responses", Type: stacked,
		},
		{
			ID:    ResponseCodesDetailed.ID + "_other",
			Title: "Detailed Response Codes Other", Units: "requests/s", Fam: "responses", Type: stacked,
		},
	}
}

func PerCategoryStats(id string) []Chart {
	return []Chart{
		{
			ID:    ResponseCodesDetailed.ID + "_" + id,
			Title: "Detailed Response Codes", Units: "requests/s", Fam: id, Ctx: "web_log.url_detailed_response_codes",
			Type: stacked,
		},
		{
			ID:    Bandwidth.ID + "_" + id,
			Title: "Bandwidth", Units: "kilobits/s", Fam: id, Ctx: "web_log.url_bandwidth",
			Type: area,
			Dims: Dims{
				{ID: id + "_resp_length", Name: "received", Algo: incremental, Mul: 8, Div: 1000},
				{ID: id + "_bytes_sent", Name: "sent", Algo: incremental, Mul: -8, Div: 1000},
			},
		},
		{
			ID:    ResponseTime.ID + "_" + id,
			Title: "Processing Time", Units: "milliseconds", Fam: id, Ctx: "web_log.url_response_time",
			Type: area,
			Dims: Dims{
				{ID: id + "_resp_time_min", Name: "min", Algo: incremental, Div: 1000},
				{ID: id + "_resp_time_max", Name: "max", Algo: incremental, Div: 1000},
				{ID: id + "_resp_time_avg", Name: "avg", Algo: incremental, Div: 1000},
			},
		},
	}
}

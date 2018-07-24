package charts

import "github.com/l2isbad/go.d.plugin/internal/pkg/charts/raw"

type (
	Chart      = raw.Chart
	Options    = raw.Options
	Dimensions = raw.Dimensions
	Dimension  = raw.Dimension
)

var (
	RespStatuses = Chart{
		ID:      "response_statuses",
		Options: Options{"Response Statuses", "requests/s", "responses", "", raw.Stacked},
		Dimensions: Dimensions{
			Dimension{"successful_requests", "success", raw.Incremental},
			Dimension{"server_errors", "error", raw.Incremental},
			Dimension{"redirects", "redirect", raw.Incremental},
			Dimension{"bad_requests", "bad", raw.Incremental},
			Dimension{"other_requests", "other", raw.Incremental},
		},
	}
	RespCodes = Chart{
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
	}
	RespCodesDetailed = Chart{
		ID:      "detailed_response_codes",
		Options: Options{"Detailed Response Codes", "requests/s", "responses", "", raw.Stacked},
	}
	Bandwidth = Chart{
		ID:      "bandwidth",
		Options: Options{"Bandwidth", "kilobits/s", "bandwidth", "", raw.Area},
		Dimensions: Dimensions{
			Dimension{"resp_length", "received", raw.Incremental, 8, 1000},
			Dimension{"bytes_sent", "sent", raw.Incremental, -8, 1000},
		},
	}
	RespTime = Chart{
		ID:      "response_time",
		Options: Options{"Processing Time", "milliseconds", "timings", "", raw.Area},
		Dimensions: Dimensions{
			Dimension{"resp_time_min", "min", raw.Incremental, 1, 1000},
			Dimension{"resp_time_max", "max", raw.Incremental, 1, 1000},
			Dimension{"resp_time_avg", "avg", raw.Incremental, 1, 1000},
		},
	}
	RespTimeHist = Chart{
		ID:      "response_time_histogram",
		Options: Options{"Processing Time Histogram", "requests/s", "timings"},
	}
	RespTimeUpstream = Chart{
		ID:      "response_time_upstream",
		Options: Options{"Processing Time Upstream", "milliseconds", "timings", "", raw.Area},
		Dimensions: Dimensions{
			Dimension{"resp_time_upstream_min", "min", raw.Incremental, 1, 1000},
			Dimension{"resp_time_upstream_max", "max", raw.Incremental, 1, 1000},
			Dimension{"resp_time_upstream_avg", "avg", raw.Incremental, 1, 1000},
		},
	}
	RespTimeUpstreamHist = Chart{
		ID:      "response_time_upstream_histogram",
		Options: Options{"Processing Time Upstream Histogram", "requests/s", "timings"},
	}
	ReqPerURL = Chart{
		ID:      "requests_per_url",
		Options: Options{"Requests Per Url", "requests/s", "urls", "", raw.Stacked},
	}
	ReqPerUserDef = Chart{
		ID:      "requests_per_user_defined",
		Options: Options{"Requests Per User Defined Pattern", "requests/s", "user defined", "", raw.Stacked},
	}
	ReqPerHTTPMethod = Chart{
		ID:      "requests_per_http_method",
		Options: Options{"Requests Per HTTP Method", "requests/s", "http methods", "", raw.Stacked},
		Dimensions: Dimensions{
			Dimension{"GET", "", raw.Incremental},
		},
	}
	ReqPerHTTPVer = Chart{
		ID:      "requests_per_http_version",
		Options: Options{"Requests Per HTTP Version", "requests/s", "http versions", "", raw.Stacked},
	}
	ReqPerIPProto = Chart{
		ID:      "requests_per_ip_proto",
		Options: Options{"Requests Per IP Protocol", "requests/s", "ip protocols", "", raw.Stacked},
		Dimensions: Dimensions{
			Dimension{"req_ipv4", "ipv4", raw.Incremental},
			Dimension{"req_ipv6", "ipv6", raw.Incremental},
		},
	}
	ClientsCurr = Chart{
		ID:      "clients_current",
		Options: Options{"Current Poll Unique Client IPs", "unique ips", "clients", "", raw.Stacked},
		Dimensions: Dimensions{
			Dimension{"unique_cur_ipv4", "ipv4", raw.Incremental},
			Dimension{"unique_cur_ipv6", "ipv6", raw.Incremental},
		},
	}
	ClientsAll = Chart{
		ID:      "clients_all_time",
		Options: Options{"All Time Unique Client IPs", "unique ips", "clients", "", raw.Stacked},
		Dimensions: Dimensions{
			Dimension{"unique_all_ipv4", "ipv4"},
			Dimension{"unique_all_ipv6", "ipv6"},
		},
	}
)

func RespCodesDetailedPerFam() []*Chart {
	return []*Chart{
		raw.NewChart(
			RespCodesDetailed.ID+"_1xx",
			Options{"Detailed Response Codes 1xx", "requests/s", "responses", "", raw.Stacked},
		),
		raw.NewChart(
			RespCodesDetailed.ID+"_2xx",
			Options{"Detailed Response Codes 2xx", "requests/s", "responses", "", raw.Stacked},
		),
		raw.NewChart(
			RespCodesDetailed.ID+"_3xx",
			Options{"Detailed Response Codes 3xx", "requests/s", "responses", "", raw.Stacked},
		),
		raw.NewChart(
			RespCodesDetailed.ID+"_4xx",
			Options{"Detailed Response Codes 4xx", "requests/s", "responses", "", raw.Stacked},
		),
		raw.NewChart(
			RespCodesDetailed.ID+"_5xx",
			Options{"Detailed Response Codes 5xx", "requests/s", "responses", "", raw.Stacked},
		),
		raw.NewChart(
			RespCodesDetailed.ID+"_other",
			Options{"Detailed Response Codes Other", "requests/s", "responses", "", raw.Stacked},
		),
	}
}

func PerCategoryStats(id string) []*Chart {
	return []*Chart{
		raw.NewChart(
			RespCodesDetailed.ID+"_"+id,
			Options{"Detailed Response Codes", "requests/s", id, "web_log.url_detailed_response_codes", raw.Stacked},
		),
		raw.NewChart(
			Bandwidth.ID+"_"+id,
			Options{"Bandwidth", "kilobits/s", id, "web_log.url_bandwidth", raw.Area},
			Dimension{id + "_resp_length", "received", raw.Incremental, 8, 1000},
			Dimension{id + "_bytes_sent", "sent", raw.Incremental, -8, 1000},
		),
		raw.NewChart(
			RespTime.ID+"_"+id,
			Options{"Processing Time", "milliseconds", id, "web_log.url_response_time", raw.Area},
			Dimension{id + "_resp_time_min", "min", raw.Incremental, 1, 1000},
			Dimension{id + "_resp_time_max", "max", raw.Incremental, 1, 1000},
			Dimension{id + "_resp_time_avg", "avg", raw.Incremental, 1, 1000},
		),
	}
}

package weblog

//
//import (
//	"github.com/l2isbad/go.d.plugin/pkg/charts"
//)
//
//type (
//	Chart = charts.Chart
//	Opts  = charts.Opts
//	Dims  = charts.Dims
//)
//
//var (
//	chartRespStatuses = Chart{
//		ID: "response_statuses",
//		Opts: Opts{
//			Title: "Response Statuses", Units: "requests/s", Fam: "responses", Type: charts.Stacked},
//		Dims: Dims{
//			{ID: "successful_requests", Name: "success", Algo: charts.Incremental},
//			{ID: "server_errors", Name: "error", Algo: charts.Incremental},
//			{ID: "redirects", Name: "redirect", Algo: charts.Incremental},
//			{ID: "bad_requests", Name: "bad", Algo: charts.Incremental},
//			{ID: "other_requests", Name: "other", Algo: charts.Incremental},
//		},
//	}
//	chartRespCodes = Chart{
//		ID: "response_codes",
//		Opts: Opts{
//			Title: "Response Codes", Units: "requests/s", Fam: "responses", Type: charts.Stacked},
//		Dims: Dims{
//			{ID: "2xx", Algo: charts.Incremental},
//			{ID: "5xx", Algo: charts.Incremental},
//			{ID: "3xx", Algo: charts.Incremental},
//			{ID: "4xx", Algo: charts.Incremental},
//			{ID: "1xx", Algo: charts.Incremental},
//			{ID: "0xx", Algo: charts.Incremental},
//			{ID: "unmatched", Algo: charts.Incremental},
//		},
//	}
//	chartRespCodesDetailed = Chart{
//		ID: "detailed_response_codes",
//		Opts: Opts{
//			Title: "Detailed Response Codes", Units: "requests/s", Fam: "responses", Type: charts.Stacked},
//	}
//	chartBandwidth = Chart{
//		ID: "bandwidth",
//		Opts: Opts{
//			Title: "chartBandwidth", Units: "kilobits/s", Fam: "bandwidth", Type: charts.Area},
//		Dims: Dims{
//			{ID: "resp_length", Name: "received", Algo: charts.Incremental, Mul: 8, Div: 1000},
//			{ID: "bytes_sent", Name: "sent", Algo: charts.Incremental, Mul: -8, Div: 1000},
//		},
//	}
//	chartRespTime = Chart{
//		ID: "response_time",
//		Opts: Opts{
//			Title: "Processing Time", Units: "milliseconds", Fam: "timings", Type: charts.Area},
//		Dims: Dims{
//			{ID: "resp_time_min", Name: "min", Algo: charts.Incremental, Div: 1000},
//			{ID: "resp_time_max", Name: "max", Algo: charts.Incremental, Div: 1000},
//			{ID: "resp_time_avg", Name: "avg", Algo: charts.Incremental, Div: 1000},
//		},
//	}
//	chartRespTimeHist = Chart{
//		ID: "response_time_histogram",
//		Opts: Opts{
//			Title: "Processing Time Histogram", Units: "requests/s", Fam: "timings"},
//	}
//	chartRespTimeUpstream = Chart{
//		ID: "response_time_upstream",
//		Opts: Opts{
//			Title: "Processing Time Upstream", Units: "milliseconds", Fam: "timings", Type: charts.Area},
//		Dims: Dims{
//			{ID: "resp_time_upstream_min", Name: "min", Algo: charts.Incremental, Div: 1000},
//			{ID: "resp_time_upstream_max", Name: "max", Algo: charts.Incremental, Div: 1000},
//			{ID: "resp_time_upstream_avg", Name: "avg", Algo: charts.Incremental, Div: 1000},
//		},
//	}
//	chartRespTimeUpstreamHist = Chart{
//		ID: "response_time_upstream_histogram",
//		Opts: Opts{
//			Title: "Processing Time Upstream Histogram", Units: "requests/s", Fam: "timings"},
//	}
//	chartReqPerURL = Chart{
//		ID: "requests_per_url",
//		Opts: Opts{
//			Title: "Requests Per Url", Units: "requests/s", Fam: "urls", Type: charts.Stacked},
//	}
//	chartReqPerUserDef = Chart{
//		ID: "requests_per_user_defined",
//		Opts: Opts{
//			Title: "Requests Per User Defined Pattern", Units: "requests/s", Fam: "user defined", Type: charts.Stacked},
//	}
//	chartReqPerHTTPMethod = Chart{
//		ID: "requests_per_http_method",
//		Opts: Opts{
//			Title: "Requests Per HTTP Method", Units: "requests/s", Fam: "http methods", Type: charts.Stacked},
//		Dims: Dims{
//			{ID: "GET", Algo: charts.Incremental},
//		},
//	}
//	chartReqPerHTTPVer = Chart{
//		ID: "requests_per_http_version",
//		Opts: Opts{
//			Title: "Requests Per HTTP Version", Units: "requests/s", Fam: "http versions", Type: charts.Stacked},
//	}
//	chartReqPerIPProto = Chart{
//		ID: "requests_per_ip_proto",
//		Opts: Opts{
//			Title: "Requests Per IP Protocol", Units: "requests/s", Fam: "ip protocols", Type: charts.Stacked},
//		Dims: Dims{
//			{ID: "req_ipv4", Name: "ipv4", Algo: charts.Incremental},
//			{ID: "req_ipv6", Name: "ipv6", Algo: charts.Incremental},
//		},
//	}
//	chartClientsCurr = Chart{
//		ID: "clients_current",
//		Opts: Opts{
//			Title: "Current Poll Unique Client IPs", Units: "unique ips", Fam: "clients", Type: charts.Stacked},
//		Dims: Dims{
//			{ID: "unique_cur_ipv4", Name: "ipv4", Algo: charts.Incremental},
//			{ID: "unique_cur_ipv6", Name: "ipv6", Algo: charts.Incremental},
//		},
//	}
//	chartClientsAll = Chart{
//		ID: "clients_all_time",
//		Opts: Opts{
//			Title: "All Time Unique Client IPs", Units: "unique ips", Fam: "clients", Type: charts.Stacked},
//		Dims: Dims{
//			{ID: "unique_all_ipv4", Name: "ipv4"},
//			{ID: "unique_all_ipv6", Name: "ipv6"},
//		},
//	}
//)
//
//func chartRespCodesDetailedPerFam() []Chart {
//	return []Chart{
//		{
//			ID:   chartRespCodesDetailed.ID + "_1xx",
//			Opts: Opts{Title: "Detailed Response Codes 1xx", Units: "requests/s", Fam: "responses", Type: charts.Stacked},
//		},
//		{
//			ID:   chartRespCodesDetailed.ID + "_2xx",
//			Opts: Opts{Title: "Detailed Response Codes 2xx", Units: "requests/s", Fam: "responses", Type: charts.Stacked},
//		},
//		{
//			ID:   chartRespCodesDetailed.ID + "_3xx",
//			Opts: Opts{Title: "Detailed Response Codes 3xx", Units: "requests/s", Fam: "responses", Type: charts.Stacked},
//		},
//		{
//			ID:   chartRespCodesDetailed.ID + "_4xx",
//			Opts: Opts{Title: "Detailed Response Codes 4xx", Units: "requests/s", Fam: "responses", Type: charts.Stacked},
//		},
//		{
//			ID:   chartRespCodesDetailed.ID + "_5xx",
//			Opts: Opts{Title: "Detailed Response Codes 5xx", Units: "requests/s", Fam: "responses", Type: charts.Stacked},
//		},
//		{
//			ID:   chartRespCodesDetailed.ID + "_other",
//			Opts: Opts{Title: "Detailed Response Codes Other", Units: "requests/s", Fam: "responses", Type: charts.Stacked},
//		},
//	}
//}
//
//func chartPerCategoryStats(id string) []Chart {
//	return []Chart{
//		{
//			ID: chartRespCodesDetailed.ID + "_" + id,
//			Opts: Opts{
//				Title: "Detailed Response Codes", Units: "requests/s", Fam: id,
//				Ctx: "weblog.url_detailed_response_codes", Type: charts.Stacked},
//		},
//		{
//			ID: chartBandwidth.ID + "_" + id,
//			Opts: Opts{
//				Title: "chartBandwidth", Units: "kilobits/s", Fam: id, Ctx: "weblog.url_bandwidth", Type: charts.Area},
//			Dims: Dims{
//				{ID: id + "_resp_length", Name: "received", Algo: charts.Incremental, Mul: 8, Div: 1000},
//				{ID: id + "_bytes_sent", Name: "sent", Algo: charts.Incremental, Mul: -8, Div: 1000},
//			},
//		},
//		{
//			ID: chartRespTime.ID + "_" + id,
//			Opts: Opts{
//				Title: "Processing Time", Units: "milliseconds", Fam: id, Ctx: "weblog.url_response_time", Type: charts.Area},
//			Dims: Dims{
//				{ID: id + "_resp_time_min", Name: "min", Algo: charts.Incremental, Div: 1000},
//				{ID: id + "_resp_time_max", Name: "max", Algo: charts.Incremental, Div: 1000},
//				{ID: id + "_resp_time_avg", Name: "avg", Algo: charts.Incremental, Div: 1000},
//			},
//		},
//	}
//}

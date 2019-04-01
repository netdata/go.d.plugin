package weblog

import (
	"github.com/netdata/go-orchestrator/module"
)

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Chart is an alias for module.Chart
	Chart = module.Chart
	// Dims is an alias for module.Dims
	Dims = module.Dims
	// Dim is an alias for module.Dim
	Dim = module.Dim
)

//
//// NOTE: inconsistency between contexts with python web_log
var (
	responseStatuses = Chart{
		ID:    "response_statuses",
		Title: "Response Statuses",
		Units: "requests/s",
		Fam:   "responses",
		Ctx:   "web_log.response_statuses",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "resp_successful", Name: "success", Algo: module.Incremental},
			{ID: "resp_redirect", Name: "redirect", Algo: module.Incremental},
			{ID: "resp_client_error", Name: "bad", Algo: module.Incremental},
			{ID: "resp_server_error", Name: "error", Algo: module.Incremental},
		},
	}

	responseCodes = Chart{
		ID:    "response_codes",
		Title: "Response Codes",
		Units: "requests/s",
		Fam:   "responses",
		Ctx:   "web_log.response_codes",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "resp_2xx", Algo: module.Incremental},
			{ID: "resp_5xx", Algo: module.Incremental},
			{ID: "resp_3xx", Algo: module.Incremental},
			{ID: "resp_4xx", Algo: module.Incremental},
			{ID: "resp_1xx", Algo: module.Incremental},
		},
	}

	responseCodesDetailed = Chart{
		ID:    "detailed_response_codes",
		Title: "Detailed Response Codes",
		Units: "requests/s",
		Fam:   "responses",
		Ctx:   "web_log.response_codes_detailed",
		Type:  module.Stacked,
	}

	bandwidth = Chart{
		ID:    "bandwidth",
		Title: "Bandwidth",
		Units: "kilobits/s",
		Fam:   "bandwidth",
		Ctx:   "web_log.bandwidth",
		Type:  module.Area,
		Dims: Dims{
			{ID: "bytes_received", Name: "received", Algo: module.Incremental, Mul: 8, Div: 1000},
			{ID: "bytes_sent", Name: "sent", Algo: module.Incremental, Mul: -8, Div: 1000},
		},
	}

	responseTime = Chart{
		ID:    "response_time",
		Title: "Processing Time",
		Units: "milliseconds",
		Fam:   "timings",
		Ctx:   "web_log.response_time",
		Type:  module.Line,
		Dims: Dims{
			{ID: "resp_time_min", Name: "min", Mul: 1000},
			{ID: "resp_time_max", Name: "max", Mul: 1000},
			{ID: "resp_time_avg", Name: "avg", Mul: 1000},
		},
	}

	responseTimeHistogram = Chart{
		ID:    "response_time_histogram",
		Title: "Processing Time Histogram",
		Units: "requests/s",
		Fam:   "timings",
		Ctx:   "web_log.response_time_histogram",
	}

	responseTimeUpstream = Chart{
		ID:    "response_time_upstream",
		Title: "Processing Time Upstream",
		Units: "milliseconds",
		Fam:   "timings",
		Ctx:   "web_log.response_time_upstream",
		Type:  module.Line,
		Dims: Dims{
			{ID: "resp_time_upstream_min", Name: "min", Mul: 1000},
			{ID: "resp_time_upstream_max", Name: "max", Mul: 1000},
			{ID: "resp_time_upstream_avg", Name: "avg", Mul: 1000},
		},
	}

	responseTimeUpstreamHistogram = Chart{
		ID:    "response_time_upstream_histogram",
		Title: "Processing Time Upstream Histogram",
		Units: "requests/s",
		Fam:   "timings",
		Ctx:   "web_log.response_time_upstream_histogram",
	}

	requestsPerURL = Chart{
		ID:    "requests_per_url",
		Title: "Requests Per Url",
		Units: "requests/s",
		Fam:   "urls",
		Ctx:   "web_log.requests_per_url",
		Type:  module.Stacked,
	}

	requestsPerUserDefined = Chart{
		ID:    "requests_per_user_defined",
		Title: "Requests Per User Defined Pattern",
		Units: "requests/s",
		Fam:   "user defined",
		Ctx:   "web_log.requests_per_user_defined",
		Type:  module.Stacked,
	}

	requestsPerHTTPMethod = Chart{
		ID:    "requests_per_http_method",
		Title: "Requests Per HTTP Method",
		Units: "requests/s",
		Fam:   "http methods",
		Ctx:   "web_log.requests_per_http_method",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "req_method_GET", Name: "GET", Algo: module.Incremental},
		},
	}
	requestsPerHTTPVersion = Chart{
		ID:    "requests_per_http_version",
		Title: "Requests Per HTTP Version",
		Units: "requests/s",
		Fam:   "http versions",
		Ctx:   "web_log.requests_per_http_version",
		Type:  module.Stacked,
	}

	requestsPerIPProto = Chart{
		ID:    "requests_per_ip_proto",
		Title: "Requests Per IP Protocol",
		Units: "requests/s",
		Fam:   "ip protocols",
		Ctx:   "web_log.requests_per_ip_proto",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "req_ipv4", Name: "ipv4", Algo: module.Incremental},
			{ID: "req_ipv6", Name: "ipv6", Algo: module.Incremental},
		},
	}
	requestsPerVhost = Chart{
		ID:    "requests_per_vhost",
		Title: "Requests Per Vhost",
		Units: "requests/s",
		Fam:   "vhost",
		Ctx:   "web_log.requests_per_vhost",
		Type:  module.Stacked,
	}

	currentPollIPs = Chart{
		ID:    "clients_current",
		Title: "Current Poll Unique Client IPs",
		Units: "unique ips",
		Fam:   "clients",
		Ctx:   "web_log.current_poll_ips",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "unique_current_poll_ipv4", Name: "ipv4", Algo: module.Absolute},
			{ID: "unique_current_poll_ipv6", Name: "ipv6", Algo: module.Absolute},
		},
	}
)

func responseCodesDetailedPerFamily() []*Chart {
	return []*Chart{
		{
			ID:    responseCodesDetailed.ID + "_1xx",
			Title: "Detailed Response Codes 1xx",
			Units: "requests/s",
			Fam:   "responses",
			Ctx:   "web_log.response_codes_detailed_1xx",
			Type:  module.Stacked,
		},
		{
			ID:    responseCodesDetailed.ID + "_2xx",
			Title: "Detailed Response Codes 2xx",
			Units: "requests/s",
			Fam:   "responses",
			Ctx:   "web_log.response_codes_detailed_2xx",
			Type:  module.Stacked,
		},
		{
			ID:    responseCodesDetailed.ID + "_3xx",
			Title: "Detailed Response Codes 3xx",
			Units: "requests/s",
			Fam:   "responses",
			Ctx:   "web_log.response_codes_detailed_3xx",
			Type:  module.Stacked,
		},
		{
			ID:    responseCodesDetailed.ID + "_4xx",
			Title: "Detailed Response Codes 4xx",
			Units: "requests/s",
			Fam:   "responses",
			Ctx:   "web_log.response_codes_detailed_4xx",
			Type:  module.Stacked,
		},
		{
			ID:    responseCodesDetailed.ID + "_5xx",
			Title: "Detailed Response Codes 5xx",
			Units: "requests/s",
			Fam:   "responses",
			Ctx:   "web_log.response_codes_detailed_5xx",
			Type:  module.Stacked,
		},
		{
			ID:    responseCodesDetailed.ID + "_other",
			Title: "Detailed Response Codes Other",
			Units: "requests/s",
			Fam:   "responses",
			Ctx:   "web_log.response_codes_detailed_other",
			Type:  module.Stacked,
		},
	}
}

func perCategoryStats(id string) []*Chart {
	return []*Chart{
		{
			ID:    responseCodesDetailed.ID + "_" + id,
			Title: "Detailed Response Codes",
			Units: "requests/s",
			Fam:   "url " + id,
			Ctx:   "web_log.response_codes_detailed_per_url",
			Type:  module.Stacked,
		},
		{
			ID:    bandwidth.ID + "_" + id,
			Title: "Bandwidth",
			Units: "kilobits/s",
			Fam:   "url " + id,
			Ctx:   "web_log.bandwidth_per_url",
			Type:  module.Area,
			Dims: Dims{
				{ID: id + "_resp_length", Name: "received", Algo: module.Incremental, Mul: 8, Div: 1000},
				{ID: id + "_bytes_sent", Name: "sent", Algo: module.Incremental, Mul: -8, Div: 1000},
			},
		},
		{
			ID:    responseTime.ID + "_" + id,
			Title: "Processing Time",
			Units: "milliseconds",
			Fam:   "url " + id,
			Ctx:   "web_log.response_time_per_url",
			Type:  module.Area,
			Dims: Dims{
				{ID: id + "_resp_time_min", Name: "min", Algo: module.Incremental, Div: 1000},
				{ID: id + "_resp_time_max", Name: "max", Algo: module.Incremental, Div: 1000},
				{ID: id + "_resp_time_avg", Name: "avg", Algo: module.Incremental, Div: 1000},
			},
		},
	}
}

package coredns

import "github.com/netdata/go-orchestrator/module"

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

var summaryCharts = Charts{
	{
		ID:    "requests_rate_by_status",
		Title: "Requests Rate By Status",
		Units: "requests/s",
		Fam:   "summary",
		Ctx:   "coredns.requests_rate_by_status",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "request_by_status_processed", Name: "processed", Algo: module.Incremental},
			{ID: "request_by_status_dropped", Name: "dropped", Algo: module.Incremental},
		},
	},
	{
		ID:    "panics_rate",
		Title: "Panics Rate",
		Units: "events/s",
		Fam:   "summary",
		Ctx:   "coredns.panics_rate",
		Dims: Dims{
			{ID: "panic_total", Name: "panics", Algo: module.Incremental},
		},
	},
}

var serverCharts = Charts{
	{
		ID:    "%s_requests_rate_by_status",
		Title: "Requests Rate By Status",
		Units: "requests/s",
		Fam:   "server %s",
		Ctx:   "coredns.requests_rate_by_status",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "%s_request_by_status_processed", Name: "processed", Algo: module.Incremental},
			{ID: "%s_request_by_status_dropped", Name: "dropped", Algo: module.Incremental},
		},
	},
	{
		ID:    "%s_requests_rate_by_proto",
		Title: "Requests Rate By Proto",
		Units: "requests/s",
		Fam:   "server %s",
		Ctx:   "coredns.requests_rate_by_proto",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "%s_request_by_proto_udp", Name: "udp", Algo: module.Incremental},
			{ID: "%s_request_by_proto_tcp", Name: "tcp", Algo: module.Incremental},
		},
	},
	{
		ID:    "%s_requests_rate_by_ip_family",
		Title: "Requests Rate By IP Family",
		Units: "requests/s",
		Fam:   "server %s",
		Ctx:   "coredns.requests_rate_by_ip_family",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "%s_request_by_ip_family_v4", Name: "v4", Algo: module.Incremental},
			{ID: "%s_request_by_ip_family_v6", Name: "v6", Algo: module.Incremental},
		},
	},
}

//var chartReqByTypeTotal = Chart{
//	ID:    "request_type_count_total",
//	Title: "The Total Number Of Requests By Type",
//	Units: "requests/s",
//	Fam:   "requests",
//	Ctx:   "coredns.request_type_count_total",
//}
//
//var chartRespByRcodeTotal = Chart{
//	ID:    "response_rcode_count_total",
//	Title: "The Total Number Of Responses By Rcode",
//	Units: "responses/s",
//	Fam:   "responses",
//	Ctx:   "coredns.request_type_count_total",
//}

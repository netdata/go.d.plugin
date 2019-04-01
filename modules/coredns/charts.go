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

var charts = Charts{
	{
		ID:    "panic_count_total",
		Title: "The Total Number Of Panics",
		Units: "events/s",
		Fam:   "panics",
		Ctx:   "coredns.panic_count_total",
		Dims: Dims{
			{ID: "panic_count_total", Name: "panics", Algo: module.Incremental},
		},
	},
	{
		ID:    "request_count_total",
		Title: "The Total Number Of Request",
		Units: "events/s",
		Fam:   "requests",
		Ctx:   "coredns.request_count_total",
		Dims: Dims{
			{ID: "request_count_total", Name: "requests", Algo: module.Incremental},
		},
	},
}

var chartReqByTypeTotal = Chart{
	ID:    "request_type_count_total",
	Title: "The Total Number Of Requests By Type",
	Units: "requests/s",
	Fam:   "requests",
	Ctx:   "coredns.request_type_count_total",
}

var chartRespByRcodeTotal = Chart{
	ID:    "response_rcode_count_total",
	Title: "The Total Number Of Responses By Rcode",
	Units: "responses/s",
	Fam:   "responses",
	Ctx:   "coredns.request_type_count_total",
}

package solr

import (
	"github.com/netdata/go.d.plugin/modules"
)

type (
	// Charts is an alias for modules.Charts
	Charts = modules.Charts
	// Chart is an alias for modules.Chart
	Chart = modules.Chart
	// Dims is an alias for modules.Dims
	Dims = modules.Dims
)

var charts = []*Chart{
	{
		ID:    "search_requests",
		Title: "Search Requests",
		Units: "requests/s",
		Ctx:   "solr.search_requests",
		Dims: Dims{
			{ID: "query_requests_count", Name: "requests", Algo: modules.Incremental},
		},
	},
	{
		ID:    "search_errors",
		Title: "Search Errors",
		Units: "errors/s",
		Ctx:   "solr.search_errors",
		Dims: Dims{
			{ID: "query_errors_count", Name: "errors", Algo: modules.Incremental},
		},
	},
	{
		ID:    "search_errors_by_type",
		Title: "Search Errors By Types",
		Units: "errors/s",
		Ctx:   "solr.search_errors_by_type",
		Dims: Dims{
			{ID: "query_clientErrors_count", Name: "client", Algo: modules.Incremental},
			{ID: "query_serverErrors_count", Name: "server", Algo: modules.Incremental},
			{ID: "query_timeouts_count", Name: "timeouts", Algo: modules.Incremental},
		},
	},
	{
		ID:    "search_requests_processing_time",
		Title: "Search Requests Processing Time",
		Units: "milliseconds",
		Ctx:   "solr.search_requests_processing_time",
		Dims: Dims{
			{ID: "query_totalTime_count", Name: "time", Algo: modules.Incremental},
		},
	},
	{
		ID:    "search_requests_timings",
		Title: "Search Requests Timings",
		Units: "milliseconds",
		Ctx:   "solr.search_requests_timings",
		Dims: Dims{
			{ID: "query_requestTimes_min_ms", Name: "min", Div: 1000000},
			{ID: "query_requestTimes_median_ms", Name: "median", Div: 1000000},
			{ID: "query_requestTimes_mean_ms", Name: "mean", Div: 1000000},
			{ID: "query_requestTimes_max_ms", Name: "max", Div: 1000000},
		},
	},
	{
		ID:    "search_requests_duration",
		Title: "Search Requests Duration",
		Units: "milliseconds",
		Ctx:   "solr.search_requests_duration",
		Dims: Dims{
			{ID: "query_requestTimes_p75_ms", Name: "p75", Div: 1000000},
			{ID: "query_requestTimes_p95_ms", Name: "p95", Div: 1000000},
			{ID: "query_requestTimes_p99_ms", Name: "p99", Div: 1000000},
			{ID: "query_requestTimes_p999_ms", Name: "p999", Div: 1000000},
		},
	},
}

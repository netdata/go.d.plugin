package solr

import (
	"github.com/netdata/go.d.plugin/modules"
)

type (
	// Charts is an alias for modules.Charts
	Charts = modules.Charts
	// Dims is an alias for modules.Dims
	Dims = modules.Dims
)

var charts = Charts{
	{
		ID:    "search_requests",
		Title: "Search Requests",
		Units: "requests/s",
		Ctx:   "solr.search_requests",
		Dims: Dims{
			{ID: "query_requests_count", Name: "search", Algo: modules.Incremental},
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
		Title: "Search Errors By Type",
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
	{
		ID:    "update_requests",
		Title: "Update Requests",
		Units: "requests/s",
		Ctx:   "solr.update_requests",
		Dims: Dims{
			{ID: "update_requests_count", Name: "update", Algo: modules.Incremental},
		},
	},
	{
		ID:    "update_errors",
		Title: "Update Errors",
		Units: "errors/s",
		Ctx:   "solr.update_errors",
		Dims: Dims{
			{ID: "update_errors_count", Name: "errors", Algo: modules.Incremental},
		},
	},
	{
		ID:    "update_errors_by_type",
		Title: "Update Errors By Type",
		Units: "errors/s",
		Ctx:   "solr.update_errors_by_type",
		Dims: Dims{
			{ID: "update_clientErrors_count", Name: "client", Algo: modules.Incremental},
			{ID: "update_serverErrors_count", Name: "server", Algo: modules.Incremental},
			{ID: "update_timeouts_count", Name: "timeouts", Algo: modules.Incremental},
		},
	},
	{
		ID:    "update_requests_processing_time",
		Title: "Update Requests Processing Time",
		Units: "milliseconds",
		Ctx:   "solr.update_requests_processing_time",
		Dims: Dims{
			{ID: "update_totalTime_count", Name: "time", Algo: modules.Incremental},
		},
	},
	{
		ID:    "update_requests_timings",
		Title: "Update Requests Timings",
		Units: "milliseconds",
		Ctx:   "solr.update_requests_timings",
		Dims: Dims{
			{ID: "update_requestTimes_min_ms", Name: "min", Div: 1000000},
			{ID: "update_requestTimes_median_ms", Name: "median", Div: 1000000},
			{ID: "update_requestTimes_mean_ms", Name: "mean", Div: 1000000},
			{ID: "update_requestTimes_max_ms", Name: "max", Div: 1000000},
		},
	},
	{
		ID:    "update_requests_duration",
		Title: "Update Requests Duration",
		Units: "milliseconds",
		Ctx:   "solr.update_requests_duration",
		Dims: Dims{
			{ID: "update_requestTimes_p75_ms", Name: "p75", Div: 1000000},
			{ID: "update_requestTimes_p95_ms", Name: "p95", Div: 1000000},
			{ID: "update_requestTimes_p99_ms", Name: "p99", Div: 1000000},
			{ID: "update_requestTimes_p999_ms", Name: "p999", Div: 1000000},
		},
	},
}

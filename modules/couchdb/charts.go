package couchdb

import (
	"github.com/netdata/go.d.plugin/agent/module"
)

type (
	Charts = module.Charts
	Dims   = module.Dims
	Vars   = module.Vars
)

var dbActivityCharts = Charts{
	{
		ID:    "activity",
		Title: "Overall Activity",
		Units: "requests/s",
		Fam:   "dbactivity",
		Ctx:   "couchdb.activity",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "couchdb_database_reads", Name: "DB reads", Algo: module.Incremental},
			{ID: "couchdb_database_writes", Name: "DB writes", Algo: module.Incremental},
			{ID: "couchdb_httpd_view_reads", Name: "View reads", Algo: module.Incremental},
		},
	},
}

var httpTrafficBreakdownCharts = Charts{
	{
		ID:    "request_methods",
		Title: "HTTP request methods",
		Units: "requests/s",
		Fam:   "httptraffic",
		Ctx:   "couchdb.request_methods",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "couchdb_httpd_request_methods_COPY", Name: "COPY", Algo: module.Incremental},
			{ID: "couchdb_httpd_request_methods_DELETE", Name: "DELETE", Algo: module.Incremental},
			{ID: "couchdb_httpd_request_methods_GET", Name: "GET", Algo: module.Incremental},
			{ID: "couchdb_httpd_request_methods_HEAD", Name: "HEAD", Algo: module.Incremental},
			{ID: "couchdb_httpd_request_methods_OPTIONS", Name: "OPTIONS", Algo: module.Incremental},
			{ID: "couchdb_httpd_request_methods_POST", Name: "POST", Algo: module.Incremental},
			{ID: "couchdb_httpd_request_methods_PUT", Name: "PUT", Algo: module.Incremental},
		},
	},
	{
		ID:    "response_codes",
		Title: "HTTP response status codes",
		Units: "responses/s",
		Fam:   "httptraffic",
		Ctx:   "couchdb.response_codes",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "couchdb_httpd_status_codes_200", Name: "200 OK", Algo: module.Incremental},
			{ID: "couchdb_httpd_status_codes_201", Name: "201 Created", Algo: module.Incremental},
			{ID: "couchdb_httpd_status_codes_202", Name: "202 Accepted", Algo: module.Incremental},
			{ID: "couchdb_httpd_status_codes_2xx", Name: "Other 2xx Success", Algo: module.Incremental},
			{ID: "couchdb_httpd_status_codes_3xx", Name: "3xx Redirection", Algo: module.Incremental},
			{ID: "couchdb_httpd_status_codes_4xx", Name: "4xx Client error", Algo: module.Incremental},
			{ID: "couchdb_httpd_status_codes_5xx", Name: "5xx Server error", Algo: module.Incremental},
		},
	},
}

var serverOperationsCharts = Charts{
	{
		ID:    "active_tasks",
		Title: "Active task breakdown",
		Units: "tasks",
		Fam:   "ops",
		Ctx:   "couchdb.active_tasks",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "active_tasks_indexer", Name: "Indexer"},
			{ID: "active_tasks_database_compaction", Name: "DB Compaction"},
			{ID: "active_tasks_replication", Name: "Replication"},
			{ID: "active_tasks_view_compaction", Name: "View Compaction"},
		},
	},
	{
		ID:    "replicator_jobs",
		Title: "Replicator job breakdown",
		Units: "jobs",
		Fam:   "ops",
		Ctx:   "couchdb.replicator_jobs",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "couch_replicator_jobs_running", Name: "Running"},
			{ID: "couch_replicator_jobs_pending", Name: "Pending"},
			{ID: "couch_replicator_jobs_crashed", Name: "Crashed"},
			{ID: "internal_replication_jobs", Name: "Internal replication jobs"},
		},
	},
	{
		ID:    "open_files",
		Title: "Open files",
		Units: "files",
		Fam:   "ops",
		Ctx:   "couchdb.open_files",
		Dims: Dims{
			{ID: "couchdb_open_os_files", Name: "# files"},
		},
	},
}

var erlangStatisticsCharts = Charts{
	{
		ID:    "erlang_memory",
		Title: "Erlang VM memory usage",
		Units: "B",
		Fam:   "erlang",
		Ctx:   "couchdb.erlang_vm_memory",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "memory_atom", Name: "atom"},
			{ID: "memory_binary", Name: "binaries"},
			{ID: "memory_code", Name: "code"},
			{ID: "memory_ets", Name: "ets"},
			{ID: "memory_processes", Name: "procs"},
			{ID: "memory_other", Name: "other"},
		},
	},
	{
		ID:    "erlang_reductions",
		Title: "Erlang reductions",
		Units: "count",
		Fam:   "erlang",
		Ctx:   "couchdb.reductions",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "reductions", Name: "reductions", Algo: module.Incremental},
		},
	},
	{
		ID:    "erlang_proc_counts",
		Title: "Process counts",
		Units: "count",
		Fam:   "erlang",
		Ctx:   "couchdb.proccounts",
		Dims: Dims{
			{ID: "os_proc_count", Name: "OS procs"},
			{ID: "process_count", Name: "erl procs"},
		},
	},
	{
		ID:    "erlang_peak_msg_queue",
		Title: "Peak message queue size",
		Units: "count",
		Fam:   "erlang",
		Ctx:   "couchdb.peakmsgqueue",
		Dims: Dims{
			{ID: "peak_msg_queue", Name: "peak size"},
		},
	},
}

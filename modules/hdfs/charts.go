package hdfs

import "github.com/netdata/go-orchestrator/module"

type (
	Charts = module.Charts
	Dims   = module.Dims
	Vars   = module.Vars
)

var jvmCharts = Charts{
	{
		ID:    "heap_memory",
		Title: "Heap Memory",
		Units: "MiB",
		Fam:   "jvm",
		Ctx:   "hdfs.heap_memory",
		Type:  module.Area,
		Dims: Dims{
			{ID: "jvm_mem_heap_committed", Name: "committed", Div: 1000},
			{ID: "jvm_mem_heap_used", Name: "used", Div: 1000},
		},
		Vars: Vars{
			{ID: "jvm_mem_heap_max"},
		},
	},
	{
		ID:    "gc_count_total",
		Title: "Total GC Count",
		Units: "events/s",
		Fam:   "jvm",
		Ctx:   "hdfs.gc_count_total",
		Dims: Dims{
			{ID: "jvm_gc_count", Name: "gc", Algo: module.Incremental},
		},
	},
	{
		ID:    "gc_time_total",
		Title: "Total GC Time",
		Units: "ms",
		Fam:   "jvm",
		Ctx:   "hdfs.gc_time_total",
		Dims: Dims{
			{ID: "jvm_gc_time_millis", Name: "time", Algo: module.Incremental},
		},
	},
	{
		ID:    "gc_threshold",
		Title: "Number of Times That the GC Threshold is Exceeded",
		Units: "events/s",
		Fam:   "jvm",
		Ctx:   "hdfs.gc_threshold",
		Dims: Dims{
			{ID: "jvm_gc_num_info_threshold_exceeded", Name: "info", Algo: module.Incremental},
			{ID: "jvm_gc_num_warn_threshold_exceeded", Name: "warn", Algo: module.Incremental},
		},
	},
	{
		ID:    "threads",
		Title: "Current Number of Threads",
		Units: "num",
		Fam:   "jvm",
		Ctx:   "hdfs.threads",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "jvm_threads_new", Name: "new"},
			{ID: "jvm_threads_runnable", Name: "runnable"},
			{ID: "jvm_threads_blocked", Name: "blocked"},
			{ID: "jvm_threads_waiting", Name: "waiting"},
			{ID: "jvm_threads_timed_waiting", Name: "timed_waiting"},
			{ID: "jvm_threads_terminated", Name: "terminated"},
		},
	},
	{
		ID:    "logs_total",
		Title: "Total Number of Logs",
		Units: "events/s",
		Fam:   "jvm",
		Ctx:   "hdfs.logs_total",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "jvm_log_info", Name: "info", Algo: module.Incremental},
			{ID: "jvm_log_error", Name: "error", Algo: module.Incremental},
			{ID: "jvm_log_warn", Name: "warn", Algo: module.Incremental},
			{ID: "jvm_log_fatal", Name: "fatal", Algo: module.Incremental},
		},
	},
}

var fsnCharts = Charts{
	{
		ID:    "capacity",
		Title: "Current Capacity Across All Datanodes",
		Units: "KiB",
		Fam:   "fs name system",
		Ctx:   "hdfs.capacity",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "fsns_capacity_remaining", Name: "remaining", Div: 1024},
			{ID: "fsns_capacity_used", Name: "used", Div: 1024},
		},
	},
	{
		ID:    "load",
		Title: "Current Number of Connections",
		Units: "num",
		Fam:   "fs name system",
		Ctx:   "hdfs.load",
		Dims: Dims{
			{ID: "fsns_total_load", Name: "connections"},
		},
	},
	{
		ID:    "volume_failures_total",
		Title: "Total Number of Volume Failures Across All Datanodes",
		Units: "events/s",
		Fam:   "fs name system",
		Ctx:   "hdfs.volume_failures_total",
		Dims: Dims{
			{ID: "fsns_volume_failures_total", Name: "failures", Algo: module.Incremental},
		},
	},
	{
		ID:    "data_nodes",
		Title: "Number of Data Nodes",
		Units: "num",
		Fam:   "fs name system",
		Ctx:   "hdfs.data_nodes",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "fsns_num_live_data_nodes", Name: "live"},
			{ID: "fsns_num_dead_data_nodes", Name: "dead"},
		},
	},
}

func unknownNodeCharts() *Charts {
	charts := Charts{}
	panicIfError(charts.Add(*jvmCharts.Copy()...))
	return &charts
}

func dataNodeCharts() *Charts {
	charts := Charts{}
	panicIfError(charts.Add(*jvmCharts.Copy()...))
	return &charts
}

func nameNodeCharts() *Charts {
	charts := Charts{}
	panicIfError(charts.Add(*jvmCharts.Copy()...))
	panicIfError(charts.Add(*fsnCharts.Copy()...))
	return &charts
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

package logstash

import "github.com/netdata/go-orchestrator/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

var charts = Charts{
	{
		ID:    "jvm_threads",
		Title: "JVM Threads",
		Units: "count",
		Fam:   "threads",
		Dims: Dims{
			{ID: "jvm_threads_count", Name: "threads"},
		},
	},
	{
		ID:    "jvm_mem_heap_percent",
		Title: "JVM Heap Memory Percentage",
		Units: "percent",
		Fam:   "memory",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "jvm_mem_heap_used_percent", Name: "in use"},
		},
	},
	{
		ID:    "jvm_mem_heap",
		Title: "JVM Heap Memory",
		Units: "KiB",
		Fam:   "memory",
		Type:  module.Area,
		Dims: Dims{
			{ID: "jvm_mem_heap_committed_in_bytes", Name: "committed", Div: 1024},
			{ID: "jvm_mem_heap_used_in_bytes", Name: "used", Div: 1024},
		},
	},
	{
		ID:    "jvm_mem_pools_survivor",
		Title: "JVM Pool Survivor Memory",
		Units: "KiB",
		Fam:   "memory",
		Type:  module.Area,
		Dims: Dims{
			{ID: "jvm_mem_pools_survivor_committed_in_bytes", Name: "committed", Div: 1024},
			{ID: "jvm_mem_pools_survivor_used_in_bytes", Name: "used", Div: 1024},
		},
	},
	{
		ID:    "jvm_mem_pools_old",
		Title: "JVM Pool Old Memory",
		Units: "KiB",
		Fam:   "memory",
		Type:  module.Area,
		Dims: Dims{
			{ID: "jvm_mem_pools_old_committed_in_bytes", Name: "committed", Div: 1024},
			{ID: "jvm_mem_pools_old_used_in_bytes", Name: "used", Div: 1024},
		},
	},
	{
		ID:    "jvm_mem_pools_eden",
		Title: "JVM Pool Eden Memory",
		Units: "KiB",
		Fam:   "memory",
		Type:  module.Area,
		Dims: Dims{
			{ID: "jvm_mem_pools_eden_committed_in_bytes", Name: "committed", Div: 1024},
			{ID: "jvm_mem_pools_eden_used_in_bytes", Name: "used", Div: 1024},
		},
	},
	{
		ID:    "jvm_gc_collector_count",
		Title: "Garbage Collection Count",
		Units: "counts/s",
		Fam:   "garbage collection",
		Dims: Dims{
			{ID: "jvm_gc_collectors_eden_collection_count", Name: "eden", Algo: module.Incremental},
			{ID: "jvm_gc_collectors_old_collection_count", Name: "old", Algo: module.Incremental},
		},
	},
	{
		ID:    "jvm_gc_collector_time",
		Title: "Time Spent On Garbage Collection",
		Units: "ms",
		Fam:   "garbage collection",
		Dims: Dims{
			{ID: "jvm_gc_collectors_eden_collection_time_in_millis", Name: "eden", Algo: module.Incremental},
			{ID: "jvm_gc_collectors_old_collection_time_in_millis", Name: "old", Algo: module.Incremental},
		},
	},
	{
		ID:    "uptime",
		Title: "Uptime",
		Units: "seconds",
		Fam:   "uptime",
		Dims: Dims{
			{ID: "jvm_uptime_in_millis", Name: "uptime", Div: 1000},
		},
	},
}

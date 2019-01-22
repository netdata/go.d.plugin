package logstash

import "github.com/netdata/go.d.plugin/modules"

type (
	// Charts is an alias for modules.Charts
	Charts = modules.Charts
	// Dims is an alias for modules.Dims
	Dims = modules.Dims
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
		Dims: Dims{
			{ID: "jvm_mem_heap_used_percent", Name: "in use"},
		},
	},
	{
		ID:    "jvm_mem_heap",
		Title: "JVM Heap Memory",
		Units: "KiB",
		Fam:   "memory",
		Dims: Dims{
			{ID: "jvm_mem_heap_used_in_bytes", Name: "used", Div: 1024},
			{ID: "jvm_mem_heap_committed_in_bytes", Name: "committed", Div: 1024},
		},
	},
	{
		ID:    "jvm_mem_pools_survivor",
		Title: "JVM Pool Survivor Memory",
		Units: "KiB",
		Fam:   "memory",
		Dims: Dims{
			{ID: "jvm_mem_pools_survivor_used_in_bytes", Name: "used", Div: 1024},
			{ID: "jvm_mem_pools_survivor_committed_in_bytes", Name: "committed", Div: 1024},
		},
	},
	{
		ID:    "jvm_mem_pools_old",
		Title: "JVM Pool Old Memory",
		Units: "KiB",
		Fam:   "memory",
		Dims: Dims{
			{ID: "jvm_mem_pools_old_used_in_bytes", Name: "used", Div: 1024},
			{ID: "jvm_mem_pools_old_committed_in_bytes", Name: "committed", Div: 1024},
		},
	},
	{
		ID:    "jvm_mem_pools_young",
		Title: "JVM Pool Young Memory",
		Units: "KiB",
		Fam:   "memory",
		Dims: Dims{
			{ID: "jvm_mem_pools_young_used_in_bytes", Name: "used", Div: 1024},
			{ID: "jvm_mem_pools_young_committed_in_bytes", Name: "committed", Div: 1024},
		},
	},
	{
		ID:    "jvm_gc_collector_count",
		Title: "Garbage Collection Count",
		Units: "counts/s",
		Fam:   "garbage collection",
		Dims: Dims{
			{ID: "jvm_gc_collectors_young_collection_count", Name: "young", Algo: modules.Incremental},
			{ID: "jvm_gc_collectors_young_collection_count", Name: "old", Algo: modules.Incremental},
		},
	},
	{
		ID:    "jvm_gc_collector_time",
		Title: "Time Spent On Garbage Collection",
		Units: "ms",
		Fam:   "garbage collection",
		Dims: Dims{
			{ID: "jvm_gc_collectors_young_collection_time_in_millis", Name: "young", Algo: modules.Incremental},
			{ID: "jvm_gc_collectors_young_collection_time_in_millis", Name: "old", Algo: modules.Incremental},
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

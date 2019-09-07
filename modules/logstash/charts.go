package logstash

import "github.com/netdata/go-orchestrator/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

var charts = Charts{
	// thread
	{
		ID:    "jvm_threads",
		Ctx:   "logstash.jvm_threads",
		Title: "JVM Threads",
		Units: "count",
		Fam:   "threads",
		Dims: Dims{
			{ID: "jvm_threads_count", Name: "threads"},
		},
	},
	// memory
	{
		ID:    "jvm_mem_heap_used",
		Title: "JVM Heap Memory Percentage",
		Ctx:   "logstash.jvm_mem_heap_used",
		Units: "percentage",
		Fam:   "memory",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "jvm_mem_heap_used_percent", Name: "in use"},
		},
	},
	{
		ID:    "jvm_mem_heap",
		Title: "JVM Heap Memory",
		Ctx:   "logstash.jvm_mem_heap",
		Units: "KiB",
		Fam:   "memory",
		Type:  module.Area,
		Dims: Dims{
			{ID: "jvm_mem_heap_committed_in_bytes", Name: "committed", Div: 1024},
			{ID: "jvm_mem_heap_used_in_bytes", Name: "used", Div: 1024},
		},
	},
	{
		ID:    "jvm_mem_pools_eden",
		Title: "JVM Pool Eden Memory",
		Ctx:   "logstash.jvm_mem_pools_eden",
		Units: "KiB",
		Fam:   "memory",
		Type:  module.Area,
		Dims: Dims{
			{ID: "jvm_mem_pools_eden_committed_in_bytes", Name: "committed", Div: 1024},
			{ID: "jvm_mem_pools_eden_used_in_bytes", Name: "used", Div: 1024},
		},
	},
	{
		ID:    "jvm_mem_pools_survivor",
		Title: "JVM Pool Survivor Memory",
		Ctx:   "logstash.jvm_mem_pools_survivor",
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
		Ctx:   "logstash.jvm_mem_pools_old",
		Units: "KiB",
		Fam:   "memory",
		Type:  module.Area,
		Dims: Dims{
			{ID: "jvm_mem_pools_old_committed_in_bytes", Name: "committed", Div: 1024},
			{ID: "jvm_mem_pools_old_used_in_bytes", Name: "used", Div: 1024},
		},
	},
	// garbage collection
	{
		ID:    "jvm_gc_collector_count",
		Title: "Garbage Collection Count",
		Ctx:   "logstash.jvm_gc_collector_count",
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
		Ctx:   "logstash.jvm_gc_collector_time",
		Units: "ms",
		Fam:   "garbage collection",
		Dims: Dims{
			{ID: "jvm_gc_collectors_eden_collection_time_in_millis", Name: "eden", Algo: module.Incremental},
			{ID: "jvm_gc_collectors_old_collection_time_in_millis", Name: "old", Algo: module.Incremental},
		},
	},
	// processes
	{
		ID:    "open_file_descriptors",
		Title: "Open File Descriptors",
		Ctx:   "logstash.open_file_descriptors",
		Units: "count",
		Fam:   "processes",
		Dims: Dims{
			{ID: "process_open_file_descriptors", Name: "open file descriptors"},
		},
	},
	// events
	{
		ID:    "event",
		Title: "Events Overview",
		Ctx:   "logstash.event",
		Units: "events/s",
		Fam:   "events",
		Dims: Dims{
			{ID: "event_in", Name: "in", Algo: module.Incremental},
			{ID: "event_filtered", Name: "filtered", Algo: module.Incremental},
			{ID: "event_out", Name: "out", Algo: module.Incremental},
		},
	},
	{
		ID:    "event_duration",
		Title: "Events Duration",
		Ctx:   "logstash.event_duration",
		Units: "seconds",
		Fam:   "events",
		Dims: Dims{
			{ID: "event_duration_in_millis", Name: "event", Div: 1000, Algo: module.Incremental},
			{ID: "event_queue_push_duration_in_millis", Name: "queue", Div: 1000, Algo: module.Incremental},
		},
	},
	// uptime
	{
		ID:    "uptime",
		Title: "Uptime",
		Ctx:   "logstash.uptime",
		Units: "seconds",
		Fam:   "uptime",
		Dims: Dims{
			{ID: "jvm_uptime_in_millis", Name: "uptime", Div: 1000},
		},
	},
}

func createPipelineChart(id string) Charts {
	return Charts{
		{
			ID:    "pipeline_" + id + "_event",
			Title: id + " Pipeline Events",
			Ctx:   "logstash.pipeline_event",
			Units: "events/s",
			Fam:   id,
			Dims: Dims{
				{ID: "event_in", Name: "in", Algo: module.Incremental},
				{ID: "event_filtered", Name: "filtered", Algo: module.Incremental},
				{ID: "event_out", Name: "out", Algo: module.Incremental},
			},
		},
		{
			ID:    "pipeline_" + id + "_event_duration",
			Title: id + " Pipeline Events Duration",
			Ctx:   "logstash.pipeline_event_duration",
			Units: "seconds",
			Fam:   id,
			Dims: Dims{
				{ID: "event_duration_in_millis", Name: "event", Div: 1000, Algo: module.Incremental},
				{ID: "event_queue_push_duration_in_millis", Name: "queue", Div: 1000, Algo: module.Incremental},
			},
		},
	}
}

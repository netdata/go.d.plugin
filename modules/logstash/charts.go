// SPDX-License-Identifier: GPL-3.0-or-later

package logstash

import (
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

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
		Title: "JVM Threads",
		Units: "count",
		Fam:   "threads",
		Ctx:   "logstash.jvm_threads",
		Dims: Dims{
			{ID: "jvm_threads_count", Name: "threads"},
		},
	},
	// memory
	{
		ID:    "jvm_mem_heap_used",
		Title: "JVM Heap Memory Percentage",
		Units: "percentage",
		Fam:   "memory",
		Ctx:   "logstash.jvm_mem_heap_used",
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
		Ctx:   "logstash.jvm_mem_heap",
		Type:  module.Area,
		Dims: Dims{
			{ID: "jvm_mem_heap_committed_in_bytes", Name: "committed", Div: 1024},
			{ID: "jvm_mem_heap_used_in_bytes", Name: "used", Div: 1024},
		},
	},
	{
		ID:    "jvm_mem_pools_eden",
		Title: "JVM Pool Eden Memory",
		Units: "KiB",
		Fam:   "memory",
		Ctx:   "logstash.jvm_mem_pools_eden",
		Type:  module.Area,
		Dims: Dims{
			{ID: "jvm_mem_pools_eden_committed_in_bytes", Name: "committed", Div: 1024},
			{ID: "jvm_mem_pools_eden_used_in_bytes", Name: "used", Div: 1024},
		},
	},
	{
		ID:    "jvm_mem_pools_survivor",
		Title: "JVM Pool Survivor Memory",
		Units: "KiB",
		Fam:   "memory",
		Ctx:   "logstash.jvm_mem_pools_survivor",
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
		Ctx:   "logstash.jvm_mem_pools_old",
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
		Units: "counts/s",
		Fam:   "garbage collection",
		Ctx:   "logstash.jvm_gc_collector_count",
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
		Ctx:   "logstash.jvm_gc_collector_time",
		Dims: Dims{
			{ID: "jvm_gc_collectors_eden_collection_time_in_millis", Name: "eden", Algo: module.Incremental},
			{ID: "jvm_gc_collectors_old_collection_time_in_millis", Name: "old", Algo: module.Incremental},
		},
	},
	// processes
	{
		ID:    "open_file_descriptors",
		Title: "Open File Descriptors",
		Units: "fd",
		Fam:   "processes",
		Ctx:   "logstash.open_file_descriptors",
		Dims: Dims{
			{ID: "process_open_file_descriptors", Name: "open"},
		},
	},
	// events
	{
		ID:    "event",
		Title: "Events Overview",
		Units: "events/s",
		Fam:   "events",
		Ctx:   "logstash.event",
		Dims: Dims{
			{ID: "event_in", Name: "in", Algo: module.Incremental},
			{ID: "event_filtered", Name: "filtered", Algo: module.Incremental},
			{ID: "event_out", Name: "out", Algo: module.Incremental},
		},
	},
	{
		ID:    "event_duration",
		Title: "Events Duration",
		Units: "seconds",
		Fam:   "events",
		Ctx:   "logstash.event_duration",
		Dims: Dims{
			{ID: "event_duration_in_millis", Name: "event", Div: 1000, Algo: module.Incremental},
			{ID: "event_queue_push_duration_in_millis", Name: "queue", Div: 1000, Algo: module.Incremental},
		},
	},
	// uptime
	{
		ID:    "uptime",
		Title: "Uptime",
		Units: "seconds",
		Fam:   "uptime",
		Ctx:   "logstash.uptime",
		Dims: Dims{
			{ID: "jvm_uptime_in_millis", Name: "uptime", Div: 1000},
		},
	},
}

var pipelineChartsTemplate = Charts{
	{
		ID:    "pipeline_%s_event",
		Title: "%s Pipeline Events",
		Units: "events/s",
		Fam:   "%s",
		Ctx:   "logstash.pipeline_event",
		Dims: Dims{
			{ID: "pipelines_%s_event_in", Name: "in", Algo: module.Incremental},
			{ID: "pipelines_%s_event_filtered", Name: "filtered", Algo: module.Incremental},
			{ID: "pipelines_%s_event_out", Name: "out", Algo: module.Incremental},
		},
	},
	{
		ID:    "pipeline_%s_event_duration",
		Title: "%s Pipeline Events Duration",
		Units: "seconds",
		Fam:   "%s",
		Ctx:   "logstash.pipeline_event_duration",
		Dims: Dims{
			{ID: "pipelines_%s_event_duration_in_millis", Name: "event", Div: 1000, Algo: module.Incremental},
			{ID: "pipelines_%s_event_queue_push_duration_in_millis", Name: "queue", Div: 1000, Algo: module.Incremental},
		},
	},
}

func pipelineCharts(id string) *Charts {
	cs := pipelineChartsTemplate.Copy()
	for _, chart := range *cs {
		chart.ID = fmt.Sprintf(chart.ID, id)
		chart.Title = fmt.Sprintf(chart.Title, id)
		chart.Fam = fmt.Sprintf(chart.Fam, id)
		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, id)
		}
	}
	return cs
}

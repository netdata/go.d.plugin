// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

import "github.com/netdata/go.d.plugin/agent/module"

const (
	prioRequests = module.Priority + iota

	prioLatency

	prioDroppedMessages
	prioRequestsTimeouts
	prioRequestsUnavailables
	prioRequestsFailures
	prioStorageExceptions

	prioCacheHitRatio
	prioCacheHitRate
	prioCacheSize

	prioStorageLoad

	prioThreadPoolsActiveTasks
	prioThreadPoolsPendingTasks
	prioThreadPoolsBlockedTasks
	prioThreadPoolsCurrentlyBlockedTasks

	prioCompactionCompletedTasks
	prioCompactionPendingTasks
	prioCompactionBytesCompacted

	prioJVMGCCount
	prioJVMGCTime
)

var baseCharts = module.Charts{
	chartClientRequests.Copy(),

	chartClientRequestsLatency.Copy(),

	chartDroppedMessages.Copy(),
	chartClientRequestTimeouts.Copy(),
	chartClientRequestUnavailables.Copy(),
	chartClientRequestFailures.Copy(),
	chartStorageExceptions.Copy(),

	chartCacheHitRatio.Copy(),
	chartCacheHitRate.Copy(),
	chartCacheSize.Copy(),

	chartStorageLoad.Copy(),

	threadPoolsActiveTasks.Copy(),
	threadPoolsPendingTasks.Copy(),
	threadPoolsBlockedTasks.Copy(),
	threadPoolsCurrentlyBlockedTasks.Copy(),

	chartCompactionCompletedTasks.Copy(),
	chartCompactionPendingTasks.Copy(),
	chartCompactionBytesCompacted.Copy(),

	chartJVMGCCount.Copy(),
	chartJVMGCTime.Copy(),
}

var (
	chartClientRequests = module.Chart{
		ID:       "client_requests",
		Title:    "Client requests",
		Units:    "requests/s",
		Fam:      "throughput",
		Ctx:      "cassandra.client_requests",
		Priority: prioRequests,
		Dims: module.Dims{
			{ID: "client_request_latency_reads", Name: "reads", Algo: module.Incremental, Div: 1000},
			{ID: "client_request_latency_writes", Name: "writes", Algo: module.Incremental, Mul: -1, Div: 1000},
		},
	}
)

var (
	chartClientRequestsLatency = module.Chart{
		ID:       "client_requests_latency",
		Title:    "Client requests latency",
		Units:    "microseconds",
		Fam:      "latency",
		Ctx:      "cassandra.client_requests_latency",
		Priority: prioLatency,
		Dims: module.Dims{
			{ID: "client_request_total_latency_reads", Name: "read", Algo: module.Incremental, Div: 1000},
			{ID: "client_request_total_latency_writes", Name: "write", Algo: module.Incremental, Mul: -1, Div: 1000},
		},
	}
)

var (
	chartDroppedMessages = module.Chart{
		ID:       "dropped_messages_one_minute_rate",
		Title:    "Dropped messages",
		Units:    "messages/s",
		Fam:      "errors",
		Ctx:      "cassandra.dropped_messages_one_minute_rate",
		Priority: prioDroppedMessages,
		Dims: module.Dims{
			{ID: "dropped_messages_one_minute", Name: "dropped", Div: 1000},
		},
	}
	chartClientRequestTimeouts = module.Chart{
		ID:       "client_requests_timeouts",
		Title:    "Client requests timeouts",
		Units:    "timeouts/s",
		Fam:      "errors",
		Ctx:      "cassandra.client_requests_timeouts",
		Priority: prioRequestsTimeouts,
		Dims: module.Dims{
			{ID: "client_request_timeouts_reads", Name: "read", Algo: module.Incremental, Div: 1000},
			{ID: "client_request_timeouts_writes", Name: "write", Algo: module.Incremental, Mul: -1, Div: 1000},
		},
	}
	chartClientRequestUnavailables = module.Chart{
		ID:       "client_requests_unavailables",
		Title:    "Client requests unavailable exceptions",
		Units:    "exceptions/s",
		Fam:      "errors",
		Ctx:      "cassandra.client_requests_unavailables",
		Priority: prioRequestsUnavailables,
		Dims: module.Dims{
			{ID: "client_request_unavailables_reads", Name: "read", Algo: module.Incremental, Div: 1000},
			{ID: "client_request_unavailables_writes", Name: "write", Algo: module.Incremental, Mul: -1, Div: 1000},
		},
	}
	chartClientRequestFailures = module.Chart{
		ID:       "client_requests_failures",
		Title:    "Client requests failures",
		Units:    "failures/s",
		Fam:      "errors",
		Ctx:      "cassandra.client_requests_failures",
		Priority: prioRequestsFailures,
		Dims: module.Dims{
			{ID: "client_request_failures_reads", Name: "read", Algo: module.Incremental, Div: 1000},
			{ID: "client_request_failures_writes", Name: "write", Algo: module.Incremental, Mul: -1, Div: 1000},
		},
	}
	chartStorageExceptions = module.Chart{
		ID:       "storage_exceptions",
		Title:    "Storage exceptions",
		Units:    "exceptions",
		Fam:      "errors",
		Ctx:      "cassandra.storage_exceptions",
		Priority: prioStorageExceptions,
		Dims: module.Dims{
			{ID: "storage_exceptions", Name: "storage", Algo: module.Incremental},
		},
	}
)

var (
	chartCacheHitRatio = module.Chart{
		ID:       "cache_hit_ratio",
		Title:    "Cache hit ratio",
		Units:    "percentage",
		Fam:      "cache",
		Ctx:      "cassandra.cache_hit_ratio",
		Priority: prioCacheHitRatio,
		Dims: module.Dims{
			{ID: "cache_hit_ratio", Name: "hit_ratio", Div: 1000},
		},
	}
	chartCacheHitRate = module.Chart{
		ID:       "cache_hit_rate",
		Title:    "Cache hit rate",
		Units:    "events/s",
		Fam:      "cache",
		Ctx:      "cassandra.cache_hit_rate",
		Priority: prioCacheHitRate,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "cache_hits", Name: "hits", Algo: module.Incremental, Div: 1000},
			{ID: "cache_misses", Name: "misses", Algo: module.Incremental, Div: 1000},
		},
	}
	chartCacheSize = module.Chart{
		ID:       "cache_size",
		Title:    "Cache size",
		Units:    "bytes",
		Fam:      "cache",
		Ctx:      "cassandra.cache_size",
		Priority: prioCacheSize,
		Dims: module.Dims{
			{ID: "cache_size", Name: "size"},
		},
	}
)

var (
	chartStorageLoad = module.Chart{
		ID:       "storage_load",
		Title:    "Disk space used by live data on a node",
		Units:    "bytes",
		Fam:      "disk usage",
		Ctx:      "cassandra.storage_load",
		Priority: prioStorageLoad,
		Dims: module.Dims{
			{ID: "storage_load", Name: "used"},
		},
	}
)

var (
	threadPoolsActiveTasks = module.Chart{
		ID:       "thread_pools_active_tasks",
		Title:    "Active tasks",
		Units:    "tasks",
		Fam:      "thread pools",
		Ctx:      "cassandra.thread_pools_active_tasks",
		Priority: prioThreadPoolsActiveTasks,
		Dims: module.Dims{
			{ID: "thread_pools_active_tasks", Name: "active"},
		},
	}
	threadPoolsPendingTasks = module.Chart{
		ID:       "thread_pools_pending_tasks",
		Title:    "Pending tasks",
		Units:    "tasks",
		Fam:      "thread pools",
		Ctx:      "cassandra.thread_pools_pending_tasks",
		Priority: prioThreadPoolsPendingTasks,
		Dims: module.Dims{
			{ID: "thread_pools_pending_tasks", Name: "pending"},
		},
	}
	threadPoolsBlockedTasks = module.Chart{
		ID:       "thread_pools_blocked_tasks",
		Title:    "Blocked tasks",
		Units:    "tasks/s",
		Fam:      "thread pools",
		Ctx:      "cassandra.thread_pools_blocked_tasks",
		Priority: prioThreadPoolsBlockedTasks,
		Dims: module.Dims{
			{ID: "thread_pools_total_blocked_tasks", Name: "blocked", Algo: module.Incremental},
		},
	}
	threadPoolsCurrentlyBlockedTasks = module.Chart{
		ID:       "thread_pools_currently_blocked_tasks",
		Title:    "Blocked tasks",
		Units:    "tasks",
		Fam:      "thread pools",
		Ctx:      "cassandra.thread_pools_currently_blocked_tasks",
		Priority: prioThreadPoolsCurrentlyBlockedTasks,
		Dims: module.Dims{
			{ID: "thread_pools_currently_blocked_tasks", Name: "blocked"},
		},
	}
)

var (
	chartCompactionCompletedTasks = module.Chart{
		ID:       "compaction_completed_tasks",
		Title:    "Pending compactions",
		Units:    "tasks",
		Fam:      "compaction",
		Ctx:      "cassandra.compaction_completed_tasks",
		Priority: prioCompactionCompletedTasks,
		Dims: module.Dims{
			{ID: "compaction_completed_tasks", Name: "completed", Algo: module.Incremental},
		},
	}
	chartCompactionPendingTasks = module.Chart{
		ID:       "compaction_pending_tasks",
		Title:    "Pending compactions",
		Units:    "tasks",
		Fam:      "compaction",
		Ctx:      "cassandra.compaction_pending_tasks",
		Priority: prioCompactionPendingTasks,
		Dims: module.Dims{
			{ID: "compaction_pending_tasks", Name: "pending"},
		},
	}
	chartCompactionBytesCompacted = module.Chart{
		ID:       "compaction_compacted",
		Title:    "Compaction",
		Units:    "bytes/s",
		Fam:      "compaction",
		Ctx:      "cassandra.compaction_compacted",
		Priority: prioCompactionBytesCompacted,
		Dims: module.Dims{
			{ID: "compaction_bytes_compacted", Name: "compacted", Algo: module.Incremental},
		},
	}
)

var (
	chartJVMGCCount = module.Chart{
		ID:       "jvm_gc_count",
		Title:    "Garbage collections",
		Units:    "gc/s",
		Fam:      "jvm gc",
		Ctx:      "cassandra.jvm_gc_count",
		Priority: prioJVMGCCount,
		Dims: module.Dims{
			{ID: "jvm_gc_parnew_count", Name: "parnew", Algo: module.Incremental, Div: 1000},
			{ID: "jvm_gc_cms_count", Name: "cms", Algo: module.Incremental, Div: 1000},
		},
	}
	chartJVMGCTime = module.Chart{
		ID:       "jvm_gc_time",
		Title:    "Garbage collection time",
		Units:    "us",
		Fam:      "jvm gc",
		Ctx:      "cassandra.jvm_gc_time",
		Priority: prioJVMGCTime,
		Dims: module.Dims{
			{ID: "jvm_gc_parnew_time", Name: "parnew", Algo: module.Incremental, Div: 1000},
			{ID: "jvm_gc_cms_time", Name: "cms", Algo: module.Incremental, Div: 1000},
		},
	}
)

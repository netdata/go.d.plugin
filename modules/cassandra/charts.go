// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

import "github.com/netdata/go.d.plugin/agent/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Chart is an alias for module.Chart
	Chart = module.Chart
	// Dims is an alias for module.Dims
	Dims = module.Dims
	// Dim is an alias for module.Dim
	Dim = module.Dim
	// Vars is an alias for module.Vars
	Vars = module.Vars
	// Opts is an alias for module.Dim
	Opts = module.DimOpts
)

var chartCassandraThroughput = Chart{
	ID:    "throughput",
	Title: "I/O requests.",
	Units: "requests/s",
	Fam:   "throughput",
	Ctx:   "cassandra.throughput",
	Type:  module.Line,
	Dims: Dims{
		{ID: "throughput_Read", Name: "read", Algo: module.Incremental},
		{ID: "throughput_Write", Name: "write", Algo: module.Incremental},
	},
}

var chartCassandraLatency = Chart{
	ID:    "latency",
	Title: "I/O latency.",
	Units: "requests/s",
	Fam:   "latency",
	Ctx:   "cassandra.latency",
	Type:  module.Line,
	Dims: Dims{
		{ID: "read", Name: "read", Algo: module.Incremental},
		{ID: "write", Name: "write", Algo: module.Incremental},
	},
}

var chartCassandraCache = Chart{
	ID:    "cache",
	Title: "Cache Hit",
	Units: "percentage/s",
	Fam:   "cache",
	Ctx:   "cassandra.cache",
	Type:  module.Line,
	Dims: Dims{
		{ID: "cache_HitRate", Name: "ratio", Algo: module.Incremental},
	},
}

var chartCassandraDiskLoad = Chart{
	ID:    "disk_load",
	Title: "Disk Load",
	Units: "bytes/s",
	Fam:   "disk_load",
	Ctx:   "cassandra.node",
	Type:  module.Line,
	Dims: Dims{
		{ID: "disk_LiveDiskSpaceUsed", Name: "space", Algo: module.Incremental},
	},
}

var chartCassandraDiskColumn = Chart{
	ID:    "disk_column",
	Title: "Disk Column",
	Units: "bytes/s",
	Fam:   "column",
	Ctx:   "cassandra.disk_column",
	Type:  module.Line,
	Dims: Dims{
		{ID: "disk_TotalDiskSpaceUsed", Name: "space", Algo: module.Incremental},
	},
}

var chartCassandraDiskCompactionCompleted = Chart{
	ID:    "compaction_completed",
	Title: "Completed Compaction Tasks",
	Units: "events/s",
	Fam:   "compaction",
	Ctx:   "cassandra.compaction_completed",
	Type:  module.Line,
	Dims: Dims{
		{ID: "disk_CompactionBytesWritten", Name: "compaction", Algo: module.Incremental},
	},
}

var chartCassandraDiskCompactionQueue = Chart{
	ID:    "compaction_queue",
	Title: "Queued Compaction Tasks",
	Units: "events/s",
	Fam:   "queue",
	Ctx:   "cassandra.compaction_queued",
	Type:  module.Line,
	Dims: Dims{
		{ID: "disk_PendingCompactions", Name: "queue", Algo: module.Incremental},
	},
}

var chartCassandraParNewCount = Chart{
	ID:    "gc_parnew_count",
	Title: "Young-generation garbage collection counter",
	Units: "events/s",
	Fam:   "par new count",
	Ctx:   "cassandra.gc_parnew_count",
	Type:  module.Line,
	Dims: Dims{
		{ID: "java_gc_count_ParNew", Name: "parnew", Algo: module.Incremental},
	},
}

var chartCassandraParNewTime = Chart{
	ID:    "gc_parnew_time",
	Title: "Young-generation garbage collection timer",
	Units: "period of time",
	Fam:   "par new time",
	Ctx:   "cassandra.gc_parnew_time",
	Type:  module.Line,
	Dims: Dims{
		{ID: "java_gc_time_ParNew", Name: "parnew", Algo: module.Incremental},
	},
}

var chartCassandraMarkSweepCount = Chart{
	ID:    "gc_marksweep_count",
	Title: "Old-generation collection",
	Units: "events/s",
	Fam:   "mark sweep",
	Ctx:   "cassandra.gc_sweep_count",
	Type:  module.Line,
	Dims: Dims{
		{ID: "java_gc_count_ConcurrentMarkSweep", Name: "sweep", Algo: module.Incremental},
	},
}

var chartCassandraMarkSweepTime = Chart{
	ID:    "gc_marksweep_time",
	Title: "Elapsed time Old-generation collection",
	Units: "period of time",
	Fam:   "mark sweep",
	Ctx:   "cassandra.gc_sweep_time",
	Type:  module.Line,
	Dims: Dims{
		{ID: "java_gc_time_ConcurrentMarkSweep", Name: "sweep", Algo: module.Incremental},
	},
}

var chartCassandraErrorTimeout = Chart{
	ID:    "error_timeout",
	Title: "Requests not unacknowledged",
	Units: "requests/s",
	Fam:   "request timeout",
	Ctx:   "cassandra.error_timeout",
	Type:  module.Line,
	Dims: Dims{
		{ID: "error_timeout_Read", Name: "read", Algo: module.Incremental},
		{ID: "error_timeout_Write", Name: "write", Algo: module.Incremental},
	},
}

var chartCassandraErrorUnavailable = Chart{
	ID:    "error_timeout",
	Title: "Request was unavailable",
	Units: "requests/s",
	Fam:   "request unavailable",
	Ctx:   "cassandra.error_unavailable",
	Type:  module.Line,
	Dims: Dims{
		{ID: "error_unavailable_Read", Name: "read", Algo: module.Incremental},
		{ID: "error_unavailable_Write", Name: "write", Algo: module.Incremental},
	},
}

var chartCassandraPendingTasks = Chart{
	ID:    "pending_task",
	Title: "Task queued",
	Units: "tasks/s",
	Fam:   "task queued",
	Ctx:   "cassandra.task_queued",
	Type:  module.Line,
	Dims: Dims{
		{ID: "pending_tasks_tasks", Name: "task", Algo: module.Incremental},
	},
}

var chartCassandraBlockedTasks = Chart{
	ID:    "blocked_task",
	Title: "Task blocked",
	Units: "tasks/s",
	Fam:   "task queued",
	Ctx:   "cassandra.task_blocked",
	Type:  module.Line,
	Dims: Dims{
		{ID: "task", Name: "task", Algo: module.Incremental},
	},
}

func newCassandraCharts() *Charts {
	return &Charts{
		chartCassandraThroughput.Copy(),
		chartCassandraLatency.Copy(),
		chartCassandraCache.Copy(),
		chartCassandraDiskLoad.Copy(),
		chartCassandraDiskColumn.Copy(),
		chartCassandraDiskCompactionCompleted.Copy(),
		chartCassandraDiskCompactionQueue.Copy(),
		chartCassandraParNewCount.Copy(),
		chartCassandraParNewTime.Copy(),
		chartCassandraMarkSweepCount.Copy(),
		chartCassandraMarkSweepTime.Copy(),
		chartCassandraErrorTimeout.Copy(),
		chartCassandraErrorUnavailable.Copy(),
		chartCassandraPendingTasks.Copy(),
		chartCassandraBlockedTasks.Copy(),
	}
}

package cockroachdb

import "github.com/netdata/go-orchestrator/module"

type (
	Charts = module.Charts
	Chart  = module.Chart
	Dims   = module.Dims
	Vars   = module.Vars
)

var charts = Charts{
	chartLiveNodes.Copy(),
	chartHeardBeats.Copy(),

	chartCapacity.Copy(),
	chartCapacityUsability.Copy(),
	chartCapacityUsable.Copy(),
	chartCapacityUsedPercentage.Copy(),

	chartUsedLiveData.Copy(),
	chartLogicalDataBytes.Copy(),
	chartLogicalData.Copy(),

	chartRocksDBReadAmplification.Copy(),
	chartRocksDBTableOperations.Copy(),
	chartRocksDBCacheUsage.Copy(),
	chartRocksDBCacheOperations.Copy(),
	chartRocksDBCacheHitRage.Copy(),
	chartRocksDBSSTables.Copy(),

	chartTimeSeriesWrittenSamples.Copy(),
	chartTimeSeriesWriteErrors.Copy(),
	chartTimeSeriesWrittenBytes.Copy(),

	chartSQLConnections.Copy(),
	chartSQLTraffic.Copy(),
	chartSQLQueries.Copy(),
	chartSQLErroredQueries.Copy(),
	chartSQLActiveDistQueries.Copy(),
	chartSQLActiveFlowsForDistQueries.Copy(),
	chartSQLTransactions.Copy(),
	chartSQLSchemaChanges.Copy(),

	chartRanges.Copy(),
	chartRangesWithProblems.Copy(),
	chartRangesEvents.Copy(),
	chartRangesSnapshotEvents.Copy(),

	chartReplicas.Copy(),
	chartReplicasQuiescence.Copy(),
	chartReplicasLeaders.Copy(),
	chartReplicasLeaseHolder.Copy(),

	chartRebalancingQueries.Copy(),
	chartRebalancingWrites.Copy(),

	chartSlowRequests.Copy(),

	chartRPCClockOffset.Copy(),

	chartGoroutines.Copy(),
	chartGoCgoHeapMemory.Copy(),
	chartCGoCalls.Copy(),
	chartGCRuns.Copy(),
	chartGCPauseTime.Copy(),

	chartProcessCPUUsage.Copy(),
	chartProcessCPUUsagePercent.Copy(),
	chartProcessCPUUsageCombinedPercent.Copy(),
	chartProcessMemory.Copy(),
	chartProcessFDUsage.Copy(),

	chartHostDiskBandwidth.Copy(),
	chartHostDiskOperations.Copy(),
	chartHostDiskIOPS.Copy(),
	chartHostNetworkBandwidth.Copy(),
	chartHostNetworkPackets.Copy(),

	chartUptime.Copy(),
}

// Liveness
var (
	chartLiveNodes = Chart{
		ID:    "live_nodes",
		Title: "Live Nodes",
		Units: "nodes",
		Fam:   "liveness",
		Ctx:   "cockroachdb.live_nodes",
		Dims: Dims{
			{ID: metricLiveNodes, Name: "live"},
		},
	}
	chartHeardBeats = Chart{
		ID:    "heartbeats",
		Title: "HeartBeats",
		Units: "heartbeats",
		Fam:   "liveness",
		Ctx:   "cockroachdb.heartbeats",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: metricHeartBeatSuccesses, Name: "successful", Algo: module.Incremental},
			{ID: metricHeartBeatFailures, Name: "unsuccessful", Algo: module.Incremental},
		},
	}
)

// Capacity
var (
	chartCapacity = Chart{
		ID:    "total_storage_capacity",
		Title: "Total Storage Capacity",
		Units: "KiB",
		Fam:   "capacity",
		Ctx:   "cockroachdb.total_storage_capacity",
		Dims: Dims{
			{ID: metricCapacity, Name: "total", Div: 1024},
		},
	}
	chartCapacityUsability = Chart{
		ID:    "storage_capacity_usability",
		Title: "Storage Capacity Usability",
		Units: "KiB",
		Fam:   "capacity",
		Ctx:   "cockroachdb.storage_capacity_usability",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: metricCapacityUsable, Name: "usable", Div: 1024},
			{ID: metricCapacityUnusable, Name: "unusable", Div: 1024},
		},
	}
	chartCapacityUsable = Chart{
		ID:    "storage_usable_capacity",
		Title: "Storage Usable Capacity",
		Units: "KiB",
		Fam:   "capacity",
		Ctx:   "cockroachdb.storage_usable_capacity",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: metricCapacityAvailable, Name: "available", Div: 1024},
			{ID: metricCapacityUsed, Name: "used", Div: 1024},
		},
	}
	chartCapacityUsedPercentage = Chart{
		ID:    "storage_used_capacity_percentage",
		Title: "Storage Used Capacity",
		Units: "percentage",
		Fam:   "capacity",
		Ctx:   "cockroachdb.storage_used_capacity_percentage",
		Dims: Dims{
			{ID: metricCapacityUsedPercentage, Name: "total", Div: precision},
			{ID: metricCapacityUsableUsedPercentage, Name: "usable", Div: precision},
		},
	}
)

// Storage
var (
	chartUsedLiveData = Chart{
		ID:    "live_bytes",
		Title: "The Amount of Used Live Data",
		Units: "KiB",
		Fam:   "storage",
		Ctx:   "cockroachdb.live_bytes",
		Dims: Dims{
			{ID: metricLiveBytes, Name: "applications", Div: 1024},
			{ID: metricSysBytes, Name: "system", Div: 1024},
		},
	}
	chartLogicalDataBytes = Chart{
		ID:    "logical_data_usage",
		Title: "The Amount of Logical Data",
		Units: "KiB",
		Fam:   "storage",
		Ctx:   "cockroachdb.logical_data_usage",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: metricKeyBytes, Name: "keys", Div: 1024},
			{ID: metricValBytes, Name: "values", Div: 1024},
		},
	}
	chartLogicalData = Chart{
		ID:    "logical_data",
		Title: "The Amount of Logical Data",
		Units: "num",
		Fam:   "storage",
		Ctx:   "cockroachdb.logical_data",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: metricKeyCount, Name: "keys"},
			{ID: metricValCount, Name: "values"},
		},
	}
)

// RocksDB
var (
	chartRocksDBReadAmplification = Chart{
		ID:    "rocksdb_read_amplification",
		Title: "RocksDB  Read Amplification",
		Units: "reads/query",
		Fam:   "rocksdb",
		Ctx:   "cockroachdb.rocksdb_read_amplification",
		Dims: Dims{
			{ID: metricRocksDBReadAmplification, Name: "reads"},
		},
	}
	chartRocksDBTableOperations = Chart{
		ID:    "rocksdb_table_operations",
		Title: "RocksDB Table Operations",
		Units: "operations",
		Fam:   "rocksdb",
		Ctx:   "cockroachdb.rocksdb_table_operations",
		Dims: Dims{
			{ID: metricRocksDBCompactions, Name: "compactions", Algo: module.Incremental},
			{ID: metricRocksDBFlushes, Name: "flushes", Algo: module.Incremental},
		},
	}
	chartRocksDBCacheUsage = Chart{
		ID:    "rocksdb_cache_usage",
		Title: "RocksDB Block Cache Usage",
		Units: "KiB",
		Fam:   "rocksdb",
		Ctx:   "cockroachdb.rocksdb_cache_usage",
		Type:  module.Area,
		Dims: Dims{
			{ID: metricRocksDBBlockCacheUsage, Name: "used", Div: 1024},
		},
	}
	chartRocksDBCacheOperations = Chart{
		ID:    "rocksdb_cache_operations",
		Title: "RocksDB Block Cache Operations",
		Units: "operations/s",
		Fam:   "rocksdb",
		Ctx:   "cockroachdb.rocksdb_cache_operations",
		Type:  module.Area,
		Dims: Dims{
			{ID: metricRocksDBBlockCacheHits, Name: "hits", Algo: module.Incremental},
			{ID: metricRocksDBBlockCacheMisses, Name: "misses", Algo: module.Incremental},
		},
	}
	chartRocksDBCacheHitRage = Chart{
		ID:    "rocksdb_cache_hit_rate",
		Title: "RocksDB Block Cache Hit Rate",
		Units: "percentage",
		Fam:   "rocksdb",
		Ctx:   "cockroachdb.rocksdb_cache_hit_rate",
		Type:  module.Area,
		Dims: Dims{
			{ID: metricRocksDBBlockCacheHitRate, Name: "hit rate"},
		},
	}
	chartRocksDBSSTables = Chart{
		ID:    "rocksdb_sstables",
		Title: "RocksDB SSTables",
		Units: "sstables",
		Fam:   "rocksdb",
		Ctx:   "cockroachdb.rocksdb_sstables",
		Dims: Dims{
			{ID: metricRocksDBNumSSTables, Name: "sstables"},
		},
	}
)

// Time Series
var (
	chartTimeSeriesWrittenSamples = Chart{
		ID:    "timeseries_samples",
		Title: "Time Series Written Samples",
		Units: "samples/s",
		Fam:   "time series",
		Ctx:   "cockroachdb.timeseries_samples",
		Dims: Dims{
			{ID: metricTimeSeriesWriteSamples, Name: "written", Algo: module.Incremental},
		},
	}
	chartTimeSeriesWriteErrors = Chart{
		ID:    "timeseries_write_errors",
		Title: "Time Series Write Errors",
		Units: "errors/s",
		Fam:   "time series",
		Ctx:   "cockroachdb.timeseries_write_errors",
		Dims: Dims{
			{ID: metricTimeSeriesWriteErrors, Name: "write", Algo: module.Incremental},
		},
	}
	chartTimeSeriesWrittenBytes = Chart{
		ID:    "timeseries_write_bytes",
		Title: "Time Series Bytes Written",
		Units: "KiB/s",
		Fam:   "time series",
		Ctx:   "cockroachdb.timeseries_write_bytes",
		Dims: Dims{
			{ID: metricTimeSeriesWriteBytes, Name: "written", Algo: module.Incremental},
		},
	}
)

// Ranges
var (
	chartRanges = Chart{
		ID:    "ranges",
		Title: "Number of Ranges",
		Units: "ranges",
		Fam:   "ranges",
		Ctx:   "cockroachdb.ranges",
		Dims: Dims{
			{ID: metricRanges, Name: "ranges"},
		},
	}
	chartRangesWithProblems = Chart{
		ID:    "ranges_replica_problems",
		Title: "Ranges With Problems",
		Units: "ranges",
		Fam:   "ranges",
		Ctx:   "cockroachdb.ranges_replica_problems",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: metricRangesUnavailable, Name: "unavailable"},
			{ID: metricRangesUnderReplicated, Name: "under_replicated"},
			{ID: metricRangesOverReplicated, Name: "over_replicated"},
		},
	}
	chartRangesEvents = Chart{
		ID:    "ranges_events",
		Title: "Ranges Events",
		Units: "events",
		Fam:   "ranges",
		Ctx:   "cockroachdb.ranges_events",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: metricRangeSplits, Name: "split", Algo: module.Incremental},
			{ID: metricRangeAdds, Name: "add", Algo: module.Incremental},
			{ID: metricRangeRemoves, Name: "remove", Algo: module.Incremental},
			{ID: metricRangeMerges, Name: "merge", Algo: module.Incremental},
		},
	}
	chartRangesSnapshotEvents = Chart{
		ID:    "ranges_snapshots",
		Title: "Ranges Snapshots",
		Units: "snapshots",
		Fam:   "ranges",
		Ctx:   "cockroachdb.ranges_snapshots",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: metricRangeSnapshotsGenerated, Name: "generated", Algo: module.Incremental},
			{ID: metricRangeSnapshotsNormalApplied, Name: "applied", Algo: module.Incremental},
			{ID: metricRangeSnapshotsLearnerApplied, Name: "applied learner", Algo: module.Incremental},
			{ID: metricRangeSnapshotsPreemptiveApplied, Name: "applied pre-emptive", Algo: module.Incremental},
		},
	}
)

// Replicas
var (
	chartReplicas = Chart{
		ID:    "replicas",
		Title: "Number of Replicas",
		Units: "num",
		Fam:   "replicas",
		Ctx:   "cockroachdb.replicas",
		Dims: Dims{
			{ID: metricReplicas, Name: "replicas"},
		},
	}
	chartReplicasQuiescence = Chart{
		ID:    "replicas_quiescence",
		Title: "Replicas Quiescence",
		Units: "replicas",
		Fam:   "replicas",
		Ctx:   "cockroachdb.replicas_quiescence",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: metricReplicasQuiescent, Name: "quiescent"},
			{ID: metricReplicasActive, Name: "active"},
		},
	}
	chartReplicasLeaders = Chart{
		ID:    "replicas_leaders",
		Title: "Number of Raft Leaders",
		Units: "num",
		Fam:   "replicas",
		Ctx:   "cockroachdb.replicas",
		Type:  module.Area,
		Dims: Dims{
			{ID: metricReplicasLeaders, Name: "leaders"},
			{ID: metricReplicasLeadersNotLeaseholders, Name: "not leaseholders"},
		},
	}
	chartReplicasLeaseHolder = Chart{
		ID:    "replicas_leaseholders",
		Title: "Number of Leaseholders",
		Units: "num",
		Fam:   "replicas",
		Ctx:   "cockroachdb.replicas_leaseholders",
		Dims: Dims{
			{ID: metricReplicasLeaseholders, Name: "leaseholders"},
		},
	}
)

// Rebalancing
var (
	chartRebalancingQueries = Chart{
		ID:    "rebalancing_queries",
		Title: "Rebalancing Average Queries",
		Units: "queries/s",
		Fam:   "rebalancing",
		Ctx:   "cockroachdb.rebalancing_queries",
		Dims: Dims{
			{ID: metricRebalancingQueriesPerSecond, Name: "queries", Div: precision},
		},
	}
	chartRebalancingWrites = Chart{
		ID:    "rebalancing_writes",
		Title: "Rebalancing Average Number of Written Keys To The Store",
		Units: "writes/s",
		Fam:   "rebalancing",
		Ctx:   "cockroachdb.rebalancing_queries",
		Dims: Dims{
			{ID: metricRebalancingWritesPerSecond, Name: "write", Div: precision},
		},
	}
)

// Slow Requests
var (
	chartSlowRequests = Chart{
		ID:    "slow_requests",
		Title: "Number of Requests that have been stuck for a long time",
		Units: "requests",
		Fam:   "slow requests",
		Ctx:   "cockroachdb.slow_requests",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: metricRequestsSlowLatch, Name: "acquiring_latches"},
			{ID: metricRequestsSlowLease, Name: "acquiring_lease"},
			{ID: metricRequestsSlowRaft, Name: "in_raft"},
		},
	}
)

// Go/CGo
var (
	chartGoCgoHeapMemory = Chart{
		ID:    "code_heap_memory_usage",
		Title: "Go/CGo Heap Memory Usage",
		Units: "KiB",
		Fam:   "go/cgo",
		Ctx:   "cockroachdb.code_heap_memory_usage",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: metricSysGoAllocBytes, Name: "go", Div: 1024},
			{ID: metricSysCGoAllocBytes, Name: "cgo", Div: 1024},
		},
	}
	chartGoroutines = Chart{
		ID:    "goroutines_count",
		Title: "Goroutines Count",
		Units: "goroutines",
		Fam:   "go/cgo",
		Ctx:   "cockroachdb.goroutines",
		Dims: Dims{
			{ID: metricSysGoroutines, Name: "goroutines"},
		},
	}
	chartGCRuns = Chart{
		ID:    "gc_count",
		Title: "GC Runs",
		Units: "invokes/s",
		Fam:   "go/cgo",
		Ctx:   "cockroachdb.gc_count",
		Dims: Dims{
			{ID: metricSysGCCount, Name: "gc", Algo: module.Incremental},
		},
	}
	chartGCPauseTime = Chart{
		ID:    "gc_pause",
		Title: "GC Pause Time",
		Units: "us",
		Fam:   "go/cgo",
		Ctx:   "cockroachdb.gc_pause",
		Dims: Dims{
			{ID: metricSysGCPauseNs, Name: "pause", Algo: module.Incremental, Div: 1e3},
		},
	}
	chartCGoCalls = Chart{
		ID:    "cgo_calls",
		Title: "CGo Calls",
		Units: "calls",
		Fam:   "go/cgo",
		Ctx:   "cockroachdb.cgo_calls",
		Dims: Dims{
			{ID: metricSysCGoCalls, Name: "cgo", Algo: module.Incremental},
		},
	}
)

// SQL
var (
	chartSQLConnections = Chart{
		ID:    "sql_connections",
		Title: "Active SQL Connections",
		Units: "connections",
		Fam:   "sql",
		Ctx:   "cockroachdb.sql_connections",
		Dims: Dims{
			{ID: metricSQLConnections, Name: "active"},
		},
	}
	chartSQLTraffic = Chart{
		ID:    "sql_traffic",
		Title: "SQL Traffic",
		Units: "KiB",
		Fam:   "sql",
		Ctx:   "cockroachdb.sql_traffic",
		Type:  module.Area,
		Dims: Dims{
			{ID: metricSQLBytesIn, Name: "received", Div: 1024, Algo: module.Incremental},
			{ID: metricSQLBytesOut, Name: "sent", Div: -11024, Algo: module.Incremental},
		},
	}
	chartSQLQueries = Chart{
		ID:    "sql_queries",
		Title: "SQL Queries Successfully Executed",
		Units: "queries",
		Fam:   "sql",
		Ctx:   "cockroachdb.sql_queries",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: metricSQLSelectCount, Name: "SELECT", Algo: module.Incremental},
			{ID: metricSQLUpdateCount, Name: "UPDATE", Algo: module.Incremental},
			{ID: metricSQLInsertCount, Name: "INSERT", Algo: module.Incremental},
			{ID: metricSQLDeleteCount, Name: "DELETE", Algo: module.Incremental},
		},
	}
	chartSQLErroredQueries = Chart{
		ID:    "sql_errored_queries",
		Title: "SQL Queries Resulting in a Planning or Runtime Error",
		Units: "errors",
		Fam:   "sql",
		Ctx:   "cockroachdb.sql_errored_queries",
		Dims: Dims{
			{ID: metricSQLFailureCount, Name: "errors", Algo: module.Incremental},
		},
	}
	chartSQLActiveDistQueries = Chart{
		ID:    "sql_active_distributed_queries",
		Title: "Number of Distributed SQL Queries Currently Active",
		Units: "queries",
		Fam:   "sql",
		Ctx:   "cockroachdb.sql_active_distributed_queries",
		Dims: Dims{
			{ID: metricSQLDistSQLQueriesActive, Name: "active"},
		},
	}
	chartSQLActiveFlowsForDistQueries = Chart{
		ID:    "sql_active_distributed_flows",
		Title: "Number of Distributed SQL Flows Currently Active",
		Units: "flows",
		Fam:   "sql",
		Ctx:   "cockroachdb.sql_active_distributed_flows",
		Dims: Dims{
			{ID: metricSQLDistSQLFlowsActive, Name: "active"},
		},
	}
	chartSQLTransactions = Chart{
		ID:    "sql_transactions",
		Title: "SQL Transactions Successfully Executed",
		Units: "transactions",
		Fam:   "sql",
		Ctx:   "cockroachdb.sql_queries",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: metricSQLTXNBeginCount, Name: "BEGIN"},
			{ID: metricSQLTXNCommitCount, Name: "COMMIT"},
			{ID: metricSQLTXNRollbackCount, Name: "ROLLBACK"},
			{ID: metricSQLTXNAbortCount, Name: "ABORT"},
		},
	}
	chartSQLSchemaChanges = Chart{
		ID:    "sql_schema_changes",
		Title: "SQL DDL Statements Successfully Executed",
		Units: "queries",
		Fam:   "sql",
		Ctx:   "cockroachdb.sql_schema_changes",
		Dims: Dims{
			{ID: metricSQLDDLCount, Name: "DDL"},
		},
	}
)

var (
	chartProcessCPUUsage = Chart{
		ID:    "process_cpu_time",
		Title: "Process CPU Time",
		Units: "ms",
		Fam:   "process statistics",
		Ctx:   "cockroachdb.process_cpu_time",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: metricSysCPUUserNs, Name: "user", Algo: module.Incremental, Div: 1e6},
			{ID: metricSysCPUSysNs, Name: "sys", Algo: module.Incremental, Div: 1e6},
		},
	}
	chartProcessCPUUsagePercent = Chart{
		ID:    "process_cpu_time_percentage",
		Title: "Process CPU Time Percentage",
		Units: "percentage",
		Fam:   "process statistics",
		Ctx:   "cockroachdb.process_cpu_time_percentage",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: metricSysCPUUserPercent, Name: "user", Div: precision},
			{ID: metricSysCPUSysPercent, Name: "sys", Div: precision},
		},
	}
	chartProcessCPUUsageCombinedPercent = Chart{
		ID:    "process_cpu_time_combined_percentage",
		Title: "Process CPU Time Percentage",
		Units: "percentage",
		Fam:   "process statistics",
		Ctx:   "cockroachdb.process_cpu_time_combined_percentage",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: metricSysCPUCombinedPercentNormalized, Name: "usage", Div: precision},
		},
	}
	chartProcessMemory = Chart{
		ID:    "process_memory",
		Title: "Process Memory Usage",
		Units: "KiB",
		Fam:   "process statistics",
		Ctx:   "cockroachdb.process_memory",
		Dims: Dims{
			{ID: metricSysRSS, Name: "rss", Div: 1024},
		},
	}
	chartProcessFDUsage = Chart{
		ID:    "process_file_descriptors",
		Title: "File Descriptors Statistics",
		Units: "file descriptors",
		Fam:   "process statistics",
		Ctx:   "cockroachdb.process_file_descriptors",
		Dims: Dims{
			{ID: metricSysFDOpen, Name: "open"},
		},
		Vars: Vars{
			{ID: metricSysFDSoftLimit},
		},
	}
)

var chartRPCClockOffset = Chart{
	ID:    "mean_clock_offset",
	Title: "Mean Clock Offset Against With Other Nodes",
	Units: "us",
	Fam:   "rpc",
	Ctx:   "cockroachdb.mean_clock_offset",
	Dims: Dims{
		{ID: metricClockOffsetMeanNs, Name: "clock offset", Div: 1e3},
	},
}

var chartUptime = Chart{
	ID:    "system_uptime",
	Title: "Nodes",
	Units: "seconds",
	Fam:   "uptime",
	Ctx:   "cockroachdb.uptime",
	Dims: Dims{
		{ID: metricSysUptime, Name: "uptime"},
	},
}

var (
	chartHostDiskBandwidth = Chart{
		ID:    "host_disk_bandwidth",
		Title: "Host Disk Bandwidth",
		Units: "KiB/s",
		Fam:   "host statistics",
		Ctx:   "cockroachdb.host_disk_bandwidth",
		Dims: Dims{
			{ID: metricSysHostDiskReadBytes, Name: "read", Div: 1024, Algo: module.Incremental},
			{ID: metricSysHostDiskWriteBytes, Name: "write", Div: -1024, Algo: module.Incremental},
		},
	}
	chartHostDiskOperations = Chart{
		ID:    "host_disk_operations",
		Title: "Host Disk Operations",
		Units: "operations/s",
		Fam:   "host statistics",
		Ctx:   "cockroachdb.host_disk_operations",
		Dims: Dims{
			{ID: metricSysHostDiskReadCount, Name: "reads", Algo: module.Incremental},
			{ID: metricSysHostDiskWriteCount, Name: "writes", Mul: -1, Algo: module.Incremental},
		},
	}
	chartHostDiskIOPS = Chart{
		ID:    "host_disk_iops_in_progress",
		Title: "Host Disk IOPS In Progress",
		Units: "iops",
		Fam:   "host statistics",
		Ctx:   "cockroachdb.host_disk_iops_in_progress",
		Dims: Dims{
			{ID: metricSysHostDiskIOPSInProgress, Name: "in progress"},
		},
	}
	chartHostNetworkBandwidth = Chart{
		ID:    "host_network_bandwidth",
		Title: "Host Network Bandwidth",
		Units: "kilobits/s",
		Fam:   "host statistics",
		Ctx:   "cockroachdb.host_network_bandwidth",
		Dims: Dims{
			{ID: metricSysHostNetRecvBytes, Name: "received", Div: 1000, Algo: module.Incremental},
			{ID: metricSysHostNetSendBytes, Name: "sent", Div: -1000, Algo: module.Incremental},
		},
	}
	chartHostNetworkPackets = Chart{
		ID:    "host_network_packets",
		Title: "Host Network Packets",
		Units: "pps",
		Fam:   "host statistics",
		Ctx:   "cockroachdb.host_network_packets",
		Dims: Dims{
			{ID: metricSysHostNetRecvPackets, Name: "received", Algo: module.Incremental},
			{ID: metricSysHostNetSendPackets, Name: "sent", Algo: module.Incremental},
		},
	}
)

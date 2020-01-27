package cockroachdb

const (
	// Storage Dashboard
	// https://github.com/cockroachdb/cockroach/blob/master/pkg/storage/metrics.go
	metricCapacity                 = "capacity"
	metricCapacityAvailable        = "capacity_available"
	metricCapacityUsed             = "capacity_used"
	metricCapacityReserved         = "capacity_reserved"
	metricLiveBytes                = "livebytes"
	metricSysBytes                 = "sysbytes"
	metricRocksDBReadAmplification = "rocksdb_read_amplification"
	metricRocksDBNumSSTables       = "rocksdb_num_sstables"
	metricRocksDBBlockCacheUsage   = "rocksdb_block_cache_usage"
	metricRocksDBBlockCacheHits    = "rocksdb_block_cache_hits"
	metricRocksDBBlockCacheMisses  = "rocksdb_block_cache_misses"
	metricRocksDBCompactions       = "rocksdb_compactions"
	metricRocksDBFlushes           = "rocksdb_flushes"
	metricSysFDOpen                = "sys_fd_open"
	metricSysFDSoftLimit           = "sys_fd_softlimit"
	// https://github.com/cockroachdb/cockroach/blob/master/pkg/ts/metrics.go
	metricTimeSeriesWriteSamples = "timeseries_write_samples"
	metricTimeSeriesWriteErrors  = "timeseries_write_errors"
	metricTimeSeriesWriteBytes   = "timeseries_write_bytes"

	// Runtime Dashboard
	// https://github.com/cockroachdb/cockroach/blob/master/pkg/server/status/runtime.go
	metricLiveNodes         = "liveness_livenodes"
	metricSysUptime         = "sys_uptime"
	metricSysRSS            = "sys_rss"
	metricSysGoAllocBytes   = "sys_go_allocbytes"
	metricSysGoTotalBytes   = "sys_go_totalbytes"
	metricSysCGoAllocBytes  = "sys_cgo_allocbytes"
	metricSysCGoTotalBytes  = "sys_cgo_totalbytes"
	metricSysCGoCalls       = "sys_cgocalls"
	metricSysGoroutines     = "sys_goroutines"
	metricSysGCCount        = "sys_gc_count"
	metricSysGCPauseNs      = "sys_gc_pause_ns"
	metricSysCPUUserNs      = "sys_cpu_user_ns"
	metricSysCPUSysNs       = "sys_cpu_sys_ns"
	metricClockOffsetMeanNs = "clock_offset_meannanos"

	// SQL Dashboard
	// https://github.com/cockroachdb/cockroach/blob/master/pkg/sql/pgwire/server.go
	metricSQLConnections = "sql_conns"
	metricSQLBytesIn     = "sql_bytesin"
	metricSQLBytesOut    = "sql_bytesout"
	// https://github.com/cockroachdb/cockroach/blob/master/pkg/sql/exec_util.go
	metricSQLSelectCount          = "sql_select_count"
	metricSQLUpdateCount          = "sql_update_count"
	metricSQLInsertCount          = "sql_insert_count"
	metricSQLDeleteCount          = "sql_delete_count"
	metricSQLFailureCount         = "sql_failure_count"
	metricSQLTXNBeginCount        = "sql_txn_begin_count"
	metricSQLTXNAbortCount        = "sql_txn_abort_count"
	metricSQLTXNCommitCount       = "sql_txn_commit_count"
	metricSQLTXNRollbackCount     = "sql_txn_rollback_count"
	metricSQLDistSQLQueriesActive = "sql_distsql_queries_active"
	metricSQLDistSQLFlowsActive   = "sql_distsql_flows_active"
	metricSQLDDLCount             = "sql_ddl_count"

	// Replication Dashboard
	// https://github.com/cockroachdb/cockroach/blob/master/pkg/storage/metrics.go
	metricRanges                          = "ranges"
	metricRangesUnavailable               = "ranges_unavailable"
	metricRangesUnderReplicated           = "ranges_underreplicated"
	metricRangesOverReplicated            = "ranges_overreplicated"
	metricRangeSplits                     = "range_splits"
	metricRangeAdds                       = "range_adds"
	metricRangeRemoves                    = "range_removes"
	metricRangeMerges                     = "range_merges"
	metricRangeSnapshotsGenerated         = "range_snapshots_generated"
	metricRangeSnapshotsPreemptiveApplied = "range_snapshots_preemptive_applied"
	metricRangeSnapshotsLearnerApplied    = "range_snapshots_learner_applied"
	metricRangeSnapshotsNormalApplied     = "range_snapshots_normal_applied"
	metricReplicas                        = "replicas"
	metricReplicasReserved                = "replicas_reserved"
	metricReplicasLeaders                 = "replicas_leaders"
	metricReplicasLeadersNotLeaseholders  = "replicas_leaders_not_leaseholders"
	metricReplicasLeaseholders            = "replicas_leaseholders"
	metricReplicasQuiescent               = "replicas_quiescent"
	metricKeyBytes                        = "keybytes"
	metricValBytes                        = "valbytes"
	metricRebalancingQueriesPerSecond     = "rebalancing_queriespersecond"
	metricRebalancingWritesPerSecond      = "rebalancing_writespersecond"

	// Dashboard Slow Requests
	// https://github.com/cockroachdb/cockroach/blob/master/pkg/storage/metrics.go
	metricRequestsSlowLease = "requests_slow_lease"
	metricRequestsSlowLatch = "requests_slow_latch"
	metricRequestsSlowRaft  = "requests_slow_raft"

	// Dashboard Hardware
	// https://github.com/cockroachdb/cockroach/blob/master/pkg/server/status/runtime.go
	metricSysHostDiskReadBytes      = "sys_host_disk_read_bytes"
	metricSysHostDiskWriteBytes     = "sys_host_disk_write_bytes"
	metricSysHostDiskReadCount      = "sys_host_disk_read_count"
	metricSysHostDiskWriteCount     = "sys_host_disk_write_count"
	metricSysHostDiskIOPSInProgress = "sys_host_disk_iopsinprogress"
	metricSysHostNetSendBytes       = "sys_host_net_send_bytes"
	metricSysHostNetRecvBytes       = "sys_host_net_recv_bytes"
	metricSysHostNetSendPackets     = "sys_host_net_send_packets"
	metricSysHostNetRecvPackets     = "sys_host_net_recv_packets"
)

const (
	// calculated metrics
	metricCapacityUsable               = "capacity_usable"
	metricCapacityUnusable             = "capacity_unusable"
	metricCapacityUsedPercentage       = "capacity_used_percentage"
	metricCapacityUsableUsedPercentage = "capacity_usable_used_percentage"
	metricRocksDBBlockCacheHitRate     = "rocksdb_block_cache_hit_rate"
)

var metrics = []string{
	metricCapacity,
	metricCapacityAvailable,
	metricCapacityUsed,
	metricCapacityReserved,
	metricLiveBytes,
	metricSysBytes,
	metricRocksDBReadAmplification,
	metricRocksDBNumSSTables,
	metricRocksDBBlockCacheUsage,
	metricRocksDBBlockCacheHits,
	metricRocksDBBlockCacheMisses,
	metricRocksDBCompactions,
	metricRocksDBFlushes,
	metricSysFDOpen,
	metricSysFDSoftLimit,
	metricTimeSeriesWriteSamples,
	metricTimeSeriesWriteErrors,
	metricTimeSeriesWriteBytes,

	metricLiveNodes,
	metricSysUptime,
	metricSysRSS,
	metricSysGoAllocBytes,
	metricSysGoTotalBytes,
	metricSysCGoAllocBytes,
	metricSysCGoTotalBytes,
	metricSysCGoCalls,
	metricSysGoroutines,
	metricSysGCCount,
	metricSysGCPauseNs,
	metricSysCPUUserNs,
	metricSysCPUSysNs,
	metricClockOffsetMeanNs,

	metricSQLConnections,
	metricSQLBytesIn,
	metricSQLBytesOut,
	metricSQLSelectCount,
	metricSQLUpdateCount,
	metricSQLInsertCount,
	metricSQLDeleteCount,
	metricSQLFailureCount,
	metricSQLTXNAbortCount,
	metricSQLTXNBeginCount,
	metricSQLTXNCommitCount,
	metricSQLTXNRollbackCount,
	metricSQLDistSQLQueriesActive,
	metricSQLDistSQLFlowsActive,
	metricSQLDDLCount,

	metricRanges,
	metricRangesUnavailable,
	metricRangesUnderReplicated,
	metricRangesOverReplicated,
	metricRangeSplits,
	metricRangeAdds,
	metricRangeRemoves,
	metricRangeMerges,
	metricRangeSnapshotsGenerated,
	metricRangeSnapshotsPreemptiveApplied,
	metricRangeSnapshotsLearnerApplied,
	metricRangeSnapshotsNormalApplied,
	metricReplicas,
	metricReplicasReserved,
	metricReplicasLeaders,
	metricReplicasLeadersNotLeaseholders,
	metricReplicasLeaseholders,
	metricReplicasQuiescent,
	metricKeyBytes,
	metricValBytes,
	metricRebalancingQueriesPerSecond,
	metricRebalancingWritesPerSecond,

	metricRequestsSlowLease,
	metricRequestsSlowLatch,
	metricRequestsSlowRaft,
}

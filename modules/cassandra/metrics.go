// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

// https://cassandra.apache.org/doc/latest/cassandra/operating/metrics.html#table-metrics
// https://www.datadoghq.com/blog/how-to-collect-cassandra-metrics/

type cassandraMetrics struct {
	clientRequestTotalLatencyReads  *float64 `stm:"client_request_total_latency_reads,1000,1"`
	clientRequestTotalLatencyWrites *float64 `stm:"client_request_total_latency_writes,1000,1"`
	clientRequestLatencyReads       *float64 `stm:"client_request_latency_reads,1000,1"`
	clientRequestLatencyWrites      *float64 `stm:"client_request_latency_writes,1000,1"`
	clientRequestTimeoutsReads      *float64 `stm:"client_request_timeouts_reads,1000,1"`
	clientRequestTimeoutsWrites     *float64 `stm:"client_request_timeouts_writes,1000,1"`
	clientRequestUnavailablesReads  *float64 `stm:"client_request_unavailables_reads,1000,1"`
	clientRequestUnavailablesWrites *float64 `stm:"client_request_unavailables_writes,1000,1"`
	clientRequestFailuresReads      *float64 `stm:"client_request_failures_reads,1000,1"`
	clientRequestFailuresWrites     *float64 `stm:"client_request_failures_writes,1000,1"`

	cacheHits     *float64 `stm:"cache_hits,1000,1"`
	cacheMisses   *float64 `stm:"cache_misses,1000,1"`
	cacheHitRatio *float64 `stm:"cache_hit_ratio,1000,1"` // calculated
	cacheSize     *float64 `stm:"cache_size"`

	threadPoolsTotalBlockedTasks     *float64 `stm:"thread_pools_total_blocked_tasks"`
	threadPoolsCurrentlyBlockedTasks *float64 `stm:"thread_pools_currently_blocked_tasks"`

	// https://cassandra.apache.org/doc/latest/cassandra/operating/metrics.html#dropped-metrics
	droppedMsgsOneMinute *float64 `stm:"dropped_messages_one_minute,1000,1"`

	// https://cassandra.apache.org/doc/latest/cassandra/operating/metrics.html#storage-metrics
	storageLoad       *float64 `stm:"storage_load"`
	storageExceptions *float64 `stm:"storage_exceptions"`

	// https://cassandra.apache.org/doc/latest/cassandra/operating/metrics.html#compaction-metrics
	compactionBytesCompacted *float64 `stm:"compaction_bytes_compacted"`
	compactionPendingTasks   *float64 `stm:"compaction_pending_tasks"`
	compactionCompletedTasks *float64 `stm:"compaction_completed_tasks"`

	// https://cassandra.apache.org/doc/latest/cassandra/operating/metrics.html#garbagecollector
	jvmGCParNewCount *float64 `stm:"jvm_gc_parnew_count,1000,1"`
	jvmGCParNewTime  *float64 `stm:"jvm_gc_parnew_time,1000,1"`
	jvmGCCMSCount    *float64 `stm:"jvm_gc_cms_count,1000,1"`
	jvmGCCMSTime     *float64 `stm:"jvm_gc_cms_time,1000,1"`
}

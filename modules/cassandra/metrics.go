// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

// https://cassandra.apache.org/doc/latest/cassandra/operating/metrics.html#table-metrics
// https://www.datadoghq.com/blog/how-to-collect-cassandra-metrics/
// https://docs.opennms.com/horizon/29/deployment/time-series-storage/newts/cassandra-jmx.html

type cassandraMetrics struct {
	clientRequestTotalLatencyReads  *float64 `stm:"client_request_total_latency_reads"`
	clientRequestTotalLatencyWrites *float64 `stm:"client_request_total_latency_writes"`
	clientRequestLatencyReads       *float64 `stm:"client_request_latency_reads"`
	clientRequestLatencyWrites      *float64 `stm:"client_request_latency_writes"`
	clientRequestTimeoutsReads      *float64 `stm:"client_request_timeouts_reads"`
	clientRequestTimeoutsWrites     *float64 `stm:"client_request_timeouts_writes"`
	clientRequestUnavailablesReads  *float64 `stm:"client_request_unavailables_reads"`
	clientRequestUnavailablesWrites *float64 `stm:"client_request_unavailables_writes"`
	clientRequestFailuresReads      *float64 `stm:"client_request_failures_reads"`
	clientRequestFailuresWrites     *float64 `stm:"client_request_failures_writes"`

	cacheHits     *float64 `stm:"cache_hits"`
	cacheMisses   *float64 `stm:"cache_misses"`
	cacheHitRatio *float64 `stm:"cache_hit_ratio,1000,1"` // calculated
	cacheSize     *float64 `stm:"cache_size"`

	threadPoolsActiveTasks           *float64 `stm:"thread_pools_active_tasks"`
	threadPoolsPendingTasks          *float64 `stm:"thread_pools_pending_tasks"`
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
	jvmGCParNewCount *float64 `stm:"jvm_gc_parnew_count"`
	jvmGCParNewTime  *float64 `stm:"jvm_gc_parnew_time,1000,1"`
	jvmGCCMSCount    *float64 `stm:"jvm_gc_cms_count"`
	jvmGCCMSTime     *float64 `stm:"jvm_gc_cms_time,1000,1"`
}

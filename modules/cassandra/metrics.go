// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

type metrics struct {
	throughput *THROUGHPUT               `stm:"org_apache_cassandra_metrics_clientrequest_oneminuterate"`
	latency    *LATENCY                  `stm:"org_apache_cassandra_metrics_clientrequest_count"`
	cache      *CACHE                    `stm:"org_apache_cassandra_metrics_cache_count"`
	disk       *DISK                     `stm:"org_apache_cassandra_metrics_table_count"`
	gcc        *GARBAGE_COLLECTION_COUNT `stm:"jvm_gc_collection_seconds_count"`
	gct        *GARBAGE_COLLECTION_TIME  `stm:"jvm_gc_collection_seconds_sum"`
}

func (c metrics) hasThrouput() bool { return c.throughput != nil }
func (c metrics) hasLatency() bool  { return c.latency != nil }
func (c metrics) hasCache() bool    { return c.cache != nil }
func (c metrics) hasDisk() bool     { return c.disk != nil }

type (
	THROUGHPUT struct {
		read  int64 `stm:"Read"`
		write int64 `stm:"Write"`
	}
	LATENCY struct {
		read  int64 `stm:"ReadLatency"`
		write int64 `stm:"WriteLatency"`
	}
	CACHE struct {
		hit int64 `stm:"HitRate"`
	}
	DISK struct {
		load                 float64 `stm:"LiveDiskSpaceUsed"`
		used                 float64 `stm:"TotalDiskSpaceUsed"`
		compaction_completed float64 `stm:"CompactionBytesWritten"`
		compaction_queue     float64 `stm:"PendingCompactions"`
	}
	GARBAGE_COLLECTION_COUNT struct {
		parNewCount    int64 `stm:"ParNew"`
		markSweepCount int64 `stm:"ConcurrentMarkSweep"`
	}
	GARBAGE_COLLECTION_TIME struct {
		parNewTime    int64 `stm:"ParNew"`
		markSweepTime int64 `stm:"ConcurrentMarkSweep"`
	}
)

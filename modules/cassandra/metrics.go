// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

type metrics struct {
	throughput   *THROUGHPUT         `stm:"throuput"`
	latency      *LATENCY            `stm:"latency"`
	cache        *CACHE              `stm:"cache"`
	disk         *DISK               `stm:"disk"`
	gc_count     *GARBAGE_COLLECTION `stm:"java_gc_count"`
	gc_time      *GARBAGE_COLLECTION `stm:"java_gc_time"`
	etimeout     *REQUEST_ERROR      `stm:"error_timeout"`
	eunavailable *REQUEST_ERROR      `stm:"error_unavailable"`
	pending_task *PENDING_TASK       `stm:"pending_tasks"`
	blocked_task *BLOCKED_TASK       `stm:"blocked_tasks"`
}

func (c metrics) hasThrouput() bool { return c.throughput != nil }
func (c metrics) hasLatency() bool  { return c.latency != nil }
func (c metrics) hasCache() bool    { return c.cache != nil }
func (c metrics) hasDisk() bool     { return c.disk != nil }

const (
	metricRequestType = "org_apache_cassandra_metrics_clientrequest_count"
)

type (
	THROUGHPUT struct {
		read  int64 `stm:"Read"`
		write int64 `stm:"Write"`
	}
	LATENCY struct {
		read_latency  int64 `stm:"ReadLatency"`
		write_latency int64 `stm:"WriteLatency"`
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
	GARBAGE_COLLECTION struct {
		parNew    int64 `stm:"ParNew"`
		markSweep int64 `stm:"ConcurrentMarkSweep"`
	}
	REQUEST_ERROR struct {
		read_error  int64 `stm:"Read"`
		write_error int64 `stm:"Write"`
	}
	PENDING_TASK struct {
		task int64 `stm:"tasks"`
	}
	BLOCKED_TASK struct {
		task int64 `stm:"CurrentlyBlockedTasks"`
	}
)

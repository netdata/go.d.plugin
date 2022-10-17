// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

type metrics struct {
	ioThroughput *throughput        `stm:"throughput"`
	ioLatency    *latency           `stm:"latency"`
	hit          *cache             `stm:"cache"`
	hd           *disk              `stm:"disk"`
	gcCount      *garbageCollection `stm:"java_gc_count"`
	gcTime       *garbageCollection `stm:"java_gc_time"`
	etimeout     *requestError      `stm:"error_timeout"`
	eunavailable *requestError      `stm:"error_unavailable"`
	pTask        *pendingTask       `stm:"pending_tasks"`
	bTask        *blockedTask       `stm:"blocked_tasks"`
}

const (
	metricRequestType = "org_apache_cassandra_metrics_clientrequest_count"
)

type (
	throughput struct {
		read  int64 `stm:"Read"`
		write int64 `stm:"Write"`
	}
	latency struct {
		read_latency  int64 `stm:"Read"`
		write_latency int64 `stm:"Write"`
	}
	cache struct {
		hit int64 `stm:"HitRate"`
	}
	disk struct {
		load                 float64 `stm:"LiveDiskSpaceUsed"`
		used                 float64 `stm:"TotalDiskSpaceUsed"`
		compaction_completed float64 `stm:"CompactionBytesWritten"`
		compaction_queue     float64 `stm:"PendingCompactions"`
	}
	garbageCollection struct {
		parNew    int64 `stm:"ParNew"`
		markSweep int64 `stm:"ConcurrentMarkSweep"`
	}
	requestError struct {
		read_error  int64 `stm:"Read"`
		write_error int64 `stm:"Write"`
	}
	pendingTask struct {
		task int64 `stm:"tasks"`
	}
	blockedTask struct {
		task int64 `stm:"CurrentlyBlockedTasks"`
	}
)

// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

type metrics struct {
	throughput *THROUGHPUT `stm:"org_apache_cassandra_metrics_clientrequest_oneminuterate"`
	latency    *LATENCY    `stm:"org_apache_cassandra_metrics_table_count"`
	cache      *CACHE      `stm:"org_apache_cassandra_metrics_table_count"`
	disk       *DISK       `stm:"org_apache_cassandra_metrics_table_count"`
}

const (
	metricTableType = "org_apache_cassandra_metrics_table_count"
)

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
)

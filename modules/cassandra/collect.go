// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

import (
	"errors"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	suffixCount     = "_count"
	suffixValue     = "_value"
	suffixOneMinute = "_oneminuterate"
)

func (c *Cassandra) collect() (map[string]int64, error) {
	pms, err := c.prom.Scrape()
	if err != nil {
		return nil, err
	}

	if c.validateMetrics {
		if !isCassandraMetrics(pms) {
			return nil, errors.New("collected metrics aren't Cassandra metrics")
		}
		c.validateMetrics = false
	}

	mx := make(map[string]int64)

	c.resetMetrics()
	c.collectMetrics(pms)
	c.processMetric(mx)

	return mx, nil
}

func (c *Cassandra) resetMetrics() {
	cm := newCassandraMetrics()
	for key, p := range c.mx.threadPools {
		cm.threadPools[key] = &threadPoolMetrics{
			name:      p.name,
			hasCharts: p.hasCharts,
		}
	}
	c.mx = cm
}

func (c *Cassandra) processMetric(mx map[string]int64) {
	c.mx.clientReqTotalLatencyReads.write(mx, "client_request_total_latency_reads")
	c.mx.clientReqTotalLatencyWrites.write(mx, "client_request_total_latency_writes")
	c.mx.clientReqLatencyReads.write(mx, "client_request_latency_reads")
	c.mx.clientReqLatencyWrites.write(mx, "client_request_latency_writes")
	c.mx.clientReqTimeoutsReads.write(mx, "client_request_timeouts_reads")
	c.mx.clientReqTimeoutsWrites.write(mx, "client_request_timeouts_writes")
	c.mx.clientReqUnavailablesReads.write(mx, "client_request_unavailables_reads")
	c.mx.clientReqUnavailablesWrites.write(mx, "client_request_unavailables_writes")
	c.mx.clientReqFailuresReads.write(mx, "client_request_failures_reads")
	c.mx.clientReqFailuresWrites.write(mx, "client_request_failures_writes")

	c.mx.keyCacheHits.write(mx, "key_cache_hits")
	c.mx.keyCacheMisses.write(mx, "key_cache_misses")
	c.mx.keyCacheSize.write(mx, "key_cache_size")
	if c.mx.keyCacheHits.isSet && c.mx.keyCacheMisses.isSet {
		if s := c.mx.keyCacheHits.value + c.mx.keyCacheMisses.value; s > 0 {
			mx["key_cache_hit_ratio"] = int64((c.mx.keyCacheHits.value * 100 / s) * 1000)
		} else {
			mx["key_cache_hit_ratio"] = 0
		}
	}
	if c.mx.keyCacheCapacity.isSet && c.mx.keyCacheSize.isSet {
		if s := c.mx.keyCacheCapacity.value; s > 0 {
			mx["key_cache_utilization"] = int64((c.mx.keyCacheSize.value * 100 / s) * 1000)
		} else {
			mx["key_cache_utilization"] = 0
		}
	}

	c.mx.droppedMsgsOneMinute.write1k(mx, "dropped_messages_one_minute")

	c.mx.storageLoad.write(mx, "storage_load")
	c.mx.storageExceptions.write(mx, "storage_exceptions")

	c.mx.compactionBytesCompacted.write(mx, "compaction_bytes_compacted")
	c.mx.compactionPendingTasks.write(mx, "compaction_pending_tasks")
	c.mx.compactionCompletedTasks.write(mx, "compaction_completed_tasks")

	c.mx.jvmGCParNewCount.write(mx, "jvm_gc_parnew_count")
	c.mx.jvmGCParNewTime.write1k(mx, "jvm_gc_parnew_time")
	c.mx.jvmGCCMSCount.write(mx, "jvm_gc_cms_count")
	c.mx.jvmGCCMSTime.write1k(mx, "jvm_gc_cms_time")

	for _, p := range c.mx.threadPools {
		if !p.hasCharts {
			p.hasCharts = true
			c.addThreadPoolCharts(p)
		}

		px := "thread_pool_" + p.name + "_"
		p.activeTasks.write(mx, px+"active_tasks")
		p.pendingTasks.write(mx, px+"pending_tasks")
		p.blockedTasks.write(mx, px+"blocked_tasks")
		p.totalBlockedTasks.write(mx, px+"total_blocked_tasks")
	}
}

func (c *Cassandra) collectMetrics(pms prometheus.Metrics) {
	c.collectClientRequestMetrics(pms)
	c.collectDroppedMessagesMetrics(pms)
	c.collectThreadPoolsMetrics(pms)
	c.collectStorageMetrics(pms)
	c.collectCacheMetrics(pms)
	c.collectJVMMetrics(pms)
	c.collectCompactionMetrics(pms)
}

func (c *Cassandra) collectClientRequestMetrics(pms prometheus.Metrics) {
	const metric = "org_apache_cassandra_metrics_clientrequest"

	var rw struct{ r, w *metricValue }
	for _, pm := range pms.FindByName(metric + suffixCount) {
		name := pm.Labels.Get("name")
		scope := pm.Labels.Get("scope")

		switch name {
		case "TotalLatency":
			rw.r, rw.w = &c.mx.clientReqTotalLatencyReads, &c.mx.clientReqTotalLatencyWrites
		case "Latency":
			rw.r, rw.w = &c.mx.clientReqLatencyReads, &c.mx.clientReqLatencyWrites
		case "Timeouts":
			rw.r, rw.w = &c.mx.clientReqTimeoutsReads, &c.mx.clientReqTimeoutsWrites
		case "Unavailables":
			rw.r, rw.w = &c.mx.clientReqUnavailablesReads, &c.mx.clientReqUnavailablesWrites
		case "Failures":
			rw.r, rw.w = &c.mx.clientReqFailuresReads, &c.mx.clientReqFailuresWrites
		default:
			continue
		}

		switch scope {
		case "Read":
			rw.r.add(pm.Value)
		case "Write":
			rw.w.add(pm.Value)
		}
	}
}

func (c *Cassandra) collectCacheMetrics(pms prometheus.Metrics) {
	const metric = "org_apache_cassandra_metrics_cache"

	for _, pm := range pms.FindByName(metric + suffixCount) {
		name := pm.Labels.Get("name")
		scope := pm.Labels.Get("scope")
		if scope != "KeyCache" {
			continue
		}

		switch name {
		case "Hits":
			c.mx.keyCacheHits.add(pm.Value)
		case "Misses":
			c.mx.keyCacheMisses.add(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metric + suffixValue) {
		name := pm.Labels.Get("name")

		switch name {
		case "Capacity":
			c.mx.keyCacheCapacity.add(pm.Value)
		case "Size":
			c.mx.keyCacheSize.add(pm.Value)
		}
	}
}

func (c *Cassandra) collectThreadPoolsMetrics(pms prometheus.Metrics) {
	const metric = "org_apache_cassandra_metrics_threadpools"

	for _, pm := range pms.FindByName(metric + suffixValue) {
		name := pm.Labels.Get("name")
		scope := pm.Labels.Get("scope")
		pool := c.getThreadPoolMetrics(scope)

		switch name {
		case "ActiveTasks":
			pool.activeTasks.add(pm.Value)
		case "PendingTasks":
			pool.pendingTasks.add(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metric + suffixCount) {
		name := pm.Labels.Get("name")
		scope := pm.Labels.Get("scope")
		pool := c.getThreadPoolMetrics(scope)

		switch name {
		case "CompletedTasks":
			pool.totalBlockedTasks.add(pm.Value)
		case "TotalBlockedTasks":
			pool.totalBlockedTasks.add(pm.Value)
		case "CurrentlyBlockedTasks":
			pool.blockedTasks.add(pm.Value)
		}
	}
}

func (c *Cassandra) collectStorageMetrics(pms prometheus.Metrics) {
	const metric = "org_apache_cassandra_metrics_storage"

	for _, pm := range pms.FindByName(metric + suffixCount) {
		name := pm.Labels.Get("name")

		switch name {
		case "Load":
			c.mx.storageLoad.add(pm.Value)
		case "Exceptions":
			c.mx.storageExceptions.add(pm.Value)
		}
	}
}

func (c *Cassandra) collectDroppedMessagesMetrics(pms prometheus.Metrics) {
	const metric = "org_apache_cassandra_metrics_droppedmessage"

	for _, pm := range pms.FindByName(metric + suffixOneMinute) {
		c.mx.droppedMsgsOneMinute.add(pm.Value)
	}
}

func (c *Cassandra) collectJVMMetrics(pms prometheus.Metrics) {
	const metric = "jvm_gc_collection_seconds"

	for _, pm := range pms.FindByName(metric + suffixCount) {
		gc := pm.Labels.Get("gc")

		switch gc {
		case "ParNew":
			c.mx.jvmGCParNewCount.add(pm.Value)
		case "ConcurrentMarkSweep":
			c.mx.jvmGCCMSCount.add(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metric + "_sum") {
		gc := pm.Labels.Get("gc")

		switch gc {
		case "ParNew":
			c.mx.jvmGCParNewTime.add(pm.Value)
		case "ConcurrentMarkSweep":
			c.mx.jvmGCCMSTime.add(pm.Value)
		}
	}
}

func (c *Cassandra) collectCompactionMetrics(pms prometheus.Metrics) {
	const metric = "org_apache_cassandra_metrics_compaction"

	for _, pm := range pms.FindByName(metric + suffixValue) {
		name := pm.Labels.Get("name")

		switch name {
		case "CompletedTasks":
			c.mx.compactionCompletedTasks.add(pm.Value)
		case "PendingTasks":
			c.mx.compactionPendingTasks.add(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metric + suffixCount) {
		name := pm.Labels.Get("name")

		switch name {
		case "BytesCompacted":
			c.mx.compactionBytesCompacted.add(pm.Value)
		}
	}
}

func (c *Cassandra) getThreadPoolMetrics(name string) *threadPoolMetrics {
	pool, ok := c.mx.threadPools[name]
	if !ok {
		pool = &threadPoolMetrics{name: name}
		c.mx.threadPools[name] = pool
	}
	return pool
}

func isCassandraMetrics(pms prometheus.Metrics) bool {
	for _, pm := range pms {
		if strings.HasPrefix(pm.Name(), "org_apache_cassandra_metrics") {
			return true
		}
	}
	return false
}

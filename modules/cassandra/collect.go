// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

import (
	"errors"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/stm"
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

	var mx cassandraMetrics
	c.collectMetrics(&mx, pms)

	if mx.cacheMisses != nil && mx.cacheHits != nil {
		var hitRatio float64
		if total := *mx.cacheMisses + *mx.cacheHits; total > 0 {
			hitRatio = *mx.cacheHits * 100 / total
		}
		mx.cacheHitRatio = &hitRatio
	}

	return stm.ToMap(mx), nil
}

func (c *Cassandra) collectMetrics(mx *cassandraMetrics, pms prometheus.Metrics) {
	c.collectClientRequestMetrics(mx, pms)
	c.collectDroppedMessagesMetrics(mx, pms)
	c.collectThreadPoolsMetrics(mx, pms)
	c.collectStorageMetrics(mx, pms)
	c.collectMetricsCacheMetrics(mx, pms)
	c.collectJVMMetrics(mx, pms)
	c.collectCompactionMetrics(mx, pms)
}

func (c *Cassandra) collectClientRequestMetrics(mx *cassandraMetrics, pms prometheus.Metrics) {
	const metric = "org_apache_cassandra_metrics_clientrequest"

	var rw struct{ r, w **float64 }
	for _, pm := range pms.FindByName(metric + suffixCount) {
		name := pm.Labels.Get("name")
		scope := pm.Labels.Get("scope")

		switch name {
		case "TotalLatency":
			rw.r, rw.w = &mx.clientRequestTotalLatencyReads, &mx.clientRequestTotalLatencyWrites
		case "Latency":
			rw.r, rw.w = &mx.clientRequestLatencyReads, &mx.clientRequestLatencyWrites
		case "Timeouts":
			rw.r, rw.w = &mx.clientRequestTimeoutsReads, &mx.clientRequestTimeoutsWrites
		case "Unavailables":
			rw.r, rw.w = &mx.clientRequestUnavailablesReads, &mx.clientRequestUnavailablesWrites
		case "Failures":
			rw.r, rw.w = &mx.clientRequestFailuresReads, &mx.clientRequestFailuresWrites
		default:
			continue
		}

		switch scope {
		case "Read":
			addValue(rw.r, pm.Value)
		case "Write":
			addValue(rw.w, pm.Value)
		}

	}
}

func (c *Cassandra) collectMetricsCacheMetrics(mx *cassandraMetrics, pms prometheus.Metrics) {
	const metric = "org_apache_cassandra_metrics_cache"

	for _, pm := range pms.FindByName(metric + suffixCount) {
		name := pm.Labels.Get("name")

		switch name {
		case "Misses":
			addValue(&mx.cacheMisses, pm.Value)
		case "Hits":
			addValue(&mx.cacheHits, pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metric + suffixValue) {
		name := pm.Labels.Get("name")

		switch name {
		case "Size":
			addValue(&mx.cacheSize, pm.Value)
		}
	}
}

func (c *Cassandra) collectThreadPoolsMetrics(mx *cassandraMetrics, pms prometheus.Metrics) {
	const metric = "org_apache_cassandra_metrics_threadpools"

	for _, pm := range pms.FindByName(metric + suffixCount) {
		name := pm.Labels.Get("name")

		switch name {
		case "TotalBlockedTasks":
			addValue(&mx.threadPoolsTotalBlockedTasks, pm.Value)
		case "CurrentlyBlockedTasks":
			addValue(&mx.threadPoolsCurrentlyBlockedTasks, pm.Value)
		}
	}
}

func (c *Cassandra) collectStorageMetrics(mx *cassandraMetrics, pms prometheus.Metrics) {
	const metric = "org_apache_cassandra_metrics_storage"

	for _, pm := range pms.FindByName(metric + suffixCount) {
		name := pm.Labels.Get("name")

		switch name {
		case "Load":
			addValue(&mx.storageLoad, pm.Value)
		case "Exceptions":
			addValue(&mx.storageExceptions, pm.Value)
		}
	}
}

func (c *Cassandra) collectDroppedMessagesMetrics(mx *cassandraMetrics, pms prometheus.Metrics) {
	const metric = "org_apache_cassandra_metrics_droppedmessage"

	for _, pm := range pms.FindByName(metric + suffixOneMinute) {
		addValue(&mx.droppedMsgsOneMinute, pm.Value)
	}
}

func (c *Cassandra) collectJVMMetrics(mx *cassandraMetrics, pms prometheus.Metrics) {
	const metric = "jvm_gc_collection_seconds"

	for _, pm := range pms.FindByName(metric + suffixCount) {
		gc := pm.Labels.Get("gc")

		switch gc {
		case "ParNew":
			addValue(&mx.jvmGCParNewCount, pm.Value)
		case "ConcurrentMarkSweep":
			addValue(&mx.jvmGCCMSCount, pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metric + "_sum") {
		gc := pm.Labels.Get("gc")

		switch gc {
		case "ParNew":
			addValue(&mx.jvmGCParNewTime, pm.Value)
		case "ConcurrentMarkSweep":
			addValue(&mx.jvmGCCMSTime, pm.Value)
		}
	}
}

func (c *Cassandra) collectCompactionMetrics(mx *cassandraMetrics, pms prometheus.Metrics) {
	const metric = "org_apache_cassandra_metrics_compaction"

	for _, pm := range pms.FindByName(metric + suffixValue) {
		name := pm.Labels.Get("name")

		switch name {
		case "CompletedTasks":
			addValue(&mx.compactionCompletedTasks, pm.Value)
		case "PendingTasks":
			addValue(&mx.compactionPendingTasks, pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metric + suffixCount) {
		name := pm.Labels.Get("name")

		switch name {
		case "BytesCompacted":
			addValue(&mx.compactionBytesCompacted, pm.Value)
		}
	}
}

func isCassandraMetrics(pms prometheus.Metrics) bool {
	for _, pm := range pms {
		if strings.HasPrefix(pm.Name(), "org_apache_cassandra_metrics") {
			return true
		}
	}
	return false
}

func addValue(current **float64, value float64) {
	if *current == nil {
		*current = &value
	} else {
		**current += value
	}
}

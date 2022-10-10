// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	collectorGCC   = "ParNew"
	metricGCCount = "jvm_gc_collection_seconds_count"
)

func doCollectGCCount(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, metricGCCount, collectorGCC, false)
	return enabled && success
}

func collectGCCount(pms prometheus.Metrics) *GARBAGE_COLLECTION_COUNT {
	if !doCollectGCCount(pms) {
		return nil
	}

	var gcc GARBAGE_COLLECTION_COUNT
	collectGCCountByType(&gcc, pms)

	return &gcc
}

func collectGCCountByType(gcc *GARBAGE_COLLECTION_COUNT, pms prometheus.Metrics) {
	for _, pm := range pms.FindByName(metricGCCount) {
		metricName := pm.Labels.Get("gc")
		assignGCCountMetric(gcc, metricName, pm.Value)
	}
}

func assignGCCountMetric(gcc *GARBAGE_COLLECTION_COUNT, scope string, value float64) {
	switch scope {
	default:
	case "ParNew":
		gcc.parNewCount = int64(value)
	case "ConcurrentMarkSweep":
		gcc.markSweepCount = int64(value)
	}
}

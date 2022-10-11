// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	collectorGCC  = "ParNew"
	metricGCCount = "jvm_gc_collection_seconds_count"
	metricGCTime  = "jvm_gc_collection_seconds_sum"
)

func doCollectGCCount(pms prometheus.Metrics, metric string) bool {
	var tester string
	if metric == metricGCCount {
		tester = metricGCCount
	} else {
		tester = metricGCTime
	}
	enabled, success := checkCollector(pms, tester, collectorGCC, true)
	return enabled && success
}

func collectGC(pms prometheus.Metrics, metric string) *garbageCollection {
	if !doCollectGCCount(pms, metric) {
		return nil
	}

	var gcc garbageCollection
	collectGCCountByType(&gcc, pms)

	return &gcc
}

func collectGCCountByType(gcc *garbageCollection, pms prometheus.Metrics) {
	for _, pm := range pms.FindByName(metricGCCount) {
		metricName := pm.Labels.Get("gc")
		assignGCCountMetric(gcc, metricName, pm.Value)
	}
}

func assignGCCountMetric(gcc *garbageCollection, scope string, value float64) {
	switch scope {
	default:
	case "ParNew":
		gcc.parNew = int64(value)
	case "ConcurrentMarkSweep":
		gcc.markSweep = int64(value)
	}
}

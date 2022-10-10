// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	collectorGCT   = "ParNew"
	metricGCTime = "jvm_gc_collection_seconds_sum"
)

func doCollectGCTime(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, metricGCCount, collectorGCT, false)
	return enabled && success
}

func collectGCTime(pms prometheus.Metrics) *GARBAGE_COLLECTION_TIME {
	if !doCollectGCCount(pms) {
		return nil
	}

	var gct GARBAGE_COLLECTION_TIME
	collectGCTimeByType(&gct, pms)

	return &gct
}

func collectGCTimeByType(gct *GARBAGE_COLLECTION_TIME, pms prometheus.Metrics) {
	for _, pm := range pms.FindByName(metricGCCount) {
		metricName := pm.Labels.Get("gc")
		assignGCTimeMetric(gct, metricName, pm.Value*100)
	}
}

func assignGCTimeMetric(gct *GARBAGE_COLLECTION_TIME, scope string, value float64) {
	switch scope {
	default:
	case "ParNew":
		gct.parNewTime = int64(value)
	case "ConcurrentMarkSweep":
		gct.markSweepTime = int64(value)
	}
}

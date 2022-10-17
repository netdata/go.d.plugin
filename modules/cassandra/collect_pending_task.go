// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	collectorPendingTask  = "PendingTasks"
	metricPendingTaskType = "org_apache_cassandra_metrics_compaction_value"
)

func doCollectPendingTask(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, metricPendingTaskType, collectorPendingTask, false)
	return enabled && success
}

func collectPendingTask(pms prometheus.Metrics) *pendingTask {
	if !doCollectPendingTask(pms) {
		return nil
	}

	var pt pendingTask
	collectPendingTaskByType(&pt, pms)

	return &pt
}

func collectPendingTaskByType(pt *pendingTask, pms prometheus.Metrics) {
	for _, pm := range pms.FindByName(metricPendingTaskType) {
		metricName := pm.Labels.Get("name")
		assignPendingTaskMetric(pt, metricName, pm.Value)
	}
}

func assignPendingTaskMetric(pt *pendingTask, scope string, value float64) {
	switch scope {
	default:
	case "PendingTasks":
		pt.task = int64(value)
	}
}

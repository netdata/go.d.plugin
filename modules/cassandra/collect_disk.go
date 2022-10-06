// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

import (
	"fmt"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	collectorDisk  = "LiveDiskSpaceUsed"
	metricDiskType = "org_apache_cassandra_metrics_table_count"
)

func doCollectDisk(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, metricDiskType, collectorDisk, false)
	return enabled && success
}

func collectDisk(pms prometheus.Metrics) *DISK {
	fmt.Println("Testing = ", doCollectDisk(pms))
	if !doCollectDisk(pms) {
		return nil
	}

	var di DISK
	collectDiskByType(&di, pms)

	return &di
}

func collectDiskByType(di *DISK, pms prometheus.Metrics) {
	var total DISK
	for _, pm := range pms.FindByName(metricDiskType) {
		metricName := pm.Labels.Get("name")
		sumDiskMetric(&total, metricName, pm.Value)
	}

	di.load = total.load
	di.used = total.used
	di.compaction_completed = total.compaction_completed
	di.compaction_queue = total.compaction_queue
}

func sumDiskMetric(di *DISK, scope string, value float64) {
	switch scope {
	default:
	case "LiveDiskSpaceUsed":
		di.load += value
	case "TotalDiskSpaceUsed":
		di.used += value
	case "CompactionBytesWritten":
		di.compaction_completed += value
	case "PendingCompactions":
		di.compaction_queue += value
	}
}

// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	collectorProcess = "process"

	metricAppsCPUTimeTotal  = "windows_process_cpu_time_total"
	metricAppsCPUHandles    = "windows_process_handles"
	metricAppsIOBytes       = "windows_process_io_bytes_total"
	metricAppsIOOperations  = "windows_process_io_operations_total"
	metricAppsPageFaults    = "windows_process_page_faults_total"
	metricAppsPageFileBytes = "windows_process_page_file_bytes"
	metricAppsPoolBytes     = "windows_process_pool_bytes"
	metricAppsThreads       = "windows_process_threads"
)

func doCollectAPPS(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, collectorProcess)
	return enabled && success
}

func collectAPPS(pms prometheus.Metrics) *appsMetrics {
	if !doCollectCPU(pms) {
		return nil
	}

	apps := &appsMetrics{}

	return apps
}

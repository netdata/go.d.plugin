// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

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

func doCollectProcess(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, collectorProcess)
	return enabled && success
}

func collectProcess(pms prometheus.Metrics) *processesMetrics {
	if !doCollectProcess(pms) {
		return nil
	}

	procs := &processesMetrics{procs: make(map[string]*processMetrics)}
	collectProcessCpuTimeTotal(procs, pms)
	collectProcessHandles(procs, pms)
	collectProcessIOBytes(procs, pms)
	collectProcessIOOperations(procs, pms)
	collectProcessPageFaults(procs, pms)
	collectProcessPageFileBytes(procs, pms)
	collectProcessPoolBytes(procs, pms)
	collectProcessThreads(procs, pms)

	return procs
}

func collectProcessCpuTimeTotal(procs *processesMetrics, pms prometheus.Metrics) {
	var proc *processMetrics
	for _, pm := range pms.FindByName(metricAppsCPUTimeTotal) {
		name := pm.Labels.Get("process")
		if name == "" {
			continue
		}

		if proc == nil || proc.ID != name {
			proc = procs.get(name)
		}

		proc.CPUTime += pm.Value
	}
}

func collectProcessHandles(procs *processesMetrics, pms prometheus.Metrics) {
	var proc *processMetrics
	for _, pm := range pms.FindByName(metricAppsCPUHandles) {
		name := pm.Labels.Get("process")
		if name == "" {
			continue
		}

		if proc == nil || proc.ID != name {
			proc = procs.get(name)
		}

		proc.Handles += pm.Value
	}
}

func collectProcessIOBytes(procs *processesMetrics, pms prometheus.Metrics) {
	var proc *processMetrics
	for _, pm := range pms.FindByName(metricAppsIOBytes) {
		name := pm.Labels.Get("process")
		if name == "" {
			continue
		}

		if proc == nil || proc.ID != name {
			proc = procs.get(name)
		}

		proc.IOBytes += pm.Value
	}
}

func collectProcessIOOperations(procs *processesMetrics, pms prometheus.Metrics) {
	var proc *processMetrics
	for _, pm := range pms.FindByName(metricAppsIOOperations) {
		name := pm.Labels.Get("process")
		if name == "" {
			continue
		}

		if proc == nil || proc.ID != name {
			proc = procs.get(name)
		}

		proc.IOOperations += pm.Value
	}
}

func collectProcessPageFaults(procs *processesMetrics, pms prometheus.Metrics) {
	var proc *processMetrics
	for _, pm := range pms.FindByName(metricAppsPageFaults) {
		name := pm.Labels.Get("process")
		if name == "" {
			continue
		}

		if proc == nil || proc.ID != name {
			proc = procs.get(name)
		}

		proc.PageFaults += pm.Value
	}
}

func collectProcessPageFileBytes(procs *processesMetrics, pms prometheus.Metrics) {
	var proc *processMetrics
	for _, pm := range pms.FindByName(metricAppsPageFileBytes) {
		name := pm.Labels.Get("process")
		if name == "" {
			continue
		}

		if proc == nil || proc.ID != name {
			proc = procs.get(name)
		}

		proc.PageFileBytes += pm.Value
	}
}

func collectProcessPoolBytes(procs *processesMetrics, pms prometheus.Metrics) {
	var proc *processMetrics
	for _, pm := range pms.FindByName(metricAppsPoolBytes) {
		name := pm.Labels.Get("process")
		if name == "" {
			continue
		}

		if proc == nil || proc.ID != name {
			proc = procs.get(name)
		}

		proc.PoolBytes += pm.Value
	}
}

func collectProcessThreads(procs *processesMetrics, pms prometheus.Metrics) {
	var proc *processMetrics
	for _, pm := range pms.FindByName(metricAppsThreads) {
		name := pm.Labels.Get("process")
		if name == "" {
			continue
		}

		if proc == nil || proc.ID != name {
			proc = procs.get(name)
		}

		proc.Threads += pm.Value
	}
}

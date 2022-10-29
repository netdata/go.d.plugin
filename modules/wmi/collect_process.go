// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	collectorProcess = "process"

	metricProcessCPUTimeTotal  = "windows_process_cpu_time_total"
	metricProcessCPUHandles    = "windows_process_handles"
	metricProcessIOBytes       = "windows_process_io_bytes_total"
	metricProcessIOOperations  = "windows_process_io_operations_total"
	metricProcessPageFaults    = "windows_process_page_faults_total"
	metricProcessPageFileBytes = "windows_process_page_file_bytes"
	metricProcessPoolBytes     = "windows_process_pool_bytes"
	metricProcessThreads       = "windows_process_threads"
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
	for _, pm := range pms.FindByName(metricProcessCPUTimeTotal) {
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
	for _, pm := range pms.FindByName(metricProcessCPUHandles) {
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
	for _, pm := range pms.FindByName(metricProcessIOBytes) {
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
	for _, pm := range pms.FindByName(metricProcessIOOperations) {
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
	for _, pm := range pms.FindByName(metricProcessPageFaults) {
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
	for _, pm := range pms.FindByName(metricProcessPageFileBytes) {
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
	for _, pm := range pms.FindByName(metricProcessPoolBytes) {
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
	for _, pm := range pms.FindByName(metricProcessThreads) {
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

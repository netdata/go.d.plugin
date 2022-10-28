// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import (
	"sort"

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

func doCollectApps(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, collectorProcess)
	return enabled && success
}

func collectApps(pms prometheus.Metrics) *appsMetrics {
	if !doCollectApps(pms) {
		return nil
	}

	apps := &appsMetrics{}
	collectAppsCpuTimeTotal(apps, pms)
	collectAppsHandles(apps, pms)
	collectAppsIOBytes(apps, pms)
	collectAppsIOOperations(apps, pms)
	collectAppsPageFaults(apps, pms)
	collectAppsPageFileBytes(apps, pms)
	collectAppsPoolBytes(apps, pms)
	collectAppsThreads(apps, pms)
	sortApps(&apps.info)

	return apps
}

func collectAppsCpuTimeTotal(apps *appsMetrics, pms prometheus.Metrics) {
	var app *appsInfo
	for _, pm := range pms.FindByName(metricAppsCPUTimeTotal) {
		processName := pm.Labels.Get("process")
		if processName == "" {
			continue
		}

		if app == nil || app.ID != processName {
			app = apps.info.get(processName)
		}

		app.cpuTimeTotal = pm.Value
	}
}

func collectAppsHandles(apps *appsMetrics, pms prometheus.Metrics) {
	var app *appsInfo
	for _, pm := range pms.FindByName(metricAppsCPUHandles) {
		processName := pm.Labels.Get("process")
		if processName == "" {
			continue
		}

		if app == nil || app.ID != processName {
			app = apps.info.get(processName)
		}

		app.cpuHandles += pm.Value
	}
}

func collectAppsIOBytes(apps *appsMetrics, pms prometheus.Metrics) {
	var app *appsInfo
	for _, pm := range pms.FindByName(metricAppsIOBytes) {
		processName := pm.Labels.Get("process")
		if processName == "" {
			continue
		}

		if app == nil || app.ID != processName {
			app = apps.info.get(processName)
		}

		app.ioBytes += pm.Value
	}
}

func collectAppsIOOperations(apps *appsMetrics, pms prometheus.Metrics) {
	var app *appsInfo
	for _, pm := range pms.FindByName(metricAppsIOOperations) {
		processName := pm.Labels.Get("process")
		if processName == "" {
			continue
		}

		if app == nil || app.ID != processName {
			app = apps.info.get(processName)
		}

		app.ioOperations += pm.Value
	}
}

func collectAppsPageFaults(apps *appsMetrics, pms prometheus.Metrics) {
	var app *appsInfo
	for _, pm := range pms.FindByName(metricAppsPageFaults) {
		processName := pm.Labels.Get("process")
		if processName == "" {
			continue
		}

		if app == nil || app.ID != processName {
			app = apps.info.get(processName)
		}

		app.PageFaults += pm.Value
	}
}

func collectAppsPageFileBytes(apps *appsMetrics, pms prometheus.Metrics) {
	var app *appsInfo
	for _, pm := range pms.FindByName(metricAppsPageFileBytes) {
		processName := pm.Labels.Get("process")
		if processName == "" {
			continue
		}

		if app == nil || app.ID != processName {
			app = apps.info.get(processName)
		}

		app.PageFileBytes += pm.Value
	}
}

func collectAppsPoolBytes(apps *appsMetrics, pms prometheus.Metrics) {
	var app *appsInfo
	for _, pm := range pms.FindByName(metricAppsPoolBytes) {
		processName := pm.Labels.Get("process")
		if processName == "" {
			continue
		}

		if app == nil || app.ID != processName {
			app = apps.info.get(processName)
		}

		app.PoolBytes += pm.Value
	}
}

func collectAppsThreads(apps *appsMetrics, pms prometheus.Metrics) {
	var app *appsInfo
	for _, pm := range pms.FindByName(metricAppsPoolBytes) {
		processName := pm.Labels.Get("process")
		if processName == "" {
			continue
		}

		if app == nil || app.ID != processName {
			app = apps.info.get(processName)
		}

		app.Threads += pm.Value
	}
}

func sortApps(apps *appsInfos) {
	sort.Slice(*apps, func(i, j int) bool { return (*apps)[i].ID < (*apps)[j].ID })
}

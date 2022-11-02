// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import (
	"time"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	metricSysContextSwitchesTotal     = "windows_system_context_switches_total"
	metricSysExceptionDispatchesTotal = "windows_system_exception_dispatches_total"
	metricSysProcessorQueueLength     = "windows_system_processor_queue_length"
	metricSysSystemCallsTotal         = "windows_system_system_calls_total"
	metricSysSystemUpTime             = "windows_system_system_up_time"
	metricSysThreads                  = "windows_system_threads"
)

func (w *WMI) collectSystem(mx map[string]int64, pms prometheus.Metrics) {
	if !w.cache.collection[collectorSystem] {
		w.cache.collection[collectorSystem] = true
		w.addSystemCharts()
	}

	px := "system_"
	if pm := pms.FindByName(metricSysContextSwitchesTotal); pm.Len() > 0 {
		mx[px+"context_switches_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricSysExceptionDispatchesTotal); pm.Len() > 0 {
		mx[px+"exception_dispatches_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricSysProcessorQueueLength); pm.Len() > 0 {
		mx[px+"processor_queue_length"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricSysSystemCallsTotal); pm.Len() > 0 {
		mx[px+"calls_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricSysSystemUpTime); pm.Len() > 0 {
		mx[px+"up_time"] = time.Now().Unix() - int64(pm.Max())
	}
	if pm := pms.FindByName(metricSysThreads); pm.Len() > 0 {
		mx[px+"threads"] = int64(pm.Max())
	}
}

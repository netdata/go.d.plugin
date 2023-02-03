// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	metricNetFrameworkCLRLocksAndThreadsQueueLengthTotal      = "windows_netframework_clrlocksandthreads_queue_length_total"
	metricNetFrameworkCLRLocksAndThreadsCurrentLogicalThreads = "windows_netframework_clrlocksandthreads_current_logical_threads"
	metricNetFrameworkCLRLocksAndThreadsPhysicalThreads       = "windows_netframework_clrlocksandthreads_physical_threads_current"
	metricNetFrameworkCLRLocksAndThreadsRecognizedThreads     = "windows_netframework_clrlocksandthreads_recognized_threads_total"
	metricNetFrameworkCLRLocksAndThreadsContentions           = "windows_netframework_clrlocksandthreads_contentions_total"
)

func (w *WMI) collectNetFrameworkCLRLocksandthreads(mx map[string]int64, pms prometheus.Series) {
	seen := make(map[string]bool)
	for _, pm := range pms.FindByName(metricNetFrameworkCLRLocksAndThreadsQueueLengthTotal) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[netframeworkPrefix+name+"_clrlocksandthreads_queue_length_total"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRLocksAndThreadsCurrentLogicalThreads) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[netframeworkPrefix+name+"_clrlocksandthreads_current_logical_threads_total"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRLocksAndThreadsPhysicalThreads) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[netframeworkPrefix+name+"_clrlocksandthreads_current_physical_threads_total"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRLocksAndThreadsRecognizedThreads) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[netframeworkPrefix+name+"_clrlocksandthreads_recognized_threads_total"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRLocksAndThreadsContentions) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[netframeworkPrefix+name+"_clrlocksandthreads_contentions_total"] += int64(pm.Value)
		}
	}

	for proc := range seen {
		if !w.cache.netFrameworkCLRLocksandthreads[proc] {
			w.cache.netFrameworkCLRLocksandthreads[proc] = true
			w.addProcessToNetFrameworkLocksandthreadsCharts(proc)
		}
	}

	for proc := range w.cache.netFrameworkCLRLocksandthreads {
		if !seen[proc] {
			delete(w.cache.netFrameworkCLRLocksandthreads, proc)
			w.removeProcessFromNetFrameworkLocksandthreadsCharts(proc)
		}
	}
}

// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	metricNetFrameworkCLRJITMethods         = "windows_netframework_clrjit_jit_methods_total"
	metricNetFrameworkCLRJITTime            = "windows_netframework_clrjit_jit_time_percent"
	metricNetFrameworkCLRJITStandardFailure = "windows_netframework_clrjit_jit_standard_failures_total"
	metricNetFrameworkCLRJITILBytes         = "windows_netframework_clrjit_jit_il_bytes_total"
)

func (w *WMI) collectNetFrameworkCLRJIT(mx map[string]int64, pms prometheus.Series) {
	if !w.cache.collection[collectorNetFrameworkCLRJIT] {
		w.cache.collection[collectorNetFrameworkCLRJIT] = true
		w.addNetFrameworkCRLJIT()
	}

	seen := make(map[string]bool)
	px := "netframework_clrjit_"

	for _, pm := range pms.FindByName(metricNetFrameworkCLRJITMethods) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_methods"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRJITStandardFailure) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_standard_failure"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRJITTime) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_time"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRJITILBytes) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_il_bytes"] += int64(pm.Value)
		}
	}

	for proc := range seen {
		if !w.cache.netFrameworkCLRJIT[proc] {
			w.cache.netFrameworkCLRJIT[proc] = true
			w.addProcessToNetFrameworkJITCharts(proc)
		}
	}

	for proc := range w.cache.netFrameworkCLRJIT {
		if !seen[proc] {
			delete(w.cache.netFrameworkCLRJIT, proc)
			w.removeProcessFromNetFrameworkJITCharts(proc)
		}
	}
}

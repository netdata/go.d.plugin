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
	seen := make(map[string]bool)

	for _, pm := range pms.FindByName(metricNetFrameworkCLRJITMethods) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[netframeworkPrefix+name+"_clrjit_methods_total"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRJITStandardFailure) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[netframeworkPrefix+name+"_clrjit_standard_failure_total"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRJITTime) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[netframeworkPrefix+name+"_clrjit_time_percent"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRJITILBytes) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[netframeworkPrefix+name+"_clrjit_il_bytes_total"] += int64(pm.Value)
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

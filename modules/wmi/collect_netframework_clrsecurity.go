// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	metricNetFrameworkCLRSecurityLinkTimeChecksTotal = "windows_netframework_clrsecurity_link_time_checks_total"
	metricNetFrameworkCLRSecurityRTChecksTimePercent = "windows_netframework_clrsecurity_rt_checks_time_percent"
	metricNetFrameworkCLRSecurityStackWalkDepth      = "windows_netframework_clrsecurity_stack_walk_depth"
	metricNetFrameworkCLRSecurityRuntimeChecksTotal  = "windows_netframework_clrsecurity_runtime_checks_total"
)

func (w *WMI) collectNetFrameworkCLR(mx map[string]int64, pms prometheus.Series) {
	if !w.cache.collection[collectorNetFrameworkCLRExceptions] {
		w.cache.collection[collectorNetFrameworkCLRExceptions] = true
		w.addNetFrameworkCRLExceptions()
	}

	seen := make(map[string]bool)
	px := "netframework_clr"
	for _, pm := range pms.FindByName(metricNetFrameworkCLRSecurityLinkTimeChecksTotal) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"security_"+name+"_link_time_checks"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRSecurityRTChecksTimePercent) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"security_"+name+"_checks_time"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRSecurityStackWalkDepth) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"security_"+name+"_stack_walk_depth"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRSecurityRuntimeChecksTotal) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"security_"+name+"_runtime_checks"] += int64(pm.Value)
		}
	}

	for proc := range seen {
		if !w.cache.netFrameworkCLRExceptions[proc] {
			w.cache.netFrameworkCLRExceptions[proc] = true
			w.addProcessToNetFrameworkCharts(proc)
		}
	}

	for proc := range w.cache.netFrameworkCLRExceptions {
		if !seen[proc] {
			delete(w.cache.netFrameworkCLRExceptions, proc)
			w.removeProcessFromNetFrameworkCharts(proc)
		}
	}
}

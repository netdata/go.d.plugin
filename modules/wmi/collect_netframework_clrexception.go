// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	netframeworkPrefix = "netframework_"
)

const (
	metricNetFrameworkCLRExceptionsThrown          = "windows_netframework_clrexceptions_exceptions_thrown_total"
	metricNetFrameworkCLRExceptionsFilters         = "windows_netframework_clrexceptions_exceptions_filters_total"
	metricNetFrameworkCLRExceptionsFinallys        = "windows_netframework_clrexceptions_exceptions_finallys_total"
	metricNetFrameworkCLRExceptionsThrowCatchDepth = "windows_netframework_clrexceptions_throw_to_catch_depth_total"
)

func (w *WMI) collectNetFrameworkCLRExceptions(mx map[string]int64, pms prometheus.Series) {
	seen := make(map[string]bool)
	for _, pm := range pms.FindByName(metricNetFrameworkCLRExceptionsThrown) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[netframeworkPrefix+name+"_clrexception_thrown_total"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRExceptionsFilters) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[netframeworkPrefix+name+"_clrexception_filters_total"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRExceptionsFinallys) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[netframeworkPrefix+name+"_clrexception_finally_total"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRExceptionsThrowCatchDepth) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[netframeworkPrefix+name+"_clrexception_throw_catch_depth_total"] += int64(pm.Value)
		}
	}

	for proc := range seen {
		if !w.cache.netFrameworkCLRExceptions[proc] {
			w.cache.netFrameworkCLRExceptions[proc] = true
			w.addProcessToNetFrameworkExceptionsCharts(proc)
		}
	}

	for proc := range w.cache.netFrameworkCLRExceptions {
		if !seen[proc] {
			delete(w.cache.netFrameworkCLRExceptions, proc)
			w.removeProcessFromNetFrameworkExceptionsCharts(proc)
		}
	}
}

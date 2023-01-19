// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	metricNetFrameworkCLRExceptionsThrown          = "windows_netframework_clrexceptions_exceptions_thrown_total"
	metricNetFrameworkCLRExceptionsFilters         = "windows_netframework_clrexceptions_exceptions_filters_total"
	metricNetFrameworkCLRExceptionsFinallys        = "windows_netframework_clrexceptions_exceptions_finallys_total"
	metricNetFrameworkCLRExceptionsThrowCatchDepth = "windows_netframework_clrexceptions_throw_to_catch_depth_total"
)

func (w *WMI) collectNetFrameworkCLR(mx map[string]int64, pms prometheus.Series) {
	if !w.cache.collection[collectorProcess] {
		w.cache.collection[collectorNetFrameworkCLR] = true
		w.addNetFrameworkCRLExceptions()
	}

	seen := make(map[string]bool)
	px := "netframework_clrexceptions_"
	for _, pm := range pms.FindByName(metricNetFrameworkCLRExceptionsThrown) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_thrown"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRExceptionsFilters) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_filters"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRExceptionsFinallys) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_finallys"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRExceptionsThrowCatchDepth) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_throw_catch_depth"] += int64(pm.Value)
		}
	}

	for proc := range seen {
		if !w.cache.netFrameworkCLR[proc] {
			w.cache.netFrameworkCLR[proc] = true
			w.addProcessToNetFrameworkCharts(proc)
		}
	}
	for proc := range w.cache.netFrameworkCLR {
		if !seen[proc] {
			delete(w.cache.netFrameworkCLR, proc)
			w.removeProcessFromNetFrameworkCharts(proc)
		}
	}
}

// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	metricNetFrameworkCLRInteropComCallableWrapper = "windows_netframework_clrinterop_com_callable_wrappers_total"
	metricNetFrameworkCLRInteropMarshalling        = "windows_netframework_clrinterop_interop_marshalling_total"
	metricNetFrameworkCLRInteropStubsCreated       = "windows_netframework_clrinterop_interop_stubs_created_total"
)

func (w *WMI) collectNetFrameworkCLRInterop(mx map[string]int64, pms prometheus.Series) {
	if !w.cache.collection[collectorNetFrameworkCLRInterop] {
		w.cache.collection[collectorNetFrameworkCLRInterop] = true
		w.addNetFrameworkCRLInterop()
	}

	seen := make(map[string]bool)
	px := "netframework_clrinterop_"
	for _, pm := range pms.FindByName(metricNetFrameworkCLRInteropComCallableWrapper) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_com_callable_wrappers"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRInteropMarshalling) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_marshalling"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRInteropStubsCreated) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_stubs_created"] += int64(pm.Value)
		}
	}

	for proc := range seen {
		if !w.cache.netFrameworkCLRInterops[proc] {
			w.cache.netFrameworkCLRInterops[proc] = true
			w.addProcessToNetFrameworkInteropCharts(proc)
		}
	}

	for proc := range w.cache.netFrameworkCLRInterops {
		if !seen[proc] {
			delete(w.cache.netFrameworkCLRInterops, proc)
			w.removeProcessFromNetFrameworkInteropCharts(proc)
		}
	}
}

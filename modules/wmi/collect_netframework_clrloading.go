// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	metricNetFrameworkCLRLoadingLoaderHeapSize    = "windows_netframework_clrloading_loader_heap_size_bytes"
	metricNetFrameworkCLRLoadingAppDomainLoaded   = "windows_netframework_clrloading_appdomains_loaded_total"
	metricNetFrameworkCLRLoadingAppDomainUnloaded = "windows_netframework_clrloading_appdomains_unloaded_total"
	metricNetFrameworkCLRLoadingAssembliesLoaded  = "windows_netframework_clrloading_assemblies_loaded_total"
	metricNetFrameworkCLRLoadingClassesLoaded     = "windows_netframework_clrloading_classes_loaded_total"
	metricNetFrameworkCLRLoadingClassLoadFailure  = "windows_netframework_clrloading_class_load_failures_total"
)

func (w *WMI) collectNetFrameworkCLRLoading(mx map[string]int64, pms prometheus.Series) {
	seen := make(map[string]bool)
	px := "net_framework_"
	for _, pm := range pms.FindByName(metricNetFrameworkCLRLoadingLoaderHeapSize) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_clrloading_loader_heap_size_bytes"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRLoadingAppDomainLoaded) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_clrloading_appdomains_loaded_total"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRLoadingAppDomainUnloaded) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_clrloading_appdomains_unloaded_total"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRLoadingAssembliesLoaded) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_clrloading_assemblies_loaded_total"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRLoadingClassesLoaded) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_clrloading_classes_loaded_total"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRLoadingClassLoadFailure) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_clrloading_class_load_failures_total"] += int64(pm.Value)
		}
	}

	for proc := range seen {
		if !w.cache.netFrameworkCLRLoading[proc] {
			w.cache.netFrameworkCLRLoading[proc] = true
			w.addProcessToNetFrameworkLoadingCharts(proc)
		}
	}

	for proc := range w.cache.netFrameworkCLRLoading {
		if !seen[proc] {
			delete(w.cache.netFrameworkCLRLoading, proc)
			w.removeProcessFromNetFrameworkLoadingCharts(proc)
		}
	}
}

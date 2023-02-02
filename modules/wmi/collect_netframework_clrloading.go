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
	if !w.cache.collection[collectorNetFrameworkCLRLoading] {
		w.cache.collection[collectorNetFrameworkCLRLoading] = true
		w.addNetFrameworkCRLLoading()
	}

	seen := make(map[string]bool)
	px := "netframework_clrloading_"
	for _, pm := range pms.FindByName(metricNetFrameworkCLRLoadingLoaderHeapSize) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_loader_heap_size"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRLoadingAppDomainLoaded) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_app_domains_loaded"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRLoadingAppDomainUnloaded) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_app_domains_unloaded"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRLoadingAssembliesLoaded) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_assemblies_loaded"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRLoadingClassesLoaded) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_classes_loaded"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRLoadingClassLoadFailure) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_class_load_failure"] += int64(pm.Value)
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

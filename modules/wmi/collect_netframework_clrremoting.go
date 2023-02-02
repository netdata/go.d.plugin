// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	metricNetFrameworkCLRRemotingChannelsTotal             = "windows_netframework_clrremoting_channels_total"
	metricNetFrameworkCLRRemotingContextBoundClassesLoaded = "windows_netframework_clrremoting_context_bound_classes_loaded"
	metricNetFrameworkCLRRemotingContextBoundObjectsTotal  = "windows_netframework_clrremoting_context_bound_objects_total"
	metricNetFrameworkCLRRemotingContextProxiesTotal       = "windows_netframework_clrremoting_context_proxies_total"
	metricNetFrameworkCLRRemotingContexts                  = "windows_netframework_clrremoting_contexts"
	metricNetFrameworkCLRRemotingRemoteCallsTotal          = "windows_netframework_clrremoting_remote_calls_total"
)

func (w *WMI) collectNetFrameworkCLRRemoting(mx map[string]int64, pms prometheus.Series) {
	if !w.cache.collection[collectorNetFrameworkCLRExceptions] {
		w.cache.collection[collectorNetFrameworkCLRExceptions] = true
		w.addNetFrameworkCRLExceptions()
	}

	seen := make(map[string]bool)
	px := "netframework_clrremoting_"
	for _, pm := range pms.FindByName(metricNetFrameworkCLRRemotingChannelsTotal) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_channels"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRRemotingContextBoundClassesLoaded) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_context_bound_classes_loaded"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRRemotingContextBoundObjectsTotal) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_context_bound_objects"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRRemotingContextProxiesTotal) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_context_proxies"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRRemotingContexts) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_contexts"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRRemotingRemoteCallsTotal) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_calls"] += int64(pm.Value)
		}
	}

	for proc := range seen {
		if !w.cache.netFrameworkCLRRemoting[proc] {
			w.cache.netFrameworkCLRRemoting[proc] = true
			w.addProcessToNetFrameworkRemotingCharts(proc)
		}
	}

	for proc := range w.cache.netFrameworkCLRRemoting {
		if !seen[proc] {
			delete(w.cache.netFrameworkCLRRemoting, proc)
			w.removeProcessFromNetFrameworkRemotingCharts(proc)
		}
	}
}

// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	metricNetFrameworkCLRMemoryAllocatedBytesTotal   = "windows_netframework_clrmemory_allocated_bytes_total"
	metricNetFrameworkCLRMemoryFinalizationSurvivors = "windows_netframework_clrmemory_finalization_survivors"
	metricNetFrameworkCLRMemoryHeapSizeBytes         = "windows_netframework_clrmemory_heap_size_bytes"
	metricNetFrameworkCLRMemoryPromotedBytes         = "windows_netframework_clrmemory_promoted_bytes"
	metricNetFrameworkCLRMemoryNumberGCHandles       = "windows_netframework_clrmemory_number_gc_handles"
	metricNetFrameworkCLRMemoryCollectionsTotal      = "windows_netframework_clrmemory_collections_total"
	metricNetFrameworkCLRMemoryInducedGCTotal        = "windows_netframework_clrmemory_induced_gc_total"
	metricNetFrameworkCLRMemoryNumberPinnedObjects   = "windows_netframework_clrmemory_number_pinned_objects"
	metricNetFrameworkCLRMemoryNumberSinkBlockInUse  = "windows_netframework_clrmemory_number_sink_blocksinuse"
	metricNetFrameworkCLRMemoryCommittedBytes        = "windows_netframework_clrmemory_committed_bytes"
	metricNetFrameworkCLRMemoryReservedBytes         = "windows_netframework_clrmemory_reserved_bytes"
	metricNetFrameworkCLRMemoryGCTimePecent          = "windows_netframework_clrmemory_gc_time_percent"

	metricNetFrameworkCLRRemotingChannelsTotal             = "windows_netframework_clrremoting_channels_total"
	metricNetFrameworkCLRRemotingContextBoundClassesLoaded = "windows_netframework_clrremoting_context_bound_classes_loaded"
	metricNetFrameworkCLRRemotingContextBoundObjectsTotal  = "windows_netframework_clrremoting_context_bound_objects_total"
	metricNetFrameworkCLRRemotingContextProxiesTotal       = "windows_netframework_clrremoting_context_proxies_total"
	metricNetFrameworkCLRRemotingContexts                  = "windows_netframework_clrremoting_contexts"
	metricNetFrameworkCLRRemotingRemoteCallsTotal          = "windows_netframework_clrremoting_remote_calls_total"

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
	for _, pm := range pms.FindByName(metricNetFrameworkCLRMemoryAllocatedBytesTotal) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"memory_"+name+"_allocated_bytes_total"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRMemoryFinalizationSurvivors) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"memory_"+name+"_finalization_survivors"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRMemoryHeapSizeBytes) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"memory_"+name+"_heap_size"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRMemoryPromotedBytes) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"memory_"+name+"_promoted"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRMemoryNumberGCHandles) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"memory_"+name+"_gc_handles"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRMemoryCollectionsTotal) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"memory_"+name+"_collection"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRMemoryInducedGCTotal) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"memory_"+name+"_induced_gc_total"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRMemoryNumberPinnedObjects) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"memory_"+name+"_pinned_objects"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRMemoryNumberSinkBlockInUse) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"memory_"+name+"_sink_block_in_use"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRMemoryCommittedBytes) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"memory_"+name+"_committed"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRMemoryReservedBytes) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"memory_"+name+"_reserved"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRMemoryGCTimePecent) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"memory_"+name+"_gc_time"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRRemotingChannelsTotal) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"remoting_"+name+"_channels"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRRemotingContextBoundClassesLoaded) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"remoting_"+name+"_context_bound_classes_loaded"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRRemotingContextBoundObjectsTotal) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"remoting_"+name+"_context_bound_objects"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRRemotingContextProxiesTotal) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"remoting_"+name+"_context_proxies"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRRemotingContexts) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"remoting_"+name+"_contexts"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRRemotingRemoteCallsTotal) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"remoting_"+name+"_calls"] += int64(pm.Value)
		}
	}

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

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
)

func (w *WMI) collectNetFrameworkCLRMemory(mx map[string]int64, pms prometheus.Series) {
	if !w.cache.collection[collectorNetFrameworkCLRMemory] {
		w.cache.collection[collectorNetFrameworkCLRMemory] = true
		w.addNetFrameworkCRLMemory()
	}

	seen := make(map[string]bool)
	px := "netframework_clrmemory_"
	for _, pm := range pms.FindByName(metricNetFrameworkCLRMemoryAllocatedBytesTotal) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_allocated_bytes_total"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRMemoryFinalizationSurvivors) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_finalization_survivors"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRMemoryHeapSizeBytes) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_heap_size"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRMemoryPromotedBytes) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_promoted"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRMemoryNumberGCHandles) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_gc_handles"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRMemoryCollectionsTotal) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_collection"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRMemoryInducedGCTotal) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_induced_gc_total"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRMemoryNumberPinnedObjects) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_pinned_objects"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRMemoryNumberSinkBlockInUse) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_sink_block_in_use"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRMemoryCommittedBytes) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_committed"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRMemoryReservedBytes) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_reserved"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRMemoryGCTimePecent) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+name+"_gc_time"] += int64(pm.Value)
		}
	}

	for proc := range seen {
		if !w.cache.netFrameworkCLRMemory[proc] {
			w.cache.netFrameworkCLRMemory[proc] = true
			w.addProcessToNetFrameworkMemoryCharts(proc)
		}
	}

	for proc := range w.cache.netFrameworkCLRMemory {
		if !seen[proc] {
			delete(w.cache.netFrameworkCLRMemory, proc)
			w.removeProcessFromNetFrameworkMemoryCharts(proc)
		}
	}
}

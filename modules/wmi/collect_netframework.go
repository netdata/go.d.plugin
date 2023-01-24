// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	metricNetFrameworkCLRExceptionsThrown          = "windows_netframework_clrexceptions_exceptions_thrown_total"
	metricNetFrameworkCLRExceptionsFilters         = "windows_netframework_clrexceptions_exceptions_filters_total"
	metricNetFrameworkCLRExceptionsFinallys        = "windows_netframework_clrexceptions_exceptions_finallys_total"
	metricNetFrameworkCLRExceptionsThrowCatchDepth = "windows_netframework_clrexceptions_throw_to_catch_depth_total"

	metricNetFrameworkCLRInteropComCallableWrapper = "windows_netframework_clrinterop_com_callable_wrappers_total"
	metricNetFrameworkCLRInteropMarshalling        = "windows_netframework_clrinterop_interop_marshalling_total"
	metricNetFrameworkCLRInteropStubsCreated       = "windows_netframework_clrinterop_interop_stubs_created_total"

	metricNetFrameworkCLRJITMethods         = "windows_netframework_clrjit_jit_methods_total"
	metricNetFrameworkCLRJITTime            = "windows_netframework_clrjit_jit_time_percent"
	metricNetFrameworkCLRJITStandardFailure = "windows_netframework_clrjit_jit_standard_failures_total"
	metricNetFrameworkCLRJITILBytes         = "windows_netframework_clrjit_jit_il_bytes_total"

	metricNetFrameworkCLRLoadingLoaderHeapSize    = "windows_netframework_clrloading_loader_heap_size_bytes"
	metricNetFrameworkCLRLoadingAppDomainLoaded   = "windows_netframework_clrloading_appdomains_loaded_total"
	metricNetFrameworkCLRLoadingAppDomainUnloaded = "windows_netframework_clrloading_appdomains_unloaded_total"
	metricNetFrameworkCLRLoadingAssembliesLoaded  = "windows_netframework_clrloading_assemblies_loaded_total"
	metricNetFrameworkCLRLoadingClassesLoaded     = "windows_netframework_clrloading_classes_loaded_total"
	metricNetFrameworkCLRLoadingClassLoadFailure  = "windows_netframework_clrloading_class_load_failures_total"

	metricNetFrameworkCLRLocksAndThreadsQueueLengthTotal      = "windows_netframework_clrlocksandthreads_queue_length_total"
	metricNetFrameworkCLRLocksAndThreadsCurrentLogicalThreads = "windows_netframework_clrlocksandthreads_current_logical_threads"
	metricNetFrameworkCLRLocksAndThreadsPhysicalThreads       = "windows_netframework_clrlocksandthreads_physical_threads_current"
	metricNetFrameworkCLRLocksAndThreadsRecognizedThreads     = "windows_netframework_clrlocksandthreads_recognized_threads_total"
	metricNetFrameworkCLRLocksAndThreadsContentions           = "windows_netframework_clrlocksandthreads_contentions_total"

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

func (w *WMI) collectNetFrameworkCLR(mx map[string]int64, pms prometheus.Series) {
	if !w.cache.collection[collectorNetFrameworkCLR] {
		w.cache.collection[collectorNetFrameworkCLR] = true
		w.addNetFrameworkCRLExceptions()
	}

	seen := make(map[string]bool)
	px := "netframework_clr"
	for _, pm := range pms.FindByName(metricNetFrameworkCLRExceptionsThrown) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"exceptions_"+name+"_thrown"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRExceptionsFilters) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"exceptions_"+name+"_filters"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRExceptionsFinallys) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"exceptions_"+name+"_finallys"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRExceptionsThrowCatchDepth) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"exceptions_"+name+"_throw_catch_depth"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRInteropComCallableWrapper) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"interop_"+name+"_com_callable_wrappers"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRInteropMarshalling) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"interop_"+name+"_marshalling"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRInteropStubsCreated) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"interop_"+name+"_stubs_created"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRJITMethods) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"jit_"+name+"_methods"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRJITStandardFailure) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"jit_"+name+"_standard_failure"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRJITTime) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"jit_"+name+"_time"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRJITILBytes) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"jit_"+name+"_il_bytes"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRLoadingLoaderHeapSize) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"loading_"+name+"_loader_heap_size"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRLoadingAppDomainLoaded) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"loading_"+name+"_app_domains_loaded"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRLoadingAppDomainUnloaded) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"loading_"+name+"_app_domains_unloaded"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRLoadingAssembliesLoaded) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"loading_"+name+"_assemblies_loaded"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRLoadingClassesLoaded) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"loading_"+name+"_classes_loaded"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRLoadingClassLoadFailure) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"loading_"+name+"_class_load_failure"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRLocksAndThreadsQueueLengthTotal) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"lockandthreads_"+name+"_queue_length"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRLocksAndThreadsCurrentLogicalThreads) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"lockandthreads_"+name+"_current_logical_threads"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRLocksAndThreadsPhysicalThreads) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"lockandthreads_"+name+"_current_physical_threads"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRLocksAndThreadsRecognizedThreads) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"lockandthreads_"+name+"_recognized_threads"] += int64(pm.Value)
		}
	}

	for _, pm := range pms.FindByName(metricNetFrameworkCLRLocksAndThreadsContentions) {
		if name := cleanProcessName(pm.Labels.Get("process")); name != "" {
			seen[name] = true
			mx[px+"lockandthreads_"+name+"_contentions"] += int64(pm.Value)
		}
	}

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

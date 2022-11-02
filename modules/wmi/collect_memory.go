// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	metricMemAvailBytes                      = "windows_memory_available_bytes"
	metricMemCacheBytes                      = "windows_memory_cache_bytes"
	metricMemCacheBytesPeak                  = "windows_memory_cache_bytes_peak"
	metricMemCacheFaultsTotal                = "windows_memory_cache_faults_total"
	metricMemCommitLimit                     = "windows_memory_commit_limit"
	metricMemCommittedBytes                  = "windows_memory_committed_bytes"
	metricMemDemandZeroFaultsTotal           = "windows_memory_demand_zero_faults_total"
	metricMemFreeAndZeroPageListBytes        = "windows_memory_free_and_zero_page_list_bytes"
	metricMemFreeSystemPageTableEntries      = "windows_memory_free_system_page_table_entries"
	metricMemModifiedPageListBytes           = "windows_memory_modified_page_list_bytes"
	metricMemPageFaultsTotal                 = "windows_memory_page_faults_total"
	metricMemSwapPageReadsTotal              = "windows_memory_swap_page_reads_total"
	metricMemSwapPagesReadTotal              = "windows_memory_swap_pages_read_total"
	metricMemSwapPagesWrittenTotal           = "windows_memory_swap_pages_written_total"
	metricMemSwapPageOperationsTotal         = "windows_memory_swap_page_operations_total"
	metricMemSwapPageWritesTotal             = "windows_memory_swap_page_writes_total"
	metricMemPoolNonPagedAllocsTotal         = "windows_memory_pool_nonpaged_allocs_total"
	metricMemPoolNonPagedBytesTotal          = "windows_memory_pool_nonpaged_bytes"
	metricMemPoolPagedAllocsTotal            = "windows_memory_pool_paged_allocs_total"
	metricMemPoolPagedBytes                  = "windows_memory_pool_paged_bytes"
	metricMemPoolPagedResidentBytes          = "windows_memory_pool_paged_resident_bytes"
	metricMemStandbyCacheCoreBytes           = "windows_memory_standby_cache_core_bytes"
	metricMemStandbyCacheNormalPriorityBytes = "windows_memory_standby_cache_normal_priority_bytes"
	metricMemStandbyCacheReserveBytes        = "windows_memory_standby_cache_reserve_bytes"
	metricMemSystemCacheResidentBytes        = "windows_memory_system_cache_resident_bytes"
	metricMemSystemCodeResidentBytes         = "windows_memory_system_code_resident_bytes"
	metricMemSystemCodeTotalBytes            = "windows_memory_system_code_total_bytes"
	metricMemSystemDriverResidentBytes       = "windows_memory_system_driver_resident_bytes"
	metricMemSystemDriverTotalBytes          = "windows_memory_system_driver_total_bytes"
	metricMemTransitionFaultsTotal           = "windows_memory_transition_faults_total"
	metricMemTransitionPagesRepurposedTotal  = "windows_memory_transition_pages_repurposed_total"
	metricMemWriteCopiesTotal                = "windows_memory_write_copies_total"
)

func (w *WMI) collectMemory(mx map[string]int64, pms prometheus.Metrics) {
	if !w.cache.collection[collectorMemory] {
		w.cache.collection[collectorMemory] = true
		w.addMemoryCharts()
	}

	if pm := pms.FindByName(metricMemAvailBytes); pm.Len() > 0 {
		mx["memory_available_bytes"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemCacheBytes); pm.Len() > 0 {
		mx["memory_cache_bytes"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemCacheBytesPeak); pm.Len() > 0 {
		mx["memory_cache_bytes_peak"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemCacheFaultsTotal); pm.Len() > 0 {
		mx["memory_cache_faults_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemCommitLimit); pm.Len() > 0 {
		mx["memory_commit_limit"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemCommittedBytes); pm.Len() > 0 {
		mx["memory_committed_bytes"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemDemandZeroFaultsTotal); pm.Len() > 0 {
		mx["memory_demand_zero_faults_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemFreeAndZeroPageListBytes); pm.Len() > 0 {
		mx["memory_free_and_zero_page_list_bytes"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemFreeSystemPageTableEntries); pm.Len() > 0 {
		mx["memory_free_system_page_table_entries"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemModifiedPageListBytes); pm.Len() > 0 {
		mx["memory_modified_page_list_bytes"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemPageFaultsTotal); pm.Len() > 0 {
		mx["memory_page_faults_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemSwapPageReadsTotal); pm.Len() > 0 {
		mx["memory_swap_page_reads_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemSwapPagesReadTotal); pm.Len() > 0 {
		mx["memory_swap_pages_read_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemSwapPagesWrittenTotal); pm.Len() > 0 {
		mx["memory_swap_pages_written_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemSwapPageOperationsTotal); pm.Len() > 0 {
		mx["memory_swap_page_operations_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemSwapPageWritesTotal); pm.Len() > 0 {
		mx["memory_swap_page_writes_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemPoolNonPagedAllocsTotal); pm.Len() > 0 {
		mx["memory_pool_nonpaged_allocs_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemPoolNonPagedBytesTotal); pm.Len() > 0 {
		mx["memory_pool_nonpaged_bytes_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemPoolPagedAllocsTotal); pm.Len() > 0 {
		mx["memory_pool_paged_allocs_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemPoolPagedBytes); pm.Len() > 0 {
		mx["memory_pool_paged_bytes"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemPoolPagedResidentBytes); pm.Len() > 0 {
		mx["memory_pool_paged_resident_bytes"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemStandbyCacheCoreBytes); pm.Len() > 0 {
		mx["memory_standby_cache_core_bytes"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemStandbyCacheNormalPriorityBytes); pm.Len() > 0 {
		mx["memory_standby_cache_normal_priority_bytes"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemStandbyCacheReserveBytes); pm.Len() > 0 {
		mx["memory_standby_cache_reserve_bytes"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemSystemCacheResidentBytes); pm.Len() > 0 {
		mx["memory_system_cache_resident_bytes"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemSystemCodeResidentBytes); pm.Len() > 0 {
		mx["memory_system_code_resident_bytes"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemSystemCodeTotalBytes); pm.Len() > 0 {
		mx["memory_system_code_total_bytes"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemSystemDriverResidentBytes); pm.Len() > 0 {
		mx["memory_system_driver_resident_bytes"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemSystemDriverTotalBytes); pm.Len() > 0 {
		mx["memory_system_driver_total_bytes"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemTransitionFaultsTotal); pm.Len() > 0 {
		mx["memory_transition_faults_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemTransitionPagesRepurposedTotal); pm.Len() > 0 {
		mx["memory_transition_pages_repurposed_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricMemWriteCopiesTotal); pm.Len() > 0 {
		mx["memory_write_copies_total"] = int64(pm.Max())
	}
}

package wmi

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	metricMemAvailBytes                      = "wmi_memory_available_bytes"
	metricMemCacheBytes                      = "wmi_memory_cache_bytes"
	metricMemCacheBytesPeak                  = "wmi_memory_cache_bytes_peak"
	metricMemCacheFaultsTotal                = "wmi_memory_cache_faults_total"
	metricMemCommitLimit                     = "wmi_memory_commit_limit"
	metricMemCommittedBytes                  = "wmi_memory_committed_bytes"
	metricMemDemandZeroFaultsTotal           = "wmi_memory_demand_zero_faults_total"
	metricMemFreeAndZeroPageListBytes        = "wmi_memory_free_and_zero_page_list_bytes"
	metricMemFreeSystemPageTableEntries      = "wmi_memory_free_system_page_table_entries"
	metricMemModifiedPageListBytes           = "wmi_memory_modified_page_list_bytes"
	metricMemPageFaultsTotal                 = "wmi_memory_page_faults_total"
	metricMemSwapPageReadsTotal              = "wmi_memory_swap_page_reads_total"
	metricMemSwapPagesReadTotal              = "wmi_memory_swap_pages_read_total"
	metricMemSwapPagesWrittenTotal           = "wmi_memory_swap_pages_written_total"
	metricMemSwapPageOperationsTotal         = "wmi_memory_swap_page_operations_total"
	metricMemSwapPageWritesTotal             = "wmi_memory_swap_page_writes_total"
	metricMemPoolNonpagedAllocsTotal         = "wmi_memory_pool_nonpaged_allocs_total"
	metricMemPoolNonpagedBytesTotal          = "wmi_memory_pool_nonpaged_bytes_total"
	metricMemPoolPagedAllocsTotal            = "wmi_memory_pool_paged_allocs_total"
	metricMemPoolPagedBytes                  = "wmi_memory_pool_paged_bytes"
	metricMemPoolPagedResidentBytes          = "wmi_memory_pool_paged_resident_bytes"
	metricMemStandbyCacheCoreBytes           = "wmi_memory_standby_cache_core_bytes"
	metricMemStandbyCacheNormalPriorityBytes = "wmi_memory_standby_cache_normal_priority_bytes"
	metricMemStandbyCacheReserveBytes        = "wmi_memory_standby_cache_reserve_bytes"
	metricMemSystemCacheResidentBytes        = "wmi_memory_system_cache_resident_bytes"
	metricMemSystemCodeResidentBytes         = "wmi_memory_system_code_resident_bytes"
	metricMemSystemCodeTotalBytes            = "wmi_memory_system_code_total_bytes"
	metricMemSystemDriverResidentBytes       = "wmi_memory_system_driver_resident_bytes"
	metricMemSystemDriverTotalBytes          = "wmi_memory_system_driver_total_bytes"
	metricMemTransitionFaultsTotal           = "wmi_memory_transition_faults_total"
	metricMemTransitionPagesRepurposedTotal  = "wmi_memory_transition_pages_repurposed_total"
	metricMemWriteCopiesTotal                = "wmi_memory_write_copies_total"
)

func (w *WMI) collectMemory(mx *metrics, pms prometheus.Metrics) {
	mx.Memory.AvailableBytes = pms.FindByName(metricMemAvailBytes).Max()
	mx.Memory.CacheBytes = pms.FindByName(metricMemCacheBytes).Max()
	mx.Memory.CacheBytesPeak = pms.FindByName(metricMemCacheBytesPeak).Max()
	mx.Memory.CacheFaultsTotal = pms.FindByName(metricMemCacheFaultsTotal).Max()
	mx.Memory.CommitLimit = pms.FindByName(metricMemCommitLimit).Max()
	mx.Memory.CommittedBytes = pms.FindByName(metricMemCommittedBytes).Max()
	mx.Memory.DemandZeroFaultsTotal = pms.FindByName(metricMemDemandZeroFaultsTotal).Max()
	mx.Memory.FreeAndZeroPageListBytes = pms.FindByName(metricMemFreeAndZeroPageListBytes).Max()
	mx.Memory.FreeSystemPageTableEntries = pms.FindByName(metricMemFreeSystemPageTableEntries).Max()
	mx.Memory.ModifiedPageListBytes = pms.FindByName(metricMemModifiedPageListBytes).Max()
	mx.Memory.PageFaultsTotal = pms.FindByName(metricMemPageFaultsTotal).Max()
	mx.Memory.SwapPageReadsTotal = pms.FindByName(metricMemSwapPageReadsTotal).Max()
	mx.Memory.SwapPagesReadTotal = pms.FindByName(metricMemSwapPagesReadTotal).Max()
	mx.Memory.SwapPagesWrittenTotal = pms.FindByName(metricMemSwapPagesWrittenTotal).Max()
	mx.Memory.SwapPageOperationsTotal = pms.FindByName(metricMemSwapPageOperationsTotal).Max()
	mx.Memory.SwapPageWritesTotal = pms.FindByName(metricMemSwapPageWritesTotal).Max()
	mx.Memory.PoolNonPagedAllocsTotal = pms.FindByName(metricMemPoolNonpagedAllocsTotal).Max()
	mx.Memory.PoolNonPagedBytes = pms.FindByName(metricMemPoolNonpagedBytesTotal).Max()
	mx.Memory.PoolPagedAllocsTotal = pms.FindByName(metricMemPoolPagedAllocsTotal).Max()
	mx.Memory.PoolPagedBytes = pms.FindByName(metricMemPoolPagedBytes).Max()
	mx.Memory.PoolPagedResidentBytes = pms.FindByName(metricMemPoolPagedResidentBytes).Max()
	mx.Memory.StandbyCacheCoreBytes = pms.FindByName(metricMemStandbyCacheCoreBytes).Max()
	mx.Memory.StandbyCacheNormalPriorityBytes = pms.FindByName(metricMemStandbyCacheNormalPriorityBytes).Max()
	mx.Memory.StandbyCacheReserveBytes = pms.FindByName(metricMemStandbyCacheReserveBytes).Max()
	mx.Memory.SystemCacheResidentBytes = pms.FindByName(metricMemSystemCacheResidentBytes).Max()
	mx.Memory.SystemCodeResidentBytes = pms.FindByName(metricMemSystemCodeResidentBytes).Max()
	mx.Memory.SystemCodeTotalBytes = pms.FindByName(metricMemSystemCodeTotalBytes).Max()
	mx.Memory.SystemDriverResidentBytes = pms.FindByName(metricMemSystemDriverResidentBytes).Max()
	mx.Memory.SystemDriverTotalBytes = pms.FindByName(metricMemSystemDriverTotalBytes).Max()
	mx.Memory.TransitionFaultsTotal = pms.FindByName(metricMemTransitionFaultsTotal).Max()
	mx.Memory.TransitionPagesRePurposedTotal = pms.FindByName(metricMemTransitionPagesRepurposedTotal).Max()
	mx.Memory.WriteCopiesTotal = pms.FindByName(metricMemWriteCopiesTotal).Max()
}

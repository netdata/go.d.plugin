package wmi

import (
	"fmt"

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

func (w *WMI) collectMemory(mx *metrics, pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, collectorMemory)
	if !(enabled && success) {
		return false
	}
	mx.Memory = &memory{}

	names := []string{
		metricMemAvailBytes,
		metricMemCacheBytes,
		metricMemCacheBytesPeak,
		metricMemCacheFaultsTotal,
		metricMemCommitLimit,
		metricMemCommittedBytes,
		metricMemDemandZeroFaultsTotal,
		metricMemFreeAndZeroPageListBytes,
		metricMemFreeSystemPageTableEntries,
		metricMemModifiedPageListBytes,
		metricMemPageFaultsTotal,
		metricMemSwapPageReadsTotal,
		metricMemSwapPagesReadTotal,
		metricMemSwapPagesWrittenTotal,
		metricMemSwapPageOperationsTotal,
		metricMemSwapPageWritesTotal,
		metricMemPoolNonpagedAllocsTotal,
		metricMemPoolNonpagedBytesTotal,
		metricMemPoolPagedAllocsTotal,
		metricMemPoolPagedBytes,
		metricMemPoolPagedResidentBytes,
		metricMemStandbyCacheCoreBytes,
		metricMemStandbyCacheNormalPriorityBytes,
		metricMemStandbyCacheReserveBytes,
		metricMemSystemCacheResidentBytes,
		metricMemSystemCodeResidentBytes,
		metricMemSystemCodeTotalBytes,
		metricMemSystemDriverResidentBytes,
		metricMemSystemDriverTotalBytes,
		metricMemTransitionFaultsTotal,
		metricMemTransitionPagesRepurposedTotal,
		metricMemWriteCopiesTotal,
	}

	for _, name := range names {
		collectMemoryAny(mx, pms, name)
	}

	return true
}

func collectMemoryAny(mx *metrics, pms prometheus.Metrics, name string) {
	value := pms.FindByName(name).Max()

	switch name {
	default:
		panic(fmt.Sprintf("unknown metric name during memory collection : %s", name))
	case metricMemAvailBytes:
		mx.Memory.AvailableBytes = value
	case metricMemCacheBytes:
		mx.Memory.CacheBytes = value
	case metricMemCacheBytesPeak:
		mx.Memory.CacheBytesPeak = value
	case metricMemCacheFaultsTotal:
		mx.Memory.CacheFaultsTotal = value
	case metricMemCommitLimit:
		mx.Memory.CommitLimit = value
	case metricMemCommittedBytes:
		mx.Memory.CommittedBytes = value
	case metricMemDemandZeroFaultsTotal:
		mx.Memory.DemandZeroFaultsTotal = value
	case metricMemFreeAndZeroPageListBytes:
		mx.Memory.FreeAndZeroPageListBytes = value
	case metricMemFreeSystemPageTableEntries:
		mx.Memory.FreeSystemPageTableEntries = value
	case metricMemModifiedPageListBytes:
		mx.Memory.ModifiedPageListBytes = value
	case metricMemPageFaultsTotal:
		mx.Memory.PageFaultsTotal = value
	case metricMemSwapPageReadsTotal:
		mx.Memory.SwapPageReadsTotal = value
	case metricMemSwapPagesReadTotal:
		mx.Memory.SwapPagesReadTotal = value
	case metricMemSwapPagesWrittenTotal:
		mx.Memory.SwapPagesWrittenTotal = value
	case metricMemSwapPageOperationsTotal:
		mx.Memory.SwapPageOperationsTotal = value
	case metricMemSwapPageWritesTotal:
		mx.Memory.SwapPageWritesTotal = value
	case metricMemPoolNonpagedAllocsTotal:
		mx.Memory.PoolNonPagedAllocsTotal = value
	case metricMemPoolNonpagedBytesTotal:
		mx.Memory.PoolNonPagedBytes = value
	case metricMemPoolPagedAllocsTotal:
		mx.Memory.PoolPagedAllocsTotal = value
	case metricMemPoolPagedBytes:
		mx.Memory.PoolPagedBytes = value
	case metricMemPoolPagedResidentBytes:
		mx.Memory.PoolPagedResidentBytes = value
	case metricMemStandbyCacheCoreBytes:
		mx.Memory.StandbyCacheCoreBytes = value
	case metricMemStandbyCacheNormalPriorityBytes:
		mx.Memory.StandbyCacheNormalPriorityBytes = value
	case metricMemStandbyCacheReserveBytes:
		mx.Memory.StandbyCacheReserveBytes = value
	case metricMemSystemCacheResidentBytes:
		mx.Memory.SystemCacheResidentBytes = value
	case metricMemSystemCodeResidentBytes:
		mx.Memory.SystemCodeResidentBytes = value
	case metricMemSystemCodeTotalBytes:
		mx.Memory.SystemCodeTotalBytes = value
	case metricMemSystemDriverResidentBytes:
		mx.Memory.SystemDriverResidentBytes = value
	case metricMemSystemDriverTotalBytes:
		mx.Memory.SystemDriverTotalBytes = value
	case metricMemTransitionFaultsTotal:
		mx.Memory.TransitionFaultsTotal = value
	case metricMemTransitionPagesRepurposedTotal:
		mx.Memory.TransitionPagesRePurposedTotal = value
	case metricMemWriteCopiesTotal:
		mx.Memory.WriteCopiesTotal = value
	}
}

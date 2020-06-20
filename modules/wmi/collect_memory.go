package wmi

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	collectorMemory = "memory"

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
	metricMemPoolNonpagedAllocsTotal         = "windows_memory_pool_nonpaged_allocs_total"
	metricMemPoolNonpagedBytesTotal          = "windows_memory_pool_nonpaged_bytes_total"
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

var memoryMetricNames = []string{
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

func doCollectMemory(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, collectorMemory)
	return enabled && success
}

func collectMemory(pms prometheus.Metrics) *memoryMetrics {
	if !doCollectMemory(pms) {
		return nil
	}

	mm := &memoryMetrics{}
	for _, name := range memoryMetricNames {
		collectMemoryMetric(mm, pms, name)
	}

	mm.NotCommittedBytes = mm.CommitLimit - mm.CommittedBytes
	mm.StandbyCacheTotal = mm.StandbyCacheReserveBytes + mm.StandbyCacheNormalPriorityBytes + mm.StandbyCacheCoreBytes
	mm.Cached = mm.StandbyCacheTotal + mm.ModifiedPageListBytes
	return mm
}

func collectMemoryMetric(mm *memoryMetrics, pms prometheus.Metrics, name string) {
	value := pms.FindByName(name).Max()
	assignMemoryMetric(mm, name, value)
}

func assignMemoryMetric(mm *memoryMetrics, name string, value float64) {
	switch name {
	case metricMemAvailBytes:
		mm.AvailableBytes = value
	case metricMemCacheBytes:
		mm.CacheBytes = value
	case metricMemCacheBytesPeak:
		mm.CacheBytesPeak = value
	case metricMemCacheFaultsTotal:
		mm.CacheFaultsTotal = value
	case metricMemCommitLimit:
		mm.CommitLimit = value
	case metricMemCommittedBytes:
		mm.CommittedBytes = value
	case metricMemDemandZeroFaultsTotal:
		mm.DemandZeroFaultsTotal = value
	case metricMemFreeAndZeroPageListBytes:
		mm.FreeAndZeroPageListBytes = value
	case metricMemFreeSystemPageTableEntries:
		mm.FreeSystemPageTableEntries = value
	case metricMemModifiedPageListBytes:
		mm.ModifiedPageListBytes = value
	case metricMemPageFaultsTotal:
		mm.PageFaultsTotal = value
	case metricMemSwapPageReadsTotal:
		mm.SwapPageReadsTotal = value
	case metricMemSwapPagesReadTotal:
		mm.SwapPagesReadTotal = value
	case metricMemSwapPagesWrittenTotal:
		mm.SwapPagesWrittenTotal = value
	case metricMemSwapPageOperationsTotal:
		mm.SwapPageOperationsTotal = value
	case metricMemSwapPageWritesTotal:
		mm.SwapPageWritesTotal = value
	case metricMemPoolNonpagedAllocsTotal:
		mm.PoolNonPagedAllocsTotal = value
	case metricMemPoolNonpagedBytesTotal:
		mm.PoolNonPagedBytes = value
	case metricMemPoolPagedAllocsTotal:
		mm.PoolPagedAllocsTotal = value
	case metricMemPoolPagedBytes:
		mm.PoolPagedBytes = value
	case metricMemPoolPagedResidentBytes:
		mm.PoolPagedResidentBytes = value
	case metricMemStandbyCacheCoreBytes:
		mm.StandbyCacheCoreBytes = value
	case metricMemStandbyCacheNormalPriorityBytes:
		mm.StandbyCacheNormalPriorityBytes = value
	case metricMemStandbyCacheReserveBytes:
		mm.StandbyCacheReserveBytes = value
	case metricMemSystemCacheResidentBytes:
		mm.SystemCacheResidentBytes = value
	case metricMemSystemCodeResidentBytes:
		mm.SystemCodeResidentBytes = value
	case metricMemSystemCodeTotalBytes:
		mm.SystemCodeTotalBytes = value
	case metricMemSystemDriverResidentBytes:
		mm.SystemDriverResidentBytes = value
	case metricMemSystemDriverTotalBytes:
		mm.SystemDriverTotalBytes = value
	case metricMemTransitionFaultsTotal:
		mm.TransitionFaultsTotal = value
	case metricMemTransitionPagesRepurposedTotal:
		mm.TransitionPagesRePurposedTotal = value
	case metricMemWriteCopiesTotal:
		mm.WriteCopiesTotal = value
	}
}

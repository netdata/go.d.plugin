package wmi

import (
	"sort"
	"strconv"
)

func newMetrics() *metrics {
	return &metrics{
		CollectDuration: make(map[string]float64),
	}
}

type metrics struct {
	// https://github.com/martinlindhe/wmi_exporter/tree/master/docs
	CPU             *cpu               `stm:"cpu"`
	Net             *network           `stm:"net"`
	Memory          *memory            `stm:"memory"`
	CS              *cs                `stm:"cs"`
	OS              *os                `stm:"os"`
	System          *system            `stm:"system"`
	CollectDuration map[string]float64 `stm:",1000,1"`
}

type (
	cpu struct {
		cpuTimeTotal `stm:""`
		Cores        cpuCores `stm:"core"`
	}

	cpuCores []*cpuCore

	// Win32_PerfRawData_PerfOS_Processor
	// https://msdn.microsoft.com/en-us/ie/aa394317(v=vs.94)
	cpuCore struct {
		STMKey                string
		ID                    string
		id                    int // for sorting
		cpuCStateSecondsTotal `stm:""`
		cpuTimeTotal          `stm:""`
		DPCsTotal             float64 `stm:"dpcs,1000,1"`       // DPCsQueuedPersec
		InterruptsTotal       float64 `stm:"interrupts,1000,1"` // InterruptsPersec
	}

	cpuTimeTotal struct {
		Idle       float64 `stm:"idle,1000,1"`       // PercentIdleTime
		Interrupt  float64 `stm:"interrupt,1000,1"`  // PercentInterruptTime
		DPC        float64 `stm:"dpc,1000,1"`        // PercentDPCTime
		Privileged float64 `stm:"privileged,1000,1"` // PercentPrivilegedTime
		User       float64 `stm:"user,1000,1"`       // PercentUserTime
	}

	cpuCStateSecondsTotal struct {
		C1 float64 `stm:"c1,1000,1"` // PercentC1Time
		C2 float64 `stm:"c2,1000,1"` // PercentC2Time
		C3 float64 `stm:"c3,1000,1"` // PercentC3Time
	}
)

type (
	network struct {
		NICs netNICs `stm:""`
	}

	netNICs []*netNIC

	// Win32_PerfRawData_Tcpip_NetworkInterface
	// https://docs.microsoft.com/en-us/previous-versions/aa394293(v%3Dvs.85)
	netNIC struct {
		STMKey string
		ID     string

		BytesReceivedTotal       float64 `stm:"bytes_received,1000,1"` // BytesReceivedPersec
		BytesSentTotal           float64 `stm:"bytes_sent,1000,1"`     // BytesSentPersec
		BytesTotal               float64 `stm:"bytes_total,1000,1"`    // BytesTotalPersec
		PacketsOutboundDiscarded float64 `stm:"packets_outbound_discarded,1000,1"`
		PacketsOutboundErrors    float64 `stm:"packets_outbound_errors,1000,1"`
		PacketsTotal             float64 `stm:"packets_total,1000,1"` // PacketsPersec
		PacketsReceivedDiscarded float64 `stm:"packets_received_discarded,1000,1"`
		PacketsReceivedErrors    float64 `stm:"packets_received_errors,1000,1"`
		PacketsReceivedTotal     float64 `stm:"packets_received_total,1000,1"` // PacketsReceivedPersec
		PacketsReceivedUnknown   float64 `stm:"packets_received_unknown,1000,1"`
		PacketsSentTotal         float64 `stm:"packets_sent_total,1000,1"` // PacketsSentPersec
		CurrentBandwidth         float64 `stm:"current_bandwidth"`
	}
)

type (
	// Win32_PerfRawData_PerfOS_Memory
	// https://technet.microsoft.com/en-ca/aa394314(v=vs.71)
	// http://wutils.com/wmi/root/cimv2/win32_perfrawdata_perfos_memory/
	memory struct {
		AvailableBytes                  float64 `stm:"available_bytes,1000,1"`
		CacheBytes                      float64 `stm:"cache_bytes,1000,1"`
		CacheBytesPeak                  float64 `stm:"cache_bytes_peak,1000,1"`
		CacheFaultsTotal                float64 `stm:"cache_faults_total,1000,1"` // CacheFaultsPersec
		CommitLimit                     float64 `stm:"commit_limit,1000,1"`
		CommittedBytes                  float64 `stm:"committed_bytes,1000,1"`
		DemandZeroFaultsTotal           float64 `stm:"demand_zero_faults_total,1000,1"` // DemandZeroFaultsPersec
		FreeAndZeroPageListBytes        float64 `stm:"free_and_zero_page_list_bytes,1000,1"`
		FreeSystemPageTableEntries      float64 `stm:"free_system_page_table_entries,1000,1"`
		ModifiedPageListBytes           float64 `stm:"modified_page_list_bytes,1000,1"`
		PageFaultsTotal                 float64 `stm:"page_faults_total,1000,1"`          // PageFaultsPersec
		SwapPageReadsTotal              float64 `stm:"swap_page_reads_total,1000,1"`      // PageReadsPersec
		SwapPagesReadTotal              float64 `stm:"swap_pages_read_total,1000,1"`      // PagesInputPersec
		SwapPagesWrittenTotal           float64 `stm:"swap_pages_written_total,1000,1"`   // PagesOutputPersec
		SwapPageOperationsTotal         float64 `stm:"swap_page_operations_total,1000,1"` // PagesPersec
		SwapPageWritesTotal             float64 `stm:"swap_page_writes_total,1000,1"`     // PageWritesPersec
		PoolNonPagedAllocsTotal         float64 `stm:"pool_nonpaged_allocs_total,1000,1"` // PoolNonPagedAllocs
		PoolNonPagedBytes               float64 `stm:"pool_nonpaged_bytes_total,1000,1"`
		PoolPagedAllocsTotal            float64 `stm:"pool_paged_allocs_total,1000,1"` // PoolPagedAllocs
		PoolPagedBytes                  float64 `stm:"pool_paged_bytes,1000,1"`
		PoolPagedResidentBytes          float64 `stm:"pool_paged_resident_bytes,1000,1"`
		StandbyCacheCoreBytes           float64 `stm:"standby_cache_core_bytes,1000,1"`
		StandbyCacheNormalPriorityBytes float64 `stm:"standby_cache_normal_priority_bytes,1000,1"`
		StandbyCacheReserveBytes        float64 `stm:"standby_cache_reserve_bytes,1000,1"`
		SystemCacheResidentBytes        float64 `stm:"system_cache_resident_bytes,1000,1"`
		SystemCodeResidentBytes         float64 `stm:"system_code_resident_bytes,1000,1"`
		SystemCodeTotalBytes            float64 `stm:"system_code_total_bytes,1000,1"`
		SystemDriverResidentBytes       float64 `stm:"system_driver_resident_bytes,1000,1"`
		SystemDriverTotalBytes          float64 `stm:"system_driver_total_bytes,1000,1"`
		TransitionFaultsTotal           float64 `stm:"transition_faults_total,1000,1"`           // TransitionFaultsPersec
		TransitionPagesRePurposedTotal  float64 `stm:"transition_pages_repurposed_total,1000,1"` // TransitionPagesRePurposedPersec
		WriteCopiesTotal                float64 `stm:"write_copies_total,1000,1"`                // WriteCopiesPersec
	}
	// Win32_ComputerSystem
	// https://msdn.microsoft.com/en-us/library/aa394102
	cs struct {
		LogicalProcessors   float64 `stm:"logical_processors"`           // NumberOfLogicalProcessors
		PhysicalMemoryBytes float64 `stm:"physical_memory_bytes,1000,1"` // TotalPhysicalMemory
	}
	// Win32_PerfRawData_PerfOS_System
	// https://docs.microsoft.com/en-us/previous-versions/aa394272(v%3Dvs.85)
	system struct {
		ContextSwitchesTotal     float64 `stm:"context_switches_total,1000,1"`     // ContextSwitchesPersec
		ExceptionDispatchesTotal float64 `stm:"exception_dispatches_total,1000,1"` // ExceptionDispatchesPersec
		ProcessorQueueLength     float64 `stm:"processor_queue_length"`
		SystemCallsTotal         float64 `stm:"system_calls_total,1000,1"` // SystemCallsPersec
		SystemUpTime             float64 `stm:"system_up_time"`
		Threads                  float64 `stm:"system_threads"`
	}
	// Win32_OperatingSystem
	// https://docs.microsoft.com/en-us/windows/desktop/CIMWin32Prov/win32-operatingsystem
	os struct {
		PhysicalMemoryFreeBytes float64 `stm:"physical_memory_free_bytes,1000,1"` // FreePhysicalMemory
		PagingFreeBytes         float64 `stm:"paging_free_bytes,1000,1"`          // FreeSpaceInPagingFiles
		VirtualMemoryFreeBytes  float64 `stm:"virtual_memory_free_bytes,1000,1"`  // FreeVirtualMemory
		ProcessesLimit          float64 `stm:"processes_limit"`                   // MaxNumberOfProcesses
		ProcessMemoryLimitBytes float64 `stm:"process_memory_limit_bytes,1000,1"` // MaxProcessMemorySize
		Processes               float64 `stm:"processes"`                         // NumberOfProcesses
		Users                   float64 `stm:"users"`                             // NumberOfUsers
		PagingLimitBytes        float64 `stm:"paging_limit_bytes,1000,1"`         // SizeStoredInPagingFiles
		VirtualMemoryBytes      float64 `stm:"virtual_memory_bytes,1000,1"`       // TotalVirtualMemorySize
		VisibleMemoryBytes      float64 `stm:"visible_memory_bytes,1000,1"`       // TotalVisibleMemorySize
		Time                    float64 `stm:"time"`                              // LocalDateTime
		// Timezone                float64 `stm:"timezone"`                          // LocalDateTime
	}
)

func newCPUCore(id string) *cpuCore { return &cpuCore{STMKey: id, ID: id, id: getCPUIntID(id)} }

func (cc *cpuCores) get(id string, createIfNotExist bool) (core *cpuCore) {
	for _, c := range *cc {
		if c.ID == id {
			return c
		}
	}
	if createIfNotExist {
		core = newCPUCore(id)
		*cc = append(*cc, core)
	}
	return core
}

func (cc *cpuCores) sort() { sort.Slice(*cc, func(i, j int) bool { return (*cc)[i].id < (*cc)[j].id }) }

func getCPUIntID(id string) int {
	if id == "" {
		return -1
	}
	v, err := strconv.Atoi(string(id[len(id)-1]))
	if err != nil {
		return -1
	}
	return v
}

func newNIC(id string) *netNIC { return &netNIC{STMKey: id, ID: id} }

func (ns *netNICs) get(id string, createIfNotExist bool) (nic *netNIC) {
	for _, n := range *ns {
		if n.ID == id {
			return n
		}
	}
	if createIfNotExist {
		nic = newNIC(id)
		*ns = append(*ns, nic)
	}
	return nic
}

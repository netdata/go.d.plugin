// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import "strconv"

type metrics struct {
	// https://github.com/prometheus-community/windows_exporter/tree/master/docs#documentation
	CPU         *cpuMetrics         `stm:"cpu"`
	Memory      *memoryMetrics      `stm:"memory"`
	Net         *networkMetrics     `stm:"net"`
	LogicalDisk *logicalDiskMetrics `stm:"logical_disk"`
	OS          *osMetrics          `stm:"os"`
	System      *systemMetrics      `stm:"system"`
	Logon       *logonMetrics       `stm:"logon"`
	ThermalZone *thermalZoneMetrics `stm:"thermalzone"`
	Collectors  *collectors         `stm:""`
}

func (m metrics) hasCPU() bool         { return m.CPU != nil }
func (m metrics) hasMemory() bool      { return m.Memory != nil }
func (m metrics) hasNet() bool         { return m.Net != nil }
func (m metrics) hasLogicalDisk() bool { return m.LogicalDisk != nil }
func (m metrics) hasOS() bool          { return m.OS != nil }
func (m metrics) hasSystem() bool      { return m.System != nil }
func (m metrics) hasLogon() bool       { return m.Logon != nil }
func (m metrics) hasThermalZone() bool { return m.ThermalZone != nil }
func (m metrics) hasCollectors() bool  { return m.Collectors != nil }

// cpu
type (
	cpuMetrics struct {
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

// Win32_PerfRawData_PerfOS_Memory
// https://technet.microsoft.com/en-ca/aa394314(v=vs.71)
// http://wutils.com/wmi/root/cimv2/win32_perfrawdata_perfos_memory/
type memoryMetrics struct {
	UsedBytes         *float64 `stm:"used_bytes,1000,1"`          // os.VisibleMemoryBytes - AvailableBytes
	NotCommittedBytes float64  `stm:"not_committed_bytes,1000,1"` // CommitLimit - CommittedBytes
	StandbyCacheTotal float64  `stm:"standby_cache_total,1000,1"` // StandbyCacheCoreBytes + StandbyCacheNormalPriorityBytes + StandbyCacheReserveBytes
	Cached            float64  `stm:"cache_total,1000,1"`         // StandbyCacheTotal + ModifiedPageListBytes

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

// network
type (
	networkMetrics struct {
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

// logical disk
type (
	logicalDiskMetrics struct {
		Volumes volumes `stm:""`
	}

	volumes []*volume

	// Win32_PerfRawData_PerfDisk_LogicalDisk
	// https://msdn.microsoft.com/en-us/windows/hardware/aa394307(v=vs.71)
	volume struct {
		STMKey string
		ID     string

		UsedSpace float64 `stm:"used_space,1000,1"` // TotalSpace - FreeSpace

		RequestsQueued  float64 `stm:"requests_queued"`            // CurrentDiskQueueLength
		ReadBytesTotal  float64 `stm:"read_bytes_total,1000,1"`    // DiskReadBytesPerSec
		ReadsTotal      float64 `stm:"reads_total"`                // DiskReadsPerSec
		WriteBytesTotal float64 `stm:"write_bytes_total,1000,1"`   // DiskWriteBytesPerSec
		WritesTotal     float64 `stm:"writes_total"`               // DiskWritesPerSec
		ReadTime        float64 `stm:"read_seconds_total,1000,1"`  // PercentDiskReadTime
		WriteTime       float64 `stm:"write_seconds_total,1000,1"` // PercentDiskWriteTime
		TotalSpace      float64 `stm:"total_space,1000,1"`         // PercentFreeSpace_Base
		FreeSpace       float64 `stm:"free_space,1000,1"`          // PercentFreeSpace
		IdleTime        float64 `stm:"idle_seconds_total,1000,1"`  // PercentIdleTime
		SplitIOs        float64 `stm:"split_ios_total"`            // SplitIOPerSec
		ReadLatency     float64 `stm:"read_latency,1000,1"`        //AvgDiskSecPerRead
		WriteLatency    float64 `stm:"write_latency,1000,1"`       // AvgDiskSecPerWrite
	}
)

// Win32_perfrawdata_counters_thermalzoneinformation
// https://wutils.com/wmi/root/cimv2/win32_perfrawdata_counters_thermalzoneinformation/#temperature_properties
type (
	thermalZoneMetrics struct {
		Zones thermalZones `stm:""`
	}

	thermalZones []*thermalZone

	thermalZone struct {
		STMKey string
		ID     string

		Temperature float64 `stm:"temperature,1000,1"`
	}
)

// Win32_PerfRawData_PerfOS_System
// https://docs.microsoft.com/en-us/previous-versions/aa394272(v%3Dvs.85)
type systemMetrics struct {
	SystemUpTime int64 `stm:"up_time"`

	ContextSwitchesTotal     float64 `stm:"context_switches_total,1000,1"`     // ContextSwitchesPersec
	ExceptionDispatchesTotal float64 `stm:"exception_dispatches_total,1000,1"` // ExceptionDispatchesPersec
	ProcessorQueueLength     float64 `stm:"processor_queue_length"`
	SystemCallsTotal         float64 `stm:"calls_total,1000,1"` // SystemCallsPersec
	SystemBootTime           float64 `stm:"boot_time"`
	Threads                  float64 `stm:"threads"`
}

// Win32_OperatingSystem
// https://docs.microsoft.com/en-us/windows/desktop/CIMWin32Prov/win32-operatingsystem
type osMetrics struct {
	PagingLimitBytes float64 `stm:"paging_limit_bytes,1000,1"` // SizeStoredInPagingFiles
	PagingFreeBytes  float64 `stm:"paging_free_bytes,1000,1"`  // FreeSpaceInPagingFiles
	PagingUsedBytes  float64 `stm:"paging_used_bytes,1000,1"`  // PagingLimitBytes - PagingFreeBytes

	VisibleMemoryBytes      float64 `stm:"visible_memory_bytes,1000,1"`       // TotalVisibleMemorySize
	PhysicalMemoryFreeBytes float64 `stm:"physical_memory_free_bytes,1000,1"` // FreePhysicalMemory
	VisibleMemoryUsedBytes  float64 `stm:"visible_memory_used_bytes,1000,1"`  // VisibleMemoryBytes - PhysicalMemoryFreeBytes

	VirtualMemoryBytes     float64 `stm:"virtual_memory_bytes,1000,1"`      // TotalVirtualMemorySize
	VirtualMemoryFreeBytes float64 `stm:"virtual_memory_free_bytes,1000,1"` // FreeVirtualMemory

	ProcessesLimit          float64 `stm:"processes_limit"`                   // MaxNumberOfProcesses
	Processes               float64 `stm:"processes"`                         // NumberOfProcesses
	ProcessMemoryLimitBytes float64 `stm:"process_memory_limit_bytes,1000,1"` // MaxProcessMemorySize

	Users float64 `stm:"users"` // NumberOfUsers
	Time  float64 `stm:"time"`  // LocalDateTime
	// Timezone                float64 `stm:"timezone"`                          // LocalDateTime
}

// Win32_LogonSession
// https://docs.microsoft.com/en-us/windows/win32/cimwin32prov/win32-logonsession
type logonMetrics struct {
	Type struct {
		System                  float64 `stm:"system"`
		Interactive             float64 `stm:"interactive"`
		Network                 float64 `stm:"network"`
		Batch                   float64 `stm:"batch"`
		Service                 float64 `stm:"service"`
		Proxy                   float64 `stm:"proxy"`
		Unlock                  float64 `stm:"unlock"`
		NetworkCleartext        float64 `stm:"network_clear_text"`
		NewCredentials          float64 `stm:"new_credentials"`
		RemoteInteractive       float64 `stm:"remote_interactive"`
		CachedInteractive       float64 `stm:"cached_interactive"`
		CachedRemoteInteractive float64 `stm:"cached_remote_interactive"`
		CachedUnlock            float64 `stm:"cached_unlock"`
	} `stm:"type"`
}

type (
	collectors []*collector
	collector  struct {
		STMKey string
		ID     string

		Duration float64 `stm:"collection_duration,1000,1"`
		Success  bool    `stm:"collection_success"`
	}
)

func newCollector(id string) *collector     { return &collector{STMKey: id, ID: id} }
func newCPUCore(id string) *cpuCore         { return &cpuCore{STMKey: id, ID: id, id: getCPUIntID(id)} }
func newNIC(id string) *netNIC              { return &netNIC{STMKey: id, ID: id} }
func newVolume(id string) *volume           { return &volume{STMKey: id, ID: id} }
func newThermalZone(id string) *thermalZone { return &thermalZone{STMKey: id, ID: id} }

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

func (cs *collectors) get(id string) *collector {
	for _, cr := range *cs {
		if cr.ID == id {
			return cr
		}
	}
	cr := newCollector(id)
	*cs = append(*cs, cr)
	return cr
}

func (cc *cpuCores) get(id string) *cpuCore {
	for _, core := range *cc {
		if core.ID == id {
			return core
		}
	}
	core := newCPUCore(id)
	*cc = append(*cc, core)
	return core
}

func (ns *netNICs) get(id string) *netNIC {
	for _, nic := range *ns {
		if nic.ID == id {
			return nic
		}
	}
	nic := newNIC(id)
	*ns = append(*ns, nic)
	return nic
}

func (vs *volumes) get(id string) *volume {
	for _, v := range *vs {
		if v.ID == id {
			return v
		}
	}
	vol := newVolume(id)
	*vs = append(*vs, vol)
	return vol
}

func (tz *thermalZones) get(id string) *thermalZone {
	for _, v := range *tz {
		if v.ID == id {
			return v
		}
	}
	v := newThermalZone(id)
	*tz = append(*tz, v)
	return v
}

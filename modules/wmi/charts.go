// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import (
	"fmt"
	"sort"
	"strings"

	"github.com/netdata/go.d.plugin/agent/module"
)

const (
	prioCPUUtil = module.Priority + iota
	prioCPUDPCs
	prioCPUInterrupts
	prioCPUCoreUtil
	prioCPUCoreCState

	prioMemUtil
	prioMemPageFaults
	prioMemSwapUtil
	prioMemSwapOperations
	prioMemSwapPages
	prioMemCache
	prioMemCacheFaults
	prioMemSystemPool

	prioNICBandwidth
	prioNICPackets
	prioNICErrors
	prioNICDiscards

	prioDiskUtil
	prioDiskBandwidth
	prioDiskOperations
	prioDiskAvgLatency

	prioOSProcesses
	prioOSUsers
	prioOSVisibleMemoryUsage
	prioOSPagingUsage

	prioSystemThreads
	prioSystemUptime

	prioLogonSessions

	prioTCPConnsEstablished
	prioTCPConnsActive
	prioTCPConnsPassive
	prioTCPConnsFailure
	prioTCPConnsReset
	prioTCPSegmentsReceived
	prioTCPSegmentsSent
	prioTCPSegmentsRetransmitted

	prioThermalzoneTemperature

	prioCollectionDuration
	prioCollectionStatus

	prioProcessesCPUTimeTotal
	prioProcessesHandles
	prioProcessesIOBytes
	prioProcessesIOOperations
	prioProcessesPageFaults
	prioProcessesPageFileBytes
	prioProcessesPoolBytes
	prioProcessesThreads

	prioServiceStateContinuePending
	prioServiceStatePausePending
	prioServiceStatePaused
	prioServiceStateRunning
	prioServiceStateStartPending
	prioServiceStateStopPending
	prioServiceStateStopped
	prioServiceStateUnknown
	prioServiceStatusDegraded
	prioServiceStatusError
	prioServiceStatusLostConn
	prioServiceStatusNoContact
	prioServiceStatusNonRecover
	prioServiceStatusOK
	prioServiceStatusPredFail
	prioServiceStatusService
	prioServiceStatusStarting
	prioServiceStatusStopping
	prioServiceStatusStressed
	prioServiceStatusUnknown
)

func newProcessesCharts() module.Charts {
	return module.Charts{
		processesCPUTimeTotalChart.Copy(),
		processesHandlesChart.Copy(),
		processesIOBytesChart.Copy(),
		processesIOOperationsChart.Copy(),
		processesPageFaultsChart.Copy(),
		processesPageFileBytes.Copy(),
		processesPoolBytes.Copy(),
		processesThreads.Copy(),
	}
}

var (
	processesCPUTimeTotalChart = module.Chart{
		ID:       "processes_cpu_time",
		Title:    "CPU usage",
		Units:    "percentage",
		Fam:      "processes",
		Ctx:      "wmi.processes_cpu_time",
		Type:     module.Stacked,
		Priority: prioProcessesCPUTimeTotal,
	}
	processesHandlesChart = module.Chart{
		ID:       "processes_handles",
		Title:    "Number of handles open",
		Units:    "handles",
		Fam:      "processes",
		Ctx:      "wmi.processes_handles",
		Type:     module.Stacked,
		Priority: prioProcessesHandles,
	}
	processesIOBytesChart = module.Chart{
		ID:       "processes_io_bytes",
		Title:    "Total of IO bytes (read, write, other)",
		Units:    "bytes/s",
		Fam:      "processes",
		Ctx:      "wmi.processes_io_bytes",
		Type:     module.Stacked,
		Priority: prioProcessesIOBytes,
	}
	processesIOOperationsChart = module.Chart{
		ID:       "processes_io_operations",
		Title:    "Total of IO events (read, write, other)",
		Units:    "operations/s",
		Fam:      "processes",
		Ctx:      "wmi.processes_io_operations",
		Type:     module.Stacked,
		Priority: prioProcessesIOOperations,
	}
	processesPageFaultsChart = module.Chart{
		ID:       "processes_page_faults",
		Title:    "Number of page faults",
		Units:    "pgfaults/s",
		Fam:      "processes",
		Ctx:      "wmi.processes_page_faults",
		Type:     module.Stacked,
		Priority: prioProcessesPageFaults,
	}
	processesPageFileBytes = module.Chart{
		ID:       "processes_page_file_bytes",
		Title:    "Bytes used in page file(s)",
		Units:    "bytes",
		Fam:      "processes",
		Ctx:      "wmi.processes_file_bytes",
		Type:     module.Stacked,
		Priority: prioProcessesPageFileBytes,
	}
	processesPoolBytes = module.Chart{
		ID:       "processes_pool_bytes",
		Title:    "Last observed bytes in paged",
		Units:    "bytes",
		Fam:      "processes",
		Ctx:      "wmi.processes_pool_bytes",
		Type:     module.Stacked,
		Priority: prioProcessesPoolBytes,
	}
	processesThreads = module.Chart{
		ID:       "processes_threads",
		Title:    "Active threads",
		Units:    "threads",
		Fam:      "processes",
		Ctx:      "wmi.processes_threads",
		Type:     module.Stacked,
		Priority: prioProcessesThreads,
	}
)

func newCPUCharts() module.Charts {
	return module.Charts{
		cpuUtilChart.Copy(),
		cpuDPCsChart.Copy(),
		cpuInterruptsChart.Copy(),
	}
}

var (
	cpuUtilChart = module.Chart{
		ID:       "cpu_utilization_total",
		Title:    "Total CPU Utilization (all cores)",
		Units:    "percentage",
		Fam:      "cpu",
		Ctx:      "wmi.cpu_utilization_total",
		Type:     module.Stacked,
		Priority: prioCPUUtil,
		Dims: module.Dims{
			{ID: "cpu_idle", Name: "idle", Algo: module.PercentOfIncremental, Div: 1000, DimOpts: module.DimOpts{Hidden: true}},
			{ID: "cpu_dpc", Name: "dpc", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_user", Name: "user", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_privileged", Name: "privileged", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_interrupt", Name: "interrupt", Algo: module.PercentOfIncremental, Div: 1000},
		},
	}
	cpuDPCsChart = module.Chart{
		ID:       "cpu_dpcs",
		Title:    "Received and Serviced Deferred Procedure Calls (DPC)",
		Units:    "dpc/s",
		Fam:      "cpu",
		Ctx:      "wmi.cpu_dpcs",
		Type:     module.Stacked,
		Priority: prioCPUDPCs,
	}
	cpuInterruptsChart = module.Chart{
		ID:       "cpu_interrupts",
		Title:    "Received and Serviced Hardware Interrupts",
		Units:    "interrupts/s",
		Fam:      "cpu",
		Ctx:      "wmi.cpu_interrupts",
		Type:     module.Stacked,
		Priority: prioCPUInterrupts,
	}
)

func newCPUCoreCharts() module.Charts {
	return module.Charts{
		cpuCoreUtilChart.Copy(),
		cpuCoreCStateChart.Copy(),
	}
}

var (
	cpuCoreUtilChart = module.Chart{
		ID:       "core_%s_cpu_utilization",
		Title:    "Core%s CPU Utilization",
		Units:    "percentage",
		Fam:      "cpu",
		Ctx:      "wmi.cpu_utilization",
		Type:     module.Stacked,
		Priority: prioCPUCoreUtil,
		Dims: module.Dims{
			{ID: "cpu_core_%s_idle", Name: "idle", Algo: module.PercentOfIncremental, Div: 1000, DimOpts: module.DimOpts{Hidden: true}},
			{ID: "cpu_core_%s_dpc", Name: "dpc", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_core_%s_user", Name: "user", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_core_%s_privileged", Name: "privileged", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_core_%s_interrupt", Name: "interrupt", Algo: module.PercentOfIncremental, Div: 1000},
		},
	}
	cpuCoreCStateChart = module.Chart{
		ID:       "core_%s_cpu_cstate",
		Title:    "Core%s Time Spent in Low-Power Idle State",
		Units:    "percentage",
		Fam:      "cpu",
		Ctx:      "wmi.cpu_cstate",
		Type:     module.Stacked,
		Priority: prioCPUCoreCState,
		Dims: module.Dims{
			{ID: "cpu_core_%s_c1", Name: "c1", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_core_%s_c2", Name: "c2", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_core_%s_c3", Name: "c3", Algo: module.PercentOfIncremental, Div: 1000},
		},
	}
)

func newMemCharts() module.Charts {
	return module.Charts{
		memUtilChart.Copy(),
		memPageFaultsChart.Copy(),
		memSwapUtilChart.Copy(),
		memSwapOperationsChart.Copy(),
		memSwapPagesChart.Copy(),
		memCacheChart.Copy(),
		memCacheFaultsChart.Copy(),
		memSystemPoolChart.Copy(),
	}
}

var (
	memUtilChart = module.Chart{
		ID:       "memory_utilization",
		Title:    "Memory Utilization",
		Units:    "KiB",
		Fam:      "mem",
		Ctx:      "wmi.memory_utilization",
		Type:     module.Stacked,
		Priority: prioMemUtil,
		Dims: module.Dims{
			{ID: "memory_available_bytes", Name: "available", Div: 1000 * 1024},
			{ID: "memory_used_bytes", Name: "used", Div: 1000 * 1024},
		},
		Vars: module.Vars{
			{ID: "os_visible_memory_bytes"},
		},
	}
	memPageFaultsChart = module.Chart{
		ID:       "memory_page_faults",
		Title:    "Memory Page Faults",
		Units:    "events/s",
		Fam:      "mem",
		Ctx:      "wmi.memory_page_faults",
		Priority: prioMemPageFaults,
		Dims: module.Dims{
			{ID: "memory_page_faults_total", Name: "page faults", Algo: module.Incremental, Div: 1000},
		},
	}
	memSwapUtilChart = module.Chart{
		ID:       "memory_swap_utilization",
		Title:    "Swap Utilization",
		Units:    "KiB",
		Fam:      "mem",
		Ctx:      "wmi.memory_swap_utilization",
		Type:     module.Stacked,
		Priority: prioMemSwapUtil,
		Dims: module.Dims{
			{ID: "memory_not_committed_bytes", Name: "available", Div: 1000 * 1024},
			{ID: "memory_committed_bytes", Name: "used", Div: 1000 * 1024},
		},
		Vars: module.Vars{
			{ID: "memory_commit_limit"},
		},
	}
	memSwapOperationsChart = module.Chart{
		ID:       "memory_swap_operations",
		Title:    "Swap Operations",
		Units:    "operations/s",
		Fam:      "mem",
		Ctx:      "wmi.memory_swap_operations",
		Type:     module.Area,
		Priority: prioMemSwapOperations,
		Dims: module.Dims{
			{ID: "memory_swap_page_reads_total", Name: "read", Algo: module.Incremental, Div: 1000},
			{ID: "memory_swap_page_writes_total", Name: "write", Algo: module.Incremental, Div: -1000},
		},
	}
	memSwapPagesChart = module.Chart{
		ID:       "memory_swap_pages",
		Title:    "Swap Pages",
		Units:    "pages/s",
		Fam:      "mem",
		Ctx:      "wmi.memory_swap_pages",
		Type:     module.Area,
		Priority: prioMemSwapPages,
		Dims: module.Dims{
			{ID: "memory_swap_pages_read_total", Name: "read", Algo: module.Incremental, Div: 1000},
			{ID: "memory_swap_pages_written_total", Name: "written", Algo: module.Incremental, Div: -1000},
		},
	}
	memCacheChart = module.Chart{
		ID:       "memory_cached",
		Title:    "Cached",
		Units:    "KiB",
		Fam:      "mem",
		Ctx:      "wmi.memory_cached",
		Priority: prioMemCache,
		Dims: module.Dims{
			{ID: "memory_cache_total", Name: "cached", Div: 1000 * 1024},
		},
	}
	memCacheFaultsChart = module.Chart{
		ID:       "memory_cache_faults",
		Title:    "Cache Faults",
		Units:    "events/s",
		Fam:      "mem",
		Ctx:      "wmi.memory_cache_faults",
		Priority: prioMemCacheFaults,
		Dims: module.Dims{
			{ID: "memory_cache_faults_total", Name: "cache faults", Algo: module.Incremental, Div: 1000},
		},
	}
	memSystemPoolChart = module.Chart{
		ID:       "memory_system_pool",
		Title:    "System Memory Pool",
		Units:    "KiB",
		Fam:      "mem",
		Ctx:      "wmi.memory_system_pool",
		Type:     module.Stacked,
		Priority: prioMemSystemPool,
		Dims: module.Dims{
			{ID: "memory_pool_paged_bytes", Name: "paged", Div: 1000 * 1024},
			{ID: "memory_pool_nonpaged_bytes_total", Name: "non-paged", Div: 1000 * 1024},
		},
	}
)

func newNICCharts() module.Charts {
	return module.Charts{
		nicBandwidthChart.Copy(),
		nicPacketsChart.Copy(),
		nicErrorsChart.Copy(),
		nicDiscardsChart.Copy(),
	}
}

var (
	nicBandwidthChart = module.Chart{
		ID:       "nic_%s_bandwidth",
		Title:    "Bandwidth %s",
		Units:    "kilobits/s",
		Fam:      "net",
		Ctx:      "wmi.net_bandwidth",
		Type:     module.Area,
		Priority: prioNICBandwidth,
		Dims: module.Dims{
			{ID: "net_%s_bytes_received", Name: "received", Algo: module.Incremental, Div: 1000 * 125},
			{ID: "net_%s_bytes_sent", Name: "sent", Algo: module.Incremental, Div: -1000 * 125},
		},
		Vars: module.Vars{
			{ID: "net_%s_current_bandwidth"},
		},
	}
	nicPacketsChart = module.Chart{
		ID:       "nic_%s_packets",
		Title:    "Packets %s",
		Units:    "packets/s",
		Fam:      "net",
		Ctx:      "wmi.net_packets",
		Type:     module.Area,
		Priority: prioNICPackets,
		Dims: module.Dims{
			{ID: "net_%s_packets_received_total", Name: "received", Algo: module.Incremental, Div: 1000},
			{ID: "net_%s_packets_sent_total", Name: "sent", Algo: module.Incremental, Div: -1000},
		},
	}
	nicErrorsChart = module.Chart{
		ID:       "nic_%s_errors",
		Title:    "Errors %s",
		Units:    "errors/s",
		Fam:      "net",
		Ctx:      "wmi.net_errors",
		Type:     module.Area,
		Priority: prioNICErrors,
		Dims: module.Dims{
			{ID: "net_%s_packets_received_errors", Name: "inbound", Algo: module.Incremental, Div: 1000},
			{ID: "net_%s_packets_outbound_errors", Name: "outbound", Algo: module.Incremental, Div: -1000},
		},
	}
	nicDiscardsChart = module.Chart{
		ID:       "nic_%s_discarded",
		Title:    "Discards %s",
		Units:    "discards/s",
		Fam:      "net",
		Ctx:      "wmi.net_discarded",
		Type:     module.Area,
		Priority: prioNICDiscards,
		Dims: module.Dims{
			{ID: "net_%s_packets_received_discarded", Name: "inbound", Algo: module.Incremental, Div: 1000},
			{ID: "net_%s_packets_outbound_discarded", Name: "outbound", Algo: module.Incremental, Div: -1000},
		},
	}
)

func newTCPCharts() module.Charts {
	return module.Charts{
		tcpConnsActiveChart.Copy(),
		tcpConnsEstablishedChart.Copy(),
		tcpConnsFailuresChart.Copy(),
		tcpConnsPassiveChart.Copy(),
		tcpConnsResetsChart.Copy(),
		tcpSegmentsReceivedChart.Copy(),
		tcpSegmentsRetransmittedChart.Copy(),
		tcpSegmentsSentChart.Copy(),
	}
}

var (
	tcpConnsEstablishedChart = module.Chart{
		ID:       "tcp_conns_established",
		Title:    "TCP established connections",
		Units:    "connections",
		Fam:      "tcp",
		Ctx:      "wmi.tcp_conns_established",
		Priority: prioTCPConnsEstablished,
		Dims: module.Dims{
			{ID: "tcp_conns_established_ipv4", Name: "ipv4"},
			{ID: "tcp_conns_established_ipv6", Name: "ipv6"},
		},
	}
	tcpConnsActiveChart = module.Chart{
		ID:       "tcp_conns_active",
		Title:    "TCP active connections",
		Units:    "connections/s",
		Fam:      "tcp",
		Ctx:      "wmi.tcp_conns_active",
		Priority: prioTCPConnsActive,
		Dims: module.Dims{
			{ID: "tcp_conns_active_ipv4", Name: "ipv4", Algo: module.Incremental},
			{ID: "tcp_conns_active_ipv6", Name: "ipv6", Algo: module.Incremental},
		},
	}
	tcpConnsPassiveChart = module.Chart{
		ID:       "tcp_conns_passive",
		Title:    "TCP passive connections",
		Units:    "connections/s",
		Fam:      "tcp",
		Ctx:      "wmi.tcp_conns_passive",
		Priority: prioTCPConnsPassive,
		Dims: module.Dims{
			{ID: "tcp_conns_passive_ipv4", Name: "ipv4", Algo: module.Incremental},
			{ID: "tcp_conns_passive_ipv6", Name: "ipv6", Algo: module.Incremental},
		},
	}
	tcpConnsFailuresChart = module.Chart{
		ID:       "tcp_conns_failures",
		Title:    "TCP connection failures",
		Units:    "failures/s",
		Fam:      "tcp",
		Ctx:      "wmi.tcp_conns_failures",
		Priority: prioTCPConnsFailure,
		Dims: module.Dims{
			{ID: "tcp_conns_failures_ipv4", Name: "ipv4", Algo: module.Incremental},
			{ID: "tcp_conns_failures_ipv6", Name: "ipv6", Algo: module.Incremental},
		},
	}
	tcpConnsResetsChart = module.Chart{
		ID:       "tcp_conns_resets",
		Title:    "TCP connections resets",
		Units:    "resets/s",
		Fam:      "tcp",
		Ctx:      "wmi.tcp_conns_resets",
		Priority: prioTCPConnsReset,
		Dims: module.Dims{
			{ID: "tcp_conns_resets_ipv4", Name: "ipv4", Algo: module.Incremental},
			{ID: "tcp_conns_resets_ipv6", Name: "ipv6", Algo: module.Incremental},
		},
	}
	tcpSegmentsReceivedChart = module.Chart{
		ID:       "tcp_segments_received",
		Title:    "Number of TCP segments received",
		Units:    "segments/s",
		Fam:      "tcp",
		Ctx:      "wmi.tcp_segments_received",
		Priority: prioTCPSegmentsReceived,
		Dims: module.Dims{
			{ID: "tcp_segments_received_ipv4", Name: "ipv4", Algo: module.Incremental},
			{ID: "tcp_segments_received_ipv6", Name: "ipv6", Algo: module.Incremental},
		},
	}
	tcpSegmentsSentChart = module.Chart{
		ID:       "tcp_segments_sent",
		Title:    "Number of TCP segments sent",
		Units:    "segments/s",
		Fam:      "tcp",
		Ctx:      "wmi.tcp_segments_sent",
		Priority: prioTCPSegmentsSent,
		Dims: module.Dims{
			{ID: "tcp_segments_sent_ipv4", Name: "ipv4", Algo: module.Incremental},
			{ID: "tcp_segments_sent_ipv6", Name: "ipv6", Algo: module.Incremental},
		},
	}
	tcpSegmentsRetransmittedChart = module.Chart{
		ID:       "tcp_segments_retransmitted",
		Title:    "Number of TCP segments retransmitted",
		Units:    "segments/s",
		Fam:      "tcp",
		Ctx:      "wmi.tcp_segments_retransmitted",
		Priority: prioTCPSegmentsRetransmitted,
		Dims: module.Dims{
			{ID: "tcp_segments_retransmitted_ipv4", Name: "ipv4", Algo: module.Incremental},
			{ID: "tcp_segments_retransmitted_ipv6", Name: "ipv6", Algo: module.Incremental},
		},
	}
)

func newDiskCharts() module.Charts {
	return module.Charts{
		diskUtilChart.Copy(),
		diskBandwidthChart.Copy(),
		diskOperationsChart.Copy(),
		diskAvgLatencyChart.Copy(),
	}
}

var (
	diskUtilChart = module.Chart{
		ID:       "logical_disk_%s_utilization",
		Title:    "Utilization Disk %s",
		Units:    "KiB",
		Fam:      "disk",
		Ctx:      "wmi.logical_disk_utilization",
		Type:     module.Stacked,
		Priority: prioDiskUtil,
		Dims: module.Dims{
			{ID: "logical_disk_%s_free_space", Name: "free", Div: 1000 * 1024},
			{ID: "logical_disk_%s_used_space", Name: "used", Div: 1000 * 1024},
		},
	}
	diskBandwidthChart = module.Chart{
		ID:       "logical_disk_%s_bandwidth",
		Title:    "Bandwidth Disk %s",
		Units:    "KiB/s",
		Fam:      "disk",
		Ctx:      "wmi.logical_disk_bandwidth",
		Type:     module.Area,
		Priority: prioDiskBandwidth,
		Dims: module.Dims{
			{ID: "logical_disk_%s_read_bytes_total", Name: "read", Algo: module.Incremental, Div: 1000 * 1024},
			{ID: "logical_disk_%s_write_bytes_total", Name: "write", Algo: module.Incremental, Div: -1000 * 1024},
		},
	}
	diskOperationsChart = module.Chart{
		ID:       "logical_disk_%s_operations",
		Title:    "Operations Disk %s",
		Units:    "operations/s",
		Fam:      "disk",
		Ctx:      "wmi.logical_disk_operations",
		Type:     module.Area,
		Priority: prioDiskOperations,
		Dims: module.Dims{
			{ID: "logical_disk_%s_reads_total", Name: "reads", Algo: module.Incremental},
			{ID: "logical_disk_%s_writes_total", Name: "writes", Algo: module.Incremental, Mul: -1},
		},
	}
	diskAvgLatencyChart = module.Chart{
		ID:       "logical_disk_%s_latency",
		Title:    "Average Read/Write Latency Disk %s",
		Units:    "milliseconds",
		Fam:      "disk",
		Ctx:      "wmi.logical_disk_latency",
		Priority: prioDiskAvgLatency,
		Dims: module.Dims{
			{ID: "logical_disk_%s_read_latency", Name: "read", Algo: module.Incremental},
			{ID: "logical_disk_%s_write_latency", Name: "write", Algo: module.Incremental},
		},
	}
)

func newOSCharts() module.Charts {
	return module.Charts{
		osProcessesChart.Copy(),
		osUsersChart.Copy(),
		osMemoryUsage.Copy(),
		osPagingFilesUsageChart.Copy(),
	}
}

var (
	osProcessesChart = module.Chart{
		ID:       "os_processes",
		Title:    "Processes",
		Units:    "number",
		Fam:      "os",
		Ctx:      "wmi.os_processes",
		Priority: prioOSProcesses,
		Dims: module.Dims{
			{ID: "os_processes", Name: "processes"},
		},
		Vars: module.Vars{
			{ID: "os_processes_limit"},
		},
	}
	osUsersChart = module.Chart{
		ID:       "os_users",
		Title:    "Number of Users",
		Units:    "users",
		Fam:      "os",
		Ctx:      "wmi.os_users",
		Priority: prioOSUsers,
		Dims: module.Dims{
			{ID: "os_users", Name: "users"},
		},
	}
	osMemoryUsage = module.Chart{
		ID:       "os_visible_memory_usage",
		Title:    "Visible Memory Usage",
		Units:    "bytes",
		Fam:      "os",
		Ctx:      "wmi.os_visible_memory_usage",
		Type:     module.Stacked,
		Priority: prioOSVisibleMemoryUsage,
		Dims: module.Dims{
			{ID: "os_physical_memory_free_bytes", Name: "free", Div: 1000},
			{ID: "os_visible_memory_used_bytes", Name: "used", Div: 1000},
		},
		Vars: module.Vars{
			{ID: "os_visible_memory_bytes"},
		},
	}
	osPagingFilesUsageChart = module.Chart{
		ID:       "os_paging_files_usage",
		Title:    "Paging Files Usage",
		Units:    "bytes",
		Fam:      "os",
		Ctx:      "wmi.os_paging_files_usage",
		Type:     module.Stacked,
		Priority: prioOSPagingUsage,
		Dims: module.Dims{
			{ID: "os_paging_free_bytes", Name: "free", Div: 1000},
			{ID: "os_paging_used_bytes", Name: "used", Div: 1000},
		},
		Vars: module.Vars{
			{ID: "os_paging_limit_bytes"},
		},
	}
)

func newSystemCharts() module.Charts {
	return module.Charts{
		systemThreadsChart.Copy(),
		systemUptimeChart.Copy(),
	}
}

var (
	systemThreadsChart = module.Chart{
		ID:       "system_threads",
		Title:    "Threads",
		Units:    "number",
		Fam:      "system",
		Ctx:      "wmi.system_threads",
		Priority: prioSystemThreads,
		Dims: module.Dims{
			{ID: "system_threads", Name: "threads"},
		},
	}
	systemUptimeChart = module.Chart{
		ID:       "system_uptime",
		Title:    "Uptime",
		Units:    "seconds",
		Fam:      "system",
		Ctx:      "wmi.system_uptime",
		Priority: prioSystemUptime,
		Dims: module.Dims{
			{ID: "system_up_time", Name: "time"},
		},
	}
)

func newLogonCharts() module.Charts {
	return module.Charts{
		logonSessionsChart.Copy(),
	}
}

var (
	logonSessionsChart = module.Chart{
		ID:       "logon_active_sessions_by_type",
		Title:    "Active User Logon Sessions By Type",
		Units:    "sessions",
		Fam:      "logon",
		Ctx:      "wmi.logon_type_sessions",
		Type:     module.Stacked,
		Priority: prioLogonSessions,
		Dims: module.Dims{
			{ID: "logon_type_system", Name: "system"},
			{ID: "logon_type_interactive", Name: "interactive"},
			{ID: "logon_type_network", Name: "network"},
			{ID: "logon_type_batch", Name: "batch"},
			{ID: "logon_type_service", Name: "service"},
			{ID: "logon_type_proxy", Name: "proxy"},
			{ID: "logon_type_unlock", Name: "unlock"},
			{ID: "logon_type_network_clear_text", Name: "network_clear_text"},
			{ID: "logon_type_new_credentials", Name: "new_credentials"},
			{ID: "logon_type_remote_interactive", Name: "remote_interactive"},
			{ID: "logon_type_cached_interactive", Name: "cached_interactive"},
			{ID: "logon_type_cached_remote_interactive", Name: "cached_remote_interactive"},
			{ID: "logon_type_cached_unlock", Name: "cached_unlock"},
		},
	}
)

func newThermalzoneCharts() module.Charts {
	return module.Charts{
		thermalzoneTemperatureChart.Copy(),
	}
}

var (
	thermalzoneTemperatureChart = module.Chart{
		ID:       "thermalzone_temperature",
		Title:    "Thermal zone temperature",
		Units:    "celsius",
		Fam:      "thermalzone",
		Ctx:      "wmi.thermalzone_temperature",
		Type:     module.Area,
		Priority: prioThermalzoneTemperature,
	}
)

func newCollectionCharts() *module.Charts {
	return &module.Charts{
		collectionDurationChart.Copy(),
		collectionsStatusChart.Copy(),
	}
}

var (
	collectionDurationChart = module.Chart{
		ID:       "collector_duration",
		Title:    "Duration",
		Units:    "ms",
		Fam:      "collection",
		Ctx:      "cpu.collector_duration",
		Priority: prioCollectionDuration,
	}
	collectionsStatusChart = module.Chart{
		ID:       "collector_success",
		Title:    "Collection Status",
		Units:    "bool",
		Fam:      "collection",
		Ctx:      "cpu.collector_success",
		Priority: prioCollectionStatus,
	}
)

func newServicesCharts() *module.Charts {
	return &module.Charts{
		servicesStateContinuePendingChart.Copy(),
		servicesStatePausePendingChart.Copy(),
		servicesStatePausedChart.Copy(),
		servicesStateRunningChart.Copy(),
		servicesStateStartPendingChart.Copy(),
		servicesStateStopPendingChart.Copy(),
		servicesStateStoppedChart.Copy(),
		servicesStateUnknownChart.Copy(),
		servicesStatusDegradedChart.Copy(),
		servicesStatusErrorChart.Copy(),
		servicesStatusLostCommChart.Copy(),
		servicesStatusNoContactChart.Copy(),
		servicesStatusNonRecoverChart.Copy(),
		servicesStatusOKChart.Copy(),
		servicesStatusPredFailChart.Copy(),
		servicesStatusServiceChart.Copy(),
		servicesStatusStartingChart.Copy(),
		servicesStatusStoppingChart.Copy(),
		servicesStatusStressedChart.Copy(),
		servicesStatusUnknownChart.Copy(),
	}
}

var (
	// State
	servicesStateContinuePendingChart = module.Chart{
		ID:       "services_state_continue_pending",
		Title:    "Services with Continue Pending state",
		Units:    "bool",
		Fam:      "services",
		Ctx:      "wmi.services_state_continue_pending",
		Priority: prioServiceStateContinuePending,
	}
	servicesStatePausePendingChart = module.Chart{
		ID:       "services_state_pause_pending",
		Title:    "Services with Pause Pending state",
		Units:    "bool",
		Fam:      "services",
		Ctx:      "wmi.services_state_pause_pending",
		Priority: prioServiceStatePausePending,
	}
	servicesStatePausedChart = module.Chart{
		ID:       "services_state_paused",
		Title:    "Services with Paused state",
		Units:    "bool",
		Fam:      "services",
		Ctx:      "wmi.services_state_paused",
		Priority: prioServiceStatePaused,
	}
	servicesStateRunningChart = module.Chart{
		ID:       "services_state_running",
		Title:    "Services with Running state",
		Units:    "bool",
		Fam:      "services",
		Ctx:      "wmi.services_state_running",
		Priority: prioServiceStateRunning,
	}
	servicesStateStartPendingChart = module.Chart{
		ID:       "services_state_start_pending",
		Title:    "Services with Start Pending state",
		Units:    "bool",
		Fam:      "services",
		Ctx:      "wmi.services_state_start_pending",
		Priority: prioServiceStateStartPending,
	}
	servicesStateStopPendingChart = module.Chart{
		ID:       "services_state_stop_pending",
		Title:    "Services with Stop Pending state",
		Units:    "bool",
		Fam:      "services",
		Ctx:      "wmi.services_state_stop_pending",
		Priority: prioServiceStateStopPending,
	}
	servicesStateStoppedChart = module.Chart{
		ID:       "services_state_stopped",
		Title:    "Services with Stopped state",
		Units:    "bool",
		Fam:      "services",
		Ctx:      "wmi.services_state_stopped",
		Priority: prioServiceStateStopPending,
	}
	servicesStateUnknownChart = module.Chart{
		ID:       "services_state_unkown",
		Title:    "Services with Unknown state",
		Units:    "bool",
		Fam:      "services",
		Ctx:      "wmi.services_state_unknown",
		Priority: prioServiceStateUnknown,
	}
	// Status
	servicesStatusDegradedChart = module.Chart{
		ID:       "services_state_status_degraded",
		Title:    "Services with status Degraded",
		Units:    "bool",
		Fam:      "services",
		Ctx:      "wmi.services_status_degraded",
		Priority: prioServiceStatusDegraded,
	}
	servicesStatusErrorChart = module.Chart{
		ID:       "services_state_status_error",
		Title:    "Services with status Error",
		Units:    "bool",
		Fam:      "services",
		Ctx:      "wmi.services_status_error",
		Priority: prioServiceStatusError,
	}
	servicesStatusLostCommChart = module.Chart{
		ID:       "services_state_status_lost_comm",
		Title:    "Services with status Lost Comm",
		Units:    "bool",
		Fam:      "services",
		Ctx:      "wmi.services_status_lost_comm",
		Priority: prioServiceStatusLostConn,
	}
	servicesStatusNoContactChart = module.Chart{
		ID:       "services_state_status_no_contact",
		Title:    "Services with status No Contact",
		Units:    "bool",
		Fam:      "services",
		Ctx:      "wmi.services_status_no_contact",
		Priority: prioServiceStatusNoContact,
	}
	servicesStatusNonRecoverChart = module.Chart{
		ID:       "services_state_status_non_recover",
		Title:    "Services with status Non Recover",
		Units:    "bool",
		Fam:      "services",
		Ctx:      "wmi.services_status_non_recover",
		Priority: prioServiceStatusNonRecover,
	}
	servicesStatusOKChart = module.Chart{
		ID:       "services_state_status_ok",
		Title:    "Services with status OK",
		Units:    "bool",
		Fam:      "services",
		Ctx:      "wmi.services_status_ok",
		Priority: prioServiceStatusOK,
	}
	servicesStatusPredFailChart = module.Chart{
		ID:       "services_state_status_pred_fail",
		Title:    "Services with status Pred Fail",
		Units:    "bool",
		Fam:      "services",
		Ctx:      "wmi.services_status_pred_fail",
		Priority: prioServiceStatusPredFail,
	}
	servicesStatusServiceChart = module.Chart{
		ID:       "services_state_status_service",
		Title:    "Services with status Service",
		Units:    "bool",
		Fam:      "services",
		Ctx:      "wmi.services_status_service",
		Priority: prioServiceStatusPredFail,
	}
	servicesStatusStartingChart = module.Chart{
		ID:       "services_state_status_starting",
		Title:    "Services with status Starting",
		Units:    "bool",
		Fam:      "services",
		Ctx:      "wmi.services_status_starting",
		Priority: prioServiceStatusStarting,
	}
	servicesStatusStoppingChart = module.Chart{
		ID:       "services_state_status_stopping",
		Title:    "Services with status Stopping",
		Units:    "bool",
		Fam:      "services",
		Ctx:      "wmi.services_status_stopping",
		Priority: prioServiceStatusStopping,
	}
	servicesStatusStressedChart = module.Chart{
		ID:       "services_state_status_stressed",
		Title:    "Services with status Stressed",
		Units:    "bool",
		Fam:      "services",
		Ctx:      "wmi.services_status_stressed",
		Priority: prioServiceStatusStressed,
	}
	servicesStatusUnknownChart = module.Chart{
		ID:       "services_state_status_unknown",
		Title:    "Services with status Unknown",
		Units:    "bool",
		Fam:      "services",
		Ctx:      "wmi.services_status_unknown",
		Priority: prioServiceStatusUnknown,
	}
)

func newChartFromTemplate(template module.Chart, id string) *module.Chart {
	chart := template.Copy()
	chart.ID = fmt.Sprintf(chart.ID, id)
	chart.Title = fmt.Sprintf(chart.Title, id)
	for _, dim := range chart.Dims {
		dim.ID = fmt.Sprintf(dim.ID, id)
	}
	for _, v := range chart.Vars {
		v.ID = fmt.Sprintf(v.ID, id)
	}
	return chart
}

func (w *WMI) updateCharts(mx *metrics) {
	w.updateCollectionCharts(mx)
	w.updateCPUCharts(mx)
	w.updateMemoryCharts(mx)
	w.updateNetCharts(mx)
	w.updateLogicalDisksCharts(mx)
	w.updateSystemCharts(mx)
	w.updateOSCharts(mx)
	w.updateLogonCharts(mx)
	w.updateThermalzoneCharts(mx)
	w.updateTCPCharts(mx)
	w.updateProcessesCharts(mx)
	w.updateServicesCharts(mx)
}

func (w *WMI) updateServicesCharts(mx *metrics) {
	if !mx.hasServices() {
		return
	}

	if !w.cache.collectors[collectorService] {
		w.cache.collectors[collectorService] = true

		if err := w.Charts().Add(*newServicesCharts()...); err != nil {
			w.Warning(err)
		}
	}

	servs := make([]string, 0, len(mx.Services.servs))
	for _, serv := range mx.Services.servs {
		servs = append(servs, serv.ID)
	}
	sort.Slice(servs, func(i, j int) bool { return servs[i] < servs[j] })

	for _, id := range servs {
		if w.cache.servs[id] {
			continue
		}
		w.cache.servs[id] = true

		// State
		if err := addDimToServicesStateContinuePendingChart(w.Charts(), id); err != nil {
			w.Warning(err)
		}
		if err := addDimToServicesStatePausePendingChart(w.Charts(), id); err != nil {
			w.Warning(err)
		}
		if err := addDimToServicesStatePausedChart(w.Charts(), id); err != nil {
			w.Warning(err)
		}
		if err := addDimToServicesStateRunningChart(w.Charts(), id); err != nil {
			w.Warning(err)
		}
		if err := addDimToServicesStartPendingChart(w.Charts(), id); err != nil {
			w.Warning(err)
		}
		if err := addDimToServicesStopPendingChart(w.Charts(), id); err != nil {
			w.Warning(err)
		}
		if err := addDimToServicesStoppedChart(w.Charts(), id); err != nil {
			w.Warning(err)
		}
		if err := addDimToServicesUnknownChart(w.Charts(), id); err != nil {
			w.Warning(err)
		}
		// Status
		if err := addDimToServicesStatusDegradedChart(w.Charts(), id); err != nil {
			w.Warning(err)
		}
		if err := addDimToServicesStatusErrorChart(w.Charts(), id); err != nil {
			w.Warning(err)
		}
		if err := addDimToServicesStatusLostCommChart(w.Charts(), id); err != nil {
			w.Warning(err)
		}
		if err := addDimToServicesStatusNoContactChart(w.Charts(), id); err != nil {
			w.Warning(err)
		}
		if err := addDimToServicesStatusNonRecoverChart(w.Charts(), id); err != nil {
			w.Warning(err)
		}
		if err := addDimToServicesStatusOKChart(w.Charts(), id); err != nil {
			w.Warning(err)
		}
		if err := addDimToServicesStatusPredFailChart(w.Charts(), id); err != nil {
			w.Warning(err)
		}
		if err := addDimToServicesStatusServiceChart(w.Charts(), id); err != nil {
			w.Warning(err)
		}
		if err := addDimToServicesStatusStartingChart(w.Charts(), id); err != nil {
			w.Warning(err)
		}
		if err := addDimToServicesStatusStoppingChart(w.Charts(), id); err != nil {
			w.Warning(err)
		}
		if err := addDimToServicesStatusStressedChart(w.Charts(), id); err != nil {
			w.Warning(err)
		}
		if err := addDimToServicesStatusUnknownChart(w.Charts(), id); err != nil {
			w.Warning(err)
		}
	}
}

func (w *WMI) updateProcessesCharts(mx *metrics) {
	if !mx.hasProcesses() {
		return
	}

	if !w.cache.collectors[collectorProcess] {
		w.cache.collectors[collectorProcess] = true

		if err := w.Charts().Add(newProcessesCharts()...); err != nil {
			w.Warning(err)
		}
	}

	procs := make([]string, 0, len(mx.Processes.procs))
	for _, proc := range mx.Processes.procs {
		procs = append(procs, proc.ID)
	}
	sort.Slice(procs, func(i, j int) bool { return procs[i] < procs[j] })

	for _, id := range procs {
		if w.cache.procs[id] {
			continue
		}
		w.cache.procs[id] = true

		if err := addDimToProcessesCPUTimeTotalChart(w.Charts(), id); err != nil {
			w.Warning(err)
		}
		if err := addDimToProcessesHandlesChart(w.Charts(), id); err != nil {
			w.Warning(err)
		}
		if err := addDimToProcessesIOBytesChart(w.Charts(), id); err != nil {
			w.Warning(err)
		}
		if err := addDimToProcessesIOOperationsChart(w.Charts(), id); err != nil {
			w.Warning(err)
		}
		if err := addDimToProcessesPageFaultsChart(w.Charts(), id); err != nil {
			w.Warning(err)
		}
		if err := addDimToProcessesPageFileBytes(w.Charts(), id); err != nil {
			w.Warning(err)
		}
		if err := addDimToProcessesPoolBytes(w.Charts(), id); err != nil {
			w.Warning(err)
		}
		if err := addDimToProcessesThreads(w.Charts(), id); err != nil {
			w.Warning(err)
		}
	}

}

func (w *WMI) updateCollectionCharts(mx *metrics) {
	if !mx.hasCollectors() {
		return
	}
	for _, c := range *mx.Collectors {
		if w.cache.collection[c.ID] {
			continue
		}
		w.cache.collection[c.ID] = true
		if err := addDimToCollectionDurationChart(w.Charts(), c.ID); err != nil {
			w.Warning(err)
		}
		if err := addDimToCollectionStatusChart(w.Charts(), c.ID); err != nil {
			w.Warning(err)
		}
	}
}

func (w *WMI) updateCPUCharts(mx *metrics) {
	if !mx.hasCPU() {
		return
	}
	if !w.cache.collectors[collectorCPU] {
		w.cache.collectors[collectorCPU] = true
		if err := w.Charts().Add(newCPUCharts()...); err != nil {
			w.Warning(err)
		}
	}
	for _, core := range mx.CPU.Cores {
		if w.cache.cores[core.ID] {
			continue
		}
		w.cache.cores[core.ID] = true
		if err := addCPUCoreCharts(w.Charts(), core.ID); err != nil {
			w.Warning(err)
		}
		if err := addDimToCPUDPCsChart(w.Charts(), core.ID); err != nil {
			w.Warning(err)
		}
		if err := addDimToCPUInterruptsChart(w.Charts(), core.ID); err != nil {
			w.Warning(err)
		}
	}
}

func (w *WMI) updateMemoryCharts(mx *metrics) {
	if !mx.hasMemory() {
		return
	}
	if w.cache.collectors[collectorMemory] {
		return
	}
	w.cache.collectors[collectorMemory] = true
	if err := w.Charts().Add(newMemCharts()...); err != nil {
		w.Warning(err)
	}
}

func (w *WMI) updateNetCharts(mx *metrics) {
	if !mx.hasNet() {
		return
	}
	for _, nic := range mx.Net.NICs {
		if w.cache.nics[nic.ID] {
			continue
		}
		w.cache.nics[nic.ID] = true
		if err := addNICCharts(w.Charts(), nic.ID); err != nil {
			w.Warning(err)
		}
	}
}

func (w *WMI) updateLogicalDisksCharts(mx *metrics) {
	if !mx.hasLogicalDisk() {
		return
	}
	set := make(map[string]bool)
	for _, vol := range mx.LogicalDisk.Volumes {
		set[vol.ID] = true
		if w.cache.volumes[vol.ID] {
			continue
		}
		w.cache.volumes[vol.ID] = true
		if err := addLogicalDiskCharts(w.Charts(), vol.ID); err != nil {
			w.Warning(err)
		}
	}
	for vol := range w.cache.volumes {
		if set[vol] {
			continue
		}
		delete(w.cache.volumes, vol)
		removeLogicalDiskFromCharts(w.Charts(), vol)
	}
}

func (w *WMI) updateSystemCharts(mx *metrics) {
	if !mx.hasSystem() {
		return
	}
	if w.cache.collectors[collectorSystem] {
		return
	}
	w.cache.collectors[collectorSystem] = true
	if err := w.Charts().Add(newSystemCharts()...); err != nil {
		w.Warning(err)
	}
}

func (w *WMI) updateOSCharts(mx *metrics) {
	if !mx.hasOS() {
		return
	}
	if w.cache.collectors[collectorOS] {
		return
	}
	w.cache.collectors[collectorOS] = true
	if err := w.Charts().Add(newOSCharts()...); err != nil {
		w.Warning(err)
	}
}

func (w *WMI) updateLogonCharts(mx *metrics) {
	if !mx.hasLogon() {
		return
	}
	if w.cache.collectors[collectorLogon] {
		return
	}
	w.cache.collectors[collectorLogon] = true
	if err := w.Charts().Add(newLogonCharts()...); err != nil {
		w.Warning(err)
	}
}

func (w *WMI) updateThermalzoneCharts(mx *metrics) {
	if !mx.hasThermalZone() {
		return
	}

	if !w.cache.collectors[collectorThermalzone] {
		w.cache.collectors[collectorThermalzone] = true
		if err := w.Charts().Add(newThermalzoneCharts()...); err != nil {
			w.Warning(err)
		}
	}

	for _, zone := range mx.ThermalZone.Zones {
		if w.cache.thermalZones[zone.ID] {
			continue
		}
		w.cache.thermalZones[zone.ID] = true
		if err := addDimToThermalzoneTemperatureChart(w.Charts(), zone.ID); err != nil {
			w.Warning(err)
		}
	}

}

func (w *WMI) updateTCPCharts(mx *metrics) {
	if !mx.hasTCP() {
		return
	}
	if w.cache.collectors[collectorTCP] {
		return
	}
	w.cache.collectors[collectorTCP] = true
	if err := w.Charts().Add(newTCPCharts()...); err != nil {
		w.Warning(err)
	}
}

func addCPUCoreCharts(charts *module.Charts, coreID string) error {
	for _, chart := range newCPUCoreCharts() {
		chart = newChartFromTemplate(*chart, coreID)
		if err := charts.Add(chart); err != nil {
			return err
		}
	}
	return nil
}

func addNICCharts(charts *module.Charts, nicID string) error {
	for _, chart := range newNICCharts() {
		chart = newChartFromTemplate(*chart, nicID)
		if err := charts.Add(chart); err != nil {
			return err
		}
	}
	return nil
}

func addLogicalDiskCharts(charts *module.Charts, diskID string) error {
	for _, chart := range newDiskCharts() {
		chart = newChartFromTemplate(*chart, diskID)
		if err := charts.Add(chart); err != nil {
			return err
		}
	}
	return nil
}

func removeLogicalDiskFromCharts(charts *module.Charts, diskID string) {
	for _, chart := range *charts {
		if !strings.HasPrefix(chart.ID, fmt.Sprintf("logical_disk_%s", diskID)) {
			continue
		}
		chart.MarkRemove()
		chart.MarkNotCreated()
	}
}

func addDimToCPUDPCsChart(charts *module.Charts, coreID string) error {
	chart := charts.Get(cpuDPCsChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", cpuDPCsChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("cpu_core_%s_dpc", coreID),
		Name: "core" + coreID,
		Algo: module.Incremental,
		Div:  1000,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToCPUInterruptsChart(charts *module.Charts, coreID string) error {
	chart := charts.Get(cpuInterruptsChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", cpuInterruptsChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("cpu_core_%s_interrupts", coreID),
		Name: "core" + coreID,
		Algo: module.Incremental,
		Div:  1000,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToThermalzoneTemperatureChart(charts *module.Charts, zoneName string) error {
	chart := charts.Get(thermalzoneTemperatureChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", thermalzoneTemperatureChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("thermalzone_%s_temperature", zoneName),
		Name: zoneName,
		Algo: module.Absolute,
		Div:  1000,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToCollectionDurationChart(charts *module.Charts, colName string) error {
	chart := charts.Get(collectionDurationChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", collectionDurationChart.ID)
	}
	dim := &module.Dim{
		ID:   colName + "_collection_duration",
		Name: colName,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToCollectionStatusChart(charts *module.Charts, colName string) error {
	chart := charts.Get(collectionsStatusChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", collectionsStatusChart.ID)
	}
	dim := &module.Dim{
		ID:   colName + "_collection_success",
		Name: colName,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToProcessesCPUTimeTotalChart(charts *module.Charts, procID string) error {
	chart := charts.Get(processesCPUTimeTotalChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", processesCPUTimeTotalChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("process_%s_cpu_time", procID),
		Name: procID,
		Algo: module.Incremental,
		Div:  1000,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToProcessesHandlesChart(charts *module.Charts, procID string) error {
	chart := charts.Get(processesHandlesChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", processesHandlesChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("process_%s_handles", procID),
		Name: procID,
		Algo: module.Absolute,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToProcessesIOBytesChart(charts *module.Charts, procID string) error {
	chart := charts.Get(processesIOBytesChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", processesIOBytesChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("process_%s_io_bytes", procID),
		Name: procID,
		Algo: module.Incremental,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToProcessesIOOperationsChart(charts *module.Charts, procID string) error {
	chart := charts.Get(processesIOOperationsChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", processesIOOperationsChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("process_%s_io_operations", procID),
		Name: procID,
		Algo: module.Incremental,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToProcessesPageFaultsChart(charts *module.Charts, procID string) error {
	chart := charts.Get(processesPageFaultsChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", processesPageFaultsChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("process_%s_page_faults", procID),
		Name: procID,
		Algo: module.Incremental,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToProcessesPageFileBytes(charts *module.Charts, procID string) error {
	chart := charts.Get(processesPageFileBytes.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", processesPageFileBytes.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("process_%s_page_file_bytes", procID),
		Name: procID,
		Algo: module.Absolute,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToProcessesPoolBytes(charts *module.Charts, procID string) error {
	chart := charts.Get(processesPoolBytes.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", processesPoolBytes.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("process_%s_pool_bytes", procID),
		Name: procID,
		Algo: module.Absolute,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToProcessesThreads(charts *module.Charts, procID string) error {
	chart := charts.Get(processesThreads.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", processesThreads.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("process_%s_threads", procID),
		Name: procID,
		Algo: module.Absolute,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToServicesStateContinuePendingChart(charts *module.Charts, servID string) error {
	chart := charts.Get(servicesStateContinuePendingChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", servicesStateContinuePendingChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("service_%s_state_continue_pending", servID),
		Name: servID,
		Algo: module.Absolute,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToServicesStatePausePendingChart(charts *module.Charts, servID string) error {
	chart := charts.Get(servicesStateContinuePendingChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", servicesStatePausePendingChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("service_%s_state_pause_pending", servID),
		Name: servID,
		Algo: module.Absolute,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToServicesStatePausedChart(charts *module.Charts, servID string) error {
	chart := charts.Get(servicesStatePausedChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", servicesStatePausedChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("service_%s_state_paused", servID),
		Name: servID,
		Algo: module.Absolute,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToServicesStateRunningChart(charts *module.Charts, servID string) error {
	chart := charts.Get(servicesStateRunningChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", servicesStateRunningChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("service_%s_state_running", servID),
		Name: servID,
		Algo: module.Absolute,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToServicesStartPendingChart(charts *module.Charts, servID string) error {
	chart := charts.Get(servicesStateStartPendingChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", servicesStateStartPendingChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("service_%s_state_start_pending", servID),
		Name: servID,
		Algo: module.Absolute,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToServicesStopPendingChart(charts *module.Charts, servID string) error {
	chart := charts.Get(servicesStateStopPendingChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", servicesStateStopPendingChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("service_%s_state_stop_pending", servID),
		Name: servID,
		Algo: module.Absolute,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToServicesStoppedChart(charts *module.Charts, servID string) error {
	chart := charts.Get(servicesStateStoppedChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", servicesStateStoppedChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("service_%s_state_stopped", servID),
		Name: servID,
		Algo: module.Absolute,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToServicesUnknownChart(charts *module.Charts, servID string) error {
	chart := charts.Get(servicesStateUnknownChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", servicesStateUnknownChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("service_%s_state_unknown", servID),
		Name: servID,
		Algo: module.Absolute,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToServicesStatusDegradedChart(charts *module.Charts, servID string) error {
	chart := charts.Get(servicesStatusDegradedChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", servicesStatusDegradedChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("service_%s_status_degraded", servID),
		Name: servID,
		Algo: module.Absolute,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToServicesStatusErrorChart(charts *module.Charts, servID string) error {
	chart := charts.Get(servicesStatusErrorChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", servicesStatusErrorChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("service_%s_status_error", servID),
		Name: servID,
		Algo: module.Absolute,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToServicesStatusLostCommChart(charts *module.Charts, servID string) error {
	chart := charts.Get(servicesStatusLostCommChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", servicesStatusLostCommChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("service_%s_status_lost_comm", servID),
		Name: servID,
		Algo: module.Absolute,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToServicesStatusNoContactChart(charts *module.Charts, servID string) error {
	chart := charts.Get(servicesStatusNoContactChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", servicesStatusNoContactChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("service_%s_status_no_contact", servID),
		Name: servID,
		Algo: module.Absolute,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToServicesStatusNonRecoverChart(charts *module.Charts, servID string) error {
	chart := charts.Get(servicesStatusNonRecoverChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", servicesStatusNonRecoverChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("service_%s_status_nonrecover", servID),
		Name: servID,
		Algo: module.Absolute,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToServicesStatusOKChart(charts *module.Charts, servID string) error {
	chart := charts.Get(servicesStatusOKChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", servicesStatusOKChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("service_%s_status_ok", servID),
		Name: servID,
		Algo: module.Absolute,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToServicesStatusPredFailChart(charts *module.Charts, servID string) error {
	chart := charts.Get(servicesStatusPredFailChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", servicesStatusPredFailChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("service_%s_status_pred_fail", servID),
		Name: servID,
		Algo: module.Absolute,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToServicesStatusServiceChart(charts *module.Charts, servID string) error {
	chart := charts.Get(servicesStatusServiceChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", servicesStatusServiceChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("service_%s_status_service", servID),
		Name: servID,
		Algo: module.Absolute,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToServicesStatusStartingChart(charts *module.Charts, servID string) error {
	chart := charts.Get(servicesStatusStartingChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", servicesStatusStartingChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("service_%s_status_starting", servID),
		Name: servID,
		Algo: module.Absolute,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToServicesStatusStoppingChart(charts *module.Charts, servID string) error {
	chart := charts.Get(servicesStatusStoppingChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", servicesStatusStoppingChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("service_%s_status_stopping", servID),
		Name: servID,
		Algo: module.Absolute,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToServicesStatusStressedChart(charts *module.Charts, servID string) error {
	chart := charts.Get(servicesStatusStressedChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", servicesStatusStressedChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("service_%s_status_stressed", servID),
		Name: servID,
		Algo: module.Absolute,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToServicesStatusUnknownChart(charts *module.Charts, servID string) error {
	chart := charts.Get(servicesStatusUnknownChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", servicesStatusUnknownChart.ID)
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("service_%s_status_unknown", servID),
		Name: servID,
		Algo: module.Absolute,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

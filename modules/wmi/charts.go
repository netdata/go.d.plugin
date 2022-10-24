// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import (
	"fmt"
	"strings"

	"github.com/netdata/go.d.plugin/agent/module"
)

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Chart is an alias for module.Chart
	Chart = module.Chart
	// Dims is an alias for module.Dims
	Dims = module.Dims
	// Dim is an alias for module.Dim
	Dim = module.Dim
	// Vars is an alias for module.Vars
	Vars = module.Vars
	// Opts is an alias for module.Dim
	Opts = module.DimOpts
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
)

func newCPUCharts() Charts {
	return Charts{
		cpuUtilChart.Copy(),
		cpuDPCsChart.Copy(),
		cpuInterruptsChart.Copy(),
	}
}

var (
	cpuUtilChart = Chart{
		ID:       "cpu_utilization_total",
		Title:    "Total CPU Utilization (all cores)",
		Units:    "percentage",
		Fam:      "cpu",
		Ctx:      "wmi.cpu_utilization_total",
		Type:     module.Stacked,
		Priority: prioCPUUtil,
		Dims: Dims{
			{ID: "cpu_idle", Name: "idle", Algo: module.PercentOfIncremental, Div: 1000, DimOpts: Opts{Hidden: true}},
			{ID: "cpu_dpc", Name: "dpc", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_user", Name: "user", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_privileged", Name: "privileged", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_interrupt", Name: "interrupt", Algo: module.PercentOfIncremental, Div: 1000},
		},
	}
	cpuDPCsChart = Chart{
		ID:       "cpu_dpcs",
		Title:    "Received and Serviced Deferred Procedure Calls (DPC)",
		Units:    "dpc/s",
		Fam:      "cpu",
		Ctx:      "wmi.cpu_dpcs",
		Type:     module.Stacked,
		Priority: prioCPUDPCs,
	}
	cpuInterruptsChart = Chart{
		ID:       "cpu_interrupts",
		Title:    "Received and Serviced Hardware Interrupts",
		Units:    "interrupts/s",
		Fam:      "cpu",
		Ctx:      "wmi.cpu_interrupts",
		Type:     module.Stacked,
		Priority: prioCPUInterrupts,
	}
)

func newCPUCoreCharts() Charts {
	return Charts{
		cpuCoreUtilChart.Copy(),
		cpuCoreCStateChart.Copy(),
	}
}

var (
	cpuCoreUtilChart = Chart{
		ID:       "core_%s_cpu_utilization",
		Title:    "Core%s CPU Utilization",
		Units:    "percentage",
		Fam:      "cpu",
		Ctx:      "wmi.cpu_utilization",
		Type:     module.Stacked,
		Priority: prioCPUCoreUtil,
		Dims: Dims{
			{ID: "cpu_core_%s_idle", Name: "idle", Algo: module.PercentOfIncremental, Div: 1000, DimOpts: Opts{Hidden: true}},
			{ID: "cpu_core_%s_dpc", Name: "dpc", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_core_%s_user", Name: "user", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_core_%s_privileged", Name: "privileged", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_core_%s_interrupt", Name: "interrupt", Algo: module.PercentOfIncremental, Div: 1000},
		},
	}
	cpuCoreCStateChart = Chart{
		ID:       "core_%s_cpu_cstate",
		Title:    "Core%s Time Spent in Low-Power Idle State",
		Units:    "percentage",
		Fam:      "cpu",
		Ctx:      "wmi.cpu_cstate",
		Type:     module.Stacked,
		Priority: prioCPUCoreCState,
		Dims: Dims{
			{ID: "cpu_core_%s_c1", Name: "c1", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_core_%s_c2", Name: "c2", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_core_%s_c3", Name: "c3", Algo: module.PercentOfIncremental, Div: 1000},
		},
	}
)

func newMemCharts() Charts {
	return Charts{
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
	memUtilChart = Chart{
		ID:       "memory_utilization",
		Title:    "Memory Utilization",
		Units:    "KiB",
		Fam:      "mem",
		Ctx:      "wmi.memory_utilization",
		Type:     module.Stacked,
		Priority: prioMemUtil,
		Dims: Dims{
			{ID: "memory_available_bytes", Name: "available", Div: 1000 * 1024},
			{ID: "memory_used_bytes", Name: "used", Div: 1000 * 1024},
		},
		Vars: Vars{
			{ID: "os_visible_memory_bytes"},
		},
	}
	memPageFaultsChart = Chart{
		ID:       "memory_page_faults",
		Title:    "Memory Page Faults",
		Units:    "events/s",
		Fam:      "mem",
		Ctx:      "wmi.memory_page_faults",
		Priority: prioMemPageFaults,
		Dims: Dims{
			{ID: "memory_page_faults_total", Name: "page faults", Algo: module.Incremental, Div: 1000},
		},
	}
	memSwapUtilChart = Chart{
		ID:       "memory_swap_utilization",
		Title:    "Swap Utilization",
		Units:    "KiB",
		Fam:      "mem",
		Ctx:      "wmi.memory_swap_utilization",
		Type:     module.Stacked,
		Priority: prioMemSwapUtil,
		Dims: Dims{
			{ID: "memory_not_committed_bytes", Name: "available", Div: 1000 * 1024},
			{ID: "memory_committed_bytes", Name: "used", Div: 1000 * 1024},
		},
		Vars: Vars{
			{ID: "memory_commit_limit"},
		},
	}
	memSwapOperationsChart = Chart{
		ID:       "memory_swap_operations",
		Title:    "Swap Operations",
		Units:    "operations/s",
		Fam:      "mem",
		Ctx:      "wmi.memory_swap_operations",
		Type:     module.Area,
		Priority: prioMemSwapOperations,
		Dims: Dims{
			{ID: "memory_swap_page_reads_total", Name: "read", Algo: module.Incremental, Div: 1000},
			{ID: "memory_swap_page_writes_total", Name: "write", Algo: module.Incremental, Div: -1000},
		},
	}
	memSwapPagesChart = Chart{
		ID:       "memory_swap_pages",
		Title:    "Swap Pages",
		Units:    "pages/s",
		Fam:      "mem",
		Ctx:      "wmi.memory_swap_pages",
		Type:     module.Area,
		Priority: prioMemSwapPages,
		Dims: Dims{
			{ID: "memory_swap_pages_read_total", Name: "read", Algo: module.Incremental, Div: 1000},
			{ID: "memory_swap_pages_written_total", Name: "written", Algo: module.Incremental, Div: -1000},
		},
	}
	memCacheChart = Chart{
		ID:       "memory_cached",
		Title:    "Cached",
		Units:    "KiB",
		Fam:      "mem",
		Ctx:      "wmi.memory_cached",
		Priority: prioMemCache,
		Dims: Dims{
			{ID: "memory_cache_total", Name: "cached", Div: 1000 * 1024},
		},
	}
	memCacheFaultsChart = Chart{
		ID:       "memory_cache_faults",
		Title:    "Cache Faults",
		Units:    "events/s",
		Fam:      "mem",
		Ctx:      "wmi.memory_cache_faults",
		Priority: prioMemCacheFaults,
		Dims: Dims{
			{ID: "memory_cache_faults_total", Name: "cache faults", Algo: module.Incremental, Div: 1000},
		},
	}
	memSystemPoolChart = Chart{
		ID:       "memory_system_pool",
		Title:    "System Memory Pool",
		Units:    "KiB",
		Fam:      "mem",
		Ctx:      "wmi.memory_system_pool",
		Type:     module.Stacked,
		Priority: prioMemSystemPool,
		Dims: Dims{
			{ID: "memory_pool_paged_bytes", Name: "paged", Div: 1000 * 1024},
			{ID: "memory_pool_nonpaged_bytes_total", Name: "non-paged", Div: 1000 * 1024},
		},
	}
)

func newNICCharts() Charts {
	return Charts{
		nicBandwidthChart.Copy(),
		nicPacketsChart.Copy(),
		nicErrorsChart.Copy(),
		nicDiscardsChart.Copy(),
	}
}

var (
	nicBandwidthChart = Chart{
		ID:       "nic_%s_bandwidth",
		Title:    "Bandwidth %s",
		Units:    "kilobits/s",
		Fam:      "net",
		Ctx:      "wmi.net_bandwidth",
		Type:     module.Area,
		Priority: prioNICBandwidth,
		Dims: Dims{
			{ID: "net_%s_bytes_received", Name: "received", Algo: module.Incremental, Div: 1000 * 125},
			{ID: "net_%s_bytes_sent", Name: "sent", Algo: module.Incremental, Div: -1000 * 125},
		},
		Vars: Vars{
			{ID: "net_%s_current_bandwidth"},
		},
	}
	nicPacketsChart = Chart{
		ID:       "nic_%s_packets",
		Title:    "Packets %s",
		Units:    "packets/s",
		Fam:      "net",
		Ctx:      "wmi.net_packets",
		Type:     module.Area,
		Priority: prioNICPackets,
		Dims: Dims{
			{ID: "net_%s_packets_received_total", Name: "received", Algo: module.Incremental, Div: 1000},
			{ID: "net_%s_packets_sent_total", Name: "sent", Algo: module.Incremental, Div: -1000},
		},
	}
	nicErrorsChart = Chart{
		ID:       "nic_%s_errors",
		Title:    "Errors %s",
		Units:    "errors/s",
		Fam:      "net",
		Ctx:      "wmi.net_errors",
		Type:     module.Area,
		Priority: prioNICErrors,
		Dims: Dims{
			{ID: "net_%s_packets_received_errors", Name: "inbound", Algo: module.Incremental, Div: 1000},
			{ID: "net_%s_packets_outbound_errors", Name: "outbound", Algo: module.Incremental, Div: -1000},
		},
	}
	nicDiscardsChart = Chart{
		ID:       "nic_%s_discarded",
		Title:    "Discards %s",
		Units:    "discards/s",
		Fam:      "net",
		Ctx:      "wmi.net_discarded",
		Type:     module.Area,
		Priority: prioNICDiscards,
		Dims: Dims{
			{ID: "net_%s_packets_received_discarded", Name: "inbound", Algo: module.Incremental, Div: 1000},
			{ID: "net_%s_packets_outbound_discarded", Name: "outbound", Algo: module.Incremental, Div: -1000},
		},
	}
)

func newTCPCharts() Charts {
	return Charts{
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
	tcpConnsEstablishedChart = Chart{
		ID:       "tcp_conns_established",
		Title:    "TCP established connections",
		Units:    "connections",
		Fam:      "tcp",
		Ctx:      "wmi.tcp_conns_established",
		Priority: prioTCPConnsEstablished,
		Dims: Dims{
			{ID: "tcp_conns_established_ipv4", Name: "ipv4"},
			{ID: "tcp_conns_established_ipv6", Name: "ipv6"},
		},
	}
	tcpConnsActiveChart = Chart{
		ID:       "tcp_conns_active",
		Title:    "TCP active connections",
		Units:    "connections/s",
		Fam:      "tcp",
		Ctx:      "wmi.tcp_conns_active",
		Priority: prioTCPConnsActive,
		Dims: Dims{
			{ID: "tcp_conns_active_ipv4", Name: "ipv4", Algo: module.Incremental},
			{ID: "tcp_conns_active_ipv6", Name: "ipv6", Algo: module.Incremental},
		},
	}
	tcpConnsPassiveChart = Chart{
		ID:       "tcp_conns_passive",
		Title:    "TCP passive connections",
		Units:    "connections/s",
		Fam:      "tcp",
		Ctx:      "wmi.tcp_conns_passive",
		Priority: prioTCPConnsPassive,
		Dims: Dims{
			{ID: "tcp_conns_passive_ipv4", Name: "ipv4", Algo: module.Incremental},
			{ID: "tcp_conns_passive_ipv6", Name: "ipv6", Algo: module.Incremental},
		},
	}
	tcpConnsFailuresChart = Chart{
		ID:       "tcp_conns_failures",
		Title:    "TCP connection failures",
		Units:    "failures/s",
		Fam:      "tcp",
		Ctx:      "wmi.tcp_conns_failures",
		Priority: prioTCPConnsFailure,
		Dims: Dims{
			{ID: "tcp_conns_failures_ipv4", Name: "ipv4", Algo: module.Incremental},
			{ID: "tcp_conns_failures_ipv6", Name: "ipv6", Algo: module.Incremental},
		},
	}
	tcpConnsResetsChart = Chart{
		ID:       "tcp_conns_reset",
		Title:    "TCP connections reseted",
		Units:    "resets/s",
		Fam:      "tcp",
		Ctx:      "wmi.tcp_conns_reset",
		Priority: prioTCPConnsReset,
		Dims: Dims{
			{ID: "tcp_conns_resets_ipv4", Name: "ipv4", Algo: module.Incremental},
			{ID: "tcp_conns_resets_ipv6", Name: "ipv6", Algo: module.Incremental},
		},
	}
	tcpSegmentsReceivedChart = Chart{
		ID:       "tcp_segments_received",
		Title:    "Number of TCP segments received",
		Units:    "segments/s",
		Fam:      "tcp",
		Ctx:      "wmi.tcp_segments_received",
		Priority: prioTCPSegmentsReceived,
		Dims: Dims{
			{ID: "tcp_segments_received_ipv4", Name: "ipv4", Algo: module.Incremental},
			{ID: "tcp_segments_received_ipv6", Name: "ipv6", Algo: module.Incremental},
		},
	}
	tcpSegmentsSentChart = Chart{
		ID:       "tcp_segments_sent",
		Title:    "Number of TCP segments sent",
		Units:    "segments/s",
		Fam:      "tcp",
		Ctx:      "wmi.tcp_segments_sent",
		Priority: prioTCPSegmentsSent,
		Dims: Dims{
			{ID: "tcp_segments_sent_ipv4", Name: "ipv4", Algo: module.Incremental},
			{ID: "tcp_segments_sent_ipv6", Name: "ipv6", Algo: module.Incremental},
		},
	}
	tcpSegmentsRetransmittedChart = Chart{
		ID:       "tcp_segments_retransmitted",
		Title:    "Number of TCP segments retransmitted",
		Units:    "segments/s",
		Fam:      "tcp",
		Ctx:      "wmi.tcp_segments_retransmitted",
		Priority: prioTCPSegmentsRetransmitted,
		Dims: Dims{
			{ID: "tcp_segments_retransmitted_ipv4", Name: "ipv4", Algo: module.Incremental},
			{ID: "tcp_segments_retransmitted_ipv6", Name: "ipv6", Algo: module.Incremental},
		},
	}
)

func newDiskCharts() Charts {
	return Charts{
		diskUtilChart.Copy(),
		diskBandwidthChart.Copy(),
		diskOperationsChart.Copy(),
		diskAvgLatencyChart.Copy(),
	}
}

var (
	diskUtilChart = Chart{
		ID:       "logical_disk_%s_utilization",
		Title:    "Utilization Disk %s",
		Units:    "KiB",
		Fam:      "disk",
		Ctx:      "wmi.logical_disk_utilization",
		Type:     module.Stacked,
		Priority: prioDiskUtil,
		Dims: Dims{
			{ID: "logical_disk_%s_free_space", Name: "free", Div: 1000 * 1024},
			{ID: "logical_disk_%s_used_space", Name: "used", Div: 1000 * 1024},
		},
	}
	diskBandwidthChart = Chart{
		ID:       "logical_disk_%s_bandwidth",
		Title:    "Bandwidth Disk %s",
		Units:    "KiB/s",
		Fam:      "disk",
		Ctx:      "wmi.logical_disk_bandwidth",
		Type:     module.Area,
		Priority: prioDiskBandwidth,
		Dims: Dims{
			{ID: "logical_disk_%s_read_bytes_total", Name: "read", Algo: module.Incremental, Div: 1000 * 1024},
			{ID: "logical_disk_%s_write_bytes_total", Name: "write", Algo: module.Incremental, Div: -1000 * 1024},
		},
	}
	diskOperationsChart = Chart{
		ID:       "logical_disk_%s_operations",
		Title:    "Operations Disk %s",
		Units:    "operations/s",
		Fam:      "disk",
		Ctx:      "wmi.logical_disk_operations",
		Type:     module.Area,
		Priority: prioDiskOperations,
		Dims: Dims{
			{ID: "logical_disk_%s_reads_total", Name: "reads", Algo: module.Incremental},
			{ID: "logical_disk_%s_writes_total", Name: "writes", Algo: module.Incremental, Mul: -1},
		},
	}
	diskAvgLatencyChart = Chart{
		ID:       "logical_disk_%s_latency",
		Title:    "Average Read/Write Latency Disk %s",
		Units:    "milliseconds",
		Fam:      "disk",
		Ctx:      "wmi.logical_disk_latency",
		Priority: prioDiskAvgLatency,
		Dims: Dims{
			{ID: "logical_disk_%s_read_latency", Name: "read", Algo: module.Incremental},
			{ID: "logical_disk_%s_write_latency", Name: "write", Algo: module.Incremental},
		},
	}
)

func newOSCharts() Charts {
	return Charts{
		osProcessesChart.Copy(),
		osUsersChart.Copy(),
		osMemoryUsage.Copy(),
		osPagingFilesUsageChart.Copy(),
	}
}

var (
	osProcessesChart = Chart{
		ID:       "os_processes",
		Title:    "Processes",
		Units:    "number",
		Fam:      "os",
		Ctx:      "wmi.os_processes",
		Priority: prioOSProcesses,
		Dims: Dims{
			{ID: "os_processes", Name: "processes"},
		},
		Vars: Vars{
			{ID: "os_processes_limit"},
		},
	}
	osUsersChart = Chart{
		ID:       "os_users",
		Title:    "Number of Users",
		Units:    "users",
		Fam:      "os",
		Ctx:      "wmi.os_users",
		Priority: prioOSUsers,
		Dims: Dims{
			{ID: "os_users", Name: "users"},
		},
	}
	osMemoryUsage = Chart{
		ID:       "os_visible_memory_usage",
		Title:    "Visible Memory Usage",
		Units:    "bytes",
		Fam:      "os",
		Ctx:      "wmi.os_visible_memory_usage",
		Type:     module.Stacked,
		Priority: prioOSVisibleMemoryUsage,
		Dims: Dims{
			{ID: "os_physical_memory_free_bytes", Name: "free", Div: 1000},
			{ID: "os_visible_memory_used_bytes", Name: "used", Div: 1000},
		},
		Vars: Vars{
			{ID: "os_visible_memory_bytes"},
		},
	}
	osPagingFilesUsageChart = Chart{
		ID:       "os_paging_files_usage",
		Title:    "Paging Files Usage",
		Units:    "bytes",
		Fam:      "os",
		Ctx:      "wmi.os_paging_files_usage",
		Type:     module.Stacked,
		Priority: prioOSPagingUsage,
		Dims: Dims{
			{ID: "os_paging_free_bytes", Name: "free", Div: 1000},
			{ID: "os_paging_used_bytes", Name: "used", Div: 1000},
		},
		Vars: Vars{
			{ID: "os_paging_limit_bytes"},
		},
	}
)

func newSystemCharts() Charts {
	return Charts{
		systemThreadsChart.Copy(),
		systemUptimeChart.Copy(),
	}
}

var (
	systemThreadsChart = Chart{
		ID:       "system_threads",
		Title:    "Threads",
		Units:    "number",
		Fam:      "system",
		Ctx:      "wmi.system_threads",
		Priority: prioSystemThreads,
		Dims: Dims{
			{ID: "system_threads", Name: "threads"},
		},
	}
	systemUptimeChart = Chart{
		ID:       "system_uptime",
		Title:    "Uptime",
		Units:    "seconds",
		Fam:      "system",
		Ctx:      "wmi.system_uptime",
		Priority: prioSystemUptime,
		Dims: Dims{
			{ID: "system_up_time", Name: "time"},
		},
	}
)

func newLogonCharts() Charts {
	return Charts{
		logonSessionsChart.Copy(),
	}
}

var (
	logonSessionsChart = Chart{
		ID:       "logon_active_sessions_by_type",
		Title:    "Active User Logon Sessions By Type",
		Units:    "sessions",
		Fam:      "logon",
		Ctx:      "wmi.logon_type_sessions",
		Type:     module.Stacked,
		Priority: prioLogonSessions,
		Dims: Dims{
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

func newThermalzoneCharts() Charts {
	return Charts{
		thermalzoneTemperatureChart.Copy(),
	}
}

var (
	thermalzoneTemperatureChart = Chart{
		ID:       "thermalzone_temperature",
		Title:    "Thermal zone temperature",
		Units:    "celsius",
		Fam:      "thermalzone",
		Ctx:      "wmi.thermalzone_temperature",
		Type:     module.Area,
		Priority: prioThermalzoneTemperature,
	}
)

func newCollectionCharts() *Charts {
	return &Charts{
		collectionDurationChart.Copy(),
		collectionsStatusChart.Copy(),
	}
}

var (
	collectionDurationChart = Chart{
		ID:       "collector_duration",
		Title:    "Duration",
		Units:    "ms",
		Fam:      "collection",
		Ctx:      "cpu.collector_duration",
		Priority: prioCollectionDuration,
	}
	collectionsStatusChart = Chart{
		ID:       "collector_success",
		Title:    "Collection Status",
		Units:    "bool",
		Fam:      "collection",
		Ctx:      "cpu.collector_success",
		Priority: prioCollectionStatus,
	}
)

func newChartFromTemplate(template Chart, id string) *Chart {
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

func addCPUCoreCharts(charts *Charts, coreID string) error {
	for _, chart := range newCPUCoreCharts() {
		chart = newChartFromTemplate(*chart, coreID)
		if err := charts.Add(chart); err != nil {
			return err
		}
	}
	return nil
}

func addNICCharts(charts *Charts, nicID string) error {
	for _, chart := range newNICCharts() {
		chart = newChartFromTemplate(*chart, nicID)
		if err := charts.Add(chart); err != nil {
			return err
		}
	}
	return nil
}

func addLogicalDiskCharts(charts *Charts, diskID string) error {
	for _, chart := range newDiskCharts() {
		chart = newChartFromTemplate(*chart, diskID)
		if err := charts.Add(chart); err != nil {
			return err
		}
	}
	return nil
}

func removeLogicalDiskFromCharts(charts *Charts, diskID string) {
	for _, chart := range *charts {
		if !strings.HasPrefix(chart.ID, fmt.Sprintf("logical_disk_%s", diskID)) {
			continue
		}
		chart.MarkRemove()
		chart.MarkNotCreated()
	}
}

func addDimToCPUDPCsChart(charts *Charts, coreID string) error {
	chart := charts.Get(cpuDPCsChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", cpuDPCsChart.ID)
	}
	dim := &Dim{
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

func addDimToCPUInterruptsChart(charts *Charts, coreID string) error {
	chart := charts.Get(cpuInterruptsChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", cpuInterruptsChart.ID)
	}
	dim := &Dim{
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

func addDimToThermalzoneTemperatureChart(charts *Charts, zoneName string) error {
	chart := charts.Get(thermalzoneTemperatureChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", thermalzoneTemperatureChart.ID)
	}
	dim := &Dim{
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

func addDimToCollectionDurationChart(charts *Charts, colName string) error {
	chart := charts.Get(collectionDurationChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", collectionDurationChart.ID)
	}
	dim := &Dim{
		ID:   colName + "_collection_duration",
		Name: colName,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

func addDimToCollectionStatusChart(charts *Charts, colName string) error {
	chart := charts.Get(collectionsStatusChart.ID)
	if chart == nil {
		return fmt.Errorf("chart '%s' is not in charts", collectionsStatusChart.ID)
	}
	dim := &Dim{
		ID:   colName + "_collection_success",
		Name: colName,
	}
	if err := chart.AddDim(dim); err != nil {
		return err
	}
	chart.MarkNotCreated()
	return nil
}

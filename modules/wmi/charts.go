// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import (
	"fmt"
	"strings"

	"github.com/netdata/go.d.plugin/agent/module"
)

const (
	prioCPUUtil = module.Priority + iota
	prioCPUCoreUtil
	prioCPUInterrupts
	prioCPUDPCs
	prioCPUCoreCState

	prioMemUtil
	prioMemPageFaults
	prioMemSwapUtil
	prioMemSwapOperations
	prioMemSwapPages
	prioMemCache
	prioMemCacheFaults
	prioMemSystemPool

	prioDiskSpaceUsage
	prioDiskBandwidth
	prioDiskOperations
	prioDiskAvgLatency

	prioNICBandwidth
	prioNICPackets
	prioNICErrors
	prioNICDiscards

	prioTCPConnsEstablished
	prioTCPConnsActive
	prioTCPConnsPassive
	prioTCPConnsFailure
	prioTCPConnsReset
	prioTCPSegmentsReceived
	prioTCPSegmentsSent
	prioTCPSegmentsRetransmitted

	prioOSProcesses
	prioOSUsers
	prioOSVisibleMemoryUsage
	prioOSPagingUsage

	prioSystemThreads
	prioSystemUptime

	prioLogonSessions

	prioThermalzoneTemperature

	prioProcessesCPUUtilization
	prioProcessesMemoryUsage
	prioProcessesIOBytes
	prioProcessesIOOperations
	prioProcessesPageFaults
	prioProcessesPageFileBytes
	prioProcessesThreads
	prioProcessesHandles

	prioIISWebsiteTraffic
	prioIISWebsiteRequestsRate
	prioIISWebsiteActiveConnectionsCount
	prioIISWebsiteUsersCount
	prioIISWebsiteConnectionAttemptsRate
	prioIISWebsiteISAPIExtRequestsCount
	prioIISWebsiteISAPIExtRequestsRate
	prioIISWebsiteFTPFileTransferRate
	prioIISWebsiteLogonAttemptsRate
	prioIISWebsiteErrorsRate
	prioIISWebsiteUptime

	prioMSSQLAccessMethodPageSplits
	prioMSSQLCacheHitRatio
	prioMSSQLBufferCheckpointPages
	prioMSSQLBufferPageLifeExpectancy
	prioMSSQLBufManIOPS
	prioMSSQLDatabaseActiveTransactions
	prioMSSQLDatabaseBackupRestoreOperations
	prioMSSQLDatabaseDataFileSize
	prioMSSQLDatabaseLogFlushed
	prioMSSQLDatabaseLogFlushes
	prioMSSQLDatabaseTransactions
	prioMSSQLDatabaseWriteTransactions
	prioMSSQLBlockedProcess
	prioMSSQLUserConnections
	prioMSSQLLocksLockWait
	prioMSSQLMemmgrPendingMemoryGrants
	prioMSSQLMemTotalServer
	prioMSSQLStatsAutoParameterization
	prioMSSQLStatsSafeAutoParameterization
	prioMSSQLStatsCompilations
	prioMSSQLStatsRecompilations

	prioServiceState
	prioServiceStatus

	prioADDRAReplicationIntersiteCompressedTraffic
	prioADDRAReplicationIntrasiteCompressedTraffic
	prioADDRAReplicationSyncObjectsRemaining
	prioADDRAReplicationPropertiesUpdated
	prioADDRAReplicationPropertiesFiltered
	prioADDRAReplicationObjectsFiltered
	prioADReplicationPendingSyncs
	prioADDRASyncRequests
	prioADDirectoryServiceThreadsInUse
	prioADLDAPBindTime
	prioADBindsTotal
	prioADLDAPSearchesTotal

	prioADCSCertTemplateRequests
	prioADCSCertTemplateRequestProcessingTime
	prioADCSCertTemplateRetrievals
	prioADCSCertTemplateRetrievalProcessingTime
	prioADCSCertTemplateFailedRequests
	prioADCSCertTemplateIssuesRequests
	prioADCSCertTemplatePendingRequests
	prioADCSCertTemplateRequestCryptoSigningTime
	prioADCSCertTemplateRequestPolicyModuleProcessingTime
	prioADCSCertTemplateChallengeResponses
	prioADCSCertTemplateChallengeResponseProcessingTime
	prioADCSCertTemplateSignedCertificateTimestampLists
	prioADCSCertTemplateSignedCertificateTimestampListProcessingTime

	prioADFSADLoginConnectionFailures
	prioADFSCertificateAuthentications
	prioADFSDBArtifactFailures
	prioADFSDBArtifactQueryTimeSeconds
	prioADFSDBConfigFailures
	prioADFSDBConfigQueryTimeSeconds
	prioADFSDeviceAuthentications
	prioADFSExternalAuthentications
	prioADFSFederatedAuthentications
	prioADFSFederationMetadataRequests
	prioADFSOauthAuthorizationRequests
	prioADFSOauthClientAuthentications
	prioADFSOauthClientCredentials
	prioADFSOauthClientPrivkeyJwtAuthentication
	prioADFSOauthClientSecretBasicAuthentications
	prioADFSOauthClientSecretPostAuthentications
	prioADFSOauthClientWindowsAuthentications
	prioADFSOauthLogonCertificateRequests
	prioADFSOauthPasswordGrantRequests
	prioADFSOauthTokenRequestsSuccess
	prioADFSPassiveRequests
	prioADFSPassportAuthentications
	prioADFSPasswordChangeRequests
	prioADFSSAMLPTokenRequests
	prioADFSSSOAuthentications
	prioADFSTokenRequests
	prioADFSUserPasswordAuthentications
	prioADFSWindowsIntegratedAuthentications
	prioADFSWSFedTokenRequestsSuccess
	prioADFSWSTrustTokenRequestsSuccess

	prioCollectorDuration
	prioCollectorStatus
)

// CPU
var (
	cpuCharts = module.Charts{
		cpuUtilChart.Copy(),
	}
	cpuUtilChart = module.Chart{
		ID:       "cpu_utilization_total",
		Title:    "Total CPU Utilization (all cores)",
		Units:    "percentage",
		Fam:      "cpu",
		Ctx:      "wmi.cpu_utilization_total",
		Type:     module.Stacked,
		Priority: prioCPUUtil,
		Dims: module.Dims{
			{ID: "cpu_idle_time", Name: "idle", Algo: module.PercentOfIncremental, Div: 1000, DimOpts: module.DimOpts{Hidden: true}},
			{ID: "cpu_dpc_time", Name: "dpc", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_user_time", Name: "user", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_privileged_time", Name: "privileged", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_interrupt_time", Name: "interrupt", Algo: module.PercentOfIncremental, Div: 1000},
		},
	}
)

// CPU core
var (
	cpuCoreChartsTmpl = module.Charts{
		cpuCoreUtilChartTmpl.Copy(),
		cpuCoreInterruptsChartTmpl.Copy(),
		cpuDPCsChartTmpl.Copy(),
		cpuCoreCStateChartTmpl.Copy(),
	}
	cpuCoreUtilChartTmpl = module.Chart{
		ID:       "core_%s_cpu_utilization",
		Title:    "Core CPU Utilization",
		Units:    "percentage",
		Fam:      "cpu",
		Ctx:      "wmi.cpu_core_utilization",
		Type:     module.Stacked,
		Priority: prioCPUCoreUtil,
		Dims: module.Dims{
			{ID: "cpu_core_%s_idle_time", Name: "idle", Algo: module.PercentOfIncremental, Div: precision, DimOpts: module.DimOpts{Hidden: true}},
			{ID: "cpu_core_%s_dpc_time", Name: "dpc", Algo: module.PercentOfIncremental, Div: precision},
			{ID: "cpu_core_%s_user_time", Name: "user", Algo: module.PercentOfIncremental, Div: precision},
			{ID: "cpu_core_%s_privileged_time", Name: "privileged", Algo: module.PercentOfIncremental, Div: precision},
			{ID: "cpu_core_%s_interrupt_time", Name: "interrupt", Algo: module.PercentOfIncremental, Div: precision},
		},
	}
	cpuCoreInterruptsChartTmpl = module.Chart{
		ID:       "cpu_core_%s_interrupts",
		Title:    "Received and Serviced Hardware Interrupts",
		Units:    "interrupts/s",
		Fam:      "cpu",
		Ctx:      "wmi.cpu_core_interrupts",
		Priority: prioCPUInterrupts,
		Dims: module.Dims{
			{ID: "cpu_core_%s_interrupts", Name: "interrupts", Algo: module.Incremental},
		},
	}
	cpuDPCsChartTmpl = module.Chart{
		ID:       "cpu_core_%s_dpcs",
		Title:    "Received and Serviced Deferred Procedure Calls (DPC)",
		Units:    "dpc/s",
		Fam:      "cpu",
		Ctx:      "wmi.cpu_core_dpcs",
		Priority: prioCPUDPCs,
		Dims: module.Dims{
			{ID: "cpu_core_%s_dpcs", Name: "dpcs", Algo: module.Incremental},
		},
	}
	cpuCoreCStateChartTmpl = module.Chart{
		ID:       "cpu_core_%s_cpu_cstate",
		Title:    "Core Time Spent in Low-Power Idle State",
		Units:    "percentage",
		Fam:      "cpu",
		Ctx:      "wmi.cpu_core_cstate",
		Type:     module.Stacked,
		Priority: prioCPUCoreCState,
		Dims: module.Dims{
			{ID: "cpu_core_%s_cstate_c1", Name: "c1", Algo: module.PercentOfIncremental, Div: precision},
			{ID: "cpu_core_%s_cstate_c2", Name: "c2", Algo: module.PercentOfIncremental, Div: precision},
			{ID: "cpu_core_%s_cstate_c3", Name: "c3", Algo: module.PercentOfIncremental, Div: precision},
		},
	}
)

// Memory
var (
	memCharts = module.Charts{
		memUtilChart.Copy(),
		memPageFaultsChart.Copy(),
		memSwapUtilChart.Copy(),
		memSwapOperationsChart.Copy(),
		memSwapPagesChart.Copy(),
		memCacheChart.Copy(),
		memCacheFaultsChart.Copy(),
		memSystemPoolChart.Copy(),
	}
	memUtilChart = module.Chart{
		ID:       "memory_utilization",
		Title:    "Memory Utilization",
		Units:    "bytes",
		Fam:      "mem",
		Ctx:      "wmi.memory_utilization",
		Type:     module.Stacked,
		Priority: prioMemUtil,
		Dims: module.Dims{
			{ID: "memory_available_bytes", Name: "available"},
			{ID: "memory_used_bytes", Name: "used"},
		},
	}
	memPageFaultsChart = module.Chart{
		ID:       "memory_page_faults",
		Title:    "Memory Page Faults",
		Units:    "pgfaults/s",
		Fam:      "mem",
		Ctx:      "wmi.memory_page_faults",
		Priority: prioMemPageFaults,
		Dims: module.Dims{
			{ID: "memory_page_faults_total", Name: "page_faults", Algo: module.Incremental},
		},
	}
	memSwapUtilChart = module.Chart{
		ID:       "memory_swap_utilization",
		Title:    "Swap Utilization",
		Units:    "bytes",
		Fam:      "mem",
		Ctx:      "wmi.memory_swap_utilization",
		Type:     module.Stacked,
		Priority: prioMemSwapUtil,
		Dims: module.Dims{
			{ID: "memory_not_committed_bytes", Name: "available"},
			{ID: "memory_committed_bytes", Name: "used"},
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
			{ID: "memory_swap_page_reads_total", Name: "read", Algo: module.Incremental},
			{ID: "memory_swap_page_writes_total", Name: "write", Algo: module.Incremental, Mul: -1},
		},
	}
	memSwapPagesChart = module.Chart{
		ID:       "memory_swap_pages",
		Title:    "Swap Pages",
		Units:    "pages/s",
		Fam:      "mem",
		Ctx:      "wmi.memory_swap_pages",
		Priority: prioMemSwapPages,
		Dims: module.Dims{
			{ID: "memory_swap_pages_read_total", Name: "read", Algo: module.Incremental},
			{ID: "memory_swap_pages_written_total", Name: "written", Algo: module.Incremental, Mul: -1},
		},
	}
	memCacheChart = module.Chart{
		ID:       "memory_cached",
		Title:    "Cached",
		Units:    "bytes",
		Fam:      "mem",
		Ctx:      "wmi.memory_cached",
		Type:     module.Area,
		Priority: prioMemCache,
		Dims: module.Dims{
			{ID: "memory_cache_total", Name: "cached"},
		},
	}
	memCacheFaultsChart = module.Chart{
		ID:       "memory_cache_faults",
		Title:    "Cache Faults",
		Units:    "faults/s",
		Fam:      "mem",
		Ctx:      "wmi.memory_cache_faults",
		Priority: prioMemCacheFaults,
		Dims: module.Dims{
			{ID: "memory_cache_faults_total", Name: "cache_faults", Algo: module.Incremental},
		},
	}
	memSystemPoolChart = module.Chart{
		ID:       "memory_system_pool",
		Title:    "System Memory Pool",
		Units:    "bytes",
		Fam:      "mem",
		Ctx:      "wmi.memory_system_pool",
		Type:     module.Stacked,
		Priority: prioMemSystemPool,
		Dims: module.Dims{
			{ID: "memory_pool_paged_bytes", Name: "paged"},
			{ID: "memory_pool_nonpaged_bytes_total", Name: "non-paged"},
		},
	}
)

// Logical Disks
var (
	diskChartsTmpl = module.Charts{
		diskSpaceUsageChartTmpl.Copy(),
		diskBandwidthChartTmpl.Copy(),
		diskOperationsChartTmpl.Copy(),
		diskAvgLatencyChartTmpl.Copy(),
	}
	diskSpaceUsageChartTmpl = module.Chart{
		ID:       "logical_disk_%s_space_usage",
		Title:    "Space usage",
		Units:    "bytes",
		Fam:      "disk",
		Ctx:      "wmi.logical_disk_space_usage",
		Type:     module.Stacked,
		Priority: prioDiskSpaceUsage,
		Dims: module.Dims{
			{ID: "logical_disk_%s_free_space", Name: "free"},
			{ID: "logical_disk_%s_used_space", Name: "used"},
		},
	}
	diskBandwidthChartTmpl = module.Chart{
		ID:       "logical_disk_%s_bandwidth",
		Title:    "Bandwidth",
		Units:    "bytes/s",
		Fam:      "disk",
		Ctx:      "wmi.logical_disk_bandwidth",
		Type:     module.Area,
		Priority: prioDiskBandwidth,
		Dims: module.Dims{
			{ID: "logical_disk_%s_read_bytes_total", Name: "read", Algo: module.Incremental},
			{ID: "logical_disk_%s_write_bytes_total", Name: "write", Algo: module.Incremental, Mul: -1},
		},
	}
	diskOperationsChartTmpl = module.Chart{
		ID:       "logical_disk_%s_operations",
		Title:    "Operations",
		Units:    "operations/s",
		Fam:      "disk",
		Ctx:      "wmi.logical_disk_operations",
		Priority: prioDiskOperations,
		Dims: module.Dims{
			{ID: "logical_disk_%s_reads_total", Name: "reads", Algo: module.Incremental},
			{ID: "logical_disk_%s_writes_total", Name: "writes", Algo: module.Incremental, Mul: -1},
		},
	}
	diskAvgLatencyChartTmpl = module.Chart{
		ID:       "logical_disk_%s_latency",
		Title:    "Average Read/Write Latency",
		Units:    "seconds",
		Fam:      "disk",
		Ctx:      "wmi.logical_disk_latency",
		Priority: prioDiskAvgLatency,
		Dims: module.Dims{
			{ID: "logical_disk_%s_read_latency", Name: "read", Algo: module.Incremental, Div: precision},
			{ID: "logical_disk_%s_write_latency", Name: "write", Algo: module.Incremental, Div: precision},
		},
	}
)

// Network interfaces
var (
	nicChartsTmpl = module.Charts{
		nicBandwidthChartTmpl.Copy(),
		nicPacketsChartTmpl.Copy(),
		nicErrorsChartTmpl.Copy(),
		nicDiscardsChartTmpl.Copy(),
	}
	nicBandwidthChartTmpl = module.Chart{
		ID:       "nic_%s_bandwidth",
		Title:    "Bandwidth",
		Units:    "kilobits/s",
		Fam:      "net",
		Ctx:      "wmi.net_nic_bandwidth",
		Type:     module.Area,
		Priority: prioNICBandwidth,
		Dims: module.Dims{
			{ID: "net_nic_%s_bytes_received", Name: "received", Algo: module.Incremental, Div: 1000},
			{ID: "net_nic_%s_bytes_sent", Name: "sent", Algo: module.Incremental, Mul: -1, Div: 1000},
		},
	}
	nicPacketsChartTmpl = module.Chart{
		ID:       "nic_%s_packets",
		Title:    "Packets",
		Units:    "packets/s",
		Fam:      "net",
		Ctx:      "wmi.net_nic_packets",
		Priority: prioNICPackets,
		Dims: module.Dims{
			{ID: "net_nic_%s_packets_received_total", Name: "received", Algo: module.Incremental},
			{ID: "net_nic_%s_packets_sent_total", Name: "sent", Algo: module.Incremental, Mul: -1},
		},
	}
	nicErrorsChartTmpl = module.Chart{
		ID:       "nic_%s_errors",
		Title:    "Errors",
		Units:    "errors/s",
		Fam:      "net",
		Ctx:      "wmi.net_nic_errors",
		Priority: prioNICErrors,
		Dims: module.Dims{
			{ID: "net_nic_%s_packets_received_errors", Name: "inbound", Algo: module.Incremental},
			{ID: "net_nic_%s_packets_outbound_errors", Name: "outbound", Algo: module.Incremental, Mul: -1},
		},
	}
	nicDiscardsChartTmpl = module.Chart{
		ID:       "nic_%s_discarded",
		Title:    "Discards",
		Units:    "discards/s",
		Fam:      "net",
		Ctx:      "wmi.net_nic_discarded",
		Priority: prioNICDiscards,
		Dims: module.Dims{
			{ID: "net_nic_%s_packets_received_discarded", Name: "inbound", Algo: module.Incremental},
			{ID: "net_nic_%s_packets_outbound_discarded", Name: "outbound", Algo: module.Incremental, Mul: -1},
		},
	}
)

// TCP
var (
	tcpCharts = module.Charts{
		tcpConnsActiveChart.Copy(),
		tcpConnsEstablishedChart.Copy(),
		tcpConnsFailuresChart.Copy(),
		tcpConnsPassiveChart.Copy(),
		tcpConnsResetsChart.Copy(),
		tcpSegmentsReceivedChart.Copy(),
		tcpSegmentsRetransmittedChart.Copy(),
		tcpSegmentsSentChart.Copy(),
	}
	tcpConnsEstablishedChart = module.Chart{
		ID:       "tcp_conns_established",
		Title:    "TCP established connections",
		Units:    "connections",
		Fam:      "tcp",
		Ctx:      "wmi.tcp_conns_established",
		Priority: prioTCPConnsEstablished,
		Dims: module.Dims{
			{ID: "tcp_ipv4_conns_established", Name: "ipv4"},
			{ID: "tcp_ipv6_conns_established", Name: "ipv6"},
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
			{ID: "tcp_ipv4_conns_active", Name: "ipv4", Algo: module.Incremental},
			{ID: "tcp_ipv6_conns_active", Name: "ipv6", Algo: module.Incremental},
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
			{ID: "tcp_ipv4_conns_passive", Name: "ipv4", Algo: module.Incremental},
			{ID: "tcp_ipv6_conns_passive", Name: "ipv6", Algo: module.Incremental},
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
			{ID: "tcp_ipv4_conns_failures", Name: "ipv4", Algo: module.Incremental},
			{ID: "tcp_ipv6_conns_failures", Name: "ipv6", Algo: module.Incremental},
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
			{ID: "tcp_ipv4_conns_resets", Name: "ipv4", Algo: module.Incremental},
			{ID: "tcp_ipv6_conns_resets", Name: "ipv6", Algo: module.Incremental},
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
			{ID: "tcp_ipv4_segments_received", Name: "ipv4", Algo: module.Incremental},
			{ID: "tcp_ipv6_segments_received", Name: "ipv6", Algo: module.Incremental},
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
			{ID: "tcp_ipv4_segments_sent", Name: "ipv4", Algo: module.Incremental},
			{ID: "tcp_ipv6_segments_sent", Name: "ipv6", Algo: module.Incremental},
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
			{ID: "tcp_ipv4_segments_retransmitted", Name: "ipv4", Algo: module.Incremental},
			{ID: "tcp_ipv6_segments_retransmitted", Name: "ipv6", Algo: module.Incremental},
		},
	}
)

// OS
var (
	osCharts = module.Charts{
		osProcessesChart.Copy(),
		osUsersChart.Copy(),
		osMemoryUsage.Copy(),
		osPagingFilesUsageChart.Copy(),
	}
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
			{ID: "os_physical_memory_free_bytes", Name: "free"},
			{ID: "os_visible_memory_used_bytes", Name: "used"},
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
			{ID: "os_paging_free_bytes", Name: "free"},
			{ID: "os_paging_used_bytes", Name: "used"},
		},
		Vars: module.Vars{
			{ID: "os_paging_limit_bytes"},
		},
	}
)

// System
var (
	systemCharts = module.Charts{
		systemThreadsChart.Copy(),
		systemUptimeChart.Copy(),
	}
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

// IIS
var (
	iisWebsiteChartsTmpl = module.Charts{
		iisWebsiteTrafficChartTempl.Copy(),
		iisWebsiteRequestsRateChartTmpl.Copy(),
		iisWebsiteActiveConnectionsCountChartTmpl.Copy(),
		iisWebsiteUsersCountChartTmpl.Copy(),
		iisWebsiteConnectionAttemptsRate.Copy(),
		iisWebsiteISAPIExtRequestsCountChartTmpl.Copy(),
		iisWebsiteISAPIExtRequestsRateChartTmpl.Copy(),
		iisWebsiteFTPFileTransferRateChartTempl.Copy(),
		iisWebsiteLogonAttemptsRateChartTmpl.Copy(),
		iisWebsiteErrorsRateChart.Copy(),
		iisWebsiteUptimeChartTmpl.Copy(),
	}
	iisWebsiteTrafficChartTempl = module.Chart{
		ID:       "iis_website_%s_traffic",
		Title:    "Website traffic",
		Units:    "bytes/s",
		Fam:      "iis",
		Ctx:      "wmi.iis_website_traffic",
		Type:     module.Area,
		Priority: prioIISWebsiteTraffic,
		Dims: module.Dims{
			{ID: "iis_website_%s_received_bytes_total", Name: "received", Algo: module.Incremental},
			{ID: "iis_website_%s_sent_bytes_total", Name: "sent", Algo: module.Incremental, Mul: -1},
		},
	}
	iisWebsiteRequestsRateChartTmpl = module.Chart{
		ID:       "iis_website_%s_requests_rate",
		Title:    "Website requests rate",
		Units:    "requests/s",
		Fam:      "iis",
		Ctx:      "wmi.iis_website_requests_rate",
		Priority: prioIISWebsiteRequestsRate,
		Dims: module.Dims{
			{ID: "iis_website_%s_requests_total", Name: "requests", Algo: module.Incremental},
		},
	}
	iisWebsiteActiveConnectionsCountChartTmpl = module.Chart{
		ID:       "iis_website_%s_active_connections_count",
		Title:    "Website active connections",
		Units:    "connections",
		Fam:      "iis",
		Ctx:      "wmi.iis_website_active_connections_count",
		Priority: prioIISWebsiteActiveConnectionsCount,
		Dims: module.Dims{
			{ID: "iis_website_%s_current_connections", Name: "active"},
		},
	}
	iisWebsiteUsersCountChartTmpl = module.Chart{
		ID:       "iis_website_%s_users_count",
		Title:    "Website users with pending requests",
		Units:    "users",
		Fam:      "iis",
		Ctx:      "wmi.iis_website_users_count",
		Type:     module.Stacked,
		Priority: prioIISWebsiteUsersCount,
		Dims: module.Dims{
			{ID: "iis_website_%s_current_anonymous_users", Name: "anonymous"},
			{ID: "iis_website_%s_current_non_anonymous_users", Name: "non_anonymous"},
		},
	}
	iisWebsiteConnectionAttemptsRate = module.Chart{
		ID:       "iis_website_%s_connection_attempts_rate",
		Title:    "Website connections attempts",
		Units:    "attempts/s",
		Fam:      "iis",
		Ctx:      "wmi.iis_website_connection_attempts_rate",
		Priority: prioIISWebsiteConnectionAttemptsRate,
		Dims: module.Dims{
			{ID: "iis_website_%s_connection_attempts_all_instances_total", Name: "connection", Algo: module.Incremental},
		},
	}
	iisWebsiteISAPIExtRequestsCountChartTmpl = module.Chart{
		ID:       "iis_website_%s_isapi_extension_requests_count",
		Title:    "ISAPI extension requests",
		Units:    "requests",
		Fam:      "iis",
		Ctx:      "wmi.iis_website_isapi_extension_requests_count",
		Priority: prioIISWebsiteISAPIExtRequestsCount,
		Dims: module.Dims{
			{ID: "iis_website_%s_current_isapi_extension_requests", Name: "isapi"},
		},
	}
	iisWebsiteISAPIExtRequestsRateChartTmpl = module.Chart{
		ID:       "iis_website_%s_isapi_extension_requests_rate",
		Title:    "Website extensions request",
		Units:    "requests/s",
		Fam:      "iis",
		Ctx:      "wmi.iis_website_isapi_extension_requests_rate",
		Priority: prioIISWebsiteISAPIExtRequestsRate,
		Dims: module.Dims{
			{ID: "iis_website_%s_isapi_extension_requests_total", Name: "isapi", Algo: module.Incremental},
		},
	}
	iisWebsiteFTPFileTransferRateChartTempl = module.Chart{
		ID:       "iis_website_%s_ftp_file_transfer_rate",
		Title:    "Website FTP file transfer rate",
		Units:    "files/s",
		Fam:      "iis",
		Ctx:      "wmi.iis_website_ftp_file_transfer_rate",
		Priority: prioIISWebsiteFTPFileTransferRate,
		Dims: module.Dims{
			{ID: "iis_website_%s_files_received_total", Name: "received", Algo: module.Incremental},
			{ID: "iis_website_%s_files_sent_total", Name: "sent", Algo: module.Incremental},
		},
	}
	iisWebsiteLogonAttemptsRateChartTmpl = module.Chart{
		ID:       "iis_website_%s_logon_attempts_rate",
		Title:    "Website logon attempts",
		Units:    "attempts/s",
		Fam:      "iis",
		Ctx:      "wmi.iis_website_logon_attempts_rate",
		Priority: prioIISWebsiteLogonAttemptsRate,
		Dims: module.Dims{
			{ID: "iis_website_%s_logon_attempts_total", Name: "logon", Algo: module.Incremental},
		},
	}
	iisWebsiteErrorsRateChart = module.Chart{
		ID:       "iis_website_%s_errors_rate",
		Title:    "Website errors",
		Units:    "errors/s",
		Fam:      "iis",
		Ctx:      "wmi.iis_website_errors_rate",
		Type:     module.Stacked,
		Priority: prioIISWebsiteErrorsRate,
		Dims: module.Dims{
			{ID: "iis_website_%s_locked_errors_total", Name: "document_locked", Algo: module.Incremental},
			{ID: "iis_website_%s_not_found_errors_total", Name: "document_not_found", Algo: module.Incremental},
		},
	}
	iisWebsiteUptimeChartTmpl = module.Chart{
		ID:       "iis_website_%s_uptime",
		Title:    "Website uptime",
		Units:    "seconds",
		Fam:      "iis",
		Ctx:      "wmi.iis_website_uptime",
		Priority: prioIISWebsiteUptime,
		Dims: module.Dims{
			{ID: "iis_website_%s_service_uptime", Name: "uptime"},
		},
	}
)

// MS-SQL
var (
	mssqlInstanceChartsTmpl = module.Charts{
		mssqlAccessMethodPageSplitsChart.Copy(),
		mssqlCacheHitRatioChart.Copy(),
		mssqlBufferCheckpointPageChart.Copy(),
		mssqlBufferPageLifeExpectancyChart.Copy(),
		mssqlBufManIOPSChart.Copy(),
		mssqlBlockedProcessChart.Copy(),
		mssqlLocksWaitChart.Copy(),
		mssqlMemmgrPendingMemoryChart.Copy(),
		mssqlMemmgrTotalServerChart.Copy(),
		mssqlStatsAutoParamChart.Copy(),
		mssqlStatsSafeAutoChart.Copy(),
		mssqlStatsCompilationChart.Copy(),
		mssqlStatsRecompilationChart.Copy(),
		mssqlUserConnectionChart.Copy(),
	}
	mssqlDatabaseChartsTmpl = module.Charts{
		mssqlDatabaseActiveTransactionsChart.Copy(),
		mssqlDatabaseBackupRestoreOperationsChart.Copy(),
		mssqlDatabaseSizeChart.Copy(),
		mssqlDatabaseLogFlushedChart.Copy(),
		mssqlDatabaseLogFlushesChart.Copy(),
		mssqlDatabaseTransactionsChart.Copy(),
		mssqlDatabaseWriteTransactionsChart.Copy(),
	}
	// Access Method:
	// Source: https://learn.microsoft.com/en-us/sql/relational-databases/performance-monitor/sql-server-access-methods-object?view=sql-server-ver16
	mssqlAccessMethodPageSplitsChart = module.Chart{
		ID:       "mssql_instance_%s_accessmethods_page_splits",
		Title:    "Page splits",
		Units:    "splits/s",
		Fam:      "mssql",
		Ctx:      "wmi.mssql_instance_accessmethods_page_splits",
		Priority: prioMSSQLAccessMethodPageSplits,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_accessmethods_page_splits", Name: "page", Algo: module.Incremental},
		},
	}
	// Buffer Management
	// Source: https://learn.microsoft.com/en-us/sql/relational-databases/performance-monitor/sql-server-buffer-manager-object?view=sql-server-ver16
	mssqlCacheHitRatioChart = module.Chart{
		ID:       "mssql_instance_%s_cache_hit_ratio",
		Title:    "Buffer Cache hit ratio",
		Units:    "percentage",
		Fam:      "mssql",
		Ctx:      "wmi.mssql_instance_cache_hit_ratio",
		Priority: prioMSSQLCacheHitRatio,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_cache_hit_ratio", Name: "hit_ratio"},
		},
	}
	mssqlBufferCheckpointPageChart = module.Chart{
		ID:       "mssql_instance_%s_bufman_checkpoint_pages",
		Title:    "Flushed pages",
		Units:    "pages/s",
		Fam:      "mssql",
		Ctx:      "wmi.mssql_instance_bufman_checkpoint_pages",
		Priority: prioMSSQLBufferCheckpointPages,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_bufman_checkpoint_pages", Name: "flushed", Algo: module.Incremental},
		},
	}
	mssqlBufferPageLifeExpectancyChart = module.Chart{
		ID:       "mssql_instance_%s_bufman_page_life_expectancy",
		Title:    "Page life expectancy",
		Units:    "seconds",
		Fam:      "mssql",
		Ctx:      "wmi.mssql_instance_bufman_page_life_expectancy",
		Priority: prioMSSQLBufferPageLifeExpectancy,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_bufman_page_life_expectancy_seconds", Name: "life_expectancy"},
		},
	}
	mssqlBufManIOPSChart = module.Chart{
		ID:       "mssql_instance_%s_bufman_iops",
		Title:    "Number of pages input and output",
		Units:    "pages/s",
		Fam:      "mssql",
		Ctx:      "wmi.mssql_instance_bufman_iops",
		Priority: prioMSSQLBufManIOPS,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_bufman_page_reads", Name: "read", Algo: module.Incremental},
			{ID: "mssql_instance_%s_bufman_page_writes", Name: "written", Mul: -1, Algo: module.Incremental},
		},
	}
	// General Statistic
	// Source: https://learn.microsoft.com/en-us/sql/relational-databases/performance-monitor/sql-server-general-statistics-object?view=sql-server-ver16
	mssqlBlockedProcessChart = module.Chart{
		ID:       "mssql_instance_%s_blocked_process",
		Title:    "Blocked processes",
		Units:    "process",
		Fam:      "mssql",
		Ctx:      "wmi.mssql_instance_blocked_processes",
		Priority: prioMSSQLBlockedProcess,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_genstats_blocked_processes", Name: "blocked"},
		},
	}
	mssqlUserConnectionChart = module.Chart{
		ID:       "mssql_instance_%s_user_connection",
		Title:    "User connections",
		Units:    "connections",
		Fam:      "mssql",
		Ctx:      "wmi.mssql_instance_user_connection",
		Priority: prioMSSQLUserConnections,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_genstats_user_connections", Name: "user"},
		},
	}
	// Lock Wait
	// Source: https://learn.microsoft.com/en-us/sql/relational-databases/performance-monitor/sql-server-locks-object?view=sql-server-ver16
	mssqlLocksWaitChart = module.Chart{
		ID:       "mssql_instance_%s_locks_lock_wait",
		Title:    "Lock requests that required the caller to wait",
		Units:    "locks/s",
		Fam:      "mssql",
		Ctx:      "wmi.mssql_instance_locks_lock_wait",
		Priority: prioMSSQLLocksLockWait,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_resource_AllocUnit_locks_lock_wait_seconds", Name: "alloc_unit", Algo: module.Incremental},
			{ID: "mssql_instance_%s_resource_Application_locks_lock_wait_seconds", Name: "application", Algo: module.Incremental},
			{ID: "mssql_instance_%s_resource_Database_locks_lock_wait_seconds", Name: "database", Algo: module.Incremental},
			{ID: "mssql_instance_%s_resource_Extent_locks_lock_wait_seconds", Name: "extent", Algo: module.Incremental},
			{ID: "mssql_instance_%s_resource_File_locks_lock_wait_seconds", Name: "file", Algo: module.Incremental},
			{ID: "mssql_instance_%s_resource_HoBT_locks_lock_wait_seconds", Name: "hobt", Algo: module.Incremental},
			{ID: "mssql_instance_%s_resource_Key_locks_lock_wait_seconds", Name: "key", Algo: module.Incremental},
			{ID: "mssql_instance_%s_resource_Metadata_locks_lock_wait_seconds", Name: "metadata", Algo: module.Incremental},
			{ID: "mssql_instance_%s_resource_OIB_locks_lock_wait_seconds", Name: "oib", Algo: module.Incremental},
			{ID: "mssql_instance_%s_resource_Object_locks_lock_wait_seconds", Name: "object", Algo: module.Incremental},
			{ID: "mssql_instance_%s_resource_Page_locks_lock_wait_seconds", Name: "page", Algo: module.Incremental},
			{ID: "mssql_instance_%s_resource_RID_locks_lock_wait_seconds", Name: "rid", Algo: module.Incremental},
			{ID: "mssql_instance_%s_resource_RowGroup_locks_lock_wait_seconds", Name: "row_group", Algo: module.Incremental},
			{ID: "mssql_instance_%s_resource_Xact_locks_lock_wait_seconds", Name: "xact", Algo: module.Incremental},
		},
	}
	// Memory Manager
	// Source: https://learn.microsoft.com/en-us/sql/relational-databases/performance-monitor/sql-server-memory-manager-object?view=sql-server-ver16
	mssqlMemmgrPendingMemoryChart = module.Chart{
		ID:       "mssql_instance_%s_memmgr_pending_memory_grants",
		Title:    "Process waiting for memory grant",
		Units:    "process",
		Fam:      "mssql",
		Ctx:      "wmi.mssql_instance_memmgr_pending_memory_grants",
		Priority: prioMSSQLMemmgrPendingMemoryGrants,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_memmgr_pending_memory_grants", Name: "pending"},
		},
	}
	mssqlMemmgrTotalServerChart = module.Chart{
		ID:       "mssql_instance_%s_memmgr_server_memory",
		Title:    "Memory committed",
		Units:    "bytes",
		Fam:      "mssql",
		Ctx:      "wmi.mssql_instance_memmgr_server_memory",
		Priority: prioMSSQLMemTotalServer,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_memmgr_total_server_memory_bytes", Name: "memory"},
		},
	}
	// SQL Statistic
	// Source: https://learn.microsoft.com/en-us/sql/relational-databases/performance-monitor/sql-server-sql-statistics-object?view=sql-server-ver16
	mssqlStatsAutoParamChart = module.Chart{
		ID:       "mssql_instance_%s_sqlstats_auto_parameterization_attempts",
		Title:    "Failed auto-parameterization attempts",
		Units:    "attempts/s",
		Fam:      "mssql",
		Ctx:      "wmi.mssql_instance_sqlstats_auto_parameterization_attempts",
		Priority: prioMSSQLStatsAutoParameterization,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_sqlstats_auto_parameterization_attempts", Name: "failed", Algo: module.Incremental},
		},
	}
	mssqlStatsSafeAutoChart = module.Chart{
		ID:       "mssql_instance_%s_sqlstats_safe_auto_parameterization_attempts",
		Title:    "Safe auto-parameterization attempts",
		Units:    "attempts/s",
		Fam:      "mssql",
		Ctx:      "wmi.mssql_instance_sqlstats_safe_auto_parameterization_attempts",
		Priority: prioMSSQLStatsSafeAutoParameterization,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_sqlstats_safe_auto_parameterization_attempts", Name: "safe", Algo: module.Incremental},
		},
	}
	mssqlStatsCompilationChart = module.Chart{
		ID:       "mssql_instance_%s_sqlstats_sql_compilations",
		Title:    "SQL compilations",
		Units:    "compilations/s",
		Fam:      "mssql",
		Ctx:      "wmi.mssql_instance_sqlstats_sql_compilations",
		Priority: prioMSSQLStatsCompilations,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_sqlstats_sql_compilations", Name: "compilations", Algo: module.Incremental},
		},
	}
	mssqlStatsRecompilationChart = module.Chart{
		ID:       "mssql_instance_%s_sqlstats_sql_recompilations",
		Title:    "SQL re-compilations",
		Units:    "recompiles/s",
		Fam:      "mssql",
		Ctx:      "wmi.mssql_instance_sqlstats_sql_recompilations",
		Priority: prioMSSQLStatsRecompilations,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_sqlstats_sql_recompilations", Name: "recompiles", Algo: module.Incremental},
		},
	}

	// Database
	// Source: https://learn.microsoft.com/en-us/sql/relational-databases/performance-monitor/sql-server-databases-object?view=sql-server-2017
	mssqlDatabaseActiveTransactionsChart = module.Chart{
		ID:       "mssql_db_%s_instance_%s_active_transactions",
		Title:    "Active transactions per database",
		Units:    "transactions",
		Fam:      "mssql",
		Ctx:      "wmi.mssql_database_active_transactions",
		Priority: prioMSSQLDatabaseActiveTransactions,
		Dims: module.Dims{
			{ID: "mssql_db_%s_instance_%s_active_transactions", Name: "active"},
		},
	}
	mssqlDatabaseBackupRestoreOperationsChart = module.Chart{
		ID:       "mssql_db_%s_instance_%s_backup_restore_operations",
		Title:    "Backup IO per database",
		Units:    "operations/s",
		Fam:      "mssql",
		Ctx:      "wmi.mssql_database_backup_restore_operations",
		Priority: prioMSSQLDatabaseBackupRestoreOperations,
		Dims: module.Dims{
			{ID: "mssql_db_%s_instance_%s_backup_restore_operations", Name: "backup", Algo: module.Incremental},
		},
	}
	mssqlDatabaseSizeChart = module.Chart{
		ID:       "mssql_db_%s_instance_%s_data_files_size",
		Title:    "Current database size",
		Units:    "bytes",
		Fam:      "mssql",
		Ctx:      "wmi.mssql_database_data_files_size",
		Priority: prioMSSQLDatabaseDataFileSize,
		Dims: module.Dims{
			{ID: "mssql_db_%s_instance_%s_data_files_size_bytes", Name: "size"},
		},
	}
	mssqlDatabaseLogFlushedChart = module.Chart{
		ID:       "mssql_db_%s_instance_%s_log_flushed",
		Title:    "Log flushed",
		Units:    "bytes/s",
		Fam:      "mssql",
		Ctx:      "wmi.mssql_database_log_flushed",
		Priority: prioMSSQLDatabaseLogFlushed,
		Dims: module.Dims{
			{ID: "mssql_db_%s_instance_%s_log_flushed_bytes", Name: "flushed", Algo: module.Incremental},
		},
	}
	mssqlDatabaseLogFlushesChart = module.Chart{
		ID:       "mssql_db_%s_instance_%s_log_flushes",
		Title:    "Log flushes",
		Units:    "flushes/s",
		Fam:      "mssql",
		Ctx:      "wmi.mssql_database_log_flushes",
		Priority: prioMSSQLDatabaseLogFlushes,
		Dims: module.Dims{
			{ID: "mssql_db_%s_instance_%s_log_flushes", Name: "log", Algo: module.Incremental},
		},
	}
	mssqlDatabaseTransactionsChart = module.Chart{
		ID:       "mssql_db_%s_instance_%s_transactions",
		Title:    "Transactions",
		Units:    "transactions/s",
		Fam:      "mssql",
		Ctx:      "wmi.mssql_database_transactions",
		Priority: prioMSSQLDatabaseTransactions,
		Dims: module.Dims{
			{ID: "mssql_db_%s_instance_%s_transactions", Name: "transactions", Algo: module.Incremental},
		},
	}
	mssqlDatabaseWriteTransactionsChart = module.Chart{
		ID:       "mssql_db_%s_instance_%s_write_transactions",
		Title:    "Write transactions",
		Units:    "transactions/s",
		Fam:      "mssql",
		Ctx:      "wmi.mssql_instance_write_transactions",
		Priority: prioMSSQLDatabaseWriteTransactions,
		Dims: module.Dims{
			{ID: "mssql_db_%s_instance_%s_write_transactions", Name: "write", Algo: module.Incremental},
		},
	}
)

// AD
var (
	adCharts = module.Charts{
		adDRAReplicationIntersiteCompressedTrafficChart.Copy(),
		adDRAReplicationIntrasiteCompressedTrafficChart.Copy(),
		adDRAReplicationSyncObjectRemainingChart.Copy(),
		adDRAReplicationObjectsFilteredChart.Copy(),
		adDRAReplicationPropertiesUpdatedChart.Copy(),
		adDRAReplicationPropertiesFilteredChart.Copy(),
		adDRAReplicationPendingSyncsChart.Copy(),
		adDRAReplicationSyncRequestsChart.Copy(),
		adDirectoryServiceThreadsChart.Copy(),
		adLDAPLastBindTimeChart.Copy(),
		adBindsTotalChart.Copy(),
		adLDAPSearchesChart.Copy(),
	}
	adDRAReplicationIntersiteCompressedTrafficChart = module.Chart{
		ID:       "ad_dra_replication_intersite_compressed_traffic",
		Title:    "DRA replication compressed traffic withing site",
		Units:    "bytes/s",
		Fam:      "ad",
		Ctx:      "wmi.ad_dra_replication_intersite_compressed_traffic",
		Priority: prioADDRAReplicationIntersiteCompressedTraffic,
		Type:     module.Area,
		Dims: module.Dims{
			{ID: "ad_replication_data_intersite_bytes_total_inbound", Name: "inbound", Algo: module.Incremental},
			{ID: "ad_replication_data_intersite_bytes_total_outbound", Name: "outbound", Algo: module.Incremental, Mul: -1},
		},
	}
	adDRAReplicationIntrasiteCompressedTrafficChart = module.Chart{
		ID:       "ad_dra_replication_intrasite_compressed_traffic",
		Title:    "DRA replication compressed traffic between sites",
		Units:    "bytes/s",
		Fam:      "ad",
		Ctx:      "wmi.ad_dra_replication_intrasite_compressed_traffic",
		Priority: prioADDRAReplicationIntrasiteCompressedTraffic,
		Type:     module.Area,
		Dims: module.Dims{
			{ID: "ad_replication_data_intrasite_bytes_total_inbound", Name: "inbound", Algo: module.Incremental},
			{ID: "ad_replication_data_intrasite_bytes_total_outbound", Name: "outbound", Algo: module.Incremental, Mul: -1},
		},
	}
	adDRAReplicationSyncObjectRemainingChart = module.Chart{
		ID:       "ad_dra_replication_sync_objects_remaining",
		Title:    "DRA replication full sync objects remaining",
		Units:    "objects",
		Fam:      "ad",
		Ctx:      "wmi.ad_dra_replication_sync_objects_remaining",
		Priority: prioADDRAReplicationSyncObjectsRemaining,
		Dims: module.Dims{
			{ID: "ad_replication_inbound_sync_objects_remaining", Name: "inbound"},
		},
	}
	adDRAReplicationObjectsFilteredChart = module.Chart{
		ID:       "ad_dra_replication_objects_filtered",
		Title:    "DRA replication objects filtered",
		Units:    "objects/s",
		Fam:      "ad",
		Ctx:      "wmi.ad_dra_replication_objects_filtered",
		Priority: prioADDRAReplicationObjectsFiltered,
		Dims: module.Dims{
			{ID: "ad_replication_inbound_objects_filtered_total", Name: "inbound", Algo: module.Incremental},
		},
	}
	adDRAReplicationPropertiesUpdatedChart = module.Chart{
		ID:       "ad_dra_replication_properties_updated",
		Title:    "DRA replication properties updated",
		Units:    "properties/s",
		Fam:      "ad",
		Ctx:      "wmi.ad_dra_replication_properties_updated",
		Priority: prioADDRAReplicationPropertiesUpdated,
		Dims: module.Dims{
			{ID: "ad_replication_inbound_properties_updated_total", Name: "inbound", Algo: module.Incremental},
		},
	}
	adDRAReplicationPropertiesFilteredChart = module.Chart{
		ID:       "ad_dra_replication_properties_filtered",
		Title:    "DRA replication properties filtered",
		Units:    "properties/s",
		Fam:      "ad",
		Ctx:      "wmi.ad_dra_replication_properties_filtered",
		Priority: prioADDRAReplicationPropertiesFiltered,
		Dims: module.Dims{
			{ID: "ad_replication_inbound_properties_filtered_total", Name: "inbound", Algo: module.Incremental},
		},
	}
	adDRAReplicationPendingSyncsChart = module.Chart{
		ID:       "ad_dra_replication_pending_syncs",
		Title:    "DRA replication pending syncs",
		Units:    "syncs",
		Fam:      "ad",
		Ctx:      "wmi.ad_dra_replication_pending_syncs",
		Priority: prioADReplicationPendingSyncs,
		Dims: module.Dims{
			{ID: "ad_replication_pending_synchronizations", Name: "pending"},
		},
	}
	adDRAReplicationSyncRequestsChart = module.Chart{
		ID:       "ad_dra_replication_sync_requests",
		Title:    "DRA replication sync requests",
		Units:    "requests/s",
		Fam:      "ad",
		Ctx:      "wmi.ad_dra_replication_sync_requests",
		Priority: prioADDRASyncRequests,
		Dims: module.Dims{
			{ID: "ad_replication_sync_requests_total", Name: "request", Algo: module.Incremental},
		},
	}
	adDirectoryServiceThreadsChart = module.Chart{
		ID:       "ad_ds_threads",
		Title:    "Directory Service threads",
		Units:    "threads",
		Fam:      "ad",
		Ctx:      "wmi.ad_ds_threads",
		Priority: prioADDirectoryServiceThreadsInUse,
		Dims: module.Dims{
			{ID: "ad_directory_service_threads", Name: "in_use"},
		},
	}
	adLDAPLastBindTimeChart = module.Chart{
		ID:       "ad_ldap_last_bind_time",
		Title:    "LDAP last successful bind time",
		Units:    "seconds",
		Fam:      "ad",
		Ctx:      "wmi.ad_ldap_last_bind_time",
		Priority: prioADLDAPBindTime,
		Dims: module.Dims{
			{ID: "ad_ldap_last_bind_time_seconds", Name: "last_bind"},
		},
	}
	adBindsTotalChart = module.Chart{
		ID:       "ad_binds",
		Title:    "Successful binds",
		Units:    "bind/s",
		Fam:      "ad",
		Ctx:      "wmi.ad_binds",
		Priority: prioADBindsTotal,
		Dims: module.Dims{
			{ID: "ad_binds_total", Name: "binds", Algo: module.Incremental},
		},
	}
	adLDAPSearchesChart = module.Chart{
		ID:       "ad_ldap_searches",
		Title:    "LDAP client search operations",
		Units:    "searches/s",
		Fam:      "ad",
		Ctx:      "wmi.ad_ldap_searches",
		Priority: prioADLDAPSearchesTotal,
		Dims: module.Dims{
			{ID: "ad_ldap_searches_total", Name: "searches", Algo: module.Incremental},
		},
	}
)

// AD CS
var (
	adcsCertTemplateChartsTmpl = module.Charts{
		adcsCertTemplateRequestsChartTmpl.Copy(),
		adcsCertTemplateFailedRequestsChartTmpl.Copy(),
		adcsCertTemplateIssuedRequestsChartTmpl.Copy(),
		adcsCertTemplatePendingRequestsChartTmpl.Copy(),
		adcsCertTemplateRequestProcessingTimeChartTmpl.Copy(),

		adcsCertTemplateRetrievalsChartTmpl.Copy(),
		adcsCertificateRetrievalsTimeChartTmpl.Copy(),
		adcsCertTemplateRequestCryptoSigningTimeChartTmpl.Copy(),
		adcsCertTemplateRequestPolicyModuleProcessingTimeChartTmpl.Copy(),
		adcsCertTemplateChallengeResponseChartTmpl.Copy(),
		adcsCertTemplateChallengeResponseProcessingTimeChartTmpl.Copy(),
		adcsCertTemplateSignedCertificateTimestampListsChartTmpl.Copy(),
		adcsCertTemplateSignedCertificateTimestampListProcessingTimeChartTmpl.Copy(),
	}
	adcsCertTemplateRequestsChartTmpl = module.Chart{
		ID:       "adcs_cert_template%s_requests",
		Title:    "Certificate requests processed",
		Units:    "requests/s",
		Fam:      "adcs",
		Ctx:      "wmi.adcs_cert_template_requests",
		Priority: prioADCSCertTemplateRequests,
		Dims: module.Dims{
			{ID: "adcs_cert_template_%s_requests_total", Name: "requests", Algo: module.Incremental},
		},
	}
	adcsCertTemplateFailedRequestsChartTmpl = module.Chart{
		ID:       "adcs_cert_template_%s_failed_requests",
		Title:    "Certificate failed requests processed",
		Units:    "requests/s",
		Fam:      "adcs",
		Ctx:      "wmi.adcs_cert_template_failed_requests",
		Priority: prioADCSCertTemplateFailedRequests,
		Dims: module.Dims{
			{ID: "adcs_cert_template_%s_failed_requests_total", Name: "failed", Algo: module.Incremental},
		},
	}
	adcsCertTemplateIssuedRequestsChartTmpl = module.Chart{
		ID:       "adcs_cert_template_%s_issued_requests",
		Title:    "Certificate issued requests processed",
		Units:    "requests/s",
		Fam:      "adcs",
		Ctx:      "wmi.adcs_cert_template_issued_requests",
		Priority: prioADCSCertTemplateIssuesRequests,
		Dims: module.Dims{
			{ID: "adcs_cert_template_%s_issued_requests_total", Name: "issued", Algo: module.Incremental},
		},
	}
	adcsCertTemplatePendingRequestsChartTmpl = module.Chart{
		ID:       "adcs_cert_template_%s_pending_requests",
		Title:    "Certificate pending requests processed",
		Units:    "requests/s",
		Fam:      "adcs",
		Ctx:      "wmi.adcs_cert_template_pending_requests",
		Priority: prioADCSCertTemplatePendingRequests,
		Dims: module.Dims{
			{ID: "adcs_cert_template_%s_pending_requests_total", Name: "pending", Algo: module.Incremental},
		},
	}
	adcsCertTemplateRequestProcessingTimeChartTmpl = module.Chart{
		ID:       "adcs_cert_template_%s_request_processing_time",
		Title:    "Certificate last request processing time",
		Units:    "seconds",
		Fam:      "adcs",
		Ctx:      "wmi.adcs_cert_template_request_processing_time",
		Priority: prioADCSCertTemplateRequestProcessingTime,
		Dims: module.Dims{
			{ID: "adcs_cert_template_%s_request_processing_time_seconds", Name: "processing_time", Div: precision},
		},
	}
	adcsCertTemplateRetrievalsChartTmpl = module.Chart{
		ID:       "adcs_cert_template_%s_retrievals",
		Title:    "Total of certificate retrievals",
		Units:    "retrievals/s",
		Fam:      "adcs",
		Ctx:      "wmi.adcs_cert_template_retrievals",
		Priority: prioADCSCertTemplateRetrievals,
		Dims: module.Dims{
			{ID: "adcs_cert_template_%s_retrievals_total", Name: "retrievals", Algo: module.Incremental},
		},
	}
	adcsCertificateRetrievalsTimeChartTmpl = module.Chart{
		ID:       "adcs_cert_template_%s_retrievals_processing_time",
		Title:    "Certificate last retrieval processing time",
		Units:    "seconds",
		Fam:      "adcs",
		Ctx:      "wmi.adcs_cert_template_retrieval_processing_time",
		Priority: prioADCSCertTemplateRetrievalProcessingTime,
		Dims: module.Dims{
			{ID: "adcs_cert_template_%s_retrievals_processing_time_seconds", Name: "processing_time", Div: precision},
		},
	}
	adcsCertTemplateRequestCryptoSigningTimeChartTmpl = module.Chart{
		ID:       "adcs_cert_template_%s_request_cryptographic_signing_time",
		Title:    "Certificate last signing operation request time",
		Units:    "seconds",
		Fam:      "adcs",
		Ctx:      "wmi.adcs_cert_template_request_cryptographic_signing_time",
		Priority: prioADCSCertTemplateRequestCryptoSigningTime,
		Dims: module.Dims{
			{ID: "adcs_cert_template_%s_request_cryptographic_signing_time_seconds", Name: "singing_time", Div: precision},
		},
	}
	adcsCertTemplateRequestPolicyModuleProcessingTimeChartTmpl = module.Chart{
		ID:       "adcs_cert_template_%s_request_policy_module_processing_time",
		Title:    "Certificate last policy module processing request time",
		Units:    "seconds",
		Fam:      "adcs",
		Ctx:      "wmi.adcs_cert_template_request_policy_module_processing",
		Priority: prioADCSCertTemplateRequestPolicyModuleProcessingTime,
		Dims: module.Dims{
			{ID: "adcs_cert_template_%s_request_policy_module_processing_time_seconds", Name: "processing_time", Div: precision},
		},
	}
	adcsCertTemplateChallengeResponseChartTmpl = module.Chart{
		ID:       "adcs_cert_template_%s_challenge_responses",
		Title:    "Certificate challenge responses",
		Units:    "responses/s",
		Fam:      "adcs",
		Ctx:      "wmi.adcs_cert_template_challenge_responses",
		Priority: prioADCSCertTemplateChallengeResponses,
		Dims: module.Dims{
			{ID: "adcs_cert_template_%s_challenge_responses_total", Name: "challenge", Algo: module.Incremental},
		},
	}
	adcsCertTemplateChallengeResponseProcessingTimeChartTmpl = module.Chart{
		ID:       "adcs_cert_template_%s_challenge_response_processing_time",
		Title:    "Certificate last challenge response time",
		Units:    "seconds",
		Fam:      "adcs",
		Ctx:      "wmi.adcs_cert_template_challenge_response_processing_time",
		Priority: prioADCSCertTemplateChallengeResponseProcessingTime,
		Dims: module.Dims{
			{ID: "adcs_cert_template_%s_challenge_response_processing_time_seconds", Name: "processing_time", Div: precision},
		},
	}
	adcsCertTemplateSignedCertificateTimestampListsChartTmpl = module.Chart{
		ID:       "adcs_cert_template_%s_signed_certificate_timestamp_lists",
		Title:    "Certificate Signed Certificate Timestamp Lists processed",
		Units:    "lists/s",
		Fam:      "adcs",
		Ctx:      "wmi.adcs_cert_template_signed_certificate_timestamp_lists",
		Priority: prioADCSCertTemplateSignedCertificateTimestampLists,
		Dims: module.Dims{
			{ID: "adcs_cert_template_%s_signed_certificate_timestamp_lists_total", Name: "processed", Algo: module.Incremental},
		},
	}
	adcsCertTemplateSignedCertificateTimestampListProcessingTimeChartTmpl = module.Chart{
		ID:       "adcs_cert_template_%s_signed_certificate_timestamp_list_processing_time",
		Title:    "Certificate last Signed Certificate Timestamp List process time",
		Units:    "seconds",
		Fam:      "adcs",
		Ctx:      "wmi.adcs_cert_template_signed_certificate_timestamp_list_processing_time",
		Priority: prioADCSCertTemplateSignedCertificateTimestampListProcessingTime,
		Dims: module.Dims{
			{ID: "adcs_cert_template_%s_signed_certificate_timestamp_list_processing_time_seconds", Name: "processing_time", Div: precision},
		},
	}
)

// AD FS
var (
	adfsCharts = module.Charts{
		adfsADLoginConnectionFailuresChart.Copy(),
		adfsCertificateAuthenticationsChart.Copy(),
		adfsDBArtifactFailuresChart.Copy(),
		adfsDBArtifactQueryTimeSecondsChart.Copy(),
		adfsDBConfigFailuresChart.Copy(),
		adfsDBConfigQueryTimeSecondsChart.Copy(),
		adfsDeviceAuthenticationsChart.Copy(),
		adfsExternalAuthenticationsChart.Copy(),
		adfsFederatedAuthenticationsChart.Copy(),
		adfsFederationMetadataRequestsChart.Copy(),

		adfsOAuthAuthorizationRequestsChart.Copy(),
		adfsOAuthClientAuthenticationsChart.Copy(),
		adfsOAuthClientCredentialRequestsChart.Copy(),
		adfsOAuthClientPrivKeyJwtAuthenticationsChart.Copy(),
		adfsOAuthClientSecretBasicAuthenticationsChart.Copy(),
		adfsOAuthClientSecretPostAuthenticationsChart.Copy(),
		adfsOAuthClientWindowsAuthenticationsChart.Copy(),
		adfsOAuthLogonCertificateRequestsChart.Copy(),
		adfsOAuthPasswordGrantRequestsChart.Copy(),
		adfsOAuthTokenRequestsChart.Copy(),

		adfsPassiveRequestsChart.Copy(),
		adfsPassportAuthenticationsChart.Copy(),
		adfsPasswordChangeChart.Copy(),
		adfsSAMLPTokenRequestsChart.Copy(),
		adfsSSOAuthenticationsChart.Copy(),
		adfsTokenRequestsChart.Copy(),
		adfsUserPasswordAuthenticationsChart.Copy(),
		adfsWindowsIntegratedAuthenticationsChart.Copy(),
		adfsWSFedTokenRequestsSuccessChart.Copy(),
		adfsWSTrustTokenRequestsSuccessChart.Copy(),
	}

	adfsADLoginConnectionFailuresChart = module.Chart{
		ID:       "adfs_ad_login_connection_failures",
		Title:    "Connection failures",
		Units:    "failures/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_ad_login_connection_failures",
		Priority: prioADFSADLoginConnectionFailures,
		Dims: module.Dims{
			{ID: "adfs_ad_login_connection_failures_total", Name: "connection", Algo: module.Incremental},
		},
	}
	adfsCertificateAuthenticationsChart = module.Chart{
		ID:       "adfs_certificate_authentications",
		Title:    "User Certificate authentications",
		Units:    "authentications/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_certificate_authentications",
		Priority: prioADFSCertificateAuthentications,
		Dims: module.Dims{
			{ID: "adfs_certificate_authentications_total", Name: "authentications", Algo: module.Incremental},
		},
	}

	adfsDBArtifactFailuresChart = module.Chart{
		ID:       "adfs_db_artifact_failures",
		Title:    "Connection failures to the artifact database",
		Units:    "failures/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_db_artifact_failures",
		Priority: prioADFSDBArtifactFailures,
		Dims: module.Dims{
			{ID: "adfs_db_artifact_failure_total", Name: "connection", Algo: module.Incremental},
		},
	}
	adfsDBArtifactQueryTimeSecondsChart = module.Chart{
		ID:       "adfs_db_artifact_query_time_seconds",
		Title:    "Time taken for an artifact database query",
		Units:    "seconds/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_db_artifact_query_time_seconds",
		Priority: prioADFSDBArtifactQueryTimeSeconds,
		Dims: module.Dims{
			{ID: "adfs_db_artifact_query_time_seconds_total", Name: "query_time", Algo: module.Incremental, Div: precision},
		},
	}
	adfsDBConfigFailuresChart = module.Chart{
		ID:       "adfs_db_config_failures",
		Title:    "Connection failures to the configuration database",
		Units:    "failures/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_db_config_failures",
		Priority: prioADFSDBConfigFailures,
		Dims: module.Dims{
			{ID: "adfs_db_config_failure_total", Name: "connection", Algo: module.Incremental},
		},
	}
	adfsDBConfigQueryTimeSecondsChart = module.Chart{
		ID:       "adfs_db_config_query_time_seconds",
		Title:    "Time taken for a configuration database query",
		Units:    "seconds/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_db_config_query_time_seconds",
		Priority: prioADFSDBConfigQueryTimeSeconds,
		Dims: module.Dims{
			{ID: "adfs_db_config_query_time_seconds_total", Name: "query_time", Algo: module.Incremental, Div: precision},
		},
	}
	adfsDeviceAuthenticationsChart = module.Chart{
		ID:       "adfs_device_authentications",
		Title:    "Device authentications",
		Units:    "authentications/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_device_authentications",
		Priority: prioADFSDeviceAuthentications,
		Dims: module.Dims{
			{ID: "adfs_device_authentications_total", Name: "authentications", Algo: module.Incremental},
		},
	}
	adfsExternalAuthenticationsChart = module.Chart{
		ID:       "adfs_external_authentications",
		Title:    "Authentications from external MFA providers",
		Units:    "authentications/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_external_authentications",
		Priority: prioADFSExternalAuthentications,
		Dims: module.Dims{
			{ID: "adfs_external_authentications_success_total", Name: "success", Algo: module.Incremental},
			{ID: "adfs_external_authentications_failure_total", Name: "failure", Algo: module.Incremental},
		},
	}
	adfsFederatedAuthenticationsChart = module.Chart{
		ID:       "adfs_federated_authentications",
		Title:    "Authentications from Federated Sources",
		Units:    "authentications/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_federated_authentications",
		Priority: prioADFSFederatedAuthentications,
		Dims: module.Dims{
			{ID: "adfs_federated_authentications_total", Name: "authentications", Algo: module.Incremental},
		},
	}
	adfsFederationMetadataRequestsChart = module.Chart{
		ID:       "adfs_federation_metadata_requests",
		Title:    "Federation Metadata requests",
		Units:    "requests/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_federation_metadata_requests",
		Priority: prioADFSFederationMetadataRequests,
		Dims: module.Dims{
			{ID: "adfs_federation_metadata_requests_total", Name: "requests", Algo: module.Incremental},
		},
	}

	adfsOAuthAuthorizationRequestsChart = module.Chart{
		ID:       "adfs_oauth_authorization_requests",
		Title:    "Incoming requests to the OAuth Authorization endpoint",
		Units:    "requests/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_oauth_authorization_requests",
		Priority: prioADFSOauthAuthorizationRequests,
		Dims: module.Dims{
			{ID: "adfs_oauth_authorization_requests_total", Name: "requests", Algo: module.Incremental},
		},
	}
	adfsOAuthClientAuthenticationsChart = module.Chart{
		ID:       "adfs_oauth_client_authentications",
		Title:    "OAuth client authentications",
		Units:    "authentications/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_oauth_client_authentications",
		Priority: prioADFSOauthClientAuthentications,
		Dims: module.Dims{
			{ID: "adfs_oauth_client_authentication_success_total", Name: "success", Algo: module.Incremental},
			{ID: "adfs_oauth_client_authentication_failure_total", Name: "failure", Algo: module.Incremental},
		},
	}
	adfsOAuthClientCredentialRequestsChart = module.Chart{
		ID:       "adfs_oauth_client_credentials_requests",
		Title:    "OAuth client credentials requests",
		Units:    "requests/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_oauth_client_credentials_requests",
		Priority: prioADFSOauthClientCredentials,
		Dims: module.Dims{
			{ID: "adfs_oauth_client_credentials_success_total", Name: "success", Algo: module.Incremental},
			{ID: "adfs_oauth_client_credentials_failure_total", Name: "failure", Algo: module.Incremental},
		},
	}
	adfsOAuthClientPrivKeyJwtAuthenticationsChart = module.Chart{
		ID:       "adfs_oauth_client_privkey_jwt_authentications",
		Title:    "OAuth client private key JWT authentications",
		Units:    "authentications/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_oauth_client_privkey_jwt_authentications",
		Priority: prioADFSOauthClientPrivkeyJwtAuthentication,
		Dims: module.Dims{
			{ID: "adfs_oauth_client_privkey_jwt_authentications_success_total", Name: "success", Algo: module.Incremental},
			{ID: "adfs_oauth_client_privkey_jtw_authentication_failure_total", Name: "failure", Algo: module.Incremental},
		},
	}
	adfsOAuthClientSecretBasicAuthenticationsChart = module.Chart{
		ID:       "adfs_oauth_client_secret_basic_authentications",
		Title:    "OAuth client secret basic authentications",
		Units:    "authentications/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_oauth_client_secret_basic_authentications",
		Type:     module.Line,
		Priority: prioADFSOauthClientSecretBasicAuthentications,
		Dims: module.Dims{
			{ID: "adfs_oauth_client_secret_basic_authentications_success_total", Name: "success", Algo: module.Incremental},
			{ID: "adfs_oauth_client_secret_basic_authentications_failure_total", Name: "failure", Algo: module.Incremental},
		},
	}
	adfsOAuthClientSecretPostAuthenticationsChart = module.Chart{
		ID:       "adfs_oauth_client_secret_post_authentications",
		Title:    "OAuth client secret post authentications",
		Units:    "authentications/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_oauth_client_secret_post_authentications",
		Priority: prioADFSOauthClientSecretPostAuthentications,
		Dims: module.Dims{
			{ID: "adfs_oauth_client_secret_post_authentications_success_total", Name: "success", Algo: module.Incremental},
			{ID: "adfs_oauth_client_secret_post_authentications_failure_total", Name: "failure", Algo: module.Incremental},
		},
	}
	adfsOAuthClientWindowsAuthenticationsChart = module.Chart{
		ID:       "adfs_oauth_client_windows_authentications",
		Title:    "OAuth client windows integrated authentications",
		Units:    "authentications/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_oauth_client_windows_authentications",
		Priority: prioADFSOauthClientWindowsAuthentications,
		Dims: module.Dims{
			{ID: "adfs_oauth_client_windows_authentications_success_total", Name: "success", Algo: module.Incremental},
			{ID: "adfs_oauth_client_windows_authentications_failure_total", Name: "failure", Algo: module.Incremental},
		},
	}
	adfsOAuthLogonCertificateRequestsChart = module.Chart{
		ID:       "adfs_oauth_logon_certificate_requests",
		Title:    "OAuth logon certificate requests",
		Units:    "requests/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_oauth_logon_certificate_requests",
		Priority: prioADFSOauthLogonCertificateRequests,
		Dims: module.Dims{
			{ID: "adfs_oauth_logon_certificate_token_requests_success_total", Name: "success", Algo: module.Incremental},
			{ID: "adfs_oauth_logon_certificate_requests_failure_total", Name: "failure", Algo: module.Incremental},
		},
	}
	adfsOAuthPasswordGrantRequestsChart = module.Chart{
		ID:       "adfs_oauth_password_grant_requests",
		Title:    "OAuth password grant requests",
		Units:    "requests/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_oauth_password_grant_requests",
		Priority: prioADFSOauthPasswordGrantRequests,
		Dims: module.Dims{
			{ID: "adfs_oauth_password_grant_requests_success_total", Name: "success", Algo: module.Incremental},
			{ID: "adfs_oauth_password_grant_requests_failure_total", Name: "failure", Algo: module.Incremental},
		},
	}
	adfsOAuthTokenRequestsChart = module.Chart{
		ID:       "adfs_oauth_token_requests_success",
		Title:    "Successful RP token requests over OAuth protocol",
		Units:    "requests/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_oauth_token_requests_success",
		Priority: prioADFSOauthTokenRequestsSuccess,
		Dims: module.Dims{
			{ID: "adfs_oauth_token_requests_success_total", Name: "success", Algo: module.Incremental},
		},
	}

	adfsPassiveRequestsChart = module.Chart{
		ID:       "adfs_passive_requests",
		Title:    "Passive requests",
		Units:    "requests/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_passive_requests",
		Priority: prioADFSPassiveRequests,
		Dims: module.Dims{
			{ID: "adfs_passive_requests_total", Name: "passive", Algo: module.Incremental},
		},
	}
	adfsPassportAuthenticationsChart = module.Chart{
		ID:       "adfs_passport_authentications",
		Title:    "Microsoft Passport SSO authentications",
		Units:    "authentications/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_passport_authentications",
		Priority: prioADFSPassportAuthentications,
		Dims: module.Dims{
			{ID: "adfs_passport_authentications_total", Name: "passport", Algo: module.Incremental},
		},
	}
	adfsPasswordChangeChart = module.Chart{
		ID:       "adfs_password_change_requests",
		Title:    "Password change requests",
		Units:    "requests/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_password_change_requests",
		Priority: prioADFSPasswordChangeRequests,
		Dims: module.Dims{
			{ID: "adfs_password_change_succeeded_total", Name: "success", Algo: module.Incremental},
			{ID: "adfs_password_change_failed_total", Name: "failed", Algo: module.Incremental},
		},
	}
	adfsSAMLPTokenRequestsChart = module.Chart{
		ID:       "adfs_samlp_token_requests_success",
		Title:    "Successful RP token requests over SAML-P protocol",
		Units:    "requests/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_samlp_token_requests_success",
		Priority: prioADFSSAMLPTokenRequests,
		Dims: module.Dims{
			{ID: "adfs_samlp_token_requests_success_total", Name: "success", Algo: module.Incremental},
		},
	}
	adfsSSOAuthenticationsChart = module.Chart{
		ID:       "adfs_sso_authentications",
		Title:    "SSO authentications",
		Units:    "authentications/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_sso_authentications",
		Priority: prioADFSSSOAuthentications,
		Dims: module.Dims{
			{ID: "adfs_sso_authentications_success_total", Name: "success", Algo: module.Incremental},
			{ID: "adfs_sso_authentications_failure_total", Name: "failure", Algo: module.Incremental},
		},
	}
	adfsTokenRequestsChart = module.Chart{
		ID:       "adfs_token_requests",
		Title:    "Token access requests",
		Units:    "requests/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_token_requests",
		Priority: prioADFSTokenRequests,
		Dims: module.Dims{
			{ID: "adfs_token_requests_total", Name: "requests", Algo: module.Incremental},
		},
	}
	adfsUserPasswordAuthenticationsChart = module.Chart{
		ID:       "adfs_userpassword_authentications",
		Title:    "AD U/P authentications",
		Units:    "authentications/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_userpassword_authentications",
		Priority: prioADFSUserPasswordAuthentications,
		Dims: module.Dims{
			{ID: "adfs_sso_authentications_success_total", Name: "success", Algo: module.Incremental},
			{ID: "adfs_sso_authentications_failure_total", Name: "failure", Algo: module.Incremental},
		},
	}
	adfsWindowsIntegratedAuthenticationsChart = module.Chart{
		ID:       "adfs_windows_integrated_authentications",
		Title:    "f Windows integrated authentications using Kerberos or NTLM",
		Units:    "authentications/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_windows_integrated_authentications",
		Priority: prioADFSWindowsIntegratedAuthentications,
		Dims: module.Dims{
			{ID: "adfs_windows_integrated_authentications_total", Name: "authentications", Algo: module.Incremental},
		},
	}
	adfsWSFedTokenRequestsSuccessChart = module.Chart{
		ID:       "adfs_wsfed_token_requests_success",
		Title:    "Successful RP token requests over WS-Fed protocol",
		Units:    "requests/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_wsfed_token_requests_success",
		Priority: prioADFSWSFedTokenRequestsSuccess,
		Dims: module.Dims{
			{ID: "adfs_wsfed_token_requests_success_total", Name: "success", Algo: module.Incremental},
		},
	}
	adfsWSTrustTokenRequestsSuccessChart = module.Chart{
		ID:       "adfs_wstrust_token_requests_success",
		Title:    "Successful RP token requests over WS-Trust protocol",
		Units:    "requests/s",
		Fam:      "adfs",
		Ctx:      "wmi.adfs_wstrust_token_requests_success",
		Priority: prioADFSWSTrustTokenRequestsSuccess,
		Dims: module.Dims{
			{ID: "adfs_wstrust_token_requests_success_total", Name: "success", Algo: module.Incremental},
		},
	}
)

// Logon
var (
	logonCharts = module.Charts{
		logonSessionsChart.Copy(),
	}
	logonSessionsChart = module.Chart{
		ID:       "logon_active_sessions_by_type",
		Title:    "Active User Logon Sessions By Type",
		Units:    "sessions",
		Fam:      "logon",
		Ctx:      "wmi.logon_type_sessions",
		Type:     module.Stacked,
		Priority: prioLogonSessions,
		Dims: module.Dims{
			{ID: "logon_type_system_sessions", Name: "system"},
			{ID: "logon_type_proxy_sessions", Name: "proxy"},
			{ID: "logon_type_network_sessions", Name: "network"},
			{ID: "logon_type_interactive_sessions", Name: "interactive"},
			{ID: "logon_type_batch_sessions", Name: "batch"},
			{ID: "logon_type_service_sessions", Name: "service"},
			{ID: "logon_type_unlock_sessions", Name: "unlock"},
			{ID: "logon_type_network_clear_text_sessions", Name: "network_clear_text"},
			{ID: "logon_type_new_credentials_sessions", Name: "new_credentials"},
			{ID: "logon_type_remote_interactive_sessions", Name: "remote_interactive"},
			{ID: "logon_type_cached_interactive_sessions", Name: "cached_interactive"},
			{ID: "logon_type_cached_remote_interactive_sessions", Name: "cached_remote_interactive"},
			{ID: "logon_type_cached_unlock_sessions", Name: "cached_unlock"},
		},
	}
)

// Thermal zone
var (
	thermalzoneChartsTmpl = module.Charts{
		thermalzoneTemperatureChartTmpl.Copy(),
	}
	thermalzoneTemperatureChartTmpl = module.Chart{
		ID:       "thermalzone_%s_temperature",
		Title:    "Thermal zone temperature",
		Units:    "celsius",
		Fam:      "thermalzone",
		Ctx:      "wmi.thermalzone_temperature",
		Priority: prioThermalzoneTemperature,
		Dims: module.Dims{
			{ID: "thermalzone_%s_temperature", Name: "temperature"},
		},
	}
)

// Processes
var (
	processesCharts = module.Charts{
		processesCPUUtilizationTotalChart.Copy(),
		processesMemoryUsageChart.Copy(),
		processesHandlesChart.Copy(),
		processesIOBytesChart.Copy(),
		processesIOOperationsChart.Copy(),
		processesPageFaultsChart.Copy(),
		processesPageFileBytes.Copy(),
		processesThreads.Copy(),
	}
	processesCPUUtilizationTotalChart = module.Chart{
		ID:       "processes_cpu_utilization",
		Title:    "CPU usage (100% = 1 core)",
		Units:    "percentage",
		Fam:      "processes",
		Ctx:      "wmi.processes_cpu_utilization",
		Type:     module.Stacked,
		Priority: prioProcessesCPUUtilization,
	}
	processesMemoryUsageChart = module.Chart{
		ID:       "processes_memory_usage",
		Title:    "Memory usage",
		Units:    "bytes",
		Fam:      "processes",
		Ctx:      "wmi.processes_memory_usage",
		Type:     module.Stacked,
		Priority: prioProcessesMemoryUsage,
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
	processesThreads = module.Chart{
		ID:       "processes_threads",
		Title:    "Active threads",
		Units:    "threads",
		Fam:      "processes",
		Ctx:      "wmi.processes_threads",
		Type:     module.Stacked,
		Priority: prioProcessesThreads,
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
)

// Service
var (
	serviceChartsTmpl = module.Charts{
		serviceStateChartTmpl.Copy(),
		serviceStatusChartTmpl.Copy(),
	}
	serviceStateChartTmpl = module.Chart{
		ID:       "service_%s_state",
		Title:    "Service state",
		Units:    "state",
		Fam:      "services",
		Ctx:      "wmi.service_state",
		Priority: prioServiceState,
		Dims: module.Dims{
			{ID: "service_%s_state_running", Name: "running"},
			{ID: "service_%s_state_stopped", Name: "stopped"},
			{ID: "service_%s_state_start_pending", Name: "start_pending"},
			{ID: "service_%s_state_stop_pending", Name: "stop_pending"},
			{ID: "service_%s_state_continue_pending", Name: "continue_pending"},
			{ID: "service_%s_state_pause_pending", Name: "pause_pending"},
			{ID: "service_%s_state_paused", Name: "paused"},
			{ID: "service_%s_state_unknown", Name: "unknown"},
		},
	}
	serviceStatusChartTmpl = module.Chart{
		ID:       "service_%s_status",
		Title:    "Service status",
		Units:    "status",
		Fam:      "services",
		Ctx:      "wmi.service_status",
		Priority: prioServiceStatus,
		Dims: module.Dims{
			{ID: "service_%s_status_ok", Name: "ok"},
			{ID: "service_%s_status_error", Name: "error"},
			{ID: "service_%s_status_unknown", Name: "unknown"},
			{ID: "service_%s_status_degraded", Name: "degraded"},
			{ID: "service_%s_status_pred_fail", Name: "pred_fail"},
			{ID: "service_%s_status_starting", Name: "starting"},
			{ID: "service_%s_status_stopping", Name: "stopping"},
			{ID: "service_%s_status_service", Name: "service"},
			{ID: "service_%s_status_stressed", Name: "stressed"},
			{ID: "service_%s_status_nonrecover", Name: "nonrecover"},
			{ID: "service_%s_status_no_contact", Name: "no_contact"},
			{ID: "service_%s_status_lost_comm", Name: "lost_comm"},
		},
	}
)

// Collectors
var (
	collectorChartsTmpl = module.Charts{
		collectorDurationChartTmpl.Copy(),
		collectorStatusChartTmpl.Copy(),
	}
	collectorDurationChartTmpl = module.Chart{
		ID:       "collector_%s_duration",
		Title:    "Duration of a data collection",
		Units:    "seconds",
		Fam:      "collection",
		Ctx:      "wmi.collector_duration",
		Priority: prioCollectorDuration,
		Dims: module.Dims{
			{ID: "collector_%s_duration", Name: "duration", Div: precision},
		},
	}
	collectorStatusChartTmpl = module.Chart{
		ID:       "collector_%s_status",
		Title:    "Status of a data collection",
		Units:    "status",
		Fam:      "collection",
		Ctx:      "wmi.collector_status",
		Priority: prioCollectorStatus,
		Dims: module.Dims{
			{ID: "collector_%s_status_success", Name: "success"},
			{ID: "collector_%s_status_fail", Name: "fail"},
		},
	}
)

func (w *WMI) addCPUCharts() {
	charts := cpuCharts.Copy()

	if err := w.Charts().Add(*charts...); err != nil {
		w.Warning(err)
	}
}

func (w *WMI) addCPUCoreCharts(core string) {
	charts := cpuCoreChartsTmpl.Copy()

	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, core)
		chart.Labels = []module.Label{
			{Key: "core", Value: core},
		}
		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, core)
		}
	}

	if err := w.Charts().Add(*charts...); err != nil {
		w.Warning(err)
	}
}

func (w *WMI) removeCPUCoreCharts(core string) {
	px := fmt.Sprintf("cpu_core_%s", core)
	for _, chart := range *w.Charts() {
		if strings.HasPrefix(chart.ID, px) {
			chart.MarkRemove()
			chart.MarkNotCreated()
		}
	}
}

func (w *WMI) addMemoryCharts() {
	charts := memCharts.Copy()

	if err := w.Charts().Add(*charts...); err != nil {
		w.Warning(err)
	}
}

func (w *WMI) addDiskCharts(disk string) {
	charts := diskChartsTmpl.Copy()

	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, disk)
		chart.Labels = []module.Label{
			{Key: "disk", Value: disk},
		}
		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, disk)
		}
	}

	if err := w.Charts().Add(*charts...); err != nil {
		w.Warning(err)
	}
}

func (w *WMI) removeDiskCharts(disk string) {
	px := fmt.Sprintf("logical_disk_%s", disk)
	for _, chart := range *w.Charts() {
		if strings.HasPrefix(chart.ID, px) {
			chart.MarkRemove()
			chart.MarkNotCreated()
		}
	}
}

func (w *WMI) addNICCharts(nic string) {
	charts := nicChartsTmpl.Copy()

	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, nic)
		chart.Labels = []module.Label{
			{Key: "nic", Value: nic},
		}
		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, nic)
		}
	}

	if err := w.Charts().Add(*charts...); err != nil {
		w.Warning(err)
	}
}

func (w *WMI) removeNICCharts(nic string) {
	px := fmt.Sprintf("nic_%s", nic)
	for _, chart := range *w.Charts() {
		if strings.HasPrefix(chart.ID, px) {
			chart.MarkRemove()
			chart.MarkNotCreated()
		}
	}
}

func (w *WMI) addTCPCharts() {
	charts := tcpCharts.Copy()

	if err := w.Charts().Add(*charts...); err != nil {
		w.Warning(err)
	}
}

func (w *WMI) addOSCharts() {
	charts := osCharts.Copy()

	if err := w.Charts().Add(*charts...); err != nil {
		w.Warning(err)
	}
}

func (w *WMI) addSystemCharts() {
	charts := systemCharts.Copy()

	if err := w.Charts().Add(*charts...); err != nil {
		w.Warning(err)
	}
}

func (w *WMI) addLogonCharts() {
	charts := logonCharts.Copy()

	if err := w.Charts().Add(*charts...); err != nil {
		w.Warning(err)
	}
}

func (w *WMI) addADFSCharts() {
	charts := adfsCharts.Copy()

	if err := w.Charts().Add(*charts...); err != nil {
		w.Warning(err)
	}
}

func (w *WMI) addThermalZoneCharts(zone string) {
	charts := thermalzoneChartsTmpl.Copy()

	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, zone)
		chart.Labels = []module.Label{
			{Key: "thermalzone", Value: zone},
		}
		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, zone)
		}
	}

	if err := w.Charts().Add(*charts...); err != nil {
		w.Warning(err)
	}
}

func (w *WMI) removeThermalZoneCharts(zone string) {
	px := fmt.Sprintf("thermalzone_%s", zone)
	for _, chart := range *w.Charts() {
		if strings.HasPrefix(chart.ID, px) {
			chart.MarkRemove()
			chart.MarkNotCreated()
		}
	}
}

func (w *WMI) addIISWebsiteCharts(website string) {
	charts := iisWebsiteChartsTmpl.Copy()

	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, website)
		chart.Labels = []module.Label{
			{Key: "website", Value: website},
		}
		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, website)
		}
	}

	if err := w.Charts().Add(*charts...); err != nil {
		w.Warning(err)
	}
}

func (w *WMI) removeIIWebsiteSCharts(website string) {
	px := fmt.Sprintf("iis_website_%s", website)
	for _, chart := range *w.Charts() {
		if strings.HasPrefix(chart.ID, px) {
			chart.MarkRemove()
			chart.MarkNotCreated()
		}
	}
}

func (w *WMI) addMSSQLDBCharts(instance string, dbname string) {
	charts := mssqlDatabaseChartsTmpl.Copy()

	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, dbname, instance)
		chart.Labels = []module.Label{
			{Key: "mssql_instance", Value: instance},
			{Key: "database", Value: dbname},
		}
		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, dbname, instance)
		}
	}

	if err := w.Charts().Add(*charts...); err != nil {
		w.Warning(err)
	}
}

func (w *WMI) removeMSSQLDBCharts(instance string, dbname string) {
	px := fmt.Sprintf("mssql_db_%s_instance_%s", dbname, instance)
	for _, chart := range *w.Charts() {
		if strings.HasPrefix(chart.ID, px) {
			chart.MarkRemove()
			chart.MarkNotCreated()
		}
	}
}

func (w *WMI) addMSSQLInstanceCharts(instance string) {
	charts := mssqlInstanceChartsTmpl.Copy()

	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, instance)
		chart.Labels = []module.Label{
			{Key: "mssql_instance", Value: instance},
		}
		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, instance)
		}
	}

	if err := w.Charts().Add(*charts...); err != nil {
		w.Warning(err)
	}
}

func (w *WMI) removeMSSQLInstanceCharts(instance string) {
	px := fmt.Sprintf("mssql_instance_%s", instance)
	for _, chart := range *w.Charts() {
		if strings.HasPrefix(chart.ID, px) {
			chart.MarkRemove()
			chart.MarkNotCreated()
		}
	}
}

func (w *WMI) addProcessesCharts() {
	charts := processesCharts.Copy()

	if err := w.Charts().Add(*charts...); err != nil {
		w.Warning(err)
	}
}

func (w *WMI) addADCharts() {
	charts := adCharts.Copy()

	if err := w.Charts().Add(*charts...); err != nil {
		w.Warning(err)
	}
}

func (w *WMI) addCertificateTemplateCharts(template string) {
	charts := adcsCertTemplateChartsTmpl.Copy()

	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, template)
		chart.Labels = []module.Label{
			{Key: "cert_template", Value: template},
		}
		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, template)
		}
	}

	if err := w.Charts().Add(*charts...); err != nil {
		w.Warning(err)
	}
}

func (w *WMI) removeCertificateTemplateCharts(template string) {
	px := fmt.Sprintf("adcs_cert_template_%s", template)
	for _, chart := range *w.Charts() {
		if strings.HasPrefix(chart.ID, px) {
			chart.MarkRemove()
			chart.MarkNotCreated()
		}
	}
}

func (w *WMI) addProcessToCharts(procID string) {
	for _, chart := range *w.Charts() {
		var dim *module.Dim
		switch chart.ID {
		case processesCPUUtilizationTotalChart.ID:
			id := fmt.Sprintf("process_%s_cpu_time", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental, Div: 1000, Mul: 100}
			if procID == "Idle" {
				dim.Hidden = true
			}
		case processesMemoryUsageChart.ID:
			id := fmt.Sprintf("process_%s_working_set_private_bytes", procID)
			dim = &module.Dim{ID: id, Name: procID}
		case processesIOBytesChart.ID:
			id := fmt.Sprintf("process_%s_io_bytes", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case processesIOOperationsChart.ID:
			id := fmt.Sprintf("process_%s_io_operations", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case processesPageFaultsChart.ID:
			id := fmt.Sprintf("process_%s_page_faults", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case processesPageFileBytes.ID:
			id := fmt.Sprintf("process_%s_page_file_bytes", procID)
			dim = &module.Dim{ID: id, Name: procID}
		case processesThreads.ID:
			id := fmt.Sprintf("process_%s_threads", procID)
			dim = &module.Dim{ID: id, Name: procID}
		case processesHandlesChart.ID:
			id := fmt.Sprintf("process_%s_handles", procID)
			dim = &module.Dim{ID: id, Name: procID}
		default:
			continue
		}

		if dim == nil {
			continue
		}
		if err := chart.AddDim(dim); err != nil {
			w.Warning(err)
			continue
		}
		chart.MarkNotCreated()
	}
}

func (w *WMI) removeProcessFromCharts(procID string) {
	for _, chart := range *w.Charts() {
		var id string
		switch chart.ID {
		case processesCPUUtilizationTotalChart.ID:
			id = fmt.Sprintf("process_%s_cpu_time", procID)
		case processesMemoryUsageChart.ID:
			id = fmt.Sprintf("process_%s_working_set_private_bytes", procID)
		case processesIOBytesChart.ID:
			id = fmt.Sprintf("process_%s_io_bytes", procID)
		case processesIOOperationsChart.ID:
			id = fmt.Sprintf("process_%s_io_operations", procID)
		case processesPageFaultsChart.ID:
			id = fmt.Sprintf("process_%s_page_faults", procID)
		case processesPageFileBytes.ID:
			id = fmt.Sprintf("process_%s_page_file_bytes", procID)
		case processesThreads.ID:
			id = fmt.Sprintf("process_%s_threads", procID)
		case processesHandlesChart.ID:
			id = fmt.Sprintf("process_%s_handles", procID)
		default:
			continue
		}

		if err := chart.MarkDimRemove(id, false); err != nil {
			w.Warning(err)
			continue
		}
		chart.MarkNotCreated()
	}
}

func (w *WMI) addServiceCharts(svc string) {
	charts := serviceChartsTmpl.Copy()

	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, svc)
		chart.Labels = []module.Label{
			{Key: "service", Value: svc},
		}
		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, svc)
		}
	}

	if err := w.Charts().Add(*charts...); err != nil {
		w.Warning(err)
	}
}

func (w *WMI) removeServiceCharts(svc string) {
	px := fmt.Sprintf("service_%s", svc)
	for _, chart := range *w.Charts() {
		if strings.HasPrefix(chart.ID, px) {
			chart.MarkRemove()
			chart.MarkNotCreated()
		}
	}
}

func (w *WMI) addCollectorCharts(name string) {
	charts := collectorChartsTmpl.Copy()

	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, name)
		chart.Labels = []module.Label{
			{Key: "collector", Value: name},
		}
		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, name)
		}
	}

	if err := w.Charts().Add(*charts...); err != nil {
		w.Warning(err)
	}
}

func (w *WMI) removeCollectorCharts(name string) {
	px := fmt.Sprintf("collector_%s", name)
	for _, chart := range *w.Charts() {
		if strings.HasPrefix(chart.ID, px) {
			chart.MarkRemove()
			chart.MarkNotCreated()
		}
	}
}

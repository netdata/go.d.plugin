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
	prioIISWebsiteFTPFileTransferRate
	prioIISWebsiteActiveConnectionsCount
	prioIISWebsiteRequestsRate
	prioIISWebsiteConnectionAttemptsRate
	prioIISWebsiteUsersCount
	prioIISWebsiteISAPIExtRequestsCount
	prioIISWebsiteISAPIExtRequestsRate
	prioIISWebsiteErrorsRate
	prioIISWebsiteLogonAttemptsRate
	prioIISWebsiteUptime

	// Connections
	prioMSSQLUserConnections

	// Transactions
	prioMSSQLDatabaseTransactions
	prioMSSQLDatabaseActiveTransactions
	prioMSSQLDatabaseWriteTransactions
	prioMSSQLDatabaseBackupRestoreOperations
	prioMSSQLDatabaseLogFlushes
	prioMSSQLDatabaseLogFlushed

	// Size
	prioMSSQLDatabaseDataFileSize

	// SQL activity
	prioMSSQLStatsBatchRequests
	prioMSSQLStatsCompilations
	prioMSSQLStatsRecompilations
	prioMSSQLStatsAutoParameterization
	prioMSSQLStatsSafeAutoParameterization

	// Processes
	prioMSSQLBlockedProcess

	// Buffer Cache
	prioMSSQLCacheHitRatio
	prioMSSQLBufManIOPS
	prioMSSQLBufferCheckpointPages
	prioMSSQLAccessMethodPageSplits
	prioMSSQLBufferPageLifeExpectancy

	// Memory
	prioMSSQLMemmgrConnectionMemoryBytes
	prioMSSQLMemTotalServer
	prioMSSQLMemmgrExternalBenefitOfMemory
	prioMSSQLMemmgrPendingMemoryGrants

	// Locks
	prioMSSQLLocksLockWait
	prioMSSQLLocksDeadLocks

	// Error
	prioMSSQLSqlErrorsTotal

	// SQL stats
	prioMSSQLStatsAutoParameterization
	prioMSSQLStatsBatchRequests
	prioMSSQLStatsSafeAutoParameterization
	prioMSSQLStatsCompilations
	prioMSSQLStatsRecompilations

	prioNETFrameworkCLRExceptionsThrown
	prioNETFrameworkCLRExceptionsFilters
	prioNETFrameworkCLRExceptionsFinallys
	prioNETFrameworkCLRExceptionsThrowCatchDepth
	prioNETFrameworkCLRInteropCOMCalableWrapper
	prioNETFrameworkCLRInteropMarshalling
	prioNETFrameworkCLRInteropStubsCreated
	prioNETFrameworkCLRJITMethods
	prioNETFrameworkCLRJITTime
	prioNETFrameworkCLRJITStandardFailures
	prioNETFrameworkCLRJITILBytes
	prioNETFrameworkCLRLoadingLoaderHeapSize
	prioNETFrameworkCLRLoadingAppDomainsLoaded
	prioNETFrameworkCLRLoadingAppDomainsUnloaded
	prioNETFrameworkCLRLoadingAssembliesLoaded
	prioNETFrameworkCLRLoadingClassesLoaded
	prioNETFrameworkCLRLoadingClassLoadFailure
	prioNETFrameworkCLRLockAndThreadsQueueLength
	prioNETFrameworkCLRLockAndThreadsCurrentLogicalThreads
	prioNETFrameworkCLRLockAndThreadsCurrentPhysicalThreads
	prioNETFrameworkCLRLockAndThreadsRecognizedThreads
	prioNETFrameworkCLRLockAndThreadsContentions
	prioNETFrameworkCLRMemoryAllocatedBytes
	prioNETFrameworkCLRMemoryFinalizationSurvivors
	prioNETFrameworkCLRMemoryHeapSize
	prioNETFrameworkCLRMemoryPromoted
	prioNETFrameworkCLRMemoryNumberGCHandles
	prioNETFrameworkCLRMemoryCollections
	prioNETFrameworkCLRMemoryInducedGC
	prioNETFrameworkCLRMemoryNumberPinnedObjects
	prioNETFrameworkCLRMemoryNumberSinkBlocksInUse
	prioNETFrameworkCLRMemoryCommitted
	prioNETFrameworkCLRMemoryReserved
	prioNETFrameworkCLRMemoryGCTime
	prioNETFrameworkCLRRemotingChannels
	prioNETFrameworkCLRRemotingContextBoundClassesLoaded
	prioNETFrameworkCLRRemotingContextBoundObjects
	prioNETFrameworkCLRRemotingContextProxies
	prioNETFrameworkCLRRemotingContexts
	prioNETFrameworkCLRRemotingRemoteCalls
	prioNETFrameworkCLRSecurityLinkTimeChecks
	prioNETFrameworkCLRSecurityRTChecksTime
	prioNETFrameworkCLRSecurityStackWalkDepth
	prioNETFrameworkCLRSecurityRuntimeChecks

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

	// Requests
	prioADCSCertTemplateRequests
	prioADCSCertTemplateRequestProcessingTime
	prioADCSCertTemplateRetrievals
	prioADCSCertTemplateFailedRequests
	prioADCSCertTemplateIssuesRequests
	prioADCSCertTemplatePendingRequests

	// Response
	prioADCSCertTemplateChallengeResponses

	// Retrieval
	prioADCSCertTemplateRetrievalProcessingTime

	// Timing
	prioADCSCertTemplateRequestCryptoSigningTime
	prioADCSCertTemplateRequestPolicyModuleProcessingTime
	prioADCSCertTemplateChallengeResponseProcessingTime
	prioADCSCertTemplateSignedCertificateTimestampLists
	prioADCSCertTemplateSignedCertificateTimestampListProcessingTime

	// ADFS
	// AD
	prioADFSADLoginConnectionFailures

	// DB Artifacts
	prioADFSDBArtifactFailures
	prioADFSDBArtifactQueryTimeSeconds

	// DB Config
	prioADFSDBConfigFailures
	prioADFSDBConfigQueryTimeSeconds

	// Auth
	prioADFSDeviceAuthentications
	prioADFSExternalAuthentications
	prioADFSOauthAuthorizationRequests
	prioADFSCertificateAuthentications
	prioADFSOauthClientAuthentications
	prioADFSPassportAuthentications
	prioADFSSSOAuthentications
	prioADFSUserPasswordAuthentications
	prioADFSWindowsIntegratedAuthentications

	// OAuth
	prioADFSOauthClientCredentials
	prioADFSOauthClientPrivkeyJwtAuthentication
	prioADFSOauthClientSecretBasicAuthentications
	prioADFSOauthClientSecretPostAuthentications
	prioADFSOauthClientWindowsAuthentications
	prioADFSOauthLogonCertificateRequests
	prioADFSOauthPasswordGrantRequests
	prioADFSOauthTokenRequestsSuccess
	prioADFSFederatedAuthentications

	// Requests
	prioADFSFederationMetadataRequests
	prioADFSPassiveRequests
	prioADFSPasswordChangeRequests
	prioADFSSAMLPTokenRequests
	prioADFSWSTrustTokenRequestsSuccess
	prioADFSTokenRequests
	prioADFSWSFedTokenRequestsSuccess

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
		OverModule: "iis",
		ID:         "iis_website_%s_traffic",
		Title:      "Website traffic",
		Units:      "bytes/s",
		Fam:        "traffic",
		Ctx:        "iis.website_traffic",
		Type:       module.Area,
		Priority:   prioIISWebsiteTraffic,
		Dims: module.Dims{
			{ID: "iis_website_%s_received_bytes_total", Name: "received", Algo: module.Incremental},
			{ID: "iis_website_%s_sent_bytes_total", Name: "sent", Algo: module.Incremental, Mul: -1},
		},
	}
	iisWebsiteFTPFileTransferRateChartTempl = module.Chart{
		OverModule: "iis",
		ID:         "iis_website_%s_ftp_file_transfer_rate",
		Title:      "Website FTP file transfer rate",
		Units:      "files/s",
		Fam:        "traffic",
		Ctx:        "iis.website_ftp_file_transfer_rate",
		Priority:   prioIISWebsiteFTPFileTransferRate,
		Dims: module.Dims{
			{ID: "iis_website_%s_files_received_total", Name: "received", Algo: module.Incremental},
			{ID: "iis_website_%s_files_sent_total", Name: "sent", Algo: module.Incremental},
		},
	}
	iisWebsiteActiveConnectionsCountChartTmpl = module.Chart{
		OverModule: "iis",
		ID:         "iis_website_%s_active_connections_count",
		Title:      "Website active connections",
		Units:      "connections",
		Fam:        "connections",
		Ctx:        "iis.website_active_connections_count",
		Priority:   prioIISWebsiteActiveConnectionsCount,
		Dims: module.Dims{
			{ID: "iis_website_%s_current_connections", Name: "active"},
		},
	}
	iisWebsiteConnectionAttemptsRate = module.Chart{
		OverModule: "iis",
		ID:         "iis_website_%s_connection_attempts_rate",
		Title:      "Website connections attempts",
		Units:      "attempts/s",
		Fam:        "connections",
		Ctx:        "iis.website_connection_attempts_rate",
		Priority:   prioIISWebsiteConnectionAttemptsRate,
		Dims: module.Dims{
			{ID: "iis_website_%s_connection_attempts_all_instances_total", Name: "connection", Algo: module.Incremental},
		},
	}
	iisWebsiteRequestsRateChartTmpl = module.Chart{
		OverModule: "iis",
		ID:         "iis_website_%s_requests_rate",
		Title:      "Website requests rate",
		Units:      "requests/s",
		Fam:        "requests",
		Ctx:        "iis.website_requests_rate",
		Priority:   prioIISWebsiteRequestsRate,
		Dims: module.Dims{
			{ID: "iis_website_%s_requests_total", Name: "requests", Algo: module.Incremental},
		},
	}
	iisWebsiteUsersCountChartTmpl = module.Chart{
		OverModule: "iis",
		ID:         "iis_website_%s_users_count",
		Title:      "Website users with pending requests",
		Units:      "users",
		Fam:        "requests",
		Ctx:        "iis.website_users_count",
		Type:       module.Stacked,
		Priority:   prioIISWebsiteUsersCount,
		Dims: module.Dims{
			{ID: "iis_website_%s_current_anonymous_users", Name: "anonymous"},
			{ID: "iis_website_%s_current_non_anonymous_users", Name: "non_anonymous"},
		},
	}
	iisWebsiteISAPIExtRequestsCountChartTmpl = module.Chart{
		OverModule: "iis",
		ID:         "iis_website_%s_isapi_extension_requests_count",
		Title:      "ISAPI extension requests",
		Units:      "requests",
		Fam:        "requests",
		Ctx:        "iis.website_isapi_extension_requests_count",
		Priority:   prioIISWebsiteISAPIExtRequestsCount,
		Dims: module.Dims{
			{ID: "iis_website_%s_current_isapi_extension_requests", Name: "isapi"},
		},
	}
	iisWebsiteISAPIExtRequestsRateChartTmpl = module.Chart{
		OverModule: "iis",
		ID:         "iis_website_%s_isapi_extension_requests_rate",
		Title:      "Website extensions request",
		Units:      "requests/s",
		Fam:        "requests",
		Ctx:        "iis.website_isapi_extension_requests_rate",
		Priority:   prioIISWebsiteISAPIExtRequestsRate,
		Dims: module.Dims{
			{ID: "iis_website_%s_isapi_extension_requests_total", Name: "isapi", Algo: module.Incremental},
		},
	}
	iisWebsiteErrorsRateChart = module.Chart{
		OverModule: "iis",
		ID:         "iis_website_%s_errors_rate",
		Title:      "Website errors",
		Units:      "errors/s",
		Fam:        "requests",
		Ctx:        "iis.website_errors_rate",
		Type:       module.Stacked,
		Priority:   prioIISWebsiteErrorsRate,
		Dims: module.Dims{
			{ID: "iis_website_%s_locked_errors_total", Name: "document_locked", Algo: module.Incremental},
			{ID: "iis_website_%s_not_found_errors_total", Name: "document_not_found", Algo: module.Incremental},
		},
	}
	iisWebsiteLogonAttemptsRateChartTmpl = module.Chart{
		OverModule: "iis",
		ID:         "iis_website_%s_logon_attempts_rate",
		Title:      "Website logon attempts",
		Units:      "attempts/s",
		Fam:        "logon",
		Ctx:        "iis.website_logon_attempts_rate",
		Priority:   prioIISWebsiteLogonAttemptsRate,
		Dims: module.Dims{
			{ID: "iis_website_%s_logon_attempts_total", Name: "logon", Algo: module.Incremental},
		},
	}
	iisWebsiteUptimeChartTmpl = module.Chart{
		OverModule: "iis",
		ID:         "iis_website_%s_uptime",
		Title:      "Website uptime",
		Units:      "seconds",
		Fam:        "uptime",
		Ctx:        "iis.website_uptime",
		Priority:   prioIISWebsiteUptime,
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
		mssqlDeadLocksChart.Copy(),
		mssqlMemmgrConnectionMemoryBytesChart.Copy(),
		mssqlMemmgrExternalBenefitOfMemoryChart.Copy(),
		mssqlMemmgrPendingMemoryChart.Copy(),
		mssqlMemmgrTotalServerChart.Copy(),
		mssqlSQLErrorsTotalChart.Copy(),
		mssqlStatsAutoParamChart.Copy(),
		mssqlStatsBatchRequestsChart.Copy(),
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
		OverModule: "mssql",
		ID:         "mssql_instance_%s_accessmethods_page_splits",
		Title:      "Page splits",
		Units:      "splits/s",
		Fam:        "buffer cache",
		Ctx:        "mssql.instance_accessmethods_page_splits",
		Priority:   prioMSSQLAccessMethodPageSplits,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_accessmethods_page_splits", Name: "page", Algo: module.Incremental},
		},
	}
	// Buffer Management
	// Source: https://learn.microsoft.com/en-us/sql/relational-databases/performance-monitor/sql-server-buffer-manager-object?view=sql-server-ver16
	mssqlCacheHitRatioChart = module.Chart{
		OverModule: "mssql",
		ID:         "mssql_instance_%s_cache_hit_ratio",
		Title:      "Buffer Cache hit ratio",
		Units:      "percentage",
		Fam:        "buffer cache",
		Ctx:        "mssql.instance_cache_hit_ratio",
		Priority:   prioMSSQLCacheHitRatio,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_cache_hit_ratio", Name: "hit_ratio"},
		},
	}
	mssqlBufferCheckpointPageChart = module.Chart{
		OverModule: "mssql",
		ID:         "mssql_instance_%s_bufman_checkpoint_pages",
		Title:      "Flushed pages",
		Units:      "pages/s",
		Fam:        "buffer cache",
		Ctx:        "mssql.instance_bufman_checkpoint_pages",
		Priority:   prioMSSQLBufferCheckpointPages,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_bufman_checkpoint_pages", Name: "flushed", Algo: module.Incremental},
		},
	}
	mssqlBufferPageLifeExpectancyChart = module.Chart{
		OverModule: "mssql",
		ID:         "mssql_instance_%s_bufman_page_life_expectancy",
		Title:      "Page life expectancy",
		Units:      "seconds",
		Fam:        "buffer cache",
		Ctx:        "mssql.instance_bufman_page_life_expectancy",
		Priority:   prioMSSQLBufferPageLifeExpectancy,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_bufman_page_life_expectancy_seconds", Name: "life_expectancy"},
		},
	}
	mssqlBufManIOPSChart = module.Chart{
		OverModule: "mssql",
		ID:         "mssql_instance_%s_bufman_iops",
		Title:      "Number of pages input and output",
		Units:      "pages/s",
		Fam:        "buffer cache",
		Ctx:        "mssql.instance_bufman_iops",
		Priority:   prioMSSQLBufManIOPS,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_bufman_page_reads", Name: "read", Algo: module.Incremental},
			{ID: "mssql_instance_%s_bufman_page_writes", Name: "written", Mul: -1, Algo: module.Incremental},
		},
	}
	// General Statistic
	// Source: https://learn.microsoft.com/en-us/sql/relational-databases/performance-monitor/sql-server-general-statistics-object?view=sql-server-ver16
	mssqlBlockedProcessChart = module.Chart{
		OverModule: "mssql",
		ID:         "mssql_instance_%s_blocked_process",
		Title:      "Blocked processes",
		Units:      "process",
		Fam:        "processes",
		Ctx:        "mssql.instance_blocked_processes",
		Priority:   prioMSSQLBlockedProcess,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_genstats_blocked_processes", Name: "blocked"},
		},
	}
	mssqlUserConnectionChart = module.Chart{
		OverModule: "mssql",
		ID:         "mssql_instance_%s_user_connection",
		Title:      "User connections",
		Units:      "connections",
		Fam:        "connections",
		Ctx:        "mssql.instance_user_connection",
		Priority:   prioMSSQLUserConnections,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_genstats_user_connections", Name: "user"},
		},
	}
	// Lock Wait
	// Source: https://learn.microsoft.com/en-us/sql/relational-databases/performance-monitor/sql-server-locks-object?view=sql-server-ver16
	mssqlLocksWaitChart = module.Chart{
		OverModule: "mssql",
		ID:         "mssql_instance_%s_locks_lock_wait",
		Title:      "Lock requests that required the caller to wait",
		Units:      "locks/s",
		Fam:        "locks",
		Ctx:        "mssql.instance_locks_lock_wait",
		Priority:   prioMSSQLLocksLockWait,
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
	mssqlDeadLocksChart = module.Chart{
		OverModule: "mssql",
		ID:         "mssql_instance_%s_locks_deadlocks",
		Title:      "Lock requests that resulted in deadlock.",
		Units:      "locks/s",
		Fam:        "locks",
		Ctx:        "mssql.instance_locks_deadlocks",
		Priority:   prioMSSQLLocksLockWait,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_resource_AllocUnit_locks_deadlocks", Name: "alloc_unit", Algo: module.Incremental},
			{ID: "mssql_instance_%s_resource_Application_locks_deadlocks", Name: "application", Algo: module.Incremental},
			{ID: "mssql_instance_%s_resource_Database_locks_deadlocks", Name: "database", Algo: module.Incremental},
			{ID: "mssql_instance_%s_resource_Extent_locks_deadlocks", Name: "extent", Algo: module.Incremental},
			{ID: "mssql_instance_%s_resource_File_locks_deadlocks", Name: "file", Algo: module.Incremental},
			{ID: "mssql_instance_%s_resource_HoBT_locks_deadlocks", Name: "hobt", Algo: module.Incremental},
			{ID: "mssql_instance_%s_resource_Key_locks_deadlocks", Name: "key", Algo: module.Incremental},
			{ID: "mssql_instance_%s_resource_Metadata_locks_deadlocks", Name: "metadata", Algo: module.Incremental},
			{ID: "mssql_instance_%s_resource_OIB_locks_deadlocks", Name: "oib", Algo: module.Incremental},
			{ID: "mssql_instance_%s_resource_Object_locks_deadlocks", Name: "object", Algo: module.Incremental},
			{ID: "mssql_instance_%s_resource_Page_locks_deadlocks", Name: "page", Algo: module.Incremental},
			{ID: "mssql_instance_%s_resource_RID_locks_deadlocks", Name: "rid", Algo: module.Incremental},
			{ID: "mssql_instance_%s_resource_RowGroup_locks_deadlocks", Name: "row_group", Algo: module.Incremental},
			{ID: "mssql_instance_%s_resource_Xact_locks_deadlocks", Name: "xact", Algo: module.Incremental},
		},
	}

	// Memory Manager
	// Source: https://learn.microsoft.com/en-us/sql/relational-databases/performance-monitor/sql-server-memory-manager-object?view=sql-server-ver16
	mssqlMemmgrConnectionMemoryBytesChart = module.Chart{
		OverModule: "mssql",
		ID:         "mssql_instance_%s_memmgr_connection_memory_bytes",
		Title:      "Amount of dynamic memory to maintain connections",
		Units:      "bytes",
		Fam:        "memory",
		Ctx:        "mssql.instance_memmgr_connection_memory_bytes",
		Priority:   prioMSSQLMemmgrConnectionMemoryBytes,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_memmgr_connection_memory_bytes", Name: "memory", Algo: module.Incremental},
		},
	}
	mssqlMemmgrExternalBenefitOfMemoryChart = module.Chart{
		OverModule: "mssql",
		ID:         "mssql_instance_%s_memmgr_external_benefit_of_memory",
		Title:      "Performance benefit from adding memory to a specific cache",
		Units:      "bytes",
		Fam:        "memory",
		Ctx:        "mssql.instance_memmgr_external_benefit_of_memory",
		Priority:   prioMSSQLMemmgrExternalBenefitOfMemory,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_memmgr_external_benefit_of_memory", Name: "benefit", Algo: module.Incremental},
		},
	}
	mssqlMemmgrPendingMemoryChart = module.Chart{
		OverModule: "mssql",
		ID:         "mssql_instance_%s_memmgr_pending_memory_grants",
		Title:      "Process waiting for memory grant",
		Units:      "process",
		Fam:        "memory",
		Ctx:        "mssql.instance_memmgr_pending_memory_grants",
		Priority:   prioMSSQLMemmgrPendingMemoryGrants,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_memmgr_pending_memory_grants", Name: "pending"},
		},
	}
	mssqlMemmgrTotalServerChart = module.Chart{
		OverModule: "mssql",
		ID:         "mssql_instance_%s_memmgr_server_memory",
		Title:      "Memory committed",
		Units:      "bytes",
		Fam:        "memory",
		Ctx:        "mssql.instance_memmgr_server_memory",
		Priority:   prioMSSQLMemTotalServer,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_memmgr_total_server_memory_bytes", Name: "memory"},
		},
	}

	// SQL errors
	// Source: https://learn.microsoft.com/en-us/sql/relational-databases/performance-monitor/sql-server-sql-errors-object?view=sql-server-ver16
	mssqlSQLErrorsTotalChart = module.Chart{
		OverModule: "mssql",
		ID:         "mssql_instance_%s_sql_errors_total",
		Title:      "Errors",
		Units:      "errors/s",
		Fam:        "errors",
		Ctx:        "mssql.instance_sql_errors",
		Priority:   prioMSSQLSqlErrorsTotal,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_sql_errors_total_db_offline_errors", Name: "db_offline", Algo: module.Incremental},
			{ID: "mssql_instance_%s_sql_errors_total_info_errors", Name: "info", Algo: module.Incremental},
			{ID: "mssql_instance_%s_sql_errors_total_kill_connection_errors", Name: "kill_connection", Algo: module.Incremental},
			{ID: "mssql_instance_%s_sql_errors_total_user_errors", Name: "user", Algo: module.Incremental},
		},
	}

	// SQL Statistic
	// Source: https://learn.microsoft.com/en-us/sql/relational-databases/performance-monitor/sql-server-sql-statistics-object?view=sql-server-ver16
	mssqlStatsAutoParamChart = module.Chart{
		OverModule: "mssql",
		ID:         "mssql_instance_%s_sqlstats_auto_parameterization_attempts",
		Title:      "Failed auto-parameterization attempts",
		Units:      "attempts/s",
		Fam:        "sql activity",
		Ctx:        "mssql.instance_sqlstats_auto_parameterization_attempts",
		Priority:   prioMSSQLStatsAutoParameterization,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_sqlstats_auto_parameterization_attempts", Name: "failed", Algo: module.Incremental},
		},
	}
	mssqlStatsBatchRequestsChart = module.Chart{
		OverModule: "mssql",
		ID:         "mssql_instance_%s_sqlstats_batch_requests",
		Title:      "Total of batches requests",
		Units:      "requests/s",
		Fam:        "sql activity",
		Ctx:        "mssql.instance_sqlstats_batch_requests",
		Priority:   prioMSSQLStatsBatchRequests,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_sqlstats_batch_requests", Name: "batch", Algo: module.Incremental},
		},
	}
	mssqlStatsSafeAutoChart = module.Chart{
		OverModule: "mssql",
		ID:         "mssql_instance_%s_sqlstats_safe_auto_parameterization_attempts",
		Title:      "Safe auto-parameterization attempts",
		Units:      "attempts/s",
		Fam:        "sql activity",
		Ctx:        "mssql.instance_sqlstats_safe_auto_parameterization_attempts",
		Priority:   prioMSSQLStatsSafeAutoParameterization,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_sqlstats_safe_auto_parameterization_attempts", Name: "safe", Algo: module.Incremental},
		},
	}
	mssqlStatsCompilationChart = module.Chart{
		OverModule: "mssql",
		ID:         "mssql_instance_%s_sqlstats_sql_compilations",
		Title:      "SQL compilations",
		Units:      "compilations/s",
		Fam:        "sql activity",
		Ctx:        "mssql.instance_sqlstats_sql_compilations",
		Priority:   prioMSSQLStatsCompilations,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_sqlstats_sql_compilations", Name: "compilations", Algo: module.Incremental},
		},
	}
	mssqlStatsRecompilationChart = module.Chart{
		OverModule: "mssql",
		ID:         "mssql_instance_%s_sqlstats_sql_recompilations",
		Title:      "SQL re-compilations",
		Units:      "recompiles/s",
		Fam:        "sql activity",
		Ctx:        "mssql.instance_sqlstats_sql_recompilations",
		Priority:   prioMSSQLStatsRecompilations,
		Dims: module.Dims{
			{ID: "mssql_instance_%s_sqlstats_sql_recompilations", Name: "recompiles", Algo: module.Incremental},
		},
	}

	// Database
	// Source: https://learn.microsoft.com/en-us/sql/relational-databases/performance-monitor/sql-server-databases-object?view=sql-server-2017
	mssqlDatabaseActiveTransactionsChart = module.Chart{
		OverModule: "mssql",
		ID:         "mssql_db_%s_instance_%s_active_transactions",
		Title:      "Active transactions per database",
		Units:      "transactions",
		Fam:        "transactions",
		Ctx:        "mssql.database_active_transactions",
		Priority:   prioMSSQLDatabaseActiveTransactions,
		Dims: module.Dims{
			{ID: "mssql_db_%s_instance_%s_active_transactions", Name: "active"},
		},
	}
	mssqlDatabaseBackupRestoreOperationsChart = module.Chart{
		OverModule: "mssql",
		ID:         "mssql_db_%s_instance_%s_backup_restore_operations",
		Title:      "Backup IO per database",
		Units:      "operations/s",
		Fam:        "transactions",
		Ctx:        "mssql.database_backup_restore_operations",
		Priority:   prioMSSQLDatabaseBackupRestoreOperations,
		Dims: module.Dims{
			{ID: "mssql_db_%s_instance_%s_backup_restore_operations", Name: "backup", Algo: module.Incremental},
		},
	}
	mssqlDatabaseSizeChart = module.Chart{
		OverModule: "mssql",
		ID:         "mssql_db_%s_instance_%s_data_files_size",
		Title:      "Current database size",
		Units:      "bytes",
		Fam:        "size",
		Ctx:        "mssql.database_data_files_size",
		Priority:   prioMSSQLDatabaseDataFileSize,
		Dims: module.Dims{
			{ID: "mssql_db_%s_instance_%s_data_files_size_bytes", Name: "size"},
		},
	}
	mssqlDatabaseLogFlushedChart = module.Chart{
		OverModule: "mssql",
		ID:         "mssql_db_%s_instance_%s_log_flushed",
		Title:      "Log flushed",
		Units:      "bytes/s",
		Fam:        "transactions",
		Ctx:        "mssql.database_log_flushed",
		Priority:   prioMSSQLDatabaseLogFlushed,
		Dims: module.Dims{
			{ID: "mssql_db_%s_instance_%s_log_flushed_bytes", Name: "flushed", Algo: module.Incremental},
		},
	}
	mssqlDatabaseLogFlushesChart = module.Chart{
		OverModule: "mssql",
		ID:         "mssql_db_%s_instance_%s_log_flushes",
		Title:      "Log flushes",
		Units:      "flushes/s",
		Fam:        "transactions",
		Ctx:        "mssql.database_log_flushes",
		Priority:   prioMSSQLDatabaseLogFlushes,
		Dims: module.Dims{
			{ID: "mssql_db_%s_instance_%s_log_flushes", Name: "log", Algo: module.Incremental},
		},
	}
	mssqlDatabaseTransactionsChart = module.Chart{
		OverModule: "mssql",
		ID:         "mssql_db_%s_instance_%s_transactions",
		Title:      "Transactions",
		Units:      "transactions/s",
		Fam:        "transactions",
		Ctx:        "mssql.database_transactions",
		Priority:   prioMSSQLDatabaseTransactions,
		Dims: module.Dims{
			{ID: "mssql_db_%s_instance_%s_transactions", Name: "transactions", Algo: module.Incremental},
		},
	}
	mssqlDatabaseWriteTransactionsChart = module.Chart{
		OverModule: "mssql",
		ID:         "mssql_db_%s_instance_%s_write_transactions",
		Title:      "Write transactions",
		Units:      "transactions/s",
		Fam:        "transactions",
		Ctx:        "mssql.instance_write_transactions",
		Priority:   prioMSSQLDatabaseWriteTransactions,
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
		OverModule: "ad",
		ID:         "ad_dra_replication_intersite_compressed_traffic",
		Title:      "DRA replication compressed traffic withing site",
		Units:      "bytes/s",
		Fam:        "replication",
		Ctx:        "ad.dra_replication_intersite_compressed_traffic",
		Priority:   prioADDRAReplicationIntersiteCompressedTraffic,
		Type:       module.Area,
		Dims: module.Dims{
			{ID: "ad_replication_data_intersite_bytes_total_inbound", Name: "inbound", Algo: module.Incremental},
			{ID: "ad_replication_data_intersite_bytes_total_outbound", Name: "outbound", Algo: module.Incremental, Mul: -1},
		},
	}
	adDRAReplicationIntrasiteCompressedTrafficChart = module.Chart{
		OverModule: "ad",
		ID:         "ad_dra_replication_intrasite_compressed_traffic",
		Title:      "DRA replication compressed traffic between sites",
		Units:      "bytes/s",
		Fam:        "replication",
		Ctx:        "ad.dra_replication_intrasite_compressed_traffic",
		Priority:   prioADDRAReplicationIntrasiteCompressedTraffic,
		Type:       module.Area,
		Dims: module.Dims{
			{ID: "ad_replication_data_intrasite_bytes_total_inbound", Name: "inbound", Algo: module.Incremental},
			{ID: "ad_replication_data_intrasite_bytes_total_outbound", Name: "outbound", Algo: module.Incremental, Mul: -1},
		},
	}
	adDRAReplicationSyncObjectRemainingChart = module.Chart{
		OverModule: "ad",
		ID:         "ad_dra_replication_sync_objects_remaining",
		Title:      "DRA replication full sync objects remaining",
		Units:      "objects",
		Fam:        "replication",
		Ctx:        "ad.dra_replication_sync_objects_remaining",
		Priority:   prioADDRAReplicationSyncObjectsRemaining,
		Dims: module.Dims{
			{ID: "ad_replication_inbound_sync_objects_remaining", Name: "inbound"},
		},
	}
	adDRAReplicationObjectsFilteredChart = module.Chart{
		OverModule: "ad",
		ID:         "ad_dra_replication_objects_filtered",
		Title:      "DRA replication objects filtered",
		Units:      "objects/s",
		Fam:        "replication",
		Ctx:        "ad.dra_replication_objects_filtered",
		Priority:   prioADDRAReplicationObjectsFiltered,
		Dims: module.Dims{
			{ID: "ad_replication_inbound_objects_filtered_total", Name: "inbound", Algo: module.Incremental},
		},
	}
	adDRAReplicationPropertiesUpdatedChart = module.Chart{
		OverModule: "ad",
		ID:         "ad_dra_replication_properties_updated",
		Title:      "DRA replication properties updated",
		Units:      "properties/s",
		Fam:        "replication",
		Ctx:        "ad.dra_replication_properties_updated",
		Priority:   prioADDRAReplicationPropertiesUpdated,
		Dims: module.Dims{
			{ID: "ad_replication_inbound_properties_updated_total", Name: "inbound", Algo: module.Incremental},
		},
	}
	adDRAReplicationPropertiesFilteredChart = module.Chart{
		OverModule: "ad",
		ID:         "ad_dra_replication_properties_filtered",
		Title:      "DRA replication properties filtered",
		Units:      "properties/s",
		Fam:        "replication",
		Ctx:        "ad.dra_replication_properties_filtered",
		Priority:   prioADDRAReplicationPropertiesFiltered,
		Dims: module.Dims{
			{ID: "ad_replication_inbound_properties_filtered_total", Name: "inbound", Algo: module.Incremental},
		},
	}
	adDRAReplicationPendingSyncsChart = module.Chart{
		OverModule: "ad",
		ID:         "ad_dra_replication_pending_syncs",
		Title:      "DRA replication pending syncs",
		Units:      "syncs",
		Fam:        "replication",
		Ctx:        "ad.dra_replication_pending_syncs",
		Priority:   prioADReplicationPendingSyncs,
		Dims: module.Dims{
			{ID: "ad_replication_pending_synchronizations", Name: "pending"},
		},
	}
	adDRAReplicationSyncRequestsChart = module.Chart{
		OverModule: "ad",
		ID:         "ad_dra_replication_sync_requests",
		Title:      "DRA replication sync requests",
		Units:      "requests/s",
		Fam:        "replication",
		Ctx:        "ad.dra_replication_sync_requests",
		Priority:   prioADDRASyncRequests,
		Dims: module.Dims{
			{ID: "ad_replication_sync_requests_total", Name: "request", Algo: module.Incremental},
		},
	}
	adDirectoryServiceThreadsChart = module.Chart{
		OverModule: "ad",
		ID:         "ad_ds_threads",
		Title:      "Directory Service threads",
		Units:      "threads",
		Fam:        "replication",
		Ctx:        "ad.ds_threads",
		Priority:   prioADDirectoryServiceThreadsInUse,
		Dims: module.Dims{
			{ID: "ad_directory_service_threads", Name: "in_use"},
		},
	}
	adLDAPLastBindTimeChart = module.Chart{
		OverModule: "ad",
		ID:         "ad_ldap_last_bind_time",
		Title:      "LDAP last successful bind time",
		Units:      "seconds",
		Fam:        "bind",
		Ctx:        "ad.ldap_last_bind_time",
		Priority:   prioADLDAPBindTime,
		Dims: module.Dims{
			{ID: "ad_ldap_last_bind_time_seconds", Name: "last_bind"},
		},
	}
	adBindsTotalChart = module.Chart{
		OverModule: "ad",
		ID:         "ad_binds",
		Title:      "Successful binds",
		Units:      "bind/s",
		Fam:        "bind",
		Ctx:        "ad.binds",
		Priority:   prioADBindsTotal,
		Dims: module.Dims{
			{ID: "ad_binds_total", Name: "binds", Algo: module.Incremental},
		},
	}
	adLDAPSearchesChart = module.Chart{
		OverModule: "ad",
		ID:         "ad_ldap_searches",
		Title:      "LDAP client search operations",
		Units:      "searches/s",
		Fam:        "ldap",
		Ctx:        "ad.ldap_searches",
		Priority:   prioADLDAPSearchesTotal,
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
		OverModule: "adcs",
		ID:         "adcs_cert_template%s_requests",
		Title:      "Certificate requests processed",
		Units:      "requests/s",
		Fam:        "requests",
		Ctx:        "adcs.cert_template_requests",
		Priority:   prioADCSCertTemplateRequests,
		Dims: module.Dims{
			{ID: "adcs_cert_template_%s_requests_total", Name: "requests", Algo: module.Incremental},
		},
	}
	adcsCertTemplateFailedRequestsChartTmpl = module.Chart{
		OverModule: "adcs",
		ID:         "adcs_cert_template_%s_failed_requests",
		Title:      "Certificate failed requests processed",
		Units:      "requests/s",
		Fam:        "requests",
		Ctx:        "adcs.cert_template_failed_requests",
		Priority:   prioADCSCertTemplateFailedRequests,
		Dims: module.Dims{
			{ID: "adcs_cert_template_%s_failed_requests_total", Name: "failed", Algo: module.Incremental},
		},
	}
	adcsCertTemplateIssuedRequestsChartTmpl = module.Chart{
		OverModule: "adcs",
		ID:         "adcs_cert_template_%s_issued_requests",
		Title:      "Certificate issued requests processed",
		Units:      "requests/s",
		Fam:        "requests",
		Ctx:        "adcs.cert_template_issued_requests",
		Priority:   prioADCSCertTemplateIssuesRequests,
		Dims: module.Dims{
			{ID: "adcs_cert_template_%s_issued_requests_total", Name: "issued", Algo: module.Incremental},
		},
	}
	adcsCertTemplatePendingRequestsChartTmpl = module.Chart{
		OverModule: "adcs",
		ID:         "adcs_cert_template_%s_pending_requests",
		Title:      "Certificate pending requests processed",
		Units:      "requests/s",
		Fam:        "requests",
		Ctx:        "adcs.cert_template_pending_requests",
		Priority:   prioADCSCertTemplatePendingRequests,
		Dims: module.Dims{
			{ID: "adcs_cert_template_%s_pending_requests_total", Name: "pending", Algo: module.Incremental},
		},
	}
	adcsCertTemplateRequestProcessingTimeChartTmpl = module.Chart{
		OverModule: "adcs",
		ID:         "adcs_cert_template_%s_request_processing_time",
		Title:      "Certificate last request processing time",
		Units:      "seconds",
		Fam:        "requests",
		Ctx:        "adcs.cert_template_request_processing_time",
		Priority:   prioADCSCertTemplateRequestProcessingTime,
		Dims: module.Dims{
			{ID: "adcs_cert_template_%s_request_processing_time_seconds", Name: "processing_time", Div: precision},
		},
	}
	adcsCertTemplateChallengeResponseChartTmpl = module.Chart{
		OverModule: "adcs",
		ID:         "adcs_cert_template_%s_challenge_responses",
		Title:      "Certificate challenge responses",
		Units:      "responses/s",
		Fam:        "responses",
		Ctx:        "adcs.cert_template_challenge_responses",
		Priority:   prioADCSCertTemplateChallengeResponses,
		Dims: module.Dims{
			{ID: "adcs_cert_template_%s_challenge_responses_total", Name: "challenge", Algo: module.Incremental},
		},
	}
	adcsCertTemplateRetrievalsChartTmpl = module.Chart{
		OverModule: "adcs",
		ID:         "adcs_cert_template_%s_retrievals",
		Title:      "Total of certificate retrievals",
		Units:      "retrievals/s",
		Fam:        "retrievals",
		Ctx:        "adcs.cert_template_retrievals",
		Priority:   prioADCSCertTemplateRetrievals,
		Dims: module.Dims{
			{ID: "adcs_cert_template_%s_retrievals_total", Name: "retrievals", Algo: module.Incremental},
		},
	}
	adcsCertificateRetrievalsTimeChartTmpl = module.Chart{
		OverModule: "adcs",
		ID:         "adcs_cert_template_%s_retrievals_processing_time",
		Title:      "Certificate last retrieval processing time",
		Units:      "seconds",
		Fam:        "retrievals",
		Ctx:        "adcs.cert_template_retrieval_processing_time",
		Priority:   prioADCSCertTemplateRetrievalProcessingTime,
		Dims: module.Dims{
			{ID: "adcs_cert_template_%s_retrievals_processing_time_seconds", Name: "processing_time", Div: precision},
		},
	}
	adcsCertTemplateRequestCryptoSigningTimeChartTmpl = module.Chart{
		OverModule: "adcs",
		ID:         "adcs_cert_template_%s_request_cryptographic_signing_time",
		Title:      "Certificate last signing operation request time",
		Units:      "seconds",
		Fam:        "timings",
		Ctx:        "adcs.cert_template_request_cryptographic_signing_time",
		Priority:   prioADCSCertTemplateRequestCryptoSigningTime,
		Dims: module.Dims{
			{ID: "adcs_cert_template_%s_request_cryptographic_signing_time_seconds", Name: "singing_time", Div: precision},
		},
	}
	adcsCertTemplateRequestPolicyModuleProcessingTimeChartTmpl = module.Chart{
		OverModule: "adcs",
		ID:         "adcs_cert_template_%s_request_policy_module_processing_time",
		Title:      "Certificate last policy module processing request time",
		Units:      "seconds",
		Fam:        "timings",
		Ctx:        "adcs.cert_template_request_policy_module_processing",
		Priority:   prioADCSCertTemplateRequestPolicyModuleProcessingTime,
		Dims: module.Dims{
			{ID: "adcs_cert_template_%s_request_policy_module_processing_time_seconds", Name: "processing_time", Div: precision},
		},
	}
	adcsCertTemplateChallengeResponseProcessingTimeChartTmpl = module.Chart{
		OverModule: "adcs",
		ID:         "adcs_cert_template_%s_challenge_response_processing_time",
		Title:      "Certificate last challenge response time",
		Units:      "seconds",
		Fam:        "timings",
		Ctx:        "adcs.cert_template_challenge_response_processing_time",
		Priority:   prioADCSCertTemplateChallengeResponseProcessingTime,
		Dims: module.Dims{
			{ID: "adcs_cert_template_%s_challenge_response_processing_time_seconds", Name: "processing_time", Div: precision},
		},
	}
	adcsCertTemplateSignedCertificateTimestampListsChartTmpl = module.Chart{
		OverModule: "adcs",
		ID:         "adcs_cert_template_%s_signed_certificate_timestamp_lists",
		Title:      "Certificate Signed Certificate Timestamp Lists processed",
		Units:      "lists/s",
		Fam:        "timings",
		Ctx:        "adcs.cert_template_signed_certificate_timestamp_lists",
		Priority:   prioADCSCertTemplateSignedCertificateTimestampLists,
		Dims: module.Dims{
			{ID: "adcs_cert_template_%s_signed_certificate_timestamp_lists_total", Name: "processed", Algo: module.Incremental},
		},
	}
	adcsCertTemplateSignedCertificateTimestampListProcessingTimeChartTmpl = module.Chart{
		OverModule: "adcs",
		ID:         "adcs_cert_template_%s_signed_certificate_timestamp_list_processing_time",
		Title:      "Certificate last Signed Certificate Timestamp List process time",
		Units:      "seconds",
		Fam:        "timings",
		Ctx:        "adcs.cert_template_signed_certificate_timestamp_list_processing_time",
		Priority:   prioADCSCertTemplateSignedCertificateTimestampListProcessingTime,
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
		OverModule: "adfs",
		ID:         "adfs_ad_login_connection_failures",
		Title:      "Connection failures",
		Units:      "failures/s",
		Fam:        "ad",
		Ctx:        "adfs.ad_login_connection_failures",
		Priority:   prioADFSADLoginConnectionFailures,
		Dims: module.Dims{
			{ID: "adfs_ad_login_connection_failures_total", Name: "connection", Algo: module.Incremental},
		},
	}
	adfsCertificateAuthenticationsChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_certificate_authentications",
		Title:      "User Certificate authentications",
		Units:      "authentications/s",
		Fam:        "auth",
		Ctx:        "adfs.certificate_authentications",
		Priority:   prioADFSCertificateAuthentications,
		Dims: module.Dims{
			{ID: "adfs_certificate_authentications_total", Name: "authentications", Algo: module.Incremental},
		},
	}

	adfsDBArtifactFailuresChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_db_artifact_failures",
		Title:      "Connection failures to the artifact database",
		Units:      "failures/s",
		Fam:        "db artifact",
		Ctx:        "adfs.db_artifact_failures",
		Priority:   prioADFSDBArtifactFailures,
		Dims: module.Dims{
			{ID: "adfs_db_artifact_failure_total", Name: "connection", Algo: module.Incremental},
		},
	}
	adfsDBArtifactQueryTimeSecondsChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_db_artifact_query_time_seconds",
		Title:      "Time taken for an artifact database query",
		Units:      "seconds/s",
		Fam:        "db artifact",
		Ctx:        "adfs.db_artifact_query_time_seconds",
		Priority:   prioADFSDBArtifactQueryTimeSeconds,
		Dims: module.Dims{
			{ID: "adfs_db_artifact_query_time_seconds_total", Name: "query_time", Algo: module.Incremental, Div: precision},
		},
	}
	adfsDBConfigFailuresChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_db_config_failures",
		Title:      "Connection failures to the configuration database",
		Units:      "failures/s",
		Fam:        "db config",
		Ctx:        "adfs.db_config_failures",
		Priority:   prioADFSDBConfigFailures,
		Dims: module.Dims{
			{ID: "adfs_db_config_failure_total", Name: "connection", Algo: module.Incremental},
		},
	}
	adfsDBConfigQueryTimeSecondsChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_db_config_query_time_seconds",
		Title:      "Time taken for a configuration database query",
		Units:      "seconds/s",
		Fam:        "db config",
		Ctx:        "adfs.db_config_query_time_seconds",
		Priority:   prioADFSDBConfigQueryTimeSeconds,
		Dims: module.Dims{
			{ID: "adfs_db_config_query_time_seconds_total", Name: "query_time", Algo: module.Incremental, Div: precision},
		},
	}
	adfsDeviceAuthenticationsChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_device_authentications",
		Title:      "Device authentications",
		Units:      "authentications/s",
		Fam:        "auth",
		Ctx:        "adfs.device_authentications",
		Priority:   prioADFSDeviceAuthentications,
		Dims: module.Dims{
			{ID: "adfs_device_authentications_total", Name: "authentications", Algo: module.Incremental},
		},
	}
	adfsExternalAuthenticationsChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_external_authentications",
		Title:      "Authentications from external MFA providers",
		Units:      "authentications/s",
		Fam:        "auth",
		Ctx:        "adfs.external_authentications",
		Priority:   prioADFSExternalAuthentications,
		Dims: module.Dims{
			{ID: "adfs_external_authentications_success_total", Name: "success", Algo: module.Incremental},
			{ID: "adfs_external_authentications_failure_total", Name: "failure", Algo: module.Incremental},
		},
	}
	adfsFederatedAuthenticationsChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_federated_authentications",
		Title:      "Authentications from Federated Sources",
		Units:      "authentications/s",
		Fam:        "auth",
		Ctx:        "adfs.federated_authentications",
		Priority:   prioADFSFederatedAuthentications,
		Dims: module.Dims{
			{ID: "adfs_federated_authentications_total", Name: "authentications", Algo: module.Incremental},
		},
	}
	adfsFederationMetadataRequestsChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_federation_metadata_requests",
		Title:      "Federation Metadata requests",
		Units:      "requests/s",
		Fam:        "requests",
		Ctx:        "adfs.federation_metadata_requests",
		Priority:   prioADFSFederationMetadataRequests,
		Dims: module.Dims{
			{ID: "adfs_federation_metadata_requests_total", Name: "requests", Algo: module.Incremental},
		},
	}

	adfsOAuthAuthorizationRequestsChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_oauth_authorization_requests",
		Title:      "Incoming requests to the OAuth Authorization endpoint",
		Units:      "requests/s",
		Fam:        "oauth",
		Ctx:        "adfs.oauth_authorization_requests",
		Priority:   prioADFSOauthAuthorizationRequests,
		Dims: module.Dims{
			{ID: "adfs_oauth_authorization_requests_total", Name: "requests", Algo: module.Incremental},
		},
	}
	adfsOAuthClientAuthenticationsChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_oauth_client_authentications",
		Title:      "OAuth client authentications",
		Units:      "authentications/s",
		Fam:        "oauth",
		Ctx:        "adfs.oauth_client_authentications",
		Priority:   prioADFSOauthClientAuthentications,
		Dims: module.Dims{
			{ID: "adfs_oauth_client_authentication_success_total", Name: "success", Algo: module.Incremental},
			{ID: "adfs_oauth_client_authentication_failure_total", Name: "failure", Algo: module.Incremental},
		},
	}
	adfsOAuthClientCredentialRequestsChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_oauth_client_credentials_requests",
		Title:      "OAuth client credentials requests",
		Units:      "requests/s",
		Fam:        "oauth",
		Ctx:        "adfs.oauth_client_credentials_requests",
		Priority:   prioADFSOauthClientCredentials,
		Dims: module.Dims{
			{ID: "adfs_oauth_client_credentials_success_total", Name: "success", Algo: module.Incremental},
			{ID: "adfs_oauth_client_credentials_failure_total", Name: "failure", Algo: module.Incremental},
		},
	}
	adfsOAuthClientPrivKeyJwtAuthenticationsChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_oauth_client_privkey_jwt_authentications",
		Title:      "OAuth client private key JWT authentications",
		Units:      "authentications/s",
		Fam:        "oauth",
		Ctx:        "adfs.oauth_client_privkey_jwt_authentications",
		Priority:   prioADFSOauthClientPrivkeyJwtAuthentication,
		Dims: module.Dims{
			{ID: "adfs_oauth_client_privkey_jwt_authentications_success_total", Name: "success", Algo: module.Incremental},
			{ID: "adfs_oauth_client_privkey_jtw_authentication_failure_total", Name: "failure", Algo: module.Incremental},
		},
	}
	adfsOAuthClientSecretBasicAuthenticationsChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_oauth_client_secret_basic_authentications",
		Title:      "OAuth client secret basic authentications",
		Units:      "authentications/s",
		Fam:        "oauth",
		Ctx:        "adfs.oauth_client_secret_basic_authentications",
		Type:       module.Line,
		Priority:   prioADFSOauthClientSecretBasicAuthentications,
		Dims: module.Dims{
			{ID: "adfs_oauth_client_secret_basic_authentications_success_total", Name: "success", Algo: module.Incremental},
			{ID: "adfs_oauth_client_secret_basic_authentications_failure_total", Name: "failure", Algo: module.Incremental},
		},
	}
	adfsOAuthClientSecretPostAuthenticationsChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_oauth_client_secret_post_authentications",
		Title:      "OAuth client secret post authentications",
		Units:      "authentications/s",
		Fam:        "oauth",
		Ctx:        "adfs.oauth_client_secret_post_authentications",
		Priority:   prioADFSOauthClientSecretPostAuthentications,
		Dims: module.Dims{
			{ID: "adfs_oauth_client_secret_post_authentications_success_total", Name: "success", Algo: module.Incremental},
			{ID: "adfs_oauth_client_secret_post_authentications_failure_total", Name: "failure", Algo: module.Incremental},
		},
	}
	adfsOAuthClientWindowsAuthenticationsChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_oauth_client_windows_authentications",
		Title:      "OAuth client windows integrated authentications",
		Units:      "authentications/s",
		Fam:        "oauth",
		Ctx:        "adfs.oauth_client_windows_authentications",
		Priority:   prioADFSOauthClientWindowsAuthentications,
		Dims: module.Dims{
			{ID: "adfs_oauth_client_windows_authentications_success_total", Name: "success", Algo: module.Incremental},
			{ID: "adfs_oauth_client_windows_authentications_failure_total", Name: "failure", Algo: module.Incremental},
		},
	}
	adfsOAuthLogonCertificateRequestsChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_oauth_logon_certificate_requests",
		Title:      "OAuth logon certificate requests",
		Units:      "requests/s",
		Fam:        "oauth",
		Ctx:        "adfs.oauth_logon_certificate_requests",
		Priority:   prioADFSOauthLogonCertificateRequests,
		Dims: module.Dims{
			{ID: "adfs_oauth_logon_certificate_token_requests_success_total", Name: "success", Algo: module.Incremental},
			{ID: "adfs_oauth_logon_certificate_requests_failure_total", Name: "failure", Algo: module.Incremental},
		},
	}
	adfsOAuthPasswordGrantRequestsChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_oauth_password_grant_requests",
		Title:      "OAuth password grant requests",
		Units:      "requests/s",
		Fam:        "oauth",
		Ctx:        "adfs.oauth_password_grant_requests",
		Priority:   prioADFSOauthPasswordGrantRequests,
		Dims: module.Dims{
			{ID: "adfs_oauth_password_grant_requests_success_total", Name: "success", Algo: module.Incremental},
			{ID: "adfs_oauth_password_grant_requests_failure_total", Name: "failure", Algo: module.Incremental},
		},
	}
	adfsOAuthTokenRequestsChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_oauth_token_requests_success",
		Title:      "Successful RP token requests over OAuth protocol",
		Units:      "requests/s",
		Fam:        "oauth",
		Ctx:        "adfs.oauth_token_requests_success",
		Priority:   prioADFSOauthTokenRequestsSuccess,
		Dims: module.Dims{
			{ID: "adfs_oauth_token_requests_success_total", Name: "success", Algo: module.Incremental},
		},
	}

	adfsPassiveRequestsChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_passive_requests",
		Title:      "Passive requests",
		Units:      "requests/s",
		Fam:        "requests",
		Ctx:        "adfs.passive_requests",
		Priority:   prioADFSPassiveRequests,
		Dims: module.Dims{
			{ID: "adfs_passive_requests_total", Name: "passive", Algo: module.Incremental},
		},
	}
	adfsPassportAuthenticationsChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_passport_authentications",
		Title:      "Microsoft Passport SSO authentications",
		Units:      "authentications/s",
		Fam:        "auth",
		Ctx:        "adfs.passport_authentications",
		Priority:   prioADFSPassportAuthentications,
		Dims: module.Dims{
			{ID: "adfs_passport_authentications_total", Name: "passport", Algo: module.Incremental},
		},
	}
	adfsPasswordChangeChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_password_change_requests",
		Title:      "Password change requests",
		Units:      "requests/s",
		Fam:        "requests",
		Ctx:        "adfs.password_change_requests",
		Priority:   prioADFSPasswordChangeRequests,
		Dims: module.Dims{
			{ID: "adfs_password_change_succeeded_total", Name: "success", Algo: module.Incremental},
			{ID: "adfs_password_change_failed_total", Name: "failed", Algo: module.Incremental},
		},
	}
	adfsSAMLPTokenRequestsChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_samlp_token_requests_success",
		Title:      "Successful RP token requests over SAML-P protocol",
		Units:      "requests/s",
		Fam:        "requests",
		Ctx:        "adfs.samlp_token_requests_success",
		Priority:   prioADFSSAMLPTokenRequests,
		Dims: module.Dims{
			{ID: "adfs_samlp_token_requests_success_total", Name: "success", Algo: module.Incremental},
		},
	}
	adfsWSTrustTokenRequestsSuccessChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_wstrust_token_requests_success",
		Title:      "Successful RP token requests over WS-Trust protocol",
		Units:      "requests/s",
		Fam:        "requests",
		Ctx:        "adfs.wstrust_token_requests_success",
		Priority:   prioADFSWSTrustTokenRequestsSuccess,
		Dims: module.Dims{
			{ID: "adfs_wstrust_token_requests_success_total", Name: "success", Algo: module.Incremental},
		},
	}
	adfsSSOAuthenticationsChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_sso_authentications",
		Title:      "SSO authentications",
		Units:      "authentications/s",
		Fam:        "auth",
		Ctx:        "adfs.sso_authentications",
		Priority:   prioADFSSSOAuthentications,
		Dims: module.Dims{
			{ID: "adfs_sso_authentications_success_total", Name: "success", Algo: module.Incremental},
			{ID: "adfs_sso_authentications_failure_total", Name: "failure", Algo: module.Incremental},
		},
	}
	adfsTokenRequestsChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_token_requests",
		Title:      "Token access requests",
		Units:      "requests/s",
		Fam:        "requests",
		Ctx:        "adfs.token_requests",
		Priority:   prioADFSTokenRequests,
		Dims: module.Dims{
			{ID: "adfs_token_requests_total", Name: "requests", Algo: module.Incremental},
		},
	}
	adfsUserPasswordAuthenticationsChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_userpassword_authentications",
		Title:      "AD U/P authentications",
		Units:      "authentications/s",
		Fam:        "auth",
		Ctx:        "adfs.userpassword_authentications",
		Priority:   prioADFSUserPasswordAuthentications,
		Dims: module.Dims{
			{ID: "adfs_sso_authentications_success_total", Name: "success", Algo: module.Incremental},
			{ID: "adfs_sso_authentications_failure_total", Name: "failure", Algo: module.Incremental},
		},
	}
	adfsWindowsIntegratedAuthenticationsChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_windows_integrated_authentications",
		Title:      "f Windows integrated authentications using Kerberos or NTLM",
		Units:      "authentications/s",
		Fam:        "auth",
		Ctx:        "adfs.windows_integrated_authentications",
		Priority:   prioADFSWindowsIntegratedAuthentications,
		Dims: module.Dims{
			{ID: "adfs_windows_integrated_authentications_total", Name: "authentications", Algo: module.Incremental},
		},
	}
	adfsWSFedTokenRequestsSuccessChart = module.Chart{
		OverModule: "adfs",
		ID:         "adfs_wsfed_token_requests_success",
		Title:      "Successful RP token requests over WS-Fed protocol",
		Units:      "requests/s",
		Fam:        "requests",
		Ctx:        "adfs.wsfed_token_requests_success",
		Priority:   prioADFSWSFedTokenRequestsSuccess,
		Dims: module.Dims{
			{ID: "adfs_wsfed_token_requests_success_total", Name: "success", Algo: module.Incremental},
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

// .NET
var (
	netFrameworkCLRCharts = module.Charts{
		netFrameworkCLRExceptionsThrown.Copy(),
		netFrameworkCLRExceptionsFilters.Copy(),
		netFrameworkCLRExceptionsFinallys.Copy(),
		netFrameworkCLRExceptionsThrowCatchDepth.Copy(),
		netFrameworkCLRInteropCOMCallableWrapper.Copy(),
		netFrameworkCLRInteropMarshalling.Copy(),
		netFrameworkCLRInteropStubsCreated.Copy(),
		netFrameworkCLRJITMethods.Copy(),
		netFrameworkCLRJITTime.Copy(),
		netFrameworkCLRJITStandardFailure.Copy(),
		netFrameworkCLRJITILByes.Copy(),
		netFrameworkCLRLoadingLoaderHeapSize.Copy(),
		netFrameworkCLRLoadingAppDomainsLoaded.Copy(),
		netFrameworkCLRLoadingAppDomainsUnloaded.Copy(),
		netFrameworkCLRLoadingAssembliesLoaded.Copy(),
		netFrameworkCLRLoadingClassesLoaded.Copy(),
		netFrameworkCLRLoadingClassLoadFailure.Copy(),
		netFrameworkCLRLockAndThreadsQueueLength.Copy(),
		netFrameworkCLRLockAndThreadsCurrentLogicalThreads.Copy(),
		netFrameworkCLRLockAndThreadsCurrentPhysicalThreads.Copy(),
		netFrameworkCLRLockAndThreadsRecognizedThreads.Copy(),
		netFrameworkCLRLockAndThreadsContentions.Copy(),
		netFrameworkCLRMemoryAllocatedBytes.Copy(),
		netFrameworkCLRMemoryFinalizationSurvivors.Copy(),
		netFrameworkCLRMemoryHeapSize.Copy(),
		netFrameworkCLRMemoryPromoted.Copy(),
		netFrameworkCLRMemoryNumberGCHandles.Copy(),
		netFrameworkCLRMemoryCollections.Copy(),
		netFrameworkCLRMemoryInducedGC.Copy(),
		netFrameworkCLRMemoryNumberPinnedObjects.Copy(),
		netFrameworkCLRMemoryNumberSinkBlocksInUse.Copy(),
		netFrameworkCLRMemoryCommitted.Copy(),
		netFrameworkCLRMemoryReserved.Copy(),
		netFrameworkCLRMemoryGCTime.Copy(),
		netFrameworkCLRRemotingChannels.Copy(),
		netFrameworkCLRContextBoundClassesLoaded.Copy(),
		netFrameworkCLRContextBoundObjects.Copy(),
		netFrameworkCLRContextProxies.Copy(),
		netFrameworkCLRContexts.Copy(),
		netFrameworkCLRRemotingCalls.Copy(),
		netFrameworkCLRSecurityLinkTimeChecks.Copy(),
		netFrameworkCLRSecurityChecksTime.Copy(),
		netFrameworkCLRSecurityStackWalkDepth.Copy(),
		netFrameworkCLRSecurityRuntimeChecks.Copy(),
	}

	// Exceptions
	netFrameworkCLRExceptionsThrown = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_exception_thrown",
		Title:      "Number of exceptions thrown since start.",
		Units:      "exceptions",
		Fam:        "net_framework",
		Ctx:        "net_framework.exception_thrown",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRExceptionsThrown,
	}
	netFrameworkCLRExceptionsFilters = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_exception_filters",
		Title:      "Number of exceptions filter executed.",
		Units:      "exceptions",
		Fam:        "net_framework",
		Ctx:        "net_framework.exception_filters",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRExceptionsFilters,
	}
	netFrameworkCLRExceptionsFinallys = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_exception_finally",
		Title:      "Number of finally blocks executed.",
		Units:      "blocks",
		Fam:        "net_framework",
		Ctx:        "net_framework.exception_finally",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRExceptionsFinallys,
	}
	netFrameworkCLRExceptionsThrowCatchDepth = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_exception_throw_catch_depth",
		Title:      "Number of stack frames transversed.",
		Units:      "stack frames",
		Fam:        "net_framework",
		Ctx:        "net_framework.exception_thrown",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRExceptionsThrowCatchDepth,
	}

	// Interop
	netFrameworkCLRInteropCOMCallableWrapper = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrinterop_com_callable_wrapper",
		Title:      "Number of COM callable wrappers(CCW).",
		Units:      "ccw",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrinterop_com_callable_wrapper",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRInteropCOMCalableWrapper,
	}
	netFrameworkCLRInteropMarshalling = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrinterop_marshalling",
		Title:      "Number of times arguments or return value have been marshaled.",
		Units:      "values marshaled",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrinterop_marshalling",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRInteropMarshalling,
	}
	netFrameworkCLRInteropStubsCreated = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrinterop_stubs_created",
		Title:      "Number of stubs created.",
		Units:      "stubs",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrinterop_stubs_created",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRInteropStubsCreated,
	}

	// JIT
	netFrameworkCLRJITMethods = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrjit_methods",
		Title:      "Number of JIT methods compiled.",
		Units:      "JIT methods",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrjit_methods",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRJITMethods,
	}
	netFrameworkCLRJITTime = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrjit_time",
		Title:      "Percentage time spent in JIT compilation.",
		Units:      "percentage",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrjit_time",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRJITTime,
	}
	netFrameworkCLRJITStandardFailure = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrjit_standard_failure",
		Title:      "Number of methods the JIT compiler has failed.",
		Units:      "Failures",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrjit_standard_failure",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRJITStandardFailures,
	}
	netFrameworkCLRJITILByes = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrjit_il_bytes",
		Title:      "Number of Microsoft intermediate language (MSIL) bytes compiled.",
		Units:      "bytes",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrjit_il_bytes",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRJITILBytes,
	}

	// Loading
	netFrameworkCLRLoadingLoaderHeapSize = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrloading_loader_heap_size",
		Title:      "Bytes committed by class loader.",
		Units:      "bytes",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrloading_loader_heap_size",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRLoadingLoaderHeapSize,
	}
	netFrameworkCLRLoadingAppDomainsLoaded = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrloading_appdomains_loaded",
		Title:      "Application domain loaded since start.",
		Units:      "application domain",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrloading_appdomains_loaded",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRLoadingAppDomainsLoaded,
	}
	netFrameworkCLRLoadingAppDomainsUnloaded = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrloading_appdomains_unloaded",
		Title:      "Application domain unloaded since start.",
		Units:      "application domain",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrloading_appdomains_unloaded",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRLoadingAppDomainsUnloaded,
	}
	netFrameworkCLRLoadingAssembliesLoaded = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrloading_assemblies_loaded",
		Title:      "Assemblies loaded since start.",
		Units:      "assemblies",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrloading_assemblies_loaded",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRLoadingAssembliesLoaded,
	}
	netFrameworkCLRLoadingClassesLoaded = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrloading_classes_loaded",
		Title:      "Classes loaded in all assemblies.",
		Units:      "classes",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrloading_classes_loaded",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRLoadingClassesLoaded,
	}
	netFrameworkCLRLoadingClassLoadFailure = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrloading_class_load_failure",
		Title:      "Classes that have failed to load.",
		Units:      "classes",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrloading_class_load_failure",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRLoadingClassLoadFailure,
	}

	// Lock and Threads
	netFrameworkCLRLockAndThreadsQueueLength = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrlockandthreads_queue_length",
		Title:      "Threads waiting for lock.",
		Units:      "threads",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrlockandthreads_queue_length",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRLockAndThreadsQueueLength,
	}
	netFrameworkCLRLockAndThreadsCurrentLogicalThreads = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrlockandthreads_current_logical_threads",
		Title:      "Number of running and stopped threads.",
		Units:      "threads",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrlockandthreads_current_logical_threads",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRLockAndThreadsCurrentLogicalThreads,
	}
	netFrameworkCLRLockAndThreadsCurrentPhysicalThreads = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrlockandthreads_current_physical_threads",
		Title:      "Number of running and stopped threads.",
		Units:      "threads",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrlockandthreads_current_physical_threads",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRLockAndThreadsCurrentPhysicalThreads,
	}
	netFrameworkCLRLockAndThreadsRecognizedThreads = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrlockandthreads_recognized_threads",
		Title:      "Number of running and stopped threads.",
		Units:      "threads",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrlockandthreads_recognized_threads",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRLockAndThreadsRecognizedThreads,
	}
	netFrameworkCLRLockAndThreadsContentions = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrlockandthreads_contentions",
		Title:      "Number of fails to acquire managed lock.",
		Units:      "threads",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrlockandthreads_contentions",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRLockAndThreadsContentions,
	}

	// Memory
	netFrameworkCLRMemoryAllocatedBytes = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrmemory_allocated_bytes",
		Title:      "Total of bytes allocated on the garbage collection.",
		Units:      "bytes",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrmemory_allocated_bytes",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRMemoryAllocatedBytes,
	}
	netFrameworkCLRMemoryFinalizationSurvivors = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrmemory_finalization_survivors",
		Title:      "Number of garbage-collected objects that survived.",
		Units:      "objects",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrmemory_finalization_survivors",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRMemoryFinalizationSurvivors,
	}
	netFrameworkCLRMemoryHeapSize = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrmemory_heap_size",
		Title:      "Maximum bytes can be allocated.",
		Units:      "bytes",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrmemory_heap_size",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRMemoryHeapSize,
	}
	netFrameworkCLRMemoryPromoted = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrmemory_promoted",
		Title:      "Bytes promoted from the generation to next.",
		Units:      "bytes",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrmemory_promoted",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRMemoryPromoted,
	}
	netFrameworkCLRMemoryNumberGCHandles = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrmemory_number_gc_handles",
		Title:      "Number of Garbage Collection Handle used.",
		Units:      "handles",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrmemory_number_gc_handles",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRMemoryNumberGCHandles,
	}
	netFrameworkCLRMemoryCollections = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrmemory_collections",
		Title:      "Number of times Garbage Collection (GC) happens.",
		Units:      "GC",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrmemory_collections",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRMemoryCollections,
	}
	netFrameworkCLRMemoryInducedGC = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrmemory_induced_gc",
		Title:      "Peak number of times garbage collection was performed.",
		Units:      "GC",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrmemory_induced_gc",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRMemoryInducedGC,
	}
	netFrameworkCLRMemoryNumberPinnedObjects = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrmemory_number_pinned_objects",
		Title:      "Pinned objects in the last garbage collection.",
		Units:      "objects",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrmemory_number_pinned_objects",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRMemoryNumberPinnedObjects,
	}
	netFrameworkCLRMemoryNumberSinkBlocksInUse = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrmemory_number_sink_blocks_in_use",
		Title:      "Number of synchronization blocks in use.",
		Units:      "blocks",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrmemory_number_sink_blocks_in_use",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRMemoryNumberSinkBlocksInUse,
	}
	netFrameworkCLRMemoryCommitted = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrmemory_committed",
		Title:      "Amount of virtual memory currently committed.",
		Units:      "bytes",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrmemory_committed",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRMemoryCommitted,
	}
	netFrameworkCLRMemoryReserved = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrmemory_reserved",
		Title:      "Amount of virtual memory currently reserved.",
		Units:      "bytes",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrmemory_reserved",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRMemoryReserved,
	}
	netFrameworkCLRMemoryGCTime = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrmemory_gc_time",
		Title:      "Amount of virtual memory currently reserved.",
		Units:      "percentage",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrmemory_gc_time",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRMemoryGCTime,
	}

	// Remoting
	netFrameworkCLRRemotingChannels = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrremoting_channels",
		Title:      "Remoting channels registered.",
		Units:      "channels",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrremoting_channels",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRRemotingChannels,
	}
	netFrameworkCLRContextBoundClassesLoaded = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrremoting_context_bound_classes_loaded",
		Title:      "Context-bound classes loaded.",
		Units:      "context-bound",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrremoting_context_bound_classes_loaded",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRRemotingContextBoundClassesLoaded,
	}
	netFrameworkCLRContextBoundObjects = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrremoting_context_bound_objects",
		Title:      "Total Context-bound objects allocated.",
		Units:      "objects",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrremoting_context_bound_objects",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRRemotingContextBoundObjects,
	}
	netFrameworkCLRContextProxies = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrremoting_context_proxies",
		Title:      "Total remoting proxy objects.",
		Units:      "objects",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrremoting_context_proxies",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRRemotingContextProxies,
	}
	netFrameworkCLRContexts = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrremoting_contexts",
		Title:      "Total of remoting contexts.",
		Units:      "contexts",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrremoting_contexts",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRRemotingContexts,
	}
	netFrameworkCLRRemotingCalls = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrremoting_calls",
		Title:      "Total Remote Procedure Calls (RPC) invoked.",
		Units:      "RPC",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrremoting_calls",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRRemotingRemoteCalls,
	}

	// Security
	netFrameworkCLRSecurityLinkTimeChecks = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrsecurity_link_time_checks",
		Title:      "Number of Link-time code access security.",
		Units:      "access",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrsecurity_link_time_checks",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRSecurityLinkTimeChecks,
	}
	netFrameworkCLRSecurityChecksTime = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrsecurity_checks_time",
		Title:      "Percentage of time to access security checks.",
		Units:      "percentage",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrsecurity_checks_time",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRSecurityRTChecksTime,
	}
	netFrameworkCLRSecurityStackWalkDepth = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrsecurity_stack_walk_depth",
		Title:      "Depth of stack.",
		Units:      "depth",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrsecurity_stack_walk_depth",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRSecurityStackWalkDepth,
	}
	netFrameworkCLRSecurityRuntimeChecks = module.Chart{
		OverModule: "NET Framework",
		ID:         "net_framework_clrsecurity_runtime_checks",
		Title:      "Total number of runtime code access.",
		Units:      "access",
		Fam:        "net_framework",
		Ctx:        "net_framework.clrsecurity_checks_time",
		Type:       module.Stacked,
		Priority:   prioNETFrameworkCLRSecurityRuntimeChecks,
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

func (w *WMI) addNetFrameworkCRLExceptions() {
	charts := netFrameworkCLRCharts.Copy()

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

func (w *WMI) addProcessToNetFrameworkCharts(procID string) {
	for _, chart := range *w.Charts() {
		var dim *module.Dim
		switch chart.ID {
		case netFrameworkCLRExceptionsThrown.ID:
			id := fmt.Sprintf("netframework_clrexceptions_%s_thrown", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case netFrameworkCLRExceptionsFilters.ID:
			id := fmt.Sprintf("netframework_clrexceptions_%s_filters", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case netFrameworkCLRExceptionsFinallys.ID:
			id := fmt.Sprintf("netframework_clrexceptions_%s_finallys", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case netFrameworkCLRExceptionsThrowCatchDepth.ID:
			id := fmt.Sprintf("netframework_clrexceptions_%s_throw_catch_depth", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case netFrameworkCLRInteropCOMCallableWrapper.ID:
			id := fmt.Sprintf("netframework_clrinterop_%s_com_callable_wrappers", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case netFrameworkCLRInteropMarshalling.ID:
			id := fmt.Sprintf("netframework_clrinterop_%s_marshalling", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case netFrameworkCLRInteropStubsCreated.ID:
			id := fmt.Sprintf("netframework_clrinterop_%s_stubs_created", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case netFrameworkCLRJITMethods.ID:
			id := fmt.Sprintf("netframework_clrjit_%s_methods", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case netFrameworkCLRJITTime.ID:
			id := fmt.Sprintf("netframework_clrjit_%s_time", procID)
			dim = &module.Dim{ID: id, Name: procID}
		case netFrameworkCLRJITStandardFailure.ID:
			id := fmt.Sprintf("netframework_clrjit_%s_standard_failure", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case netFrameworkCLRJITILByes.ID:
			id := fmt.Sprintf("netframework_clrjit_%s_il_bytes", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case netFrameworkCLRLoadingLoaderHeapSize.ID:
			id := fmt.Sprintf("netframework_clrloading_%s_loader_heap_size", procID)
			dim = &module.Dim{ID: id, Name: procID}
		case netFrameworkCLRLoadingAppDomainsLoaded.ID:
			id := fmt.Sprintf("netframework_clrloading_%s_app_domains_loaded", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case netFrameworkCLRLoadingAppDomainsUnloaded.ID:
			id := fmt.Sprintf("netframework_clrloading_%s_app_domains_unloaded", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case netFrameworkCLRLoadingAssembliesLoaded.ID:
			id := fmt.Sprintf("netframework_clrloading_%s_assemblies_loaded", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case netFrameworkCLRLoadingClassesLoaded.ID:
			id := fmt.Sprintf("netframework_clrloading_%s_classes_loaded", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case netFrameworkCLRLoadingClassLoadFailure.ID:
			id := fmt.Sprintf("netframework_clrloading_%s_class_load_failure", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case netFrameworkCLRLockAndThreadsQueueLength.ID:
			id := fmt.Sprintf("netframework_clrlockandthreads_%s_queue_length", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case netFrameworkCLRLockAndThreadsCurrentLogicalThreads.ID:
			id := fmt.Sprintf("netframework_clrlockandthreads_%s_current_logical_threads", procID)
			dim = &module.Dim{ID: id, Name: procID}
		case netFrameworkCLRLockAndThreadsCurrentPhysicalThreads.ID:
			id := fmt.Sprintf("netframework_clrlockandthreads_%s_current_physical_threads", procID)
			dim = &module.Dim{ID: id, Name: procID}
		case netFrameworkCLRLockAndThreadsRecognizedThreads.ID:
			id := fmt.Sprintf("netframework_clrlockandthreads_%s_recognized_threads", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case netFrameworkCLRLockAndThreadsContentions.ID:
			id := fmt.Sprintf("netframework_clrlockandthreads_%s_contentions", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case netFrameworkCLRMemoryAllocatedBytes.ID:
			id := fmt.Sprintf("netframework_clrmemory_%s_allocated_bytes_total", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case netFrameworkCLRMemoryFinalizationSurvivors.ID:
			id := fmt.Sprintf("netframework_clrmemory_%s_finalization_survivors", procID)
			dim = &module.Dim{ID: id, Name: procID}
		case netFrameworkCLRMemoryHeapSize.ID:
			id := fmt.Sprintf("netframework_clrmemory_%s_heap_size", procID)
			dim = &module.Dim{ID: id, Name: procID}
		case netFrameworkCLRMemoryPromoted.ID:
			id := fmt.Sprintf("netframework_clrmemory_%s_promoted", procID)
			dim = &module.Dim{ID: id, Name: procID}
		case netFrameworkCLRMemoryNumberGCHandles.ID:
			id := fmt.Sprintf("netframework_clrmemory_%s_gc_handles", procID)
			dim = &module.Dim{ID: id, Name: procID}
		case netFrameworkCLRMemoryCollections.ID:
			id := fmt.Sprintf("netframework_clrmemory_%s_collection", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case netFrameworkCLRMemoryInducedGC.ID:
			id := fmt.Sprintf("netframework_clrmemory_%s_induced_gc_total", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case netFrameworkCLRMemoryNumberPinnedObjects.ID:
			id := fmt.Sprintf("netframework_clrmemory_%s_pinned_objects", procID)
			dim = &module.Dim{ID: id, Name: procID}
		case netFrameworkCLRMemoryNumberSinkBlocksInUse.ID:
			id := fmt.Sprintf("netframework_clrmemory_%s_sink_block_in_use", procID)
			dim = &module.Dim{ID: id, Name: procID}
		case netFrameworkCLRMemoryCommitted.ID:
			id := fmt.Sprintf("netframework_clrmemory_%s_committed", procID)
			dim = &module.Dim{ID: id, Name: procID}
		case netFrameworkCLRMemoryReserved.ID:
			id := fmt.Sprintf("netframework_clrmemory_%s_reserved", procID)
			dim = &module.Dim{ID: id, Name: procID}
		case netFrameworkCLRMemoryGCTime.ID:
			id := fmt.Sprintf("netframework_clrmemory_%s_gc_time", procID)
			dim = &module.Dim{ID: id, Name: procID}
		case netFrameworkCLRRemotingChannels.ID:
			id := fmt.Sprintf("netframework_clrremoting_%s_channels", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case netFrameworkCLRContextBoundClassesLoaded.ID:
			id := fmt.Sprintf("netframework_clrremoting_%s_context_bound_classes_loaded", procID)
			dim = &module.Dim{ID: id, Name: procID}
		case netFrameworkCLRContextBoundObjects.ID:
			id := fmt.Sprintf("netframework_clrremoting_%s_context_bound_objects", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case netFrameworkCLRContextProxies.ID:
			id := fmt.Sprintf("netframework_clrremoting_%s_context_proxies", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case netFrameworkCLRContexts.ID:
			id := fmt.Sprintf("netframework_clrremoting_%s_contexts", procID)
			dim = &module.Dim{ID: id, Name: procID}
		case netFrameworkCLRRemotingCalls.ID:
			id := fmt.Sprintf("netframework_clrremoting_%s_calls", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case netFrameworkCLRSecurityLinkTimeChecks.ID:
			id := fmt.Sprintf("netframework_clrsecurity_%s_link_time_checks", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		case netFrameworkCLRSecurityChecksTime.ID:
			id := fmt.Sprintf("netframework_clrsecurity_%s_checks_time", procID)
			dim = &module.Dim{ID: id, Name: procID}
		case netFrameworkCLRSecurityStackWalkDepth.ID:
			id := fmt.Sprintf("netframework_clrsecurity_%s_stack_walk_depth", procID)
			dim = &module.Dim{ID: id, Name: procID}
		case netFrameworkCLRSecurityRuntimeChecks.ID:
			id := fmt.Sprintf("netframework_clrsecurity_%s_runtime_checks", procID)
			dim = &module.Dim{ID: id, Name: procID, Algo: module.Incremental}
		default:
			dim = nil
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

func (w *WMI) removeProcessFromNetFrameworkCharts(procID string) {
	for _, chart := range *w.Charts() {
		var id string
		switch chart.ID {
		case netFrameworkCLRExceptionsThrown.ID:
			id = fmt.Sprintf("netframework_clrexceptions_%s_thrown", procID)
		case netFrameworkCLRExceptionsFilters.ID:
			id = fmt.Sprintf("netframework_clrexceptions_%s_filters", procID)
		case netFrameworkCLRExceptionsFinallys.ID:
			id = fmt.Sprintf("netframework_clrexceptions_%s_finallys", procID)
		case netFrameworkCLRExceptionsThrowCatchDepth.ID:
			id = fmt.Sprintf("netframework_clrexceptions_%s_throw_catch_depth", procID)
		case netFrameworkCLRInteropCOMCallableWrapper.ID:
			id = fmt.Sprintf("netframework_clrinterop_%s_com_callable_wrappers", procID)
		case netFrameworkCLRInteropMarshalling.ID:
			id = fmt.Sprintf("netframework_clrinterop_%s_marshalling", procID)
		case netFrameworkCLRInteropStubsCreated.ID:
			id = fmt.Sprintf("netframework_clrinterop_%s_stubs_created", procID)
		case netFrameworkCLRJITMethods.ID:
			id = fmt.Sprintf("netframework_clrjit_%s_methods", procID)
		case netFrameworkCLRJITTime.ID:
			id = fmt.Sprintf("netframework_clrjit_%s_time", procID)
		case netFrameworkCLRJITStandardFailure.ID:
			id = fmt.Sprintf("netframework_clrjit_%s_standard_failure", procID)
		case netFrameworkCLRJITILByes.ID:
			id = fmt.Sprintf("netframework_clrjit_%s_il_bytes", procID)
		case netFrameworkCLRLoadingLoaderHeapSize.ID:
			id = fmt.Sprintf("netframework_clrloading_%s_loader_heap_size", procID)
		case netFrameworkCLRLoadingAppDomainsLoaded.ID:
			id = fmt.Sprintf("netframework_clrloading_%s_app_domains_loaded", procID)
		case netFrameworkCLRLoadingAppDomainsUnloaded.ID:
			id = fmt.Sprintf("netframework_clrloading_%s_app_domains_unloaded", procID)
		case netFrameworkCLRLoadingAssembliesLoaded.ID:
			id = fmt.Sprintf("netframework_clrloading_%s_assemblies_loaded", procID)
		case netFrameworkCLRLoadingClassesLoaded.ID:
			id = fmt.Sprintf("netframework_clrloading_%s_classes_loaded", procID)
		case netFrameworkCLRLoadingClassLoadFailure.ID:
			id = fmt.Sprintf("netframework_clrloading_%s_class_load_failure", procID)
		case netFrameworkCLRLockAndThreadsQueueLength.ID:
			id = fmt.Sprintf("netframework_clrlockandthreads_%s_queue_length", procID)
		case netFrameworkCLRLockAndThreadsCurrentLogicalThreads.ID:
			id = fmt.Sprintf("netframework_clrlockandthreads_%s_current_logical_threads", procID)
		case netFrameworkCLRLockAndThreadsCurrentPhysicalThreads.ID:
			id = fmt.Sprintf("netframework_clrlockandthreads_%s_current_physical_threads", procID)
		case netFrameworkCLRLockAndThreadsRecognizedThreads.ID:
			id = fmt.Sprintf("netframework_clrlockandthreads_%s_recognized_threads", procID)
		case netFrameworkCLRLockAndThreadsContentions.ID:
			id = fmt.Sprintf("netframework_clrlockandthreads_%s_contentions", procID)
		case netFrameworkCLRMemoryAllocatedBytes.ID:
			id = fmt.Sprintf("netframework_clrmemory_%s_allocated_bytes_total", procID)
		case netFrameworkCLRMemoryFinalizationSurvivors.ID:
			id = fmt.Sprintf("netframework_clrmemory_%s_finalization_survivors", procID)
		case netFrameworkCLRMemoryHeapSize.ID:
			id = fmt.Sprintf("netframework_clrmemory_%s_heap_size", procID)
		case netFrameworkCLRMemoryPromoted.ID:
			id = fmt.Sprintf("netframework_clrmemory_%s_promoted", procID)
		case netFrameworkCLRMemoryNumberGCHandles.ID:
			id = fmt.Sprintf("netframework_clrmemory_%s_gc_handles", procID)
		case netFrameworkCLRMemoryCollections.ID:
			id = fmt.Sprintf("netframework_clrmemory_%s_collection", procID)
		case netFrameworkCLRMemoryInducedGC.ID:
			id = fmt.Sprintf("netframework_clrmemory_%s_induced_gc_total", procID)
		case netFrameworkCLRMemoryNumberPinnedObjects.ID:
			id = fmt.Sprintf("netframework_clrmemory_%s_pinned_objects", procID)
		case netFrameworkCLRMemoryNumberSinkBlocksInUse.ID:
			id = fmt.Sprintf("netframework_clrmemory_%s_sink_block_in_use", procID)
		case netFrameworkCLRMemoryCommitted.ID:
			id = fmt.Sprintf("netframework_clrmemory_%s_committed", procID)
		case netFrameworkCLRMemoryReserved.ID:
			id = fmt.Sprintf("netframework_clrmemory_%s_reserved", procID)
		case netFrameworkCLRMemoryGCTime.ID:
			id = fmt.Sprintf("netframework_clrmemory_%s_gc_time", procID)
		case netFrameworkCLRRemotingChannels.ID:
			id = fmt.Sprintf("netframework_clrremoting_%s_channels", procID)
		case netFrameworkCLRContextBoundClassesLoaded.ID:
			id = fmt.Sprintf("netframework_clrremoting_%s_context_bound_classes_loaded", procID)
		case netFrameworkCLRContextBoundObjects.ID:
			id = fmt.Sprintf("netframework_clrremoting_%s_context_bound_objects", procID)
		case netFrameworkCLRContextProxies.ID:
			id = fmt.Sprintf("netframework_clrremoting_%s_context_proxies", procID)
		case netFrameworkCLRContexts.ID:
			id = fmt.Sprintf("netframework_clrremoting_%s_contexts", procID)
		case netFrameworkCLRRemotingCalls.ID:
			id = fmt.Sprintf("netframework_clrremoting_%s_calls", procID)
		case netFrameworkCLRSecurityLinkTimeChecks.ID:
			id = fmt.Sprintf("netframework_clrsecurity_%s_link_time_checks", procID)
		case netFrameworkCLRSecurityChecksTime.ID:
			id = fmt.Sprintf("netframework_clrsecurity_%s_checks_time", procID)
		case netFrameworkCLRSecurityStackWalkDepth.ID:
			id = fmt.Sprintf("netframework_clrsecurity_%s_stack_walk_depth", procID)
		case netFrameworkCLRSecurityRuntimeChecks.ID:
			id = fmt.Sprintf("netframework_clrsecurity_%s_runtime_checks", procID)
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

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
	// Dim is an alias for module.Dim
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
	prioSystemThreads
	prioSystemUptime

	prioLogonSessions

	prioCollectionDuration
	prioCollectionStatus
)

func cpuCharts() Charts {
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

func cpuCoreCharts() Charts {
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

func memCharts() Charts {
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

func nicCharts() Charts {
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

func diskCharts() Charts {
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

func osCharts() Charts {
	return Charts{
		osProcessesChart.Copy(),
	}
}

var (
	osProcessesChart = Chart{
		ID:       "system_processes",
		Title:    "Processes",
		Units:    "number",
		Fam:      "system",
		Ctx:      "wmi.system_processes",
		Priority: prioOSProcesses,
		Dims: Dims{
			{ID: "os_processes", Name: "processes"},
		},
		Vars: Vars{
			{ID: "os_processes_limit"},
		},
	}
)

func systemCharts() Charts {
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

func logonCharts() Charts {
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

func collectionCharts() *Charts {
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
		if err := w.Charts().Add(cpuCharts()...); err != nil {
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
	if err := w.Charts().Add(memCharts()...); err != nil {
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
	if err := w.Charts().Add(systemCharts()...); err != nil {
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
	if err := w.Charts().Add(osCharts()...); err != nil {
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
	if err := w.Charts().Add(logonCharts()...); err != nil {
		w.Warning(err)
	}
}

func addCPUCoreCharts(charts *Charts, coreID string) error {
	for _, chart := range cpuCoreCharts() {
		chart = newChartFromTemplate(*chart, coreID)
		if err := charts.Add(chart); err != nil {
			return err
		}
	}
	return nil
}

func addNICCharts(charts *Charts, nicID string) error {
	for _, chart := range nicCharts() {
		chart = newChartFromTemplate(*chart, nicID)
		if err := charts.Add(chart); err != nil {
			return err
		}
	}
	return nil
}

func addLogicalDiskCharts(charts *Charts, diskID string) error {
	for _, chart := range diskCharts() {
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

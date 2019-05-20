package wmi

import (
	"fmt"
	"strings"

	"github.com/netdata/go-orchestrator"
	"github.com/netdata/go-orchestrator/module"
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
	defaultPriority     = orchestrator.DefaultJobPriority
	memoryPriority      = defaultPriority + 20
	netPriority         = defaultPriority + 40
	logicalDiskPriority = defaultPriority + 60
	systemPriority      = defaultPriority + 80
)

var charts = Charts{
	{
		ID:       "collector_duration",
		Title:    "Duration of a Collector",
		Units:    "ms",
		Fam:      "collection",
		Ctx:      "cpu.collector_duration",
		Priority: defaultPriority + 200, // last chart
		// Dims will be added during collection
	},
}

var (
	cpuCharts = Charts{
		{
			ID:    "cpu_utilization_total",
			Title: "Total CPU Utilization (all cores)",
			Units: "percentage",
			Fam:   "cpu",
			Ctx:   "wmi.cpu_utilization_total",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "cpu_idle", Name: "idle", Algo: module.PercentOfIncremental, Div: 1000, DimOpts: Opts{Hidden: true}},
				{ID: "cpu_dpc", Name: "dpc", Algo: module.PercentOfIncremental, Div: 1000},
				{ID: "cpu_user", Name: "user", Algo: module.PercentOfIncremental, Div: 1000},
				{ID: "cpu_privileged", Name: "privileged", Algo: module.PercentOfIncremental, Div: 1000},
				{ID: "cpu_interrupt", Name: "interrupt", Algo: module.PercentOfIncremental, Div: 1000},
			},
		},
		{
			ID:    "cpu_dpcs",
			Title: "Received and Serviced Deferred Procedure Calls (DPC)",
			Units: "dpc/s",
			Fam:   "cpu",
			Ctx:   "wmi.cpu_dpcs",
			Type:  module.Stacked,
			// Dims will be added during collection
		},
		{
			ID:    "cpu_interrupts",
			Title: "Received and Serviced Hardware Interrupts",
			Units: "interrupts/s",
			Fam:   "cpu",
			Ctx:   "wmi.cpu_interrupts",
			Type:  module.Stacked,
			// Dims will be added during collection
		},
	}

	// Per core charts
	cpuCoreUsageChart = Chart{

		ID:    "core_%s_cpu_utilization",
		Title: "Core%s CPU Utilization",
		Units: "percentage",
		Fam:   "cpu",
		Ctx:   "wmi.cpu_utilization",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "cpu_core_%s_idle", Name: "idle", Algo: module.PercentOfIncremental, Div: 1000, DimOpts: Opts{Hidden: true}},
			{ID: "cpu_core_%s_dpc", Name: "dpc", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_core_%s_user", Name: "user", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_core_%s_privileged", Name: "privileged", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_core_%s_interrupt", Name: "interrupt", Algo: module.PercentOfIncremental, Div: 1000},
		},
	}
	cpuCoreCStateChart = Chart{
		ID:    "core_%s_cpu_cstate",
		Title: "Core%s Time Spent in Low-Power Idle State",
		Units: "percentage",
		Fam:   "cpu",
		Ctx:   "wmi.cpu_cstate",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "cpu_core_%s_c1", Name: "c1", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_core_%s_c2", Name: "c2", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_core_%s_c3", Name: "c3", Algo: module.PercentOfIncremental, Div: 1000},
		},
	}
)

var (
	memoryCharts = Charts{
		{
			ID:       "memory_utilization",
			Title:    "Memory Utilization",
			Units:    "KiB",
			Fam:      "mem",
			Ctx:      "wmi.memory_utilization",
			Type:     module.Stacked,
			Priority: memoryPriority,
			Dims: Dims{
				{ID: "memory_available_bytes", Name: "available", Div: 1000 * 1024},
				{ID: "memory_used_bytes", Name: "used", Div: 1000 * 1024},
			},
			Vars: Vars{
				{ID: "os_visible_memory_bytes"},
			},
		},
		{
			ID:       "memory_page_faults",
			Title:    "Memory Page Faults",
			Units:    "events/s",
			Fam:      "mem",
			Ctx:      "wmi.memory_page_faults",
			Priority: memoryPriority,
			Dims: Dims{
				{ID: "memory_page_faults_total", Name: "page faults", Algo: module.Incremental, Div: 1000},
			},
		},
		{
			ID:       "memory_swap_utilization",
			Title:    "Swap Utilization",
			Units:    "KiB",
			Fam:      "mem",
			Ctx:      "wmi.memory_swap_utilization",
			Type:     module.Stacked,
			Priority: memoryPriority,
			Dims: Dims{
				{ID: "memory_not_committed_bytes", Name: "available", Div: 1000 * 1024},
				{ID: "memory_committed_bytes", Name: "used", Div: 1000 * 1024},
			},
			Vars: Vars{
				{ID: "memory_commit_limit"},
			},
		},
		{
			ID:       "memory_swap_operations",
			Title:    "Swap Operations",
			Units:    "operations/s",
			Fam:      "mem",
			Ctx:      "wmi.memory_swap_operations",
			Type:     module.Area,
			Priority: memoryPriority,
			Dims: Dims{
				{ID: "memory_swap_page_reads_total", Name: "read", Algo: module.Incremental, Div: 1000},
				{ID: "memory_swap_page_writes_total", Name: "write", Algo: module.Incremental, Div: -1000},
			},
		},
		{
			ID:       "memory_swap_pages",
			Title:    "Swap Pages",
			Units:    "pages/s",
			Fam:      "mem",
			Ctx:      "wmi.memory_swap_pages",
			Type:     module.Area,
			Priority: memoryPriority,
			Dims: Dims{
				{ID: "memory_swap_pages_read_total", Name: "read", Algo: module.Incremental, Div: 1000},
				{ID: "memory_swap_pages_written_total", Name: "written", Algo: module.Incremental, Div: -1000},
			},
		},
		{
			ID:       "memory_cached",
			Title:    "Cached",
			Units:    "KiB",
			Fam:      "mem",
			Ctx:      "wmi.memory_cached",
			Priority: memoryPriority,
			Dims: Dims{
				{ID: "memory_cache_total", Name: "cached", Div: 1000 * 1024},
			},
		},
		{
			ID:       "memory_cache_faults",
			Title:    "Cache Faults",
			Units:    "events/s",
			Fam:      "mem",
			Ctx:      "wmi.memory_cache_faults",
			Priority: memoryPriority,
			Dims: Dims{
				{ID: "memory_cache_faults_total", Name: "cache faults", Algo: module.Incremental, Div: 1000},
			},
		},
		{
			ID:       "memory_system_pool",
			Title:    "System Memory Pool",
			Units:    "KiB",
			Fam:      "mem",
			Ctx:      "wmi.memory_system_pool",
			Type:     module.Stacked,
			Priority: memoryPriority,
			Dims: Dims{
				{ID: "memory_pool_paged_bytes", Name: "paged", Div: 1000 * 1024},
				{ID: "memory_pool_nonpaged_bytes_total", Name: "non-paged", Div: 1000 * 1024},
			},
		},
	}
)

var netNICCharts = Charts{
	{
		ID:       "nic_%s_bandwidth",
		Title:    "Bandwidth %s",
		Units:    "kilobits/s",
		Fam:      "net",
		Ctx:      "wmi.net_bandwidth",
		Type:     module.Area,
		Priority: netPriority,
		Dims: Dims{
			{ID: "net_%s_bytes_received", Name: "received", Algo: module.Incremental, Div: 1000 * 125},
			{ID: "net_%s_bytes_sent", Name: "sent", Algo: module.Incremental, Div: -1000 * 125},
		},
		Vars: Vars{
			{ID: "net_%s_current_bandwidth"},
		},
	},
	{
		ID:       "nic_%s_packets",
		Title:    "Packets %s",
		Units:    "packets/s",
		Fam:      "net",
		Ctx:      "wmi.net_packets",
		Type:     module.Area,
		Priority: netPriority + 1,
		Dims: Dims{
			{ID: "net_%s_packets_received_total", Name: "received", Algo: module.Incremental, Div: 1000},
			{ID: "net_%s_packets_sent_total", Name: "sent", Algo: module.Incremental, Div: -1000},
		},
	},
	{
		ID:       "nic_%s_errors",
		Title:    "Errors %s",
		Units:    "errors/s",
		Fam:      "net",
		Ctx:      "wmi.net_errors",
		Type:     module.Area,
		Priority: netPriority + 2,
		Dims: Dims{
			{ID: "net_%s_packets_received_errors", Name: "inbound", Algo: module.Incremental, Div: 1000},
			{ID: "net_%s_packets_outbound_errors", Name: "outbound", Algo: module.Incremental, Div: -1000},
		},
	},
	{
		ID:       "nic_%s_discarded",
		Title:    "Discards %s",
		Units:    "discards/s",
		Fam:      "net",
		Ctx:      "wmi.net_discarded",
		Type:     module.Area,
		Priority: netPriority + 3,
		Dims: Dims{
			{ID: "net_%s_packets_received_discarded", Name: "inbound", Algo: module.Incremental, Div: 1000},
			{ID: "net_%s_packets_outbound_discarded", Name: "outbound", Algo: module.Incremental, Div: -1000},
		},
	},
}

var systemCharts = Charts{
	{
		ID:       "system_uptime",
		Title:    "System Uptime",
		Units:    "seconds",
		Fam:      "uptime",
		Ctx:      "wmi.system_uptime",
		Priority: systemPriority,
		Dims: Dims{
			{ID: "system_up_time", Name: "time"},
		},
	},
}

var logicalDiskCharts = Charts{
	{
		ID:       "logical_disk_%s_utilization",
		Title:    "Utilization Disk %s",
		Units:    "KiB",
		Fam:      "disk",
		Ctx:      "wmi.logical_disk_utilization",
		Type:     module.Stacked,
		Priority: logicalDiskPriority,
		Dims: Dims{
			{ID: "logical_disk_%s_free_space", Name: "free", Div: 1000 * 1024},
			{ID: "logical_disk_%s_used_space", Name: "used", Div: 1000 * 1024},
		},
	},
	{
		ID:       "logical_disk_%s_bandwidth",
		Title:    "Bandwidth Disk %s",
		Units:    "KiB/s",
		Fam:      "disk",
		Ctx:      "wmi.logical_disk_bandwidth",
		Type:     module.Area,
		Priority: logicalDiskPriority + 1,
		Dims: Dims{
			{ID: "logical_disk_%s_read_bytes_total", Name: "read", Algo: module.Incremental, Div: 1000 * 1024},
			{ID: "logical_disk_%s_write_bytes_total", Name: "write", Algo: module.Incremental, Div: -1000 * 1024},
		},
	},
	{
		ID:       "logical_disk_%s_operations",
		Title:    "Operations Disk %s",
		Units:    "operations/s",
		Fam:      "disk",
		Ctx:      "wmi.logical_disk_operations",
		Type:     module.Area,
		Priority: logicalDiskPriority + 2,
		Dims: Dims{
			{ID: "logical_disk_%s_reads_total", Name: "reads", Algo: module.Incremental},
			{ID: "logical_disk_%s_writes_total", Name: "writes", Algo: module.Incremental, Mul: -1},
		},
	},
}

func (w *WMI) updateCharts(mx *metrics) {
	w.updateCollectDurationChart(mx)

	if mx.hasCPU() {
		w.updateCPUCharts(mx)
	}

	if mx.hasMem() {
		w.updateMemoryCharts(mx)
	}

	if mx.hasNet() {
		w.updateNetCharts(mx)
	}

	if mx.hasLDisk() {
		w.updateLogicalDisksCharts(mx)
	}

	//if mx.hasSystem() {
	//	w.updateSystemCharts(mx)
	//}

	//if mx.hasOS() {
	//	w.updateOSCharts(mx)
	//}
}

func (w *WMI) updateCollectDurationChart(mx *metrics) {
	for k := range mx.CollectDuration {
		chart := w.charts.Get("collector_duration")
		if !chart.HasDim(k) {
			_ = chart.AddDim(&Dim{ID: k})
			chart.MarkNotCreated()
		}
	}
}

func (w *WMI) updateCPUCharts(mx *metrics) {
	if !w.collected.collectors[collectorCPU] {
		w.collected.collectors[collectorCPU] = true
		_ = w.charts.Add(*cpuCharts.Copy()...)
	}

	for _, core := range mx.CPU.Cores {
		if w.collected.cores[core.ID] {
			continue
		}
		chart := cpuCoreUsageChart.Copy()
		chart.ID = fmt.Sprintf(chart.ID, core.ID)
		chart.Title = fmt.Sprintf(chart.Title, core.ID)
		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, core.ID)
		}
		_ = w.charts.Add(chart)
	}

	for _, core := range mx.CPU.Cores {
		if w.collected.cores[core.ID] {
			continue
		}
		chart := cpuCoreCStateChart.Copy()
		chart.ID = fmt.Sprintf(chart.ID, core.ID)
		chart.Title = fmt.Sprintf(chart.Title, core.ID)
		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, core.ID)
		}
		_ = w.charts.Add(chart)
	}

	for _, core := range mx.CPU.Cores {
		if w.collected.cores[core.ID] {
			continue
		}
		chart := w.charts.Get("cpu_dpcs")
		dim := &Dim{
			ID:   fmt.Sprintf("cpu_core_%s_dpc", core.ID),
			Name: "core" + core.ID,
			Algo: module.Incremental,
			Div:  1000,
		}
		_ = chart.AddDim(dim)
		chart.MarkNotCreated()

		chart = w.charts.Get("cpu_interrupts")
		dim = &Dim{
			ID:   fmt.Sprintf("cpu_core_%s_interrupts", core.ID),
			Name: "core" + core.ID,
			Algo: module.Incremental,
			Div:  1000,
		}
		_ = chart.AddDim(dim)
		chart.MarkNotCreated()

		w.collected.cores[core.ID] = true
	}
}

func (w *WMI) updateMemoryCharts(mx *metrics) {
	if w.collected.collectors[collectorMemory] {
		return
	}
	w.collected.collectors[collectorMemory] = true
	charts := *memoryCharts.Copy()
	for i, chart := range charts {
		chart.Priority += i + 1
	}
	_ = w.charts.Add(charts...)
}

func (w *WMI) updateNetCharts(mx *metrics) {
	for _, nic := range mx.Net.NICs {
		if w.collected.nics[nic.ID] {
			continue
		}
		w.collected.nics[nic.ID] = true
		charts := netNICCharts.Copy()

		for _, chart := range *charts {
			chart.ID = fmt.Sprintf(chart.ID, nic.ID)
			chart.Title = fmt.Sprintf(chart.Title, nic.ID)

			for _, dim := range chart.Dims {
				dim.ID = fmt.Sprintf(dim.ID, nic.ID)
			}
			for _, v := range chart.Vars {
				v.ID = fmt.Sprintf(v.ID, nic.ID)
			}

		}
		_ = w.charts.Add(*charts...)
	}
}

func (w *WMI) updateLogicalDisksCharts(mx *metrics) {
	set := make(map[string]bool)

	for _, vol := range mx.LogicalDisk.Volumes {
		set[vol.ID] = true

		if w.collected.volumes[vol.ID] {
			continue
		}

		w.collected.volumes[vol.ID] = true
		charts := logicalDiskCharts.Copy()

		for _, chart := range *charts {
			chart.ID = fmt.Sprintf(chart.ID, vol.ID)
			chart.Title = fmt.Sprintf(chart.Title, vol.ID)

			for _, dim := range chart.Dims {
				dim.ID = fmt.Sprintf(dim.ID, vol.ID)
			}
		}
		_ = w.charts.Add(*charts...)
	}

	for volID := range w.collected.volumes {
		if !set[volID] {
			w.removeLogicalDiskFromCharts(volID)
			delete(w.collected.volumes, volID)
		}
	}
}

func (w *WMI) removeLogicalDiskFromCharts(diskID string) {
	for _, chart := range *w.charts {
		if !strings.HasPrefix(chart.ID, fmt.Sprintf("logical_disk_%s", diskID)) {
			continue
		}
		chart.MarkRemove()
		chart.MarkNotCreated()
	}
}

func (w *WMI) updateSystemCharts(mx *metrics) {
	if w.collected.collectors[collectorSystem] {
		return
	}
	w.collected.collectors[collectorSystem] = true
	charts := *systemCharts.Copy()
	for i, chart := range charts {
		chart.Priority += i + 1
	}
	_ = w.charts.Add(charts...)
}

func (w *WMI) updateOSCharts(mx *metrics) {}

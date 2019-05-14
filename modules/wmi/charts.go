package wmi

import (
	"fmt"

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
	// Opts is an alias for module.DimOpts
	Opts = module.DimOpts
)

var cpuCharts = Charts{
	{
		ID:    "cpu_usage_total",
		Title: "CPU Usage Total",
		Units: "percentage",
		Fam:   "cpu",
		Ctx:   "cpu.cpu_usage_total",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "cpu_time_idle", Name: "idle", Algo: module.PercentOfIncremental, Div: 1000, DimOpts: Opts{Hidden: true}},
			{ID: "cpu_time_dpc", Name: "dpc", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_time_user", Name: "user", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_time_privileged", Name: "privileged", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_time_interrupt", Name: "interrupt", Algo: module.PercentOfIncremental, Div: 1000},
		},
	},
	{
		ID:    "cpu_dpcs_total",
		Title: "Received and Serviced Deferred Procedure Calls (DPC)",
		Units: "dpc/s",
		Fam:   "cpu",
		Ctx:   "cpu.cpu_dpcs_total",
		Type:  module.Stacked,
		// Dims will be added during collection
	},
	{
		ID:    "cpu_interrupts_total",
		Title: "Received and Serviced Hardware Interrupts",
		Units: "interrupts/s",
		Fam:   "cpu",
		Ctx:   "cpu.cpu_interrupts_total",
		Type:  module.Stacked,
		// Dims will be added during collection
	},
}

var cpuCoreCharts = Charts{
	{
		ID:    "core_%s_cpu_usage",
		Title: "Core%s Usage",
		Units: "percentage",
		Fam:   "cpu core usage",
		Ctx:   "cpu.cpu_usage_per_core",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "cpu_core_%s_time_idle", Name: "idle", Algo: module.PercentOfIncremental, Div: 1000, DimOpts: Opts{Hidden: true}},
			{ID: "cpu_core_%s_time_dpc", Name: "dpc", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_core_%s_time_user", Name: "user", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_core_%s_time_privileged", Name: "privileged", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_core_%s_time_interrupt", Name: "interrupt", Algo: module.PercentOfIncremental, Div: 1000},
		},
	},
	{
		ID:    "core_%s_cpu_cstate",
		Title: "Core%s Time Spent in Low-Power Idle State",
		Units: "percentage",
		Fam:   "cpu core c-state",
		Ctx:   "cpu.cpu_cstate_per_core",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "cpu_core_%s_cstate_c1", Name: "c1", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_core_%s_cstate_c2", Name: "c2", Algo: module.PercentOfIncremental, Div: 1000},
			{ID: "cpu_core_%s_cstate_c3", Name: "c3", Algo: module.PercentOfIncremental, Div: 1000},
		},
	},
}

func (w *WMI) updateCharts(mx *metrics) {
	w.updateCPUCharts(mx)
}

func (w *WMI) updateCPUCharts(mx *metrics) {
	enabled := mx.CPU != nil
	if !enabled {
		return
	}

	if !w.collected.collectors[collectorCPU] {
		w.collected.collectors[collectorCPU] = true
		_ = w.charts.Add(*cpuCharts.Copy()...)
	}

	for _, core := range mx.CPU.Cores {
		if w.collected.cpuCores[core.ID] {
			continue
		}

		w.collected.cpuCores[core.ID] = true

		charts := cpuCoreCharts.Copy()

		for _, chart := range *charts {
			chart.ID = fmt.Sprintf(chart.ID, core.ID)
			chart.Title = fmt.Sprintf(chart.Title, core.ID)
			for _, dim := range chart.Dims {
				dim.ID = fmt.Sprintf(dim.ID, core.ID)
			}
		}
		_ = w.charts.Add(*charts...)

		dim := &Dim{
			ID:   fmt.Sprintf("cpu_core_%s_dpc", core.ID),
			Name: "core" + core.ID,
			Algo: module.Incremental,
			Div:  1000,
		}
		_ = w.charts.Get("cpu_dpcs_total").AddDim(dim)

		dim = &Dim{
			ID:   fmt.Sprintf("cpu_core_%s_interrupts", core.ID),
			Name: "core" + core.ID,
			Algo: module.Incremental,
			Div:  1000,
		}
		_ = w.charts.Get("cpu_interrupts_total").AddDim(dim)
	}
}

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

var charts = Charts{
	{
		ID:    "cpu_utilization_total",
		Title: "CPU Utilization Total",
		Units: "percentage",
		Fam:   "cpu",
		Ctx:   "cpu.cpu_utilization_total",
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

func (w *WMI) updateCharts(mx *metrics) {
	for _, c := range mx.CPU.Cores {
		if !w.collectedCPUCores[c.ID] {
			w.collectedCPUCores[c.ID] = true
			dim := &Dim{
				ID:   fmt.Sprintf("cpu_core_%s_dpc", c.ID),
				Name: "core" + c.ID,
				Algo: module.Incremental,
				Div:  1000,
			}
			_ = w.charts.Get("cpu_dpcs_total").AddDim(dim)

			dim = &Dim{
				ID:   fmt.Sprintf("cpu_core_%s_interrupts", c.ID),
				Name: "core" + c.ID,
				Algo: module.Incremental,
				Div:  1000,
			}
			_ = w.charts.Get("cpu_interrupts_total").AddDim(dim)
		}
	}
}

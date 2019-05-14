package wmi

import "github.com/netdata/go-orchestrator/module"

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
		ID:    "cpu_usage_total",
		Title: "CPU Usage Total",
		Units: "percentage",
		Fam:   "cpu",
		Ctx:   "cpu.cpu_usage_total",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "cpu_time_idle", Name: "idle", Algo: module.PercentOfIncremental, DimOpts: Opts{Hidden: true}},
			{ID: "cpu_time_dpc", Name: "dpc", Algo: module.PercentOfIncremental},
			{ID: "cpu_time_user", Name: "user", Algo: module.PercentOfIncremental},
			{ID: "cpu_time_privileged", Name: "privileged", Algo: module.PercentOfIncremental},
			{ID: "cpu_time_interrupt", Name: "interrupt", Algo: module.PercentOfIncremental},
		},
	},
}

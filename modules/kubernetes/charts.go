package kubernetes

import "github.com/netdata/go-orchestrator/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Chart is an alias for module.Chart
	Chart = module.Chart
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

var (
	chartCPUStats = Chart{
		ID:    "%s_cpu_stats",
		Title: "CPU Stats",
		Units: "%",
		Ctx:   "kubernetes.cpu_stats",
		Dims: Dims{
			{
				ID:   "%s_cpu_stats_usage_core_nano_seconds",
				Name: "usage",
				Mul:  100,
				Div:  1000000000,
				Algo: module.Incremental,
			},
		},
	}
	chartMemoryStatsUsage = Chart{
		ID:    "%s_memory_stats_usage",
		Title: "Memory Usage",
		Units: "KB",
		Ctx:   "kubernetes.memory_stats_usage",
		Dims: Dims{
			{
				ID:   "%s_memory_stats_available_bytes",
				Name: "available",
				Div:  1024,
				Algo: module.Incremental,
			},
			{
				ID:   "%s_memory_stats_usage_bytes",
				Name: "usage",
				Div:  1024,
				Algo: module.Incremental,
			},
			{
				ID:   "%s_memory_stats_working_set_bytes",
				Name: "working set",
				Div:  1024,
				Algo: module.Incremental,
			},
			{
				ID:   "%s_memory_stats_working_rss_bytes",
				Name: "rss",
				Div:  1024,
				Algo: module.Incremental,
			},
		},
	}
	chartMemoryStatsPageFaults = Chart{
		ID:    "%s_memory_stats_page_faults",
		Title: "Page Faults",
		Units: "KB",
		Ctx:   "kubernetes.memory_stats_page_faults",
		Dims: Dims{
			{
				ID:   "%s_memory_stats_page_faults",
				Name: "minor",
				Algo: module.Incremental,
			},
			{
				ID:   "%s_memory_stats_major_page_faults",
				Name: "major",
				Algo: module.Incremental,
			},
		},
	}
)

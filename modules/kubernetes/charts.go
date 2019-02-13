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
		Title: "Cumulative All Cores CPU Stats",
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
		Type:  module.Stacked,
		Dims: Dims{
			{
				ID:   "%s_memory_stats_available_bytes",
				Name: "available",
				Div:  1024,
			},
			{
				ID:   "%s_memory_stats_usage_bytes",
				Name: "usage",
				Div:  1024,
			},
			{
				ID:   "%s_memory_stats_working_set_bytes",
				Name: "working set",
				Div:  1024,
			},
			{
				ID:   "%s_memory_stats_rss_bytes",
				Name: "rss",
				Div:  1024,
			},
		},
	}
	chartMemoryStatsPageFaults = Chart{
		ID:    "%s_memory_stats_page_faults",
		Title: "Page Faults",
		Units: "pages",
		Ctx:   "kubernetes.memory_stats_page_faults",
		Dims: Dims{
			{
				ID:   "%s_memory_stats_page_faults",
				Name: "minor",
			},
			{
				ID:   "%s_memory_stats_major_page_faults",
				Name: "major",
			},
		},
	}
	chartInterfaceBandwidth = Chart{
		ID:    "%s_interface_stats_bandwidth",
		Title: "Interface Bandwidth",
		Units: "kilobits/s",
		Ctx:   "kubernetes.interface_stats_bandwidth",
		Type:  module.Area,
		Dims: Dims{
			{
				ID:   "%s_interface_stats_rx_bytes",
				Name: "rx",
				Algo: module.Incremental,
				Mul:  8,
				Div:  1000,
			},
			{
				ID:   "%s_interface_stats_tx_bytes",
				Name: "tx",
				Algo: module.Incremental,
				Mul:  -8,
				Div:  1000,
			},
		},
	}
	chartInterfaceErrors = Chart{
		ID:    "%s_interface_stats_errors",
		Title: "Interface Errors",
		Units: "errors/s",
		Ctx:   "kubernetes.interface_stats_errors",
		Type:  module.Area,
		Dims: Dims{
			{
				ID:   "%s_interface_stats_rx_errors",
				Name: "rx",
				Algo: module.Incremental,
				Mul:  8,
				Div:  1000,
			},
			{
				ID:   "%s_interface_stats_tx_errors",
				Name: "tx",
				Algo: module.Incremental,
				Mul:  -8,
				Div:  1000,
			},
		},
	}
)

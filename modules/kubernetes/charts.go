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
				ID:   "%s_memory_stats",
				Name: "usage",
				Mul:  100,
				Div:  1000000000,
			},
		},
	}
)

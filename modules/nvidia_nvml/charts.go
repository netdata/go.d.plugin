package nvidia_nvml

//
//import (
//	"fmt"
//
//	"github.com/netdata/go-orchestrator/module"
//)
//
//type (
//	// Charts is an alias for module.Charts
//	Charts = module.Charts
//	// Dims is an alias for module.Dims
//	Dims = module.Dims
//)
//
//var charts = Charts{
//	{
//		ID:    "%s_utilization",
//		Title: "Utilization Rates",
//		Units: "percent",
//		Ctx:   "nvidia_nvml.utilization_rates",
//		Dims: Dims{
//			{ID: "%s_gpu_util", Name: "gpu"},
//			{ID: "%s_mem_util", Name: "memory"},
//			{ID: "%s_encoder_util", Name: "encoder"},
//			{ID: "%s_decoder_util", Name: "decoder"},
//		},
//	},
//	{
//		ID:    "%s_memory",
//		Title: "Memory Usage",
//		Units: "KiB",
//		Ctx:   "nvidia_nvml.memory_usage",
//		Type:  module.Stacked,
//		Dims: Dims{
//			{ID: "%s_memory_free", Name: "free", Div: 1024},
//			{ID: "%s_memory_used", Name: "used", Div: 1024},
//		},
//	},
//	{
//		ID:    "%s_temperature",
//		Title: "Temperature",
//		Units: "celsius",
//		Ctx:   "nvidia_nvml.temperature",
//		Dims: Dims{
//			{ID: "%s_temperature", Name: "temperature"},
//		},
//	},
//	{
//		ID:    "%s_fan_speed",
//		Title: "Fan Speed",
//		Units: "percent",
//		Ctx:   "nvidia_nvml.fan_speed",
//		Dims: Dims{
//			{ID: "%s_fan_speed", Name: "speed"},
//		},
//	},
//	{
//		ID:    "%s_power_usage",
//		Title: "Power Usage",
//		Units: "watts",
//		Ctx:   "nvidia_nvml.power_usage",
//		Dims: Dims{
//			{ID: "%s_power_usage", Name: "usage", Div: 1000},
//		},
//	},
//}
//
//func createGPUCharts(gpu gpu) *Charts {
//	charts := charts.Copy()
//
//	for _, chart := range *charts {
//		chart.ID = fmt.Sprintf(chart.ID, gpu.uniqName())
//		chart.Fam = gpu.name
//		for _, dim := range chart.Dims {
//			dim.ID = fmt.Sprintf(dim.ID, gpu.uniqName())
//		}
//	}
//
//	return charts
//}

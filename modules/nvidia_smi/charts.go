// SPDX-License-Identifier: GPL-3.0-or-later

package nvidia_smi

import (
	"fmt"
	"strings"

	"github.com/netdata/go.d.plugin/agent/module"
)

const (
	prioGPUPCIBandwidthUsage = module.Priority + iota
	prioGPUPCIBandwidthUtilization
	prioGPUFanSpeed
	prioGPUUtilization
	prioGPUMemUtilization
	prioGPUDecoderUtilization
	prioGPUEncoderUtilization
	prioGPUFBMemoryUsage
	prioGPUBAR1MemoryUsage
	prioGPUTemperatureChart
	prioGPUVoltageChart
	prioGPUClockFreq
	prioGPUPowerDraw
	prioGPUPerformanceState
)

var (
	gpuXMLCharts = module.Charts{
		gpuPCIBandwidthUsageChartTmpl.Copy(),
		gpuPCIBandwidthUtilizationChartTmpl.Copy(),
		gpuFanSpeedPercChartTmpl.Copy(),
		gpuUtilizationChartTmpl.Copy(),
		gpuMemUtilizationChartTmpl.Copy(),
		gpuDecoderUtilizationChartTmpl.Copy(),
		gpuEncoderUtilizationChartTmpl.Copy(),
		gpuFrameBufferMemoryUsageChartTmpl.Copy(),
		gpuBAR1MemoryUsageChartTmpl.Copy(),
		gpuVoltageChartTmpl.Copy(),
		gpuTemperatureChartTmpl.Copy(),
		gpuClockFreqChartTmpl.Copy(),
		gpuPowerDrawChartTmpl.Copy(),
		gpuPerformanceStateChartTmpl.Copy(),
	}
	gpuCSVCharts = module.Charts{
		gpuFanSpeedPercChartTmpl.Copy(),
		gpuUtilizationChartTmpl.Copy(),
		gpuMemUtilizationChartTmpl.Copy(),
		gpuFrameBufferMemoryUsageChartTmpl.Copy(),
		gpuTemperatureChartTmpl.Copy(),
		gpuClockFreqChartTmpl.Copy(),
		gpuPowerDrawChartTmpl.Copy(),
		gpuPerformanceStateChartTmpl.Copy(),
	}
)

var (
	gpuPCIBandwidthUsageChartTmpl = module.Chart{
		ID:       "gpu_%s_pcie_bandwidth_usage",
		Title:    "PCI Express Bandwidth Usage",
		Units:    "B/s",
		Fam:      "pcie bandwidth",
		Ctx:      "nvidia_smi.gpu_pcie_bandwidth_usage",
		Type:     module.Area,
		Priority: prioGPUPCIBandwidthUsage,
		Dims: module.Dims{
			{ID: "gpu_%s_pcie_bandwidth_usage_rx", Name: "rx"},
			{ID: "gpu_%s_pcie_bandwidth_usage_tx", Name: "tx", Mul: -1},
		},
	}
	gpuPCIBandwidthUtilizationChartTmpl = module.Chart{
		ID:       "gpu_%s_pcie_bandwidth_utilization",
		Title:    "PCI Express Bandwidth Utilization",
		Units:    "percentage",
		Fam:      "pcie bandwidth",
		Ctx:      "nvidia_smi.gpu_pcie_bandwidth_utilization",
		Priority: prioGPUPCIBandwidthUtilization,
		Dims: module.Dims{
			{ID: "gpu_%s_pcie_bandwidth_utilization_rx", Name: "rx", Div: 100},
			{ID: "gpu_%s_pcie_bandwidth_utilization_tx", Name: "tx", Div: 100},
		},
	}
	gpuFanSpeedPercChartTmpl = module.Chart{
		ID:       "gpu_%s_fan_speed_perc",
		Title:    "Fan speed",
		Units:    "%",
		Fam:      "fan speed",
		Ctx:      "nvidia_smi.gpu_fan_speed_perc",
		Priority: prioGPUFanSpeed,
		Dims: module.Dims{
			{ID: "gpu_%s_fan_speed_perc", Name: "fan_speed"},
		},
	}
	gpuUtilizationChartTmpl = module.Chart{
		ID:       "gpu_%s_gpu_utilization",
		Title:    "GPU utilization",
		Units:    "%",
		Fam:      "gpu utilization",
		Ctx:      "nvidia_smi.gpu_utilization",
		Priority: prioGPUUtilization,
		Dims: module.Dims{
			{ID: "gpu_%s_gpu_utilization", Name: "gpu"},
		},
	}
	gpuMemUtilizationChartTmpl = module.Chart{
		ID:       "gpu_%s_memory_utilization",
		Title:    "Memory utilization",
		Units:    "%",
		Fam:      "mem utilization",
		Ctx:      "nvidia_smi.gpu_memory_utilization",
		Priority: prioGPUMemUtilization,
		Dims: module.Dims{
			{ID: "gpu_%s_mem_utilization", Name: "memory"},
		},
	}
	gpuDecoderUtilizationChartTmpl = module.Chart{
		ID:       "gpu_%s_decoder_utilization",
		Title:    "Decoder utilization",
		Units:    "%",
		Fam:      "dec utilization",
		Ctx:      "nvidia_smi.gpu_decoder_utilization",
		Priority: prioGPUDecoderUtilization,
		Dims: module.Dims{
			{ID: "gpu_%s_decoder_utilization", Name: "decoder"},
		},
	}
	gpuEncoderUtilizationChartTmpl = module.Chart{
		ID:       "gpu_%s_encoder_utilization",
		Title:    "Encoder utilization",
		Units:    "%",
		Fam:      "enc utilization",
		Ctx:      "nvidia_smi.gpu_encoder_utilization",
		Priority: prioGPUEncoderUtilization,
		Dims: module.Dims{
			{ID: "gpu_%s_encoder_utilization", Name: "encoder"},
		},
	}
	gpuFrameBufferMemoryUsageChartTmpl = module.Chart{
		ID:       "gpu_%s_frame_buffer_memory_usage",
		Title:    "Frame buffer memory usage",
		Units:    "B",
		Fam:      "fb mem usage",
		Ctx:      "nvidia_smi.gpu_frame_buffer_memory_usage",
		Type:     module.Stacked,
		Priority: prioGPUFBMemoryUsage,
		Dims: module.Dims{
			{ID: "gpu_%s_frame_buffer_memory_usage_free", Name: "free"},
			{ID: "gpu_%s_frame_buffer_memory_usage_used", Name: "used"},
			{ID: "gpu_%s_frame_buffer_memory_usage_reserved", Name: "reserved"},
		},
	}
	gpuBAR1MemoryUsageChartTmpl = module.Chart{
		ID:       "gpu_%s_bar1_memory_usage",
		Title:    "BAR1 memory usage",
		Units:    "B",
		Fam:      "bar1 mem usage",
		Ctx:      "nvidia_smi.gpu_bar1_memory_usage",
		Type:     module.Stacked,
		Priority: prioGPUBAR1MemoryUsage,
		Dims: module.Dims{
			{ID: "gpu_%s_bar1_memory_usage_free", Name: "free"},
			{ID: "gpu_%s_bar1_memory_usage_used", Name: "used"},
		},
	}
	gpuTemperatureChartTmpl = module.Chart{
		ID:       "gpu_%s_temperature",
		Title:    "Temperature",
		Units:    "Celsius",
		Fam:      "temperature",
		Ctx:      "nvidia_smi.gpu_temperature",
		Priority: prioGPUTemperatureChart,
		Dims: module.Dims{
			{ID: "gpu_%s_temperature", Name: "temperature"},
		},
	}
	gpuVoltageChartTmpl = module.Chart{
		ID:       "gpu_%s_voltage",
		Title:    "Voltage",
		Units:    "V",
		Fam:      "voltage",
		Ctx:      "nvidia_smi.gpu_voltage",
		Priority: prioGPUVoltageChart,
		Dims: module.Dims{
			{ID: "gpu_%s_voltage", Name: "voltage", Div: 1000}, // mV => V
		},
	}
	gpuClockFreqChartTmpl = module.Chart{
		ID:       "gpu_%s_clock_freq",
		Title:    "Clock current frequency",
		Units:    "MHz",
		Fam:      "clocks",
		Ctx:      "nvidia_smi.gpu_clock_freq",
		Priority: prioGPUClockFreq,
		Dims: module.Dims{
			{ID: "gpu_%s_graphics_clock", Name: "graphics"},
			{ID: "gpu_%s_video_clock", Name: "video"},
			{ID: "gpu_%s_sm_clock", Name: "sm"},
			{ID: "gpu_%s_mem_clock", Name: "mem"},
		},
	}
	gpuPowerDrawChartTmpl = module.Chart{
		ID:       "gpu_%s_power_draw",
		Title:    "Power draw",
		Units:    "Watts",
		Fam:      "power draw",
		Ctx:      "nvidia_smi.gpu_power_draw",
		Priority: prioGPUPowerDraw,
		Dims: module.Dims{
			{ID: "gpu_%s_power_draw", Name: "power_draw"},
		},
	}
	gpuPerformanceStateChartTmpl = module.Chart{
		ID:       "gpu_%s_performance_state",
		Title:    "Performance state",
		Units:    "state",
		Fam:      "performance state",
		Ctx:      "nvidia_smi.gpu_performance_state",
		Priority: prioGPUPerformanceState,
		Dims: module.Dims{
			{ID: "gpu_%s_performance_state_P0", Name: "P0"},
			{ID: "gpu_%s_performance_state_P1", Name: "P1"},
			{ID: "gpu_%s_performance_state_P2", Name: "P2"},
			{ID: "gpu_%s_performance_state_P3", Name: "P3"},
			{ID: "gpu_%s_performance_state_P4", Name: "P4"},
			{ID: "gpu_%s_performance_state_P5", Name: "P5"},
			{ID: "gpu_%s_performance_state_P6", Name: "P6"},
			{ID: "gpu_%s_performance_state_P7", Name: "P7"},
			{ID: "gpu_%s_performance_state_P8", Name: "P8"},
			{ID: "gpu_%s_performance_state_P9", Name: "P9"},
			{ID: "gpu_%s_performance_state_P10", Name: "P10"},
			{ID: "gpu_%s_performance_state_P11", Name: "P11"},
			{ID: "gpu_%s_performance_state_P12", Name: "P12"},
			{ID: "gpu_%s_performance_state_P13", Name: "P13"},
			{ID: "gpu_%s_performance_state_P14", Name: "P14"},
			{ID: "gpu_%s_performance_state_P15", Name: "P15"},
		},
	}
)

func (nv *NvidiaSMI) addGPUXMLCharts(gpu xmlGPUInfo) {
	charts := gpuXMLCharts.Copy()
	if !isValidValue(gpu.FanSpeed) {
		_ = charts.Remove(gpuFanSpeedPercChartTmpl.ID)
	}
	if !isValidValue(gpu.PowerReadings.PowerDraw) {
		_ = charts.Remove(gpuPowerDrawChartTmpl.ID)
	}
	if !isValidValue(gpu.Voltage.GraphicsVolt) {
		_ = charts.Remove(gpuVoltageChartTmpl.ID)
	}

	for _, c := range *charts {
		c.ID = fmt.Sprintf(c.ID, strings.ToLower(gpu.UUID))
		c.Labels = []module.Label{
			// csv output has no 'product_brand'
			{Key: "product_name", Value: gpu.ProductName},
		}
		for _, d := range c.Dims {
			d.ID = fmt.Sprintf(d.ID, gpu.UUID)
		}
	}

	if err := nv.Charts().Add(*charts...); err != nil {
		nv.Warning(err)
	}
}

func (nv *NvidiaSMI) addGPUCSVCharts(gpu csvGPUInfo) {
	charts := gpuCSVCharts.Copy()

	if !isValidValue(gpu.fanSpeed) {
		_ = charts.Remove(gpuFanSpeedPercChartTmpl.ID)
	}
	if !isValidValue(gpu.powerDraw) {
		_ = charts.Remove(gpuPowerDrawChartTmpl.ID)
	}

	for _, c := range *charts {
		c.ID = fmt.Sprintf(c.ID, strings.ToLower(gpu.uuid))
		c.Labels = []module.Label{
			{Key: "product_name", Value: gpu.name},
		}
		for _, d := range c.Dims {
			d.ID = fmt.Sprintf(d.ID, gpu.uuid)
		}
	}

	if err := nv.Charts().Add(*charts...); err != nil {
		nv.Warning(err)
	}
}

func (nv *NvidiaSMI) removeGPUCharts(uuid string) {
	prefix := "gpu_" + strings.ToLower(uuid)

	for _, c := range *nv.Charts() {
		if strings.HasPrefix(c.ID, prefix) {
			c.MarkRemove()
			c.MarkNotCreated()
		}
	}
}

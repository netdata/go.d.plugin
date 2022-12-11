// SPDX-License-Identifier: GPL-3.0-or-later

package nvidia_smi

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

func (nv *NvidiaSMI) collectGPUInfoXML(mx map[string]int64) error {
	bs, err := nv.exec.queryGPUInfoXML()
	if err != nil {
		return fmt.Errorf("error on quering XML GPU info: %v", err)
	}

	info := &xmlInfo{}
	if err := xml.Unmarshal(bs, info); err != nil {
		return fmt.Errorf("error on unmarshaling XML GPU info response: %v", err)
	}

	seen := make(map[string]bool)

	for _, gpu := range info.GPUs {
		if !isValidValue(gpu.UUID) {
			continue
		}

		seen[gpu.UUID] = true

		if !nv.gpus[gpu.UUID] {
			nv.gpus[gpu.UUID] = true
			nv.addGPUXMLCharts(gpu)
		}

		px := "gpu_" + gpu.UUID + "_"

		addMetric(mx, px+"pcie_bandwidth_usage_rx", gpu.PCI.RxUtil, 1024) // KB => bytes
		addMetric(mx, px+"pcie_bandwidth_usage_tx", gpu.PCI.TxUtil, 1024) // KB => bytes
		addMetric(mx, px+"fan_speed_perc", gpu.FanSpeed, 0)
		addMetric(mx, px+"gpu_utilization", gpu.Utilization.GpuUtil, 0)
		addMetric(mx, px+"mem_utilization", gpu.Utilization.MemoryUtil, 0)
		addMetric(mx, px+"decoder_utilization", gpu.Utilization.DecoderUtil, 0)
		addMetric(mx, px+"encoder_utilization", gpu.Utilization.EncoderUtil, 0)
		addMetric(mx, px+"frame_buffer_memory_usage_free", gpu.FBMemoryUsage.Free, 1024*1024)         // MiB => bytes
		addMetric(mx, px+"frame_buffer_memory_usage_used", gpu.FBMemoryUsage.Used, 1024*1024)         // MiB => bytes
		addMetric(mx, px+"frame_buffer_memory_usage_reserved", gpu.FBMemoryUsage.Reserved, 1024*1024) // MiB => bytes
		addMetric(mx, px+"bar1_memory_usage_free", gpu.Bar1MemoryUsage.Free, 1024*1024)               // MiB => bytes
		addMetric(mx, px+"bar1_memory_usage_used", gpu.Bar1MemoryUsage.Used, 1024*1024)               // MiB => bytes
		addMetric(mx, px+"temperature", gpu.Temperature.GpuTemp, 0)
		addMetric(mx, px+"graphics_clock", gpu.Clocks.GraphicsClock, 0)
		addMetric(mx, px+"video_clock", gpu.Clocks.VideoClock, 0)
		addMetric(mx, px+"sm_clock", gpu.Clocks.SmClock, 0)
		addMetric(mx, px+"mem_clock", gpu.Clocks.MemClock, 0)
		addMetric(mx, px+"power_draw", gpu.PowerReadings.PowerDraw, 0)
		for i := 0; i < 16; i++ {
			if s := "P" + strconv.Itoa(i); gpu.PerformanceState == s {
				mx[px+"performance_state_"+s] = 1
			} else {
				mx[px+"performance_state_"+s] = 0
			}
		}
	}

	for uuid := range nv.gpus {
		if !seen[uuid] {
			delete(nv.gpus, uuid)
			nv.removeGPUCharts(uuid)
		}
	}

	return nil
}

type (
	xmlInfo struct {
		GPUs []xmlGPUInfo `xml:"gpu"`
	}
	xmlGPUInfo struct {
		ID                  string `xml:"id,attr"`
		ProductName         string `xml:"product_name"`
		ProductBrand        string `xml:"product_brand"`
		ProductArchitecture string `xml:"product_architecture"`
		UUID                string `xml:"uuid"`
		FanSpeed            string `xml:"fan_speed"`
		PerformanceState    string `xml:"performance_state"`
		PCI                 struct {
			TxUtil string `xml:"tx_util"`
			RxUtil string `xml:"rx_util"`
		} `xml:"pci"`
		Utilization struct {
			GpuUtil     string `xml:"gpu_util"`
			MemoryUtil  string `xml:"memory_util"`
			EncoderUtil string `xml:"encoder_util"`
			DecoderUtil string `xml:"decoder_util"`
		} `xml:"utilization"`
		FBMemoryUsage struct {
			Total    string `xml:"total"`
			Reserved string `xml:"reserved"`
			Used     string `xml:"used"`
			Free     string `xml:"free"`
		} `xml:"fb_memory_usage"`
		Bar1MemoryUsage struct {
			Total string `xml:"total"`
			Used  string `xml:"used"`
			Free  string `xml:"free"`
		} `xml:"bar1_memory_usage"`
		Temperature struct {
			GpuTemp                string `xml:"gpu_temp"`
			GpuTempMaxThreshold    string `xml:"gpu_temp_max_threshold"`
			GpuTempSlowThreshold   string `xml:"gpu_temp_slow_threshold"`
			GpuTempMaxGpuThreshold string `xml:"gpu_temp_max_gpu_threshold"`
			GpuTargetTemperature   string `xml:"gpu_target_temperature"`
			MemoryTemp             string `xml:"memory_temp"`
			GpuTempMaxMemThreshold string `xml:"gpu_temp_max_mem_threshold"`
		} `xml:"temperature"`
		Clocks struct {
			GraphicsClock string `xml:"graphics_clock"`
			SmClock       string `xml:"sm_clock"`
			MemClock      string `xml:"mem_clock"`
			VideoClock    string `xml:"video_clock"`
		} `xml:"clocks"`
		PowerReadings struct {
			PowerState         string `xml:"power_state"`
			PowerManagement    string `xml:"power_management"`
			PowerDraw          string `xml:"power_draw"`
			PowerLimit         string `xml:"power_limit"`
			DefaultPowerLimit  string `xml:"default_power_limit"`
			EnforcedPowerLimit string `xml:"enforced_power_limit"`
			MinPowerLimit      string `xml:"min_power_limit"`
			MaxPowerLimit      string `xml:"max_power_limit"`
		} `xml:"power_readings"`
		Processes struct {
			ProcessInfo []struct {
				PID         string `xml:"pid"`
				ProcessName string `xml:"process_name"`
				UsedMemory  string `xml:"used_memory"`
			} `sml:"process_info"`
		} `xml:"processes"`
	}
)

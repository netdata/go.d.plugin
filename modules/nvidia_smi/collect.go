package nvidia_smi

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func (nv *NvidiaSMI) collect() (map[string]int64, error) {
	if nv.exec == nil {
		return nil, errors.New("")
	}

	mx := make(map[string]int64)

	bs, err := nv.exec.queryXML()
	if err != nil {
		return nil, fmt.Errorf("error on quering GPU info: %v", err)
	}

	info := &xmlInfo{}
	if err := xml.Unmarshal(bs, info); err != nil {
		return nil, fmt.Errorf("error on unmarshaling GPU info response: %v", err)
	}

	nv.collectXMLInfo(mx, info)

	return mx, nil
}

func (nv *NvidiaSMI) collectXMLInfo(mx map[string]int64, info *xmlInfo) {
	seen := make(map[string]bool)

	for _, gpu := range info.GPUs {
		if !isValidValue(gpu.UUID) {
			continue
		}

		seen[gpu.UUID] = true

		if !nv.gpus[gpu.UUID] {
			nv.gpus[gpu.UUID] = true
			nv.addGPUCharts(gpu)
		}

		px := "gpu_" + gpu.UUID + "_"

		addMetric(mx, px+"pcie_bandwidth_usage_rx", gpu.PCI.RxUtil, 1024)
		addMetric(mx, px+"pcie_bandwidth_usage_tx", gpu.PCI.TxUtil, 1024)
		addMetric(mx, px+"fan_speed_perc", gpu.FanSpeed, 0)
		addMetric(mx, px+"gpu_utilization", gpu.Utilization.GpuUtil, 0)
		addMetric(mx, px+"mem_utilization", gpu.Utilization.MemoryUtil, 0)
		addMetric(mx, px+"decoder_utilization", gpu.Utilization.DecoderUtil, 0)
		addMetric(mx, px+"encoder_utilization", gpu.Utilization.EncoderUtil, 0)
		addMetric(mx, px+"frame_buffer_memory_usage_free", gpu.FBMemoryUsage.Free, 1024*1024)
		addMetric(mx, px+"frame_buffer_memory_usage_used", gpu.FBMemoryUsage.Used, 1024*1024)
		addMetric(mx, px+"frame_buffer_memory_usage_reserved", gpu.FBMemoryUsage.Reserved, 1024*1024)
		addMetric(mx, px+"bar1_memory_usage_free", gpu.Bar1MemoryUsage.Free, 1024*1024)
		addMetric(mx, px+"bar1_memory_usage_used", gpu.Bar1MemoryUsage.Used, 1024*1024)
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
}

func addMetric(mx map[string]int64, key, value string, mul int) {
	if !isValidValue(value) {
		return
	}

	var i int
	if i = strings.IndexByte(value, ' '); i == -1 {
		return
	}
	value = value[:i]

	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return
	}

	if mul > 0 {
		v *= float64(mul)
	}

	mx[key] = int64(v)
}

func isValidValue(v string) bool {
	return v != "" && v != "N/A"
}

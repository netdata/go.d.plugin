package nvidia_smi

import (
	"context"
	"os/exec"
	"time"
)

func newNvidiaSMIExec(path string, cfg Config) (*nvidiaSMIExec, error) {
	return &nvidiaSMIExec{
		binPath: path,
		timeout: cfg.Timeout.Duration,
	}, nil
}

type nvidiaSMIExec struct {
	binPath string
	timeout time.Duration
}

func (e *nvidiaSMIExec) queryXML() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return exec.CommandContext(ctx, e.binPath, "-x", "-q").Output()
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

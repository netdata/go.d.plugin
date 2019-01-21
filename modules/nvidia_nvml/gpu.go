package nvidia_nvml

import (
	"fmt"
	"github.com/mindprince/gonvml"
)

type gpu struct {
	uuid     string
	name     string
	minorNum uint
	stats
}

func (g gpu) uniqName() string {
	return fmt.Sprintf("gpu%d", g.minorNum)
}

type stat struct {
	value int64
	exist bool
}

func (s *stat) set(v int64) {
	s.value = v
	s.exist = true
}

type stats struct {
	temperature stat
	powerUsage  stat
	fanSpeed    stat
	memoryTotal stat
	memoryUsed  stat
	gpuUtil     stat
	memUtil     stat
	encoderUtil stat
	decoderUtil stat
}

func (s stats) asMap() map[string]int64 {
	m := make(map[string]int64)

	if s.temperature.exist {
		m["temperature"] = s.temperature.value
	}
	if s.powerUsage.exist {
		m["power_usage"] = s.powerUsage.value
	}
	if s.fanSpeed.exist {
		m["fan_speed"] = s.fanSpeed.value
	}
	if s.memoryUsed.exist {
		m["memory_used"] = s.memoryUsed.value
	}
	if s.memoryTotal.exist && s.memoryUsed.exist {
		m["memory_free"] = s.memoryTotal.value - s.memoryUsed.value
	}
	if s.gpuUtil.exist {
		m["gpu_util"] = s.gpuUtil.value
	}
	if s.memUtil.exist {
		m["mem_util"] = s.memUtil.value
	}
	if s.encoderUtil.exist {
		m["encoder_util"] = s.encoderUtil.value
	}
	if s.decoderUtil.exist {
		m["decoder_util"] = s.decoderUtil.value
	}

	return m
}

func getGPUs() ([]gpu, error) {
	count, err := gonvml.DeviceCount()

	if err != nil {
		return nil, err
	}

	var gpus []gpu

	for i := 0; i < int(count); i++ {
		gpu, err := getGPUByIndex(i)
		if err != nil {
			return nil, err
		}
		gpus = append(gpus, gpu)
	}

	return gpus, nil

}

func getGPUByIndex(idx int) (gpu, error) {
	var g gpu

	dev, err := gonvml.DeviceHandleByIndex(uint(idx))

	if err != nil {
		return g, err
	}

	if g.name, err = dev.Name(); err != nil {
		return g, err
	}

	if g.uuid, err = dev.UUID(); err != nil {
		return g, err
	}

	if g.minorNum, err = dev.MinorNumber(); err != nil {
		return g, err
	}

	if temp, err := dev.Temperature(); err == nil {
		g.temperature.set(int64(temp))
	}

	if pow, err := dev.PowerUsage(); err == nil {
		g.powerUsage.set(int64(pow))
	}

	if fan, err := dev.FanSpeed(); err == nil {
		g.fanSpeed.set(int64(fan))
	}

	if gpuUtil, memUtil, err := dev.UtilizationRates(); err == nil {
		g.gpuUtil.set(int64(gpuUtil))
		g.memUtil.set(int64(memUtil))
	}

	if memTotal, memUsed, err := dev.MemoryInfo(); err == nil {
		g.memoryTotal.set(int64(memTotal))
		g.memoryUsed.set(int64(memUsed))
	}

	if encUtil, _, err := dev.EncoderUtilization(); err == nil {
		g.encoderUtil.set(int64(encUtil))
	}

	if decUtil, _, err := dev.EncoderUtilization(); err == nil {
		g.decoderUtil.set(int64(decUtil))
	}

	return g, nil
}

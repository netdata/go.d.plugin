package phpfpm

import (
	"math"

	"github.com/netdata/go.d.plugin/pkg/stm"
)

func (p *Phpfpm) collect() (map[string]int64, error) {
	status, err := p.client.Status()
	if err != nil {
		return nil, err
	}

	data := stm.ToMap(status)
	if len(status.Processes) == 0 {
		return data, nil
	}

	statProcesses(data, status.Processes, "ReqDur", func(p proc) int64 { return p.Duration })
	statProcesses(data, status.Processes, "ReqCpu", func(p proc) int64 { return int64(p.CPU) })
	statProcesses(data, status.Processes, "ReqMem", func(p proc) int64 { return p.Memory })

	return data, nil
}

type accessor func(p proc) int64

func statProcesses(m map[string]int64, procs []proc, met string, acc accessor) {
	var sum, count, min, max int64
	for _, proc := range procs {
		if proc.State != "Idle" {
			continue
		}

		val := acc(proc)
		sum += val
		count += 1
		if count == 1 {
			min, max = val, val
			continue
		}
		min = int64(math.Min(float64(min), float64(val)))
		max = int64(math.Max(float64(max), float64(val)))
	}

	m["min"+met] = min
	m["max"+met] = max
	m["avg"+met] = sum / count
}

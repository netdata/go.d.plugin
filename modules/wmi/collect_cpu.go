package wmi

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	metricCPUCstate     = "wmi_cpu_cstate_seconds_total"
	metricCPUDPCs       = "wmi_cpu_dpcs_total"
	metricCPUInterrupts = "wmi_cpu_interrupts_total"
	metricCPUTime       = "wmi_cpu_time_total"
)

func (w *WMI) collectCPU(mx *metrics, pms prometheus.Metrics) {
	collectCPUCoresCStates(mx, pms)
	collectCPUCoresDPCs(mx, pms)
	collectCPUCoresInterrupts(mx, pms)
	collectCPUCoresUsage(mx, pms)

	mx.CPU.Cores.sort()
	collectCPUSummary(mx)
}

func collectCPUSummary(mx *metrics) {
	for _, c := range mx.CPU.Cores {
		mx.CPU.PercentUserTime.Add(c.PercentUserTime.Value())
		mx.CPU.PercentPrivilegedTime.Add(c.PercentPrivilegedTime.Value())
		mx.CPU.PercentInterruptTime.Add(c.PercentInterruptTime.Value())
		mx.CPU.PercentIdleTime.Add(c.PercentIdleTime.Value())
		mx.CPU.PercentDPCTime.Add(c.PercentDPCTime.Value())
	}
}

func collectCPUCoresCStates(mx *metrics, pms prometheus.Metrics) {
	core := newCPUCore("")

	for _, pm := range pms.FindByName(metricCPUCstate) {
		var (
			coreID = pm.Labels.Get("core")
			state  = pm.Labels.Get("state")
			value  = pm.Value
		)
		if coreID == "" || state == "" {
			continue
		}
		if core.ID != coreID {
			core = mx.CPU.Cores.get(coreID, true)
		}
		switch state {
		default:
		case "c1":
			core.PercentC1Time.Set(value)
		case "c2":
			core.PercentC2Time.Set(value)
		case "c3":
			core.PercentC3Time.Set(value)
		}
	}
}

func collectCPUCoresInterrupts(mx *metrics, pms prometheus.Metrics) {
	core := newCPUCore("")

	for _, pm := range pms.FindByName(metricCPUInterrupts) {
		var (
			coreID = pm.Labels.Get("core")
			value  = pm.Value
		)
		if coreID == "" {
			continue
		}
		if core.ID != coreID {
			core = mx.CPU.Cores.get(coreID, true)
		}
		core.InterruptsPerSec.Set(value)
	}
}

func collectCPUCoresUsage(mx *metrics, pms prometheus.Metrics) {
	core := newCPUCore("")

	for _, pm := range pms.FindByName(metricCPUTime) {
		var (
			coreID = pm.Labels.Get("core")
			mode   = pm.Labels.Get("mode")
			value  = pm.Value
		)
		if coreID == "" || mode == "" {
			continue
		}
		if core.ID != coreID {
			core = mx.CPU.Cores.get(coreID, true)
		}
		switch mode {
		default:
		case "dpc":
			core.PercentDPCTime.Set(value)
		case "idle":
			core.PercentIdleTime.Set(value)
		case "interrupt":
			core.PercentInterruptTime.Set(value)
		case "privileged":
			core.PercentPrivilegedTime.Set(value)
		case "user":
			core.PercentUserTime.Set(value)
		}
	}
}

func collectCPUCoresDPCs(mx *metrics, pms prometheus.Metrics) {
	core := newCPUCore("")

	for _, pm := range pms.FindByName(metricCPUDPCs) {
		var (
			coreID = pm.Labels.Get("core")
			value  = pm.Value
		)
		if coreID == "" {
			continue
		}
		if core.ID != coreID {
			core = mx.CPU.Cores.get(coreID, true)
		}
		core.DPCsQueuedPerSec.Set(value)
	}
}

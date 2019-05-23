package wmi

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	collectorCPU = "cpu"

	metricCPUCstateTotal     = "wmi_cpu_cstate_seconds_total"
	metricCPUDPCsTotal       = "wmi_cpu_dpcs_total"
	metricCPUInterruptsTotal = "wmi_cpu_interrupts_total"
	metricCPUTimeTotal       = "wmi_cpu_time_total"
)

func collectCPU(mx *metrics, pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, collectorCPU)
	if !(enabled && success) {
		return false
	}
	mx.CPU = &cpu{}

	collectCPUCoresCStates(mx, pms)
	collectCPUCoresDPCs(mx, pms)
	collectCPUCoresInterrupts(mx, pms)
	collectCPUCoresUsage(mx, pms)

	mx.CPU.Cores.sort()
	collectCPUSummary(mx)

	return true
}

func collectCPUSummary(mx *metrics) {
	for _, c := range mx.CPU.Cores {
		mx.CPU.User += c.User
		mx.CPU.Privileged += c.Privileged
		mx.CPU.Interrupt += c.Interrupt
		mx.CPU.Idle += c.Idle
		mx.CPU.DPC += c.DPC
	}
}

func collectCPUCoresCStates(mx *metrics, pms prometheus.Metrics) {
	core := newCPUCore("")

	for _, pm := range pms.FindByName(metricCPUCstateTotal) {
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
			core.C1 = value
		case "c2":
			core.C2 = value
		case "c3":
			core.C3 = value
		}
	}
}

func collectCPUCoresInterrupts(mx *metrics, pms prometheus.Metrics) {
	core := newCPUCore("")

	for _, pm := range pms.FindByName(metricCPUInterruptsTotal) {
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
		core.InterruptsTotal = value
	}
}

func collectCPUCoresUsage(mx *metrics, pms prometheus.Metrics) {
	core := newCPUCore("")

	for _, pm := range pms.FindByName(metricCPUTimeTotal) {
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
			core.DPC = value
		case "idle":
			core.Idle = value
		case "interrupt":
			core.Interrupt = value
		case "privileged":
			core.Privileged = value
		case "user":
			core.User = value
		}
	}
}

func collectCPUCoresDPCs(mx *metrics, pms prometheus.Metrics) {
	core := newCPUCore("")

	for _, pm := range pms.FindByName(metricCPUDPCsTotal) {
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
		core.DPCsTotal = value
	}
}

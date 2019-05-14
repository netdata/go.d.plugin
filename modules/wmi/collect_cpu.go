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

	for _, c := range mx.CPU.Cores {
		mx.CPU.Time.User.Add(c.Time.User.Value())
		mx.CPU.Time.Privileged.Add(c.Time.Privileged.Value())
		mx.CPU.Time.Interrupt.Add(c.Time.Interrupt.Value())
		mx.CPU.Time.Idle.Add(c.Time.Idle.Value())
		mx.CPU.Time.DPC.Add(c.Time.DPC.Value())
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
			core.CState.C1.Set(value)
		case "c2":
			core.CState.C2.Set(value)
		case "c3":
			core.CState.C3.Set(value)
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
		core.Interrupts.Set(value)
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
			core.Time.DPC.Set(value)
		case "idle":
			core.Time.Idle.Set(value)
		case "interrupt":
			core.Time.Interrupt.Set(value)
		case "privileged":
			core.Time.Privileged.Set(value)
		case "user":
			core.Time.User.Set(value)
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
		core.DPCs.Set(value)
	}
}

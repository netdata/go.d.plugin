package wmi

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/stm"
)

const (
	metricCPUCstate     = "wmi_cpu_cstate_seconds_total"
	metricCPUDPCs       = "wmi_cpu_dpcs_total"
	metricCPUInterrupts = "wmi_cpu_interrupts_total"
	metricCPUTime       = "wmi_cpu_time_total"
)

func (w *WMI) collect() (map[string]int64, error) {
	scraped, err := w.prom.Scrape()
	if err != nil {
		return nil, err
	}

	mx := newMetrics()

	w.collectCPU(mx, scraped)

	return stm.ToMap(mx), nil
}

func (w *WMI) collectCPU(mx *metrics, pms prometheus.Metrics) {
	collectCPUCoresCStates(mx, pms)
	collectCPUCoresDPCs(mx, pms)
	collectCPUCoresInterrupts(mx, pms)
	collectCPUCoresUsage(mx, pms)

	for _, c := range mx.CPU.Cores {
		mx.CPU.Usage.User.Add(c.Usage.User.Value())
		mx.CPU.Usage.Privileged.Add(c.Usage.Privileged.Value())
		mx.CPU.Usage.Interrupt.Add(c.Usage.Interrupt.Value())
		mx.CPU.Usage.Idle.Add(c.Usage.Idle.Value())
		mx.CPU.Usage.DPC.Add(c.Usage.DPC.Value())
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
			core.Usage.DPC.Set(value)
		case "idle":
			core.Usage.Idle.Set(value)
		case "interrupt":
			core.Usage.Interrupt.Set(value)
		case "privileged":
			core.Usage.Privileged.Set(value)
		case "user":
			core.Usage.User.Set(value)
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

func findCollector(pms prometheus.Metrics, colName string) (exist, success bool) {
	for _, pm := range pms.FindByName("wmi_exporter_collector_success") {
		name := pm.Labels.Get("collector")
		if name == "" {
			break
		}
		if name == colName {
			return true, pm.Value == 1
		}
	}
	return false, false
}

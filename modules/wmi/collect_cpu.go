// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import (
	"sort"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	collectorCPU = "cpu"

	metricCPUCStateTotal     = "windows_cpu_cstate_seconds_total"
	metricCPUDPCsTotal       = "windows_cpu_dpcs_total"
	metricCPUInterruptsTotal = "windows_cpu_interrupts_total"
	metricCPUTimeTotal       = "windows_cpu_time_total"
)

func doCollectCPU(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, collectorCPU)
	return enabled && success
}

func collectCPU(pms prometheus.Metrics) *cpuMetrics {
	if !doCollectCPU(pms) {
		return nil
	}

	cm := &cpuMetrics{}
	collectCPUCoresCStates(cm, pms)
	collectCPUCoresDPCs(cm, pms)
	collectCPUCoresInterrupts(cm, pms)
	collectCPUCoresUsage(cm, pms)
	collectCPUSummary(cm)
	sortCPUCores(&cm.Cores)
	return cm
}

func collectCPUCoresCStates(cm *cpuMetrics, pms prometheus.Metrics) {
	var core *cpuCore

	for _, pm := range pms.FindByName(metricCPUCStateTotal) {
		coreID := pm.Labels.Get("core")
		state := pm.Labels.Get("state")
		if coreID == "" || state == "" {
			continue
		}

		if core == nil || core.ID != coreID {
			core = cm.Cores.get(coreID)
		}

		switch state {
		default:
		case "c1":
			core.C1 = pm.Value
		case "c2":
			core.C2 = pm.Value
		case "c3":
			core.C3 = pm.Value
		}
	}
}

func collectCPUCoresInterrupts(cm *cpuMetrics, pms prometheus.Metrics) {
	var core *cpuCore

	for _, pm := range pms.FindByName(metricCPUInterruptsTotal) {
		coreID := pm.Labels.Get("core")
		if coreID == "" {
			continue
		}

		if core == nil || core.ID != coreID {
			core = cm.Cores.get(coreID)
		}

		core.InterruptsTotal = pm.Value
	}
}

func collectCPUCoresUsage(cm *cpuMetrics, pms prometheus.Metrics) {
	var core *cpuCore

	for _, pm := range pms.FindByName(metricCPUTimeTotal) {
		coreID := pm.Labels.Get("core")
		mode := pm.Labels.Get("mode")
		if coreID == "" || mode == "" {
			continue
		}

		if core == nil || core.ID != coreID {
			core = cm.Cores.get(coreID)
		}

		switch mode {
		default:
		case "dpc":
			core.DPC = pm.Value
		case "idle":
			core.Idle = pm.Value
		case "interrupt":
			core.Interrupt = pm.Value
		case "privileged":
			core.Privileged = pm.Value
		case "user":
			core.User = pm.Value
		}
	}
}

func collectCPUCoresDPCs(cm *cpuMetrics, pms prometheus.Metrics) {
	var core *cpuCore

	for _, pm := range pms.FindByName(metricCPUDPCsTotal) {
		coreID := pm.Labels.Get("core")
		if coreID == "" {
			continue
		}

		if core == nil || core.ID != coreID {
			core = cm.Cores.get(coreID)
		}

		core.DPCsTotal = pm.Value
	}
}

func collectCPUSummary(cm *cpuMetrics) {
	for _, c := range cm.Cores {
		cm.User += c.User
		cm.Privileged += c.Privileged
		cm.Interrupt += c.Interrupt
		cm.Idle += c.Idle
		cm.DPC += c.DPC
	}
}

func sortCPUCores(cores *cpuCores) {
	sort.Slice(*cores, func(i, j int) bool { return (*cores)[i].id < (*cores)[j].id })
}

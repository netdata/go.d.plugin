// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	serviceStartModeAuto = iota
	serviceStartModeBoot
	serviceStartModeDisabled
	serviceStartModeManual
	serviceStartModeSystem
)

const (
	serviceStateContinuePending = iota
	serviceStatePausePending
	serviceStatePaused
	serviceStateRunning
	serviceStateStartPending
	serviceStateStopPending
	serviceStateStopped
	serviceStateUnknown
)

const (
	serviceStatusDegraded = iota
	serviceStatusError
	serviceStatusLostConn
	serviceStatusNoContact
	serviceStatusNonRecover
	serviceStatusOK
	serviceStatusPredFail
	serviceStatusService
	serviceStatusStarting
	serviceStatusStopping
	serviceStatusStressed
	serviceStatusUnkown
)

const (
	collectorService = "service"

	metricServiceStartMode = "windows_service_start_mode"
	metricServiceState     = "windows_service_state"
	metricServiceStatus    = "windows_service_status"
)

func doCollectMetrics(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, collectorService)
	return enabled && success
}

func collectService(pms prometheus.Metrics) *servicesMetrics {
	if !doCollectProcess(pms) {
		return nil
	}

	servs := &servicesMetrics{servs: make(map[string]*serviceMetrics)}
	collectServiceState(servs, pms)

	return servs
}

func collectServiceState(procs *servicesMetrics, pms prometheus.Metrics) {
	var serv *serviceMetrics
	for _, pm := range pms.FindByName(metricServiceState) {
		name := pm.Labels.Get("name")
		state := pm.Labels.Get("state")
		if name == "" {
			continue
		}

		if serv == nil || serv.ID != name {
			serv = procs.get(name)
		}

		if pm.Value == 1 {
			selectServiceState(&serv.state, state)
		}
	}
}

func collectServiceStatus(procs *servicesMetrics, pms prometheus.Metrics) {
	var serv *serviceMetrics
	for _, pm := range pms.FindByName(metricServiceStatus) {
		name := pm.Labels.Get("name")
		status := pm.Labels.Get("status")
		if name == "" {
			continue
		}

		if serv == nil || serv.ID != name {
			serv = procs.get(name)
		}

		if pm.Value == 1 {
			serv.status = float64(selectServiceStatus(status))
		}
	}
}

func selectServiceState(sse *serviceState, name string) {
	sse.continuePending = boolToFloat64(name == "continue pending")
	sse.pausePending = boolToFloat64(name == "pause pending")
	sse.paused = boolToFloat64(name == "paused")
	sse.running = boolToFloat64(name == "running")
	sse.startPending = boolToFloat64(name == "start pending")
	sse.stopPending = boolToFloat64(name == "stop pending")
	sse.stopped = boolToFloat64(name == "stopped")
	sse.unknown = boolToFloat64(name == "unknown")
}

func selectServiceStatus(name string) int32 {
	switch name {
	case "degraded":
		return serviceStatusDegraded
	case "error":
		return serviceStatusError
	case "lost comm":
		return serviceStatusLostConn
	case "no contact":
		return serviceStatusNoContact
	case "nonrecover":
		return serviceStatusNonRecover
	case "ok":
		return serviceStatusOK
	case "pred fail":
		return serviceStatusPredFail
	case "service":
		return serviceStatusService
	case "starting":
		return serviceStatusStarting
	case "stopping":
		return serviceStatusStopping
	case "stressed":
		return serviceStatusStressed
	}

	return serviceStateUnknown
}

func boolToFloat64(v bool) float64 {
	if v {
		return 1
	}
	return 0
}

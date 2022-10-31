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
	collectServiceStatus(servs, pms)

	return servs
}

func collectServiceState(procs *servicesMetrics, pms prometheus.Metrics) {
	var serv *serviceMetrics
	for _, pm := range pms.FindByName(metricServiceState) {
		name := pm.Labels.Get("name")
		state := pm.Labels.Get("state")
		if name == "" || state == "" || pm.Value == 0 {
			continue
		}

		if serv == nil || serv.ID != name {
			serv = procs.get(name)
		}

		if pm.Value == 1 {
			setServiceState(&serv.state, state)
		}
	}
}

func collectServiceStatus(procs *servicesMetrics, pms prometheus.Metrics) {
	var serv *serviceMetrics
	for _, pm := range pms.FindByName(metricServiceStatus) {
		name := pm.Labels.Get("name")
		status := pm.Labels.Get("status")
		if name == "" || status == "" || pm.Value == 0 {
			continue
		}

		if serv == nil || serv.ID != name {
			serv = procs.get(name)
		}

		setServiceStatus(&serv.status, status)
	}
}

func setServiceState(state *serviceState, name string) {
	state.continuePending = boolToFloat64(name == "continue pending")
	state.pausePending = boolToFloat64(name == "pause pending")
	state.paused = boolToFloat64(name == "paused")
	state.running = boolToFloat64(name == "running")
	state.startPending = boolToFloat64(name == "start pending")
	state.stopPending = boolToFloat64(name == "stop pending")
	state.stopped = boolToFloat64(name == "stopped")
	state.unknown = boolToFloat64(name == "unknown")
}

func setServiceStatus(status *serviceStatus, name string) {
	status.degraded = boolToFloat64(name == "degraded")
	status.errors = boolToFloat64(name == "error")
	status.lostComm = boolToFloat64(name == "lost comm")
	status.noContact = boolToFloat64(name == "no contact")
	status.nonRecover = boolToFloat64(name == "nonrecover")
	status.ok = boolToFloat64(name == "ok")
	status.predFail = boolToFloat64(name == "pred fail")
	status.service = boolToFloat64(name == "service")
	status.starting = boolToFloat64(name == "starting")
	status.stopping = boolToFloat64(name == "stopping")
	status.stressed = boolToFloat64(name == "stressed")
	status.unknown = boolToFloat64(name == "unknown")
}

func boolToFloat64(v bool) float64 {
	if v {
		return 1
	}
	return 0
}

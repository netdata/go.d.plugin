// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import "github.com/netdata/go.d.plugin/pkg/prometheus"

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
	collectServiceStartMode(servs, pms)
	collectServiceState(servs, pms)

	return servs
}

func collectServiceStartMode(procs *servicesMetrics, pms prometheus.Metrics) {
	var serv *serviceMetrics
	for _, pm := range pms.FindByName(metricServiceStartMode) {
		name := pm.Labels.Get("name")
		if name == "" {
			continue
		}

		if serv == nil || serv.ID != name {
			serv = procs.get(name)
		}

		if pm.Value == 1 {
			serv.startMode = float64(selectServiceStartMode(name))
		}
	}
}

func collectServiceState(procs *servicesMetrics, pms prometheus.Metrics) {
	var serv *serviceMetrics
	for _, pm := range pms.FindByName(metricServiceState) {
		name := pm.Labels.Get("name")
		if name == "" {
			continue
		}

		if serv == nil || serv.ID != name {
			serv = procs.get(name)
		}

		if pm.Value == 1 {
			serv.state = float64(selectServiceState(name))
		}
	}
}

func collectServiceStatus(procs *servicesMetrics, pms prometheus.Metrics) {
	var serv *serviceMetrics
	for _, pm := range pms.FindByName(metricServiceStatus) {
		name := pm.Labels.Get("name")
		if name == "" {
			continue
		}

		if serv == nil || serv.ID != name {
			serv = procs.get(name)
		}

		if pm.Value == 1 {
			serv.status = float64(selectServiceStatus(name))
		}
	}
}

func selectServiceStartMode(name string) int32 {
	switch name {
	case "auto":
		return serviceStartModeAuto
	case "boot":
		return serviceStartModeBoot
	case "disabled":
		return serviceStartModeDisabled
	case "manual":
		return serviceStartModeManual
	}

	return serviceStartModeSystem
}

func selectServiceState(name string) int32 {
	switch name {
	case "continue pending":
		return serviceStateContinuePending
	case "pause pending":
		return serviceStatePausePending
	case "paused":
		return serviceStatePaused
	case "running":
		return serviceStateRunning
	case "start pending":
		return serviceStateStartPending
	case "stop pending":
		return serviceStateStopPending
	case "stopped":
		return serviceStateStopped
	}

	return serviceStateUnknown
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

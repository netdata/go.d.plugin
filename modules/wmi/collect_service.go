// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

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

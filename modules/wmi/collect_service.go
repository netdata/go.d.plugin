package wmi

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	collectorService = "service"

	metricServiceState  = "windows_service_state"
	metricServiceStatus = "windows_service_status"
)

func doCollectService(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, collectorService)
	return enabled && success
}

func collectService(pms prometheus.Metrics) *servicesMetrics {
	if !doCollectService(pms) {
		return nil
	}

	svcs := &servicesMetrics{svcs: make(map[string]*serviceMetrics)}
	collectServiceState(svcs, pms)
	collectServiceStatus(svcs, pms)

	return svcs
}

func collectServiceState(svcs *servicesMetrics, pms prometheus.Metrics) {
	var svc *serviceMetrics
	for _, pm := range pms.FindByName(metricServiceState) {
		name := pm.Labels.Get("name")
		if name == "" {
			continue
		}

		if svc == nil || svc.ID != name {
			svc = svcs.get(name)
		}

		if pm.Value == 0 {
			continue
		}

		switch pm.Labels.Get("state") {
		case "stopped":
			svc.state.stopped = pm.Value
		case "start pending":
			svc.state.startPending = pm.Value
		case "stop pending":
			svc.state.stopPending = pm.Value
		case "running":
			svc.state.running = pm.Value
		case "continue pending":
			svc.state.continuePending = pm.Value
		case "pause pending":
			svc.state.pausePending = pm.Value
		case "paused":
			svc.state.paused = pm.Value
		case "unknown":
			svc.state.unknown = pm.Value
		}

	}
}

func collectServiceStatus(svcs *servicesMetrics, pms prometheus.Metrics) {
	var svc *serviceMetrics
	for _, pm := range pms.FindByName(metricServiceStatus) {
		name := pm.Labels.Get("name")
		if name == "" {
			continue
		}

		if svc == nil || svc.ID != name {
			svc = svcs.get(name)
		}

		if pm.Value == 0 {
			continue
		}

		switch pm.Labels.Get("status") {
		case "ok":
			svc.status.ok = pm.Value
		case "error":
			svc.status.error = pm.Value
		case "degraded":
			svc.status.degraded = pm.Value
		case "unknown":
			svc.status.unknown = pm.Value
		case "pred fail":
			svc.status.predFail = pm.Value
		case "starting":
			svc.status.starting = pm.Value
		case "stopping":
			svc.status.stopping = pm.Value
		case "service":
			svc.status.service = pm.Value
		case "stressed":
			svc.status.stressed = pm.Value
		case "nonrecover":
			svc.status.nonRecover = pm.Value
		case "no contact":
			svc.status.noContact = pm.Value
		case "lost comm":
			svc.status.lostComm = pm.Value
		}
	}
}

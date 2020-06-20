package wmi

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	collectorLogon = "logon"

	metricLogonType = "windows_logon_logon_type"
)

func doCollectLogon(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, collectorLogon)
	return enabled && success
}

func collectLogon(pms prometheus.Metrics) *logonMetrics {
	if !doCollectLogon(pms) {
		return nil
	}

	var lm logonMetrics
	collectLogonSessionsByType(&lm, pms)
	return &lm
}

func collectLogonSessionsByType(lm *logonMetrics, pms prometheus.Metrics) {
	for _, pm := range pms.FindByName(metricLogonType) {
		logonType := pm.Labels.Get("status")
		assignLogonMetric(lm, logonType, pm.Value)
	}
}

func assignLogonMetric(lm *logonMetrics, logonType string, value float64) {
	switch logonType {
	default:
	case "system":
		lm.Type.System = value
	case "interactive":
		lm.Type.Interactive = value
	case "network":
		lm.Type.Network = value
	case "batch":
		lm.Type.Batch = value
	case "service":
		lm.Type.Service = value
	case "proxy":
		lm.Type.Proxy = value
	case "unlock":
		lm.Type.Unlock = value
	case "network_clear_text":
		lm.Type.NetworkCleartext = value
	case "new_credentials":
		lm.Type.NewCredentials = value
	case "remote_interactive":
		lm.Type.RemoteInteractive = value
	case "cached_interactive":
		lm.Type.CachedInteractive = value
	case "cached_remote_interactive":
		lm.Type.CachedRemoteInteractive = value
	case "cached_unlock":
		lm.Type.CachedUnlock = value
	}
}

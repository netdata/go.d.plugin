package vernemq

import (
	"errors"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/stm"
)

func isValidVerneMQMetrics(pms prometheus.Metrics) bool {
	return pms.FindByName(metricPUBLISHError).Len() > 0 && pms.FindByName(metricRouterSubscriptions).Len() > 0
}

func (v *VerneMQ) collect() (map[string]int64, error) {
	pms, err := v.prom.Scrape()
	if err != nil {
		return nil, err
	}

	if !isValidVerneMQMetrics(pms) {
		return nil, errors.New("returned metrics aren't VerneMQ metrics")
	}

	mx := collect(pms)

	return stm.ToMap(mx), nil
}

func collect(pms prometheus.Metrics) map[string]float64 {
	mx := make(map[string]float64)
	collectSockets(mx, pms)
	collectQueues(mx, pms)
	collectSubscriptions(mx, pms)
	collectErlangVM(mx, pms)
	collectBandwidth(mx, pms)
	collectRetain(mx, pms)
	collectCluster(mx, pms)
	collectUptime(mx, pms)

	collectAUTH(mx, pms)
	collectCONNECT(mx, pms)
	collectDISCONNECT(mx, pms)
	collectSUBSCRIBE(mx, pms)
	collectUNSUBSCRIBE(mx, pms)
	collectPUBLISH(mx, pms)
	collectPING(mx, pms)
	collectMQTTInvalidMsgSize(mx, pms)
	return mx
}

func collectCONNECT(mx map[string]float64, pms prometheus.Metrics) {
	pms = pms.FindByNames(
		metricCONNECTReceived,
		metricCONNACKSent,
	)
	collectMQTT(mx, pms)
}

func collectDISCONNECT(mx map[string]float64, pms prometheus.Metrics) {
	pms = pms.FindByNames(
		metricDISCONNECTReceived,
		metricDISCONNECTSent,
	)
	collectMQTT(mx, pms)
}

func collectPUBLISH(mx map[string]float64, pms prometheus.Metrics) {
	pms = pms.FindByNames(
		metricPUBACKReceived,
		metricPUBACKSent,
		metricPUBACKInvalid,

		metricPUBCOMPReceived,
		metricPUBCOMPSent,
		metricPUNCOMPInvalid,

		metricPUBSLISHReceived,
		metricPUBSLIHSent,
		metricPUBLISHError,
		metricPUBLISHAuthError,

		metricPUBRECReceived,
		metricPUBRECSent,
		metricPUBRECInvalid,

		metricPUBRELReceived,
		metricPUBRELSent,
	)
	collectMQTT(mx, pms)
}

func collectSUBSCRIBE(mx map[string]float64, pms prometheus.Metrics) {
	pms = pms.FindByNames(
		metricSUBSCRIBEReceived,
		metricSUBACKSent,
		metricSUBSCRIBEError,
		metricSUBSCRIBEAuthError,
	)
	collectMQTT(mx, pms)
}

func collectUNSUBSCRIBE(mx map[string]float64, pms prometheus.Metrics) {
	pms = pms.FindByNames(
		metricUNSUBSCRIBEReceived,
		metricUNSUBACKSent,
		metricUNSUBSCRIBEError,
	)
	collectMQTT(mx, pms)
}

func collectPING(mx map[string]float64, pms prometheus.Metrics) {
	pms = pms.FindByNames(
		metricPINGREQReceived,
		metricPINGRESPSent,
	)
	collectMQTT(mx, pms)
}

func collectAUTH(mx map[string]float64, pms prometheus.Metrics) {
	pms = pms.FindByNames(
		metricAUTHReceived,
		metricAUTHSent,
	)
	collectMQTT(mx, pms)
}

func collectMQTTInvalidMsgSize(mx map[string]float64, pms prometheus.Metrics) {
	pms = pms.FindByName(metricMQTTInvalidMsgSizeError)
	collectMQTT(mx, pms)
}

func collectSockets(mx map[string]float64, pms prometheus.Metrics) {
	pms = pms.FindByNames(
		metricSocketClose,
		metricSocketCloseTimeout,
		metricSocketError,
		metricSocketOpen,
		metricClientKeepaliveExpired,
	)
	collectNonMQTT(mx, pms)
	mx["open_sockets"] = mx[metricSocketOpen] - mx[metricSocketClose]
}

func collectQueues(mx map[string]float64, pms prometheus.Metrics) {
	pms = pms.FindByNames(
		metricQueueInitializedFromStorage,
		metricQueueMessageDrop,
		metricQueueMessageExpired,
		metricQueueMessageIn,
		metricQueueMessageOut,
		metricQueueMessageUnhandled,
		metricQueueProcesses,
		metricQueueSetup,
		metricQueueTeardown,
	)
	collectNonMQTT(mx, pms)
	mx["queue_messages_current"] = calcQueueMessagesCurrent(mx)
}

func calcQueueMessagesCurrent(mx map[string]float64) float64 {
	undelivered := mx[metricQueueMessageDrop] + mx[metricQueueMessageExpired] + mx[metricQueueMessageUnhandled]
	out := mx[metricQueueMessageOut]
	in := mx[metricQueueMessageIn]
	return in - (out + undelivered)
}

func collectSubscriptions(mx map[string]float64, pms prometheus.Metrics) {
	pms = pms.FindByNames(
		metricRouterMatchesLocal,
		metricRouterMatchesRemote,
		metricRouterMemory,
		metricRouterSubscriptions,
	)
	collectNonMQTT(mx, pms)
}

func collectErlangVM(mx map[string]float64, pms prometheus.Metrics) {
	collectSchedulersUtilization(mx, pms)
	pms = pms.FindByNames(
		metricSystemContextSwitches,
		metricSystemGCCount,
		metricSystemIOIn,
		metricSystemIOOut,
		metricSystemProcessCount,
		metricSystemReductions,
		metricSystemRunQueue,
		metricSystemUtilization,
		metricSystemWordsReclaimedByGC,
		metricVMMemoryProcesses,
		metricVMMemorySystem,
	)
	collectNonMQTT(mx, pms)
}

func collectSchedulersUtilization(mx map[string]float64, pms prometheus.Metrics) {
	var exit bool
	for _, pm := range pms {
		switch {
		case isSchedulerUtilizationMetric(pm):
			mx[pm.Name()] += pm.Value
			exit = true
		case exit:
			return
		}
	}
}

func collectBandwidth(mx map[string]float64, pms prometheus.Metrics) {
	pms = pms.FindByNames(
		metricBytesReceived,
		metricBytesSent,
	)
	collectNonMQTT(mx, pms)
}

func collectRetain(mx map[string]float64, pms prometheus.Metrics) {
	pms = pms.FindByNames(
		metricRetainMemory,
		metricRetainMessages,
	)
	collectNonMQTT(mx, pms)
}

func collectCluster(mx map[string]float64, pms prometheus.Metrics) {
	pms = pms.FindByNames(
		metricClusterBytesDropped,
		metricClusterBytesReceived,
		metricClusterBytesSent,
		metricNetSplitDetected,
		metricNetSplitResolved,
	)
	collectNonMQTT(mx, pms)
	mx["netsplit_unresolved"] = mx[metricNetSplitDetected] - mx[metricNetSplitResolved]
}

func collectUptime(mx map[string]float64, pms prometheus.Metrics) {
	pms = pms.FindByName(metricSystemWallClock)
	collectNonMQTT(mx, pms)
}

func collectNonMQTT(mx map[string]float64, pms prometheus.Metrics) {
	for _, pm := range pms {
		mx[pm.Name()] += pm.Value
	}
}

func collectMQTT(mx map[string]float64, pms prometheus.Metrics) {
	for _, pm := range pms {
		if !isMQTTMetric(pm) {
			continue
		}
		version := versionLabelValue(pm)
		if version == "" {
			continue
		}

		mx[pm.Name()] += pm.Value
		mx[join(pm.Name(), "v", version)] += pm.Value

		if reason := reasonCodeLabelValue(pm); reason != "" {
			mx[join(pm.Name(), reason)] += pm.Value
			mx[join(pm.Name(), "v", version, reason)] += pm.Value
		}
	}
}

func isMQTTMetric(pm prometheus.Metric) bool {
	return strings.HasPrefix(pm.Name(), "mqtt_")
}

func isSchedulerUtilizationMetric(pm prometheus.Metric) bool {
	return strings.HasPrefix(pm.Name(), "system_utilization_scheduler_")
}

func reasonCodeLabelValue(pm prometheus.Metric) string {
	if v := pm.Labels.Get("reason_code"); v != "" {
		return v
	}
	// "mqtt_connack_sent" v4 has return_code
	return pm.Labels.Get("return_code")
}

func versionLabelValue(pm prometheus.Metric) string {
	return pm.Labels.Get("mqtt_version")
}

func isReasonCodeNotSuccess(name, reason string) bool {
	switch name {
	case
		metricCONNACKSent,
		metricPUBACKReceived,
		metricPUBACKSent,
		metricPUBCOMPReceived,
		metricPUBCOMPSent,
		metricPUBRECReceived,
		metricPUBRECSent,
		metricPUBRELReceived,
		metricPUBRELSent:
		return reason != "success"
	case
		metricDISCONNECTReceived,
		metricDISCONNECTSent:
		return reason != "normal_disconnect"
	}
	return false
}

func join(a, b string, rest ...string) string {
	v := a + "_" + b
	switch len(rest) {
	case 0:
		return v
	default:
		return join(v, rest[0], rest[1:]...)
	}
}

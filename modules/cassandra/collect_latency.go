// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	collectorLatency = "Read"
)

func doCollectLatency(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, metricRequestType, collectorLatency, false)
	return enabled && success
}

func collectLatency(pms prometheus.Metrics) *latency {
	if !doCollectLatency(pms) {
		return nil
	}

	var la latency
	collectLatencyByType(&la, pms)

	return &la
}

func collectLatencyByType(la *latency, pms prometheus.Metrics) {
	var total latency
	for _, pm := range pms.FindByName(metricRequestType) {
		metricName := pm.Labels.Get("name")
		metricScope := pm.Labels.Get("scope")
		// We also have a latency specific for Read/Write, but total
		// is already including it. The actual code was written considering
		// that we can show they separated one day.
		if metricName == "TotalLatency" {
			assignLatencyMetric(&total, metricScope, pm.Value)
		}
	}
	la.read_latency = total.read_latency
	la.write_latency = total.write_latency
}

func assignLatencyMetric(la *latency, scope string, value float64) {
	switch scope {
	default:
	case "Read":
		la.read_latency += int64(value)
	case "Write":
		la.write_latency += int64(value)
	}
}

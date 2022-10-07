// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	collectorLatency = "ReadLatency"
)

func doCollectLatency(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, metricTableType, collectorLatency, false)
	return enabled && success
}

func collectLatency(pms prometheus.Metrics) *LATENCY {
	if !doCollectLatency(pms) {
		return nil
	}

	var la LATENCY
	collectLatencyByType(&la, pms)

	return &la
}

func collectLatencyByType(la *LATENCY, pms prometheus.Metrics) {
	var total LATENCY
	for _, pm := range pms.FindByName(metricTableType) {
		metricName := pm.Labels.Get("name")
		assignLatencyMetric(&total, metricName, pm.Value)
	}
	la.read = total.read
	la.write = total.write
}

func assignLatencyMetric(la *LATENCY, scope string, value float64) {
	switch scope {
	default:
	case "ReadLatency":
		la.read += int64(value)
	case "WriteLatency":
		la.write += int64(value)
	}
}

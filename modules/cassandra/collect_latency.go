// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	collectorLatency  = "Read"
	metricLatencyType = "org_apache_cassandra_metrics_clientrequest_count"
)

func doCollectLatency(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, metricLatencyType, collectorLatency, false)
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
	for _, pm := range pms.FindByName(metricLatencyType) {
		metricName := pm.Labels.Get("name")
		metricScope := pm.Labels.Get("scope")
		// We also have a latency specific for Read/Write, but total
		// is already including it. The actual code was written considering
		// that we can show they separated one day.
		if metricName == "TotalLatency" {
			assignLatencyMetric(&total, metricScope, pm.Value)
		}
	}
	la.read = total.read
	la.write = total.write
}

func assignLatencyMetric(la *LATENCY, scope string, value float64) {
	switch scope {
	default:
	case "Read":
		la.read += int64(value)
	case "Write":
		la.write += int64(value)
	}
}

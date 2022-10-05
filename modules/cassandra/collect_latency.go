// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	collectorLatency = "table"
	metricLatencyType = "org_apache_cassandra_metrics_table_count"
)

func doCollectLatency(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, collectorLatency)
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
	for _, pm := range pms.FindByName(metricTroughputType) {
		metricScope := pm.Labels.Get("scope")
		metricName := pm.Labels.Get("name")
		if metricScope == "events" {
			assignLatencyMetric(la, metricName, pm.Value)
		}
	}
}

func assignLatencyMetric(la *LATENCY, scope string, value float64) {
	switch scope {
	default:
	case "ReadLatency":
		la.read = int64(value)
	case "WriteLatency":
		la.write = int64(value)
	}
}

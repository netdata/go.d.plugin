// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	collectorTroughput  = "Read"
	metricTroughputType = "org_apache_cassandra_metrics_clientrequest_oneminuterate"
)

func doCollectThroughput(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, metricTroughputType, collectorTroughput, true)
	return enabled && success
}

func collectThroughput(pms prometheus.Metrics) *throughput {
	if !doCollectThroughput(pms) {
		return nil
	}

	var tp throughput
	collectThroughputByType(&tp, pms)

	return &tp
}

func collectThroughputByType(tp *throughput, pms prometheus.Metrics) {
	for _, pm := range pms.FindByName(metricTroughputType) {
		metricScope := pm.Labels.Get("scope")
		assignThroughputMetric(tp, metricScope, pm.Value)
	}
}

func assignThroughputMetric(tp *throughput, scope string, value float64) {
	switch scope {
	default:
	case "Read":
		tp.read = int64(value * 100)
	case "Write":
		tp.write = int64(value)
	}
}

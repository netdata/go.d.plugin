// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	collectorTroughput = "clientrequest"
	metricTroughputType = "org_apache_cassandra_metrics_clientrequest_count"
)

func doCollectThroughput(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, collectorTroughput)
	return enabled && success
}

func collectThroughput(pms prometheus.Metrics) *THROUGHPUT {
	if !doCollectThroughput(pms) {
		return nil
	}

	var tp THROUGHPUT
    collectThroughputByType(&tp, pms)

    return &tp
}

func collectThroughputByType(tp *THROUGHPUT, pms prometheus.Metrics) {
	for _, pm := range pms.FindByName(metricTroughputType) {
		metricScope := pm.Labels.Get("scope")
		metricName := pm.Labels.Get("name")
		if metricName == "Latency" {
			assignThroughputMetric(tp, metricScope, pm.Value)
		}
	}
}

func assignThroughputMetric(tp *THROUGHPUT, scope string, value float64) {
	switch scope {
	default:
	case "Read":
		tp.read = int64(value);
	case "Write":
		tp.write = int64(value);
	}
}

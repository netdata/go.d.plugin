// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	collectorTimeout     = "Timeouts"
	collectorUnavailable = "Unavailables"
)

func doCollectRequestError(pms prometheus.Metrics, metric string) bool {
	var tester string
	if metric == collectorTimeout {
		tester = collectorTimeout
	} else {
		tester = collectorUnavailable
	}
	enabled, success := checkCollector(pms, metricRequestType, tester, false)
	return enabled && success
}

func collectRequestError(pms prometheus.Metrics, metric string) *REQUEST_ERROR {
	if !doCollectRequestError(pms, metric) {
		return nil
	}

	var re REQUEST_ERROR
	collectRequestErrorByType(&re, pms)

	return &re
}

func collectRequestErrorByType(re *REQUEST_ERROR, pms prometheus.Metrics) {
	for _, pm := range pms.FindByName(metricRequestType) {
		metricName := pm.Labels.Get("name")
		scopeName := pm.Labels.Get("scope")
		if metricName == "Timeouts" || metricName == "Unavailables" {
			assignRequestErrorMetric(re, scopeName, pm.Value)
		}
	}
}

func assignRequestErrorMetric(re *REQUEST_ERROR, scope string, value float64) {
	switch scope {
	default:
	case "Read":
		re.read = int64(value)
	case "Write":
		re.write = int64(value)
	}
}

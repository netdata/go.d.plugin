// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	collectorCache  = "HitRate"
	metricCacheType = "org_apache_cassandra_metrics_cache_value"
)

func doCollectCache(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, metricCacheType, collectorCache, false)
	return enabled && success
}

func collectCache(pms prometheus.Metrics) *CACHE {
	if !doCollectCache(pms) {
		return nil
	}

	var ca CACHE
	collectCacheByType(&ca, pms)

	return &ca
}

func collectCacheByType(ca *CACHE, pms prometheus.Metrics) {
	for _, pm := range pms.FindByName(metricCacheType) {
		metricScope := pm.Labels.Get("scope")
		metricName := pm.Labels.Get("name")
		if metricScope == "KeyCache" {
			assignCacheMetric(ca, metricName, pm.Value*100.0)
		}
	}
}

func assignCacheMetric(ca *CACHE, scope string, value float64) {
	switch scope {
	default:
	case "HitRate":
		ca.hit = int64(value)
	}
}

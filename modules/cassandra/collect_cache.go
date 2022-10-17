// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	collectorCache  = "Hits"
	metricCacheType = "org_apache_cassandra_metrics_cache_count"
)

func doCollectCache(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, metricCacheType, collectorCache, false)
	return enabled && success
}

func collectCache(pms prometheus.Metrics) *cache {
	if !doCollectCache(pms) {
		return nil
	}

	var ca cache
	collectCacheByType(&ca, pms)

	return &ca
}

func collectCacheByType(ca *cache, pms prometheus.Metrics) {
	var hit, requests, hits float64
	for _, pm := range pms.FindByName(metricCacheType) {
		metricName := pm.Labels.Get("name")
		// Code prepared to collect more metrics from Cache.
		if metricName == "Hits" {
			hits = pm.Value
		} else if metricName == "Requests" {
			requests = pm.Value
		}
	}
	hit = (hits / requests) * 100.0
	ca.hit = int64(hit)
}

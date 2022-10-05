// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

import (
	"errors"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/prometheus/prometheus/model/labels"
)

const (
	metricCollectorSuccess  = "peers"
)

func isValidWindowsExporterMetrics(pms prometheus.Metrics) bool {
	return pms.FindByName(metricCollectorSuccess).Len() > 0
}

func (c *Cassandra) collect() (map[string]int64, error) {
	pms, err := c.prom.Scrape()
	if err != nil {
			return nil, err
	}

	if !isValidWindowsExporterMetrics(pms) {
			return nil, errors.New("collected metrics aren't windows_exporter metrics")
	}

	mx := collect(pms)
	c.updateCharts(mx)

	return stm.ToMap(mx), nil
}

func collect(pms prometheus.Metrics) *metrics {
	mx := metrics{
		throughput:  collectThroughput(pms),
		latency:  collectLatency(pms),
		cache:  collectCache(pms),
		disk:  collectDisk(pms),
	}

	return &mx
}

func checkCollector(pms prometheus.Metrics, name string) (enabled, success bool) {
	m, err := labels.NewMatcher(labels.MatchEqual, "collector", name)
	if err != nil {
		panic(err)
	}

	pms = pms.FindByName(metricCollectorSuccess)
	ms := pms.Match(m)
	return ms.Len() > 0, ms.Max() == 1
}
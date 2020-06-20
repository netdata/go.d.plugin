package wmi

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	metricCollectorDuration = "windows_exporter_collector_duration_seconds"
	metricCollectorSuccess  = "windows_exporter_collector_success"
)

var collectorMetricsNames = []string{
	metricCollectorDuration,
	metricCollectorSuccess,
}

func collectCollection(pms prometheus.Metrics) *collectors {
	mx := &collectors{}
	for _, name := range collectorMetricsNames {
		collectCollectorMetric(mx, pms, name)
	}
	return mx
}

func collectCollectorMetric(mx *collectors, pms prometheus.Metrics, name string) {
	var col *collector

	for _, pm := range pms.FindByName(name) {
		colName := pm.Labels.Get("collector")
		if colName == "" {
			continue
		}

		if col == nil || col.ID != colName {
			col = mx.get(colName)
		}

		assignCollectorMetric(col, name, pm.Value)
	}
}

func assignCollectorMetric(col *collector, name string, value float64) {
	switch name {
	case metricCollectorDuration:
		col.Duration = value
	case metricCollectorSuccess:
		col.Success = value == 1
	}
}

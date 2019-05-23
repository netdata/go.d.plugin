package wmi

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	metricCollectorDuration = "wmi_exporter_collector_duration_seconds"
	metricCollectorSuccess  = "wmi_exporter_collector_success"
)

func collectCollection(mx *metrics, pms prometheus.Metrics) {
	mx.Collectors = &collectors{}
	collectCollectionDuration(mx, pms)
	collectCollectionSuccess(mx, pms)
}

func collectCollectionDuration(mx *metrics, pms prometheus.Metrics) {
	cr := newCollector("")
	for _, pm := range pms.FindByName(metricCollectorDuration) {
		name := pm.Labels.Get("collector")
		if name == "" {
			continue
		}
		if cr.ID != name {
			cr = mx.Collectors.get(name, true)
		}
		cr.Duration = pm.Value
	}
}

func collectCollectionSuccess(mx *metrics, pms prometheus.Metrics) {
	cr := newCollector("")
	for _, pm := range pms.FindByName(metricCollectorSuccess) {
		name := pm.Labels.Get("collector")
		if name == "" {
			continue
		}
		if cr.ID != name {
			cr = mx.Collectors.get(name, true)
		}
		cr.Success = pm.Value == 1
	}
}

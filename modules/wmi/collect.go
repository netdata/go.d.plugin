package wmi

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/prometheus/prometheus/pkg/labels"
)

const (
	collectorCPU    = "cpu"
	collectorNet    = "net"
	collectorMemory = "memory"
	collectorCS     = "cs"
	collectorOS     = "os"
	collectorSystem = "system"

	metricCollectorDuration = "wmi_exporter_collector_duration_seconds"
)

func (w *WMI) collect() (map[string]int64, error) {
	scraped, err := w.prom.Scrape()
	if err != nil {
		return nil, err
	}

	mx := newMetrics()

	w.collectScraped(mx, scraped)
	w.updateCharts(mx)

	return stm.ToMap(mx), nil
}

func (w *WMI) collectScraped(mx *metrics, scraped prometheus.Metrics) {
	collectCollectorsDuration(mx, scraped)

	enabled, success := checkCollector(scraped, collectorCPU)
	if enabled && success {
		mx.CPU = &cpu{}
		w.collectCPU(mx, scraped)
	}

	enabled, success = checkCollector(scraped, collectorNet)
	if enabled && success {
		mx.Net = &network{}
		w.collectNet(mx, scraped)
	}

	enabled, success = checkCollector(scraped, collectorMemory)
	if enabled && success {
		mx.Memory = &memory{}
		w.collectMemory(mx, scraped)
	}

	enabled, success = checkCollector(scraped, collectorCS)
	if enabled && success {
		mx.CS = &cs{}
		w.collectCS(mx, scraped)
	}

	enabled, success = checkCollector(scraped, collectorOS)
	if enabled && success {
		mx.OS = &os{}
		w.collectOS(mx, scraped)
	}

	enabled, success = checkCollector(scraped, collectorSystem)
	if enabled && success {
		mx.System = &system{}
		w.collectSystem(mx, scraped)
	}
}

func collectCollectorsDuration(mx *metrics, pms prometheus.Metrics) {
	for _, pm := range pms.FindByName(metricCollectorDuration) {
		name := pm.Labels.Get("collector")
		if name == "" {
			continue
		}
		mx.CollectDuration[name] = pm.Value
	}
}

func checkCollector(pms prometheus.Metrics, name string) (enabled, success bool) {
	m, err := labels.NewMatcher(labels.MatchEqual, "collector", name)
	if err != nil {
		panic(err)
	}
	ms := pms.Match(m)
	return ms.Len() > 0, ms.Max() == 1
}

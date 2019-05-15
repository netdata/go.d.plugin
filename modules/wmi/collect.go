package wmi

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/prometheus/prometheus/pkg/labels"
)

const (
	collectorCPU = "cpu"
	collectorNet = "net"
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
	enabled, success := findCollector(scraped, collectorCPU)
	if enabled && success {
		mx.CPU = &cpu{}
		w.collectCPU(mx, scraped)
	}

	enabled, success = findCollector(scraped, collectorNet)
	if enabled && success {
		mx.Net = &network{}
		w.collectNet(mx, scraped)
	}
}

func findCollector(pms prometheus.Metrics, name string) (enabled, success bool) {
	m, err := labels.NewMatcher(labels.MatchEqual, "collector", name)
	if err != nil {
		panic(err)
	}
	ms := pms.Match(m)
	return ms.Len() > 0, ms.Max() == 1
}

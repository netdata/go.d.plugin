package wmi

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/prometheus/prometheus/pkg/labels"
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
	collectCollection(mx, scraped)
	collectCPU(mx, scraped)
	collectOS(mx, scraped)
	collectMemory(mx, scraped)
	collectSystem(mx, scraped)
	collectNet(mx, scraped)
	collectLogicalDisk(mx, scraped)

	if mx.hasOS() && mx.hasMem() {
		v := sum(mx.OS.VisibleMemoryBytes, -mx.Memory.AvailableBytes)
		mx.Memory.UsedBytes = &v
	}
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

func sum(vs ...float64) (s float64) {
	for _, v := range vs {
		s += v
	}
	return s
}

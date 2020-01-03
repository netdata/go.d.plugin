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

	mx := collectScraped(scraped)
	w.updateCharts(mx)

	return stm.ToMap(mx), nil
}

func collectScraped(scraped prometheus.Metrics) *metrics {
	mx := metrics{
		CPU:         collectCPU(scraped),
		Memory:      collectMemory(scraped),
		Net:         collectNet(scraped),
		LogicalDisk: collectLogicalDisk(scraped),
		OS:          collectOS(scraped),
		System:      collectSystem(scraped),
		Logon:       collectLogon(scraped),
		Collectors:  collectCollection(scraped),
	}

	if mx.hasOS() && mx.hasMemory() {
		v := mx.OS.VisibleMemoryBytes - mx.Memory.AvailableBytes
		mx.Memory.UsedBytes = &v
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

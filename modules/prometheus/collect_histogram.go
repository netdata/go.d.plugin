package prometheus

import (
	"sort"
	"strconv"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

func (p *Prometheus) collectHistogram(mx map[string]int64, pms prometheus.Metrics, meta prometheus.Metadata) {
	if !pms[0].Labels.Has("le") {
		return
	}

	sortHistogram(pms)
	name := pms[0].Name()
	if !p.cache.has(name) {
		p.cache.put(name, newCacheEntry(nil))
	}

	defer p.cleanupStaleHistogramCharts(name)

	set := make(map[string]float64)
	cache := p.cache.get(name)

	for _, pm := range pms {
		chartID := defaultHistogramGrouping.chartID(pm)
		dimID := defaultHistogramGrouping.dimID(pm)
		dimName := defaultHistogramGrouping.dimName(pm)

		// {handler="/",le="0.1"} 1
		// {handler="/",le="0.2"} 2
		// {handler="/",le="0.4"} 3
		// le="0.4" = 3 - 2 (le="0.4" - le="0.2")
		if v, ok := set[chartID]; !ok {
			mx[dimID] = int64(pm.Value * precision)
		} else {
			mx[dimID] = int64((pm.Value - v) * precision)
		}
		set[chartID] = pm.Value

		if !cache.hasChart(chartID) {
			chart := histogramChart(chartID, pm, meta)
			cache.putChart(chartID, chart)
			if err := p.Charts().Add(chart); err != nil {
				p.Warning(err)
			}
		}
		if !cache.hasDim(dimID) {
			cache.putDim(dimID)
			chart := cache.getChart(chartID)
			dim := histogramChartDim(dimID, dimName)
			if err := chart.AddDim(dim); err != nil {
				p.Warning(err)
			}
			chart.MarkNotCreated()
		}
	}
}

func (p *Prometheus) cleanupStaleHistogramCharts(name string) {
	if !p.cache.has(name) {
		return
	}
	cache := p.cache.get(name)
	for _, chart := range cache.charts {
		if chart.Retries < 10 {
			continue
		}

		for _, dim := range chart.Dims {
			cache.removeDim(dim.ID)
			_ = chart.MarkDimRemove(dim.ID, true)
		}
		cache.removeChart(chart.ID)

		chart.MarkRemove()
		chart.MarkNotCreated()
	}
}

func sortHistogram(pms prometheus.Metrics) {
	sort.Slice(pms, func(i, j int) bool {
		iStr := pms[i].Labels.Get("le")
		jStr := pms[j].Labels.Get("le")
		if iStr == "+Inf" {
			return false
		}
		if jStr == "+Inf" {
			return true
		}
		iVal, _ := strconv.ParseFloat(iStr, 64)
		jVal, _ := strconv.ParseFloat(jStr, 64)
		return iVal < jVal
	})
}

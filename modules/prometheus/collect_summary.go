package prometheus

import (
	"math"
	"sort"
	"strconv"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

func (p *Prometheus) collectSummary(mx map[string]int64, pms prometheus.Metrics, meta prometheus.Metadata) {
	sortSummary(pms)
	name := pms[0].Name()
	if !p.cache.has(name) {
		p.cache.put(name, &cacheEntry{
			split:  newSummarySplit(),
			charts: make(chartsCache),
			dims:   make(dimsCache),
		})
	}

	cache := p.cache.get(name)

	for _, pm := range pms {
		if math.IsNaN(pm.Value) {
			continue
		}

		chartID := cache.split.chartID(pm)
		dimID := cache.split.dimID(pm)
		dimName := cache.split.dimName(pm)

		mx[dimID] = int64(pm.Value * precision)

		if !cache.hasChart(chartID) {
			chart := summaryChart(chartID, pm, meta)
			cache.putChart(chartID, chart)
			if err := p.Charts().Add(chart); err != nil {
				p.Warning(err)
			}
		}
		if !cache.hasDim(dimID, chartID) {
			cache.putDim(dimID, chartID)
			chart := cache.getChart(chartID)
			dim := summaryChartDimension(dimID, dimName)
			if err := chart.AddDim(dim); err != nil {
				p.Warning(err)
			}
			chart.MarkNotCreated()
		}
	}
}

func sortSummary(pms prometheus.Metrics) {
	sort.Slice(pms, func(i, j int) bool {
		iVal, _ := strconv.ParseFloat(pms[i].Labels.Get("quantile"), 64)
		jVal, _ := strconv.ParseFloat(pms[j].Labels.Get("quantile"), 64)
		return iVal < jVal
	})
}

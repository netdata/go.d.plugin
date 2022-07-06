// SPDX-License-Identifier: GPL-3.0-or-later

package prometheus

import (
	"math"
	"sort"
	"strconv"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

func (p *Prometheus) collectSummary(mx map[string]int64, pms prometheus.Metrics, meta prometheus.Metadata) {
	if !pms[0].Labels.Has("quantile") {
		return
	}

	sortSummary(pms)
	name := pms[0].Name()
	if !p.cache.has(name) {
		p.cache.put(name, newCacheEntry(nil))
	}

	defer p.cleanupStaleSummaryCharts(name)

	cache := p.cache.get(name)

	for _, pm := range pms {
		if math.IsNaN(pm.Value) {
			continue
		}

		chartID := defaultSummaryGrouping.chartID(pm)
		dimID := defaultSummaryGrouping.dimID(pm)
		dimName := defaultSummaryGrouping.dimName(pm)

		mx[dimID] = int64(pm.Value * precision)

		if !cache.hasChart(chartID) {
			chart := summaryChart(chartID, p.application(), pm, meta)
			cache.putChart(chartID, chart)
			if err := p.Charts().Add(chart); err != nil {
				p.Warning(err)
			}
		}
		if !cache.hasDim(dimID) {
			cache.putDim(dimID)
			chart := cache.getChart(chartID)
			dim := summaryChartDimension(dimID, dimName)
			if err := chart.AddDim(dim); err != nil {
				p.Warning(err)
			}
			chart.MarkNotCreated()
		}
	}
}

func (p *Prometheus) cleanupStaleSummaryCharts(name string) {
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

func sortSummary(pms prometheus.Metrics) {
	sort.Slice(pms, func(i, j int) bool {
		iVal, _ := strconv.ParseFloat(pms[i].Labels.Get("quantile"), 64)
		jVal, _ := strconv.ParseFloat(pms[j].Labels.Get("quantile"), 64)
		return iVal < jVal
	})
}

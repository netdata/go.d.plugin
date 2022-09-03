// SPDX-License-Identifier: GPL-3.0-or-later

package prometheus

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/prometheus/selector"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/prometheus/prometheus/model/labels"
)

// TODO: make it easier, very error-prone interface
type (
	collectCache map[string]*cacheEntry // map[metricName]
	cacheEntry   struct {
		groups groupCache
		charts chartsCache
		dims   dimsCache
	}
	groupCache struct {
		init       bool
		user       []optionalGrouping
		defaultGrp grouper
		series     map[uint64]optionalGrouping // map[seriesHash]
		charts     map[string]grouper          // map[chartID]
	}
	chartsCache map[string]*module.Chart // map[chartID]
	dimsCache   map[string]struct{}      // map[dimID]
)

func newCacheEntry(optGrps []optionalGrouping) *cacheEntry {
	return &cacheEntry{
		groups: groupCache{
			user:   optGrps,
			series: make(map[uint64]optionalGrouping),
			charts: make(map[string]grouper),
		},
		charts: make(chartsCache),
		dims:   make(dimsCache),
	}
}

func (c collectCache) get(name string) *cacheEntry        { return c[name] }
func (c collectCache) has(name string) bool               { _, ok := c[name]; return ok }
func (c collectCache) put(name string, entry *cacheEntry) { c[name] = entry }
func (c collectCache) remove(name string)                 { delete(c, name) }

func (ce cacheEntry) hasChart(chartID string) bool                 { _, ok := ce.charts[chartID]; return ok }
func (ce cacheEntry) putChart(chartID string, chart *module.Chart) { ce.charts[chartID] = chart }
func (ce cacheEntry) getChart(chartID string) *module.Chart        { return ce.charts[chartID] }
func (ce cacheEntry) removeChart(chartID string)                   { delete(ce.charts, chartID) }

func (ce cacheEntry) hasDim(dimID string) bool { _, ok := ce.dims[dimID]; return ok }
func (ce cacheEntry) putDim(dimID string)      { ce.dims[dimID] = struct{}{} }
func (ce cacheEntry) removeDim(dimID string)   { delete(ce.dims, dimID) }

func (ce *cacheEntry) getGrouping(pm prometheus.Metric, pms prometheus.Metrics) grouper {
	if ce.groups.defaultGrp != nil {
		return ce.groups.defaultGrp
	}
	if !ce.groups.init {
		ce.groups.init = true
		if len(ce.groups.user) == 0 || !isUserGroupingMatches(pms, ce.groups.user) {
			num := desiredNumOfCharts(len(pms))
			ce.groups.defaultGrp = newGroupingSplitN(defaultAnyGrouping, num)
			return ce.groups.defaultGrp
		}
	}

	hash := pm.Labels.Hash()
	tsGrp, ok := ce.groups.series[hash]
	if !ok {
		tsGrp = findTimeSeriesGrouping(pm.Labels, ce.groups.user)
		ce.groups.series[hash] = tsGrp
	}

	id := tsGrp.grp.chartID(pm)
	grp, ok := ce.groups.charts[id]
	if !ok {
		grp = findChartGrouping(pms, id, tsGrp)
		ce.groups.charts[id] = grp
	}
	return grp
}

func findTimeSeriesGrouping(lbs labels.Labels, userGrps []optionalGrouping) optionalGrouping {
	sr := selector.True()
	for _, item := range userGrps {
		if item.sr.Matches(lbs) {
			return item
		}
		sr = selector.And(sr, selector.Not(item.sr))
	}
	return optionalGrouping{
		sr:  sr,
		grp: defaultAnyGrouping,
	}
}

func findChartGrouping(pms prometheus.Metrics, id string, optGrp optionalGrouping) grouper {
	var numOfSeries int
	for _, pm := range pms {
		if optGrp.sr.Matches(pm.Labels) && optGrp.grp.chartID(pm) == id {
			numOfSeries++
		}
	}
	num := desiredNumOfCharts(numOfSeries)
	return newGroupingSplitN(optGrp.grp, num)
}

func isUserGroupingMatches(pms prometheus.Metrics, userGrps []optionalGrouping) bool {
	for _, pm := range pms {
		for _, grp := range userGrps {
			if grp.sr.Matches(pm.Labels) {
				return true
			}
		}
	}
	return false
}

func desiredNumOfCharts(numOfSeries int) (num uint64) {
	num = uint64(numOfSeries / desiredDim)
	if numOfSeries%desiredDim != 0 {
		num++
	}
	return num
}

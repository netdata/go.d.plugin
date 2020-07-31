package prometheus

import (
	"github.com/netdata/go-orchestrator/module"
)

// TODO: make it easier, very error-prone interface
type (
	metricsCache map[string]*cacheEntry
	cacheEntry   struct {
		split  split
		charts chartsCache
		dims   dimsCache
	}
	chartsCache map[string]*module.Chart
	dimsCache   map[string]struct{}
)

func (c metricsCache) get(name string) *cacheEntry        { return c[name] }
func (c metricsCache) has(name string) bool               { _, ok := c[name]; return ok }
func (c metricsCache) put(name string, entry *cacheEntry) { c[name] = entry }
func (c metricsCache) remove(name string)                 { delete(c, name) }

func (ce cacheEntry) hasChart(chartID string) bool                 { _, ok := ce.charts[chartID]; return ok }
func (ce cacheEntry) putChart(chartID string, chart *module.Chart) { ce.charts[chartID] = chart }
func (ce cacheEntry) getChart(chartID string) *module.Chart        { return ce.charts[chartID] }
func (ce cacheEntry) removeChart(chartID string)                   { delete(ce.charts, chartID) }

func (ce cacheEntry) hasDim(dimID string) bool { _, ok := ce.dims[dimID]; return ok }
func (ce cacheEntry) putDim(dimID string)      { ce.dims[dimID] = struct{}{} }
func (ce cacheEntry) removeDim(dimID string)   { delete(ce.dims, dimID) }

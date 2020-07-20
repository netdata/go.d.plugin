package prometheus

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"

	"github.com/netdata/go-orchestrator/module"
)

// TODO: make it easier, very error-prone interface
type (
	metricsCache map[string]*cacheEntry
	cacheEntry   struct {
		split         func(metric prometheus.Metric) string
		desiredCharts uint64
		charts        chartsCache
		dims          dimsCache
	}
	chartsCache map[string]*module.Chart
	dimsCache   map[string]map[string]struct{}
)

func (c metricsCache) get(key string) *cacheEntry        { return c[key] }
func (c metricsCache) has(key string) bool               { _, ok := c[key]; return ok }
func (c metricsCache) put(key string, entry *cacheEntry) { c[key] = entry }
func (c metricsCache) remove(key string)                 { delete(c, key) }

func (ce cacheEntry) hasChart(key string) bool                 { _, ok := ce.charts[key]; return ok }
func (ce cacheEntry) putChart(key string, chart *module.Chart) { ce.charts[key] = chart }
func (ce cacheEntry) getChart(key string) *module.Chart        { return ce.charts[key] }

func (ce cacheEntry) hasDim(dimKey, chartKey string) bool {
	if _, ok := ce.dims[dimKey]; !ok {
		return false
	}
	_, ok := ce.dims[dimKey][chartKey]
	return ok
}

func (ce cacheEntry) putDim(dimKey, chartKey string) {
	if _, ok := ce.dims[dimKey]; !ok {
		ce.dims[dimKey] = make(map[string]struct{})
	}
	ce.dims[dimKey][chartKey] = struct{}{}
}

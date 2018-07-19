package cooked

import (
	"github.com/l2isbad/go.d.plugin/internal/pkg/charts/raw"
	"github.com/l2isbad/go.d.plugin/internal/pkg/logger"
)

var initPriority = 70000

type baseConfHook interface {
	ModuleName() string // for Chart + CacheGet
	JobName() string    // for CacheGet
	FullName() string   // for Chart
	UpdateEvery() int   // for Chart
}

type ChartsMap map[string]*Chart

type charts struct {
	charts   ChartsMap
	bc       baseConfHook
	priority int
}

func NewCharts(bc baseConfHook) *charts {
	return &charts{
		charts:   make(ChartsMap),
		bc:       bc,
		priority: initPriority}
}

// AddOne adds/re-adds one raw Chart.
func (c *charts) AddOne(r *raw.Chart) error {
	newChart, err := newChart(r, c.bc, c.priority)
	if err != nil {
		logger.CacheGet(c.bc).Errorf("invalid Chart '%s' (%s)", newChart.id, err)
		return err
	}
	// re-add
	if v, ok := c.charts[newChart.id]; ok {
		newChart.priority = v.priority
		return nil
	}
	// add
	c.priority++
	c.charts[newChart.id] = newChart
	return nil
}

// AddMany adds all charts from (raw.Charts) Order if they are in Definitions.
func (c *charts) AddMany(r *raw.Charts) int {
	var added int

	for _, chartID := range r.Order {
		rawChart := r.GetChartByID(chartID)
		if rawChart == nil {
			logger.CacheGet(c.bc).Warningf("'%s' is not in Definitions, skipping it", chartID)
			continue
		}
		if err := c.AddOne(rawChart); err != nil {
			continue
		}
		added++
	}
	return added
}

// GetCharts returns charts.
func (c *charts) GetCharts() ChartsMap {
	return c.charts
}

// GetCharts returns chart by id.
func (c *charts) GetChartByID(id string) *Chart {
	return c.charts[id]
}

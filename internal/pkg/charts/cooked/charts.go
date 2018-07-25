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

type charts struct {
	items    map[string]*Chart
	bc       baseConfHook
	priority int
}

func NewCharts(bc baseConfHook) *charts {
	return &charts{
		items:    make(map[string]*Chart),
		bc:       bc,
		priority: initPriority}
}

// AddOne adds/re-adds one raw Chart.
func (c *charts) AddOne(r *raw.Chart) error {
	if err := check(r); err != nil {
		logger.CacheGet(c.bc).Errorf("invalid Chart '%s' (%s)", r.ID, err)
		return err
	}

	chart := newChart(r, c.bc, c.priority)
	// re-add
	if v, ok := c.items[chart.id]; ok {
		chart.priority = v.priority
		return nil
	}
	// add
	c.priority++
	c.items[chart.id] = chart
	return nil
}

// AddMany adds all items from (raw.Charts) Order if they are in Definitions.
func (c *charts) AddMany(r *raw.Charts) int {
	var added int

	for _, id := range r.Order {
		chart, ok := r.LookupChartByID(id)
		if !ok {
			logger.CacheGet(c.bc).Warningf("'%s' is not in Definitions, skipping it", id)
			continue
		}
		if err := c.AddOne(chart); err != nil {
			continue
		}
		added++
	}
	return added
}

// ListNames returns list of chart names.
func (c *charts) ListNames() []string {
	var rv []string
	for k := range c.items {
		rv = append(rv, k)
	}
	return rv
}

// GetCharts returns chart by id.
func (c *charts) GetChartByID(id string) *Chart {
	return c.items[id]
}

// LookupChartsByID looks up a chart by id.
func (c *charts) LookupChartByID(id string) (*Chart, bool) {
	v, ok := c.items[id]
	return v, ok
}

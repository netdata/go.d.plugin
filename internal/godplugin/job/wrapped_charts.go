package job

import (
	"github.com/l2isbad/go.d.plugin/internal/modules"
	"github.com/l2isbad/go.d.plugin/internal/pkg/charts"
)

var initPriority = 70000

type baseConfHook interface {
	ModuleName() string // for Chart + CacheGet
	JobName() string    // for CacheGet
	FullName() string   // for Chart
	UpdateEvery() int   // for Chart
}

func NewWrappedCharts(h baseConfHook) *WrappedCharts {
	return &WrappedCharts{
		items:    make(map[string]*wrappedChart),
		hook:     h,
		priority: initPriority,
	}
}

type WrappedCharts struct {
	items    map[string]*wrappedChart
	hook     baseConfHook
	priority int
}

func (w *WrappedCharts) AddChart(charts ...charts.Chart) {
	for idx := range charts {
	//	if !check(charts[idx]) {
	//		continue
	//	}
		chart := newWrappedChart(charts[idx], w.hook, w.priority)
		v, ok := w.items[chart.item.ID]

		if ok {
			chart.priority = v.priority
		} else {
			w.priority++
		}

		w.items[chart.item.ID] = chart
	}
}

func (w WrappedCharts) GetChart(id string) modules.Chart {
	return w.items[id]
}

func (w WrappedCharts) LookupChart(id string) (modules.Chart, bool) {
	v, ok := w.items[id]
	return v, ok
}

//func check(c charts.Chart) bool {
//	if !c.IsValid() {
//		return false
//	}
//	for idx := range c.Dimensions {
//		if !c.Dimensions[idx].IsValid() {
//			return false
//		}
//	}
//	for idx := range c.Variables {
//		if !c.Variables[idx].IsValid() {
//			return false
//		}
//	}
//	return true
//}

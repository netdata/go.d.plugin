package job

import (
	"github.com/l2isbad/go.d.plugin/internal/modules"
	"github.com/l2isbad/go.d.plugin/internal/pkg/charts"
)

type baseConfHook interface {
	ModuleName() string // for Chart + CacheGet
	JobName() string    // for CacheGet
	FullName() string   // for Chart
	UpdateEvery() int   // for Chart
}

func NewCharts(h baseConfHook) *Charts {
	return &Charts{
		items: make(map[string]*chart),
		hook:  h,
		prio:  70000,
	}
}

type Charts struct {
	items map[string]*chart
	hook  baseConfHook
	prio  int
}

func (w *Charts) AddChart(charts ...charts.Chart) {
	for idx := range charts {
		//	if !check(charts[idx]) {
		//		continue
		//	}
		chart := newChart(charts[idx], w.hook, w.prio)
		v, ok := w.items[chart.item.ID]

		if ok {
			chart.prio = v.prio
		} else {
			w.prio++
		}

		w.items[chart.item.ID] = chart
	}
}

func (w Charts) GetChart(id string) modules.Chart {
	return w.items[id]
}

func (w Charts) LookupChart(id string) (modules.Chart, bool) {
	v, ok := w.items[id]
	return v, ok
}

//func check(c charts.Chart) bool {
//	if !c.IsValid() {
//		return false
//	}
//	for idx := range c.Dims {
//		if !c.Dims[idx].IsValid() {
//			return false
//		}
//	}
//	for idx := range c.Vars {
//		if !c.Vars[idx].IsValid() {
//			return false
//		}
//	}
//	return true
//}

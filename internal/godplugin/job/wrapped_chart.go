package job

import (
	"fmt"
	"github.com/l2isbad/go.d.plugin/internal/pkg/charts"
)

func newWrappedChart(chart charts.Chart, hook baseConfHook, prio int) *wrappedChart {
	p := &wrappedChart{
		item:     chart,
		hook:     hook,
		priority: prio,
		flags:    chartFlags{push: true},
	}
	for _, v := range chart.Dims {
		p.dimensions = append(p.dimensions, newWrappedDim(v))
	}
	for _, v := range chart.Vars {
		p.variables = append(p.variables, &wrappedVar{v})
	}
	return p
}

type wrappedChart struct {
	item          charts.Chart
	dimensions    []*wrappedDim
	variables     []*wrappedVar
	hook          baseConfHook
	failedUpdates int
	priority      int
	flags         chartFlags
}

func (w wrappedChart) begin(sinceLast int) string {
	return fmt.Sprintf(formatChartBEGIN,
		w.hook.FullName(),
		w.item.ID,
		sinceLast,
	)
}

func (w wrappedChart) create() string {
	var dims, vars string
	chart := fmt.Sprintf(formatChartCREATE,
		w.hook.FullName(),
		w.item.ID,
		w.item.OverrideID,
		w.item.Title,
		w.item.Units,
		w.item.Family,
		w.item.Context,
		w.item.Type,
		w.priority,
		w.hook.UpdateEvery(),
		w.hook.ModuleName(),
	)

	for idx := range w.dimensions {
		dims += w.dimensions[idx].create()
	}

	for _, v := range w.variables {
		if v.item.Value != 0 {
			vars += v.set(v.item.Value)
		}
	}

	w.flags.setPush(false)
	w.flags.setCreated(true)

	return chart + dims + vars + "\n"
}

func (w *wrappedChart) obsolete() {
	w.flags.setObsoleted(true)
	if !w.flags.created {
		return
	}
	safePrint(fmt.Sprintf(formatChartOBSOLETE,
		w.hook.FullName(),
		w.item.ID,
		w.item.OverrideID,
		w.item.Title,
		w.item.Units,
		w.item.Family,
		w.item.Context,
		w.item.Type,
		w.priority,
		w.hook.UpdateEvery(),
		w.hook.ModuleName()),
	)
}

func (w *wrappedChart) refresh() {
	w.flags.setPush(true)
	if w.flags.obsoleted {
		w.failedUpdates = 0
		w.flags.setCreated(false)
		w.flags.setObsoleted(false)
	}
}

func (w wrappedChart) canBeUpdated(data map[string]int64) bool {
	for idx := range w.dimensions {
		if _, ok := data[w.dimensions[idx].item.ID]; ok {
			return true
		}
	}
	return false
}

func (w *wrappedChart) update(data map[string]int64, interval int) bool {
	var (
		dims    string
		vars    string
		success bool
	)

	for _, d := range w.dimensions {
		if value, ok := d.get(data); ok {
			dims += d.set(value)
			success = true
		} else {
			dims += d.empty()
		}
		if d.flags.push {
			w.flags.setPush(true)
			d.flags.push = false
		}
	}

	for _, v := range w.variables {
		if value, ok := data[v.item.ID]; ok {
			vars += v.set(value)
		}
	}

	if !success {
		w.failedUpdates++
		w.flags.setUpdated(false)
		return false
	}

	if !w.flags.updated {
		interval = 0
	}

	if w.flags.push {
		safePrint(w.create())
	}

	safePrint(w.begin(interval), dims, vars, "END\n\n")
	w.flags.setUpdated(true)
	w.failedUpdates = 0

	return true
}

func (w *wrappedChart) AddDim(dim charts.Dim) {
	if w.indexDim(dim.ID) == -1 {
		w.dimensions = append(w.dimensions, newWrappedDim(dim))
	}
}

func (w *wrappedChart) AddVar(v charts.Var) {
	if w.indexDim(v.ID) == -1 {
		w.variables = append(w.variables, &wrappedVar{v})
	}
}

func (w wrappedChart) indexDim(id string) int {
	for idx := range w.dimensions {
		if w.dimensions[idx].item.ID == id {
			return idx
		}
	}
	return -1
}

func (w wrappedChart) indexVar(id string) int {
	for idx := range w.variables {
		if w.variables[idx].item.ID == id {
			return idx
		}
	}
	return -1
}

// ------------------------------------------------------------------------------------
type chartFlags struct {
	push      bool
	created   bool
	updated   bool
	obsoleted bool
}

func (f *chartFlags) setPush(b bool) {
	f.push = b
}

func (f *chartFlags) setCreated(b bool) {
	f.created = b
}

func (f *chartFlags) setUpdated(b bool) {
	f.updated = b
}

func (f *chartFlags) setObsoleted(b bool) {
	f.obsoleted = b
}

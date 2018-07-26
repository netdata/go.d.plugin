package job

import (
	"fmt"
	"github.com/l2isbad/go.d.plugin/internal/pkg/charts"
)

func newChart(ch charts.Chart, hook baseConfHook, prio int) *chart {
	p := &chart{
		item: ch,
		hook: hook,
		prio: prio,
		push: true,
	}
	for _, v := range ch.Dims {
		p.dimensions = append(p.dimensions, newDim(v))
	}
	for _, v := range ch.Vars {
		p.variables = append(p.variables, &variable{v})
	}
	return p
}

type chart struct {
	item          charts.Chart
	dimensions    []*dimension
	variables     []*variable

	hook          baseConfHook
	failedUpdates int
	prio          int

	push      bool
	created   bool
	updated   bool
	obsoleted bool
}

func (w *chart) AddDim(dim charts.Dim) {
	if w.indexDim(dim.ID) == -1 {
		w.dimensions = append(w.dimensions, newDim(dim))
	}
}

func (w *chart) AddVar(v charts.Var) {
	if w.indexDim(v.ID) == -1 {
		w.variables = append(w.variables, &variable{v})
	}
}

func (w chart) begin(sinceLast int) string {
	return fmt.Sprintf(formatChartBEGIN,
		w.hook.FullName(),
		w.item.ID,
		sinceLast,
	)
}

func (w chart) create() string {
	var dims, vars string
	chart := fmt.Sprintf(formatChartCREATE,
		w.hook.FullName(),
		w.item.ID,
		w.item.OverID,
		w.item.Title,
		w.item.Units,
		w.item.Fam,
		w.item.Ctx,
		w.item.Type,
		w.prio,
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

	w.push = false
	w.created = true

	return chart + dims + vars + "\n"
}

func (w *chart) obsolete() {
	w.obsoleted = true
	if !w.created {
		return
	}
	safePrint(fmt.Sprintf(formatChartOBSOLETE,
		w.hook.FullName(),
		w.item.ID,
		w.item.OverID,
		w.item.Title,
		w.item.Units,
		w.item.Fam,
		w.item.Ctx,
		w.item.Type,
		w.prio,
		w.hook.UpdateEvery(),
		w.hook.ModuleName()),
	)
}

func (w *chart) refresh() {
	w.push = true
	if w.obsoleted {
		w.failedUpdates = 0
		w.created = false
		w.obsoleted = false
	}
}

func (w chart) canBeUpdated(data map[string]int64) bool {
	for idx := range w.dimensions {
		if _, ok := data[w.dimensions[idx].item.ID]; ok {
			return true
		}
	}
	return false
}

func (w *chart) update(data map[string]int64, interval int) bool {
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
		if d.push {
			w.push = true
			d.push = false
		}
	}

	for _, v := range w.variables {
		if value, ok := data[v.item.ID]; ok {
			vars += v.set(value)
		}
	}

	if !success {
		w.failedUpdates++
		w.updated = false
		return false
	}

	if !w.created {
		interval = 0
	}

	if w.push {
		safePrint(w.create())
	}

	safePrint(w.begin(interval), dims, vars, "END\n\n")
	w.updated = true
	w.failedUpdates = 0

	return true
}

func (w chart) indexDim(id string) int {
	for idx := range w.dimensions {
		if w.dimensions[idx].item.ID == id {
			return idx
		}
	}
	return -1
}

func (w chart) indexVar(id string) int {
	for idx := range w.variables {
		if w.variables[idx].item.ID == id {
			return idx
		}
	}
	return -1
}

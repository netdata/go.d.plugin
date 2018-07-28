package job

import (
	"github.com/l2isbad/go.d.plugin/internal/pkg/charts"
)

type baseConfHook interface {
	FullName() string
	ModuleName() string
	UpdateEvery() int
}

type (
	observer struct {
		charts *charts.Charts

		items map[string]*chart
		hook  baseConfHook
		prio  int
	}

	chart struct {
		item *charts.Chart

		hook    baseConfHook
		prio    int
		retries int

		push      bool
		created   bool
		updated   bool
		obsoleted bool
	}
)

func (c *chart) refresh() {
	c.push = true
	if c.obsoleted {
		c.retries = 0
		c.updated = false
		c.obsoleted = false
	}
}

// FIXME: DUPLICATE CHARTS, DIMS
func NewObserver(ch *charts.Charts, hook baseConfHook) *observer {
	o := &observer{
		charts: ch,
		items: make(map[string]*chart),
		prio:  70000,
		hook:  hook,
	}
	return o
}

func (o *observer) init() {
	for _, v := range *o.charts {
		o.Add(v.ID)
	}
}

func (o observer) Update(id string) {
	o.items[id].refresh()
}

func (o *observer) Obsolete(id string) {
	o.items[id].obsolete()
}

func (o *observer) Delete(id string) {
	o.items[id].obsolete()
	delete(o.items, id)
}

func (o *observer) Add(id string) {
	ch := o.charts.GetChart(id)
	ch.Register(o)

	chart := &chart{
		item: ch,
		hook: o.hook,
		prio: o.prio,
		push: true,
	}

	o.prio++

	o.items[ch.ID] = chart
}

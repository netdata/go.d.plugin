package job

import (
	"fmt"

	"github.com/l2isbad/go.d.plugin/internal/pkg/charts"
)

type (
	observer struct {
		charts *charts.Charts

		items map[string]*chart
		hook  *Config
		prio  int
	}

	chart struct {
		item *charts.Chart

		hook    *Config
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

func newObserver(hook *Config) *observer {
	o := &observer{
		items: make(map[string]*chart),
		prio:  70000,
		hook:  hook,
	}
	return o
}

func (o *observer) Set(c *charts.Charts) {
	o.charts = c
}

func (o observer) Update(id string) {
	o.items[id].refresh()
}

func (o *observer) Obsolete(id string) {
	if !o.items[id].obsoleted {
		o.items[id].obsolete()
	}
}

func (o *observer) Delete(id string) {
	if v, ok := o.items[id]; ok {
		v.obsolete()
		delete(o.items, id)
	}
}

func (o *observer) add(ch *charts.Chart) {
	ch.Register(o)

	if ch.Ctx == "" {
		ch.Ctx = fmt.Sprintf("%s.%s", o.hook.moduleName, ch.ID)
	}

	chart := &chart{
		item: ch,
		hook: o.hook,
		prio: o.prio,
		push: true,
	}

	o.prio++

	o.items[ch.ID] = chart
}

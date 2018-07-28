package job

import (
	"fmt"

	"github.com/l2isbad/go.d.plugin/internal/pkg/charts"
)

func canBeUpdated(c chart, data map[string]int64) bool {
	for _, v := range c.item.Dims {
		if _, ok := data[v.ID]; ok {
			return true
		}
	}
	return false
}

func (c chart) begin(sinceLast int, dims, vars string) {
	begin := fmt.Sprintf(formatChartBEGIN,
		c.hook.FullName(),
		c.item.ID,
		sinceLast,
	)

	safePrint(begin, dims, vars, "END\n\n")
}

func (c *chart) create() {
	var dims, vars string
	chart := fmt.Sprintf(formatChartCREATE,
		c.hook.FullName(),
		c.item.ID,
		c.item.OverID,
		c.item.Title,
		c.item.Units,
		c.item.Fam,
		c.item.Ctx,
		c.item.Type,
		c.prio,
		c.hook.UpdateEvery(),
		c.hook.ModuleName(),
	)

	for idx := range c.item.Dims {
		dims += dimCreate(c.item.Dims[idx])
	}

	for _, v := range c.item.Vars {
		if v.Value != 0 {
			vars += varSet(v.ID, v.Value)
		}
	}

	c.push = false
	c.created = true

	safePrint(chart + dims + vars + "\n")
}

func (c *chart) obsolete() {
	c.obsoleted = true
	if !c.created {
		return
	}
	c.created = false
	safePrint(fmt.Sprintf(formatChartOBSOLETE,
		c.hook.FullName(),
		c.item.ID,
		c.item.OverID,
		c.item.Title,
		c.item.Units,
		c.item.Fam,
		c.item.Ctx,
		c.item.Type,
		c.prio,
		c.hook.UpdateEvery(),
		c.hook.ModuleName()),
	)
}

func (c *chart) update(data map[string]int64, interval int) bool {
	var (
		dims    string
		vars    string
		success bool
	)

	for _, d := range c.item.Dims {
		val, ok := data[d.ID]
		if !ok {
			dims += dimEmpty(d.ID)
		} else {
			dims += dimSet(d.ID, val)
			success = true
		}
	}

	for _, v := range c.item.Vars {
		if val, ok := data[v.ID]; ok {
			vars += varSet(v.ID, val)
		}
	}

	if !success {
		c.retries++
		c.updated = false
		return false
	}

	if !c.updated {
		interval = 0
	}

	if c.push {
		c.create()
	}

	c.begin(interval, dims, vars)
	c.updated = true
	c.retries = 0

	return true
}

func dimCreate(d *charts.Dim) string {
	return fmt.Sprintf(formatDimCREATE,
		d.ID,
		d.Name,
		d.Algo,
		d.Mul,
		d.Div,
		d.Hidden,
	)
}

func dimSet(id string, val int64) string {
	return fmt.Sprintf(formatDimSET, id, val)
}

func dimEmpty(id string) string {
	return fmt.Sprintf(formatDimEmptySET, id)
}

func varSet(id string, val int64) string {
	return fmt.Sprintf(formatVarSET, id, val)
}

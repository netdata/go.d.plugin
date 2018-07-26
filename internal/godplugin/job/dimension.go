package job

import (
	"fmt"

	"github.com/l2isbad/go.d.plugin/internal/pkg/charts"
)

func newDim(d charts.Dim) *dimension {
	return &dimension{
		item: d,
		push: true,
		maxRetries: 5,
	}
}

type dimension struct {
	item       charts.Dim
	push       bool
	retries    int
	maxRetries int
}

func (d dimension) alive() bool {
	return d.retries < d.maxRetries
}

func (d *dimension) get(m map[string]int64) (v int64, ok bool) {
	v, ok = m[d.item.ID]
	if ok || !d.alive() {
		return
	}
	d.retries++
	if !d.alive() {
		d.item.Hidden = charts.Hidden
		d.push = true
	}
	return
}

func (d dimension) create() string {
	return fmt.Sprintf(formatDimCREATE,
		d.item.ID,
		d.item.Name,
		d.item.Algo,
		d.item.Mul,
		d.item.Div,
		d.item.Hidden,
	)
}

func (d *dimension) set(value int64) string {
	if !d.alive() {
		d.item.Hidden = charts.NotHidden
		d.push = true
	}
	d.retries = 0
	return fmt.Sprintf(formatDimSET,
		d.item.ID,
		value,
	)
}

func (d dimension) empty() string {
	return fmt.Sprintf(formatDimEmptySET, d.item.ID)
}

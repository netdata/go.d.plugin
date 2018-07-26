package job

import (
	"fmt"
	"github.com/l2isbad/go.d.plugin/internal/pkg/charts"
)

var retriesMax = 5

func newWrappedDim(dim charts.Dimension) *wrappedDim {
	return &wrappedDim{
		item:  dim,
		flags: dimFlags{retriesMax: retriesMax},
	}
}

type wrappedDim struct {
	item  charts.Dimension
	flags dimFlags
}

func (w *wrappedDim) setHidden(b bool) {
	if b {
		w.item.Hidden = charts.Hidden
	} else {
		w.item.Hidden = charts.NotHidden
	}
}

func (w *wrappedDim) get(m map[string]int64) (v int64, ok bool) {
	v, ok = m[w.item.ID]
	if ok || !w.flags.alive() {
		return
	}
	w.flags.retries++
	if !w.flags.alive() {
		w.setHidden(true)
		w.flags.setPush(true)
	}
	return
}

func (w wrappedDim) create() string {
	return fmt.Sprintf(formatDimCREATE,
		w.item.ID,
		w.item.Name,
		w.item.Algorithm,
		w.item.Multiplier,
		w.item.Divisor,
		w.item.Hidden,
	)
}

func (w *wrappedDim) set(value int64) string {
	if !w.flags.alive() {
		w.setHidden(false)
		w.flags.setPush(true)
	}
	w.flags.retries = 0
	return fmt.Sprintf(formatDimSET,
		w.item.ID,
		value,
	)
}

func (w wrappedDim) empty() string {
	return fmt.Sprintf(formatDimEmptySET, w.item.ID)
}

// Dimension Flag ------------------------------------------------------------------------------------------------------
type dimFlags struct {
	push       bool
	retries    int
	retriesMax int
}

func (f *dimFlags) setPush(b bool) {
	f.push = b
}

func (f *dimFlags) resetRetries() {
	f.retries = 0
}

func (f dimFlags) alive() bool {
	return f.retries < f.retriesMax
}

package example

import (
	"math/rand"

	"github.com/netdata/go.d.plugin/modules"
)

type Example struct {
	modules.Base

	data map[string]int64
}

func New() modules.Creator {
	return modules.Creator{
		Create: func() modules.Module {
			return &Example{data: make(map[string]int64)}
		},
	}
}

func (e *Example) Check() bool {
	return true
}

func (Example) GetCharts() *Charts {
	return charts.Copy()
}

func (e *Example) GatherMetrics() map[string]int64 {
	e.data["random0"] = rand.Int63n(100)
	e.data["random1"] = rand.Int63n(100)

	return e.data
}

func init() {
	modules.Register("example", New())
}

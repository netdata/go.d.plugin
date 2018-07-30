package example

import (
	"math/rand"

	"github.com/l2isbad/go.d.plugin/internal/modules"
	"github.com/l2isbad/go.d.plugin/internal/pkg/charts"
)

type Example struct {
	modules.NoConfiger

	data map[string]int64
}

func (e *Example) Check() bool {
	return true
}

func (Example) GetCharts() *charts.Charts {
	return charts.NewCharts(uCharts...)
}

func (e *Example) GetData() map[string]int64 {
	e.data["random0"] = rand.Int63n(100)
	e.data["random1"] = rand.Int63n(100)

	return e.data
}

func init() {
	modules.SetDefault().SetDisabledByDefault()

	f := func() modules.Module {
		return &Example{
			data: make(map[string]int64),
		}
	}
	modules.Add(f)
}

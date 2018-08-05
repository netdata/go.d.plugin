package example

import (
	"math/rand"

	"github.com/l2isbad/go.d.plugin/internal/modules"
	"github.com/l2isbad/go.d.plugin/internal/pkg/charts"
)

func init() {
	modules.Register("example", modules.Creator{
		DisabledByDefault: true,
		NoConfig:          true,
		Create: func() modules.Module {
			return &Example{
				data: make(map[string]int64),
			}
		},
	})
}

// Example module
type Example struct {
	modules.ModuleBase

	data map[string]int64
}

// Init Init
func (Example) Init() error { return nil }

// Check Check
func (e *Example) Check() bool {
	return true
}

// GetCharts GetCharts
func (Example) GetCharts() *charts.Charts {
	return charts.NewCharts(uCharts...)
}

// GetData GetData
func (e *Example) GetData() map[string]int64 {
	e.data["random0"] = rand.Int63n(100)
	e.data["random1"] = rand.Int63n(100)

	return e.data
}

package example

import (
	"fmt"
	"math/rand"

	"github.com/l2isbad/go.d.plugin/charts/raw"
	"github.com/l2isbad/go.d.plugin/modules"
)

type (
	Charts      = raw.Charts
	Order       = raw.Order
	Definitions = raw.Definitions
	Chart       = raw.Chart
	Options     = raw.Options
	Dimensions  = raw.Dimensions
	Dimension   = raw.Dimension
	Variables   = raw.Variables
	Variable    = raw.Variable
)

var uCharts = Charts{
	Order: Order{"chart1", "chart2"},
	Definitions: Definitions{
		Chart{
			ID:      "chart1",
			Options: Options{"Random Data", "random", "random"},
			Dimensions: Dimensions{
				Dimension{"random0"},
			},
		},
		Chart{
			ID:      "chart2",
			Options: Options{"Random Data", "random", "random"},
			Dimensions: Dimensions{
				Dimension{"random1"},
			},
		},
	},
}

type Example struct {
	modules.Charts
	modules.NoConfiger

	data map[string]int64
}

func (e *Example) Check() bool {
	e.AddMany(&uCharts)
	return true
}

func (e *Example) GetData() *map[string]int64 {
	for i := 0; i < 2; i++ {
		e.data[fmt.Sprintf("random%d", i)] = rand.Int63n(100)
	}
	return &e.data
}

func init() {
	f := func() modules.Module {
		return &Example{data: make(map[string]int64)}
	}
	modules.Add(f)
}

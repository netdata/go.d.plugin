package example

import (
	"math/rand"

	"github.com/l2isbad/go.d.plugin/internal/modules"
	"github.com/l2isbad/go.d.plugin/internal/pkg/charts/raw"
)

type (
	Charts      = raw.Charts
	Order       = raw.Order
	Definitions = raw.Definitions
	Chart       = raw.Chart
	Options     = raw.Options
	Dimensions  = raw.Dimensions
	Dimension   = raw.Dimension
)

var charts = Charts{
	Order: Order{"chart1", "chart2"},
	Definitions: Definitions{
		&Chart{
			ID:      "chart1",
			Options: Options{"Random Data", "random", "random"},
			Dimensions: Dimensions{
				Dimension{"random0"},
			},
		},
		&Chart{
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
	modules.Logger
	modules.NoConfiger

	data map[string]int64
}

func (e *Example) Check() bool {
	e.AddMany(&charts)
	return true
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

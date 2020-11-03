package example

import (
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

func (e *Example) collect() (map[string]int64, error) {
	collected := make(map[string]int64)

	var numOfDims int
	for _, chart := range *e.Charts() {
		if chart.Opts.Hidden {
			numOfDims = e.Config.HiddenCharts.Dims
		} else {
			numOfDims = e.Config.Charts.Dims
		}

		for i := 0; i < numOfDims; i++ {
			name := fmt.Sprintf("random%d", i)
			id := fmt.Sprintf("%s_%s", chart.ID, name)

			if !e.collectedDims[id] {
				e.collectedDims[id] = true

				dim := &module.Dim{ID: id, Name: name}
				if err := chart.AddDim(dim); err != nil {
					e.Warning(err)
				}
				chart.MarkNotCreated()
			}

			collected[id] = e.randInt()
		}
	}
	return collected, nil
}

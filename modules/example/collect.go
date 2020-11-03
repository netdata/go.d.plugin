package example

import (
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

func (e *Example) collect() (map[string]int64, error) {
	collected := make(map[string]int64)

	for _, chart := range *e.Charts() {
		for i := 0; i < e.NumDims; i++ {
			name := fmt.Sprintf("random%d", i)
			id := fmt.Sprintf("%s_%s", chart.ID, name)

			if !e.collectedDims[id] {
				e.collectedDims[id] = true

				dim := &module.Dim{ID: id, Name: name}
				if err := chart.AddDim(dim); err != nil {
					e.Warning(err)
				}
			}

			collected[id] = e.randInt()
		}
	}
	return collected, nil
}

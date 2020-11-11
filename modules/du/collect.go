package du

import (
	"errors"

	"github.com/netdata/go.d.plugin/agent/module"
)

func (du *Du) collect() (map[string]int64, error) {
	collected := make(map[string]int64)

	if len(*du.Charts()) <= 0 {
		return collected, errors.New("Chart not found")
	}

	chart := (*du.Charts())[0]
	for _, path := range du.Config.Paths {
		if !du.collectedDims[path] {
			du.collectedDims[path] = true

			dim := &module.Dim{ID: path, Name: path}
			if err := chart.AddDim(dim); err != nil {
				du.Warning(err)
			}
			chart.MarkNotCreated()
		}

		size, err := fileSize(path)
		if err != nil {
			// Return file size as -1 when error happens
			size = -1
		}
		collected[path] = size
	}
	return collected, nil
}

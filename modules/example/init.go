package example

import (
	"fmt"
	"github.com/netdata/go.d.plugin/agent/module"
)

func (e Example) validateConfig() error {
	if e.NumCharts <= 0 {
		return fmt.Errorf("'num_of_charts' must be > 0 (current: %d)", e.NumCharts)
	}
	if e.NumDims <= 0 {
		return fmt.Errorf("'num_of_dimensions' must be > 0 (current: %d)", e.NumDims)
	}
	return nil
}

func (e Example) initCharts() (*module.Charts, error) {
	charts := &module.Charts{}

	for i := 0; i < e.NumCharts; i++ {
		chart := chartTemplate.Copy()
		chart.ID = fmt.Sprintf(chart.ID, i)

		if err := charts.Add(chart); err != nil {
			return nil, err
		}
	}

	return charts, nil
}

package example

import (
	"errors"
	"github.com/netdata/go.d.plugin/agent/module"
)

func (e Example) validateConfig() error {
	if e.Config.Charts.Num <= 0 && e.Config.HiddenCharts.Num <= 0 {
		return errors.New("'charts->num' or `hidden_charts->num` must be > 0")
	}
	if e.Config.Charts.Num > 0 && e.Config.Charts.Dims <= 0 {
		return errors.New("'charts->dimensions' must be > 0")
	}
	if e.Config.HiddenCharts.Num > 0 && e.Config.HiddenCharts.Dims <= 0 {
		return errors.New("'hidden_charts->dimensions' must be > 0")
	}
	return nil
}

func (e Example) initCharts() (*module.Charts, error) {
	charts := &module.Charts{}

	for i := 0; i < e.Config.Charts.Num; i++ {
		chart := newChart(i, module.ChartType(e.Config.Charts.Type))

		if err := charts.Add(chart); err != nil {
			return nil, err
		}
	}

	for i := 0; i < e.Config.HiddenCharts.Num; i++ {
		chart := newHiddenChart(i, module.ChartType(e.Config.HiddenCharts.Type))

		if err := charts.Add(chart); err != nil {
			return nil, err
		}
	}

	return charts, nil
}

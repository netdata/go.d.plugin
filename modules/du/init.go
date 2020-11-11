package du

import (
	"errors"

	"github.com/netdata/go.d.plugin/agent/module"
)

func (du Du) validateConfig() error {
	if len(du.Config.Paths) <= 0 {
		return errors.New("'paths' is required")
	}
	return nil
}

func (du Du) initCharts() (*module.Charts, error) {
	charts := &module.Charts{}
	chart := chartTemplate.Copy()

	if err := charts.Add(chart); err != nil {
		return nil, err
	}

	return charts, nil
}

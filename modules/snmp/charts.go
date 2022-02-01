package snmp

import (
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

type Dims = module.Dims

var snmp_chart_template = module.Chart{
	ID:       "snmp_%d",
	Title:    "%s",
	Units:    "kilobits/s",
	Type:     "area",
	Priority: 1,
	Dims: Dims{
		{
			Name: "in",
			ID:   "1.3.6.1.2.1.2.2.1.10.2",
			Algo: module.Incremental,
			Mul:  8,
			Div:  1024,
		},
		{
			Name: "out",
			ID:   "1.3.6.1.2.1.2.2.1.16.2",
			Algo: module.Incremental,
			Mul:  -8,
			Div:  1024,
		},
	},
}

func newChart(settings []ChartsConfig) *module.Charts {
	charts := &module.Charts{}
	for i, s := range settings {
		c := snmp_chart_template.Copy()
		c.ID = fmt.Sprintf(c.ID, i)
		c.Title = fmt.Sprintf(c.Title, s.Title)
		charts.Add(c)
	}
	return charts
}

func (s SNMP) validateConfig() error {
	//TODO:
	return nil
}

package snmp

import (
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

var snmp_chart_template = module.Chart{
	ID:       "snmp_%d",
	Title:    "%s",
	Units:    "kilobits/s",
	Type:     "area",
	Priority: 1,
	Fam:      "ports",
}

var default_dims = module.Dims{
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
}

func newChart(settings []ChartsConfig) *module.Charts {
	charts := &module.Charts{}
	for i, s := range settings {
		c := snmp_chart_template.Copy()
		c.ID = fmt.Sprintf(c.ID, i)
		c.Title = fmt.Sprintf(c.Title, s.Title)
		if s.Family != "" {
			c.Fam = s.Family
		}

		if s.Units != "" {
			c.Units = s.Units
		}

		if s.Type != "" {
			c.Type = module.ChartType(s.Type)
		}

		c.Priority = s.Priority
		for _, d := range s.Dimensions {
			dim := &module.Dim{
				Name: d.Name,
				ID:   d.OID,
				Algo: module.DimAlgo(d.Algorithm),
				Mul:  d.Multiplier,
				Div:  d.Divisor,
			}
			c.AddDim(dim)
		}

		//Add default ones if no dimensions defined
		if len(c.Dims) == 0 {
			c.AddDim(default_dims[0])
			c.AddDim(default_dims[1])
		}
		charts.Add(c)
	}
	return charts
}

func (s SNMP) validateConfig() error {
	//TODO:
	return nil
}

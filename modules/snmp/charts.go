package snmp

import (
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

var defaultSNMPchart = module.Chart{
	ID:       "snmp_%d",
	Title:    "%s",
	Units:    "kilobits/s",
	Type:     "area",
	Priority: 1,
	Fam:      "ports",
}

var defaultDims = module.Dims{
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

func newChart(id, multiplier int, s ChartsConfig) *module.Chart {
	c := defaultSNMPchart.Copy()
	c.ID = fmt.Sprintf(c.ID, id)
	c.Title = fmt.Sprintf(c.Title, s.Title)

	if multiplier != 0 {
		c.ID = fmt.Sprintf("%s_%d", c.ID, multiplier)
		c.Title = fmt.Sprintf("%s %d", c.Title, multiplier)
	}

	if s.Family != nil {
		c.Fam = *s.Family
	}

	if s.Units != nil {
		c.Units = *s.Units
	}

	if s.Type != nil {
		c.Type = module.ChartType(*s.Type)
	}

	c.Priority = s.Priority
	for _, d := range s.Dimensions {
		oid := d.OID
		if multiplier != 0 {
			oid = fmt.Sprintf("%s.%d", oid, multiplier)
		}
		dim := &module.Dim{
			Name: d.Name,
			ID:   oid,
		}
		if d.Algorithm != nil {
			dim.Algo = module.DimAlgo(*d.Algorithm)
		}
		if d.Multiplier != nil {
			dim.Mul = *d.Multiplier
		}
		if d.Divisor != nil {
			dim.Div = *d.Divisor
		}
		_ = c.AddDim(dim)
	}

	return c
}

func allCharts(chartIn []ChartsConfig) *module.Charts {
	charts := &module.Charts{}
	for i, s := range chartIn {
		if s.MultiplyRange != nil {
			for j := s.MultiplyRange[0]; j <= s.MultiplyRange[1]; j++ {
				chart := newChart(i, j, s)
				_ = charts.Add(chart)
			}
		} else {
			chart := newChart(i, 0, s)
			_ = charts.Add(chart)
		}
	}
	return charts
}

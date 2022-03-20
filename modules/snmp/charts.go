package snmp

import (
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

//TODO: probably not needed
var defaultSNMPchart = module.Chart{
	Title:    "%s",
	Units:    "kilobits/s",
	Type:     "area",
	Priority: 1,
	Fam:      "ports",
}

func newCharts(configs []ChartsConfig) (*module.Charts, error) {
	charts := &module.Charts{}
	for _, cfg := range configs {
		if cfg.MultiplyRange != nil {
			for j := cfg.MultiplyRange[0]; j <= cfg.MultiplyRange[1]; j++ {
				chart, err := newChart(&j, cfg)
				if err != nil {
					return nil, err
				}
				if err = charts.Add(chart); err != nil {
					return nil, err
				}
			}
		} else {
			chart, err := newChart(nil, cfg)
			if err != nil {
				return nil, err
			}
			if err = charts.Add(chart); err != nil {
				return nil, err
			}
		}
	}
	return charts, nil
}

// newChart creates news chart based on 'ChartsConfig', 'id' and 'oidIndex'
// parameters. oidIndex is optional param, which decided whether to add an
// index to OID value or not.
func newChart(oidIndex *int, s ChartsConfig) (*module.Chart, error) {
	c := defaultSNMPchart.Copy()
	c.ID = s.ID
	c.Title = s.Title

	if oidIndex != nil {
		c.ID = fmt.Sprintf("%s_%d", c.ID, *oidIndex)
		c.Title = fmt.Sprintf("%s %d", c.Title, *oidIndex)
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

	if c.Priority = s.Priority; c.Priority < module.Priority {
		c.Priority += module.Priority
	}

	for _, d := range s.Dimensions {
		oid := d.OID
		if oidIndex != nil {
			oid = fmt.Sprintf("%s.%d", oid, *oidIndex)
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

		if err := c.AddDim(dim); err != nil {
			return nil, err
		}
	}

	return c, nil
}

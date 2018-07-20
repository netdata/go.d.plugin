package raw

import (
	"errors"
	"fmt"

	"github.com/l2isbad/go.d.plugin/internal/pkg/utils"
)

type (
	Order       = utils.StringSlice
	Definitions []Chart

	Charts struct {
		Order       Order
		Definitions Definitions
	}
)

// ---------------------------------------------------------------------------------------------------------------------

// CHART GETTER

// GetChartByID returns chart id.
func (c *Charts) GetChartByID(chartID string) *Chart {
	if idx := c.Definitions.index(chartID); idx != -1 {
		return &c.Definitions[idx]
	}
	return nil
}

// LookupChartByID looks up chart by id.
func (c *Charts) LookupChartByID(chartID string) (*Chart, bool) {
	if v := c.GetChartByID(chartID); v != nil {
		return v, true
	}
	return nil, false
}

// GetChartByIndex returns chart by index.
func (c *Charts) GetChartByIndex(idx int) *Chart {
	if idx >= 0 && idx < len(c.Definitions) {
		return &c.Definitions[idx]
	}
	return nil
}

// LookupChartByIndex looks up chart by index.
func (c *Charts) LookupChartByIndex(idx int) (*Chart, bool) {
	if v := c.GetChartByIndex(idx); v != nil {
		return v, true
	}
	return nil, false
}

// ---------------------------------------------------------------------------------------------------------------------

// CHART DELETER

// DeleteChartByID deletes chart by id.
func (c *Charts) DeleteChartByID(chartID string) error {
	if idx := c.Definitions.index(chartID); idx != -1 {
		c.Order.DeleteByID(chartID)
		c.Definitions = append(c.Definitions[:idx], c.Definitions[idx+1:]...)
		return nil
	}
	return errors.New("nonexistent chart")
}

// DeleteChartByIndex deletes chart by index.
func (c *Charts) DeleteChartByIndex(idx int) error {
	if idx >= 0 && idx < len(c.Definitions) {
		c.Order.DeleteByID(c.GetChartByIndex(idx).ID)
		c.Definitions = append(c.Definitions[:idx], c.Definitions[idx+1:]...)
		return nil
	}
	return errors.New("nonexistent chart")
}

// ---------------------------------------------------------------------------------------------------------------------

// CHART ADDER

// AddChartOrder adds valid non duplicate chart to Definitions and to Order.
func (c *Charts) AddChartNoOrder(chart Chart) error {
	if err := chart.IsValid(); err != nil {
		return fmt.Errorf("invalid chart '%s' (%s)", chart.ID, err)
	}
	if c.Definitions.index(chart.ID) != -1 {
		return fmt.Errorf("duplicate chart '%s'", chart.ID)
	}
	c.Definitions = append(c.Definitions, chart)
	return nil
}

// AddChart adds valid non duplicate chart to Definitions.
func (c *Charts) AddChart(chart Chart) error {
	if err := c.AddChartNoOrder(chart); err != nil {
		return err
	}
	if !c.Order.Include(chart.ID) {
		c.Order.Append(chart.ID)
	}
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// Copy makes a full copy of charts.
func (c *Charts) Copy() *Charts {
	rv := Charts{}
	for _, v := range c.Order {
		rv.Order.Append(v)
	}
	for _, v := range c.Definitions {
		rv.Definitions = append(rv.Definitions, v.Copy())
	}
	return &rv
}

func (d *Definitions) index(chartID string) int {
	for idx, c := range *d {
		if c.ID == chartID {
			return idx
		}
	}
	return -1
}

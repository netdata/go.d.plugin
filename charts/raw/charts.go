package raw

import (
	"errors"
	"fmt"

	"github.com/l2isbad/go.d.plugin/shared"
)

type (
	Charts struct {
		Order       Order
		Definitions Definitions
	}
	Order       = shared.StringSlice
	Definitions []Chart
)

// ---------------------------------------------------------------------------------------------------------------------

// CHART GETTER

// GetChartByID returns chart id.
func (c *Charts) GetChartByID(chartID string) *Chart {
	if idx, ok := c.Definitions.index(chartID); ok {
		return &c.Definitions[idx]
	}
	return nil
}

// GetChartByIndex returns chart by index.
func (c *Charts) GetChartByIndex(idx int) *Chart {
	if idx >= 0 && idx < len(c.Definitions) {
		return &c.Definitions[idx]
	}
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// CHART DELETER

// DeleteChartByID deletes chart by id.
func (c *Charts) DeleteChartByID(chartID string) error {
	if idx, ok := c.Definitions.index(chartID); ok {
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

// AddChart adds valid non duplicate chart to Definitions and optionally to Order.
func (c *Charts) AddChart(chart Chart, addToOrder bool) error {
	if err := chart.IsValid(); err != nil {
		return fmt.Errorf("invalid chart '%s' (%s)", chart.ID, err)
	}
	if _, ok := c.Definitions.index(chart.ID); ok {
		return fmt.Errorf("duplicate chart '%s'", chart.ID)
	}
	c.Definitions = append(c.Definitions, chart)
	if addToOrder {
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
		rv.Definitions = append(rv.Definitions, v.copy())
	}
	return &rv
}

func (d *Definitions) index(chartID string) (int, bool) {
	for idx, c := range *d {
		if c.ID == chartID {
			return idx, true
		}
	}
	return 0, false
}

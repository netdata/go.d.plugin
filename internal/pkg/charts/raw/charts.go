package raw

import (
	"fmt"
	"github.com/l2isbad/go.d.plugin/internal/pkg/utils"
)

type (
	Order       = utils.StringSlice
	Definitions []*Chart

	Charts struct {
		Order       Order
		Definitions Definitions
	}
)

// ---------------------------------------------------------------------------------------------------------------------

// CHART GETTER

// GetChartByID returns chart id.
func (c Charts) GetChartByID(id string) *Chart {
	idx := c.Definitions.index(id)
	if idx == -1 {
		return nil
	}
	return c.Definitions[idx]
}

// LookupChartByID looks up chart by id.
func (c Charts) LookupChartByID(id string) (*Chart, bool) {
	v := c.GetChartByID(id)
	if v != nil {
		return v, true
	}
	return nil, false
}

// GetChartByIndex returns chart by index.
func (c Charts) GetChartByIndex(idx int) *Chart {
	ok := idx >= 0 && idx < len(c.Definitions)
	if !ok {
		return nil
	}
	return c.Definitions[idx]
}

// LookupChartByIndex looks up chart by index.
func (c Charts) LookupChartByIndex(idx int) (*Chart, bool) {
	v := c.GetChartByIndex(idx)
	if v != nil {
		return v, true
	}
	return nil, false
}

// ---------------------------------------------------------------------------------------------------------------------

// CHART DELETER

// DeleteChartByID deletes chart by id.
func (c *Charts) DeleteChartByID(id string) bool {
	idx := c.Definitions.index(id)
	if idx == -1 {
		return false
	}

	c.Order.DeleteByID(id)
	c.Definitions = append(c.Definitions[:idx], c.Definitions[idx+1:]...)
	return true
}

// DeleteChartByIndex deletes chart by index.
func (c *Charts) DeleteChartByIndex(idx int) bool {
	ok := idx >= 0 && idx < len(c.Definitions)
	if !ok {
		return false
	}

	c.Order.DeleteByID(c.GetChartByIndex(idx).ID)
	c.Definitions = append(c.Definitions[:idx], c.Definitions[idx+1:]...)
	return true
}

// ---------------------------------------------------------------------------------------------------------------------

// CHART ADDER

// FIXME: change method name
// AddChartOrder adds valid non duplicate chart to Definitions and to Order.
func (c *Charts) AddChartNoOrder(chart *Chart) error {
	err := chart.IsValid()
	if err != nil {
		return fmt.Errorf("invalid chart '%s': %s", chart.ID, err)
	}

	idx := c.Definitions.index(chart.ID)
	if idx != -1 {
		return fmt.Errorf("duplicate chart '%s'", chart.ID)
	}

	c.Definitions = append(c.Definitions, chart)
	return nil
}

// AddChart adds valid non duplicate chart to Definitions.
func (c *Charts) AddChart(chart *Chart) error {
	err := c.AddChartNoOrder(chart)
	if err != nil {
		return err
	}

	if !c.Order.Include(chart.ID) {
		c.Order.Append(chart.ID)
	}
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// Copy makes a full copy of charts.
func (c Charts) Copy() *Charts {
	rv := Charts{}
	for idx := range c.Order {
		rv.Order.Append(c.Order[idx])
	}

	for idx := range c.Definitions {
		rv.Definitions = append(rv.Definitions, c.Definitions[idx].Copy())
	}

	return &rv
}

func (d Definitions) index(id string) int {
	for idx := range d {
		if d[idx].ID == id {
			return idx
		}
	}
	return -1
}

package raw

import (
	"errors"
	"fmt"
	"strings"
)

const (
	IdxChartTitle = iota
	IdxChartUnits
	IdxChartFamily
	IdxChartContext
	IdxChartType
	IdxChartOverrideID
)

const (
	Line    = "line"
	Area    = "area"
	Stacked = "stacked"
)

var (
	defaultChartType = Line
)

var (
	errNoID     = errors.New("no 'id' specified")
	errNoTitle  = errors.New("no 'title' specified")
	errNoUnits  = errors.New("no 'units' specified")
	errNoFamily = errors.New("no 'family' specified")
)

// NewChart creates new Chart from id, options and optionally dimensions.
func NewChart(id string, options Options, dims ...Dimension) *Chart {
	chart := Chart{ID: id, Options: options}
	for idx := range dims {
		chart.AddDim(dims[idx])
	}
	return &chart
}

type (
	Chart struct {
		ID         string
		Options    Options
		Dimensions Dimensions
		Variables  Variables
	}

	Options    [6]string
	Dimensions []Dimension
	Variables  []Variable
)

func (c Chart) IsValid() error {
	switch {
	case c.ID == "":
		return errNoID
	case c.Title() == "":
		return errNoTitle
	case c.Units() == "":
		return errNoUnits
	case c.Family() == "":
		return errNoFamily
	default:
		return nil
	}
}

// ---------------------------------------------------------------------------------------------------------------------

// FIELD GETTER

// Title returns 0 element of Options.
func (c Chart) Title() string {
	return c.Options[IdxChartTitle]
}

// Units returns 1 element of Options.
func (c Chart) Units() string {
	return c.Options[IdxChartUnits]
}

// Family returns 2 element of Options converted to lower case.
func (c Chart) Family() string {
	return strings.ToLower(c.Options[IdxChartFamily])
}

// Context returns 3 element of Options.
func (c Chart) Context() string {
	return c.Options[IdxChartContext]
}

// ChartType returns 4 element of Options.
func (c Chart) ChartType() string {
	if ValidChartType(c.Options[IdxChartType]) {
		return c.Options[IdxChartType]
	}
	return defaultChartType
}

// OverrideID returns 5 element of Options.
func (c Chart) OverrideID() string {
	return c.Options[IdxChartOverrideID]
}

// ---------------------------------------------------------------------------------------------------------------------

// FIELD SETTER

// SetTitle sets 0 element of Options.
func (c *Chart) SetTitle(s string) *Chart {
	c.Options[IdxChartTitle] = s
	return c
}

// SetUnits sets 1 element of Options.
func (c *Chart) SetUnits(s string) *Chart {
	c.Options[IdxChartUnits] = s
	return c
}

// SetFamily sets 2 element of Options.
func (c *Chart) SetFamily(s string) *Chart {
	c.Options[IdxChartFamily] = s
	return c
}

// SetContext sets 3 element of Options.
func (c *Chart) SetContext(s string) *Chart {
	c.Options[IdxChartContext] = s
	return c
}

// SetChartType sets 4 element of Options.
func (c *Chart) SetChartType(s string) *Chart {
	c.Options[IdxChartType] = s
	return c
}

// SetOverrideID sets 5 element of Options.
func (c *Chart) SetOverrideID(s string) *Chart {
	c.Options[IdxChartOverrideID] = s
	return c
}

// ---------------------------------------------------------------------------------------------------------------------

// DIM/VAR GETTER

// GetDimByID returns dimension by id.
func (c Chart) GetDimByID(dimID string) *Dimension {
	idx := c.indexDim(dimID)
	if idx == -1 {
		return nil
	}
	return &c.Dimensions[idx]
}

// LookupDimByID looks up dimension by id.
func (c Chart) LookupDimByID(dimID string) (*Dimension, bool) {
	d := c.GetDimByID(dimID)
	if d != nil {
		return d, true
	}
	return nil, false
}

// GetDimByIndex returns dimension by index.
func (c Chart) GetDimByIndex(idx int) *Dimension {
	ok := idx >= 0 && idx < len(c.Dimensions)
	if !ok {
		return nil
	}
	return &c.Dimensions[idx]
}

// LookupDimByIndex looks dimension by index.
func (c Chart) LookupDimByIndex(idx int) (*Dimension, bool) {
	d := c.GetDimByIndex(idx)
	if d != nil {
		return d, true
	}
	return nil, false
}

// GetVarByID returns variable by id.
func (c Chart) GetVarByID(varID string) *Variable {
	idx := c.indexVar(varID)
	if idx == -1 {
		return nil
	}
	return &c.Variables[idx]
}

// LookupVarByID looks up variable by id.
func (c Chart) LookupVarByID(id string) (*Variable, bool) {
	v := c.GetVarByID(id)
	if v != nil {
		return v, true
	}
	return nil, false
}

// ---------------------------------------------------------------------------------------------------------------------

// DIM/VAR DELETER

// DeleteDimByID deletes dimension by id.
func (c *Chart) DeleteDimByID(id string) bool {
	idx := c.indexDim(id)
	if idx == -1 {
		return false
	}
	c.Dimensions = append(c.Dimensions[:idx], c.Dimensions[idx+1:]...)
	return true
}

// DeleteDimByIndex deletes dimension by index.
func (c *Chart) DeleteDimByIndex(idx int) bool {
	ok := idx >= 0 && idx < len(c.Dimensions)
	if !ok {
		return false
	}
	c.Dimensions = append(c.Dimensions[:idx], c.Dimensions[idx+1:]...)
	return true
}

// DeleteVarByID deletes variable by id.
func (c *Chart) DeleteVarByID(varID string) bool {
	idx := c.indexVar(varID)
	if idx == -1 {
		return false
	}
	c.Variables = append(c.Variables[:idx], c.Variables[idx+1:]...)
	return true
}

// ---------------------------------------------------------------------------------------------------------------------

// DIM/VAR ADDER

// AddDim adds valid non duplicate dimension.
func (c *Chart) AddDim(d Dimension) error {
	err := d.IsValid()
	if err != nil {
		return fmt.Errorf("chart '%s' invalid dim: %s", c.ID, err)
	}

	idx := c.indexDim(d.ID())
	if idx != -1 {
		return fmt.Errorf("chart '%s' duplicate dim: %s", c.ID, d.ID())
	}

	c.Dimensions = append(c.Dimensions, d)
	return nil
}

// AddVar adds valid non duplicate variable.
func (c *Chart) AddVar(v Variable) error {
	err := v.IsValid()
	if err != nil {
		return fmt.Errorf("chart '%s' invalid var: %s", c.ID, err)
	}

	idx := c.indexVar(v.ID())
	if idx != -1 {
		return fmt.Errorf("chart '%s' duplicate var: %s", c.ID, v.ID())
	}

	c.Variables = append(c.Variables, v)
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

func (c Chart) indexDim(dimID string) int {
	for idx, dim := range c.Dimensions {
		if dim.ID() == dimID {
			return idx
		}
	}
	return -1
}

func (c Chart) indexVar(varID string) int {
	for idx, v := range c.Variables {
		if v.ID() == varID {
			return idx
		}
	}
	return -1
}

func (c Chart) Copy() *Chart {
	chart := Chart{ID: c.ID, Options: c.Options}
	for idx := range c.Dimensions {
		chart.AddDim(c.Dimensions[idx])
	}
	for idx := range c.Variables {
		chart.AddVar(c.Variables[idx])
	}
	return &chart
}

// ValidChartType returns whether the chart type is valid.
// Valid chart types: "line", "area", "stacked".
func ValidChartType(t string) bool {
	switch t {
	case Line, Area, Stacked:
		return true
	}
	return false
}

package cooked

import (
	"fmt"
	"strings"

	"github.com/l2isbad/go.d.plugin/charts/raw"
	"github.com/l2isbad/go.d.plugin/logger"
)

func newChart(c *raw.Chart, bc baseConfHook, priority int) (*Chart, error) {
	if err := c.IsValid(); err != nil {
		return nil, err
	}

	newChart := &Chart{
		bc:         bc,
		id:         c.ID,
		overrideID: c.OverrideID(),
		title:      c.Title(),
		units:      c.Units(),
		family:     c.Family(),
		context:    c.Context(),
		chartType:  c.ChartType(),
		priority:   priority,
		variables:  make(map[string]*variable),
		flagsChart: &flagsChart{push: true},
	}

	for _, d := range c.Dimensions {
		err := newChart.AddDim(d)
		if err != nil {
			logger.CacheGet(bc).Error(err)
			continue
		}
	}
	for _, v := range c.Variables {
		err := newChart.AddVar(v)
		if err != nil {
			logger.CacheGet(bc).Error(err)
			continue
		}
	}

	return newChart, nil
}

type (
	Chart struct {
		bc            baseConfHook
		id            string
		overrideID    string
		title         string
		units         string
		family        string
		context       string
		chartType     string
		priority      int
		FailedUpdates int
		dimensions
		variables
		*flagsChart
	}
	dimensions []*dimension // dimension order is a thing so dimensions must be a slice
	variables  map[string]*variable
)

// ---------------------------------------------------------------------------------------------------------------------

// FIELD GETTER

// ID returns Chart id.
func (c *Chart) ID() string {
	return c.id
}

// OverrideID returns Chart override id.
func (c *Chart) OverrideID() string {
	return c.overrideID
}

// Title returns Chart title.
func (c *Chart) Title() string {
	return c.title
}

// Units returns Chart units.
func (c *Chart) Units() string {
	return c.units
}

// Family returns Chart family.
func (c *Chart) Family() string {
	return c.family
}

// Context returns Chart context.
// If context not specified it returns "moduleName.chartID".
func (c *Chart) Context() string {
	if c.context == "" {
		return fmt.Sprintf("%s.%s", c.bc.ModuleName(), c.id)
	}
	return c.context
}

// ChartType returns Chart type.
func (c *Chart) ChartType() string {
	return c.chartType
}

// ---------------------------------------------------------------------------------------------------------------------

// FIELD SETTER

// SetID sets Chart id.
func (c *Chart) SetID(id string) *Chart {
	c.id = id
	return c
}

// SetOverrideID sets Chart override id.
func (c *Chart) SetOverrideID(id string) *Chart {
	c.overrideID = id
	return c
}

// SetTitle sets Chart title.
func (c *Chart) SetTitle(t string) *Chart {
	c.title = t
	return c
}

// SetUnits sets Chart units.
func (c *Chart) SetUnits(u string) *Chart {
	c.units = u
	return c
}

// SetFamily sets Chart family converted to lower case.
func (c *Chart) SetFamily(f string) *Chart {
	c.family = strings.ToLower(f)
	return c
}

// SetContext sets Chart context.
func (c *Chart) SetContext(ctx string) *Chart {
	c.context = ctx
	return c
}

// SetChartType sets Chart type (only if type is a valid Chart type).
func (c *Chart) SetChartType(t string) *Chart {
	if raw.ValidChartType(t) {
		c.chartType = t
	}
	return c
}

// ---------------------------------------------------------------------------------------------------------------------

// DIM/VAR GETTER

// GetDimByID returns dimension by id.
func (c *Chart) GetDimByID(dimID string) *dimension {
	if idx := c.index(dimID); idx != -1 {
		return c.dimensions[idx]
	}
	return nil
}

// GetDimByIndex returns dimension by index.
func (c *Chart) GetDimByIndex(idx int) *dimension {
	if idx >= 0 && idx < len(c.dimensions) {
		return c.dimensions[idx]
	}
	return nil
}

// GetVarByID returns variable by id.
func (c *Chart) GetVarByID(varID string) *variable {
	return c.variables[varID]
}

// ---------------------------------------------------------------------------------------------------------------------

// DIM/VAR ADDER

// AddDim adds valid non duplicate dimension to dimensions.
func (c *Chart) AddDim(d raw.Dimension) error {
	if c.index(d.ID()) != -1 {
		return fmt.Errorf("chart '%s': duplicate dimension %s, skipping it", c.id, d.ID())
	}
	newDim, err := newDimension(d)
	if err != nil {
		return fmt.Errorf("chart '%s': invalid dimension (%s), skipping it", c.id, err)
	}
	c.dimensions = append(c.dimensions, newDim)
	c.Refresh()
	return nil
}

// AddVar adds valid variable to variables.
func (c *Chart) AddVar(v raw.Variable) error {
	newVar, err := newVariable(v)
	if err != nil {
		return fmt.Errorf("chart '%s': invalid variable (%s), skipping it", c.id, err)
	}
	c.variables[newVar.id] = newVar
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// CONVERTER TO NETDATA FORMAT

func (c *Chart) begin(sinceLast int) string {
	return fmt.Sprintf(formatChartBEGIN,
		c.bc.FullName(),
		c.id,
		sinceLast)
}

func (c *Chart) create() string {
	var dimensions, variables string
	chart := fmt.Sprintf(formatChartCREATE,
		c.bc.FullName(),
		c.id,
		c.overrideID,
		c.title,
		c.units,
		c.family,
		c.Context(),
		c.chartType,
		c.priority,
		c.bc.UpdateEvery(),
		c.bc.ModuleName())

	for _, dim := range c.dimensions {
		dimensions += dim.create()
	}

	for _, v := range c.variables {
		if v.value != 0 {
			variables += v.set(v.value)
		}
	}
	c.setPush(false)
	c.setCreated(true)
	return chart + dimensions + variables + "\n"
}

// ---------------------------------------------------------------------------------------------------------------------

// Obsolete pushes Chart to the netdata with "obsolete" option if Chart was created.
// Regardless of the created flag it sets obsolete flag to true.
func (c *Chart) Obsolete() {
	c.setObsoleted(true)
	if !c.isCreated() {
		return
	}
	SafePrint(fmt.Sprintf(formatChartOBSOLETE,
		c.bc.FullName(),
		c.id,
		c.overrideID,
		c.title,
		c.units,
		c.family,
		c.Context(),
		c.chartType,
		c.priority,
		c.bc.UpdateEvery(),
		c.bc.ModuleName()))
}

// Refresh sets Chart push flag to true. If Chart was obsoleted it also sets obsolete and created flags to false.
func (c *Chart) Refresh() {
	c.setPush(true)
	if c.IsObsoleted() {
		c.FailedUpdates = 0
		c.setCreated(false)
		c.setObsoleted(false)
	}
}

// Update does Chart data collection, Chart creating and updating. Returns true if at least one dimension was updated.
func (c *Chart) Update(data map[string]int64, interval int) bool {
	var updDim, updVar string

	for _, d := range c.dimensions {
		if value, ok := d.get(data); ok {
			updDim += d.set(value)
		}
		if d.push {
			c.setPush(true)
			d.push = false
		}
	}

	for _, v := range c.variables {
		if value, ok := data[v.id]; ok {
			updVar += v.set(value)
		}
	}

	if updDim == "" {
		c.FailedUpdates++
		c.setUpdated(false)
		return false
	}
	if !c.isUpdated() {
		interval = 0
	}
	if c.isPush() {
		SafePrint(c.create())
	}
	SafePrint(c.begin(interval), updDim, updVar, "END\n\n")
	c.setUpdated(true)
	c.FailedUpdates = 0

	return true
}

// CanBeUpdated returns whether the Chart can be updated with current data.
func (c *Chart) CanBeUpdated(data map[string]int64) bool {
	for _, dim := range c.dimensions {
		if _, ok := data[dim.id]; ok {
			return true
		}
	}
	return false
}

// Index finds dimension index by id.
func (c *Chart) index(dimID string) int {
	for idx, dim := range c.dimensions {
		if dim.id == dimID {
			return idx
		}
	}
	return -1
}

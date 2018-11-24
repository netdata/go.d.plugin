package modules

import (
	"strconv"
	"strings"
)

type (
	chartPriority int

	dimAlgo   string
	chartType string
	dimHidden bool
	dimDivMul int
)

const (
	defChartPriority = "70000"

	// Line chart type.
	Line chartType = "line"
	// Area chart type.
	Area chartType = "area"
	// Stacked chart type.
	Stacked chartType = "stacked"

	// Absolute dimension algorithm.
	// The value is to drawn as-is (interpolated to second boundary).
	Absolute dimAlgo = "absolute"
	// Incremental dimension algorithm.
	// The value increases over time, the difference from the last value is presented in the chart,
	// the server interpolates the value and calculates a per second figure.
	Incremental dimAlgo = "incremental"
	// PercentOfAbsolute dimension algorithm.
	// The percent of this value compared to the total of all dimensions.
	PercentOfAbsolute dimAlgo = "percentage-of-absolute-row"
	// PercentOfIncremental dimension algorithm.
	// The percent of this value compared to the incremental total of all dimensions
	PercentOfIncremental dimAlgo = "percentage-of-incremental-row"
)

func (c chartPriority) String() string {
	if c > 0 {
		return strconv.Itoa(int(c))
	}
	return defChartPriority
}

func (d dimAlgo) String() string {
	switch d {
	case Absolute, Incremental, PercentOfAbsolute, PercentOfIncremental:
		return string(d)
	}
	return ""
}

func (c chartType) String() string {
	switch c {
	case Line, Area, Stacked:
		return string(c)
	}
	return ""
}

func (d dimHidden) String() string {
	if d {
		return "hidden"
	}
	return ""
}

func (d dimDivMul) String() string {
	if d != 0 {
		return strconv.Itoa(int(d))
	}
	return ""
}

type (
	// Charts is a collection of charts.
	Charts []*Chart

	// Chart represents a chart.
	// For the full description please visit https://docs.netdata.cloud/collectors/plugins.d/#chart
	Chart struct {
		// typeID is the unique identification of the chart, if not specified,
		// the plugin will use job full name + chart ID as typeID (default behaviour).
		typeID string

		ID       string
		OverID   string
		Title    string
		Units    string
		Fam      string
		Ctx      string
		Type     chartType
		Priority chartPriority
		Opts

		Dims Dims
		Vars Vars

		Retries int

		// created flag is used to indicate whether the chart needs to be created by the plugin.
		created bool
		// updated flag is used to indicate whether the chart was updated on last data collection interval.
		updated bool
	}
	// Opts represents chart options.
	Opts struct {
		Obsolete   bool
		Detail     bool
		StoreFirst bool
		Hidden     bool
	}
	// Dims is a collection of dims.
	Dims []*Dim
	// Vars is a collection of vars.
	Vars []*Var

	// Dim represents a chart dimension.
	// For detailed description please visit https://docs.netdata.cloud/collectors/plugins.d/#dimension.
	Dim struct {
		ID     string
		Name   string
		Algo   dimAlgo
		Mul    dimDivMul
		Div    dimDivMul
		Hidden dimHidden
	}

	// Var represents a chart variable.
	// For detailed description please visit https://docs.netdata.cloud/collectors/plugins.d/#variable
	Var struct {
		ID    string
		Value int64
	}
)

func (o Opts) String() string {
	var opts []string

	if o.Obsolete {
		opts = append(opts, "obsolete")
	}
	if o.Detail {
		opts = append(opts, "detail")
	}
	if o.StoreFirst {
		opts = append(opts, "store_first")
	}
	if o.Hidden {
		opts = append(opts, "hidden")
	}

	return strings.Join(opts, " ")
}

// Add adds (appends) a variable number of charts.
// Chart won't be added if:
//   * chart isn't a valid chart
//   * charts have a chart with the same ID
func (c *Charts) Add(charts ...*Chart) {
	for _, v := range charts {
		if c.index(v.ID) != -1 || !v.IsValid() {
			continue
		}
		*c = append(*c, v)
	}
}

// Get returns the chart by ID.
func (c Charts) Get(chartID string) *Chart {
	idx := c.index(chartID)
	if idx == -1 {
		return nil
	}
	return c[idx]
}

// Has returns true if charts contain the chart with the given ID, false otherwise.
func (c Charts) Has(chartID string) bool {
	idx := c.index(chartID)
	if idx == -1 {
		return false
	}
	return true
}

// Remove removes the chart from charts by ID,
// it returns true if charts contain the chart with the given ID, false otherwise.
// Avoid to use it in runtime.
func (c *Charts) Remove(chartID string) bool {
	idx := c.index(chartID)
	if idx == -1 {
		return false
	}
	*c = append((*c)[:idx], (*c)[idx+1:]...)
	return true
}

// Copy returns a deep copy of charts.
func (c Charts) Copy() *Charts {
	charts := Charts{}
	for idx := range c {
		charts = append(charts, c[idx].Copy())
	}
	return &charts
}

func (c Charts) index(chartID string) int {
	for idx := range c {
		if c[idx].ID == chartID {
			return idx
		}
	}
	return -1
}

// MarkNotCreated changes 'created' chart flag to false.
// Use it to add dimension in runtime.
func (c *Chart) MarkNotCreated() {
	c.created = false
}

// AddDim adds new dimension to the chart dimensions.
func (c *Chart) AddDim(newDim *Dim) bool {
	if c.indexDim(newDim.ID) != -1 || !newDim.IsValid() {
		return false
	}
	c.Dims = append(c.Dims, newDim)

	return true
}

// AddVar adds new variable to the chart variables.
func (c *Chart) AddVar(newVar *Var) bool {
	if c.indexVar(newVar.ID) != -1 || !newVar.IsValid() {
		return false
	}
	c.Vars = append(c.Vars, newVar)

	return true
}

// GetDim returns dimension by ID.
func (c *Chart) GetDim(dimID string) *Dim {
	if idx := c.indexDim(dimID); idx != -1 {
		return c.Dims[idx]
	}
	return nil
}

// RemoveDim removes dimension by ID,
// it returns true if the chart contains dimension with the given ID, false otherwise.
// Avoid to use it in runtime.
func (c *Chart) RemoveDim(dimID string) bool {
	idx := c.indexDim(dimID)
	if idx == -1 {
		return false
	}
	c.Dims = append(c.Dims[:idx], c.Dims[idx+1:]...)

	return true
}

// HasDim returns true if the chart contains dimension with the given ID, false otherwise.
func (c Chart) HasDim(dimID string) bool {
	return c.indexDim(dimID) != -1
}

// Copy returns a deep copy of the chart.
func (c Chart) Copy() *Chart {
	chart := c
	chart.Dims = Dims{}
	chart.Vars = Vars{}

	for idx := range c.Dims {
		chart.Dims = append(chart.Dims, c.Dims[idx].copy())
	}
	for idx := range c.Vars {
		chart.Vars = append(chart.Vars, c.Vars[idx].copy())
	}

	return &chart
}

func (c Chart) indexDim(dimID string) int {
	for idx := range c.Dims {
		if c.Dims[idx].ID == dimID {
			return idx
		}
	}
	return -1
}

func (c Chart) indexVar(varID string) int {
	for idx := range c.Vars {
		if c.Vars[idx].ID == varID {
			return idx
		}
	}
	return -1
}

// IsValid returns whether the chart is valid.
// Chart is valid if it has non empty ID, Title and Units.
func (c Chart) IsValid() bool {
	return c.ID != "" && c.Title != "" && c.Units != ""
}

func (d Dim) copy() *Dim {
	return &d
}

// IsValid returns whether the dimension is valid.
// Dimension is valid if it has non empty ID.
func (d Dim) IsValid() bool {
	return d.ID != ""
}

func (v Var) copy() *Var {
	return &v
}

// IsValid returns whether the variable is valid.
// Variable is valid if it has non empty ID.
func (v Var) IsValid() bool {
	return v.ID != ""
}

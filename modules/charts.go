package modules

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type (
	dimAlgo   string
	chartType string
	dimHidden bool
	dimDivMul int
)

const (
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
	// Charts is a collection of ChartsFunc.
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
		Priority int
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

// Add adds (appends) a variable number of Charts.
func (c *Charts) Add(charts ...*Chart) error {
	for _, chart := range charts {
		if err := checkChart(chart); err != nil {
			return fmt.Errorf("error on adding chart : %s", err)
		}
		if c.index(chart.ID) != -1 {
			return fmt.Errorf("error on adding chart : '%s' is already in charts", chart.ID)
		}
		*c = append(*c, chart)
	}

	return nil
}

// Get returns the chart by ID.
func (c Charts) Get(chartID string) *Chart {
	idx := c.index(chartID)
	if idx == -1 {
		return nil
	}
	return c[idx]
}

// Has returns true if ChartsFunc contain the chart with the given ID, false otherwise.
func (c Charts) Has(chartID string) bool {
	idx := c.index(chartID)
	if idx == -1 {
		return false
	}
	return true
}

// Remove removes the chart from Charts by ID.
// Avoid to use it in runtime.
func (c *Charts) Remove(chartID string) error {
	idx := c.index(chartID)
	if idx == -1 {
		return fmt.Errorf("error on removing chart : '%s' is not in charts", chartID)
	}
	*c = append((*c)[:idx], (*c)[idx+1:]...)
	return nil
}

// Copy returns a deep copy of ChartsFunc.
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
func (c *Chart) AddDim(newDim *Dim) error {
	if err := checkDim(newDim); err != nil {
		return fmt.Errorf("error on adding dim to chart '%s' : %s", c.ID, err)
	}
	if c.indexDim(newDim.ID) != -1 {
		return fmt.Errorf("error on adding dim : '%s' is already in chart '%s' dims", newDim.ID, c.ID)
	}
	c.Dims = append(c.Dims, newDim)

	return nil
}

// AddVar adds new variable to the chart variables.
func (c *Chart) AddVar(newVar *Var) error {
	if err := checkVar(newVar); err != nil {
		return fmt.Errorf("error on adding var to chart '%s' : %s", c.ID, err)
	}
	if c.indexVar(newVar.ID) != -1 {
		return fmt.Errorf("error on adding var : '%s' is already in chart '%s' vars", newVar.ID, c.ID)
	}
	c.Vars = append(c.Vars, newVar)

	return nil
}

// GetDim returns dimension by ID.
func (c *Chart) GetDim(dimID string) *Dim {
	if idx := c.indexDim(dimID); idx != -1 {
		return c.Dims[idx]
	}
	return nil
}

// RemoveDim removes dimension by ID.
// Avoid to use it in runtime.
func (c *Chart) RemoveDim(dimID string) error {
	idx := c.indexDim(dimID)
	if idx == -1 {
		return fmt.Errorf("error on removing dim : '%s' isn't in chart '%s'", dimID, c.ID)
	}
	c.Dims = append(c.Dims[:idx], c.Dims[idx+1:]...)

	return nil
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

func (d Dim) copy() *Dim {
	return &d
}

func (v Var) copy() *Var {
	return &v
}

// CheckCharts checks charts
func CheckCharts(charts ...*Chart) error {
	for _, chart := range charts {
		if err := checkChart(chart); err != nil {
			return err
		}
	}
	return nil
}

func checkChart(chart *Chart) error {
	if chart.ID == "" {
		return errors.New("empty chart id")
	}

	if chart.Title == "" {
		return fmt.Errorf("empty title in chart '%s'", chart.ID)
	}

	if chart.Units == "" {
		return fmt.Errorf("empty units in chart '%s'", chart.ID)
	}

	if id := checkID(chart.ID); id != -1 {
		return fmt.Errorf("unacceptable symbol in chart id '%s' : '%s'", chart.ID, string(id))
	}

	set := make(map[string]bool)

	for _, d := range chart.Dims {
		if err := checkDim(d); err != nil {
			return err
		}
		if set[d.ID] {
			return fmt.Errorf("duplicate dim '%s' in chart '%s'", d.ID, chart.ID)
		}
		set[d.ID] = true
	}

	set = make(map[string]bool)

	for _, v := range chart.Vars {
		if err := checkVar(v); err != nil {
			return err
		}
		if set[v.ID] {
			return fmt.Errorf("duplicate var '%s' in chart '%s'", v.ID, chart.ID)
		}
		set[v.ID] = true
	}
	return nil
}

func checkDim(d *Dim) error {
	if d.ID == "" {
		return errors.New("empty dim id")
	}
	if id := checkID(d.ID); id != -1 {
		return fmt.Errorf("unacceptable symbol in dim id '%s' : '%s'", d.ID, string(id))
	}
	return nil
}

func checkVar(v *Var) error {
	if v.ID == "" {
		return errors.New("empty var id")
	}
	if id := checkID(v.ID); id != -1 {
		return fmt.Errorf("unacceptable symbol in var id '%s' : '%s'", v.ID, string(id))
	}
	return nil
}

func checkID(id string) int {
	for _, r := range id {
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_') {
			return int(r)
		}
	}
	return -1
}

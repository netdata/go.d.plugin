package modules

import (
	"strconv"
	"strings"
)

type (
	dimAlgo   string
	chartType string
	dimHidden bool
	dimDivMul int
)

const (
	Line    chartType = "line"
	Area    chartType = "area"
	Stacked chartType = "stacked"

	Absolute             dimAlgo = "absolute"
	Incremental          dimAlgo = "incremental"
	PercentOfAbsolute    dimAlgo = "percentage-of-absolute-row"
	PercentOfIncremental dimAlgo = "percentage-of-incremental-row"
)

func (c chartType) String() string {
	switch c {
	case Line, Area, Stacked:
		return string(c)
	}
	return ""
}

func (d dimAlgo) String() string {
	switch d {
	case Absolute, Incremental, PercentOfAbsolute, PercentOfIncremental:
		return string(d)
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
	if d == 0 {
		return ""
	}
	return strconv.Itoa(int(d))
}

type (
	Charts []*Chart

	Chart struct {
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

		pushed  bool
		updated bool
	}
	Opts struct {
		Obsolete   bool
		Detail     bool
		StoreFirst bool
		Hidden     bool
	}
	Dims []*Dim
	Vars []*Var

	Dim struct {
		ID     string
		Name   string
		Algo   dimAlgo
		Mul    dimDivMul
		Div    dimDivMul
		Hidden dimHidden
	}

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

//------------------------------------------------------Charts----------------------------------------------------------

func NewCharts(charts ...*Chart) *Charts {
	c := new(Charts)
	c.Add(charts...)
	return c
}

func (c *Charts) Add(charts ...*Chart) {
	for _, v := range charts {
		if c.index(v.ID) != -1 || !v.isValid() {
			continue
		}
		*c = append(*c, v.Copy())
	}
}

func (c Charts) Get(chartID string) *Chart {
	idx := c.index(chartID)
	if idx == -1 {
		return nil
	}
	return c[idx]
}

func (c Charts) Has(chartID string) bool {
	idx := c.index(chartID)
	if idx == -1 {
		return false
	}
	return true
}

func (c *Charts) Remove(chartID string) bool {
	idx := c.index(chartID)
	if idx == -1 {
		return false
	}
	*c = append((*c)[:idx], (*c)[idx+1:]...)
	return true
}

func (c Charts) Copy() Charts {
	charts := Charts{}
	for idx := range c {
		charts = append(charts, c[idx].Copy())
	}
	return charts
}

func (c Charts) index(chartID string) int {
	for idx := range c {
		if c[idx].ID == chartID {
			return idx
		}
	}
	return -1
}

//------------------------------------------------------chart-----------------------------------------------------------

func (c *Chart) MarkPush() {
	c.pushed = false
}

func (c *Chart) AddDim(newDim *Dim) bool {
	if c.indexDim(newDim.ID) != -1 || !newDim.isValid() {
		return false
	}
	c.Dims = append(c.Dims, newDim)

	return true
}

func (c *Chart) AddVar(newVar *Var) bool {
	if c.indexVar(newVar.ID) != -1 || !newVar.isValid() {
		return false
	}
	c.Vars = append(c.Vars, newVar)

	return true
}

func (c *Chart) RemoveDim(dimID string) bool {
	idx := c.indexDim(dimID)
	if idx == -1 {
		return false
	}
	c.Dims = append(c.Dims[:idx], c.Dims[idx+1:]...)

	return true
}

func (c Chart) HasDim(dimID string) bool {
	idx := c.indexDim(dimID)
	if idx == -1 {
		return false
	}
	return true
}

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

func (c Chart) isValid() bool {
	return c.ID != "" && c.Title != "" && c.Units != ""
}

//------------------------------------------------------dimension-------------------------------------------------------

func (d Dim) copy() *Dim {
	return &d
}

func (d Dim) isValid() bool {
	return d.ID != ""
}

//------------------------------------------------------Variable--------------------------------------------------------

func (v Var) copy() *Var {
	return &v
}

func (v Var) isValid() bool {
	return v.ID != ""
}

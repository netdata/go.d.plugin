package modules

type chartType string

const (
	Line    chartType = "line"
	Area    chartType = "area"
	Stacked chartType = "stacked"
)

type dimAlgo string

const (
	Absolute             dimAlgo = "absolute"
	Incremental          dimAlgo = "incremental"
	PercentOfAbsolute    dimAlgo = "percentage-of-absolute-row"
	PercentOfIncremental dimAlgo = "percentage-of-incremental-row"
)

type (
	Charts []*Chart

	Chart struct {
		ID string
		Opts
		Dims Dims
		Vars Vars

		state    state
		priority int
		retries  int
	}
	Opts struct {
		Title  string
		Units  string
		Fam    string
		Ctx    string
		Type   chartType
		OverID string
	}
	Dims []*Dim
	Vars []*Var

	Dim struct {
		ID     string
		Name   string
		Algo   dimAlgo
		Mul    int
		Div    int
		Hidden bool
	}

	Var struct {
		ID    string
		Value int64
	}
)

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
		*c = append(*c, v)
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

func (c *Chart) AddDim(newDim *Dim) bool {
	if c.indexDim(newDim.ID) != -1 || !newDim.isValid() {
		return false
	}
	c.Dims = append(c.Dims, newDim)
	c.dispatch(renewTrigger)

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
	c.dispatch(rebuildTrigger)

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
	chart := new(Chart)
	chart.ID = c.ID
	chart.Opts = c.Opts

	for idx := range c.Dims {
		chart.Dims = append(chart.Dims, c.Dims[idx].copy())
	}
	for idx := range c.Vars {
		chart.Vars = append(chart.Vars, c.Vars[idx].copy())
	}

	return chart
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

func (c *Chart) dispatch(trigger trigger) {
	c.state = c.state.dispatch(trigger)
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

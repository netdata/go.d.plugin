package charts

var (
	Line    = chartType{"line"}
	Area    = chartType{"area"}
	Stacked = chartType{"stacked"}
)

type (
	Chart struct {
		ID   string
		Opts
		Dims Dims
		Vars Vars
	}
	Opts struct {
		Title      string
		Units      string
		Family     string
		Context    string
		Type       chartType
		OverrideID string
	}
	Dims []Dim
	Vars []Var
	chartType  struct {
		t string
	}
)

func (t chartType) String() string {
	return t.t
}

// DIM/VAR ADDER -------------------------------------------------------------------------------------------------------

func (c *Chart) AddDim(d Dim) {
	if c.indexDim(d.ID) != -1 {
		return
	}
	c.Dims = append(c.Dims, d)
}

func (c *Chart) AddVar(v Var) {
	if c.indexVar(v.ID) != -1 {
		return
	}
	c.Vars = append(c.Vars, v)
}

// DIM/VAR GETTER ------------------------------------------------------------------------------------------------------

func (c Chart) GetDimByID(id string) *Dim {
	idx := c.indexDim(id)
	if idx == -1 {
		return nil
	}
	return &c.Dims[idx]
}

func (c Chart) LookupDimByID(id string) (*Dim, bool) {
	d := c.GetDimByID(id)
	if d == nil {
		return nil, false
	}
	return d, true
}

func (c *Chart) DeleteDimByID(id string) bool {
	idx := c.indexDim(id)
	if idx == -1 {
		return false
	}
	c.Dims = append(c.Dims[:idx], c.Dims[idx+1:]...)
	return true
}

func (c Chart) Copy() Chart {
	chart := c
	c.Dims = nil
	c.Vars = nil

	for idx := range c.Dims {
		chart.Dims = append(chart.Dims, c.Dims[idx].copy())
	}
	for idx := range c.Vars {
		chart.Vars = append(chart.Vars, c.Vars[idx].copy())
	}

	return chart
}

func (c Chart) indexDim(id string) int {
	for idx := range c.Dims {
		if c.Dims[idx].ID == id {
			return idx
		}
	}
	return -1
}

func (c Chart) indexVar(id string) int {
	for idx := range c.Vars {
		if c.Vars[idx].ID == id {
			return idx
		}
	}
	return -1
}

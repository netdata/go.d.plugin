package charts

var (
	Line    = chartType{"line"}
	Area    = chartType{"area"}
	Stacked = chartType{"stacked"}
)

type observer interface {
	Add(string)
	Delete(string)
	Update(string)
	Obsolete(string)
}

type chartType struct {
	t string
}

func (t chartType) String() string {
	return t.t
}

type (
	Chart struct {
		ID string
		Opts
		Dims Dims
		Vars Vars

		obs observer
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
)

// ---------------------------------------------------------------------------------------------------------------------
func (c *Chart) Register(o observer) {
	c.obs = o
}

func (c Chart) runtime() bool {
	return c.obs != nil
}

func (c Chart) Refresh() {
	if c.runtime() {
		c.obs.Update(c.ID)
	}
}

// ---------------------------------------------------------------------------------------------------------------------

func (c *Chart) AddDim(d *Dim) {
	if c.indexDim(d.ID) != -1 || !d.isValid() {
		return
	}
	c.Refresh()
	c.Dims = append(c.Dims, d)
}

func (c *Chart) AddVar(v *Var) {
	if c.indexVar(v.ID) != -1 || !v.isValid() {
		return
	}
	c.Refresh()
	c.Vars = append(c.Vars, v)
}

func (c Chart) GetDimByID(id string) *Dim {
	idx := c.indexDim(id)
	if idx == -1 {
		return nil
	}
	return c.Dims[idx]
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
	if c.runtime() {
		c.obs.Obsolete(c.ID)
		c.Refresh()
	}
	c.Dims = append(c.Dims[:idx], c.Dims[idx+1:]...)
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

func (c Chart) isValid() bool {
	return c.ID != "" && c.Title != "" && c.Units != ""
}

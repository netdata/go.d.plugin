package charts

var (
	Line    = chartType{"line"}
	Area    = chartType{"area"}
	Stacked = chartType{"stacked"}
)

type (
	Chart struct {
		ID string
		Options
		Dimensions Dimensions
		Variables  Variables
	}
	Options struct {
		Title      string
		Units      string
		Family     string
		Context    string
		Type       chartType
		OverrideID string
	}
	Dimensions []Dimension
	Variables  []Variable
	chartType  struct {
		t string
	}
)

func (t chartType) String() string {
	return t.t
}

// DIM/VAR ADDER -------------------------------------------------------------------------------------------------------

func (c *Chart) AddDim(d Dimension) {
	if c.indexDim(d.ID) != -1 {
		return
	}
	c.Dimensions = append(c.Dimensions, d)
}

func (c *Chart) AddVar(v Variable) {
	if c.indexVar(v.ID) != -1 {
		return
	}
	c.Variables = append(c.Variables, v)
}

// DIM/VAR GETTER ------------------------------------------------------------------------------------------------------

func (c Chart) GetDimByID(id string) *Dimension {
	idx := c.indexDim(id)
	if idx == -1 {
		return nil
	}
	return &c.Dimensions[idx]
}

func (c Chart) LookupDimByID(id string) (*Dimension, bool) {
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
	c.Dimensions = append(c.Dimensions[:idx], c.Dimensions[idx+1:]...)
	return true
}

func (c Chart) Copy() Chart {
	chart := c
	c.Dimensions = nil
	c.Variables = nil

	for idx := range c.Dimensions {
		chart.Dimensions = append(chart.Dimensions, c.Dimensions[idx].copy())
	}
	for idx := range c.Variables {
		chart.Variables = append(chart.Variables, c.Variables[idx].copy())
	}

	return chart
}

func (c Chart) indexDim(id string) int {
	for idx := range c.Dimensions {
		if c.Dimensions[idx].ID == id {
			return idx
		}
	}
	return -1
}

func (c Chart) indexVar(id string) int {
	for idx := range c.Variables {
		if c.Variables[idx].ID == id {
			return idx
		}
	}
	return -1
}

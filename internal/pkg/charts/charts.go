package charts

type Charts []*Chart

func New() *Charts {
	return &Charts{}
}

func (c *Charts) AddChart(charts ...*Chart) {
	for _, v := range charts {
		if c.index(v.ID) != -1 || !v.isValid() {
			continue
		}
		*c = append(*c, v)
	}
}

func (c Charts) GetChart(id string) *Chart {
	idx := c.index(id)
	if idx == -1 {
		return nil
	}
	return c[idx]
}

func (c Charts) LookupChart(id string) (*Chart, bool) {
	v := c.GetChart(id)
	if v == nil {
		return nil, false
	}
	return v, true
}

func (c *Charts) DeleteChart(id string) bool {
	idx := c.index(id)
	if idx == -1 {
		return false
	}
	v := (*c)[idx]

	if v.runtime() {
		v.obs.Delete(v.ID)
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

func (c Charts) index(id string) int {
	for idx := range c {
		if c[idx].ID == id {
			return idx
		}
	}
	return -1
}

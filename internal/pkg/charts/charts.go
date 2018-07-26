package charts

type Charts []Chart

func (c *Charts) Add(charts ...Chart) {
	for idx := range charts {
		*c = append(*c, charts[idx])
	}
}

func (c Charts) Get(id string) *Chart {
	idx := c.index(id)
	if idx == -1 {
		return nil
	}
	return &c[idx]
}

func (c Charts) Lookup(id string) (*Chart, bool) {
	v := c.Get(id)
	if v == nil {
		return nil, false
	}
	return v, true
}

func (c *Charts) Delete(id string) bool {
	idx := c.index(id)
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

func (c Charts) index(id string) int {
	for idx := range c {
		if c[idx].ID == id {
			return idx
		}
	}
	return -1
}

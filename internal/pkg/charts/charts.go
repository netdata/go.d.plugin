package charts

type Charts []*Chart

func NewCharts(charts ...*Chart) *Charts {
	c := new(Charts)
	c.Add(charts...)
	return c
}

func (c *Charts) AddAfter(id string, charts ...*Chart) {
	c.addAfter(id, false, charts...)
}

func (c *Charts) AddBefore(id string, charts ...*Chart) {
	c.addAfter(id, true, charts...)
}

func (c *Charts) addAfter(id string, before bool, charts ...*Chart) {
	idx := c.index(id)
	if idx == -1 {
		return
	}

	ch := NewCharts()
	for idx := range charts {
		if c.index(charts[idx].ID) == -1 && charts[idx].isValid() {
			ch.Add(charts[idx])
		}
	}

	if before {
		*c = append((*c)[:idx], append(*ch, (*c)[idx:]...)...)
	} else {
		*c = append((*c)[:idx+1], append(*ch, (*c)[idx+1:]...)...)
	}
}

func (c *Charts) Add(charts ...*Chart) {
	for _, v := range charts {
		if c.index(v.ID) != -1 || !v.isValid() {
			continue
		}
		*c = append(*c, v)
	}
}

func (c Charts) Get(id string) *Chart {
	idx := c.index(id)
	if idx == -1 {
		return nil
	}
	return c[idx]
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

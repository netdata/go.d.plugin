package modules

type S interface {
	SetUpdateEvery(int)
	SetChartsCleanup(int)
	SetDisabledByDefault()
}

type G interface {
	GetUpdateEvery() (int, bool)
	GetChartsCleanup() (int, bool)
	GetDisabledByDefault() bool
}

type moduleDefault struct {
	u *int
	c *int
	d bool
}

func (m *moduleDefault) SetUpdateEvery(v int) {
	if m.u == nil {
		m.u = new(int)
	}
	*m.u = v
}

func (m *moduleDefault) SetChartsCleanup(v int) {
	if m.c == nil {
		m.c = new(int)
	}
	*m.c = v
}

func (m *moduleDefault) SetDisabledByDefault() {
	m.d = true
}

func (m *moduleDefault) GetUpdateEvery() (int, bool) {
	if m.u == nil {
		return 0, false
	}
	return *m.u, true
}

func (m *moduleDefault) GetChartsCleanup() (int, bool) {
	if m.c == nil {
		return 0, false
	}
	return *m.c, true
}

func (m *moduleDefault) GetDisabledByDefault() bool {
	return m.d
}

var moduleDefaults = map[string]*moduleDefault{"_": {}}

func SetDefault() S {
	name := getFileName(2)
	if _, ok := moduleDefaults[name]; !ok {
		moduleDefaults[name] = new(moduleDefault)
	}
	return moduleDefaults[name]
}

func GetDefault(n string) G {
	v, ok := moduleDefaults[n]
	if !ok {
		return moduleDefaults["_"]
	}
	return v
}

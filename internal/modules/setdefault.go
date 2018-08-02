package modules

type S interface {
	SetUpdateEvery(int)
	SetChartsCleanup(int)
	SetDisabledByDefault()
}

type G interface {
	UpdateEvery() (int, bool)
	ChartsCleanup() (int, bool)
	DisabledByDefault() bool
}

type moduleDefault struct {
	updateEvery       *int
	chartsCleanup     *int
	disabledByDefault bool
}

func (m *moduleDefault) SetUpdateEvery(v int) {
	if m.updateEvery == nil {
		m.updateEvery = new(int)
	}
	*m.updateEvery = v
}

func (m *moduleDefault) SetChartsCleanup(v int) {
	if m.chartsCleanup == nil {
		m.chartsCleanup = new(int)
	}
	*m.chartsCleanup = v
}

func (m *moduleDefault) SetDisabledByDefault() {
	m.disabledByDefault = true
}

func (m *moduleDefault) UpdateEvery() (int, bool) {
	if m.updateEvery == nil {
		return 0, false
	}
	return *m.updateEvery, true
}

func (m *moduleDefault) ChartsCleanup() (int, bool) {
	if m.chartsCleanup == nil {
		return 0, false
	}
	return *m.chartsCleanup, true
}

func (m *moduleDefault) DisabledByDefault() bool {
	return m.disabledByDefault
}

var moduleDefaults = map[string]*moduleDefault{"_": {}}

func SetDefault() S {
	name := getFileName(2)
	if _, ok := moduleDefaults[name]; !ok {
		moduleDefaults[name] = new(moduleDefault)
	}
	return moduleDefaults[name]
}

func GetDefault(moduleName string) G {
	v, ok := moduleDefaults[moduleName]
	if !ok {
		return moduleDefaults["_"]
	}
	return v
}

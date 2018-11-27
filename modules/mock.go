package modules

type MockModule struct {
	Base
	InitFunc      func() bool
	CheckFunc     func() bool
	GetChartsFunc func() *Charts
	GetDataDunc   func() map[string]int64
	CleanupDone   bool
}

func (m MockModule) Init() bool {
	return m.InitFunc()
}

func (m MockModule) Check() bool {
	return m.CheckFunc()
}

func (m MockModule) GetCharts() *Charts {
	return m.GetChartsFunc()
}

func (m MockModule) GetData() map[string]int64 {
	return m.GetDataDunc()
}

func (m *MockModule) Cleanup() {
	m.CleanupDone = true
}

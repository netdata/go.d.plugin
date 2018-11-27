package modules

// MockModule MockModule
type MockModule struct {
	Base
	InitFunc      func() bool
	CheckFunc     func() bool
	GetChartsFunc func() *Charts
	GetDataDunc   func() map[string]int64
	CleanupDone   bool
}

// Init invokes InitFunc.
func (m MockModule) Init() bool {
	return m.InitFunc()
}

// Check invokes CheckFunc.
func (m MockModule) Check() bool {
	return m.CheckFunc()
}

// GetCharts invokes GetChartsFunc.
func (m MockModule) GetCharts() *Charts {
	return m.GetChartsFunc()
}

// GetData invokes GetDataDunc.
func (m MockModule) GetData() map[string]int64 {
	return m.GetDataDunc()
}

// Cleanup sets CleanupDone to true.
func (m *MockModule) Cleanup() {
	m.CleanupDone = true
}

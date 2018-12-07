package modules

// MockModule MockModule
type MockModule struct {
	Base

	InitFunc          func() bool
	CheckFunc         func() bool
	ChartsFunc        func() *Charts
	GatherMetricsFunc func() map[string]int64
	CleanupDone       bool
}

// Init invokes InitFunc
func (m MockModule) Init() bool {
	if m.InitFunc == nil {
		return true
	}
	return m.InitFunc()
}

// Check invokes CheckFunc
func (m MockModule) Check() bool {
	if m.CheckFunc == nil {
		return true
	}
	return m.CheckFunc()
}

// Charts invokes ChartsFunc
func (m MockModule) Charts() *Charts {
	if m.ChartsFunc == nil {
		return nil
	}
	return m.ChartsFunc()
}

// GatherMetrics invokes GetDataDunc
func (m MockModule) GatherMetrics() map[string]int64 {
	if m.GatherMetricsFunc == nil {
		return nil
	}
	return m.GatherMetricsFunc()
}

// Cleanup sets CleanupDone to true
func (m *MockModule) Cleanup() {
	m.CleanupDone = true
}

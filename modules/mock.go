package modules

// mockModule mockModule
type mockModule struct {
	Base
	initFunc          func() bool
	checkFunc         func() bool
	chartsFunc        func() *Charts
	gatherMetricsFunc func() map[string]int64
	cleanupDone       bool
}

// Init invokes initFunc.
func (m mockModule) Init() bool {
	return m.initFunc()
}

// Check invokes checkFunc.
func (m mockModule) Check() bool {
	return m.checkFunc()
}

// Charts invokes chartsFunc.
func (m mockModule) Charts() *Charts {
	return m.chartsFunc()
}

// GatherMetrics invokes GetDataDunc.
func (m mockModule) GatherMetrics() map[string]int64 {
	return m.gatherMetricsFunc()
}

// Cleanup sets cleanupDone to true.
func (m *mockModule) Cleanup() {
	m.cleanupDone = true
}

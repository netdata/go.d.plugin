package modules

// mockModule mockModule
type mockModule struct {
	Base
	init          func() bool
	check         func() bool
	charts        func() *Charts
	gatherMetrics func() map[string]int64
	cleanupDone   bool
}

// Init invokes init.
func (m mockModule) Init() bool {
	return m.init()
}

// Check invokes check.
func (m mockModule) Check() bool {
	return m.check()
}

// Charts invokes charts.
func (m mockModule) Charts() *Charts {
	return m.charts()
}

// GatherMetrics invokes GetDataDunc.
func (m mockModule) GatherMetrics() map[string]int64 {
	return m.gatherMetrics()
}

// Cleanup sets cleanupDone to true.
func (m *mockModule) Cleanup() {
	m.cleanupDone = true
}

package godplugin

type mockJob struct {
	fullName   func() string
	moduleName func() string
	name       func() string

	autoDetectionRetry func() int

	panicked func() bool

	init      func() bool
	check     func() bool
	postCheck func() bool

	tick func(int)

	start func()
	stop  func()
}

// FullName invokes fullName
func (m mockJob) FullName() string {
	return m.fullName()
}

// ModuleName invokes moduleName
func (m mockJob) ModuleName() string {
	return m.moduleName()
}

// Name invokes name
func (m mockJob) Name() string {
	return m.name()
}

// AutoDetectionRetry invokes autoDetectionRetry
func (m mockJob) AutoDetectionRetry() int {
	return m.autoDetectionRetry()
}

// Panicked invokes panicked
func (m mockJob) Panicked() bool {
	return m.panicked()
}

// Init invokes init
func (m mockJob) Init() bool {
	return m.init()
}

// Check invokes check
func (m mockJob) Check() bool {
	return m.check()
}

// PostCheck invokes postCheck
func (m mockJob) PostCheck() bool {
	return m.postCheck()
}

// Tick invokes tick
func (m mockJob) Tick(clock int) {
	m.tick(clock)
}

// Start invokes start
func (m mockJob) Start() {
	m.start()
}

// Stop invokes stop
func (m mockJob) Stop() {
	m.stop()
}

package godplugin

type mockJob struct {
	fullName   func() string
	moduleName func() string
	name       func() string

	autoDetectionRetry func() int
	panicked           func() bool

	init      func() bool
	check     func() bool
	postCheck func() bool

	tick func(int)

	start func()
	stop  func()
}

// FullName returns mock job full name
func (m mockJob) FullName() string {
	if m.fullName == nil {
		return "mock"
	}
	return m.fullName()
}

// ModuleName returns mock job module name
func (m mockJob) ModuleName() string {
	if m.moduleName == nil {
		return "mock"
	}
	return m.moduleName()
}

// Name returns mock job name
func (m mockJob) Name() string {
	if m.name == nil {
		return "mock"
	}
	return m.name()
}

// AutoDetectionRetry returns mock job autoDetectionRetry
func (m mockJob) AutoDetectionRetry() int {
	if m.autoDetectionRetry == nil {
		return 0
	}
	return m.autoDetectionRetry()
}

// Panicked return whether the mock job is panicked
func (m mockJob) Panicked() bool {
	if m.panicked == nil {
		return false
	}
	return m.panicked()
}

// Init invokes mock job init
func (m mockJob) Init() bool {
	if m.init == nil {
		return true
	}
	return m.init()
}

// Check invokes mock job check
func (m mockJob) Check() bool {
	if m.check == nil {
		return true
	}
	return m.check()
}

// PostCheck invokes mock job postCheck
func (m mockJob) PostCheck() bool {
	if m.postCheck == nil {
		return true
	}
	return m.postCheck()
}

// Tick invokes mock job tick
func (m mockJob) Tick(clock int) {
	if m.tick == nil {
		return
	}
	m.tick(clock)
}

// Start invokes mock job start
func (m mockJob) Start() {
	if m.start == nil {
		return
	}
	m.start()
}

// Stop invokes mock job stop
func (m mockJob) Stop() {
	if m.stop == nil {
		return
	}
	m.stop()
}

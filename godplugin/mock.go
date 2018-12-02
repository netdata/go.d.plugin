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

func (m mockJob) FullName() string {
	return m.fullName()
}

func (m mockJob) ModuleName() string {
	return m.moduleName()
}

func (m mockJob) Name() string {
	return m.name()
}

func (m mockJob) AutoDetectionRetry() int {
	return m.autoDetectionRetry()
}

func (m mockJob) Panicked() bool {
	return m.panicked()
}

func (m mockJob) Init() bool {
	return m.init()
}

func (m mockJob) Check() bool {
	return m.check()
}

func (m mockJob) PostCheck() bool {
	return m.postCheck()
}

func (m mockJob) Tick(clock int) {

}

func (m mockJob) Start() {

}

func (m mockJob) Stop() {

}

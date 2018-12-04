package godplugin

const (
	mockFullName           = "mock"
	mockModuleName         = "mock"
	mockName               = "mock"
	mockAutoDetectionRetry = 0
	mockPanicked           = false
	mockInit               = true
	mockCheck              = true
	mockPostCheck          = true
)

type mockJob struct{}

// FullName returns "mock"
func (mockJob) FullName() string {
	return mockFullName
}

// ModuleName returns "mock"
func (mockJob) ModuleName() string {
	return mockModuleName
}

// Name returns "mock"
func (mockJob) Name() string {
	return mockName
}

// AutoDetectionRetry returns 0
func (mockJob) AutoDetectionRetry() int {
	return mockAutoDetectionRetry
}

// Panicked returns false
func (mockJob) Panicked() bool {
	return mockPanicked
}

// Init returns true
func (mockJob) Init() bool {
	return mockInit
}

// Check returns true
func (mockJob) Check() bool {
	return mockCheck
}

// PostCheck returns true
func (mockJob) PostCheck() bool {
	return mockPostCheck
}

// Tick does nothing
func (m mockJob) Tick(int) {}

// Start does nothing
func (m mockJob) Start() {}

// Stop does nothing
func (m mockJob) Stop() {}

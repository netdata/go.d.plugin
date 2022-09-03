// SPDX-License-Identifier: GPL-3.0-or-later

package job

type MockJob struct {
	FullNameFunc           func() string
	ModuleNameFunc         func() string
	NameFunc               func() string
	AutoDetectionFunc      func() bool
	AutoDetectionEveryFunc func() int
	RetryAutoDetectionFunc func() bool
	TickFunc               func(int)
	StartFunc              func()
	StopFunc               func()
}

// FullName returns mock job full name.
func (m MockJob) FullName() string {
	if m.FullNameFunc == nil {
		return "mock"
	}
	return m.FullNameFunc()
}

// ModuleName returns mock job module name.
func (m MockJob) ModuleName() string {
	if m.ModuleNameFunc == nil {
		return "mock"
	}
	return m.ModuleNameFunc()
}

// Name returns mock job name.
func (m MockJob) Name() string {
	if m.NameFunc == nil {
		return "mock"
	}
	return m.NameFunc()
}

// AutoDetectionEvery returns mock job AutoDetectionEvery.
func (m MockJob) AutoDetectionEvery() int {
	if m.AutoDetectionEveryFunc == nil {
		return 0
	}
	return m.AutoDetectionEveryFunc()
}

// AutoDetection returns mock job AutoDetection.
func (m MockJob) AutoDetection() bool {
	if m.AutoDetectionFunc == nil {
		return true
	}
	return m.AutoDetectionFunc()
}

// RetryAutoDetection invokes mock job RetryAutoDetection.
func (m MockJob) RetryAutoDetection() bool {
	if m.RetryAutoDetectionFunc == nil {
		return true
	}
	return m.RetryAutoDetectionFunc()
}

// Tick invokes mock job Tick.
func (m MockJob) Tick(clock int) {
	if m.TickFunc != nil {
		m.TickFunc(clock)
	}
}

// Start invokes mock job Start.
func (m MockJob) Start() {
	if m.StartFunc != nil {
		m.StartFunc()
	}
}

// Stop invokes mock job Stop.
func (m MockJob) Stop() {
	if m.StopFunc != nil {
		m.StopFunc()
	}
}

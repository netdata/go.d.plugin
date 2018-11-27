package modules

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockModule struct {
	Base
	initFunc      func() bool
	checkFunc     func() bool
	getChartsFunc func() *Charts
	getDataDunc   func() map[string]int64
	cleanupDone   bool
}

func (m mockModule) Init() bool {
	return m.initFunc()
}

func (m mockModule) Check() bool {
	return m.checkFunc()
}

func (m mockModule) GetCharts() *Charts {
	return m.getChartsFunc()
}

func (m mockModule) GetData() map[string]int64 {
	return m.getDataDunc()
}

func (m *mockModule) Cleanup() {
	m.cleanupDone = true
}

func TestNewJob(t *testing.T) {
	assert.IsType(
		t,
		(*Job)(nil),
		NewJob("example", nil, ioutil.Discard, nil),
	)
}

func TestJob_FullName(t *testing.T) {
	job := NewJob("modName", &mockModule{}, ioutil.Discard, nil)

	assert.Equal(t, job.FullName(), "modName_modName")
	job.Nam = "jobName"
	assert.Equal(t, job.FullName(), "modName_jobName")

}

func TestJob_ModuleName(t *testing.T) {
	job := NewJob("modName", &mockModule{}, ioutil.Discard, nil)

	assert.Equal(t, job.ModuleName(), "modName")
}

func TestJob_Name(t *testing.T) {
	job := NewJob("modName", &mockModule{}, ioutil.Discard, nil)

	assert.Equal(t, job.Name(), "modName")
	job.Nam = "jobName"
	assert.Equal(t, job.Name(), "jobName")
}

func TestJob_Initialized(t *testing.T) {
	job := NewJob("modName", &mockModule{}, ioutil.Discard, nil)

	assert.Equal(t, job.Initialized(), job.initialized)
	job.initialized = true
	assert.Equal(t, job.Initialized(), job.initialized)

}

func TestJob_Panicked(t *testing.T) {
	job := NewJob("modName", &mockModule{}, ioutil.Discard, nil)

	assert.Equal(t, job.Panicked(), job.panicked)
	job.panicked = true
	assert.Equal(t, job.Panicked(), job.panicked)

}

func TestJob_AutoDetectionRetry(t *testing.T) {
	job := NewJob("modName", &mockModule{}, ioutil.Discard, nil)

	assert.Equal(t, job.AutoDetectionRetry(), job.AutoDetectRetry)
	job.AutoDetectRetry = 1
	assert.Equal(t, job.AutoDetectionRetry(), job.AutoDetectRetry)

}

func TestJob_Init(t *testing.T) {
	okMockModule := &mockModule{
		initFunc: func() bool { return true },
	}
	job := NewJob("modName", okMockModule, ioutil.Discard, nil)
	assert.True(t, job.Init())
	assert.True(t, job.Initialized())
	assert.False(t, job.Panicked())
	assert.False(t, okMockModule.cleanupDone)

	panicMockModule := &mockModule{
		initFunc: func() bool { panic("panic in init") },
	}
	job = NewJob("modName", panicMockModule, ioutil.Discard, nil)
	assert.False(t, job.Init())
	assert.False(t, job.Initialized())
	assert.True(t, job.Panicked())
	assert.True(t, panicMockModule.cleanupDone)
}

func TestJob_Check(t *testing.T) {
	okMockModule := &mockModule{
		checkFunc: func() bool { return true },
	}
	job := NewJob("modName", okMockModule, ioutil.Discard, nil)
	assert.True(t, job.Check())
	assert.False(t, job.Panicked())
	assert.False(t, okMockModule.cleanupDone)

	panicMockModule := &mockModule{
		checkFunc: func() bool { panic("panic in check") },
	}
	job = NewJob("modName", panicMockModule, ioutil.Discard, nil)
	assert.False(t, job.Check())
	assert.True(t, job.Panicked())
	assert.True(t, panicMockModule.cleanupDone)
}

func TestJob_PostCheck(t *testing.T) {
	okMockModule := &mockModule{
		getChartsFunc: func() *Charts { return &Charts{} },
	}
	job := NewJob("modName", okMockModule, ioutil.Discard, nil)
	assert.True(t, job.PostCheck())

	ngMockModule := &mockModule{
		getChartsFunc: func() *Charts { return nil },
	}
	job = NewJob("modName", ngMockModule, ioutil.Discard, nil)
	assert.False(t, job.PostCheck())
}

func TestJob_Start(t *testing.T) {

}

func TestJob_Stop(t *testing.T) {

}

func TestJob_Tick(t *testing.T) {

}

func TestJob_MainLoop(t *testing.T) {

}

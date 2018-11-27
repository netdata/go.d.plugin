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

func TestNewJob(t *testing.T) {
	assert.IsType(
		t,
		(*Job)(nil),
		NewJob("example", nil, ioutil.Discard, nil),
	)
}

func TestJob_FullName(t *testing.T) {
	job := NewJob(
		"modName",
		&mockModule{},
		ioutil.Discard,
		nil,
	)
	assert.Equal(t, job.FullName(), "modName_modName")
	job.Nam = "jobName"
	assert.Equal(t, job.FullName(), "modName_jobName")

}

func TestJob_ModuleName(t *testing.T) {
	job := NewJob(
		"modName",
		&mockModule{},
		ioutil.Discard,
		nil,
	)
	assert.Equal(t, job.ModuleName(), "modName")
}

func TestJob_Name(t *testing.T) {
	job := NewJob(
		"modName",
		&mockModule{},
		ioutil.Discard,
		nil,
	)
	assert.Equal(t, job.Name(), "modName")
	job.Nam = "jobName"
	assert.Equal(t, job.Name(), "jobName")
}

func TestJob_Initialized(t *testing.T) {

}

func TestJob_Panicked(t *testing.T) {

}

func TestJob_AutoDetectionRetry(t *testing.T) {

}

func TestJob_Init(t *testing.T) {

}

func TestJob_Check(t *testing.T) {

}

func TestJob_PostCheck(t *testing.T) {

}

func TestJob_Start(t *testing.T) {

}

func TestJob_Stop(t *testing.T) {

}

func TestJob_Tick(t *testing.T) {

}

func TestJob_MainLoop(t *testing.T) {

}

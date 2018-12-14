package modules

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/netdata/go.d.plugin/logger"

	"github.com/stretchr/testify/assert"
)

var (
	testModName = "modName"
	testJobName = "jobName"
)

func testNewJob() *Job {
	return NewJob(testModName, nil, ioutil.Discard, nil)
}

func TestNewJob(t *testing.T) {
	assert.IsType(t, (*Job)(nil), testNewJob())
}

func TestJob_FullName(t *testing.T) {
	job := testNewJob()

	assert.Equal(t, job.FullName(), testModName)
	job.Nam = testModName
	assert.Equal(t, job.FullName(), testModName)
	job.Nam = testJobName
	assert.Equal(t, job.FullName(), fmt.Sprintf("%s_%s", testModName, testJobName))

}

func TestJob_ModuleName(t *testing.T) {
	job := testNewJob()

	assert.Equal(t, job.ModuleName(), testModName)
}

func TestJob_Name(t *testing.T) {
	job := testNewJob()

	assert.Equal(t, job.Name(), testModName)
	job.Nam = testJobName
	assert.Equal(t, job.Name(), testJobName)
}

func TestJob_Panicked(t *testing.T) {
	job := testNewJob()

	assert.Equal(t, job.Panicked(), job.panicked)
	job.panicked = true
	assert.Equal(t, job.Panicked(), job.panicked)

}

func TestJob_AutoDetectionRetry(t *testing.T) {
	job := testNewJob()

	assert.Equal(t, job.AutoDetectionRetry(), job.AutoDetectRetry)
	job.AutoDetectRetry = 1
	assert.Equal(t, job.AutoDetectionRetry(), job.AutoDetectRetry)

}

func TestJob_Init(t *testing.T) {
	// OK case
	m := &MockModule{
		InitFunc: func() bool { return true },
	}
	job := testNewJob()
	job.module = m

	assert.True(t, job.Init())
	assert.True(t, job.initialized)
	assert.False(t, job.Panicked())
	assert.False(t, m.CleanupDone)

	// NG case
	m = &MockModule{
		InitFunc: func() bool { return false },
	}
	job = testNewJob()
	job.module = m

	assert.False(t, job.Init())
	assert.False(t, job.initialized)
	assert.False(t, job.Panicked())
	assert.False(t, m.CleanupDone)

	// PANIC case
	m = &MockModule{
		InitFunc: func() bool { panic("panic in InitFunc") },
	}
	job = testNewJob()
	job.module = m

	assert.False(t, job.Init())
	assert.False(t, job.initialized)
	assert.True(t, job.Panicked())
	assert.True(t, m.CleanupDone)
}

func TestJob_Check(t *testing.T) {
	// OK case
	m := &MockModule{
		CheckFunc: func() bool { return true },
	}
	job := testNewJob()
	job.module = m

	assert.True(t, job.Check())
	assert.False(t, job.Panicked())
	assert.False(t, m.CleanupDone)

	// NG case
	m = &MockModule{
		CheckFunc: func() bool { return false },
	}
	job = testNewJob()
	job.module = m

	assert.False(t, job.Check())
	assert.False(t, job.Panicked())
	assert.False(t, m.CleanupDone)

	// PANIC case
	m = &MockModule{
		CheckFunc: func() bool { panic("panic in InitFunc") },
	}
	job = testNewJob()
	job.module = m

	assert.False(t, job.Check())
	assert.False(t, job.initialized)
	assert.True(t, job.Panicked())
	assert.True(t, m.CleanupDone)
}

func TestJob_PostCheck(t *testing.T) {
	// OK case
	m := &MockModule{
		ChartsFunc: func() *Charts { return &Charts{} },
	}
	job := testNewJob()
	job.module = m

	assert.True(t, job.PostCheck())

	// NG case
	m = &MockModule{
		ChartsFunc: func() *Charts { return nil },
	}
	job = testNewJob()
	job.module = m

	assert.False(t, job.PostCheck())
}

func TestJob_MainLoop(t *testing.T) {
	m := &MockModule{
		ChartsFunc: func() *Charts {
			return &Charts{
				&Chart{
					ID:    "id",
					Title: "title",
					Units: "units",
					Dims: Dims{
						{ID: "id1"},
						{ID: "id2"},
					},
				},
			}
		},
		CollectFunc: func() map[string]int64 {
			return map[string]int64{
				"id1": 1,
				"id2": 2,
			}
		},
	}
	job := testNewJob()
	job.module = m
	job.charts = job.module.Charts()
	job.UpdateEvery = 1

	go func() {
		for i := 1; i < 3; i++ {
			job.Tick(i)
			time.Sleep(time.Second)
		}
		job.Stop()
	}()

	job.MainLoop()

	assert.True(t, m.CleanupDone)
}

func TestJob_MainLoop_Panic(t *testing.T) {
	m := &MockModule{
		CollectFunc: func() map[string]int64 {
			panic("panic in Collect")
		},
	}
	job := testNewJob()
	job.module = m
	obs := testObserver(false)
	job.observer = &obs
	job.UpdateEvery = 1

	go func() {
		for i := 1; i < 3; i++ {
			time.Sleep(time.Second)
			job.Tick(i)
		}
		job.Stop()
	}()

	job.MainLoop()

	assert.True(t, job.Panicked())
	assert.True(t, bool(*job.observer.(*testObserver)))
}

func TestJob_Tick(t *testing.T) {
	job := NewJob(testModName, nil, ioutil.Discard, nil)
	for i := 0; i < 3; i++ {
		job.Tick(i)
	}
}

type testObserver bool

func (t *testObserver) RemoveFromQueue(name string) {
	*t = true
}

func TestJob_Start(t *testing.T) {
	job := testNewJob()
	job.module = &MockModule{}

	go func() {
		time.Sleep(time.Second)
		job.Stop()
	}()

	job.Start()

}

func TestBase_SetLogger(t *testing.T) {
	var b Base
	b.SetLogger(&logger.Logger{})

	assert.NotNil(t, b.Logger)
}

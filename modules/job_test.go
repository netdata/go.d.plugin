package modules

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"time"
)

var (
	testModName = "testModName"
	testJobName = "testJobName"
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

func TestJob_Initialized(t *testing.T) {
	job := testNewJob()

	assert.Equal(t, job.Initialized(), job.initialized)
	job.initialized = true
	assert.Equal(t, job.Initialized(), job.initialized)

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
	okMockModule := &mockModule{
		initFunc: func() bool { return true },
	}
	job := testNewJob()
	job.module = okMockModule

	assert.True(t, job.Init())
	assert.True(t, job.Initialized())
	assert.False(t, job.Panicked())
	assert.False(t, okMockModule.cleanupDone)

	panicMockModule := &mockModule{
		initFunc: func() bool { panic("panic in init") },
	}
	job = testNewJob()
	job.module = panicMockModule

	assert.False(t, job.Init())
	assert.False(t, job.Initialized())
	assert.True(t, job.Panicked())
	assert.True(t, panicMockModule.cleanupDone)
}

func TestJob_Check(t *testing.T) {
	okMockModule := &mockModule{
		checkFunc: func() bool { return true },
	}
	job := testNewJob()
	job.module = okMockModule

	assert.True(t, job.Check())
	assert.False(t, job.Panicked())
	assert.False(t, okMockModule.cleanupDone)

	panicMockModule := &mockModule{
		checkFunc: func() bool { panic("panic in check") },
	}
	job = testNewJob()
	job.module = panicMockModule

	assert.False(t, job.Check())
	assert.True(t, job.Panicked())
	assert.True(t, panicMockModule.cleanupDone)
}

func TestJob_PostCheck(t *testing.T) {
	okMockModule := &mockModule{
		chartsFunc: func() *Charts { return &Charts{} },
	}
	job := testNewJob()
	job.module = okMockModule

	assert.True(t, job.PostCheck())

	ngMockModule := &mockModule{
		chartsFunc: func() *Charts { return nil },
	}
	job = testNewJob()
	job.module = ngMockModule

	assert.False(t, job.PostCheck())
}

func TestJob_MainLoop(t *testing.T) {
	module := &mockModule{
		chartsFunc: func() *Charts {
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
		gatherMetricsFunc: func() map[string]int64 {
			return map[string]int64{
				"id1": 1,
				"id2": 2,
			}
		},
	}
	job := testNewJob()
	job.module = module
	job.charts = job.module.Charts()
	job.UpdateEvery = 1

	go func() {
		for i := 1; i < 4; i++ {
			job.Tick(i)
			time.Sleep(time.Second)
		}
		job.Stop()
	}()

	job.MainLoop()
}

func TestJob_MainLoop_Panic(t *testing.T) {
	module := &mockModule{
		gatherMetricsFunc: func() map[string]int64 {
			panic("panic in GatherMetrics")
		},
	}
	job := testNewJob()
	job.module = module
	obs := testObserver(false)
	job.observer = &obs
	job.UpdateEvery = 1

	go func() {
		for i := 1; i < 4; i++ {
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

package godplugin

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_initJob(t *testing.T) {
	p := New()
	job := &mockJob{}

	// OK case
	job.init = func() bool { return true }
	assert.True(t, p.initJob(job))

	// NG case
	job.init = func() bool { return false }
	assert.False(t, p.initJob(job))
}

func Test_jobCheck(t *testing.T) {
	p := New()

	job := &mockJob{}

	// OK case
	job.check = func() bool { return true }
	assert.True(t, p.checkJob(job))

	// NG case
	job.check = func() bool { return false }
	assert.False(t, p.checkJob(job))

	// Panic case
	job.check = func() bool { return true }
	job.panicked = func() bool { return true }
	assert.False(t, p.checkJob(job))

	// AutoDetectionRetry case
	job.check = func() bool { return false }
	job.panicked = func() bool { return false }
	job.autoDetectionRetry = func() int { return 1 }
	assert.False(t, p.checkJob(job))

	wait := time.NewTimer(time.Second * 2)
	defer wait.Stop()

	select {
	case <-wait.C:
		t.Error("auto detection retry test failed")
	case <-p.jobCh:
	}
}

func Test_jobPostCheck(t *testing.T) {
	p := New()

	job := &mockJob{}

	// OK case
	job.postCheck = func() bool { return true }
	assert.True(t, p.postCheckJob(job))

	// NG case
	job.postCheck = func() bool { return false }
	assert.False(t, p.postCheckJob(job))
}

func Test_jobStartLoop(t *testing.T) {
	p := New()

	go p.jobStartLoop()

	job := &mockJob{}

	p.jobCh <- job
	p.jobCh <- job
	p.jobCh <- job
	p.jobStartShutdown <- struct{}{}

	assert.Equal(t, 1, len(p.loopQueue.queue))

	for _, j := range p.loopQueue.queue {
		j.Stop()
	}
}

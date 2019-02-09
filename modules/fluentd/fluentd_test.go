package fluentd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.IsType(t, (*Fluentd)(nil), New())
}

func TestFluentd_Init(t *testing.T) {
	job := New()

	assert.True(t, job.Init())
}

func TestFluentd_Check(t *testing.T) {
	job := New()

	assert.True(t, job.Check())
}

func TestFluentd_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestFluentd_Cleanup(t *testing.T) {
	New().Cleanup()
}

func TestFluentd_Collect(t *testing.T) {
	job := New()

	assert.NotNil(t, job.Collect())
}

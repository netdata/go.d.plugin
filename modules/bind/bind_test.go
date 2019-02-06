package bind

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.IsType(t, (*Bind)(nil), New())
}

func TestBind_Init(t *testing.T) {
	job := New()

	assert.True(t, job.Init())
}

func TestBind_Check(t *testing.T) {
	job := New()

	assert.True(t, job.Check())
}

func TestBind_Charts(t *testing.T) {
	job := New()

	assert.NotNil(t, job.Charts())
}

func TestBind_Cleanup(t *testing.T) {
	job := New()

	job.Cleanup()
}

func TestBind_Collect(t *testing.T) {
	job := New()

	assert.NotNil(t, job.Collect())
}

package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) { assert.IsType(t, (*Kubernetes)(nil), New()) }

func TestKubernetes_Init(t *testing.T) {
	job := New()
	assert.True(t, job.Init())
}

func TestKubernetes_Check(t *testing.T) {
	job := New()
	assert.True(t, job.Check())
}

func TestKubernetes_Charts(t *testing.T) {
	job := New()
	assert.NotNil(t, job.Charts())
}

func TestKubernetes_Cleanup(t *testing.T) {
	job := New()
	job.Cleanup()
}

func TestKubernetes_Collect(t *testing.T) {
	mod := New()
	assert.NotNil(t, mod.Collect())
}

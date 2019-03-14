package x509check

import (
	"github.com/netdata/go-orchestrator/module"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	job := New()

	assert.Implements(t, (*module.Module)(nil), job)
	assert.Equal(t, defaultConnTimeout, job.Config.Timeout.Duration)
	assert.Equal(t, defaultDaysUntilWarn, job.Config.DaysUntilWarn)
	assert.Equal(t, defaultDaysUntilCrit, job.Config.DaysUntilCrit)
}

func TestX509Check_Cleanup(t *testing.T) { New().Cleanup() }

func TestX509Check_Charts(t *testing.T) {
	job := New()

	assert.NotNil(t, job.Charts())

}

func TestX509Check_Init(t *testing.T) {
	job := New()
	assert.False(t, job.Init())

	job = New()
	job.Source = "wrong"
	assert.False(t, job.Init())

	job = New()
	job.Source = "https://example.org:443"
	assert.True(t, job.Init())
}

// TODO:
func TestX509Check_Check(t *testing.T) {}

func TestX509Check_Collect(t *testing.T) {}

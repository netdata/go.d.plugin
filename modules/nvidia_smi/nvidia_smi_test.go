// +build linux

package nvidia_smi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.IsType(t, (*Nvsmi)(nil), New())
}

func TestNvsmi_Init(t *testing.T) {
	mod := New()

	assert.True(t, mod.Init())
}

func TestNvsmi_Check(t *testing.T) {
	mod := New()

	assert.True(t, mod.Check())
}

func TestNvsmi_Charts(t *testing.T) {
	mod := New()

	assert.NotNil(t, mod.Charts())
}

func TestNvsmi_Cleanup(t *testing.T) {
	mod := New()

	mod.Cleanup()
}

func TestNvsmi_Collect(t *testing.T) {
	mod := New()

	assert.NotNil(t, mod.Collect())
}

package chrony

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.IsType(t, (*chronyCollector)(nil), New())
}

func TestChrony_Init(t *testing.T) {
	mod := New()

	assert.True(t, mod.Init())
}

func TestChrony_Check(t *testing.T) {
	mod := New()

	assert.True(t, mod.Check())
}

func TestChrony_Charts(t *testing.T) {
	mod := New()

	assert.NotNil(t, mod.Charts())
}

func TestChrony_Cleanup(t *testing.T) {
	mod := New()

	mod.Cleanup()
}

func TestChrony_Collect(t *testing.T) {
	mod := New()

	assert.NotNil(t, mod.Collect())
}

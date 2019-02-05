package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.IsType(t, (*Example)(nil), New())
}

func TestExample_Init(t *testing.T) {
	mod := New()

	assert.True(t, mod.Init())
}

func TestExample_Check(t *testing.T) {
	mod := New()

	assert.True(t, mod.Check())
}

func TestExample_Charts(t *testing.T) {
	mod := New()

	assert.NotNil(t, mod.Charts())
}

func TestExample_Cleanup(t *testing.T) {
	mod := New()

	mod.Cleanup()
}

func TestExample_Collect(t *testing.T) {
	mod := New()

	assert.NotNil(t, mod.Collect())
}

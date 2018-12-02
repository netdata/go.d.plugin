package modules

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	modName := "modName"
	registry := make(Registry)

	// add "modName" to the register
	assert.NotPanics(
		t,
		func() {
			register(registry, modName, Creator{})
		})

	_, exist := registry[modName]

	require.True(t, exist)

	// re-add "modName" to the register
	assert.Panics(
		t,
		func() {
			register(registry, modName, Creator{})
		})

}

package modules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	modName := "modName"
	registry := make(Registry)

	assert.NotPanics(t, func() {
		register(registry, modName, Creator{})
	})

	_, ok := registry[modName]

	assert.True(t, ok)

	assert.Panics(t, func() {
		register(registry, modName, Creator{})
	})

}

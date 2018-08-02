package godplugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	goDConfFull    = "tests/go.d.conf-full.yml"
	goDConfMinimal = "tests/go.d.conf-minimal.yml"
	goDConfNeg     = "tests/go.d.conf-neg.yml"
)

func TestConfigFull(t *testing.T) {
	config := NewConfig()
	config.Load(goDConfFull)

	assert.True(t, config.Enabled)
	assert.True(t, config.DefaultRun)
	assert.Equal(t, 10, config.MaxProcs)
	assert.False(t, config.IsModuleEnabled("example", false))
	assert.True(t, config.IsModuleEnabled("foo", false))
	assert.True(t, config.IsModuleEnabled("bar", false))

	assert.False(t, config.IsModuleEnabled("example", true))
	assert.True(t, config.IsModuleEnabled("foo", true))
	assert.False(t, config.IsModuleEnabled("bar", true))
}

func TestConfigMinimal(t *testing.T) {
	config := NewConfig()
	config.Load(goDConfMinimal)

	assert.True(t, config.Enabled)
	assert.True(t, config.DefaultRun)
	assert.Equal(t, 1, config.MaxProcs)
	assert.True(t, config.IsModuleEnabled("example", false))
	assert.True(t, config.IsModuleEnabled("foo", false))
	assert.True(t, config.IsModuleEnabled("bar", false))

	assert.False(t, config.IsModuleEnabled("example", true))
	assert.False(t, config.IsModuleEnabled("foo", true))
	assert.False(t, config.IsModuleEnabled("bar", true))
}

func TestConfigNeg(t *testing.T) {
	config := NewConfig()
	config.Load(goDConfNeg)

	assert.True(t, config.Enabled)
	assert.True(t, config.DefaultRun)
	assert.Equal(t, 1, config.MaxProcs)
	assert.True(t, config.IsModuleEnabled("example", false))
	assert.True(t, config.IsModuleEnabled("foo", false))
	assert.True(t, config.IsModuleEnabled("bar", false))

	assert.False(t, config.IsModuleEnabled("example", true))
	assert.True(t, config.IsModuleEnabled("foo", true))
	assert.False(t, config.IsModuleEnabled("bar", true))
}

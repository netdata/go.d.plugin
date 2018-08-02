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
	assert.False(t, config.IsModuleEnabled("example"))
	assert.True(t, config.IsModuleEnabled("foo"))
	assert.True(t, config.IsModuleEnabled("bar"))
}

func TestConfigMinimal(t *testing.T) {
	config := NewConfig()
	config.Load(goDConfMinimal)

	assert.True(t, config.Enabled)
	assert.True(t, config.DefaultRun)
	assert.Equal(t, 1, config.MaxProcs)
	assert.True(t, config.IsModuleEnabled("example"))
	assert.True(t, config.IsModuleEnabled("foo"))
	assert.True(t, config.IsModuleEnabled("bar"))
}

func TestConfigNeg(t *testing.T) {
	config := NewConfig()
	config.Load(goDConfNeg)

	assert.True(t, config.Enabled)
	assert.True(t, config.DefaultRun)
	assert.Equal(t, 1, config.MaxProcs)
	assert.True(t, config.IsModuleEnabled("example"))
	assert.True(t, config.IsModuleEnabled("foo"))
	assert.True(t, config.IsModuleEnabled("bar"))
}

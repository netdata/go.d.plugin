package godplugin

import (
	"github.com/stretchr/testify/require"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	goDConfNotExist = "tests/go.d.conf-not-exist.yml"
	goDConf         = "tests/go.d.conf.yml"
	goDConfNeg      = "tests/go.d.conf-neg.yml"
	goDConfEmpty    = "tests/go.d.conf-empty.yml"

	foo = "foo"
	bar = "bar"
	baz = "baz"
)

func TestConfig_Load(t *testing.T) {
	config := NewConfig()

	require.Error(t, config.Load(goDConfNotExist))
	require.NoError(t, config.Load(goDConf))
}

func TestNewConfig(t *testing.T) {
	config := NewConfig()
	assert.IsType(t, config, (*Config)(nil))

	assert.Equal(
		t,
		config,
		&Config{
			Enabled:    true,
			DefaultRun: true,
		},
	)
}

func TestConfig_isModuleEnabled(t *testing.T) {
	config := Config{
		Modules: map[string]bool{
			foo: true,
			bar: false,
			// baz: true,
		},
	}

	config.DefaultRun = true

	assert.True(t, config.isModuleEnabled(foo, false))
	assert.True(t, config.isModuleEnabled(foo, true))

	assert.False(t, config.isModuleEnabled(bar, false))
	assert.False(t, config.isModuleEnabled(bar, true))

	assert.Equal(t, config.isModuleEnabled(baz, false), config.DefaultRun)
	assert.False(t, config.isModuleEnabled(baz, true))

	config.DefaultRun = false

	assert.True(t, config.isModuleEnabled(foo, false))
	assert.True(t, config.isModuleEnabled(foo, true))

	assert.False(t, config.isModuleEnabled(bar, false))
	assert.False(t, config.isModuleEnabled(bar, true))

	assert.Equal(t, config.isModuleEnabled(baz, false), config.DefaultRun)
	assert.False(t, config.isModuleEnabled(baz, true))
}

func TestConfigNone(t *testing.T) {
	config := NewConfig()
	err := config.Load(goDConfNotExist)
	assert.Error(t, err)
}

func TestConfigEmpty(t *testing.T) {
	config := NewConfig()
	err := config.Load(goDConfEmpty)
	assert.Equal(t, err, io.EOF)

	assert.True(t, config.Enabled)
	assert.True(t, config.DefaultRun)
	assert.Equal(t, 0, config.MaxProcs)

	assert.True(t, config.isModuleEnabled(foo, false))
	assert.True(t, config.isModuleEnabled(bar, false))
	assert.True(t, config.isModuleEnabled(baz, false))
}

func TestConfig(t *testing.T) {
	config := NewConfig()
	err := config.Load(goDConf)
	assert.NoError(t, err)

	assert.True(t, config.Enabled)
	assert.True(t, config.DefaultRun)
	assert.Equal(t, 10, config.MaxProcs)

	assert.True(t, config.isModuleEnabled(foo, false))
	assert.False(t, config.isModuleEnabled(bar, false))
	assert.True(t, config.isModuleEnabled(baz, false))
}

func TestConfigNeg(t *testing.T) {
	config := NewConfig()
	err := config.Load(goDConfNeg)
	assert.NoError(t, err)

	assert.True(t, config.Enabled)
	assert.False(t, config.DefaultRun)
	assert.Equal(t, 10, config.MaxProcs)

	assert.True(t, config.isModuleEnabled(foo, false))
	assert.False(t, config.isModuleEnabled(bar, false))
	assert.False(t, config.isModuleEnabled(baz, false))
}

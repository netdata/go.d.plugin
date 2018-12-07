package godplugin

import (
	"io/ioutil"
	"runtime"
	"testing"

	"github.com/netdata/go.d.plugin/cli"
	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/multipath"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlugin_SetupNoConfig(t *testing.T) {
	p := New()
	assert.False(t, p.Setup())
}

func TestPlugin_SetupBrokenConfig(t *testing.T) {
	p := New()

	p.ConfigPath = multipath.New("./tests")
	p.confName = "go.d.conf-broken.yml"

	assert.False(t, p.Setup())
}

func TestPlugin_SetupEmptyConfig(t *testing.T) {
	p := New()

	p.Option, _ = cli.Parse([]string{})
	p.ConfigPath = multipath.New("./tests")
	p.confName = "go.d.conf-empty.yml"

	assert.True(t, p.Setup())
}

func TestPlugin_SetupDisabledInConfig(t *testing.T) {
	p := New()
	p.Out = ioutil.Discard

	p.ConfigPath = multipath.New("./tests")
	p.confName = "go.d.conf-disabled.yml"

	assert.False(t, p.Setup())
}

func TestPlugin_SetupNoModulesToRun(t *testing.T) {
	p := New()

	// registry is empty
	reg := make(modules.Registry)

	p.ConfigPath = multipath.New("./tests")
	p.Option = &cli.Option{Module: "all"}
	p.confName = "go.d.conf.yml"
	p.registry = reg

	assert.False(t, p.Setup())
}

func TestPlugin_SetupSetGOMAXPROCS(t *testing.T) {
	p := New()

	reg := make(modules.Registry)
	reg.Register("module1", modules.Creator{})
	reg.Register("module2", modules.Creator{})

	p.ConfigPath = multipath.New("./tests")
	p.Option = &cli.Option{Module: "all"}
	p.config.MaxProcs = 1
	p.confName = "go.d.conf.yml"
	p.registry = reg

	assert.True(t, p.Setup())
	assert.Equal(t, p.config.MaxProcs, runtime.GOMAXPROCS(0))
}

func TestPlugin_Setup(t *testing.T) {
	p := New()

	reg := make(modules.Registry)
	reg.Register("module1", modules.Creator{})
	reg.Register("module2", modules.Creator{})

	p.ConfigPath = multipath.New("./tests")
	p.Option = &cli.Option{Module: "all"}
	p.confName = "go.d.conf.yml"
	p.registry = reg

	assert.True(t, p.Setup())
	assert.Equal(t, 2, len(p.modules))
}

func TestPlugin_populateActiveModulesAll(t *testing.T) {
	p := New()

	reg := make(modules.Registry)
	p.registry = reg

	reg.Register("module1", modules.Creator{})
	reg.Register("module2", modules.Creator{})

	require.Len(t, p.modules, 0)

	p.Option = &cli.Option{Module: "all"}
	p.populateActiveModules()

	require.Len(t, p.modules, 2)
}

func TestPlugin_populateActiveModulesWithDisabledByDefault(t *testing.T) {
	p := New()

	reg := make(modules.Registry)
	p.registry = reg

	reg.Register("module1", modules.Creator{})
	reg.Register("module2", modules.Creator{DisabledByDefault: true})

	require.Len(t, p.modules, 0)

	p.Option = &cli.Option{Module: "all"}
	p.populateActiveModules()

	require.Len(t, p.modules, 1)
}

func TestPlugin_populateActiveModulesSpecific(t *testing.T) {
	p := New()

	reg := make(modules.Registry)
	p.registry = reg

	reg.Register("module1", modules.Creator{})
	reg.Register("module2", modules.Creator{})

	require.Len(t, p.modules, 0)

	p.Option = &cli.Option{Module: "module1"}
	p.populateActiveModules()

	require.Len(t, p.modules, 1)
}

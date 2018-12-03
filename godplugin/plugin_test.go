package godplugin

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/netdata/go.d.plugin/cli"
	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/multipath"
)

func TestNew(t *testing.T) {
	assert.IsType(t, (*Plugin)(nil), New())
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

func TestPlugin_SetupSpecificModule(t *testing.T) {
	p := New()

	reg := make(modules.Registry)
	reg.Register("module1", modules.Creator{})
	reg.Register("module2", modules.Creator{})

	p.ConfigPath = multipath.New("./tests")
	p.Option = &cli.Option{Module: "module1"}
	p.confName = "go.d.conf.yml"
	p.registry = reg

	assert.True(t, p.Setup())
	assert.Equal(t, 1, len(p.modules))
}

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

	p.ConfigPath = multipath.New("./tests")
	p.confName = "go.d.conf-empty.yml"

	assert.False(t, p.Setup())
}

func TestPlugin_SetupNoModulesToRun(t *testing.T) {
	p := New()

	reg := make(modules.Registry)

	p.ConfigPath = multipath.New("./tests")
	p.Option = &cli.Option{Module: "all"}
	p.confName = "go.d.conf.yml"
	p.registry = reg

	assert.False(t, p.Setup())
}

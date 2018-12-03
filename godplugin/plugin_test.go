package godplugin

import (
	"os"
	"testing"
	"time"

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

func TestPlugin_Serve(t *testing.T) {
	p := New()
	p.Out = os.Stdout

	module := &modules.MockModule{
		InitFunc:  func() bool { return true },
		CheckFunc: func() bool { return true },
		ChartsFunc: func() *modules.Charts {
			return &modules.Charts{
				&modules.Chart{
					ID:    "id",
					Title: "title",
					Units: "units",
					Dims: modules.Dims{
						{ID: "id1"},
						{ID: "id2"},
					},
				},
			}
		},
		GatherMetricsFunc: func() map[string]int64 {
			return map[string]int64{
				"id1": 1,
				"id2": 2,
			}
		},
	}

	reg := make(modules.Registry)
	reg.Register("module1", modules.Creator{Create: func() modules.Module { return module }})
	reg.Register("module2", modules.Creator{Create: func() modules.Module { return module }})

	p.ConfigPath = multipath.New("./tests")
	p.Option = &cli.Option{Module: "all"}
	p.confName = "go.d.conf.yml"
	p.registry = reg

	p.Setup()
	go p.Serve()

	time.Sleep(time.Second * 3)
	//close(p.checkCh)
}

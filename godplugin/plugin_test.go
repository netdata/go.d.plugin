package godplugin

import (
	"os"
	"testing"
	"time"

	"github.com/netdata/go.d.plugin/cli"
	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/multipath"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.IsType(t, (*Plugin)(nil), New())
}

func TestPlugin_Serve(t *testing.T) {
	p := New()
	p.Out = os.Stdout

	mod := func() modules.Module {
		return &modules.MockModule{
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
			CollectFunc: func() map[string]int64 {
				return map[string]int64{
					"id1": 1,
					"id2": 2,
				}
			},
		}
	}

	reg := make(modules.Registry)
	reg.Register("module1", modules.Creator{Create: func() modules.Module { return mod() }})
	reg.Register("module2", modules.Creator{Create: func() modules.Module { return mod() }})

	p.ConfigPath = multipath.New("./tests")
	p.Option = &cli.Option{Module: "all"}
	p.confName = "go.d.conf.yml"
	p.registry = reg

	p.Setup()

	go p.Serve()

	time.Sleep(time.Second * 3)

	p.stop()

	for _, job := range p.loopQueue.queue {
		job.Stop()
	}
}

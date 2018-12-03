package godplugin

import (
	"github.com/netdata/go.d.plugin/modules"
	"testing"

	"github.com/netdata/go.d.plugin/cli"
	"github.com/netdata/go.d.plugin/pkg/multipath"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.IsType(t, (*Plugin)(nil), New())
}

func TestPlugin_Setup(t *testing.T) {
	p := New()

	reg := make(modules.Registry)
	reg.Register("module", modules.Creator{})

	p.ConfigPath = multipath.New("./tests")
	p.Option = &cli.Option{Module: "all"}
	p.confName = "go.d.conf.yml"
	p.registry = reg

	p.Setup()
}

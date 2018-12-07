package godplugin

import (
	"testing"

	"github.com/netdata/go.d.plugin/cli"
	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/multipath"

	"github.com/stretchr/testify/assert"
)

func Test_loadModuleConfigNoConfig(t *testing.T) {
	p := New()

	p.ConfigPath = multipath.New("./tests")

	assert.Nil(t, p.loadModuleConfig("no config"))
}

func Test_loadModuleConfigBrokenConfig(t *testing.T) {
	p := New()

	p.ConfigPath = multipath.New("./tests")

	assert.Nil(t, p.loadModuleConfig("module-broken"))
}

func Test_loadModuleConfigNoJobs(t *testing.T) {
	p := New()

	p.ConfigPath = multipath.New("./tests")

	assert.Nil(t, p.loadModuleConfig("module-no-jobs"))
}

func Test_loadModuleConfig(t *testing.T) {
	p := New()

	p.ConfigPath = multipath.New("./tests")

	conf := p.loadModuleConfig("module1")

	assert.NotNil(t, conf)

	assert.Equal(t, 3, len(conf.Jobs))
}

func Test_createModuleJobs(t *testing.T) {
	p := New()

	p.ConfigPath = multipath.New("./tests")
	p.Option = &cli.Option{}

	reg := make(modules.Registry)
	reg.Register(
		"module1",
		modules.Creator{Create: func() modules.Module { return &modules.MockModule{} }},
	)

	p.registry = reg

	conf := &moduleConfig{Jobs: []map[string]interface{}{{}, {}, {}}}
	conf.name = "module1"
	assert.Len(t, p.createModuleJobs(conf), 3)
}

package godplugin

import (
	"testing"

	"bytes"

	"github.com/golang/mock/gomock"
	"github.com/l2isbad/go.d.plugin/internal/pkg/cli"
	"github.com/l2isbad/go.d.plugin/internal/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestPlugin(t *testing.T) {
	logger.SetLevel(logger.DEBUG)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	buf := &bytes.Buffer{}

	plg := NewPlugin()
	plg.Config = NewConfig()
	plg.Config.Load("tests/go.d.conf-example-only.yml")
	plg.Option = &cli.Option{
		Module:      "all",
		UpdateEvery: 1,
	}
	plg.Out = buf

	assert.True(t, plg.Setup())
	plg.CheckJobs()
}

func TestPlugin_Disabled(t *testing.T) {
	buf := &bytes.Buffer{}
	plg := NewPlugin()
	plg.Config = NewConfig()
	plg.Config.Enabled = false
	plg.Out = buf
	assert.False(t, plg.Setup())
	assert.Equal(t, "DISABLE\n", buf.String())
}

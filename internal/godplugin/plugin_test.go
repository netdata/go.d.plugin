package godplugin

import (
	"testing"

	"bytes"

	"time"

	"github.com/golang/mock/gomock"
	"github.com/l2isbad/go.d.plugin/internal/mock"
	"github.com/l2isbad/go.d.plugin/internal/modules"
	"github.com/l2isbad/go.d.plugin/internal/pkg/charts"
	. "github.com/l2isbad/go.d.plugin/internal/pkg/charts"
	"github.com/l2isbad/go.d.plugin/internal/pkg/cli"
	"github.com/l2isbad/go.d.plugin/internal/pkg/logger"
	"github.com/stretchr/testify/assert"
)

var regi = map[string]modules.Creator{}

func TestPlugin(t *testing.T) {
	logger.SetLevel(logger.DEBUG)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockModules(t, ctrl)

	buf := &bytes.Buffer{}

	plg := NewPlugin()
	plg.Config = &Config{
		Enabled:    true,
		DefaultRun: true,
		MaxProcs:   2,
		Modules:    map[string]bool{},
	}
	plg.Option = &cli.Option{
		Module:      "all",
		UpdateEvery: 1,
	}
	plg.registry = regi
	plg.Out = buf

	assert.True(t, plg.Setup())
	plg.CheckJobs()
	plg.MainLoop()
	time.Sleep(2 * time.Second)
	plg.Shutdown()
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

func mockModules(t *testing.T, ctrl *gomock.Controller) {
	regi["example1"] = modules.Creator{
		NoConfig: true,
		Create: func() modules.Module {
			mod := mock.NewMockModule(ctrl)
			mod.EXPECT().UpdateEvery().Return(1).AnyTimes()
			mod.EXPECT().SetUpdateEvery(gomock.Any()).AnyTimes()
			mod.EXPECT().ModuleName().Return("example1").AnyTimes()
			mod.EXPECT().SetModuleName(gomock.Any()).AnyTimes()
			mod.EXPECT().SetLogger(gomock.Any()).AnyTimes()
			mod.EXPECT().Init().Return(nil).AnyTimes()
			mod.EXPECT().Check().Return(true).AnyTimes()
			mod.EXPECT().GetCharts().Return(charts.NewCharts(
				&Chart{
					ID:   "chart1",
					Opts: Opts{Title: "Random Data 1", Units: "random", Fam: "random"},
					Dims: Dims{
						{ID: "random0", Name: "random"},
					},
				},
				&Chart{
					ID:   "chart2",
					Opts: Opts{Title: "Random Data 2", Units: "random", Fam: "random", Type: charts.Area},
					Dims: Dims{
						{ID: "random1", Name: "random"},
					},
				},
			)).AnyTimes()
			mod.EXPECT().GetData().Return(map[string]int64{
				"random0": 1,
				"random1": 2,
			}).AnyTimes()
			return mod
		},
	}
}

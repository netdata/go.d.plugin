package godplugin

import (
	"io/ioutil"
	"path"
	"runtime"
	"sync"

	"github.com/l2isbad/yaml"

	_ "github.com/l2isbad/go.d.plugin/internal/modules/all"
	"github.com/l2isbad/go.d.plugin/internal/pkg/cli"
	"github.com/l2isbad/go.d.plugin/internal/pkg/logger"
)

type P interface {
	Start()
}

func New(p, m string) P {
	return &goDPlugin{
		dir:  dir{p, m},
		conf: newConfig(),
		cli:  cli.Parse(),
	}

}

type dir struct {
	pluginConf  string
	modulesConf string
}

type goDPlugin struct {
	dir  dir
	conf config
	cli  cli.ParsedCLI
	wg   sync.WaitGroup
}

func (gd *goDPlugin) Start() {
	err := gd.loadConfig()

	if err != nil {
		log.Critical(err)
	}

	if !gd.conf.Enabled {
		log.Info("disabled in configuration file")
		return
	}

	if gd.cli.Debug {
		log.SetLevel(logger.DEBUG)
	}

	if gd.conf.MaxProcs != 0 {
		log.Warningf("setting GOMAXPROCS to %d", gd.conf.MaxProcs)
		runtime.GOMAXPROCS(gd.conf.MaxProcs)
	}

	gd.jobsRun(gd.jobsSet(gd.jobsCreate()))

	gd.wg.Wait()
}

func (gd *goDPlugin) loadConfig() error {
	f, err := ioutil.ReadFile(path.Join(gd.dir.pluginConf, "go.d.conf"))

	if err != nil {
		log.Error(err)
		return nil
	}

	return yaml.Unmarshal(f, &gd.conf)
}

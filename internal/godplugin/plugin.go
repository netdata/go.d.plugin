package godplugin

import (
	"io/ioutil"
	"path"
	"runtime"
	"sync"

	"fmt"

	_ "github.com/l2isbad/go.d.plugin/internal/modules/all" // load all modules
	"github.com/l2isbad/go.d.plugin/internal/pkg/cli"
	"github.com/l2isbad/go.d.plugin/internal/pkg/logger"
)

var (
	pluginConf = "go.d.conf"
	modConfDir = "go.d/"
)

type (
	// GoDPlugin GoDPlugin
	GoDPlugin interface {
		// LoadConfig Load go.d.conf
		LoadConfig(confDir string) error
		InitModules() error
		CheckModules() error
		MainLoop()
		Shutdown()
	}

	plugin struct {
		PluginConf     string
		ModulesConfDir string
		Config         *Config
		cmd            cli.ParsedCMD
		wg             sync.WaitGroup
	}
)

// func NewGoDPlugin(args []string) *GoDPlugin {
// 	p := getConfigDir()
// 	return &GoDPlugin{
// 		dir:    dir{p, path.Join(p, modConfDir)},
// 		config: newPluginConfig(),
// 		cmd:    cli.Parse(args),
// 	}
// }

func (p *plugin) LoadConfig(confDir string) {
	f, err := ioutil.ReadFile(path.Join(gd.dir.pluginConf, pluginConf))
}

func (gd *GoDPlugin) Start() {
	err := gd.loadConfig()

	if err != nil {
		log.Critical(err)
	}

	if !gd.config.Enabled {
		log.Info("disabled in configuration file")
		return
	}

	if gd.cmd.Debug {
		logger.SetLevel(logger.DEBUG)
	}

	if gd.config.MaxProcs != 0 {
		log.Warningf("setting GOMAXPROCS to %d", gd.config.MaxProcs)
		runtime.GOMAXPROCS(gd.config.MaxProcs)
	}

	jobs := gd.jobsCreate()
	gd.jobsStart(gd.jobsCheck(gd.jobsInit(jobs)))

	gd.wg.Wait()
	fmt.Println("DISABLE")
}

// func (gd *GoDPlugin) loadConfig() error {
// 	f, err := ioutil.ReadFile(path.Join(gd.dir.pluginConf, pluginConf))

// 	if err != nil {
// 		log.Error(err)
// 		return nil
// 	}

// 	return yaml.Unmarshal(f, &gd.conf)
// }

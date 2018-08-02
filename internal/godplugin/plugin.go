package godplugin

import (
	"runtime"

	"fmt"

	"time"

	"errors"

	"github.com/l2isbad/go.d.plugin/internal/godplugin/ticker"
	_ "github.com/l2isbad/go.d.plugin/internal/modules/all" // load all modules
	"github.com/l2isbad/go.d.plugin/internal/pkg/cli"
	"github.com/l2isbad/go.d.plugin/internal/pkg/logger"
)

var (
	modConfDir = "go.d/"
)

var ErrDisabled = errors.New("disabled in configuration file")

var log = logger.New("plugin", "main")

type (
	// GoDPlugin GoDPlugin
	Plugin struct {
		Option        *cli.Option
		Config        *Config
		ModuleConfDir string
		shutdownHook  chan int
	}

	//plugin struct {
	//	PluginConf     string
	//	ModulesConfDir string
	//	Config         *Config
	//	opt            *cli.Option
	//	wg             sync.WaitGroup
	//}
)

func NewPlugin() *Plugin {
	return &Plugin{
		shutdownHook: make(chan int, 1),
	}
}

func (p *Plugin) Setup() error {
	if !p.Config.Enabled {
		return ErrDisabled
	}

	if p.Config.MaxProcs > 0 {
		log.Infof("setting GOMAXPROCS to %d", p.Config.MaxProcs)
		runtime.GOMAXPROCS(p.Config.MaxProcs)
	}

	jobs := p.createJobs()
	return nil
}

func (p *Plugin) MainLoop() {
	tk := ticker.New(time.Second)
LOOP:
	for {
		select {
		case <-p.shutdownHook:
			break LOOP
		case <-tk.C:
		}
		// run job
	}
}

func (p *Plugin) Shutdown() {
	p.shutdownHook <- 1
}

// func NewGoDPlugin(args []string) *GoDPlugin {
// 	p := getConfigDir()
// 	return &GoDPlugin{
// 		dir:    dir{p, path.Join(p, modConfDir)},
// 		config: newPluginConfig(),
// 		cmd:    cli.Parse(args),
// 	}
// }

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

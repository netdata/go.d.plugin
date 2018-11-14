package godplugin

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/go-playground/validator"
	"github.com/l2isbad/go.d.plugin/cli"
	"github.com/l2isbad/go.d.plugin/godplugin/ticker"
	"github.com/l2isbad/go.d.plugin/logger"
	"github.com/l2isbad/go.d.plugin/modules"

	_ "github.com/l2isbad/go.d.plugin/modules/all"
)

type Job interface {
	ModuleName() string
	Name() string
	AutoDetectionRetry() int

	Initialized() bool
	Panicked() bool

	Init() bool
	Check() bool
	PostCheck() bool
	//Cleanup()

	Tick(clock int)
	MainLoop()
	Shutdown()
}

var log = logger.New("plugin", "main")
var validate = validator.New()

func New() *Plugin {
	return &Plugin{
		modules:   make(modules.Registry),
		loopQueue: make([]Job, 0),
		checkCh:   make(chan Job),
	}
}

type (
	// Plugin Plugin
	Plugin struct {
		Option        *cli.Option
		Config        *PluginConfig
		ModuleConfDir string
		Out           io.Writer

		modules   modules.Registry
		loopQueue []Job
		checkCh   chan Job
	}
)

func (p *Plugin) Setup() bool {
	if !p.Config.Enabled {
		fmt.Fprintln(p.Out, "DISABLE")
		log.Info("disabled in configuration file")
		return false
	}

	if p.Option.Module != "all" {
		if creator, exist := modules.DefaultRegistry[p.Option.Module]; exist {
			p.modules[p.Option.Module] = creator
		}
	} else {
		for name, creator := range modules.DefaultRegistry {
			if creator.DisabledByDefault && !p.Config.IsModuleEnabled(name, true) {
				log.Infof("'%s' disabled by default", name)
				continue
			}
			if !p.Config.IsModuleEnabled(name, false) {
				log.Infof("'%s' disabled in configuration file", name)
				continue
			}
			p.modules[name] = creator
		}
	}

	if len(p.modules) == 0 {
		log.Info("no modules to run")
		return false
	}

	if p.Config.MaxProcs > 0 {
		log.Infof("setting GOMAXPROCS to %d", p.Config.MaxProcs)
		runtime.GOMAXPROCS(p.Config.MaxProcs)
	}

	return true
}

func (p *Plugin) Serve() {
	go shutdownTask()
	go p.checkJobs()

	for _, job := range p.createJobs() {
		p.checkCh <- job
	}

	p.MainLoop()
}

func (p *Plugin) MainLoop() {
	log.Info("start main loop")
	var clock int
	tk := ticker.New(time.Second)

	for {
		select {
		case clock = <-tk.C:
			log.Debugf("tick %d", clock)
			for _, job := range p.loopQueue {
				log.Debugf("tick job: %s[%s]", job.ModuleName(), job.Name())
				job.Tick(clock)
			}
		}
	}
}

func (p *Plugin) checkJobs() {
	for job := range p.checkCh {
		if !job.Initialized() && !job.Init() {
			log.Errorf("%s[%s] Init failed", job.ModuleName(), job.Name())
			continue
		}

		ok := job.Check()

		if job.Panicked() {
			continue
		}

		if !ok && job.AutoDetectionRetry() > 0 {
			go recheckTask(p.checkCh, job)
			continue
		}

		if !ok {
			log.Errorf("%s[%s] Check failed", job.ModuleName(), job.Name())
			continue
		}

		if !job.PostCheck() {
			log.Errorf("%s[%s] PostCheck failed", job.ModuleName(), job.Name())
		}

		p.loopQueue = append(p.loopQueue, job)
		go job.MainLoop()
	}
}

func (p *Plugin) createJobs() []Job {
	var jobs []Job
	for name, creator := range p.modules {

		if creator.DisabledByDefault {
			continue
		}

		var modConfig moduleConfig

		// FIXME:
		err := modConfig.load(fmt.Sprintf("/opt/go.d/%s.conf", name))

		// FIXME:
		if err != nil {
			continue
		}

		// FIXME:
		if len(modConfig.Jobs) == 0 {
			continue
		}

		modConfig.updateJobConfigs()

		for _, conf := range modConfig.Jobs {
			mod := creator.Create()
			yaml.Unmarshal(asBytes(conf), mod)

			job := modules.NewJob(name, mod, p.Out)
			yaml.Unmarshal(asBytes(conf), job)

			jobs = append(jobs, job)
		}
	}
	return jobs
}

func asBytes(i interface{}) []byte {
	v, _ := yaml.Marshal(i)
	return v
}

func shutdownTask() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	for {
		switch <-signalChan {
		case syscall.SIGINT:
			log.Info("SIGINT received. Terminating...")
			os.Exit(0)
		}
	}
}

func recheckTask(ch chan Job, job Job) {
	time.Sleep(time.Second * time.Duration(job.AutoDetectionRetry()))
	ch <- job
}

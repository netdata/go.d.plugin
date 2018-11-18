package godplugin

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/go-playground/validator"
	"gopkg.in/yaml.v2"

	"github.com/l2isbad/go.d.plugin/cli"
	"github.com/l2isbad/go.d.plugin/godplugin/ticker"
	"github.com/l2isbad/go.d.plugin/logger"
	"github.com/l2isbad/go.d.plugin/modules"

	_ "github.com/l2isbad/go.d.plugin/modules/all"
)

type Job interface {
	FullName() string
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
		modules: make(modules.Registry),
		checkCh: make(chan Job, 1),
	}
}

type (
	// Plugin Plugin
	Plugin struct {
		Option        *cli.Option
		Config        *Config
		ModuleConfDir string
		Out           io.Writer

		modules   modules.Registry
		loopQueue []Job
		checkCh   chan Job
	}
)

func (p *Plugin) populateActiveModules() {
	if p.Option.Module != "all" {
		if creator, exist := modules.DefaultRegistry[p.Option.Module]; exist {
			p.modules[p.Option.Module] = creator
		}
		return
	}

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

func (p *Plugin) Setup() bool {
	if !p.Config.Enabled {
		fmt.Fprintln(p.Out, "DISABLE")
		log.Info("disabled in configuration file")
		return false
	}

	p.populateActiveModules()

	if len(p.modules) == 0 {
		log.Info("no modules to run")
		return false
	}

	if p.Config.MaxProcs > 0 {
		log.Infof("setting GOMAXPROCS to %d", p.Config.MaxProcs)
		runtime.GOMAXPROCS(p.Config.MaxProcs)
	}

	log.Infof("minimum update every is set to %d", p.Option.UpdateEvery)

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
		clock = <-tk.C

		log.Debugf("tick %d", clock)
		for _, job := range p.loopQueue {
			log.Debugf("tick job: %s[%s]", job.ModuleName(), job.Name())
			job.Tick(clock)
		}
	}
}

func (p *Plugin) checkJobs() {
	started := make(map[string]bool)

	for job := range p.checkCh {
		if started[job.FullName()] {
			log.Warningf("skipping %s[%s]: already served by another job", job.ModuleName(), job.Name())
			continue
		}

		if !job.Initialized() && !job.Init() {
			log.Errorf("%s[%s] Init failed", job.ModuleName(), job.Name())
			continue
		}

		ok := job.Check()

		if job.Panicked() {
			continue
		}

		if !ok {
			log.Errorf("%s[%s] Check failed", job.ModuleName(), job.Name())
			if job.AutoDetectionRetry() > 0 {
				go recheckTask(p.checkCh, job)
			}
			continue
		}

		if !job.PostCheck() {
			log.Errorf("%s[%s] PostCheck failed", job.ModuleName(), job.Name())
			continue
		}

		started[job.FullName()] = true

		log.Infof("%s[%s]: Check OK", job.ModuleName(), job.Name())
		// FIXME:
		p.loopQueue = append(p.loopQueue, job)
		go job.MainLoop()
	}
}

func (p *Plugin) createJobs() []Job {
	var jobs []Job
	for name, creator := range p.modules {
		var modConfig moduleConfig

		// FIXME:
		err := modConfig.load(fmt.Sprintf("/opt/go.d/%s.conf", name))

		if err != nil {
			log.Errorf("skipping %s: %v", name, err)
			continue
		}

		if len(modConfig.Jobs) == 0 {
			log.Errorf("skipping %s: config 'Jobs' section is empty or not exist", name)
			continue
		}

		jobName := func(conf rawConfig) interface{} {
			if name := conf["name"]; name != nil {
				return name
			}
			return "unnamed"
		}
		modConfig.updateJobs(modConfigDefaults(), creator.UpdateEvery, p.Option.UpdateEvery)

		for _, conf := range modConfig.Jobs {
			mod := creator.Create()

			if err := unmarshalAndValidate(conf, mod); err != nil {
				log.Errorf("skipping %s[%s]: %s", name, jobName(conf), err)
				continue
			}

			job := modules.NewJob(name, mod, p.Out)

			if err := unmarshalAndValidate(conf, job); err != nil {
				log.Errorf("skipping %s[%s]: %s", name, jobName(conf), err)
				continue
			}

			jobs = append(jobs, job)
		}
	}
	return jobs
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
	log.Infof("%s[%s] scheduling next check in %d seconds",
		job.ModuleName(),
		job.Name(),
		job.AutoDetectionRetry(),
	)
	time.Sleep(time.Second * time.Duration(job.AutoDetectionRetry()))
	ch <- job
}

func unmarshalAndValidate(rawConf map[string]interface{}, module interface{}) error {
	b, _ := yaml.Marshal(rawConf)
	if err := yaml.Unmarshal(b, module); err != nil {
		return err
	}
	if err := validate.Struct(module); err != nil {
		return err
	}
	return nil
}

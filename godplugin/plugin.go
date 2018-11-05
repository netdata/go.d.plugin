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

var log = logger.New("plugin", "main")
var validate = validator.New()

func New() *Plugin {
	return &Plugin{
		modules:   make(modules.Registry),
		loopQueue: make([]modules.Job, 0),
	}
}

type (
	// Plugin Plugin
	Plugin struct {
		Option        *cli.Option
		Config        *PluginConfig
		ModuleConfDir string
		Out           io.Writer
		modules       modules.Registry
		loopQueue     []modules.Job
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
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT)
		for {
			switch <-signalChan {
			case syscall.SIGINT:
				log.Info("SIGINT received. Terminating...")
				os.Exit(0)
			}
		}
	}()

	checkCh := p.createCheckTask()
	p.createJobs(checkCh)

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

func (p *Plugin) createCheckTask() chan modules.Job {
	ch := make(chan modules.Job)
	go func() {
		for job := range ch {
			if !job.Inited() && !job.Init() {
				log.Errorf("%s[%s] Init failed", job.ModuleName(), job.Name())
				continue
			}

			ok := job.Check()

			if ok {
				if job.PostCheck() {
					p.loopQueue = append(p.loopQueue, job)
					go job.MainLoop()
				}
				continue
			}

			if job.Panicked() {
				continue
			}

			if job.AutoDetectionRetry() > 0 {
				go func(j modules.Job) {
					time.Sleep(time.Second * time.Duration(j.AutoDetectionRetry()))
					ch <- j
				}(job)
			}
			log.Errorf("%s[%s] Check failed", job.ModuleName(), job.Name())
		}
	}()
	return ch
}

func (p *Plugin) createJobs(ch chan modules.Job) {
	for modName, creator := range p.modules {
		var rawConfigs rawModConfig
		err := rawConfigs.Load(fmt.Sprintf("/opt/go.d/%s.conf", modName))

		if err != nil && !(os.IsNotExist(err) || os.IsPermission(err)) {
			log.Errorf("skipping '%s': %s", modName, err)
			continue
		} else if err != nil {
			log.Errorf("'%s': %s", modName, err)
		}

		for _, rawConf := range rawConfigs {
			conf := modules.JobNewConfig()
			mod := creator.Create()
			b, _ := yaml.Marshal(rawConf)

			yaml.Unmarshal(b, conf)
			if err := validate.Struct(conf); err != nil {
				log.Error(err)
				continue
			}

			yaml.Unmarshal(b, mod)
			if err := validate.Struct(mod); err != nil {
				log.Error(err)
				continue
			}

			job := modules.NewJob(modName, mod, conf, os.Stdout)
			ch <- job
		}
	}
}

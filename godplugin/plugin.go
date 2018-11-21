package godplugin

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
	"sync"
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

	Tick(clock int)

	Start()
	Stop()
}

var log = logger.New("plugin", "main")
var validate = validator.New()

func New() *Plugin {
	return &Plugin{
		modules: make(modules.Registry),
		checkCh: make(chan Job, 1),
		loopQueue: &jobQueue{
			mux: &sync.Mutex{},
		},
	}
}

type jobQueue struct {
	mux   *sync.Mutex
	queue []Job
}

func (q *jobQueue) add(job Job) {
	q.mux.Lock()
	defer q.mux.Unlock()

	q.queue = append(q.queue, job)
}

func (q *jobQueue) pop(fullName string) Job {
	q.mux.Lock()
	defer q.mux.Unlock()

	for i, job := range q.queue {
		if job.FullName() == fullName {
			q.queue = append(q.queue[:i], q.queue[i+1:]...)
			return job
		}
	}
	return nil
}

func (q *jobQueue) notify(clock int) {
	q.mux.Lock()
	defer q.mux.Unlock()

	for _, job := range q.queue {
		job.Tick(clock)
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
		checkCh   chan Job
		loopQueue *jobQueue
	}
)

func (p *Plugin) RemoveFromQueue(fullName string) {
	job := p.loopQueue.pop(fullName)
	job.Stop()
}

func (p *Plugin) populateActiveModules() {
	if p.Option.Module != "all" {
		if creator, exist := modules.DefaultRegistry[p.Option.Module]; exist {
			p.modules[p.Option.Module] = creator
		}
		return
	}

	for name, creator := range modules.DefaultRegistry {
		if creator.DisabledByDefault && !p.Config.isModuleEnabled(name, true) {
			log.Infof("'%s' disabled by default", name)
			continue
		}
		if !p.Config.isModuleEnabled(name, false) {
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

	p.mainLoop()
}

func (p *Plugin) mainLoop() {
	log.Info("start main loop")
	var clock int
	tk := ticker.New(time.Second)

	for {
		clock = <-tk.C
		log.Debugf("tick %d", clock)
		p.loopQueue.notify(clock)
	}
}

func (p *Plugin) checkJobs() {
	started := make(map[string]bool)

	for job := range p.checkCh {
		if started[job.FullName()] {
			log.Infof("skipping %s[%s]: already served by another job", job.ModuleName(), job.Name())
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

		go job.Start()
		p.loopQueue.add(job)
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
		modConfig.updateJobs(creator.UpdateEvery, p.Option.UpdateEvery)

		for _, conf := range modConfig.Jobs {
			mod := creator.Create()

			if err := unmarshalAndValidate(conf, mod); err != nil {
				log.Errorf("skipping %s[%s]: %s", name, jobName(conf), err)
				continue
			}

			job := modules.NewJob(name, mod, p.Out, p)

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

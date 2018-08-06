package godplugin

import (
	"fmt"
	"io"
	"runtime"
	"sync"
	"time"

	"github.com/l2isbad/go.d.plugin/internal/godplugin/job"
	"github.com/l2isbad/go.d.plugin/internal/godplugin/ticker"
	"github.com/l2isbad/go.d.plugin/internal/modules"
	_ "github.com/l2isbad/go.d.plugin/internal/modules/all" // load all modules
	"github.com/l2isbad/go.d.plugin/internal/pkg/cli"
	"github.com/l2isbad/go.d.plugin/internal/pkg/logger"
)

var log = logger.New("plugin", "main")

type (
	// Plugin Plugin
	Plugin struct {
		Option        *cli.Option
		Config        *Config
		ModuleConfDir string
		Out           io.Writer
		registry      modules.Registry
		newJobFunc    job.Factory
		shutdownHook  chan int
		recheckJobs   jobSet
		runningJobs   jobSet
	}
)

func NewPlugin() *Plugin {
	return &Plugin{
		shutdownHook: make(chan int, 1),
		recheckJobs:  jobSet{},
		runningJobs:  jobSet{KeyFunc: keyFuncFullName},
		registry:     modules.DefaultRegistry,
		newJobFunc:   job.New,
	}
}

func (p *Plugin) Setup() bool {
	if !p.Config.Enabled {
		fmt.Fprintln(p.Out, "DISABLE")
		log.Info("disabled in configuration file")
		return false
	}

	if p.Config.MaxProcs > 0 {
		log.Infof("setting GOMAXPROCS to %d", p.Config.MaxProcs)
		runtime.GOMAXPROCS(p.Config.MaxProcs)
	}

	return true
}

func (p *Plugin) CheckJobs() {
	jobs := p.createJobs()
	passed := make(chan job.Job, len(jobs))
	recheck := make(chan job.Job, len(jobs))
	done := make(chan int)

	wg := sync.WaitGroup{}
	for _, job := range jobs {
		job := job
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := job.Init(); err != nil {
				log.Warningf("%v: init failed: %v", job, err)
				return
			}
			if !job.Check() {
				log.Infof("%v: check failed", job)
				if job.AutoDetectionRetry() > 0 {
					log.Infof("%v: added to recheck list", job)
					recheck <- job
				}
				return
			}
			if !job.PostCheck() {
				log.Warningf("%v: post-check failed", job)
				return
			}
			log.Debugf("%v: added to passed list", job)
			passed <- job
		}()
	}

	go func() {
		for job := range passed {
			if p.runningJobs.PutIfNotExist(job) {
				log.Debugf("%v: added to running list", job)
			} else {
				log.Debugf("%v: skipped due to a same name job has been added.", job)
			}
		}
		for job := range recheck {
			if !p.runningJobs.Exist(job) {
				p.recheckJobs.PutIfNotExist(job)
			}
		}
		done <- 1
	}()

	wg.Wait()
	close(passed)
	close(recheck)
	<-done
}

func (p *Plugin) MainLoop() {
	log.Info("start main loop")
	var clock int
	tk := ticker.New(time.Second)
LOOP:
	for {
		select {
		case <-p.shutdownHook:
			log.Debug("caught shutdown")
			break LOOP
		case clock = <-tk.C:
			log.Debugf("tick %d", clock)
		}
		p.runningJobs.Range(func(job job.Job) bool {
			log.Debugf("tick job: %s[%s]", job.ModuleName(), job.JobName())
			job.Tick(clock)
			return true
		})
		p.recheckJobs.Range(func(job job.Job) bool {
			if clock%job.AutoDetectionRetry() == 0 {
				log.Infof("recheck job: %s[%s]", job.ModuleName(), job.JobName())
				go func() {
					if !job.Check() {
						return
					}
					p.recheckJobs.Delete(job)
					if !job.PostCheck() {
						return
					}
					p.runningJobs.PutIfNotExist(job)
				}()
			}
			return true
		})
	}
	p.runningJobs.Range(func(job job.Job) bool {
		job.Shutdown()
		return true
	})
}

func (p *Plugin) Shutdown() {
	select {
	case p.shutdownHook <- 1:
	default:
	}
}

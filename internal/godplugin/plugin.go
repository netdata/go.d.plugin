package godplugin

import (
	"runtime"

	"fmt"

	"time"

	"sync"

	"io"

	"github.com/l2isbad/go.d.plugin/internal/godplugin/job"
	"github.com/l2isbad/go.d.plugin/internal/godplugin/ticker"
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
	passed := make(chan *job.Job, len(jobs))
	recheck := make(chan *job.Job, len(jobs))
	done := make(chan int)

	wg := sync.WaitGroup{}
	wg.Add(len(jobs))
	for _, job := range jobs {
		job := job

		go func() {
			defer wg.Done()
			if err := job.Init(); err != nil {
				log.Warningf("module: %s, job: %s: init failed: %v", job.ModuleName(), job.JobName(), err)
				return
			}
			if !job.Check() {
				if job.AutoDetectionRetry > 0 {
					recheck <- job
				}
				return
			}
			if !job.PostCheck() {
				return
			}
			passed <- job
		}()
	}

	go func() {
		for job := range passed {
			p.runningJobs.PutIfNotExist(job)
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
	var clock int
	tk := ticker.New(time.Second)
LOOP:
	for {
		select {
		case <-p.shutdownHook:
			break LOOP
		case clock = <-tk.C:
		}
		p.runningJobs.Range(func(job *job.Job) bool {
			select {
			case job.Tick <- clock:
			default:
			}
			return true
		})
		p.recheckJobs.Range(func(job *job.Job) bool {
			if clock%job.AutoDetectionRetry == 0 {
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
}

func (p *Plugin) Shutdown() {
	p.shutdownHook <- 1
}

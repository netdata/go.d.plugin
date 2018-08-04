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
	// GoDPlugin GoDPlugin
	Plugin struct {
		Option        *cli.Option
		Config        *Config
		ModuleConfDir string
		Out           io.Writer
		shutdownHook  chan int
		jobs          []*job.Job
	}
)

func NewPlugin() *Plugin {
	return &Plugin{
		shutdownHook: make(chan int, 1),
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
	passed := make(chan *job.Job, 10)
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
				return
			}
			passed <- job
		}()
	}

	go func() {
		jobMap := map[string]*job.Job{}
		for job := range passed {
			name := job.FullName()
			if _, exist := jobMap[name]; !exist {
				jobMap[name] = job
				p.jobs = append(p.jobs, job)
			}
		}
		done <- 1
	}()

	wg.Wait()
	close(passed)
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
		for _, job := range p.jobs {
			select {
			case job.Tick <- clock:
			default:
			}
		}
	}
}

func (p *Plugin) Shutdown() {
	p.shutdownHook <- 1
}

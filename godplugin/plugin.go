package godplugin

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/netdata/go.d.plugin/cli"
	"github.com/netdata/go.d.plugin/godplugin/ticker"
	"github.com/netdata/go.d.plugin/logger"
	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/multipath"
	// add modules to the registry
	_ "github.com/netdata/go.d.plugin/modules/all"

	"github.com/go-playground/validator"
)

// Job is an interface that represents a job.
type Job interface {
	FullName() string
	ModuleName() string
	Name() string

	AutoDetectionRetry() int

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

// New creates Plugin with default values.
func New() *Plugin {
	return &Plugin{
		modules:  make(modules.Registry),
		config:   newConfig(),
		confName: "go.d.conf",
		registry: modules.DefaultRegistry,

		jobCh: make(chan Job),

		jobStartShutdown: make(chan struct{}),
		mainShutdown:     make(chan struct{}),
	}
}

type (
	// Plugin represents go.d.plugin
	Plugin struct {
		Option     *cli.Option
		ConfigPath multipath.MultiPath
		Out        io.Writer

		confName string
		config   *config
		registry modules.Registry
		modules  modules.Registry

		jobStartShutdown chan struct{}
		jobCh            chan Job

		mainShutdown chan struct{}
		loopQueue    loopQueue
	}
)

// RemoveFromQueue removes job from the loop queue by full name.
func (p *Plugin) RemoveFromQueue(fullName string) {
	if job := p.loopQueue.remove(fullName); job != nil {
		job.Stop()
	}
}

// Serve Serve
func (p *Plugin) Serve() {
	go shutdownTask()
	go heartbeatTask()

	go p.jobStartLoop()

	for _, job := range p.createJobs() {
		p.jobCh <- job
	}

	p.mainLoop()
}

func (p *Plugin) mainLoop() {
	log.Info("start main loop")
	tk := ticker.New(time.Second)

LOOP:
	for {
		select {
		case <-p.mainShutdown:
			break LOOP
		case clock := <-tk.C:
			p.runOnce(clock)
		}
	}
}

func (p *Plugin) runOnce(clock int) {
	log.Debugf("tick %d", clock)
	p.loopQueue.notify(clock)
}

func (p *Plugin) stop() {
	p.jobStartShutdown <- struct{}{}
	p.mainShutdown <- struct{}{}
}

func shutdownTask() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGPIPE)

	switch <-signalChan {
	case syscall.SIGINT:
		log.Info("SIGINT received. Terminating...")
	case syscall.SIGHUP:
		log.Info("SIGHUP received. Terminating...")
	case syscall.SIGPIPE:
		log.Critical("SIGPIPE received. Terminating...")
		os.Exit(1)
	}
	os.Exit(0)
}

func heartbeatTask() {
	t := time.Tick(time.Second)
	for range t {
		_, _ = fmt.Fprint(os.Stdout, "\n")
	}
}

package godplugin

import (
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator"

	"github.com/netdata/go.d.plugin/cli"
	"github.com/netdata/go.d.plugin/godplugin/ticker"
	"github.com/netdata/go.d.plugin/logger"
	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/multipath"

	_ "github.com/netdata/go.d.plugin/modules/all"
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
		checkCh:  make(chan Job, 1),
		config:   newConfig(),
		confName: "go.d.conf",
		registry: modules.DefaultRegistry,
	}
}

type (
	// Plugin represents go.d.plugin
	Plugin struct {
		Option     *cli.Option
		ConfigPath multipath.MultiPath
		Out        io.Writer

		confName  string
		config    *config
		registry  modules.Registry
		modules   modules.Registry
		checkCh   chan Job
		loopQueue loopQueue
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

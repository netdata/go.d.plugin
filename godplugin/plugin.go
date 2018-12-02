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

func New() *Plugin {
	return &Plugin{
		modules: make(modules.Registry),
		checkCh: make(chan Job, 1),
		config:  newConfig(),
	}
}

type (
	// config config
	Plugin struct {
		Option     *cli.Option
		ConfigPath multipath.MultiPath
		Out        io.Writer

		config    *config
		modules   modules.Registry
		checkCh   chan Job
		loopQueue jobQueue
	}
)

func (p *Plugin) RemoveFromQueue(fullName string) {
	if job := p.loopQueue.remove(fullName); job != nil {
		job.Stop()
	}
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

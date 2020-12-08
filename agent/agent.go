package agent

import (
	"context"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/netdata/go.d.plugin/agent/job/build"
	"github.com/netdata/go.d.plugin/agent/job/confgroup"
	"github.com/netdata/go.d.plugin/agent/job/discovery"
	"github.com/netdata/go.d.plugin/agent/job/registry"
	"github.com/netdata/go.d.plugin/agent/job/run"
	"github.com/netdata/go.d.plugin/agent/job/state"
	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/agent/netdataapi"
	"github.com/netdata/go.d.plugin/logger"
	"github.com/netdata/go.d.plugin/pkg/multipath"

	"github.com/mattn/go-isatty"
)

var isTerminal = isatty.IsTerminal(os.Stdout.Fd())

// Config is an Agent configuration.
type Config struct {
	Name              string
	ConfDir           []string
	ModulesConfDir    []string
	ModulesSDConfPath []string
	StateFile         string
	LockDir           string
	ModuleRegistry    module.Registry
	RunModule         string
	MinUpdateEvery    int
}

// Agent represents orchestrator.
type Agent struct {
	Name              string
	ConfDir           multipath.MultiPath
	ModulesConfDir    multipath.MultiPath
	ModulesSDConfPath []string
	StateFile         string
	LockDir           string
	RunModule         string
	MinUpdateEvery    int
	ModuleRegistry    module.Registry
	Out               io.Writer
	api               *netdataapi.API
	*logger.Logger
}

// New creates a new Agent.
func New(cfg Config) *Agent {
	p := &Agent{
		Name:              cfg.Name,
		ConfDir:           cfg.ConfDir,
		ModulesConfDir:    cfg.ModulesConfDir,
		ModulesSDConfPath: cfg.ModulesSDConfPath,
		StateFile:         cfg.StateFile,
		LockDir:           cfg.LockDir,
		RunModule:         cfg.RunModule,
		MinUpdateEvery:    cfg.MinUpdateEvery,
		ModuleRegistry:    module.DefaultRegistry,
		Out:               os.Stdout,
	}

	logger.Prefix = p.Name
	p.Logger = logger.New("main", "main")
	p.api = netdataapi.New(p.Out)

	return p
}

// Run
func (a *Agent) Run() {
	go a.signalHandling()
	go a.keepAlive()
	serve(a)
}

func serve(p *Agent) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP)
	var wg sync.WaitGroup

	for {
		ctx, cancel := context.WithCancel(context.Background())

		wg.Add(1)
		go func() { defer wg.Done(); p.run(ctx) }()

		sig := <-ch
		p.Infof("received %s signal (%d), stopping running instance", sig, sig)
		cancel()
		wg.Wait()
		time.Sleep(time.Second)
	}
}

func (a *Agent) run(ctx context.Context) {
	a.Info("instance is started")
	defer func() { a.Info("instance is stopped") }()

	cfg := a.loadPluginConfig()
	a.Infof("using config: %s", cfg)
	if !cfg.Enabled {
		a.Info("plugin is disabled in the configuration file, exiting...")
		if isTerminal {
			os.Exit(0)
		}
		_ = a.api.DISABLE()
		return
	}

	enabled := a.loadEnabledModules(cfg)
	if len(enabled) == 0 {
		a.Info("no modules to run")
		if isTerminal {
			os.Exit(0)
		}
		_ = a.api.DISABLE()
		return
	}

	discCfg := a.buildDiscoveryConf(enabled)

	discoverer, err := discovery.NewManager(discCfg)
	if err != nil {
		a.Error(err)
		if isTerminal {
			os.Exit(0)
		}
		return
	}

	runner := run.NewManager()

	builder := build.NewManager()
	builder.Runner = runner
	builder.PluginName = a.Name
	builder.Out = a.Out
	builder.Modules = enabled

	if a.LockDir != "" {
		builder.Registry = registry.NewFileLockRegistry(a.LockDir)
	}

	var saver *state.Manager
	if !isTerminal && a.StateFile != "" {
		saver = state.NewManager(a.StateFile)
		builder.CurState = saver
		if store, err := state.Load(a.StateFile); err != nil {
			a.Warningf("couldn't load state file: %v", err)
		} else {
			builder.PrevState = store
		}
	}

	in := make(chan []*confgroup.Group)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() { defer wg.Done(); runner.Run(ctx) }()

	wg.Add(1)
	go func() { defer wg.Done(); builder.Run(ctx, in) }()

	wg.Add(1)
	go func() { defer wg.Done(); discoverer.Run(ctx, in) }()

	if saver != nil {
		wg.Add(1)
		go func() { defer wg.Done(); saver.Run(ctx) }()
	}

	wg.Wait()
	<-ctx.Done()
	runner.Cleanup()
}

func (a *Agent) signalHandling() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGPIPE)

	sig := <-ch
	a.Infof("received %s signal (%d). Terminating...", sig, sig)
	os.Exit(0)
}

func (a *Agent) keepAlive() {
	if isTerminal {
		return
	}

	tk := time.NewTicker(time.Second)
	defer tk.Stop()

	for range tk.C {
		_ = a.api.EMPTYLINE()
	}
}

// SPDX-License-Identifier: GPL-3.0-or-later

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
	"github.com/netdata/go.d.plugin/agent/job/vnode"
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
	VnodesConfDir     []string
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
	VnodesConfDir     multipath.MultiPath
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
		VnodesConfDir:     cfg.VnodesConfDir,
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

// Run starts the Agent.
func (a *Agent) Run() {
	go a.keepAlive()
	serve(a)
}

func serve(p *Agent) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup

	var exit bool

	for {
		ctx, cancel := context.WithCancel(context.Background())

		wg.Add(1)
		go func() { defer wg.Done(); p.run(ctx) }()

		switch sig := <-ch; sig {
		case syscall.SIGHUP:
			p.Infof("received %s signal (%d). Restarting running instance", sig, sig)
		default:
			p.Infof("received %s signal (%d). Terminating...", sig, sig)
			module.DontObsoleteCharts()
			exit = true
		}

		cancel()

		func() {
			timeout := time.Second * 15
			t := time.NewTimer(timeout)
			defer t.Stop()
			done := make(chan struct{})

			go func() { wg.Wait(); close(done) }()

			select {
			case <-t.C:
				p.Errorf("stopping all goroutines timed out after %s. Exiting...", timeout)
				os.Exit(0)
			case <-done:
			}
		}()

		if exit {
			os.Exit(0)
		}

		time.Sleep(time.Second)
	}
}

func (a *Agent) run(ctx context.Context) {
	a.Info("instance is started")
	defer func() { a.Info("instance is stopped") }()

	cfg := a.loadPluginConfig()
	a.Infof("using config: %s", cfg.String())
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

	if reg := a.setupVnodeRegistry(); reg == nil || reg.Len() == 0 {
		vnode.Disabled = true
	} else {
		builder.VNodeRegistry = reg
	}

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

func (a *Agent) keepAlive() {
	if isTerminal {
		return
	}

	tk := time.NewTicker(time.Second)
	defer tk.Stop()

	var n int
	for range tk.C {
		if err := a.api.EMPTYLINE(); err != nil {
			a.Infof("keepAlive: %v", err)
			n++
		} else {
			n = 0
		}
		if n == 3 {
			a.Info("too many keepAlive errors. Terminating...")
			os.Exit(0)
		}
	}
}

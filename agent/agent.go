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

	"github.com/netdata/go.d.plugin/agent/confgroup"
	"github.com/netdata/go.d.plugin/agent/discovery"
	"github.com/netdata/go.d.plugin/agent/discovery/dyncfg"
	"github.com/netdata/go.d.plugin/agent/filelock"
	"github.com/netdata/go.d.plugin/agent/filestatus"
	"github.com/netdata/go.d.plugin/agent/functions"
	"github.com/netdata/go.d.plugin/agent/jobmgr"
	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/agent/netdataapi"
	"github.com/netdata/go.d.plugin/agent/safewriter"
	"github.com/netdata/go.d.plugin/agent/vnodes"
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
	*logger.Logger

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

	api *netdataapi.API
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
		Out:               safewriter.New(os.Stdout),
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

	discoveryManager, err := discovery.NewManager(discCfg)
	if err != nil {
		a.Error(err)
		if isTerminal {
			os.Exit(0)
		}
		return
	}

	functionsManager := functions.NewManager()

	dyncfgDiscovery, _ := dyncfg.NewDiscovery(dyncfg.Config{
		PluginName:       a.Name,
		Out:              a.Out,
		Modules:          enabled,
		ModuleDefaults:   discCfg.Registry,
		FunctionRegistry: functionsManager,
	})

	discoveryManager.Add(dyncfgDiscovery)

	jobsManager := jobmgr.NewManager()
	jobsManager.Dyncfg = dyncfgDiscovery
	jobsManager.PluginName = a.Name
	jobsManager.Out = a.Out
	jobsManager.Modules = enabled

	if reg := a.setupVnodeRegistry(); reg == nil || reg.Len() == 0 {
		vnodes.Disabled = true
	} else {
		jobsManager.Vnodes = reg
	}

	if a.LockDir != "" {
		jobsManager.FileLock = filelock.New(a.LockDir)
	}

	var statusSaveManager *filestatus.Manager
	if !isTerminal && a.StateFile != "" {
		statusSaveManager = filestatus.NewManager(a.StateFile)
		jobsManager.StatusSaver = statusSaveManager
		if store, err := filestatus.LoadStore(a.StateFile); err != nil {
			a.Warningf("couldn't load state file: %v", err)
		} else {
			jobsManager.StatusStore = store
		}
	}

	in := make(chan []*confgroup.Group)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() { defer wg.Done(); functionsManager.Run(ctx) }()

	wg.Add(1)
	go func() { defer wg.Done(); jobsManager.Run(ctx, in) }()

	wg.Add(1)
	go func() { defer wg.Done(); discoveryManager.Run(ctx, in) }()

	if statusSaveManager != nil {
		wg.Add(1)
		go func() { defer wg.Done(); statusSaveManager.Run(ctx) }()
	}

	wg.Wait()
	<-ctx.Done()
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

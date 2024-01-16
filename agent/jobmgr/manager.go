// SPDX-License-Identifier: GPL-3.0-or-later

package jobmgr

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/netdata/go.d.plugin/agent/confgroup"
	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/agent/netdataapi"
	"github.com/netdata/go.d.plugin/agent/safewriter"
	"github.com/netdata/go.d.plugin/agent/ticker"
	"github.com/netdata/go.d.plugin/logger"

	"gopkg.in/yaml.v2"
)

func New() *Manager {
	mgr := &Manager{
		Logger: logger.New().With(
			slog.String("component", "job manager"),
		),
		Out:             io.Discard,
		FileLock:        noop{},
		FileStatus:      noop{},
		FileStatusStore: noop{},
		Vnodes:          noop{},
		FnReg:           noop{},

		discoveredConfigs: newDiscoveredConfigsCache(),
		seenConfigs:       newSeenConfigCache(),
		exposedConfigs:    newExposedConfigCache(),
		runningJobs:       newRunningJobsCache(),
		retryingTasks:     newRetryingTasksCache(),

		retryCh: make(chan confgroup.Config),
		api:     netdataapi.New(safewriter.Stdout),
		mux:     sync.Mutex{},
		started: make(chan struct{}),
	}

	return mgr
}

type Manager struct {
	*logger.Logger

	PluginName     string
	Out            io.Writer
	Modules        module.Registry
	ConfigDefaults confgroup.Registry

	FileLock        FileLocker
	FileStatus      FileStatus
	FileStatusStore FileStatusStore
	Vnodes          Vnodes
	FnReg           FunctionRegistry

	discoveredConfigs *discoveredConfigs
	seenConfigs       *seenConfigs
	exposedConfigs    *exposedConfigs
	retryingTasks     *retryingTasks
	runningJobs       *runningJobs

	api     DyncfgAPI
	ctx     context.Context
	retryCh chan confgroup.Config
	mux     sync.Mutex

	started chan struct{}
}

func (m *Manager) Run(ctx context.Context, in chan []*confgroup.Group) {
	m.Info("instance is started")
	defer func() { m.cleanup(); m.Info("instance is stopped") }()
	m.ctx = ctx

	m.FnReg.Register("config", m.dyncfgConfig)

	for name := range m.Modules {
		m.dyncfgModuleCreate(name)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() { defer wg.Done(); m.runProcessDiscoveredConfigs(in) }()

	wg.Add(1)
	go func() { defer wg.Done(); m.runNotifyRunningJobs() }()

	close(m.started)

	wg.Wait()
	<-m.ctx.Done()
}

func (m *Manager) runProcessDiscoveredConfigs(in chan []*confgroup.Group) {
	for {
		select {
		case <-m.ctx.Done():
			return
		case groups := <-in:
			m.processDiscoveredConfigGroups(groups)
		case cfg := <-m.retryCh:
			m.addDiscoveredConfig(cfg)
		}
	}
}

func (m *Manager) processDiscoveredConfigGroups(groups []*confgroup.Group) {
	for _, gr := range groups {
		a, r := m.discoveredConfigs.add(gr)
		m.Debugf("received configs: %d/+%d/-%d (group '%s')", len(gr.Configs), len(a), len(r), gr.Source)
		for _, cfg := range r {
			m.removeDiscoveredConfig(cfg)
		}
		for _, cfg := range a {
			m.addDiscoveredConfig(cfg)
		}
	}
}

func (m *Manager) addDiscoveredConfig(cfg confgroup.Config) {
	m.mux.Lock()
	defer m.mux.Unlock()

	task, isRetry := m.retryingTasks.lookup(cfg)
	if isRetry {
		m.retryingTasks.remove(cfg)
	}

	scfg, ok := m.seenConfigs.lookup(cfg)
	if !ok {
		scfg = &seenConfig{cfg: cfg}
		m.seenConfigs.add(scfg)
	}

	ecfg, ok := m.exposedConfigs.lookup(cfg)
	if !ok {
		ecfg = scfg
		m.exposedConfigs.add(ecfg)
	}

	if ok {
		sp, ep := scfg.cfg.SourceTypePriority(), ecfg.cfg.SourceTypePriority()
		if ep > sp || (ep == sp && ecfg.status == dyncfgRunning) {
			return
		}
		m.stopRunningJob(ecfg.cfg.FullName())
		m.exposedConfigs.add(scfg) // replace
		ecfg = scfg
	}

	job, err := m.createCollectorJob(ecfg.cfg)
	if err != nil {
		ecfg.status = dyncfgFailed
		if !isStock(ecfg.cfg) {
			m.dyncfgJobCreate(ecfg.cfg, ecfg.status)
		}
		return
	}

	if isRetry {
		job.AutoDetectEvery = task.timeout
		job.AutoDetectTries = task.retries
	} else if job.AutoDetectionEvery() == 0 {
		if m.FileStatusStore.Contains(ecfg.cfg, "") {

		}
	}

	if err := job.AutoDetection(); err != nil {
		job.Cleanup()
		ecfg.status = dyncfgFailed
		if !isStock(ecfg.cfg) {
			m.dyncfgJobCreate(ecfg.cfg, ecfg.status)
		}
		if job.RetryAutoDetection() {
			ctx, cancel := context.WithCancel(m.ctx)
			r := &retryTask{cancel: cancel, timeout: job.AutoDetectionEvery(), retries: job.AutoDetectTries}
			m.retryingTasks.add(cfg, r)
			go runRetryTask(ctx, m.retryCh, ecfg.cfg)
		}
		return
	}

	ecfg.status = dyncfgRunning
	m.startRunningJob(job)
	m.dyncfgJobCreate(ecfg.cfg, ecfg.status)
}

func runRetryTask(ctx context.Context, out chan<- confgroup.Config, cfg confgroup.Config) {
	t := time.NewTimer(time.Second * time.Duration(cfg.AutoDetectionRetry()))
	defer t.Stop()

	select {
	case <-ctx.Done():
	case <-t.C:
		select {
		case <-ctx.Done():
		case out <- cfg:
		}
	}
}

func (m *Manager) removeDiscoveredConfig(cfg confgroup.Config) {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.retryingTasks.remove(cfg)

	scfg, ok := m.seenConfigs.lookup(cfg)
	if !ok {
		return
	}
	m.seenConfigs.remove(cfg)

	ecfg, ok := m.exposedConfigs.lookup(cfg)
	if !ok {
		return
	}
	if scfg.cfg.UID() == ecfg.cfg.UID() {
		m.exposedConfigs.remove(cfg)
		m.stopRunningJob(cfg.FullName())
		m.dyncfgJobRemove(cfg)
	}

	return
}

func (m *Manager) runNotifyRunningJobs() {
	tk := ticker.New(time.Second)
	defer tk.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case clock := <-tk.C:
			m.runningJobs.lock()
			m.runningJobs.forEach(func(_ string, job *module.Job) {
				job.Tick(clock)
			})
			m.runningJobs.unlock()
		}
	}
}

func (m *Manager) cleanup() {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.FnReg.Unregister("config")

	m.runningJobs.lock()
	defer m.runningJobs.unlock()

	m.runningJobs.forEach(func(key string, job *module.Job) {
		job.Stop()
		m.runningJobs.remove(key)
	})
}

func (m *Manager) startRunningJob(job *module.Job) {
	m.runningJobs.lock()
	defer m.runningJobs.unlock()

	if job, ok := m.runningJobs.lookup(job.FullName()); ok {
		job.Stop()
	}

	go job.Start()
	m.runningJobs.add(job.FullName(), job)
}

func (m *Manager) stopRunningJob(name string) {
	m.runningJobs.lock()
	defer m.runningJobs.unlock()

	if job, ok := m.runningJobs.lookup(name); ok {
		job.Stop()
		m.runningJobs.remove(name)
	}
}

func (m *Manager) createCollectorJob(cfg confgroup.Config) (*module.Job, error) {
	creator, ok := m.Modules[cfg.Module()]
	if !ok {
		return nil, fmt.Errorf("can not find %s module", cfg.Module())
	}

	var vnode struct {
		guid     string
		hostname string
		labels   map[string]string
	}

	if cfg.Vnode() != "" {
		n, ok := m.Vnodes.Lookup(cfg.Vnode())
		if !ok {
			return nil, fmt.Errorf("vnode '%s' is not found", cfg.Vnode())
		}

		vnode.guid = n.GUID
		vnode.hostname = n.Hostname
		vnode.labels = n.Labels
	}

	m.Debugf("creating %s[%s] job, config: %v", cfg.Module(), cfg.Name(), cfg)

	mod := creator.Create()

	if err := applyConfig(cfg, mod); err != nil {
		return nil, err
	}

	jobCfg := module.JobConfig{
		PluginName:      m.PluginName,
		Name:            cfg.Name(),
		ModuleName:      cfg.Module(),
		FullName:        cfg.FullName(),
		UpdateEvery:     cfg.UpdateEvery(),
		AutoDetectEvery: cfg.AutoDetectionRetry(),
		Priority:        cfg.Priority(),
		Labels:          makeLabels(cfg),
		IsStock:         cfg.SourceType() == "stock",
		Module:          mod,
		Out:             m.Out,
		VnodeGUID:       vnode.guid,
		VnodeHostname:   vnode.hostname,
		VnodeLabels:     vnode.labels,
	}

	job := module.NewJob(jobCfg)

	return job, nil
}

func isStock(cfg confgroup.Config) bool {
	return cfg.SourceType() == "stock"
}

func applyConfig(cfg confgroup.Config, module any) error {
	bs, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bs, module)
}

func isTooManyOpenFiles(err error) bool {
	return err != nil && strings.Contains(err.Error(), "too many open files")
}

func makeLabels(cfg confgroup.Config) map[string]string {
	labels := make(map[string]string)
	for name, value := range cfg.Labels() {
		n, ok1 := name.(string)
		v, ok2 := value.(string)
		if ok1 && ok2 {
			labels[n] = v
		}
	}
	return labels
}

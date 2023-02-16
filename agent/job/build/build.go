// SPDX-License-Identifier: GPL-3.0-or-later

package build

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	jobpkg "github.com/netdata/go.d.plugin/agent/job"
	"github.com/netdata/go.d.plugin/agent/job/confgroup"
	"github.com/netdata/go.d.plugin/agent/job/vnode"
	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/logger"

	"gopkg.in/yaml.v2"
)

type Runner interface {
	Start(job jobpkg.Job)
	Stop(fullName string)
}

type StateSaver interface {
	Save(cfg confgroup.Config, state string)
	Remove(cfg confgroup.Config)
}

type State interface {
	Contains(cfg confgroup.Config, states ...string) bool
}

type Registry interface {
	Register(name string) (bool, error)
	Unregister(name string) error
}

type VNodeRegistry interface {
	Lookup(key string) (*vnode.VirtualNode, bool)
}

type (
	noopSaver         struct{}
	noopState         struct{}
	noopRegistry      struct{}
	noopVnodeRegistry struct{}
)

func (n noopSaver) Save(_ confgroup.Config, _ string) {}
func (n noopSaver) Remove(_ confgroup.Config)         {}

func (n noopState) Contains(_ confgroup.Config, _ ...string) bool { return false }

func (n noopRegistry) Register(_ string) (bool, error) { return true, nil }
func (n noopRegistry) Unregister(_ string) error       { return nil }

func (n noopVnodeRegistry) Lookup(_ string) (*vnode.VirtualNode, bool) { return nil, false }

type state = string

const (
	success           state = "success"            // successfully started
	retry             state = "retry"              // failed, but we need keep trying auto-detection
	failed            state = "failed"             // failed
	duplicateLocal    state = "duplicate_local"    // a job with the same FullName is started
	duplicateGlobal   state = "duplicate_global"   // a job with the same FullName is registered by another plugin
	registrationError state = "registration_error" // an error during registration (only 'too many open files')
	buildError        state = "build_error"        // an error during building
)

type (
	Manager struct {
		PluginName string
		Out        io.Writer
		Modules    module.Registry
		*logger.Logger

		Runner        Runner
		CurState      StateSaver
		PrevState     State
		Registry      Registry
		VNodeRegistry VNodeRegistry

		grpCache   *groupCache
		startCache *startedCache
		retryCache *retryCache

		addCh    chan []confgroup.Config
		removeCh chan []confgroup.Config
		retryCh  chan confgroup.Config
	}
)

func NewManager() *Manager {
	mgr := &Manager{
		CurState:      noopSaver{},
		PrevState:     noopState{},
		Registry:      noopRegistry{},
		VNodeRegistry: noopVnodeRegistry{},
		Out:           io.Discard,
		Logger:        logger.New("build", "manager"),
		grpCache:      newGroupCache(),
		startCache:    newStartedCache(),
		retryCache:    newRetryCache(),
		addCh:         make(chan []confgroup.Config),
		removeCh:      make(chan []confgroup.Config),
		retryCh:       make(chan confgroup.Config),
	}
	return mgr
}

func (m *Manager) Run(ctx context.Context, in chan []*confgroup.Group) {
	m.Info("instance is started")
	defer func() { m.cleanup(); m.Info("instance is stopped") }()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() { defer wg.Done(); m.runGroupProcessing(ctx, in) }()

	wg.Add(1)
	go func() { defer wg.Done(); m.runConfigProcessing(ctx) }()

	wg.Wait()
	<-ctx.Done()
}

func (m *Manager) cleanup() {
	for _, task := range *m.retryCache {
		task.cancel()
	}
	for name := range *m.startCache {
		_ = m.Registry.Unregister(name)
	}
}

func (m *Manager) runGroupProcessing(ctx context.Context, in <-chan []*confgroup.Group) {
	for {
		select {
		case <-ctx.Done():
			return
		case groups := <-in:
			for _, group := range groups {
				select {
				case <-ctx.Done():
					return
				default:
					m.processGroup(ctx, group)
				}
			}
		}
	}
}

func (m *Manager) processGroup(ctx context.Context, group *confgroup.Group) {
	if group == nil {
		return
	}
	added, removed := m.grpCache.put(group)
	m.Debugf("received config group ('%s'): %d jobs (added: %d, removed: %d)",
		group.Source, len(group.Configs), len(added), len(removed))

	select {
	case <-ctx.Done():
		return
	case m.removeCh <- removed:
	}

	select {
	case <-ctx.Done():
		return
	case m.addCh <- added:
	}
}

func (m *Manager) runConfigProcessing(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case cfgs := <-m.addCh:
			m.handleAdd(ctx, cfgs)
		case cfgs := <-m.removeCh:
			m.handleRemove(ctx, cfgs)
		case cfg := <-m.retryCh:
			m.handleAddCfg(ctx, cfg)
		}
	}
}

func (m *Manager) handleAdd(ctx context.Context, cfgs []confgroup.Config) {
	for _, cfg := range cfgs {
		select {
		case <-ctx.Done():
			return
		default:
			m.handleAddCfg(ctx, cfg)
		}
	}
}

func (m *Manager) handleRemove(ctx context.Context, cfgs []confgroup.Config) {
	for _, cfg := range cfgs {
		select {
		case <-ctx.Done():
			return
		default:
			m.handleRemoveCfg(cfg)
		}
	}
}

func (m *Manager) handleAddCfg(ctx context.Context, cfg confgroup.Config) {
	if m.startCache.has(cfg) {
		m.Infof("%s[%s] job is being served by another job, skipping it", cfg.Module(), cfg.Name())
		m.CurState.Save(cfg, duplicateLocal)
		return
	}

	task, isRetry := m.retryCache.lookup(cfg)
	if isRetry {
		task.cancel()
		m.retryCache.remove(cfg)
	}

	job, err := m.buildJob(cfg)
	if err != nil {
		m.Warningf("couldn't build %s[%s]: %v", cfg.Module(), cfg.Name(), err)
		m.CurState.Save(cfg, buildError)
		return
	}
	cleanupJob := true
	defer func() {
		if cleanupJob {
			job.Cleanup()
		}
	}()

	if isRetry {
		job.AutoDetectEvery = task.timeout
		job.AutoDetectTries = task.retries
	} else if job.AutoDetectionEvery() == 0 {
		switch {
		case m.PrevState.Contains(cfg, success, retry):
			m.Infof("%s[%s] job last state is active/retry, applying recovering settings", cfg.Module(), cfg.Name())
			job.AutoDetectEvery = 30
			job.AutoDetectTries = 11
		case isInsideK8sCluster() && cfg.Provider() == "file watcher":
			m.Infof("%s[%s] is k8s job, applying recovering settings", cfg.Module(), cfg.Name())
			job.AutoDetectEvery = 10
			job.AutoDetectTries = 7
		}
	}

	switch detection(job) {
	case success:
		if ok, err := m.Registry.Register(cfg.FullName()); ok || err != nil && !isTooManyOpenFiles(err) {
			m.CurState.Save(cfg, success)
			m.Runner.Start(job)
			m.startCache.put(cfg)
			cleanupJob = false
		} else if isTooManyOpenFiles(err) {
			m.Error(err)
			m.CurState.Save(cfg, registrationError)
		} else {
			m.Infof("%s[%s] job is being served by another plugin, skipping it", cfg.Module(), cfg.Name())
			m.CurState.Save(cfg, duplicateGlobal)
		}
	case retry:
		m.Infof("%s[%s] job detection failed, will retry in %d seconds",
			cfg.Module(), cfg.Name(), job.AutoDetectionEvery())
		m.CurState.Save(cfg, retry)
		ctx, cancel := context.WithCancel(ctx)
		m.retryCache.put(cfg, retryTask{
			cancel:  cancel,
			timeout: job.AutoDetectionEvery(),
			retries: job.AutoDetectTries,
		})
		timeout := time.Second * time.Duration(job.AutoDetectionEvery())
		go runRetryTask(ctx, m.retryCh, cfg, timeout)
	case failed:
		m.CurState.Save(cfg, failed)
	default:
		m.Warningf("%s[%s] job detection: unknown state", cfg.Module(), cfg.Name())
	}
}

func (m *Manager) handleRemoveCfg(cfg confgroup.Config) {
	defer m.CurState.Remove(cfg)

	if m.startCache.has(cfg) {
		m.Runner.Stop(cfg.FullName())
		_ = m.Registry.Unregister(cfg.FullName())
		m.startCache.remove(cfg)
	}

	if task, ok := m.retryCache.lookup(cfg); ok {
		task.cancel()
		m.retryCache.remove(cfg)
	}
}

func (m *Manager) buildJob(cfg confgroup.Config) (*module.Job, error) {
	creator, ok := m.Modules[cfg.Module()]
	if !ok {
		return nil, fmt.Errorf("can not find %s module", cfg.Module())
	}

	m.Debugf("building %s[%s] job, config: %v", cfg.Module(), cfg.Name(), cfg)
	mod := creator.Create()
	if err := unmarshal(cfg, mod); err != nil {
		return nil, err
	}

	labels := make(map[string]string)
	for name, value := range cfg.Labels() {
		n, ok1 := name.(string)
		v, ok2 := value.(string)
		if ok1 && ok2 {
			labels[n] = v
		}
	}

	jobCfg := module.JobConfig{
		PluginName:      m.PluginName,
		Name:            cfg.Name(),
		ModuleName:      cfg.Module(),
		FullName:        cfg.FullName(),
		UpdateEvery:     cfg.UpdateEvery(),
		AutoDetectEvery: cfg.AutoDetectionRetry(),
		Priority:        cfg.Priority(),
		Labels:          labels,
		Module:          mod,
		Out:             m.Out,
	}

	if cfg.Vnode() != "" {
		n, ok := m.VNodeRegistry.Lookup(cfg.Vnode())
		if !ok {
			return nil, fmt.Errorf("vnode '%s' is not found", cfg.Vnode())
		}

		jobCfg.VnodeGUID = n.GUID
		jobCfg.VnodeHostname = n.Hostname
		jobCfg.VnodeLabels = n.Labels
	}

	job := module.NewJob(jobCfg)

	return job, nil
}

func detection(job jobpkg.Job) state {
	if !job.AutoDetection() {
		if job.RetryAutoDetection() {
			return retry
		} else {
			return failed
		}
	}
	return success
}

func runRetryTask(ctx context.Context, in chan<- confgroup.Config, cfg confgroup.Config, timeout time.Duration) {
	t := time.NewTimer(timeout)
	defer t.Stop()

	select {
	case <-ctx.Done():
	case <-t.C:
		select {
		case <-ctx.Done():
		case in <- cfg:
		}
	}
}

func unmarshal(conf interface{}, module interface{}) error {
	bs, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bs, module)
}

func isInsideK8sCluster() bool {
	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
	return host != "" && port != ""
}

func isTooManyOpenFiles(err error) bool {
	return err != nil && strings.Contains(err.Error(), "too many open files")
}

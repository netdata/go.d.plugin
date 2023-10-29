// SPDX-License-Identifier: GPL-3.0-or-later

package sd

import (
	"context"
	"sync"

	"github.com/netdata/go.d.plugin/agent/confgroup"
	"github.com/netdata/go.d.plugin/agent/discovery/sd/pipeline"
	"github.com/netdata/go.d.plugin/logger"

	"gopkg.in/yaml.v2"
)

func NewManager() (*Manager, error) {
	return nil, nil
}

type (
	Manager struct {
		*logger.Logger

		newPipeline func(config pipeline.Config) (sdPipeline, error)
		confProv    ConfigFileProvider

		cache     map[string]uint64
		pipelines map[string]func()
	}
	sdPipeline interface {
		Run(ctx context.Context, in chan<- []*confgroup.Group)
	}
)

func (m *Manager) Run(ctx context.Context, in chan<- []*confgroup.Group) {
	m.Info("instance is started")
	defer m.Info("instance is stopped")
	defer m.cleanup()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() { defer wg.Done(); m.confProv.Run(ctx) }()

	for {
		select {
		case <-ctx.Done():
			return
		case cf := <-m.confProv.Configs():
			if cf.Source == "" {
				continue
			}
			if len(cf.Data) == 0 {
				delete(m.cache, cf.Source)
				m.removePipeline(cf)
			} else if hash, ok := m.cache[cf.Source]; !ok || hash != cf.Hash() {
				m.cache[cf.Source] = cf.Hash()
				m.addPipeline(ctx, cf, in)
			}
		}
	}
}

func (m *Manager) addPipeline(ctx context.Context, cf ConfigFile, in chan<- []*confgroup.Group) {
	var cfg pipeline.Config

	if err := yaml.Unmarshal(cf.Data, &cfg); err != nil {
		m.Error(err)
		return
	}

	pl, err := m.newPipeline(cfg)
	if err != nil {
		m.Error(err)
		return
	}

	if stop, ok := m.pipelines[cf.Source]; ok {
		stop()
	}

	var wg sync.WaitGroup
	plCtx, cancel := context.WithCancel(ctx)

	wg.Add(1)
	go func() { defer wg.Done(); pl.Run(plCtx, in) }()
	stop := func() { cancel(); wg.Wait() }

	m.pipelines[cf.Source] = stop
}

func (m *Manager) removePipeline(cf ConfigFile) {
	if stop, ok := m.pipelines[cf.Source]; ok {
		delete(m.pipelines, cf.Source)
		stop()
	}
}

func (m *Manager) cleanup() {
	for _, stop := range m.pipelines {
		stop()
	}
}

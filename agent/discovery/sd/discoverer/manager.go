// SPDX-License-Identifier: GPL-3.0-or-later

package discoverer

import (
	"context"
	"sync"
	"time"

	"github.com/netdata/go.d.plugin/agent/discovery/sd/discoverer/kubernetes"
	"github.com/netdata/go.d.plugin/agent/discovery/sd/model"

	"github.com/netdata/go.d.plugin/logger"
)

type Discoverer interface {
	Discover(ctx context.Context, in chan<- []model.TargetGroup)
}

func New(cfg Config) (*Manager, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	mgr := &Manager{
		log:         logger.New("discovery manager", ""),
		send:        make(chan struct{}, 1),
		sendEvery:   5 * time.Second,
		discoverers: make([]Discoverer, 0),
		mux:         &sync.Mutex{},
		groups:      make(map[string]model.TargetGroup),
	}

	if err := mgr.registerDiscoverers(cfg); err != nil {
		return nil, err
	}

	mgr.log.Infof("registered: %v", mgr.discoverers)

	return mgr, nil
}

type Manager struct {
	log         *logger.Logger
	discoverers []Discoverer
	send        chan struct{}
	sendEvery   time.Duration
	mux         *sync.Mutex
	groups      map[string]model.TargetGroup
}

func (m *Manager) registerDiscoverers(conf Config) error {
	for _, cfg := range conf.K8S {
		d, err := kubernetes.NewDiscovery(cfg)
		if err != nil {
			return err
		}
		m.discoverers = append(m.discoverers, d)
	}
	return nil
}

func (m *Manager) Discover(ctx context.Context, in chan<- []model.TargetGroup) {
	m.log.Info("instance is started")
	defer m.log.Info("instance is stopped")

	var wg sync.WaitGroup

	for _, d := range m.discoverers {
		wg.Add(1)
		go func(d Discoverer) { defer wg.Done(); m.runDiscoverer(ctx, d) }(d)
	}

	wg.Add(1)
	go func() { defer wg.Done(); m.run(ctx, in) }()

	wg.Wait()
	<-ctx.Done()
}

func (m *Manager) runDiscoverer(ctx context.Context, d Discoverer) {
	updates := make(chan []model.TargetGroup)
	go d.Discover(ctx, updates)

	for {
		select {
		case <-ctx.Done():
			return
		case groups, ok := <-updates:
			if !ok {
				return
			}
			func() {
				m.mux.Lock()
				defer m.mux.Unlock()

				m.groupsUpdate(groups)
				m.triggerSend()
			}()
		}
	}
}

func (m *Manager) run(ctx context.Context, in chan<- []model.TargetGroup) {
	tk := time.NewTicker(m.sendEvery)
	defer tk.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tk.C:
			select {
			case <-m.send:
				m.trySend(in)
			default:
			}
		}
	}
}

func (m *Manager) trySend(in chan<- []model.TargetGroup) {
	m.mux.Lock()
	defer m.mux.Unlock()

	select {
	case in <- m.groupsAsList():
		m.groupsReset()
	default:
		m.triggerSend()
	}
}

func (m *Manager) triggerSend() {
	select {
	case m.send <- struct{}{}:
	default:
	}
}

func (m *Manager) groupsUpdate(groups []model.TargetGroup) {
	for _, group := range groups {
		if group != nil {
			m.groups[group.Source()] = group
		}
	}
}

func (m *Manager) groupsReset() {
	for key := range m.groups {
		delete(m.groups, key)
	}
}

func (m *Manager) groupsAsList() []model.TargetGroup {
	groups := make([]model.TargetGroup, 0, len(m.groups))
	for _, group := range m.groups {
		groups = append(groups, group)
	}
	return groups
}

// SPDX-License-Identifier: GPL-3.0-or-later

package state

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/netdata/go.d.plugin/agent/job/confgroup"
	"github.com/netdata/go.d.plugin/logger"
)

type Manager struct {
	path    string
	store   *Store
	flushCh chan struct{}
	*logger.Logger
}

func NewManager(path string) *Manager {
	return &Manager{
		store:   &Store{},
		path:    path,
		flushCh: make(chan struct{}, 1),
		Logger:  logger.New("state save", "manager"),
	}
}

func (m *Manager) Run(ctx context.Context) {
	m.Info("instance is started")
	defer func() { m.Info("instance is stopped") }()

	tk := time.NewTicker(time.Second * 5)
	defer tk.Stop()
	defer m.flush()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tk.C:
			select {
			case <-m.flushCh:
				m.flush()
			default:
			}
		}
	}
}

func (m *Manager) Save(cfg confgroup.Config, state string) {
	if st, ok := m.store.lookup(cfg); !ok || state != st {
		m.store.add(cfg, state)
		m.triggerFlush()
	}
}

func (m *Manager) Remove(cfg confgroup.Config) {
	if _, ok := m.store.lookup(cfg); ok {
		m.store.remove(cfg)
		m.triggerFlush()
	}
}

func (m *Manager) triggerFlush() {
	select {
	case m.flushCh <- struct{}{}:
	default:
	}
}

func (m *Manager) flush() {
	bs, err := m.store.bytes()
	if err != nil {
		return
	}
	f, err := os.Create(m.path)
	if err != nil {
		return
	}
	defer func() { _ = f.Close() }()
	_, _ = f.Write(bs)
}

type Store struct {
	mux   sync.Mutex
	items map[string]map[string]string // [module][name:hash]state
}

func (s *Store) Contains(cfg confgroup.Config, states ...string) bool {
	state, ok := s.lookup(cfg)
	if !ok {
		return false
	}
	for _, v := range states {
		if state == v {
			return true
		}
	}
	return false
}

func (s *Store) lookup(cfg confgroup.Config) (string, bool) {
	s.mux.Lock()
	defer s.mux.Unlock()

	v, ok := s.items[cfg.Module()]
	if !ok {
		return "", false
	}
	state, ok := v[storeKey(cfg)]
	return state, ok
}

func (s *Store) add(cfg confgroup.Config, state string) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.items == nil {
		s.items = make(map[string]map[string]string)
	}
	if s.items[cfg.Module()] == nil {
		s.items[cfg.Module()] = make(map[string]string)
	}
	s.items[cfg.Module()][storeKey(cfg)] = state
}

func (s *Store) remove(cfg confgroup.Config) {
	s.mux.Lock()
	defer s.mux.Unlock()

	delete(s.items[cfg.Module()], storeKey(cfg))
	if len(s.items[cfg.Module()]) == 0 {
		delete(s.items, cfg.Module())
	}
}

func (s *Store) bytes() ([]byte, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	return json.MarshalIndent(s.items, "", " ")
}

func Load(path string) (*Store, error) {
	var s Store
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	return &s, json.NewDecoder(f).Decode(&s.items)
}

func storeKey(cfg confgroup.Config) string {
	return fmt.Sprintf("%s:%d", cfg.Name(), cfg.Hash())
}

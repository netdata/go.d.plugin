package state

import (
	"context"
	"encoding/json"
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
		store:   &Store{mux: new(sync.Mutex)},
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
	defer f.Close()
	_, _ = f.Write(bs)
}

type Store struct {
	mux *sync.Mutex
	// TODO: we need [module][hash][name]store
	items map[string]map[string]string // [module][name]store
}

func (s Store) Contains(cfg confgroup.Config, states ...string) bool {
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
	state, ok := v[cfg.Name()]
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
	s.items[cfg.Module()][cfg.Name()] = state
}

func (s *Store) remove(cfg confgroup.Config) {
	s.mux.Lock()
	defer s.mux.Unlock()

	delete(s.items[cfg.Module()], cfg.Name())
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
	state := &Store{mux: new(sync.Mutex)}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return state, json.NewDecoder(f).Decode(&state.items)
}

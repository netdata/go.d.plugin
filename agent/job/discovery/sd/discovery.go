package sd

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/netdata/go.d.plugin/agent/job/discovery/sd/kubernetes"
	"github.com/netdata/go.d.plugin/agent/job/discovery/sd/model"
	"github.com/netdata/go.d.plugin/logger"
)

func NewDiscoveryManager(cfg DiscoveryManagerConfig) (*DiscoveryManager, error) {
	if err := validateDiscoveryManagerConfig(cfg); err != nil {
		return nil, err
	}

	d := &DiscoveryManager{
		send:        make(chan struct{}, 1),
		sendEvery:   5 * time.Second,
		discoverers: make([]Discoverer, 0),
		mux:         sync.RWMutex{},
		groups:      make(map[string]model.TargetGroup),
		Logger:      logger.New("discovery", "manager"),
	}

	if err := d.registerDiscoverers(cfg); err != nil {
		return nil, err
	}

	d.Infof("registered: %v", d.discoverers)
	return d, nil
}

type (
	DiscoveryManagerConfig struct {
		K8S []kubernetes.DiscoveryConfig `yaml:"k8s"`
	}
	DiscoveryManager struct {
		*logger.Logger

		discoverers []Discoverer
		send        chan struct{}
		sendEvery   time.Duration

		mux    sync.RWMutex
		groups map[string]model.TargetGroup
	}
)

func (dm *DiscoveryManager) registerDiscoverers(cfg DiscoveryManagerConfig) error {
	for _, cfg := range cfg.K8S {
		d, err := kubernetes.NewDiscovery(cfg)
		if err != nil {
			return err
		}
		dm.discoverers = append(dm.discoverers, d)
	}
	return nil
}

func (dm *DiscoveryManager) Discover(ctx context.Context, in chan<- []model.TargetGroup) {
	dm.Info("instance is started")
	defer dm.Info("instance is stopped")

	var wg sync.WaitGroup

	for _, d := range dm.discoverers {
		wg.Add(1)
		go func(d Discoverer) { defer wg.Done(); dm.runDiscoverer(ctx, d) }(d)
	}

	wg.Add(1)
	go func() { defer wg.Done(); dm.run(ctx, in) }()
	wg.Wait()

	<-ctx.Done()
}

func (dm *DiscoveryManager) runDiscoverer(ctx context.Context, d Discoverer) {
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
				dm.mux.Lock()
				defer dm.mux.Unlock()

				dm.addTargetGroups(groups)
				dm.triggerSend()
			}()
		}
	}
}

func (dm *DiscoveryManager) run(ctx context.Context, in chan<- []model.TargetGroup) {
	tk := time.NewTicker(dm.sendEvery)
	defer tk.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tk.C:
			select {
			case <-dm.send:
				dm.trySend(in)
			default:
			}
		}
	}
}

func (dm *DiscoveryManager) trySend(in chan<- []model.TargetGroup) {
	dm.mux.Lock()
	defer dm.mux.Unlock()

	select {
	case in <- dm.listTargetGroups():
		dm.resetTargetGroups()
	default:
		dm.triggerSend()
	}
}

func (dm *DiscoveryManager) triggerSend() {
	select {
	case dm.send <- struct{}{}:
	default:
	}
}

func (dm *DiscoveryManager) addTargetGroups(groups []model.TargetGroup) {
	for _, group := range groups {
		if group != nil {
			dm.groups[group.Source()] = group
		}
	}
}

func (dm *DiscoveryManager) resetTargetGroups() {
	for key := range dm.groups {
		delete(dm.groups, key)
	}
}

func (dm *DiscoveryManager) listTargetGroups() []model.TargetGroup {
	groups := make([]model.TargetGroup, 0, len(dm.groups))
	for _, group := range dm.groups {
		groups = append(groups, group)
	}
	return groups
}

func validateDiscoveryManagerConfig(cfg DiscoveryManagerConfig) error {
	if len(cfg.K8S) == 0 {
		return errors.New("empty config")
	}
	return nil
}

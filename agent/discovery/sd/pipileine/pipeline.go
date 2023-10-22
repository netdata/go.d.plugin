// SPDX-License-Identifier: GPL-3.0-or-later

package pipileine

import (
	"context"
	"sync"

	"github.com/netdata/go.d.plugin/agent/confgroup"
	"github.com/netdata/go.d.plugin/agent/discovery/sd/kubernetes"
	"github.com/netdata/go.d.plugin/agent/discovery/sd/model"
	"github.com/netdata/go.d.plugin/logger"
)

func New(cfg Config) (*Pipeline, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	p := &Pipeline{
		Logger:      logger.New("sd pipeline", cfg.Name),
		discoverers: make([]accumulateTask, 0),
		items:       make(map[string]map[uint64][]confgroup.Config),
	}

	if err := p.registerDiscoverers(cfg); err != nil {
		return nil, err
	}

	return p, nil
}

type Pipeline struct {
	*logger.Logger

	discoverers []accumulateTask

	clr *classificator
	cmr *composer

	items map[string]map[uint64][]confgroup.Config // [source][targetHash]
}

func (p *Pipeline) registerDiscoverers(conf Config) error {
	for _, cfg := range conf.Discovery.K8s {
		tags, _ := model.ParseTags(cfg.Tags)

		td, err := kubernetes.NewTargetDiscoverer(cfg.Config)
		if err != nil {
			return err
		}

		p.discoverers = append(p.discoverers, accumulateTask{
			disc: td,
			tags: tags,
		})
	}

	return nil
}

func (p *Pipeline) Discover(ctx context.Context, in chan<- []*confgroup.Group) {
	p.Info("instance is started")
	defer p.Info("instance is stopped")

	accum := newAccumulator()
	accum.tasks = p.discoverers

	updates := make(chan []model.TargetGroup)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() { defer wg.Done(); accum.run(ctx, updates) }()

	wg.Add(1)
	go func() { defer wg.Done(); p.run(ctx, updates, in) }()

	wg.Wait()
	<-ctx.Done()
}

func (p *Pipeline) run(ctx context.Context, updates <-chan []model.TargetGroup, in chan<- []*confgroup.Group) {
	for {
		select {
		case <-ctx.Done():
			return
		case groups := <-updates:
			p.Infof("received %d target groups", len(groups))

			var cfgGroups []*confgroup.Group

			// updates come from the accumulator, this ensures that all groups have different sources
			for _, group := range groups {
				p.Infof("processing group '%s' with %d target(s)", group.Source(), len(group.Targets()))

				if v := p.processGroup(group); v != nil {
					cfgGroups = append(cfgGroups, v)
				}
			}

			if len(cfgGroups) > 0 {
				send(ctx, in, cfgGroups)
			}
		}
	}
}

func (p *Pipeline) processGroup(tgg model.TargetGroup) *confgroup.Group {
	if len(tgg.Targets()) == 0 {
		if _, ok := p.items[tgg.Source()]; !ok {
			return nil
		}
		delete(p.items, tgg.Source())
		return &confgroup.Group{Source: tgg.Source()}
	}

	targetsCache, ok := p.items[tgg.Source()]
	if !ok {
		targetsCache = make(map[uint64][]confgroup.Config)
		p.items[tgg.Source()] = targetsCache
	}

	var changed bool
	seen := make(map[uint64]bool)

	for _, tgt := range tgg.Targets() {
		if tgt == nil {
			continue
		}

		hash := tgt.Hash()
		seen[hash] = true

		if _, ok := targetsCache[hash]; ok {
			continue
		}

		if tags := p.clr.classify(tgt); len(tags) > 0 {
			tgt.Tags().Merge(tags)

			if configs := p.cmr.compose(tgt); len(configs) > 0 {
				for _, cfg := range configs {
					cfg.SetProvider(tgg.Provider())
					cfg.SetSource(tgg.Source())
				}
				targetsCache[hash] = configs
				changed = true
			}
		} else {
			p.Infof("target '%s' classify: fail", tgt.TUID())
		}
	}

	for hash := range targetsCache {
		if seen[hash] {
			continue
		}
		if configs := targetsCache[hash]; len(configs) > 0 {
			changed = true
		}
		delete(targetsCache, hash)
	}

	if !changed {
		return nil
	}

	// TODO: deepcopy?
	cfgGroup := &confgroup.Group{Source: tgg.Source()}
	for _, cfgs := range targetsCache {
		cfgGroup.Configs = append(cfgGroup.Configs, cfgs...)
	}

	return cfgGroup
}

func send(ctx context.Context, in chan<- []*confgroup.Group, configs []*confgroup.Group) {
	select {
	case <-ctx.Done():
		return
	case in <- configs:
	}
}

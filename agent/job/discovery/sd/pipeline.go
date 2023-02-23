package sd

import (
	"context"
	"github.com/netdata/go.d.plugin/agent/job/discovery/sd/model"
	"sync"

	"github.com/netdata/go.d.plugin/logger"
)

type Discoverer interface {
	Discover(ctx context.Context, in chan<- []model.TargetGroup)
}

type Tagger interface {
	Tag(model.Target)
}

type Builder interface {
	Build(model.Target) []model.Config
}

type Exporter interface {
	Export(ctx context.Context, out <-chan []model.Config)
}

type (
	Pipeline struct {
		*logger.Logger
		Discoverer
		Tagger
		Builder
		Exporter

		cache pipelineCache
	}
	pipelineCache map[string]groupCache // source:hash:configs
	groupCache    map[uint64][]model.Config
)

func NewPipeline(discoverer Discoverer, tagger Tagger, builder Builder, exporter Exporter) *Pipeline {
	return &Pipeline{
		Discoverer: discoverer,
		Tagger:     tagger,
		Builder:    builder,
		Exporter:   exporter,
		cache:      make(pipelineCache),
		Logger:     logger.New("discovery", "pipeline"),
	}
}

func (p *Pipeline) Run(ctx context.Context) {
	p.Info("instance is started")
	defer p.Info("instance is stopped")

	var wg sync.WaitGroup
	disc := make(chan []model.TargetGroup)
	exp := make(chan []model.Config)

	wg.Add(1)
	go func() { defer wg.Done(); p.Discover(ctx, disc) }()

	wg.Add(1)
	go func() { defer wg.Done(); p.run(ctx, disc, exp) }()

	wg.Add(1)
	go func() { defer wg.Done(); p.Export(ctx, exp) }()

	wg.Wait()
	<-ctx.Done()
}

func (p *Pipeline) run(ctx context.Context, disc chan []model.TargetGroup, export chan []model.Config) {
	for {
		select {
		case <-ctx.Done():
			return
		case groups := <-disc:
			if configs := p.process(groups); len(configs) > 0 {
				select {
				case <-ctx.Done():
				case export <- configs:
				}
			}
		}
	}
}

func (p *Pipeline) process(groups []model.TargetGroup) (configs []model.Config) {
	p.Infof("received '%d' group(s)", len(groups))

	for _, group := range groups {
		p.Infof("processing group '%s' with %d target(s)", group.Source(), len(group.Targets()))

		if len(group.Targets()) == 0 {
			if remove := p.handleEmpty(group); len(remove) > 0 {
				p.Infof("group '%s': stale config(s) %d", group.Source(), len(remove))

				configs = append(configs, remove...)
			}
		} else {
			if add, remove := p.handleNotEmpty(group); len(add) > 0 || len(remove) > 0 {
				p.Infof("group '%s': new/stale config(s) %d/%d", group.Source(), len(add), len(remove))

				configs = append(configs, append(add, remove...)...)
			}
		}
	}
	return configs
}

func (p *Pipeline) handleEmpty(group model.TargetGroup) (remove []model.Config) {
	grpCache, exist := p.cache[group.Source()]
	if !exist {
		return
	}
	delete(p.cache, group.Source())

	for hash, cfgs := range grpCache {
		delete(grpCache, hash)
		remove = append(remove, cfgs...)
	}

	return stale(remove)
}

func (p *Pipeline) handleNotEmpty(group model.TargetGroup) (add, remove []model.Config) {
	grpCache, exist := p.cache[group.Source()]
	if !exist {
		grpCache = make(map[uint64][]model.Config)
		p.cache[group.Source()] = grpCache
	}

	seen := make(map[uint64]bool)
	for _, target := range group.Targets() {
		if target == nil {
			continue
		}
		seen[target.Hash()] = true

		if _, ok := grpCache[target.Hash()]; ok {
			continue
		}

		p.Tag(target)
		cfgs := p.Build(target)

		grpCache[target.Hash()] = cfgs
		add = append(add, cfgs...)
	}

	if !exist {
		return
	}

	for hash, cfgs := range grpCache {
		if !seen[hash] {
			delete(grpCache, hash)
			remove = append(remove, stale(cfgs)...)
		}
	}
	return add, remove
}

func stale(configs []model.Config) []model.Config {
	for i := range configs {
		configs[i].Stale = true
	}
	return configs
}

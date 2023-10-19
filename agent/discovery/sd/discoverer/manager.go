// SPDX-License-Identifier: GPL-3.0-or-later

package discoverer

import (
	"bytes"
	"context"
	"strings"
	"sync"
	"text/template"

	"github.com/netdata/go.d.plugin/agent/confgroup"
	"github.com/netdata/go.d.plugin/agent/discovery/sd/discoverer/kubernetes"
	"github.com/netdata/go.d.plugin/agent/discovery/sd/model"
	"github.com/netdata/go.d.plugin/logger"

	"gopkg.in/yaml.v2"
)

type Discoverer interface {
	Discover(ctx context.Context, in chan<- []model.TargetGroup)
}

func New(cfg Config) (*Manager, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	mgr := &Manager{
		Logger:      logger.New("discovery manager", ""),
		discoverers: make([]Discoverer, 0),
		accum:       newAccumulator(),
		items:       make(map[string]map[uint64][]confgroup.Config),
	}

	if err := mgr.registerDiscoverers(cfg); err != nil {
		return nil, err
	}

	mgr.Infof("registered: %v", mgr.discoverers)

	return mgr, nil
}

type (
	Manager struct {
		*logger.Logger

		discoverers []Discoverer
		accum       *accumulator

		tagRules   []*tagRule
		buildRules []*buildRule
		buf        bytes.Buffer

		items map[string]map[uint64][]confgroup.Config // [source][targetHash]
	}

	tagRule struct {
		id    int
		name  string
		sr    selector
		tags  model.Tags
		match []*tagRuleMatch
	}
	tagRuleMatch struct {
		id   int
		sr   selector
		tags model.Tags
		expr *template.Template
	}

	buildRule struct {
		name  string
		id    int
		sr    selector
		tags  model.Tags
		apply []*ruleApply
	}
	ruleApply struct {
		id   int
		sr   selector
		tags model.Tags
		tmpl *template.Template
	}
)

func (m *Manager) registerDiscoverers(conf Config) error {
	for _, cfg := range conf.Discovery.K8S {
		d, err := kubernetes.NewDiscovery(cfg)
		if err != nil {
			return err
		}
		m.discoverers = append(m.discoverers, d)
	}
	return nil
}

func (m *Manager) Discover(ctx context.Context, in chan<- []*confgroup.Group) {
	m.Info("instance is started")
	defer m.Info("instance is stopped")

	m.accum.discoverers = m.discoverers

	updates := make(chan []model.TargetGroup)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() { defer wg.Done(); m.accum.run(ctx, updates) }()

	wg.Add(1)
	go func() { defer wg.Done(); m.run(ctx, updates, in) }()

	wg.Wait()
	<-ctx.Done()
}

func (m *Manager) run(ctx context.Context, updates <-chan []model.TargetGroup, in chan<- []*confgroup.Group) {
	for {
		select {
		case <-ctx.Done():
			return
		case groups := <-updates:
			m.Infof("received %d target groups", len(groups))

			var cfgGroups []*confgroup.Group

			// updates come from the accumulator, this ensures that all groups have different sources
			for _, group := range groups {
				m.Infof("processing group '%s' with %d target(s)", group.Source(), len(group.Targets()))

				if v := m.processGroup(group); v != nil {
					cfgGroups = append(cfgGroups, v)
				}
			}

			if len(cfgGroups) > 0 {
				send(ctx, in, cfgGroups)
			}
		}
	}
}

func (m *Manager) processGroup(group model.TargetGroup) *confgroup.Group {
	if len(group.Targets()) == 0 {
		if _, ok := m.items[group.Source()]; !ok {
			return nil
		}
		delete(m.items, group.Source())
		return &confgroup.Group{Source: group.Source()}
	}

	targetsCache, ok := m.items[group.Source()]
	if !ok {
		targetsCache = make(map[uint64][]confgroup.Config)
		m.items[group.Source()] = targetsCache
	}

	var changed bool
	seen := make(map[uint64]bool)

	for _, tgt := range group.Targets() {
		if tgt == nil {
			continue
		}

		hash := tgt.Hash()
		seen[hash] = true

		if _, ok := targetsCache[hash]; ok {
			continue
		}

		m.tag(tgt)
		configs := m.createConfigs(group, tgt)
		targetsCache[hash] = configs
		changed = true
	}

	for hash := range targetsCache {
		if seen[hash] {
			continue
		}
		delete(targetsCache, hash)
		changed = true
	}

	if !changed {
		return nil
	}

	// TODO: deepcopy?
	cfgGroup := &confgroup.Group{Source: group.Source()}
	for _, cfgs := range targetsCache {
		cfgGroup.Configs = append(cfgGroup.Configs, cfgs...)
	}

	return cfgGroup
}

func (m *Manager) tag(tgt model.Target) {
	for _, rule := range m.tagRules {
		if !rule.sr.matches(tgt.Tags()) {
			continue
		}

		for _, match := range rule.match {
			if !match.sr.matches(tgt.Tags()) {
				continue
			}

			m.buf.Reset()

			if err := match.expr.Execute(&m.buf, tgt); err != nil {
				m.Warningf("failed to execute rule match '%d/%d' on target '%s'", rule.id, match.id, tgt.TUID())
				continue
			}
			if strings.TrimSpace(m.buf.String()) != "true" {
				continue
			}

			tgt.Tags().Merge(rule.tags)
			tgt.Tags().Merge(match.tags)

			m.Debugf("matched target '%s', tags: %s", tgt.TUID(), tgt.Tags())
		}
	}
}

func (m *Manager) createConfigs(group model.TargetGroup, tgt model.Target) []confgroup.Config {
	var configs []confgroup.Config

	for _, rule := range m.buildRules {
		if !rule.sr.matches(tgt.Tags()) {
			continue
		}

		for _, apply := range rule.apply {
			if !apply.sr.matches(tgt.Tags()) {
				continue
			}

			m.buf.Reset()

			if err := apply.tmpl.Execute(&m.buf, tgt); err != nil {
				m.Warningf("failed to execute rule apply '%d/%d' on target '%s'", rule.id, apply.id, tgt.TUID())
				continue
			}
			if m.buf.Len() == 0 {
				continue
			}

			var cfg confgroup.Config

			if err := yaml.Unmarshal(m.buf.Bytes(), &cfg); err != nil {
				m.Warningf("failed on yaml unmarshalling: %v", err)
				continue
			}

			cfg.SetProvider(group.Provider())
			cfg.SetSource(group.Source())

			configs = append(configs, cfg)
		}
	}

	if len(configs) > 0 {
		m.Infof("created %d config(s) for target '%s'", len(configs), tgt.TUID())
	}
	return configs
}

func send(ctx context.Context, in chan<- []*confgroup.Group, configs []*confgroup.Group) {
	select {
	case <-ctx.Done():
		return
	case in <- configs:
	}
}

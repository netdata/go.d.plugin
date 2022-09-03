// SPDX-License-Identifier: GPL-3.0-or-later

package build

import (
	"context"

	"github.com/netdata/go.d.plugin/agent/job/confgroup"
)

type (
	fullName  = string
	grpSource = string
	cfgHash   = uint64
	cfgCount  = uint

	startedCache map[fullName]struct{}
	retryCache   map[cfgHash]retryTask
	groupCache   struct {
		global map[cfgHash]cfgCount
		source map[grpSource]map[cfgHash]confgroup.Config
	}
	retryTask struct {
		cancel  context.CancelFunc
		timeout int
		retries int
	}
)

func newStartedCache() *startedCache {
	return &startedCache{}
}

func newRetryCache() *retryCache {
	return &retryCache{}
}

func newGroupCache() *groupCache {
	return &groupCache{
		global: make(map[cfgHash]cfgCount),
		source: make(map[grpSource]map[cfgHash]confgroup.Config),
	}
}

func (c startedCache) put(cfg confgroup.Config) {
	c[cfg.FullName()] = struct{}{}
}
func (c startedCache) remove(cfg confgroup.Config) {
	delete(c, cfg.FullName())
}
func (c startedCache) has(cfg confgroup.Config) bool {
	_, ok := c[cfg.FullName()]
	return ok
}

func (c retryCache) put(cfg confgroup.Config, retry retryTask) {
	c[cfg.Hash()] = retry
}
func (c retryCache) remove(cfg confgroup.Config) {
	delete(c, cfg.Hash())
}
func (c retryCache) lookup(cfg confgroup.Config) (retryTask, bool) {
	v, ok := c[cfg.Hash()]
	return v, ok
}

func (c *groupCache) put(group *confgroup.Group) (added, removed []confgroup.Config) {
	if group == nil {
		return
	}
	if len(group.Configs) == 0 {
		return c.putEmpty(group)
	}
	return c.putNotEmpty(group)
}

func (c *groupCache) putEmpty(group *confgroup.Group) (added, removed []confgroup.Config) {
	set, ok := c.source[group.Source]
	if !ok {
		return nil, nil
	}

	for hash, cfg := range set {
		c.global[hash]--
		if c.global[hash] == 0 {
			removed = append(removed, cfg)
		}
		delete(set, hash)
	}
	delete(c.source, group.Source)
	return nil, removed
}

func (c *groupCache) putNotEmpty(group *confgroup.Group) (added, removed []confgroup.Config) {
	set, ok := c.source[group.Source]
	if !ok {
		set = make(map[cfgHash]confgroup.Config)
		c.source[group.Source] = set
	}

	seen := make(map[uint64]struct{})

	for _, cfg := range group.Configs {
		hash := cfg.Hash()
		seen[hash] = struct{}{}

		if _, ok := set[hash]; ok {
			continue
		}

		set[hash] = cfg
		if c.global[hash] == 0 {
			added = append(added, cfg)
		}
		c.global[hash]++
	}

	if !ok {
		return added, nil
	}

	for hash, cfg := range set {
		if _, ok := seen[hash]; ok {
			continue
		}

		delete(set, hash)
		c.global[hash]--
		if c.global[hash] == 0 {
			removed = append(removed, cfg)
		}
	}

	if ok && len(set) == 0 {
		delete(c.source, group.Source)
	}

	return added, removed
}

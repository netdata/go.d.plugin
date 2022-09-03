// SPDX-License-Identifier: GPL-3.0-or-later

package build

import (
	"sort"
	"testing"

	"github.com/netdata/go.d.plugin/agent/job/confgroup"

	"github.com/stretchr/testify/assert"
)

func TestJobCache_put(t *testing.T) {
	tests := map[string]struct {
		prepareGroups  []confgroup.Group
		groups         []confgroup.Group
		expectedAdd    []confgroup.Config
		expectedRemove []confgroup.Config
	}{
		"new group, new configs": {
			groups: []confgroup.Group{
				prepareGroup("source", prepareCfg("name", "module")),
			},
			expectedAdd: []confgroup.Config{
				prepareCfg("name", "module"),
			},
		},
		"several equal updates for the same group": {
			groups: []confgroup.Group{
				prepareGroup("source", prepareCfg("name", "module")),
				prepareGroup("source", prepareCfg("name", "module")),
				prepareGroup("source", prepareCfg("name", "module")),
				prepareGroup("source", prepareCfg("name", "module")),
				prepareGroup("source", prepareCfg("name", "module")),
			},
			expectedAdd: []confgroup.Config{
				prepareCfg("name", "module"),
			},
		},
		"empty group update for cached group": {
			prepareGroups: []confgroup.Group{
				prepareGroup("source", prepareCfg("name1", "module"), prepareCfg("name2", "module")),
			},
			groups: []confgroup.Group{
				prepareGroup("source"),
			},
			expectedRemove: []confgroup.Config{
				prepareCfg("name1", "module"),
				prepareCfg("name2", "module"),
			},
		},
		"changed group update for cached group": {
			prepareGroups: []confgroup.Group{
				prepareGroup("source", prepareCfg("name1", "module"), prepareCfg("name2", "module")),
			},
			groups: []confgroup.Group{
				prepareGroup("source", prepareCfg("name2", "module")),
			},
			expectedRemove: []confgroup.Config{
				prepareCfg("name1", "module"),
			},
		},
		"empty group update for uncached group": {
			groups: []confgroup.Group{
				prepareGroup("source"),
				prepareGroup("source"),
			},
		},
		"several updates with different source but same context": {
			groups: []confgroup.Group{
				prepareGroup("source1", prepareCfg("name1", "module"), prepareCfg("name2", "module")),
				prepareGroup("source2", prepareCfg("name1", "module"), prepareCfg("name2", "module")),
			},
			expectedAdd: []confgroup.Config{
				prepareCfg("name1", "module"),
				prepareCfg("name2", "module"),
			},
		},
		"have equal configs from 2 sources, get empty group for the 1st source": {
			prepareGroups: []confgroup.Group{
				prepareGroup("source1", prepareCfg("name1", "module"), prepareCfg("name2", "module")),
				prepareGroup("source2", prepareCfg("name1", "module"), prepareCfg("name2", "module")),
			},
			groups: []confgroup.Group{
				prepareGroup("source2"),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cache := newGroupCache()

			for _, group := range test.prepareGroups {
				cache.put(&group)
			}

			var added, removed []confgroup.Config
			for _, group := range test.groups {
				a, r := cache.put(&group)
				added = append(added, a...)
				removed = append(removed, r...)
			}

			sortConfigs(added)
			sortConfigs(removed)
			sortConfigs(test.expectedAdd)
			sortConfigs(test.expectedRemove)

			assert.Equalf(t, test.expectedAdd, added, "added configs")
			assert.Equalf(t, test.expectedRemove, removed, "removed configs")
		})
	}
}

func prepareGroup(source string, cfgs ...confgroup.Config) confgroup.Group {
	return confgroup.Group{
		Configs: cfgs,
		Source:  source,
	}
}

func prepareCfg(name, module string) confgroup.Config {
	return confgroup.Config{
		"name":   name,
		"module": module,
	}
}

func sortConfigs(cfgs []confgroup.Config) {
	if len(cfgs) == 0 {
		return
	}
	sort.Slice(cfgs, func(i, j int) bool { return cfgs[i].FullName() < cfgs[j].FullName() })
}

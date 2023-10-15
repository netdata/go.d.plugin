// SPDX-License-Identifier: GPL-3.0-or-later

package discoverer

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/netdata/go.d.plugin/agent/discovery/sd/model"
)

func TestNew(t *testing.T) {

}

func TestManager_Discover(t *testing.T) {
	tests := map[string]func() discoverySim{
		"2 discoverers unique groups with delayed collect": func() discoverySim {
			const numGroups, numTargets = 2, 2
			d1 := prepareMockDiscoverer("test1", numGroups, numTargets)
			d2 := prepareMockDiscoverer("test2", numGroups, numTargets)
			mgr := prepareManager(d1, d2)
			expected := combineGroups(d1.groups, d2.groups)

			sim := discoverySim{
				mgr:            mgr,
				collectDelay:   mgr.sendEvery + time.Second,
				expectedGroups: expected,
			}
			return sim
		},
		"2 discoverers unique groups": func() discoverySim {
			const numGroups, numTargets = 2, 2
			d1 := prepareMockDiscoverer("test1", numGroups, numTargets)
			d2 := prepareMockDiscoverer("test2", numGroups, numTargets)
			mgr := prepareManager(d1, d2)
			expected := combineGroups(d1.groups, d2.groups)

			sim := discoverySim{
				mgr:            mgr,
				expectedGroups: expected,
			}
			return sim
		},
		"2 discoverers same groups": func() discoverySim {
			const numGroups, numTargets = 2, 2
			d1 := prepareMockDiscoverer("test1", numGroups, numTargets)
			mgr := prepareManager(d1, d1)
			expected := combineGroups(d1.groups)

			sim := discoverySim{
				mgr:            mgr,
				expectedGroups: expected,
			}
			return sim
		},
		"2 discoverers empty groups": func() discoverySim {
			const numGroups, numTargets = 1, 0
			d1 := prepareMockDiscoverer("test1", numGroups, numTargets)
			d2 := prepareMockDiscoverer("test2", numGroups, numTargets)
			mgr := prepareManager(d1, d2)
			expected := combineGroups(d1.groups, d2.groups)

			sim := discoverySim{
				mgr:            mgr,
				expectedGroups: expected,
			}
			return sim
		},
		"2 discoverers nil groups": func() discoverySim {
			const numGroups, numTargets = 0, 0
			d1 := prepareMockDiscoverer("test1", numGroups, numTargets)
			d2 := prepareMockDiscoverer("test2", numGroups, numTargets)
			mgr := prepareManager(d1, d2)

			sim := discoverySim{
				mgr:            mgr,
				expectedGroups: nil,
			}
			return sim
		},
	}

	for name, sim := range tests {
		t.Run(name, func(t *testing.T) { sim().run(t) })
	}
}

func prepareMockDiscoverer(source string, groups, targets int) mockDiscoverer {
	d := mockDiscoverer{}

	for i := 0; i < groups; i++ {
		group := mockGroup{
			source: fmt.Sprintf("%s_group_%d", source, i+1),
		}
		for j := 0; j < targets; j++ {
			group.targets = append(group.targets,
				mockTarget{Name: fmt.Sprintf("%s_group_%d_target_%d", source, i+1, j+1)})
		}
		d.groups = append(d.groups, group)
	}
	return d
}

func prepareManager(discoverers ...Discoverer) *Manager {
	mgr := &Manager{
		send:        make(chan struct{}, 1),
		sendEvery:   2 * time.Second,
		discoverers: discoverers,
		mux:         &sync.Mutex{},
		groups:      make(map[string]model.TargetGroup),
	}
	return mgr
}

type mockDiscoverer struct {
	groups []model.TargetGroup
}

func (md mockDiscoverer) Discover(ctx context.Context, out chan<- []model.TargetGroup) {
	for {
		select {
		case <-ctx.Done():
			return
		case out <- md.groups:
			return
		}
	}
}

type mockGroup struct {
	targets []model.Target
	source  string
}

func (mg mockGroup) Targets() []model.Target { return mg.targets }
func (mg mockGroup) Source() string          { return mg.source }

type mockTarget struct {
	Name string
}

func (mt mockTarget) Tags() model.Tags { return model.Tags{} }
func (mt mockTarget) TUID() string     { return "" }
func (mt mockTarget) Hash() uint64     { return 0 }

func combineGroups(groups ...[]model.TargetGroup) (combined []model.TargetGroup) {
	for _, set := range groups {
		combined = append(combined, set...)
	}
	return combined
}

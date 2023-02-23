package sd

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/netdata/go.d.plugin/agent/job/discovery/sd/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscoveryManager_Discover(t *testing.T) {
	tests := map[string]func() kubernetesDiscoverySim{
		"2 discoverers unique groups with delayed collect": func() kubernetesDiscoverySim {
			const numGroups, numTargets = 2, 2
			d1 := prepareMockDiscoverer("test1", numGroups, numTargets)
			d2 := prepareMockDiscoverer("test2", numGroups, numTargets)
			mgr := prepareKubernetesDiscoverer(d1, d2)
			expected := combineGroups(d1.groups, d2.groups)

			sim := kubernetesDiscoverySim{
				mgr:            mgr,
				collectDelay:   mgr.sendEvery + time.Second,
				expectedGroups: expected,
			}
			return sim
		},
		"2 discoverers unique groups": func() kubernetesDiscoverySim {
			const numGroups, numTargets = 2, 2
			d1 := prepareMockDiscoverer("test1", numGroups, numTargets)
			d2 := prepareMockDiscoverer("test2", numGroups, numTargets)
			mgr := prepareKubernetesDiscoverer(d1, d2)
			expected := combineGroups(d1.groups, d2.groups)

			sim := kubernetesDiscoverySim{
				mgr:            mgr,
				expectedGroups: expected,
			}
			return sim
		},
		"2 discoverers same groups": func() kubernetesDiscoverySim {
			const numGroups, numTargets = 2, 2
			d1 := prepareMockDiscoverer("test1", numGroups, numTargets)
			mgr := prepareKubernetesDiscoverer(d1, d1)
			expected := combineGroups(d1.groups)

			sim := kubernetesDiscoverySim{
				mgr:            mgr,
				expectedGroups: expected,
			}
			return sim
		},
		"2 discoverers empty groups": func() kubernetesDiscoverySim {
			const numGroups, numTargets = 1, 0
			d1 := prepareMockDiscoverer("test1", numGroups, numTargets)
			d2 := prepareMockDiscoverer("test2", numGroups, numTargets)
			mgr := prepareKubernetesDiscoverer(d1, d2)
			expected := combineGroups(d1.groups, d2.groups)

			sim := kubernetesDiscoverySim{
				mgr:            mgr,
				expectedGroups: expected,
			}
			return sim
		},
		"2 discoverers nil groups": func() kubernetesDiscoverySim {
			const numGroups, numTargets = 0, 0
			d1 := prepareMockDiscoverer("test1", numGroups, numTargets)
			d2 := prepareMockDiscoverer("test2", numGroups, numTargets)
			mgr := prepareKubernetesDiscoverer(d1, d2)

			sim := kubernetesDiscoverySim{
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

func prepareKubernetesDiscoverer(discoverers ...Discoverer) *DiscoveryManager {
	mgr := &DiscoveryManager{
		send:        make(chan struct{}, 1),
		sendEvery:   2 * time.Second,
		discoverers: discoverers,
		mux:         sync.RWMutex{},
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

type kubernetesDiscoverySim struct {
	mgr            *DiscoveryManager
	collectDelay   time.Duration
	expectedGroups []model.TargetGroup
}

func (sim kubernetesDiscoverySim) run(t *testing.T) {
	t.Helper()
	require.NotNil(t, sim.mgr)

	in, out := make(chan []model.TargetGroup), make(chan []model.TargetGroup)
	go sim.collectGroups(t, in, out)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go sim.mgr.Discover(ctx, in)

	actualGroups := <-out

	sortGroups(sim.expectedGroups)
	sortGroups(actualGroups)

	assert.Equal(t, sim.expectedGroups, actualGroups)
}

func (sim kubernetesDiscoverySim) collectGroups(t *testing.T, in, out chan []model.TargetGroup) {
	time.Sleep(sim.collectDelay)

	timeout := sim.mgr.sendEvery + time.Second*2
	var groups []model.TargetGroup
loop:
	for {
		select {
		case inGroups := <-in:
			if groups = append(groups, inGroups...); len(groups) >= len(sim.expectedGroups) {
				break loop
			}
		case <-time.After(timeout):
			t.Logf("discovery %s timed out after %s, got %d groups, expected %d, some events are skipped",
				sim.mgr.discoverers, timeout, len(groups), len(sim.expectedGroups))
			break loop
		}
	}
	out <- groups
}

func sortGroups(groups []model.TargetGroup) {
	sort.Slice(groups, func(i, j int) bool { return groups[i].Source() < groups[j].Source() })
}

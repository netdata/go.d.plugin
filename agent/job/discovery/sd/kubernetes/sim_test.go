package kubernetes

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/netdata/go.d.plugin/agent/job/discovery/sd/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/tools/cache"
)

const (
	startWaitTimeout  = time.Second * 3
	finishWaitTimeout = time.Second * 5
)

type discoverySim struct {
	discovery        *Discovery
	runAfterSync     func(ctx context.Context)
	sortBeforeVerify bool
	expectedGroups   []model.TargetGroup
}

func (sim discoverySim) run(t *testing.T) []model.TargetGroup {
	t.Helper()
	require.NotNil(t, sim.discovery)
	require.NotEmpty(t, sim.expectedGroups)

	in, out := make(chan []model.TargetGroup), make(chan []model.TargetGroup)
	go sim.collectGroups(t, in, out)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	go sim.discovery.Discover(ctx, in)

	select {
	case <-sim.discovery.started:
	case <-time.After(startWaitTimeout):
		t.Fatalf("discovery %s filed to start in %s", sim.discovery.discoverers, startWaitTimeout)
	}

	synced := cache.WaitForCacheSync(ctx.Done(), sim.discovery.hasSynced)
	require.Truef(t, synced, "discovery %s failed to sync", sim.discovery.discoverers)

	if sim.runAfterSync != nil {
		sim.runAfterSync(ctx)
	}

	groups := <-out

	if sim.sortBeforeVerify {
		sortGroups(groups)
	}

	sim.verifyResult(t, groups)
	return groups
}

func (sim discoverySim) collectGroups(t *testing.T, in, out chan []model.TargetGroup) {
	var groups []model.TargetGroup
loop:
	for {
		select {
		case inGroups := <-in:
			if groups = append(groups, inGroups...); len(groups) >= len(sim.expectedGroups) {
				break loop
			}
		case <-time.After(finishWaitTimeout):
			t.Logf("discovery %s timed out after %s, got %d groups, expected %d, some events are skipped",
				sim.discovery.discoverers, finishWaitTimeout, len(groups), len(sim.expectedGroups))
			break loop
		}
	}
	out <- groups
}

func (sim discoverySim) verifyResult(t *testing.T, result []model.TargetGroup) {
	var expected, actual any

	if len(sim.expectedGroups) == len(result) {
		expected = sim.expectedGroups
		actual = result
	} else {
		want := make(map[string]model.TargetGroup)
		for _, group := range sim.expectedGroups {
			want[group.Source()] = group
		}
		got := make(map[string]model.TargetGroup)
		for _, group := range result {
			got[group.Source()] = group
		}
		expected, actual = want, got
	}

	assert.Equal(t, expected, actual)
}

type hasSynced interface {
	hasSynced() bool
}

var (
	_ hasSynced = &Discovery{}
	_ hasSynced = &PodDiscovery{}
	_ hasSynced = &ServiceDiscovery{}
)

func (d *Discovery) hasSynced() bool {
	for _, dd := range d.discoverers {
		v, ok := dd.(hasSynced)
		if !ok || !v.hasSynced() {
			return false
		}
	}
	return true
}

func (pd *PodDiscovery) hasSynced() bool {
	return pd.podInformer.HasSynced() && pd.cmapInformer.HasSynced() && pd.secretInformer.HasSynced()
}

func (sd *ServiceDiscovery) hasSynced() bool {
	return sd.informer.HasSynced()
}

func sortGroups(groups []model.TargetGroup) {
	sort.Slice(groups, func(i, j int) bool { return groups[i].Source() < groups[j].Source() })
}

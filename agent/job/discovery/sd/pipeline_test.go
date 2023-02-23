package sd

import (
	"context"
	"github.com/netdata/go.d.plugin/agent/job/discovery/sd/model"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/ilyam8/hashstructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPipeline_Run(t *testing.T) {
	tests := map[string]func() pipelineSim{
		"new group with no targets": func() pipelineSim {
			g1 := mockPipelineGroup{source: "s1"}

			sim := pipelineSim{
				discoveredGroups:   []model.TargetGroup{g1},
				expectedCacheItems: 0,
			}
			return sim
		},
		"new group with targets": func() pipelineSim {
			t1 := mockPipelineTarget{Name: "t1"}
			t2 := mockPipelineTarget{Name: "t2"}
			g1 := mockPipelineGroup{targets: []model.Target{t1, t2}, source: "s1"}

			sim := pipelineSim{
				discoveredGroups:   []model.TargetGroup{g1},
				expectedTag:        []model.Target{t1, t2},
				expectedBuild:      []model.Target{t1, t2},
				expectedExport:     []model.Config{{Conf: "t1"}, {Conf: "t2"}},
				expectedCacheItems: 1,
			}
			return sim
		},
		"existing group with same targets": func() pipelineSim {
			t1 := mockPipelineTarget{Name: "t1"}
			t2 := mockPipelineTarget{Name: "t2"}
			g1 := mockPipelineGroup{targets: []model.Target{t1, t2}, source: "s1"}

			sim := pipelineSim{
				discoveredGroups:   []model.TargetGroup{g1, g1},
				expectedTag:        []model.Target{t1, t2},
				expectedBuild:      []model.Target{t1, t2},
				expectedExport:     []model.Config{{Conf: "t1"}, {Conf: "t2"}},
				expectedCacheItems: 1,
			}
			return sim
		},
		"existing group with no targets": func() pipelineSim {
			t1 := mockPipelineTarget{Name: "t1"}
			t2 := mockPipelineTarget{Name: "t2"}
			g1 := mockPipelineGroup{targets: []model.Target{t1, t2}, source: "s1"}
			g2 := mockPipelineGroup{source: "s1"}

			sim := pipelineSim{
				discoveredGroups: []model.TargetGroup{g1, g2},
				expectedTag:      []model.Target{t1, t2},
				expectedBuild:    []model.Target{t1, t2},
				expectedExport: []model.Config{
					{Conf: "t1"}, {Conf: "t2"}, {Conf: "t1", Stale: true}, {Conf: "t2", Stale: true},
				},
				expectedCacheItems: 0,
			}
			return sim
		},
		"existing group with old and new targets": func() pipelineSim {
			t1 := mockPipelineTarget{Name: "t1"}
			t2 := mockPipelineTarget{Name: "t2"}
			t3 := mockPipelineTarget{Name: "t3"}
			g1 := mockPipelineGroup{targets: []model.Target{t1, t2}, source: "s1"}
			g2 := mockPipelineGroup{targets: []model.Target{t1, t3}, source: "s1"}

			sim := pipelineSim{
				discoveredGroups: []model.TargetGroup{g1, g2},
				expectedTag:      []model.Target{t1, t2, t3},
				expectedBuild:    []model.Target{t1, t2, t3},
				expectedExport: []model.Config{
					{Conf: "t1"}, {Conf: "t2"}, {Conf: "t3"}, {Conf: "t2", Stale: true}},
				expectedCacheItems: 1,
			}
			return sim
		},
		"existing group with new targets only": func() pipelineSim {
			t1 := mockPipelineTarget{Name: "t1"}
			t2 := mockPipelineTarget{Name: "t2"}
			t3 := mockPipelineTarget{Name: "t3"}
			g1 := mockPipelineGroup{targets: []model.Target{t1, t2}, source: "s1"}
			g2 := mockPipelineGroup{targets: []model.Target{t3}, source: "s1"}

			sim := pipelineSim{
				discoveredGroups: []model.TargetGroup{g1, g2},
				expectedTag:      []model.Target{t1, t2, t3},
				expectedBuild:    []model.Target{t1, t2, t3},
				expectedExport: []model.Config{
					{Conf: "t1"}, {Conf: "t2"}, {Conf: "t3"},
					{Conf: "t1", Stale: true}, {Conf: "t2", Stale: true}},
				expectedCacheItems: 1,
			}
			return sim
		},
	}

	for name, sim := range tests {
		t.Run(name, func(t *testing.T) { sim().run(t) })
	}
}

type (
	mockPipelineDiscoverer struct {
		send []model.TargetGroup
	}
	mockTagger struct {
		seen []model.Target
	}
	mockBuilder struct {
		seen []model.Target
	}
	mockExporter struct {
		seen []model.Config
	}
)

func (d mockPipelineDiscoverer) Discover(ctx context.Context, in chan<- []model.TargetGroup) {
	select {
	case <-ctx.Done():
	case in <- d.send:
	}
	<-ctx.Done()
}

func (t *mockTagger) Tag(target model.Target) {
	t.seen = append(t.seen, target)
}

func (b *mockBuilder) Build(target model.Target) []model.Config {
	b.seen = append(b.seen, target)
	return []model.Config{{Conf: target.TUID()}}
}

func (e *mockExporter) Export(ctx context.Context, out <-chan []model.Config) {
	select {
	case <-ctx.Done():
	case cfgs := <-out:
		e.seen = append(e.seen, cfgs...)
	}
	<-ctx.Done()
}

type (
	mockPipelineGroup struct {
		targets []model.Target
		source  string
	}
	mockPipelineTarget struct {
		Name string
	}
)

func (mg mockPipelineGroup) Targets() []model.Target { return mg.targets }
func (mg mockPipelineGroup) Source() string          { return mg.source }

func (mt mockPipelineTarget) Tags() model.Tags { return nil }
func (mt mockPipelineTarget) TUID() string     { return mt.Name }
func (mt mockPipelineTarget) Hash() uint64     { h, _ := hashstructure.Hash(mt, nil); return h }

type pipelineSim struct {
	discoveredGroups   []model.TargetGroup
	expectedTag        []model.Target
	expectedBuild      []model.Target
	expectedExport     []model.Config
	expectedCacheItems int
}

func (sim pipelineSim) run(t *testing.T) {
	require.NotEmpty(t, sim.discoveredGroups)

	discoverer := &mockPipelineDiscoverer{send: sim.discoveredGroups}
	tagger := &mockTagger{}
	builder := &mockBuilder{}
	exporter := &mockExporter{}

	p := NewPipeline(discoverer, tagger, builder, exporter)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() { defer wg.Done(); p.Run(ctx) }()

	time.Sleep(time.Second)
	cancel()
	wg.Wait()

	sortStaleConfigs(sim.expectedExport)
	sortStaleConfigs(exporter.seen)

	assert.Equal(t, sim.expectedTag, tagger.seen)
	assert.Equal(t, sim.expectedBuild, builder.seen)
	assert.Equal(t, sim.expectedExport, exporter.seen)
	if sim.expectedCacheItems >= 0 {
		assert.Equal(t, sim.expectedCacheItems, len(p.cache))
	}
}

func sortStaleConfigs(cfgs []model.Config) {
	sort.Slice(cfgs, func(i, j int) bool {
		return cfgs[i].Stale && cfgs[j].Stale && cfgs[i].Conf < cfgs[j].Conf
	})
}

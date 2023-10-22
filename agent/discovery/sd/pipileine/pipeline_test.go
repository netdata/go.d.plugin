// SPDX-License-Identifier: GPL-3.0-or-later

package pipileine

import (
	"context"
	"fmt"
	"github.com/netdata/go.d.plugin/agent/confgroup"
	"github.com/netdata/go.d.plugin/agent/discovery/sd/model"
	"github.com/netdata/go.d.plugin/logger"
	"strings"
	"testing"
	"time"

	"github.com/ilyam8/hashstructure"
)

// import (
//
//	"context"
//	"fmt"
//	"testing"
//	"time"
//
//	"github.com/netdata/go.d.plugin/agent/discovery/sd/model"
//
// )
//
// func TestNew(t *testing.T) {
//
// }
//
//	func TestManager_Discover(t *testing.T) {
//		tests := map[string]func() discoverySim{
//			"2 tasks unique groups with delayed collect": func() discoverySim {
//				const numGroups, numTargets = 2, 2
//				d1 := prepareMockDiscoverer("test1", numGroups, numTargets)
//				d2 := prepareMockDiscoverer("test2", numGroups, numTargets)
//				mgr := prepareManager(d1, d2)
//				expected := combineGroups(d1.groups, d2.groups)
//
//				sim := discoverySim{
//					mgr:            mgr,
//					maxRuntime:   mgr.accum.sendEvery + time.Second,
//					wantConfGroups: expected,
//				}
//				return sim
//			},
//			"2 tasks unique groups": func() discoverySim {
//				const numGroups, numTargets = 2, 2
//				d1 := prepareMockDiscoverer("test1", numGroups, numTargets)
//				d2 := prepareMockDiscoverer("test2", numGroups, numTargets)
//				mgr := prepareManager(d1, d2)
//				expected := combineGroups(d1.groups, d2.groups)
//
//				sim := discoverySim{
//					mgr:            mgr,
//					wantConfGroups: expected,
//				}
//				return sim
//			},
//			"2 tasks same groups": func() discoverySim {
//				const numGroups, numTargets = 2, 2
//				d1 := prepareMockDiscoverer("test1", numGroups, numTargets)
//				mgr := prepareManager(d1, d1)
//				expected := combineGroups(d1.groups)
//
//				sim := discoverySim{
//					mgr:            mgr,
//					wantConfGroups: expected,
//				}
//				return sim
//			},
//			"2 tasks empty groups": func() discoverySim {
//				const numGroups, numTargets = 1, 0
//				d1 := prepareMockDiscoverer("test1", numGroups, numTargets)
//				d2 := prepareMockDiscoverer("test2", numGroups, numTargets)
//				mgr := prepareManager(d1, d2)
//				expected := combineGroups(d1.groups, d2.groups)
//
//				sim := discoverySim{
//					mgr:            mgr,
//					wantConfGroups: expected,
//				}
//				return sim
//			},
//			"2 tasks nil groups": func() discoverySim {
//				const numGroups, numTargets = 0, 0
//				d1 := prepareMockDiscoverer("test1", numGroups, numTargets)
//				d2 := prepareMockDiscoverer("test2", numGroups, numTargets)
//				mgr := prepareManager(d1, d2)
//
//				sim := discoverySim{
//					mgr:            mgr,
//					wantConfGroups: nil,
//				}
//				return sim
//			},
//		}
//
//		for name, sim := range tests {
//			t.Run(name, func(t *testing.T) { sim().run(t) })
//		}
//	}
//
//	func prepareMockDiscoverer(source string, groups, targets int) mockDiscoverer {
//		d := mockDiscoverer{}
//
//		for i := 0; i < groups; i++ {
//			group := mockTargetGroup{
//				source: fmt.Sprintf("%s_group_%d", source, i+1),
//			}
//			for j := 0; j < targets; j++ {
//				group.targets = append(group.targets,
//					mockTarget{Name: fmt.Sprintf("%s_group_%d_target_%d", source, i+1, j+1)})
//			}
//			d.groups = append(d.groups, group)
//		}
//		return d
//	}
//
//	func prepareManager(tasks ...Discoverer) *Pipeline {
//		mgr := &Pipeline{
//			accum:       newAccumulator(),
//			tasks: tasks,
//		}
//		return mgr
//	}

func TestNew(t *testing.T) {
	cfg := `
classify:
  - selector: "test"
    tags: "apps"
    match:
      - tags: "name1"
        expr: '{{ eq .Name "name1" }}'
      - tags: "name2"
        expr: '{{ eq .Name "name2" }}'
compose:
  - selector: "apps"
    config:
      - selector: "name1"
        template: |
          module: name1
          name: qqq-{{.TUID}}
      - selector: "name2"
        template: |
          module: name2
          name: www-{{.TUID}}
`

	p := &Pipeline{
		Logger:      logger.New("sd pipeline", "test"),
		discoverers: make([]accumulateTask, 0),
		items:       make(map[string]map[uint64][]confgroup.Config),
	}
	p.discoverers = []accumulateTask{
		{
			disc: &mockDiscoverer{
				tggs: []model.TargetGroup{
					&mockTargetGroup{
						source: "test",
						targets: []model.Target{
							&mockTarget{Name: "name1"},
							&mockTarget{Name: "name2"},
						},
					},
				},
			},
			tags: map[string]struct{}{"test": {}},
		},
	}

	sim := &discoverySim{
		pl:             p,
		config:         cfg,
		maxRuntime:     time.Second * 10,
		wantConfGroups: nil,
	}
	sim.run(t)

}

type mockDiscoverer struct {
	tggs []model.TargetGroup
}

func (md mockDiscoverer) Discover(ctx context.Context, out chan<- []model.TargetGroup) {
	for {
		select {
		case <-ctx.Done():
			return
		case out <- md.tggs:
			return
		}
	}
}

type mockTargetGroup struct {
	targets []model.Target
	source  string
}

func (mg mockTargetGroup) Targets() []model.Target { return mg.targets }
func (mg mockTargetGroup) Source() string          { return mg.source }
func (mg mockTargetGroup) Provider() string        { return "mock" }

func newMockTarget(name string, tags ...string) *mockTarget {
	m := &mockTarget{Name: name}
	v, _ := model.ParseTags(strings.Join(tags, " "))
	m.Tags().Merge(v)
	return m
}

type mockTarget struct {
	model.Base
	Name string
}

func (mt mockTarget) TUID() string { return mt.Name }
func (mt mockTarget) Hash() uint64 { return mustCalcHash(mt.Name) }

func mustParseTags(line string) model.Tags {
	v, err := model.ParseTags(line)
	if err != nil {
		panic(fmt.Sprintf("mustParseTags: %v", err))
	}
	return v
}

func mustCalcHash(obj any) uint64 {
	hash, err := hashstructure.Hash(obj, nil)
	if err != nil {
		panic(fmt.Sprintf("hash calculation: %v", err))
	}
	return hash
}

//
//func combineGroups(groups ...[]model.TargetGroup) (combined []model.TargetGroup) {
//	for _, set := range groups {
//		combined = append(combined, set...)
//	}
//	return combined
//}

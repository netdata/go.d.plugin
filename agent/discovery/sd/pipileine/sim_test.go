// SPDX-License-Identifier: GPL-3.0-or-later

package pipileine

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/netdata/go.d.plugin/agent/confgroup"
	"github.com/stretchr/testify/require"

	"gopkg.in/yaml.v2"
)

type discoverySim struct {
	pl             *Pipeline
	config         string
	maxRuntime     time.Duration
	wantConfGroups []*confgroup.Config
}

func (sim discoverySim) run(t *testing.T) {
	t.Helper()
	require.NotNil(t, sim.pl)

	var cfg Config
	err := yaml.Unmarshal([]byte(sim.config), &cfg)
	require.Nilf(t, err, "cfg unmarshal")

	clr, err := newClassificator(cfg.Classify)
	require.Nilf(t, err, "classify %v", err)

	cmr, err := newComposer(cfg.Compose)
	require.Nilf(t, err, "compose")

	clr.Logger = sim.pl.Logger
	cmr.Logger = sim.pl.Logger

	sim.pl.clr = clr
	sim.pl.cmr = cmr

	in := make(chan []*confgroup.Group)

	ctx, cancel := context.WithCancel(context.Background())
	go sim.pl.Discover(ctx, in)

	confGroups := sim.collectGroups(t, in)
	cancel()

	time.Sleep(time.Second * 3)

	for _, g := range confGroups {
		for _, c := range g.Configs {
			fmt.Println(c)
		}
	}
}

func (sim discoverySim) collectGroups(t *testing.T, in chan []*confgroup.Group) []*confgroup.Group {
	var groups []*confgroup.Group
loop:
	for {
		select {
		case inGroups := <-in:
			if groups = append(groups, inGroups...); len(groups) >= len(sim.wantConfGroups) {
				break loop
			}
		case <-time.After(sim.maxRuntime):
			t.Logf("discovery timed out after %s, got %d groups, expected %d, some events are skipped",
				sim.maxRuntime, len(groups), len(sim.wantConfGroups))
			break loop
		}
	}
	return groups
}

// SPDX-License-Identifier: GPL-3.0-or-later

package jobmgr

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/netdata/go.d.plugin/agent/confgroup"
	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/agent/netdataapi"
	"github.com/netdata/go.d.plugin/agent/safewriter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type runSim struct {
	do func(mgr *Manager)

	wantDiscovered []confgroup.Config
	wantSeen       []seenConfig
	wantExposed    []seenConfig
	wantRunning    []string
	wantDyncfg     string
}

func (s *runSim) run(t *testing.T) {
	t.Helper()

	require.NotNil(t, s.do, "s.do is nil")

	var buf bytes.Buffer
	mgr := New()
	mgr.api = netdataapi.New(safewriter.New(&buf))
	mgr.Modules = prepareMockRegistry()

	done := make(chan struct{})
	grpCh := make(chan []*confgroup.Group)
	ctx, cancel := context.WithCancel(context.Background())

	go func() { defer close(done); close(grpCh); mgr.Run(ctx, grpCh) }()

	timeout := time.Second * 5

	select {
	case <-mgr.started:
	case <-time.After(timeout):
		t.Errorf("failed to start work in %s", timeout)
	}

	s.do(mgr)
	cancel()

	select {
	case <-done:
	case <-time.After(timeout):
		t.Errorf("failed to finish work in %s", timeout)
	}

	var lines []string
	for _, s := range strings.Split(buf.String(), "\n") {
		if strings.HasPrefix(s, "CONFIG") && strings.Contains(s, " template ") {
			continue
		}
		if strings.HasPrefix(s, "FUNCTION_RESULT_BEGIN") {
			parts := strings.Fields(s)
			s = strings.Join(parts[:len(parts)-1], " ") // remove timestamp
		}
		lines = append(lines, s)
	}
	wantDyncfg, gotDyncfg := strings.TrimSpace(s.wantDyncfg), strings.TrimSpace(strings.Join(lines, "\n"))

	//fmt.Println(gotDyncfg)

	assert.Equal(t, wantDyncfg, gotDyncfg, "dyncfg commands")

	var n int
	for _, cfgs := range mgr.discoveredConfigs.items {
		n += len(cfgs)
	}

	require.Len(t, s.wantDiscovered, n, "discoveredConfigs: different len")

	for _, cfg := range s.wantDiscovered {
		cfgs, ok := mgr.discoveredConfigs.items[cfg.Source()]
		require.Truef(t, ok, "discoveredConfigs: source %s is not found", cfg.Source())
		_, ok = cfgs[cfg.Hash()]
		require.Truef(t, ok, "discoveredConfigs: source %s config %d is not found", cfg.Source(), cfg.Hash())
	}

	require.Len(t, s.wantSeen, len(mgr.seenConfigs.items), "seenConfigs: different len")

	for _, scfg := range s.wantSeen {
		v, ok := mgr.seenConfigs.lookup(scfg.cfg)
		require.Truef(t, ok, "seenConfigs: config '%s' is not found", scfg.cfg.UID())
		require.Truef(t, scfg.status == v.status, "seenConfigs: wrong status, want %s got %s", scfg.status, v.status)
	}

	require.Len(t, s.wantExposed, len(mgr.exposedConfigs.items), "exposedConfigs: different len")

	for _, scfg := range s.wantExposed {
		v, ok := mgr.exposedConfigs.lookup(scfg.cfg)
		require.Truef(t, ok && scfg.cfg.UID() == v.cfg.UID(), "exposedConfigs: config '%s' is not found", scfg.cfg.UID())
		require.Truef(t, scfg.status == v.status, "exposedConfigs: wrong status, want %s got %s", scfg.status, v.status)
	}
}

func prepareMockRegistry() module.Registry {
	reg := module.Registry{}

	reg.Register("success", module.Creator{
		JobConfigSchema: module.MockConfigSchema,
		Create: func() module.Module {
			return &module.MockModule{
				ChartsFunc: func() *module.Charts {
					return &module.Charts{&module.Chart{ID: "id", Title: "title", Units: "units", Dims: module.Dims{{ID: "id1"}}}}
				},
				CollectFunc: func() map[string]int64 { return map[string]int64{"id1": 1} },
			}
		},
	})
	reg.Register("fail", module.Creator{
		Create: func() module.Module {
			return &module.MockModule{
				InitFunc: func() error { return errors.New("mock failed init") },
			}
		},
	})

	return reg
}

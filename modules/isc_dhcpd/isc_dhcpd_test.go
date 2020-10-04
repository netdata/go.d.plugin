package isc_dhcpd

import (
	"testing"

	"github.com/netdata/go.d.plugin/agent/module"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}

func TestDHCPd_Cleanup(t *testing.T) {
	assert.NotPanics(t, New().Cleanup)
}

// TODO: finish
func TestDHCPd_Init(t *testing.T) {
	tests := map[string]struct {
		config   Config
		wantFail bool
	}{
		"default":               {},
		"'leases_path' not set": {},
		"'pools' not set":       {},
		"ok config ('leases_path' and 'pools' are set)": {},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			dhcpd := New()
			dhcpd.Config = test.config

			if test.wantFail {
				assert.False(t, dhcpd.Init())
			} else {
				assert.True(t, dhcpd.Init())
			}
		})
	}
}

// TODO: finish
func TestDHCPd_Check(t *testing.T) {
	tests := map[string]struct {
		prepare  func() *DHCPd
		wantFail bool
	}{
		"leases file doesn't exist": {
			prepare: func() *DHCPd {
				dhcpd := New()
				return dhcpd
			},
		},
		"empty leases file": {
			prepare: func() *DHCPd {
				dhcpd := New()
				return dhcpd
			},
		},
		"leases file with active leases": {
			prepare: func() *DHCPd {
				dhcpd := New()
				return dhcpd
			},
		},
		"leases file without active leases": {
			prepare: func() *DHCPd {
				dhcpd := New()
				return dhcpd
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			dhcpd := test.prepare()
			require.True(t, dhcpd.Init())

			if test.wantFail {
				assert.False(t, dhcpd.Init())
			} else {
				assert.True(t, dhcpd.Init())
			}
		})
	}
}

func TestDHCPd_Charts(t *testing.T) {
	dhcpd := New()
	dhcpd.LeasesPath = "leases_path"
	dhcpd.Pools = []PoolConfig{
		{Name: "name", Networks: "192.0.2.0/24"},
	}
	require.True(t, dhcpd.Init())

	assert.NotNil(t, dhcpd.Charts())
}

// TODO: finish
func TestDHCPd_Collect(t *testing.T) {
	tests := map[string]struct {
		prepare       func() *DHCPd
		wantCollected map[string]int64
	}{
		"dhcp_v4": {
			prepare: func() *DHCPd {
				dhcpd := New()
				return dhcpd
			},
			wantCollected: map[string]int64{},
		},
		"dhcp_v4_backup": {
			prepare: func() *DHCPd {
				dhcpd := New()
				return dhcpd
			},
			wantCollected: map[string]int64{},
		},
		"dhcp_v6": {
			prepare: func() *DHCPd {
				dhcpd := New()
				return dhcpd
			},
			wantCollected: map[string]int64{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			dhcpd := test.prepare()
			require.True(t, dhcpd.Init())

			collected := dhcpd.Collect()

			assert.Equal(t, test.wantCollected, collected)
			ensureCollectedHasAllChartsDimsVarsIDs(t, dhcpd, collected)
		})
	}
}

func ensureCollectedHasAllChartsDimsVarsIDs(t *testing.T, dhcpd *DHCPd, collected map[string]int64) {
	for _, chart := range *dhcpd.Charts() {
		if chart.Obsolete {
			continue
		}
		for _, dim := range chart.Dims {
			_, ok := collected[dim.ID]
			assert.Truef(t, ok, "collected metrics has no data for dim '%s' chart '%s'", dim.ID, chart.ID)
		}
		for _, v := range chart.Vars {
			_, ok := collected[v.ID]
			assert.Truef(t, ok, "collected metrics has no data for var '%s' chart '%s'", v.ID, chart.ID)
		}
	}
}

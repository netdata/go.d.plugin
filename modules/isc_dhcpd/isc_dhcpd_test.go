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

func TestDHCPd_Init(t *testing.T) {
	tests := map[string]struct {
		config   Config
		wantFail bool
	}{
		"default":               {
			config: New().Config,
			wantFail: true,
		},
		"'leases_path' not set": {
			config: Config{
				LeasesPath: "",
				Pools: []PoolConfig{
					{
						Name: "test",
						Networks: "10.220.252.0/24",
					},
				},
			},
			wantFail: true,
		},
		"'pools' not set":       {
			config: Config{
				LeasesPath: "testdata/dhcpd.leases_ipv4",
				Pools: []PoolConfig{
					{
						Name: "test",
						Networks: "",
					},
				},
			},
			wantFail: true,
		},
		"ok config ('leases_path' and 'pools' are set)": {
			config: Config{
				LeasesPath: "testdata/dhcpd.leases_ipv4",
				Pools: []PoolConfig{
					{
						Name: "test",
						Networks: "10.220.252.0/24",
					},
				},
			},
			wantFail: false,
		},
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

func TestDHCPd_Check(t *testing.T) {
	tests := map[string]struct {
		prepare  func() *DHCPd
		wantFail bool
	}{
		"leases file doesn't exist": {
			prepare: prepareDHCPdLeaseWithoutFile,
			wantFail: true,
		},
		"empty leases file": {
			prepare: prepareDHCPdCleanLease,
			wantFail: false,
		},
		"leases file with active leases": {
			prepare: prepareDHCPdIPV4Lease,
			wantFail: false,
		},
		"leases file without active leases": {
			prepare: prepareDHCPdIPV4BkpLease,
			wantFail: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			dhcpd := test.prepare()
			dhcpd.Init()

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

func TestDHCPd_Collect(t *testing.T) {
	tests := map[string]struct {
		prepare       func() *DHCPd
		wantCollected map[string]int64
	}{
		"dhcp_v4": {
			prepare: prepareDHCPdIPV4Lease,
			wantCollected: map[string]int64{
				"active_leases_total" : 2,
				"pool_name_active_leases" : 2,
				"pool_name_utilization" : 78,
			},
		},
		"dhcp_v4_backup": {
			prepare: prepareDHCPdIPV4BkpLease,
			wantCollected: map[string]int64{
				"active_leases_total" : 0,
				"pool_name_active_leases" : 0,
				"pool_name_utilization" : 0,
			},
		},
		"dhcp_v6": {
			prepare: prepareDHCPdIPV6Lease,
			wantCollected: map[string]int64{
				"active_leases_total" : 6,
				"pool_name_active_leases" : 6,
				"pool_name_utilization" : 3529,
			},
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

func prepareDHCPdLeaseWithoutFile() *DHCPd {
	dhdcpd := New()
	dhdcpd.Config = Config{
		LeasesPath: "testdata/no_file.lease",
		Pools: nil,
	}
	
	return dhdcpd
}

func prepareDHCPdCleanLease() *DHCPd {
	dhdcpd := New()
	dhdcpd.Config = Config{
		LeasesPath: "testdata/clean.lease",
		Pools: []PoolConfig {
			{
				Name: "name", 
				Networks: "10.220.252.0/24",
			},
		},
	}
	
	return dhdcpd
}

func prepareDHCPdIPV4Lease() *DHCPd {
	dhdcpd := New()
	dhdcpd.Config = Config{
		LeasesPath: "testdata/dhcpd.leases_ipv4",
		Pools: []PoolConfig {
			{
				Name: "name", 
				Networks: "10.220.252.0/24",
			},
		},
	}
	
	return dhdcpd
}

func prepareDHCPdIPV4BkpLease() *DHCPd {
	dhdcpd := New()
	dhdcpd.Config = Config{
		LeasesPath: "testdata/dhcpd.leases_ipv4_backup",
		Pools: []PoolConfig {
			{
				Name: "name", 
				Networks: "192.168.0.0/24",
			},
		},
	}
	
	return dhdcpd
}

func prepareDHCPdIPV6Lease() *DHCPd {
	dhdcpd := New()
	dhdcpd.Config = Config{
		LeasesPath: "testdata/dhcpd.leases_ipv6",
		Pools: []PoolConfig {
			{
				Name: "name", 
				Networks: "1985:470:1f0b:c9a::000-1985:470:1f0b:c9a::010",
			},
		},
	}
	
	return dhdcpd
}
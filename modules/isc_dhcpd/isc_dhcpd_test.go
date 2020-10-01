package isc_dhcpd

import (
	"testing"
	"net"

	"github.com/netdata/go.d.plugin/pkg/iprange"
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
		config Config
		wantNumOfCharts int
		wantFail bool
	} {
		"default" : {
			config: New().Config,
			wantFail: true,
		},
		"empty Lease file and pools" : {
			config: Config {
				LeaseFile : "",
				LastModification : 0,
				Pools : nil,
				Dim : nil,
				data : nil,
			},
			wantFail: true,
		},
		"With lease file" : {
			config : Config {
				LeaseFile : "testdata/ipv4_dhcpd1.leases",
				LastModification : 0,
				Pools : nil,
				Dim : nil,
				data : nil,
			},
			wantFail: true,
		},
		"only one host" : {
			config : Config {
				LeaseFile : "testdata/ipv4_dhcpd1.leases",
				Pools : map[string]string{
					"office" : "192.168.0.0-192.168.0.254",
				},
				Dim : map[string]Dimensions{
						"office" : Dimensions{ Values : 
							iprange.New(net.ParseIP("192.168.0.0"),
							net.ParseIP("192.168.0.254")),
					},
				},
			},
			wantFail: false,
			wantNumOfCharts: 3,
		},
		"four hosts" : {
			config : Config {
				LeaseFile : "testdata/ipv4_dhcpd4.leases",
				Pools : map[string]string{
					"office" : "10.220.252.0-10.220.252.254",
				},
				Dim : map[string]Dimensions{
						"office" : Dimensions{ Values : 
							iprange.New(net.ParseIP("10.220.252.0"),
							net.ParseIP("10.220.252.254")),
					},
				},
			},
			wantFail: false,
			wantNumOfCharts: 3,
		},
		"ipv6" : {
			config : Config {
				LeaseFile : "testdata/ipv6_dhcpd.leases",
				Pools : map[string]string{
					"office" : "1985:470:1f0b:c9a::000-1985:470:1f0b:c9a::255",
				},
				Dim : map[string]Dimensions{
						"office" : Dimensions{ Values : 
							iprange.New(net.ParseIP("1985:470:1f0b:c9a::000"),
							net.ParseIP("1985:470:1f0b:c9a::255")),
					},
				},
			},
			wantFail: false,
			wantNumOfCharts: 3,
		},
		"backup" : {
			config : Config {
				LeaseFile : "testdata/backup.leases",
				Pools : map[string]string{
					"office" : "192.168.0.0-192.168.0.254",
				},
				Dim : map[string]Dimensions{
						"office" : Dimensions{ Values : 
							iprange.New(net.ParseIP("192.168.0.0"),
							net.ParseIP("192.168.0.254")),
					},
				},
			},
			wantFail: false,
			wantNumOfCharts: 3,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			d := New()
			d.Config = test.config

			if test.wantFail {
				assert.False(t, d.Init())
			} else {
				require.True(t, d.Init())
				assert.Equal(t, test.wantNumOfCharts, len(*d.Charts()))
			}
		})
	}
}

func TestDHCPd_Check(t *testing.T) {
	tests := map[string]struct {
		lease func() *DHCPd
	} {
		"lease file 1" : {lease : ipv4_leaseOne},
		"lease file 4" : {lease : ipv4_leaseFour},
		"lease ipv6 file" : {lease : ipv6_leaseSix},
		"backup" : {lease : ipv4_leaseBkp},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			d := test.lease()

			require.True(t, d.Init())
			assert.True(t, d.Check())
		})
	}
}

func TestDHCPd_Collect(t *testing.T) {
	tests := map[string]struct {
		lease func() *DHCPd
		wantCollected map[string]int64
	} {
		"lease file 1" : {
			lease : ipv4_leaseOne,
			wantCollected : map[string]int64{
				"office_active" : 0,
				"office_total" : 1,
				"office_utilization" : 3,
			},
		},
		"lease file 4" : {
			lease : ipv4_leaseFour,
			wantCollected : map[string]int64{
				"office_active" : 2,
				"office_total" : 4,
				"office_utilization" : 15,
			},
		},
		"ipv6" : {
			lease : ipv6_leaseSix,
			wantCollected : map[string]int64{
				"office_active" : 6,
				"office_total" : 6,
				"office_utilization" : 10,
			},
		},
		"backup" : {
			lease : ipv4_leaseOne,
			wantCollected : map[string]int64{
				"office_active" : 0,
				"office_total" : 1,
				"office_utilization" : 3,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			d := test.lease()

			require.True(t, d.Init())

			collected := d.Collect()

			assert.Equal(t, test.wantCollected, collected)
		})
	}
}

func ipv4_leaseOne() *DHCPd {
	d := New()

	d.Config.LeaseFile = "testdata/ipv4_dhcpd1.leases"
	d.Config.Pools = map[string]string{
					"office" : "192.168.0.0-192.168.0.254",
				}
	d.Config.Dim = map[string]Dimensions{
						"office" : Dimensions{ Values : 
							iprange.New(net.ParseIP("192.168.0.0"),
							net.ParseIP("192.168.0.254")),
					},
				}

	return d
}

func ipv4_leaseFour() *DHCPd {
	d := New()

	d.Config.LeaseFile = "testdata/ipv4_dhcpd4.leases"
	d.Config.Pools = map[string]string{
					"office" : "10.220.252.0-10.220.252.254",
				}
	d.Config.Dim = map[string]Dimensions{
						"office" : Dimensions{ Values : 
							iprange.New(net.ParseIP("10.220.252.0"),
							net.ParseIP("10.220.252.254")),
					},
				}


	return d
}

func ipv6_leaseSix() *DHCPd {
	d := New()

	d.Config.LeaseFile = "testdata/ipv6_dhcpd.leases"
	d.Config.Pools = map[string]string{
					"office" : "1985:470:1f0b:c9a::000-1985:470:1f0b:c9a::255",
				}
	d.Config.Dim = map[string]Dimensions{
						"office" : Dimensions{ Values : 
							iprange.New(net.ParseIP("1985:470:1f0b:c9a::000"),
							net.ParseIP("1985:470:1f0b:c9a::255")),
					},
				}


	return d
}

func ipv4_leaseBkp() *DHCPd {
	d := New()

	d.Config.LeaseFile = "testdata/backup.leases"
	d.Config.Pools = map[string]string{
					"office" : "192.168.0.0-192.168.0.254",
				}
	d.Config.Dim = map[string]Dimensions{
						"office" : Dimensions{ Values : 
							iprange.New(net.ParseIP("192.168.0.0"),
							net.ParseIP("192.168.0.254")),
					},
				}

	return d
}
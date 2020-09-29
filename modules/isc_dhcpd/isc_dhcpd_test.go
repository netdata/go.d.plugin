package isc_dhcpd

/*
import (
	"testing"
	"net"

	"github.com/netdata/go.d.plugin/pkg/ip"
	"github.com/netdata/go-orchestrator/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}

func TestDHCPD_Cleanup(t *testing.T) {
	assert.NotPanics(t, New().Cleanup)
}

func TestDHCPD_Init(t *testing.T) {
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
				Pools : nil,
			},
			wantFail: true,
		},
		"With lease file" : {
			config : Config {
				LeaseFile : "testdata/dhcpd1.leases",
				Pools : nil,
			},
			wantFail: true,
		},
		"only one" : {
			config : Config {
				LeaseFile : "testdata/dhcpd1.leases",
				Pools : map[string]string{
					"office" : "192.168.0.0-192.168.0.254",
				},
				Dim : map[string]Dimensions{
						"office" : Dimensions{ Values : ip.Range {
							Start : net.ParseIP("192.168.0.0"),
							End :  net.ParseIP("192.168.0.254"),
						},
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


func TestDHCPD_Check(t *testing.T) {
	tests := map[string]struct {
		lease func() *DHCPD
	} {
		"lease file 1" : {lease : leaseOne},
		"lease file 4" : {lease : leaseFour},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			d := test.lease()

			require.True(t, d.Init())
			assert.True(t, d.Check())
		})
	}
}

func TestDHCPD_Collect(t *testing.T) {
	tests := map[string]struct {
		lease func() *DHCPD
		wantCollected map[string]int64
	} {
		"lease file 1" : {
			lease : leaseOne,
			wantCollected : map[string]int64{
				"office_active" : 1,
				"office_total" : 1,
				"office_utilization" : 3,
			},
		},
		"lease file 4" : {
			lease : leaseOne,
			wantCollected : map[string]int64{
				"office_active" : 1,
				"office_total" : 1,
				"office_utilization" : 11,
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

func leaseOne() *DHCPD {
	d := New()

	d.Config.LeaseFile = "testdata/dhcpd1.leases"

	return d
}

func leaseFour() *DHCPD {
	d := New()

	d.Config.LeaseFile = "testdata/dhcpd4.leases"

	return d
}
*/
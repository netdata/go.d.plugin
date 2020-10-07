package isc_dhcpd

import (
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
)

func init() {
	module.Register("isc_dhcpd", module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 1,
		},
		Create: func() module.Module { return New() },
	})
}

type (
	Config struct {
		LeasesPath string       `yaml:"leases_path"`
		Pools      []PoolConfig `yaml:"pools"`
	}
	PoolConfig struct {
		Name     string `yaml:"name"`
		Networks string `yaml:"networks"`
	}
)

type DHCPd struct {
	module.Base
	Config `yaml:",inline"`

	charts        *module.Charts
	pools         []ipPool
	leasesModTime time.Time
	collected     map[string]int64
}

func New() *DHCPd {
	return &DHCPd{
		Config: Config{
			LeasesPath: "/var/lib/dhcp/dhcpd.leases",
		},

		collected: make(map[string]int64),
	}
}

func (DHCPd) Cleanup() {}

func (d *DHCPd) Init() bool {
	err := d.validateConfig()
	if err != nil {
		d.Errorf("config validation: %v", err)
		return false
	}

	pools, err := d.initPools()
	if err != nil {
		d.Errorf("ip pools init: %v", err)
		return false
	}
	d.pools = pools

	charts, err := d.initCharts(pools)
	if err != nil {
		d.Errorf("charts init: %v", err)
		return false
	}
	d.charts = charts

	d.Debugf("monitoring leases file: %v", d.Config.LeasesPath)
	d.Debugf("monitoring ip pools: %v", d.Config.Pools)
	return true
}

func (d *DHCPd) Check() bool {
	return len(d.Collect()) > 0
}

func (d *DHCPd) Charts() *module.Charts {
	return d.charts
}

func (d *DHCPd) Collect() map[string]int64 {
	mx, err := d.collect()
	if err != nil {
		d.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}

	return mx
}

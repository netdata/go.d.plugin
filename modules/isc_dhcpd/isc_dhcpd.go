package isc_dhcpd

import (
	"github.com/netdata/go.d.plugin/pkg/ip"

	"github.com/netdata/go.d.plugin/agent/module"
)

type (
	Config struct {
		LeaseFile string            `yaml:"leases_path"`
		Pools     map[string]string `yaml:"pools"`
		Dim       map[string]Dimensions
	}

	Dimensions struct {
		Values ip.IRange
		Name  string
	}
)

type DHCPD struct {
	module.Base
	Config `yaml:",inline"`

	collectedLeases bool
	charts          *module.Charts
}

func init() {
	module.Register("isc_dhcpd", module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 10,
		},
		Create: func() module.Module { return New() },
	})
}

func New() *DHCPD {
	return &DHCPD{
		Config: Config{
			LeaseFile: "",
			Pools: nil,
			Dim: make(map[string]Dimensions),
		},
	}
}

func (DHCPD) Cleanup() {
}

func (d *DHCPD) Init() bool {
	err := d.validateConfig()
	if err != nil {
		d.Errorf("Error on validate config: %v", err)
		return false
	}

	for i, v := range d.Config.Pools {
		r := ip.ParseRange(v)
		if r != nil {
			d.Config.Dim[i] = Dimensions{Values: r, Name: v}
		}
	}

	charts, err := d.initCharts()
	if err != nil {
		d.Errorf("Error on chart initialization: %v", err)
		return false
	}
	d.charts = charts

	d.Debugf("Monitoring lease file %v", d.Config.LeaseFile)
	//	d.Debugf("Monitoring pools %v", d.Config.Pools.pools)
	return true
}

func (d *DHCPD) Check() bool {
	return len(d.Collect()) > 0
}

func (d *DHCPD) Charts() *module.Charts {
	return d.charts
}

func (d *DHCPD) Collect() map[string]int64 {
	mx, err := d.collect()
	if err != nil {
		d.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}

	return mx
}

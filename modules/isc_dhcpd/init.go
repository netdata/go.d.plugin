package isc_dhcpd

import (
	"errors"

	"github.com/netdata/go.d.plugin/agent/module"
)

func (d *DHCPD) validateConfig() error {
	if d.Config.LeaseFile == "" || len(d.Config.Pools) == 0 {
		return errors.New("neither pools nor 'lease file' is defined")
	}

	return nil
}

func (d *DHCPD) initCharts() (*module.Charts, error) {
	charts := &module.Charts{}

	if len(d.Config.LeaseFile) > 0 {
		if err := charts.Add(*dhcpdCharts.Copy()...); err != nil {
			return nil, err
		}
	}

	if len(*charts) == 0 {
		return nil, errors.New("empty charts")
	}

	return charts, nil
}

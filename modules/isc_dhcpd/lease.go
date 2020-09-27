package isc_dhcpd

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/netdata/go-orchestrator/module"
	leases "github.com/npotts/go-dhcpd-leases"
)

type LeaseFile struct {
	IP    net.IP
	Ends  time.Time
	State string
}

func parseDHCPLease(filename string) ([]LeaseFile, error) {
	set := make(map[string]int)
	var list []LeaseFile

	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("Cannot open file")
	}
	defer f.Close()

	l := leases.Parse(f)
	if len(l) > 0 {
		for _, v := range l {
			index := v.IP.String() + v.BindingState
			if idx, ok := set[index]; ok {
				list[idx] = LeaseFile{IP: v.IP, Ends: v.Ends, State: v.BindingState}
			} else {
				set[index] = len(list)
				list = append(list, LeaseFile{IP: v.IP, Ends: v.Ends, State: v.BindingState})
			}
		}
	}

	return list, nil
}

func (d *DHCPD) parseLease(c map[string]int64) {
	if !d.collectedLeases {
		d.collectedLeases = true
		d.addPoolsToCharts()
	}

	// In order to process both DHCPv4 and DHCPv6 messages you will need to run two separate instances of the dhcpd process.
	// Each of these instances will need it's own lease file.
	l, err := parseDHCPLease(d.LeaseFile)
	if err != nil {
		return
	}

	if len(l) > 0 {
		for _, v := range l {
			fmt.Println(v.IP.String())
			c["Total"] = 1
		}
	}
}

func (d *DHCPD) addPoolsToCharts() {
	for _, v := range d.Config.Pools {
		d.addPoolToCharts(v)
	}
}

func (d *DHCPD) addPoolToCharts(str string) {
	for _, val := range dhcpdCharts {
		chart := d.Charts().Get(val.ID)
		if chart == nil {
			d.Warningf("add dimension: couldn't find '%s' chart", val.ID)
			continue
		}

		var id string
		switch chart.ID {
		case dhcpPollsUtilization.ID:
			id = str + " utilization"
		case dhcpPollsActiveLeases.ID:
			id = str + " active"
		case dhcpPollsTotalLeases.ID:
			id = str + " total"
		}

		dim := &module.Dim{ID: id, Name: str}

		if err := chart.AddDim(dim); err != nil {
			d.Warning(err)
			continue
		}

		chart.MarkNotCreated()
	}
}

package isc_dhcpd

import (
	"errors"
	"net"
	"os"
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
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

	l, err := parseDHCPLease(d.LeaseFile)
	if err != nil {
		return
	}

	// The test has a problem when we compare the timeto set it as active.
	currTime := time.Now()
	if len(l) > 0 {
		for _, v := range l {
			prefix := d.getDimensionPrefix(v.IP)
			if prefix != "" {
				if _, ok := c[prefix + "_active"] ; ok {
					c[prefix + "_active"] = markActive(c[prefix + "_active"], currTime, v.Ends, v.State)
					c[prefix + "_total"] = incrementValues(c[prefix + "_total"])
				} else {
					c[prefix + "_active"] = markActive(0, currTime, v.Ends, v.State)
					c[prefix + "_total"] = 1
				}
			}
		}
	}

	for idx, v := range d.Config.Dim {
		i := *v.Values.Hosts()
		f := (float64(c[idx + "_total"])/float64(i.Uint64()))*1000
		c[idx + "_utilization"] = int64(f)
	}
}

func incrementValues(prev int64)  int64{
	prev++
	return prev
}

func markActive(prev int64, curr time.Time, old time.Time, state string) int64 {
	if state == "active" {
		test := curr.Unix() - old.Unix()
		if test >= 0 {
			prev++
		}
	}
	return prev
}

func (d *DHCPD) getDimensionPrefix(ip net.IP) string {
	for idx, v := range d.Config.Dim {
		if (v.Values.Contains(ip)) {
			return idx
		}
	}
	return ""
}

func (d *DHCPD) addPoolsToCharts() {
	for idx, _ := range d.Config.Dim {
		d.addPoolToCharts(idx)
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
			id = str + "_utilization"
		case dhcpPollsActiveLeases.ID:
			id = str + "_active"
		case dhcpPollsTotalLeases.ID:
			id = str + "_total"
		}

		dim := &module.Dim{ID: id, Name: str}

		if err := chart.AddDim(dim); err != nil {
			d.Warning(err)
			continue
		}

		chart.MarkNotCreated()
	}
}

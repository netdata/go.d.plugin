package isc_dhcpd

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
)

func (d *DHCPd) collect() (map[string]int64, error) {
	cm := make(map[string]int64)

	d.parseLease()
	d.fillDimension(cm)

	return cm, nil
}

func (d *DHCPd) fillDimension(c map[string]int64) {
	l := d.leases

	currTime := time.Now()
	if len(l) > 0 {
		for _, v := range l {
			prefix := d.getDimensionPrefix(net.ParseIP(v.ip))
			if prefix != "" {
				if _, ok := c[prefix + "_active"] ; ok {
					c[prefix + "_active"] = markActive(c[prefix + "_active"], currTime, v.ends, v.bindingState)
					c[prefix + "_total"] = incrementValues(c[prefix + "_total"])
				} else {
					c[prefix + "_active"] = markActive(0, currTime, v.ends, v.bindingState)
					c[prefix + "_total"] = 1
				}
			}
		}
	}

	for idx, v := range d.Dim {
		i := *v.Values.Size()
		f := (float64(c[idx + "_total"])/float64(i.Uint64()))*1000
		c[idx + "_utilization"] = int64(f)
	}
}

func incrementValues(prev int64)  int64{
	prev++
	return prev
}

func markActive(prev int64, curr time.Time, oldTime string, state string) int64 {
	if state == "active" {
		b := []byte(oldTime)
		var old time.Time
		if bytes.HasPrefix(b, []byte("epoch")) {
			val, _ := strconv.ParseInt(oldTime[6:], 10, 64)
			old = time.Unix(val, 0 )
		} else {
			old, _ = time.Parse("2006/01/02 15:04:05", oldTime[2:])
		}
		fmt.Println(old)

		test := curr.Unix() - old.Unix()
		if test >= 0 {
			prev++
		}
	}
	return prev
}

func (d *DHCPd) getDimensionPrefix(ip net.IP) string {
	for idx, v := range d.Dim {
		if (v.Values.Contains(ip)) {
			return idx
		}
	}
	return ""
}


func (d *DHCPd) addPoolsToCharts() {
	for idx := range d.Dim {
		d.addPoolToCharts(idx)
	}
}

func (d *DHCPd) addPoolToCharts(str string) {
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
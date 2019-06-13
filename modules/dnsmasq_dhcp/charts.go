package dnsmasq_dhcp

import (
	"fmt"

	"github.com/netdata/go.d.plugin/modules/dnsmasq_dhcp/ip"

	"github.com/netdata/go-orchestrator/module"
)

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dim is an alias for module.Dim
	Dim = module.Dim
)

var charts = Charts{
	{
		ID:    "ipv4_active_leases",
		Title: "DHCP Range Active Leases",
		Units: "active leases",
		Fam:   "ipv4",
		Ctx:   "dhcp_range_active_leases",
	},
	{
		ID:    "ipv4_utilization",
		Title: "DHCP Range Utilization",
		Units: "percentage",
		Fam:   "ipv4",
		Ctx:   "dhcp_range_utilization",
	},
	{
		ID:    "ipv6_active_leases",
		Title: "DHCP Range Active Leases",
		Units: "active leases",
		Fam:   "ipv6",
		Ctx:   "dhcp_range_active_leases",
	},
	{
		ID:    "ipv6_utilization",
		Title: "DHCP Range Utilization",
		Units: "percentage",
		Fam:   "ipv6",
		Ctx:   "dhcp_range_utilization",
	},
}

func (d DnsmasqDHCP) charts() *Charts {
	cs := charts.Copy()

	for _, r := range d.ranges {
		addRangeToCharts(cs, r)
	}

	rv := &Charts{}

	for _, c := range *cs {
		if len(c.Dims) == 0 {
			continue
		}
		panicIf(rv.Add(c))
	}

	return rv
}

func addRangeToCharts(cs *Charts, r ip.IRange) {
	var prefix string

	switch r.Family() {
	default:
		panic(fmt.Sprintf("invalid ip range '%s'", r))
	case ip.V4Family:
		prefix = "ipv4"
	case ip.V6Family:
		prefix = "ipv6"
	}

	name := r.String()
	panicIf(cs.Get(prefix + "_active_leases").AddDim(&Dim{ID: name}))
	panicIf(cs.Get(prefix + "_utilization").AddDim(&Dim{ID: name + "_percentage", Name: name}))
}

func panicIf(err error) {
	if err == nil {
		return
	}
	panic(err)
}

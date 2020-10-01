package dnsmasq_dhcp

import (
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/iprange"
)

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dim is an alias for module.Dim
	Dim = module.Dim
)

var charts = Charts{
	{
		ID:    "%s_utilization",
		Title: "DHCP Range Utilization",
		Units: "percentage",
		Ctx:   "dnsmasq_dhcp.dhcp_range_utilization",
	},
	{
		ID:    "%s_allocated_leases",
		Title: "DHCP Range Allocated Leases",
		Units: "leases",
		Ctx:   "dnsmasq_dhcp.dhcp_range_allocated_leases",
	},
}

func (d DnsmasqDHCP) charts() *Charts {
	cs := &Charts{}

	for _, r := range d.ranges {
		panicIf(cs.Add(*addRangeCharts(r)...))
	}

	return cs
}

func addRangeCharts(r iprange.Range) *Charts {
	cs := charts.Copy()

	name := r.String()

	c := cs.Get("%s_utilization")
	c.ID = fmt.Sprintf(c.ID, name)
	c.Fam = name
	panicIf(c.AddDim(&Dim{ID: name + "_percentage", Name: "used"}))

	c = cs.Get("%s_allocated_leases")
	c.ID = fmt.Sprintf(c.ID, name)
	c.Fam = name
	panicIf(c.AddDim(&Dim{ID: name, Name: "allocated"}))

	return cs
}

func panicIf(err error) {
	if err == nil {
		return
	}
	panic(err)
}

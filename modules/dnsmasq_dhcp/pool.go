package dnsmasq_dhcp

import "github.com/netdata/go.d.plugin/modules/dnsmasq_dhcp/ip"

type pool struct {
	name string
	ips  int64
	ip.IRange
}

func (p pool) utilization() float64 {
	total := p.Hosts()
	if !total.IsInt64() {
		return 0
	}
	return float64(p.ips) * 100 / float64(total.Int64())
}

func (p pool) leasedIPs() int64 { return p.ips }

func (p pool) hasValidSize() bool { return p.Hosts().IsInt64() }

func (p *pool) resetIPs() { p.ips = 0 }

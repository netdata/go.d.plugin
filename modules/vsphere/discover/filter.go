package discover

import (
	rs "github.com/netdata/go.d.plugin/modules/vsphere/resources"
)

func (d VSphereDiscoverer) matchHost(host *rs.Host) bool {
	if d.HostMatcher == nil {
		return true
	}
	return d.HostMatcher.Match(host)
}

func (d VSphereDiscoverer) matchVM(vm *rs.VM) bool {
	if d.VMMatcher == nil {
		return true
	}
	return d.VMMatcher.Match(vm)
}

func (d VSphereDiscoverer) removeUnmatched(res *rs.Resources) (removed int) {
	d.Debug("discovering : starting filtering resources")
	removed += d.removeUnmatchedHosts(res.Hosts)
	removed += d.removeUnmatchedVMs(res.VMs)
	return
}

func (d VSphereDiscoverer) removeUnmatchedHosts(hosts rs.Hosts) (removed int) {
	for _, v := range hosts {
		if !d.matchHost(v) {
			removed++
			hosts.Remove(v.ID)
		}
	}
	d.Debugf("discovering : removed %d unmatched hosts", removed)
	return removed
}

func (d VSphereDiscoverer) removeUnmatchedVMs(vms rs.VMs) (removed int) {
	for _, v := range vms {
		if !d.matchVM(v) {
			removed++
			vms.Remove(v.ID)
		}
	}
	d.Debugf("discovering : removed %d unmatched vms", removed)
	return removed
}

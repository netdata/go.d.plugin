package discover

import rs "github.com/netdata/go.d.plugin/modules/vsphere/resources"

func (d vSphereDiscoverer) matchHost(host rs.Host) bool {
	if d.HostMatcher == nil {
		return true
	}
	return d.HostMatcher.Match(host)
}

func (d vSphereDiscoverer) matchVM(vm rs.VM) bool {
	if d.VMMatcher == nil {
		return true
	}
	return d.VMMatcher.Match(vm)
}

func (d vSphereDiscoverer) removeUnmatched(res *rs.Resources) (removed int) {
	d.Debug("discovering : starting filtering resources")
	removed += d.removeUnmatchedHosts(res.Hosts)
	removed += d.removeUnmatchedVMs(res.VMs)
	return removed
}

func (d vSphereDiscoverer) removeUnmatchedHosts(hosts rs.Hosts) (removed int) {
	for _, v := range hosts {
		if !d.matchHost(*v) {
			removed++
			hosts.Remove(v.ID)
		}
	}
	if removed > 0 {
		d.Infof("discovering : found %d unmatched hosts, removing them", removed)
	}
	return removed
}

func (d vSphereDiscoverer) removeUnmatchedVMs(vms rs.VMs) (removed int) {
	for _, v := range vms {
		if !d.matchVM(*v) {
			removed++
			vms.Remove(v.ID)
		}
	}
	if removed > 0 {
		d.Infof("discovering : found %d unmatched vms, removing them", removed)
	}
	return removed
}

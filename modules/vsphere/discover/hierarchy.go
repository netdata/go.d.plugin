package discover

import (
	rs "github.com/netdata/go.d.plugin/modules/vsphere/resources"

	"fmt"
)

func (d VSphereDiscoverer) setHierarchy(res *rs.Resources) error {
	d.Debug("discovering : starting hierarchy set process")

	c := d.setClustersHierarchy(res)
	h := d.setHostsHierarchy(res)
	v := d.setVMsHierarchy(res)

	d.Infof("discovering : hierarchy properly set for %d/%d clusters, %d/%d hosts, %d/%d vms",
		c, len(res.Clusters),
		h, len(res.Hosts),
		v, len(res.VMs),
	)

	notSet := len(res.Clusters) + len(res.Hosts) + len(res.VMs) - (c + h + v)
	if notSet > 0 {
		return fmt.Errorf("discovering : %d resources have no hierarchy set", notSet)
	}

	return nil
}

func (d VSphereDiscoverer) setClustersHierarchy(res *rs.Resources) (set int) {
	for _, cluster := range res.Clusters {
		if setClusterHierarchy(cluster, res) {
			set++
		}
	}
	return set
}

func (d VSphereDiscoverer) setHostsHierarchy(res *rs.Resources) (set int) {
	for _, host := range res.Hosts {
		if setHostHierarchy(host, res) {
			set++
		}
	}
	return set
}

func (d VSphereDiscoverer) setVMsHierarchy(res *rs.Resources) (set int) {
	for _, vm := range res.VMs {
		if setVMHierarchy(vm, res) {
			set++
		}
	}
	return set
}

func setClusterHierarchy(cluster *rs.Cluster, res *rs.Resources) bool {
	dc := res.Dcs.Get(cluster.ParentID)
	if dc == nil {
		return false
	}
	cluster.Hier.Dc.Set(dc.ID, dc.Name)
	return cluster.Hier.IsSet()
}

func setHostHierarchy(host *rs.Host, res *rs.Resources) bool {
	cr := res.Clusters.Get(host.ParentID)
	if cr == nil {
		return false
	}
	host.Hier.Cluster.Set(cr.ID, cr.Name)

	dc := res.Dcs.Get(cr.ParentID)
	if dc == nil {
		return false
	}
	host.Hier.Dc.Set(dc.ID, dc.Name)
	return host.Hier.IsSet()
}

func setVMHierarchy(vm *rs.VM, res *rs.Resources) bool {
	h := res.Hosts.Get(vm.ParentID)
	if h == nil {
		return false
	}
	vm.Hier.Host.Set(h.ID, h.Name)

	cr := res.Clusters.Get(h.ParentID)
	if cr == nil {
		return false
	}
	vm.Hier.Cluster.Set(cr.ID, cr.Name)

	dc := res.Dcs.Get(cr.ParentID)
	if dc == nil {
		return false
	}
	vm.Hier.Dc.Set(dc.ID, dc.Name)
	return vm.Hier.IsSet()
}

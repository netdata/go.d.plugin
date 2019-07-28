package discover

import (
	rs "github.com/netdata/go.d.plugin/modules/vsphere/resources"

	"github.com/vmware/govmomi/vim25/mo"
)

func (d vSphereDiscoverer) build(raw *resources) *rs.Resources {
	d.Debug("discovering : starting building resources")

	var res rs.Resources
	res.Dcs = d.buildDatacenters(raw.dcs)
	res.Folders = d.buildFolders(raw.folders)
	res.Clusters = d.buildClusters(raw.clusters)
	fixClustersParentID(&res)
	res.Hosts = d.buildHosts(raw.hosts)
	res.VMs = d.buildVMs(raw.vms)

	d.Debugf("discovering : built %d datacenters, %d folders, %d clusters, %d hosts, %d vms",
		len(res.Dcs),
		len(res.Folders),
		len(res.Clusters),
		len(res.Hosts),
		len(res.VMs),
	)
	return &res
}

// cluster parent is folder by default
// should be called after buildDatacenters, buildFolders and buildClusters
func fixClustersParentID(res *rs.Resources) {
	for _, c := range res.Clusters {
		c.ParentID = findClusterDcID(c.ParentID, res.Folders)
	}
}

func findClusterDcID(parentID string, folders rs.Folders) string {
	f := folders.Get(parentID)
	if f == nil {
		return parentID
	}
	return findClusterDcID(f.ParentID, folders)
}

func (vSphereDiscoverer) buildDatacenters(raw []mo.Datacenter) rs.Dcs {
	dcs := make(rs.Dcs)
	for _, d := range raw {
		dcs.Put(newDC(d))
	}
	return dcs
}

func newDC(raw mo.Datacenter) *rs.Datacenter {
	// Datacenter1 datacenter-2 group-h4 group-v3
	return &rs.Datacenter{
		Name: raw.Name,
		ID:   raw.Reference().Value,
	}
}

func (vSphereDiscoverer) buildFolders(raw []mo.Folder) rs.Folders {
	fs := make(rs.Folders)
	for _, d := range raw {
		fs.Put(newFolder(d))
	}
	return fs
}

func newFolder(raw mo.Folder) *rs.Folder {
	// vm group-v55 datacenter-54
	// host group-h56 datacenter-54
	// datastore group-s57 datacenter-54
	// network group-n58 datacenter-54
	return &rs.Folder{
		Name:     raw.Name,
		ID:       raw.Reference().Value,
		ParentID: raw.Parent.Value,
	}
}

func (vSphereDiscoverer) buildClusters(raw []mo.ComputeResource) rs.Clusters {
	clusters := make(rs.Clusters)
	for _, c := range raw {
		clusters.Put(newCluster(c))
	}
	return clusters
}

func newCluster(raw mo.ComputeResource) *rs.Cluster {
	// s - dummy cluster, c - created by user cluster
	// 192.168.0.201 domain-s61 group-h4
	// New Cluster1 domain-c52 group-h67
	return &rs.Cluster{
		Name:     raw.Name,
		ID:       raw.Reference().Value,
		ParentID: raw.Parent.Value,
	}
}

const (
	poweredOn = "poweredOn"
)

func (d vSphereDiscoverer) buildHosts(raw []mo.HostSystem) rs.Hosts {
	var num int
	hosts := make(rs.Hosts)
	for _, h := range raw {
		//	poweredOn | poweredOff | standBy | unknown
		if h.Runtime.PowerState != poweredOn {
			num++
			continue
		}
		// connected | notResponding | disconnected
		//if v.Runtime.ConnectionState == "" {
		//
		//}
		hosts.Put(newHost(h))
	}
	if num > 0 {
		d.Infof("discovering : found %d not powered on hosts, removing them", num)
	}
	return hosts
}

func newHost(raw mo.HostSystem) *rs.Host {
	// 192.168.0.201 host-22 domain-s61
	// 192.168.0.202 host-28 domain-c52
	// 192.168.0.203 host-33 domain-c52
	return &rs.Host{
		Name:     raw.Name,
		ID:       raw.Reference().Value,
		ParentID: raw.Parent.Value,
		Ref:      raw.Reference(),
	}
}

func (d vSphereDiscoverer) buildVMs(raw []mo.VirtualMachine) rs.VMs {
	var num int
	vms := make(rs.VMs)
	for _, v := range raw {
		//  poweredOff | poweredOn | suspended
		if v.Runtime.PowerState != poweredOn {
			num++
			continue
		}
		// connected | disconnected | orphaned | inaccessible | invalid
		//if v.Runtime.ConnectionState == "" {
		//
		//}
		vms.Put(newVM(v))
	}
	if num > 0 {
		d.Infof("discovering : found %d not powered on vms, removing them", num)
	}
	return vms
}

func newVM(raw mo.VirtualMachine) *rs.VM {
	// deb91 vm-25 group-v3 host-22
	return &rs.VM{
		Name:     raw.Name,
		ID:       raw.Reference().Value,
		ParentID: raw.Runtime.Host.Value,
		Ref:      raw.Reference(),
	}
}

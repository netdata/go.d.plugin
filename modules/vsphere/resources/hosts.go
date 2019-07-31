package resources

import (
	"github.com/vmware/govmomi/performance"
	"github.com/vmware/govmomi/vim25/types"
)

type (
	HostHierarchy struct {
		Dc      HierarchyValue
		Cluster HierarchyValue
	}

	// A host is the virtual representation of the computing and memory resources of a physical machine running ESXi.
	Host struct {
		Name          string
		ID            string
		ParentID      string
		Hier          HostHierarchy
		OverallStatus string
		MetricList    performance.MetricList
		Ref           types.ManagedObjectReference
	}

	Hosts map[string]*Host
)

func (h HostHierarchy) IsSet() bool {
	return h.Dc.IsSet() && h.Cluster.IsSet()
}

func (hs Hosts) Put(host *Host) {
	hs[host.ID] = host
}

func (hs Hosts) Remove(id string) {
	delete(hs, id)
}

func (hs Hosts) Get(id string) *Host {
	return hs[id]
}

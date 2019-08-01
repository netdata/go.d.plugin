package resources

import (
	"github.com/vmware/govmomi/performance"
	"github.com/vmware/govmomi/vim25/types"
)

type HierarchyValue struct {
	ID, Name string
}

func (v HierarchyValue) IsSet() bool {
	return v.ID != "" && v.Name != ""
}

func (v *HierarchyValue) Set(id, name string) {
	v.ID = id
	v.Name = name
}

type (
	VMHierarchy struct {
		Dc      HierarchyValue
		Cluster HierarchyValue
		Host    HierarchyValue
	}

	VM struct {
		Name          string
		ID            string
		ParentID      string
		Hier          VMHierarchy
		OverallStatus string
		MetricList    performance.MetricList
		Ref           types.ManagedObjectReference
	}

	VMs map[string]*VM
)

func (h VMHierarchy) IsSet() bool {
	return h.Dc.IsSet() && h.Cluster.IsSet() && h.Host.IsSet()
}

func (vs VMs) Put(vm *VM) {
	vs[vm.ID] = vm
}

func (vs VMs) Remove(id string) {
	delete(vs, id)
}

func (vs VMs) Get(id string) *VM {
	return vs[id]
}

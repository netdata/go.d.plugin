package vsphere

import (
	"errors"
	"fmt"

	rs "github.com/netdata/go.d.plugin/modules/vsphere/resources"

	"github.com/vmware/govmomi/performance"
)

func (vs *VSphere) collectVMs(mx map[string]int64) error {
	// NOTE: returns unsorted if at least one types.PerfMetricId Instance is not ""
	ems := vs.CollectVMsMetrics(vs.resources.VMs)
	if len(ems) == 0 {
		return errors.New("failed to gather vms metrics")
	}

	vs.collectVMsMetrics(mx, ems)
	return nil
}

func (vs *VSphere) collectVMsMetrics(mx map[string]int64, ems []performance.EntityMetric) {
	vs.nilVMsMetrics()

	for _, em := range ems {
		vm := vs.resources.VMs.Get(em.Entity.Value)
		if vm == nil {
			continue
		}

		vm.Metrics = em.Value
		writeVMMetricsTo(mx, vm)
	}
}

func (vs *VSphere) nilVMsMetrics() {
	for _, v := range vs.resources.VMs {
		v.Metrics = nil
	}
}

func writeVMMetricsTo(to map[string]int64, vm *rs.VM) {
	for _, m := range vm.Metrics {
		if len(m.Value) == 0 || m.Value[0] == -1 {
			continue
		}
		key := buildVMKey(vm, m.Instance, m.Name)
		to[key] = m.Value[0]
	}
}

func buildVMKey(vm *rs.VM, instance string, metricName string) string {
	// NOTE: name is not unique
	if instance == "" {
		return fmt.Sprintf("%s_%s_%s", vm.ID, vm.Name, metricName)
	}
	return fmt.Sprintf("%s_%s_%s_%s", vm.ID, vm.Name, metricName, instance)
}

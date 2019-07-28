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
		return errors.New("failed to collect vms metrics")
	}

	vs.processVMsMetrics(mx, ems)
	return nil
}

func (vs *VSphere) processVMsMetrics(mx map[string]int64, ems []performance.EntityMetric) {
	updated := make(map[string]bool)
	for _, em := range ems {
		vm := vs.resources.VMs.Get(em.Entity.Value)
		if vm == nil {
			continue
		}
		writeVMMetrics(mx, vm, em.Value)
		updated[vm.ID] = true
		vs.collectedVMs[vm.ID] = 0
	}

	for k := range vs.collectedVMs {
		if updated[k] {
			continue
		}
		vs.collectedVMs[k] += 1
	}
}

func writeVMMetrics(dst map[string]int64, vm *rs.VM, metrics []performance.MetricSeries) {
	for _, m := range metrics {
		if len(m.Value) == 0 || m.Value[0] == -1 {
			continue
		}
		key := vmMetricKey(vm, m.Instance, m.Name)
		dst[key] = m.Value[0]
	}
}

func vmMetricKey(vm *rs.VM, instance, metricName string) string {
	if instance == "" {
		return fmt.Sprintf("%s_%s_%s", vm.ID, vm.Name, metricName)
	}
	return fmt.Sprintf("%s_%s_%s_%s", vm.ID, vm.Name, metricName, instance)
}

package vsphere

import (
	"errors"
	"fmt"

	rs "github.com/netdata/go.d.plugin/modules/vsphere/resources"

	"github.com/vmware/govmomi/performance"
)

func (vs *VSphere) collectVMs(mx map[string]int64) error {
	// NOTE: returns unsorted if at least one types.PerfMetricId Instance is not ""
	ems := vs.ScrapeVMsMetrics(vs.resources.VMs)
	if len(ems) == 0 {
		return errors.New("failed to scrape vms metrics")
	}

	collected := vs.collectVMsMetrics(mx, ems)
	vs.updateVMsCharts(collected)
	return nil
}

func (vs *VSphere) collectVMsMetrics(mx map[string]int64, ems []performance.EntityMetric) map[string]bool {
	collected := make(map[string]bool)
	for _, em := range ems {
		vm := vs.resources.VMs.Get(em.Entity.Value)
		if vm == nil {
			continue
		}
		writeVMMetrics(mx, vm, em.Value)
		collected[vm.ID] = true
		vs.discoveredVMs[vm.ID] = 0
	}

	for k := range vs.discoveredVMs {
		if collected[k] {
			continue
		}
		vs.discoveredVMs[k] += 1
	}
	return collected
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

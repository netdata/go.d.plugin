// SPDX-License-Identifier: GPL-3.0-or-later

package vsphere

import (
	"errors"
	"fmt"
	"time"

	rs "github.com/netdata/go.d.plugin/modules/vsphere/resources"

	"github.com/vmware/govmomi/performance"
)

// ManagedEntityStatus
var overallStatuses = []string{"green", "red", "yellow", "gray"}

func (vs *VSphere) collect() (map[string]int64, error) {
	vs.collectionLock.Lock()
	defer vs.collectionLock.Unlock()
	defer vs.removeStale()

	vs.Debug("starting collection process")
	t := time.Now()
	mx := make(map[string]int64)

	err := vs.collectHosts(mx)
	if err != nil {
		return mx, err
	}

	err = vs.collectVMs(mx)
	if err != nil {
		return mx, err
	}

	vs.Debugf("metrics collected, process took %s", time.Since(t))
	return mx, nil
}

func (vs *VSphere) collectHosts(mx map[string]int64) error {
	if len(vs.resources.Hosts) == 0 {
		return nil
	}
	// NOTE: returns unsorted if at least one types.PerfMetricId Instance is not ""
	metrics := vs.scraper.ScrapeHosts(vs.resources.Hosts)
	if len(metrics) == 0 {
		return errors.New("failed to scrape hosts metrics")
	}

	hosts := vs.collectHostsMetrics(mx, metrics)

	vs.updateDiscoveredHosts(hosts)
	vs.updateHostsCharts(hosts)

	return nil
}

func (vs *VSphere) collectHostsMetrics(mx map[string]int64, metrics []performance.EntityMetric) map[string]bool {
	hosts := make(map[string]bool)
	for _, metric := range metrics {
		host := vs.resources.Hosts.Get(metric.Entity.Value)
		if host == nil {
			continue
		}
		writeHostMetrics(mx, host, metric.Value)
		hosts[host.ID] = true
	}
	return hosts
}

func writeHostMetrics(mx map[string]int64, host *rs.Host, metrics []performance.MetricSeries) {
	for _, metric := range metrics {
		if len(metric.Value) == 0 || metric.Value[0] == -1 {
			continue
		}
		key := fmt.Sprintf("host_%s_%s", host.ID, metric.Name)
		mx[key] = metric.Value[0]
	}
	for _, v := range overallStatuses {
		key := fmt.Sprintf("host_%s_overall.status.%s", host.ID, v)
		mx[key] = boolToInt(host.OverallStatus == v)
	}
}

func (vs *VSphere) updateDiscoveredHosts(discoveredHosts map[string]bool) {
	for _, h := range vs.resources.Hosts {
		if _, ok := discoveredHosts[h.ID]; !ok {
			vs.discoveredHosts[h.ID] += 1
		} else {
			vs.discoveredHosts[h.ID] = 0
		}
	}
}

func (vs *VSphere) collectVMs(mx map[string]int64) error {
	if len(vs.resources.VMs) == 0 {
		return nil
	}
	// NOTE: returns unsorted if at least one types.PerfMetricId Instance is not ""
	ems := vs.scraper.ScrapeVMs(vs.resources.VMs)
	if len(ems) == 0 {
		return errors.New("failed to scrape vms metrics")
	}

	vms := vs.collectVMsMetrics(mx, ems)
	vs.updateDiscoveredVMs(vms)
	vs.updateVMsCharts(vms)

	return nil
}

func (vs *VSphere) collectVMsMetrics(mx map[string]int64, metrics []performance.EntityMetric) map[string]bool {
	vms := make(map[string]bool)
	for _, metric := range metrics {
		vm := vs.resources.VMs.Get(metric.Entity.Value)
		if vm == nil {
			continue
		}

		writeVMMetrics(mx, vm, metric.Value)
		vms[vm.ID] = true
	}
	return vms
}

func writeVMMetrics(mx map[string]int64, vm *rs.VM, metrics []performance.MetricSeries) {
	for _, metric := range metrics {
		if len(metric.Value) == 0 || metric.Value[0] == -1 {
			continue
		}
		key := fmt.Sprintf("vm_%s_%s", vm.ID, metric.Name)
		mx[key] = metric.Value[0]
	}
	for _, v := range overallStatuses {
		key := fmt.Sprintf("vm_%s_overall.status.%s", vm.ID, v)
		mx[key] = boolToInt(vm.OverallStatus == v)
	}
}

func (vs *VSphere) updateDiscoveredVMs(discoveredVMs map[string]bool) {
	for _, vm := range vs.resources.VMs {
		if _, ok := discoveredVMs[vm.ID]; !ok {
			vs.discoveredVMs[vm.ID] += 1
		} else {
			vs.discoveredVMs[vm.ID] = 0
		}
	}
}

const failedUpdatesLimit = 10

func (vs *VSphere) removeStale() {
	for hostID, fails := range vs.discoveredHosts {
		if fails < failedUpdatesLimit {
			continue
		}
		vs.removeFromCharts("host_" + hostID)
		delete(vs.charted, hostID)
		delete(vs.discoveredHosts, hostID)
	}
	for vmID, fails := range vs.discoveredVMs {
		if fails < failedUpdatesLimit {
			continue
		}
		vs.removeFromCharts("vm_" + vmID)
		delete(vs.charted, vmID)
		delete(vs.discoveredVMs, vmID)
	}
}

func boolToInt(v bool) int64 {
	if v {
		return 1
	}
	return 0
}

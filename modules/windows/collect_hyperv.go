// SPDX-License-Identifier: GPL-3.0-or-later

package windows

import (
	"strings"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	metricHypervHealthCritical = "windows_hyperv_health_critical"
	metricHypervHealthOK = "windows_hyperv_health_ok"

	metricHypervHypervisorLogicalProcessors = "windows_hyperv_hypervisor_logical_processors"
	metricHypervHypervisorVirtualProcessors = "windows_hyperv_hypervisor_virtual_processors"

	metricHypervRootPartition1GDevicePages = "windows_hyperv_root_partition_1G_device_pages"
	metricHypervRootPartition1GGPAPages = "windows_hyperv_root_partition_1G_gpa_pages"
	metricHypervRootPartition2MDevicePages = "windows_hyperv_root_partition_2M_device_pages"
	metricHypervRootPartition2MGPAPages = "windows_hyperv_root_partition_2M_gpa_pages"
	metricHypervRootPartition4KDevicePages = "windows_hyperv_root_partition_4K_device_pages"
	metricHypervRootPartition4KGPAPages = "windows_hyperv_root_partition_4K_gpa_pages"
	metricHypervRootPartitionAddressSpace = "windows_hyperv_root_partition_address_spaces"
	metricHypervRootPartitionAttachedDevices = "windows_hyperv_root_partition_attached_devices"
	metricHypervRootPartitionDepositedPages = "windows_hyperv_root_partition_deposited_pages"
	metricHypervRootPartitionDeviceDMAErrors = "windows_hyperv_root_partition_device_dma_errors"
	metricHypervRootPartitionDeviceInterruptErrors = "windows_hyperv_root_partition_device_interrupt_errors"
	metricHypervRootPartitionDeviceInterruptThrottleEvents = "windows_hyperv_root_partition_device_interrupt_throttle_events"
	metricHypervRootPartitionGPASpaceModifications = "windows_hyperv_root_partition_gpa_space_modifications"
	metricHypervRootPartitionIOTLBFlush = "windows_hyperv_root_partition_io_tlb_flush"
	metricHypervRootPartitionIOTLBFlushCost = "windows_hyperv_root_partition_io_tlb_flush_cost"
	metricHypervRootPartitionVirtualTLBFlushEntries = "windows_hyperv_root_partition_virtual_tlb_flush_entires"
	metricHypervRootPartitionVirtualTLBPages = "windows_hyperv_root_partition_virtual_tlb_pages"
	metricHypervRootPartitionVirtualTLBSize = "windows_hyperv_root_partition_recommended_virtual_tlb_size"
	metricHypervRootPartitionViirtualTLBFlushEntries = "windows_hyperv_root_partition_virtual_tlb_flush_entires"
	metricHypervRootPartition = "windows_hyperv_root_partition_virtual_tlb_pages"
)

var hypervMetrics = []string {
	metricHypervHealthCritical,
	metricHypervHealthOK,
	metricHypervHypervisorLogicalProcessors,
	metricHypervHypervisorVirtualProcessors,
	metricHypervRootPartition1GDevicePages,
	metricHypervRootPartition1GGPAPages,
	metricHypervRootPartition2MDevicePages,
	metricHypervRootPartition2MGPAPages,
	metricHypervRootPartition4KDevicePages,
	metricHypervRootPartition4KGPAPages,
	metricHypervRootPartitionAddressSpace,
	metricHypervRootPartitionAttachedDevices,
	metricHypervRootPartitionDepositedPages,
	metricHypervRootPartitionDeviceDMAErrors,
	metricHypervRootPartitionDeviceInterruptErrors,
	metricHypervRootPartitionDeviceInterruptThrottleEvents,
	metricHypervRootPartitionGPASpaceModifications,
	metricHypervRootPartitionIOTLBFlush,
	metricHypervRootPartitionIOTLBFlushCost,
	metricHypervRootPartitionVirtualTLBFlushEntries,
	metricHypervRootPartitionVirtualTLBPages,
	metricHypervRootPartitionVirtualTLBSize,
	metricHypervRootPartitionViirtualTLBFlushEntries,
	metricHypervRootPartition,
}

func (w *Windows) collectHyperv(mx map[string]int64, pms prometheus.Series) {
	if !w.cache.collection[collectorHyperv] {
		w.cache.collection[collectorHyperv] = true
   		w.addHypervCharts()
	}

	for _, pm := range pms.FindByNames(hypervMetrics...) {
		name := strings.TrimPrefix(pm.Name(), "windows_")
		v := pm.Value
		mx[name] = int64(v)
	}
}

// SPDX-License-Identifier: GPL-3.0-or-later

package windows

import (
	"strings"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	metricHypervHealthCritical = "windows_hyperv_health_critical"
	metricHypervHealthOK       = "windows_hyperv_health_ok"

	metricHypervHypervisorLogicalProcessors = "windows_hyperv_hypervisor_logical_processors"
	metricHypervHypervisorVirtualProcessors = "windows_hyperv_hypervisor_virtual_processors"

	metricHypervRootPartition1GDevicePages                 = "windows_hyperv_root_partition_1G_device_pages"
	metricHypervRootPartition1GGPAPages                    = "windows_hyperv_root_partition_1G_gpa_pages"
	metricHypervRootPartition2MDevicePages                 = "windows_hyperv_root_partition_2M_device_pages"
	metricHypervRootPartition2MGPAPages                    = "windows_hyperv_root_partition_2M_gpa_pages"
	metricHypervRootPartition4KDevicePages                 = "windows_hyperv_root_partition_4K_device_pages"
	metricHypervRootPartition4KGPAPages                    = "windows_hyperv_root_partition_4K_gpa_pages"
	metricHypervRootPartitionAddressSpace                  = "windows_hyperv_root_partition_address_spaces"
	metricHypervRootPartitionAttachedDevices               = "windows_hyperv_root_partition_attached_devices"
	metricHypervRootPartitionDepositedPages                = "windows_hyperv_root_partition_deposited_pages"
	metricHypervRootPartitionDeviceDMAErrors               = "windows_hyperv_root_partition_device_dma_errors"
	metricHypervRootPartitionDeviceInterruptErrors         = "windows_hyperv_root_partition_device_interrupt_errors"
	metricHypervRootPartitionDeviceInterruptThrottleEvents = "windows_hyperv_root_partition_device_interrupt_throttle_events"
	metricHypervRootPartitionGPASpaceModifications         = "windows_hyperv_root_partition_gpa_space_modifications"
	metricHypervRootPartitionIOTLBFlush                    = "windows_hyperv_root_partition_io_tlb_flush"
	metricHypervRootPartitionIOTLBFlushCost                = "windows_hyperv_root_partition_io_tlb_flush_cost"
	metricHypervRootPartitionVirtualTLBPages               = "windows_hyperv_root_partition_virtual_tlb_pages"
	metricHypervRootPartitionVirtualTLBSize                = "windows_hyperv_root_partition_recommended_virtual_tlb_size"
	metricHypervRootPartitionViirtualTLBFlushEntires       = "windows_hyperv_root_partition_virtual_tlb_flush_entires"
	metricHypervRootPartition                              = "windows_hyperv_root_partition_virtual_tlb_pages"

	metricHypervVMDevicesBytesRead         = "windows_hyperv_vm_device_bytes_read"
	metricHypervVMDevicesOperationsRead    = "windows_hyperv_vm_device_operations_read"
	metricHypervVMDevicesBytesWritten      = "windows_hyperv_vm_device_bytes_written"
	metricHypervVMDevicesOperationsWritten = "windows_hyperv_vm_device_queue_length"
	metricHypervVMDevicesErrorCount        = "windows_hyperv_vm_device_error_count"
	metricHypervVMDevicesQueueLength       = "windows_hyperv_vm_device_operations_written"

	metricHypervVMInterfacesBytesReceived          = "windows_hyperv_vm_interface_bytes_received"
	metricHypervVMInterfacesBytesSent              = "windows_hyperv_vm_interface_bytes_sent"
	metricHypervVMInterfacesPacketsIncomingDropped = "windows_hyperv_vm_interface_packets_incoming_dropped"
	metricHypervVMInterfacesPacketsOutgoingDropped = "windows_hyperv_vm_interface_packets_outgoing_dropped"
	metricHypervVMInterfacesPacketsReceived        = "windows_hyperv_vm_interface_packets_received"
	metricHypervVMInterfacesPacketsSent            = "windows_hyperv_vm_interface_packets_sent"

	metricHypervHostCPUGuestRunTime            = "windows_hyperv_host_cpu_guest_run_time"
	metricHypervHostCPUHypervisorRunTime       = "windows_hyperv_host_cpu_hypervisor_run_time"
	metricHypervHostCPURemoteRunTime           = "windows_hyperv_host_cpu_remote_run_time"
	metricHypervHostCPUTotalRunTime            = "windows_hyperv_host_cpu_total_run_time"
	metricHypervHostLPGuestRunTimePercent      = "windows_hyperv_host_lp_guest_run_time_percent"
	metricHypervHostLPHypervisorRunTimePercent = "windows_hyperv_host_lp_hypervisor_run_time_percent"
	metricHypervHostLPTotalRunTimePercent      = "windows_hyperv_host_lp_total_run_time_percent"
)

var hypervMetrics = []string{
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
	metricHypervRootPartitionVirtualTLBPages,
	metricHypervRootPartitionVirtualTLBSize,
	metricHypervRootPartitionViirtualTLBFlushEntires,
	metricHypervRootPartition,
}

func (w *Windows) collectHyperv(mx map[string]int64, pms prometheus.Series) {
	if !w.cache.collection[collectorHyperv] {
		w.cache.collection[collectorHyperv] = true
		w.addHypervCharts()
	}

	devices := make(map[string]bool)
	interfaces := make(map[string]bool)
	cores := make(map[string]bool)
	px := "hyperv_vm_device_"

	for _, pm := range pms.FindByNames(hypervMetrics...) {
		name := strings.TrimPrefix(pm.Name(), "windows_")
		v := pm.Value
		mx[strings.ToLower(name)+"_total"] = int64(v)
	}

	for _, pm := range pms.FindByName(metricHypervVMDevicesBytesRead) {
		if name := pm.Labels.Get("vm_device"); name != "" {
			parsed_name := hypervParsenames(name)
			devices[parsed_name] = true
			mx[px+parsed_name+"_bytes_read_counter"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMDevicesOperationsRead) {
		if name := pm.Labels.Get("vm_device"); name != "" {
			parsed_name := hypervParsenames(name)
			devices[parsed_name] = true
			mx[px+parsed_name+"_operation_read_counter"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMDevicesBytesWritten) {
		if name := pm.Labels.Get("vm_device"); name != "" {
			parsed_name := hypervParsenames(name)
			devices[parsed_name] = true
			mx[px+parsed_name+"_bytes_written_counter"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMDevicesOperationsWritten) {
		if name := pm.Labels.Get("vm_device"); name != "" {
			parsed_name := hypervParsenames(name)
			devices[parsed_name] = true
			mx[px+parsed_name+"_operation_written_counter"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMDevicesErrorCount) {
		if name := pm.Labels.Get("vm_device"); name != "" {
			parsed_name := hypervParsenames(name)
			devices[parsed_name] = true
			mx[px+parsed_name+"_error_counter"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMDevicesQueueLength) {
		if name := pm.Labels.Get("vm_device"); name != "" {
			parsed_name := hypervParsenames(name)
			devices[parsed_name] = true
			mx[px+parsed_name+"_queue_length_total"] = int64(pm.Value)
		}
	}

	px = "hyperv_vm_interface_"
	for _, pm := range pms.FindByName(metricHypervVMInterfacesBytesReceived) {
		if name := pm.Labels.Get("vm_interface"); name != "" {
			parsed_name := hypervParsenames(name)
			interfaces[parsed_name] = true
			mx[px+parsed_name+"_bytes_received_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMInterfacesBytesSent) {
		if name := pm.Labels.Get("vm_interface"); name != "" {
			parsed_name := hypervParsenames(name)
			interfaces[parsed_name] = true
			mx[px+parsed_name+"_bytes_sent_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMInterfacesPacketsIncomingDropped) {
		if name := pm.Labels.Get("vm_interface"); name != "" {
			parsed_name := hypervParsenames(name)
			interfaces[parsed_name] = true
			mx[px+parsed_name+"_packets_incoming_dropped_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMInterfacesPacketsOutgoingDropped) {
		if name := pm.Labels.Get("vm_interface"); name != "" {
			parsed_name := hypervParsenames(name)
			interfaces[parsed_name] = true
			mx[px+parsed_name+"_packets_outgoing_dropped_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMInterfacesPacketsReceived) {
		if name := pm.Labels.Get("vm_interface"); name != "" {
			parsed_name := hypervParsenames(name)
			interfaces[parsed_name] = true
			mx[px+parsed_name+"_packets_received_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMInterfacesPacketsSent) {
		if name := pm.Labels.Get("vm_interface"); name != "" {
			parsed_name := hypervParsenames(name)
			interfaces[parsed_name] = true
			mx[px+parsed_name+"_packets_sent_total"] = int64(pm.Value)
		}
	}

	px = "hyperv_host_cpu_"
	for _, pm := range pms.FindByName(metricHypervHostCPUGuestRunTime) {
		if name := pm.Labels.Get("core"); name != "" {
			parsed_name := hypervParsenames(name)
			cores[parsed_name] = true
			mx[px+parsed_name+"_guest_run_time_period"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervHostCPUHypervisorRunTime) {
		if name := pm.Labels.Get("core"); name != "" {
			parsed_name := hypervParsenames(name)
			cores[parsed_name] = true
			mx[px+parsed_name+"_hypervisor_run_time_period"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervHostCPURemoteRunTime) {
		if name := pm.Labels.Get("core"); name != "" {
			parsed_name := hypervParsenames(name)
			cores[parsed_name] = true
			mx[px+parsed_name+"_remote_run_time_period"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervHostCPUTotalRunTime) {
		if name := pm.Labels.Get("core"); name != "" {
			parsed_name := hypervParsenames(name)
			cores[parsed_name] = true
			mx[px+parsed_name+"_total_run_time_period"] = int64(pm.Value)
		}
	}

	px = "hyperv_host_lp_"
	for _, pm := range pms.FindByName(metricHypervHostLPGuestRunTimePercent) {
		if name := pm.Labels.Get("core"); name != "" {
			parsed_name := hypervParsenames(name)
			cores[parsed_name] = true
			mx[px+parsed_name+"_guest_run_time_period"] = int64(pm.Value) * 100
		}
	}
	for _, pm := range pms.FindByName(metricHypervHostLPHypervisorRunTimePercent) {
		if name := pm.Labels.Get("core"); name != "" {
			parsed_name := hypervParsenames(name)
			cores[parsed_name] = true
			mx[px+parsed_name+"_hypervisor_run_time_period"] = int64(pm.Value) * 100
		}
	}
	for _, pm := range pms.FindByName(metricHypervHostLPTotalRunTimePercent) {
		if name := pm.Labels.Get("core"); name != "" {
			parsed_name := hypervParsenames(name)
			cores[parsed_name] = true
			mx[px+parsed_name+"_total_run_time_period"] = int64(pm.Value) * 100
		}
	}

	for v := range devices {
		if !w.cache.hypervDevices[v] {
			w.cache.hypervDevices[v] = true
			w.addHypervDeviceCharts(v)
		}
	}
	for v := range interfaces {
		if !w.cache.hypervInterfaces[v] {
			w.cache.hypervInterfaces[v] = true
			w.addHypervInterfaceCharts(v)
		}
	}
	for v := range cores {
		if !w.cache.hypervCores[v] {
			w.cache.hypervCores[v] = true
			w.addHypervCoreCharts(v)
		}
	}
}

func hypervParsenames(name string) string {
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "?", "_")
	name = strings.ReplaceAll(name, ":", "_")
	return strings.ToLower(name)
}

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

	metricHypervVSwitchBroadcastPacketsReceivedTotal         = "windows_hyperv_vswitch_broadcast_packets_received_total"
	metricHypervVSwitchBroadcastPacketsSentTotal             = "windows_hyperv_vswitch_broadcast_packets_sent_total"
	metricHypervVSwitchBytesReceivedTotal                    = "windows_hyperv_vswitch_bytes_received_total"
	metricHypervVSwitchBytesSentTotal                        = "windows_hyperv_vswitch_bytes_sent_total"
	metricHypervVSwitchBytesTotal                            = "windows_hyperv_vswitch_bytes_total"
	metricHypervVSwitchDirectedPacketsReceivedTotal          = "windows_hyperv_vswitch_directed_packets_received_total"
	metricHypervVSwitchDirectedPacketsSendTotal              = "windows_hyperv_vswitch_directed_packets_send_total"
	metricHypervVSwitchDroppedPacketsIncomingTotal           = "windows_hyperv_vswitch_dropped_packets_incoming_total"
	metricHypervVSwitchDroppedPacketsOutcomingTotal          = "windows_hyperv_vswitch_dropped_packets_outcoming_total"
	metricHypervVSwitchExtensionDroppedAttacksIncomingTotal  = "windows_hyperv_vswitch_extensions_dropped_packets_incoming_total"
	metricHypervVSwitchExtensionDroppedPacketsOutcomingTotal = "windows_hyperv_vswitch_extensions_dropped_packets_outcoming_total"
	metricHypervVSwitchLearnedMACAddressTotal                = "windows_hyperv_vswitch_learned_mac_addresses_total"
	metricHypervVSwitchMulticastPacketsReceivedTotal         = "windows_hyperv_vswitch_multicast_packets_received_total"
	metricHypervVSwitchMulticastPacketsSentTotal             = "windows_hyperv_vswitch_multicast_packets_sent_total"
	metricHypervVSwitchNumberOfSendChannelMovesTotal         = "windows_hyperv_vswitch_number_of_send_channel_moves_total"
	metricHypervVSwitchNumberOfVMQMovesTottal                = "windows_hyperv_vswitch_number_of_vmq_moves_total"
	metricHypervVSwitchPacketsFloodedTotal                   = "windows_hyperv_vswitch_packets_flooded_total"
	metricHypervVSwitchPacketsReceivedTotal                  = "windows_hyperv_vswitch_packets_received_total"
	metricHypervVSwitchPacketsTotal                          = "windows_hyperv_vswitch_packets_total"
	metricHypervVSwitchPurgedMACAddresses                    = "windows_hyperv_vswitch_purged_mac_addresses_total"

	metricHyperVVIDPhysicalPagesAllocated = "windows_hyperv_vid_physical_pages_allocated"
	metricHyperVVIDPreferredNumaNodeIndex = "windows_hyperv_vid_preferred_numa_node_index"
	metricHyperVVIDRemotePhysicalPages    = "windows_hyperv_vid_remote_physical_pages"

	metricHypervVMMemoryAddOperationsTotal   = "windows_hyperv_vm_memory_add_operations_total"
	metricHypervVMMemoryAddedTotal           = "windows_hyperv_vm_memory_added_total"
	metricHypervVMMemoryPhysical             = "windows_hyperv_vm_memory_physical"
	metricHypervVMMemoryPhysicalGuestVisible = "windows_hyperv_vm_memory_physical_guest_visible"
	metricHypervVMMemoryPressureAverage      = "windows_hyperv_vm_memory_pressure_average"
	metricHypervVMMemoryPressureCurrent      = "windows_hyperv_vm_memory_pressure_current"
	metricHypervVMMemoryPressureMaximum      = "windows_hyperv_vm_memory_pressure_maximum"
	metricHypervVMMemoryPressureMinimum      = "windows_hyperv_vm_memory_pressure_minimum"
	metricHypervVMMemoryRemoveOperatonsTotal = "windows_hyperv_vm_memory_remove_operations_total"
	metricHypervVMMemoryRemovedTotal         = "windows_hyperv_vm_memory_removed_total"

	metricsHypervVMCPUGuestRunTime      = "windows_hyperv_vm_cpu_guest_run_time"
	metricsHypervVMCPUHypervisorRunTime = "windows_hyperv_vm_cpu_hypervisor_run_time"
	metricsHypervVMCPURemoteRunTime     = "windows_hyperv_vm_cpu_remote_run_time"
	metricsHypervVMCPUTotalRunTime      = "windows_hyperv_vm_cpu_total_run_time"
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
	vswitches := make(map[string]bool)
	vm := make(map[string]bool)

	px := "hyperv_vm_device_"
	for _, pm := range pms.FindByNames(hypervMetrics...) {
		name := strings.TrimPrefix(pm.Name(), "windows_")
		v := pm.Value
		mx[strings.ToLower(name)+"_total"] = int64(v)
	}

	for _, pm := range pms.FindByName(metricHypervVMDevicesBytesRead) {
		if name := pm.Labels.Get("vm_device"); name != "" {
			parsed_name := hypervParseNames(name)
			devices[parsed_name] = true
			mx[px+parsed_name+"_bytes_read_counter"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMDevicesOperationsRead) {
		if name := pm.Labels.Get("vm_device"); name != "" {
			parsed_name := hypervParseNames(name)
			devices[parsed_name] = true
			mx[px+parsed_name+"_operation_read_counter"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMDevicesBytesWritten) {
		if name := pm.Labels.Get("vm_device"); name != "" {
			parsed_name := hypervParseNames(name)
			devices[parsed_name] = true
			mx[px+parsed_name+"_bytes_written_counter"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMDevicesOperationsWritten) {
		if name := pm.Labels.Get("vm_device"); name != "" {
			parsed_name := hypervParseNames(name)
			devices[parsed_name] = true
			mx[px+parsed_name+"_operation_written_counter"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMDevicesErrorCount) {
		if name := pm.Labels.Get("vm_device"); name != "" {
			parsed_name := hypervParseNames(name)
			devices[parsed_name] = true
			mx[px+parsed_name+"_error_counter"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMDevicesQueueLength) {
		if name := pm.Labels.Get("vm_device"); name != "" {
			parsed_name := hypervParseNames(name)
			devices[parsed_name] = true
			mx[px+parsed_name+"_queue_length_total"] = int64(pm.Value)
		}
	}

	px = "hyperv_vm_interface_"
	for _, pm := range pms.FindByName(metricHypervVMInterfacesBytesReceived) {
		if name := pm.Labels.Get("vm_interface"); name != "" {
			parsed_name := hypervParseNames(name)
			interfaces[parsed_name] = true
			mx[px+parsed_name+"_bytes_received_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMInterfacesBytesSent) {
		if name := pm.Labels.Get("vm_interface"); name != "" {
			parsed_name := hypervParseNames(name)
			interfaces[parsed_name] = true
			mx[px+parsed_name+"_bytes_sent_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMInterfacesPacketsIncomingDropped) {
		if name := pm.Labels.Get("vm_interface"); name != "" {
			parsed_name := hypervParseNames(name)
			interfaces[parsed_name] = true
			mx[px+parsed_name+"_packets_incoming_dropped_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMInterfacesPacketsOutgoingDropped) {
		if name := pm.Labels.Get("vm_interface"); name != "" {
			parsed_name := hypervParseNames(name)
			interfaces[parsed_name] = true
			mx[px+parsed_name+"_packets_outgoing_dropped_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMInterfacesPacketsReceived) {
		if name := pm.Labels.Get("vm_interface"); name != "" {
			parsed_name := hypervParseNames(name)
			interfaces[parsed_name] = true
			mx[px+parsed_name+"_packets_received_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMInterfacesPacketsSent) {
		if name := pm.Labels.Get("vm_interface"); name != "" {
			parsed_name := hypervParseNames(name)
			interfaces[parsed_name] = true
			mx[px+parsed_name+"_packets_sent_total"] = int64(pm.Value)
		}
	}

	px = "hyperv_host_cpu_"
	for _, pm := range pms.FindByName(metricHypervHostCPUGuestRunTime) {
		if name := pm.Labels.Get("core"); name != "" {
			parsed_name := hypervParseNames(name)
			cores[parsed_name] = true
			mx[px+parsed_name+"_guest_run_time_period"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervHostCPUHypervisorRunTime) {
		if name := pm.Labels.Get("core"); name != "" {
			parsed_name := hypervParseNames(name)
			cores[parsed_name] = true
			mx[px+parsed_name+"_hypervisor_run_time_period"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervHostCPURemoteRunTime) {
		if name := pm.Labels.Get("core"); name != "" {
			parsed_name := hypervParseNames(name)
			cores[parsed_name] = true
			mx[px+parsed_name+"_remote_run_time_period"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervHostCPUTotalRunTime) {
		if name := pm.Labels.Get("core"); name != "" {
			parsed_name := hypervParseNames(name)
			cores[parsed_name] = true
			mx[px+parsed_name+"_total_run_time_period"] = int64(pm.Value)
		}
	}

	px = "hyperv_host_lp_"
	for _, pm := range pms.FindByName(metricHypervHostLPGuestRunTimePercent) {
		if name := pm.Labels.Get("core"); name != "" {
			parsed_name := hypervParseNames(name)
			cores[parsed_name] = true
			mx[px+parsed_name+"_guest_run_time_period"] = int64(pm.Value) * 100
		}
	}
	for _, pm := range pms.FindByName(metricHypervHostLPHypervisorRunTimePercent) {
		if name := pm.Labels.Get("core"); name != "" {
			parsed_name := hypervParseNames(name)
			cores[parsed_name] = true
			mx[px+parsed_name+"_hypervisor_run_time_period"] = int64(pm.Value) * 100
		}
	}
	for _, pm := range pms.FindByName(metricHypervHostLPTotalRunTimePercent) {
		if name := pm.Labels.Get("core"); name != "" {
			parsed_name := hypervParseNames(name)
			cores[parsed_name] = true
			mx[px+parsed_name+"_total_run_time_period"] = int64(pm.Value) * 100
		}
	}

	px = "hyperv_vswitch_"
	for _, pm := range pms.FindByName(metricHypervVSwitchBroadcastPacketsReceivedTotal) {
		if name := pm.Labels.Get("vswitch"); name != "" {
			parsed_name := hypervParseNames(name)
			vswitches[parsed_name] = true
			mx[px+parsed_name+"_broadcast_packets_received_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVSwitchBroadcastPacketsSentTotal) {
		if name := pm.Labels.Get("vswitch"); name != "" {
			parsed_name := hypervParseNames(name)
			vswitches[parsed_name] = true
			mx[px+parsed_name+"_broadcast_packets_sent_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVSwitchBytesReceivedTotal) {
		if name := pm.Labels.Get("vswitch"); name != "" {
			parsed_name := hypervParseNames(name)
			vswitches[parsed_name] = true
			mx[px+parsed_name+"_bytes_received_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVSwitchBytesSentTotal) {
		if name := pm.Labels.Get("vswitch"); name != "" {
			parsed_name := hypervParseNames(name)
			vswitches[parsed_name] = true
			mx[px+parsed_name+"_bytes_sent_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVSwitchBytesTotal) {
		if name := pm.Labels.Get("vswitch"); name != "" {
			parsed_name := hypervParseNames(name)
			vswitches[parsed_name] = true
			mx[px+parsed_name+"_bytes_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVSwitchDirectedPacketsReceivedTotal) {
		if name := pm.Labels.Get("vswitch"); name != "" {
			parsed_name := hypervParseNames(name)
			vswitches[parsed_name] = true
			mx[px+parsed_name+"_directed_packets_received_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVSwitchDirectedPacketsSendTotal) {
		if name := pm.Labels.Get("vswitch"); name != "" {
			parsed_name := hypervParseNames(name)
			vswitches[parsed_name] = true
			mx[px+parsed_name+"_directed_packets_send_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVSwitchDroppedPacketsIncomingTotal) {
		if name := pm.Labels.Get("vswitch"); name != "" {
			parsed_name := hypervParseNames(name)
			vswitches[parsed_name] = true
			mx[px+parsed_name+"_dropped_packets_incoming_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVSwitchDroppedPacketsOutcomingTotal) {
		if name := pm.Labels.Get("vswitch"); name != "" {
			parsed_name := hypervParseNames(name)
			vswitches[parsed_name] = true
			mx[px+parsed_name+"_dropped_packets_outcoming_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVSwitchExtensionDroppedAttacksIncomingTotal) {
		if name := pm.Labels.Get("vswitch"); name != "" {
			parsed_name := hypervParseNames(name)
			vswitches[parsed_name] = true
			mx[px+parsed_name+"_extensions_dropped_packets_incoming_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVSwitchExtensionDroppedPacketsOutcomingTotal) {
		if name := pm.Labels.Get("vswitch"); name != "" {
			parsed_name := hypervParseNames(name)
			vswitches[parsed_name] = true
			mx[px+parsed_name+"_extensions_dropped_packets_outcoming_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVSwitchLearnedMACAddressTotal) {
		if name := pm.Labels.Get("vswitch"); name != "" {
			parsed_name := hypervParseNames(name)
			vswitches[parsed_name] = true
			mx[px+parsed_name+"_learned_mac_addresses_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVSwitchMulticastPacketsReceivedTotal) {
		if name := pm.Labels.Get("vswitch"); name != "" {
			parsed_name := hypervParseNames(name)
			vswitches[parsed_name] = true
			mx[px+parsed_name+"_multicast_packets_received_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVSwitchMulticastPacketsSentTotal) {
		if name := pm.Labels.Get("vswitch"); name != "" {
			parsed_name := hypervParseNames(name)
			vswitches[parsed_name] = true
			mx[px+parsed_name+"_multicast_packets_sent_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVSwitchNumberOfSendChannelMovesTotal) {
		if name := pm.Labels.Get("vswitch"); name != "" {
			parsed_name := hypervParseNames(name)
			vswitches[parsed_name] = true
			mx[px+parsed_name+"_number_of_send_channel_moves_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVSwitchNumberOfVMQMovesTottal) {
		if name := pm.Labels.Get("vswitch"); name != "" {
			parsed_name := hypervParseNames(name)
			vswitches[parsed_name] = true
			mx[px+parsed_name+"_number_of_vmq_moves_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVSwitchPacketsFloodedTotal) {
		if name := pm.Labels.Get("vswitch"); name != "" {
			parsed_name := hypervParseNames(name)
			vswitches[parsed_name] = true
			mx[px+parsed_name+"_packets_flooded_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVSwitchPacketsReceivedTotal) {
		if name := pm.Labels.Get("vswitch"); name != "" {
			parsed_name := hypervParseNames(name)
			vswitches[parsed_name] = true
			mx[px+parsed_name+"_packets_received_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVSwitchPacketsTotal) {
		if name := pm.Labels.Get("vswitch"); name != "" {
			parsed_name := hypervParseNames(name)
			vswitches[parsed_name] = true
			mx[px+parsed_name+"_packets_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVSwitchPurgedMACAddresses) {
		if name := pm.Labels.Get("vswitch"); name != "" {
			parsed_name := hypervParseNames(name)
			vswitches[parsed_name] = true
			mx[px+parsed_name+"_purged_mac_addresses"] = int64(pm.Value)
		}
	}

	px = "hyperv_vid_"
	for _, pm := range pms.FindByName(metricHyperVVIDPhysicalPagesAllocated) {
		if name := pm.Labels.Get("vm"); name != "" {
			parsed_name := hypervParseNames(name)
			vm[parsed_name] = true
			mx[px+parsed_name+"_vid_physical_pages_allocated_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHyperVVIDPreferredNumaNodeIndex) {
		if name := pm.Labels.Get("vm"); name != "" {
			parsed_name := hypervParseNames(name)
			vm[parsed_name] = true
			mx[px+parsed_name+"_vid_preferred_numa_node_index"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHyperVVIDRemotePhysicalPages) {
		if name := pm.Labels.Get("vm"); name != "" {
			parsed_name := hypervParseNames(name)
			vm[parsed_name] = true
			mx[px+parsed_name+"_vid_remote_physical_page"] = int64(pm.Value)
		}
	}

	px = "hyperv_vm_"
	for _, pm := range pms.FindByName(metricHypervVMMemoryAddOperationsTotal) {
		if name := pm.Labels.Get("vm"); name != "" {
			parsed_name := hypervParseNames(name)
			vm[parsed_name] = true
			mx[px+parsed_name+"_memory_add_operations_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMMemoryAddedTotal) {
		if name := pm.Labels.Get("vm"); name != "" {
			parsed_name := hypervParseNames(name)
			vm[parsed_name] = true
			mx[px+parsed_name+"_memory_added_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMMemoryPhysical) {
		if name := pm.Labels.Get("vm"); name != "" {
			parsed_name := hypervParseNames(name)
			vm[parsed_name] = true
			mx[px+parsed_name+"_memory_physical"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMMemoryPhysicalGuestVisible) {
		if name := pm.Labels.Get("vm"); name != "" {
			parsed_name := hypervParseNames(name)
			vm[parsed_name] = true
			mx[px+parsed_name+"_memory_physical_guest_visible"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMMemoryPressureAverage) {
		if name := pm.Labels.Get("vm"); name != "" {
			parsed_name := hypervParseNames(name)
			vm[parsed_name] = true
			mx[px+parsed_name+"_memory_pressure_average"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMMemoryPressureCurrent) {
		if name := pm.Labels.Get("vm"); name != "" {
			parsed_name := hypervParseNames(name)
			vm[parsed_name] = true
			mx[px+parsed_name+"_memory_pressure_current"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMMemoryPressureMaximum) {
		if name := pm.Labels.Get("vm"); name != "" {
			parsed_name := hypervParseNames(name)
			vm[parsed_name] = true
			mx[px+parsed_name+"_memory_pressure_maximum"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMMemoryPressureMinimum) {
		if name := pm.Labels.Get("vm"); name != "" {
			parsed_name := hypervParseNames(name)
			vm[parsed_name] = true
			mx[px+parsed_name+"_memory_pressure_minimum"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMMemoryRemoveOperatonsTotal) {
		if name := pm.Labels.Get("vm"); name != "" {
			parsed_name := hypervParseNames(name)
			vm[parsed_name] = true
			mx[px+parsed_name+"_memory_remove_operations_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricHypervVMMemoryRemovedTotal) {
		if name := pm.Labels.Get("vm"); name != "" {
			parsed_name := hypervParseNames(name)
			vm[parsed_name] = true
			mx[px+parsed_name+"_memory_removed_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricsHypervVMCPUGuestRunTime) {
		if vmname := pm.Labels.Get("vm"); vmname != "" {
			parsed_name := hypervParseNames(vmname)
			if cpu := pm.Labels.Get("core"); cpu != "" {
				vm[parsed_name] = true
				mx[px+parsed_name+"_cpu_"+cpu+"_guest_run_time_total"] = int64(pm.Value)
			}
		}
	}
	for _, pm := range pms.FindByName(metricsHypervVMCPUHypervisorRunTime) {
		if vmname := pm.Labels.Get("vm"); vmname != "" {
			parsed_name := hypervParseNames(vmname)
			if cpu := pm.Labels.Get("core"); cpu != "" {
				vm[parsed_name] = true
				mx[px+parsed_name+"_cpu_"+cpu+"_hypervisor_run_time_total"] = int64(pm.Value)
			}
		}
	}
	for _, pm := range pms.FindByName(metricsHypervVMCPURemoteRunTime) {
		if vmname := pm.Labels.Get("vm"); vmname != "" {
			parsed_name := hypervParseNames(vmname)
			if cpu := pm.Labels.Get("core"); cpu != "" {
				vm[parsed_name] = true
				mx[px+parsed_name+"_cpu_"+cpu+"_remote_run_time_total"] = int64(pm.Value)
			}
		}
	}
	for _, pm := range pms.FindByName(metricsHypervVMCPUTotalRunTime) {
		if vmname := pm.Labels.Get("vm"); vmname != "" {
			parsed_name := hypervParseNames(vmname)
			if cpu := pm.Labels.Get("core"); cpu != "" {
				vm[parsed_name] = true
				mx[px+parsed_name+"_cpu_"+cpu+"_run_time_total"] = int64(pm.Value)
			}
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
	for v := range vswitches {
		if !w.cache.hypervVswitch[v] {
			w.cache.hypervVswitch[v] = true
			w.addHypervVSwitchCharts(v)
		}
	}
	for v := range vm {
		if !w.cache.hypervVM[v] {
			w.cache.hypervVM[v] = true
			w.addHypervVIDCharts(v)
			w.addHypervVMCharts(v)
			for cpu := range cores {
				w.addHypervVMCPUCharts(v, cpu)
			}
		}
	}
}

func hypervParseNames(name string) string {
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "?", "_")
	name = strings.ReplaceAll(name, ":", "_")
	return strings.ToLower(name)
}

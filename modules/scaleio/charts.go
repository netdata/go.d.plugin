package scaleio

import (
	"fmt"

	"github.com/netdata/go.d.plugin/modules/scaleio/client"

	"github.com/netdata/go-orchestrator/module"
)

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

var systemCharts = Charts{
	// Capacity
	{
		ID:    "system_capacity_total",
		Title: "System Total Capacity",
		Units: "KB",
		Fam:   "capacity",
		Ctx:   "scaleio.system_capacity_total",
		Dims: Dims{
			{ID: "system_capacity_max_capacity", Name: "total"},
		},
	},
	{
		ID:    "system_capacity",
		Title: "System Capacity",
		Units: "KB",
		Fam:   "capacity",
		Type:  module.Stacked,
		Ctx:   "scaleio.system_capacity",
		Dims: Dims{
			{ID: "system_capacity_protected", Name: "protected"},
			{ID: "system_capacity_degraded", Name: "degraded"},
			{ID: "system_capacity_spare", Name: "spare"},
			{ID: "system_capacity_failed", Name: "failed"},
			{ID: "system_capacity_decreased", Name: "decreased"},
			{ID: "system_capacity_unreachable_unused", Name: "unavailable"},
			{ID: "system_capacity_in_maintenance", Name: "in_maintenance"},
			{ID: "system_capacity_unused", Name: "unused"},
		},
	},
	{
		ID:    "system_capacity_available_volume_allocation",
		Title: "Available For Volume Allocation",
		Units: "KB",
		Fam:   "capacity",
		Ctx:   "scaleio.system_capacity_available_volume_allocation",
		Dims: Dims{
			{ID: "system_capacity_available_for_volume_allocation", Name: "available"},
		},
	},
	{
		ID:    "system_capacity_volume_usage_by_type",
		Title: "Volume Usage By Type",
		Units: "KB",
		Fam:   "capacity",
		Ctx:   "scaleio.system_capacity_volume_usage_by_type",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "system_capacity_thick_in_use", Name: "thick"},
			{ID: "system_capacity_thin_in_use", Name: "thin"},
		},
	},
	{
		ID:    "system_capacity_thin_volume_usage",
		Title: "Thin Volume Usage",
		Units: "KB",
		Fam:   "capacity",
		Ctx:   "scaleio.system_capacity_thin_volume_usage",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "system_capacity_thin_free", Name: "free"},
			{ID: "system_capacity_thin_in_use", Name: "used"},
		},
	},
	// I/O Workload BW
	{
		ID:    "system_workload_primary_bandwidth_total",
		Title: "Primary Backend Bandwidth Total (Read and Write)",
		Units: "KB/s",
		Fam:   "i/o workload",
		Ctx:   "scaleio.system_workload_primary_bandwidth_total",
		Dims: Dims{
			{ID: "system_backend_primary_bandwidth_read_write", Name: "total", Div: 1000},
		},
	},
	{
		ID:    "system_workload_primary_bandwidth",
		Title: "Primary Backend Bandwidth",
		Units: "KB/s",
		Fam:   "i/o workload",
		Ctx:   "scaleio.system_workload_primary_bandwidth",
		Type:  module.Area,
		Dims: Dims{
			{ID: "system_backend_primary_bandwidth_read", Name: "read", Div: 1000},
			{ID: "system_backend_primary_bandwidth_write", Name: "write", Mul: -1, Div: 1000},
		},
	},
	// I/O Workload IOPS
	{
		ID:    "system_workload_primary_iops_total",
		Title: "Primary Backend IOPS Total (Read and Write)",
		Units: "iops/s",
		Fam:   "i/o workload",
		Ctx:   "scaleio.system_workload_primary_iops_total",
		Dims: Dims{
			{ID: "system_backend_primary_iops_read_write", Name: "total", Div: 1000},
		},
	},
	{
		ID:    "system_workload_primary_iops",
		Title: "Primary Backend IOPS",
		Units: "iops/s",
		Fam:   "i/o workload",
		Ctx:   "scaleio.system_workload_primary_iops",
		Type:  module.Area,
		Dims: Dims{
			{ID: "system_backend_primary_iops_read", Name: "read", Div: 1000},
			{ID: "system_backend_primary_iops_write", Name: "write", Mul: -1, Div: 1000},
		},
	},
	{
		ID:    "system_workload_primary_io_size_total",
		Title: "Primary Backend I/O Size Total (Read and Write)",
		Units: "KB",
		Fam:   "i/o workload",
		Ctx:   "scaleio.system_workload_primary_io_size_total",
		Dims: Dims{
			{ID: "system_backend_primary_io_size_read_write", Name: "io_size", Div: 1000},
		},
	},
	// Rebalance
	{
		ID:    "system_rebalance",
		Title: "System Rebalance",
		Units: "KB/s",
		Fam:   "rebalance",
		Type:  module.Area,
		Ctx:   "scaleio.system_rebalance",
		Dims: Dims{
			{ID: "system_rebalance_bandwidth_read", Name: "read", Div: 1000},
			{ID: "system_rebalance_bandwidth_write", Name: "write", Mul: -1, Div: 1000},
		},
	},
	{
		ID:    "system_rebalance_left",
		Title: "System Rebalance Pending Capacity",
		Units: "KB",
		Fam:   "rebalance",
		Ctx:   "scaleio.system_rebalance_left",
		Dims: Dims{
			{ID: "system_rebalance_pending_capacity_in_Kb", Name: "left"},
		},
	},
	// Rebuild
	{
		ID:    "system_rebuild",
		Title: "System Rebuild Bandwidth Total (Forward, Backward and Normal)",
		Units: "KB/s",
		Fam:   "rebuild",
		Ctx:   "scaleio.system_rebuild",
		Type:  module.Area,
		Dims: Dims{
			{ID: "system_rebuild_total_bandwidth_read", Name: "read", Div: 1000},
			{ID: "system_rebuild_total_bandwidth_write", Name: "write", Mul: -1, Div: 1000},
		},
	},
	{
		ID:    "system_rebuild_left",
		Title: "System Rebuild Pending Capacity Total (Forward, Backward and Normal)",
		Units: "KB",
		Fam:   "rebuild",
		Ctx:   "scaleio.system_rebuild_left",
		Dims: Dims{
			{ID: "system_rebuild_total_pending_capacity_in_Kb", Name: "left"},
		},
	},
	// Components
	{
		ID:    "system_defined_components",
		Title: "System Defined Components",
		Units: "number",
		Fam:   "components",
		Ctx:   "scaleio.system_defined_components",
		//Type:  module.Stacked,
		Dims: Dims{
			{ID: "system_num_of_devices", Name: "devices"},
			{ID: "system_num_of_fault_sets", Name: "fault_sets"},
			{ID: "system_num_of_protection_domains", Name: "protection_domains"},
			{ID: "system_num_of_rfcache_devices", Name: "rfcache_devices"},
			{ID: "system_num_of_sdc", Name: "sdc"},
			{ID: "system_num_of_sds", Name: "sds"},
			{ID: "system_num_of_snapshots", Name: "snapshots"},
			{ID: "system_num_of_storage_pools", Name: "storage_pools"},
			{ID: "system_num_of_volumes", Name: "volumes"},
			{ID: "system_num_of_vtrees", Name: "vtrees"},
		},
	},
	{
		ID:    "system_components_volumes_by_type",
		Title: "Volumes By Type",
		Units: "number",
		Fam:   "components",
		Ctx:   "scaleio.system_components_volumes_by_type",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "system_num_of_thick_base_volumes", Name: "thick"},
			{ID: "system_num_of_thin_base_volumes", Name: "thin"},
		},
	},
	{
		ID:    "system_components_volumes_by_mapping",
		Title: "Volumes By Mapping",
		Units: "number",
		Fam:   "components",
		Ctx:   "scaleio.system_components_volumes_by_mapping",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "system_num_of_mapped_volumes", Name: "mapped"},
			{ID: "system_num_of_unmapped_volumes", Name: "unmapped"},
		},
	},
}

var storagePoolCharts = Charts{
	{
		ID:    "storage_pool_%s_capacity_total",
		Title: "Total Capacity",
		Units: "KB",
		Fam:   "pool %s",
		Ctx:   "scaleio.storage_pool_%s_capacity_total",
		Dims: Dims{
			{ID: "storage_pool_%s_capacity_max_capacity", Name: "total"},
		},
	},
	{
		ID:    "storage_pool_%s_capacity",
		Title: "Capacity",
		Units: "KB",
		Fam:   "pool %s",
		Type:  module.Stacked,
		Ctx:   "scaleio.storage_pool_%s_capacity",
		Dims: Dims{
			{ID: "storage_pool_%s_capacity_protected", Name: "protected"},
			{ID: "storage_pool_%s_capacity_degraded", Name: "degraded"},
			{ID: "storage_pool_%s_capacity_spare", Name: "spare"},
			{ID: "storage_pool_%s_capacity_failed", Name: "failed"},
			{ID: "storage_pool_%s_capacity_decreased", Name: "decreased"},
			{ID: "storage_pool_%s_capacity_unreachable_unused", Name: "unavailable"},
			{ID: "storage_pool_%s_capacity_in_maintenance", Name: "in_maintenance"},
			{ID: "storage_pool_%s_capacity_unused", Name: "unused"},
		},
	},
	{
		ID:    "storage_pool_%s_available_volume_allocation",
		Title: "Available For Volume Allocation",
		Units: "KB",
		Fam:   "pool %s",
		Ctx:   "scaleio.storage_pool_%s_available_volume_allocation",
		Dims: Dims{
			{ID: "storage_pool_%s_capacity_available_for_volume_allocation", Name: "available"},
		},
	},
	{
		ID:    "storage_pool_%s_components",
		Title: "Components",
		Units: "number",
		Fam:   "pool %s",
		Ctx:   "scaleio.storage_pool_%s_components",
		//Type:  module.Stacked,
		Dims: Dims{
			{ID: "storage_pool_%s_num_of_devices", Name: "devices"},
			{ID: "storage_pool_%s_num_of_snapshots", Name: "snapshots"},
			{ID: "storage_pool_%s_num_of_volumes", Name: "volumes"},
			{ID: "storage_pool_%s_num_of_vtrees", Name: "vtrees"},
		},
	},
}

var sdcCharts = Charts{
	{
		ID:    "sdc_%s_mdm_connection_state",
		Title: "MDM Connection State",
		Units: "bool",
		Fam:   "sdc %s",
		Ctx:   "scaleio.sdc_%s_mdm_connection_state",
		Dims: Dims{
			{ID: "sdc_%s_mdm_connection_state", Name: "connected"},
		},
	},
	{
		ID:    "sdc_%s_bandwidth",
		Title: "Bandwidth",
		Units: "KB/s",
		Fam:   "sdc %s",
		Ctx:   "scaleio.sdc_%s_bandwidth",
		Type:  module.Area,
		Dims: Dims{
			{ID: "sdc_%s_bandwidth_read", Name: "read", Div: 1000},
			{ID: "sdc_%s_bandwidth_write", Name: "write", Mul: -1, Div: 1000},
		},
	},
	{
		ID:    "sdc_%s_iops",
		Title: "IOPS",
		Units: "iops/s",
		Fam:   "sdc %s",
		Ctx:   "scaleio.sdc_%s_iops",
		Type:  module.Area,
		Dims: Dims{
			{ID: "sdc_%s_iops_read", Name: "read", Div: 1000},
			{ID: "sdc_%s_iops_write", Name: "write", Mul: -1, Div: 1000},
		},
	},
	{
		ID:    "sdc_%s_io_size",
		Title: "I/O Size",
		Units: "KB",
		Fam:   "sdc %s",
		Ctx:   "scaleio.sdc_%s_io_size",
		Type:  module.Area,
		Dims: Dims{
			{ID: "sdc_%s_io_size_read", Name: "read", Div: 1000},
			{ID: "sdc_%s_io_size_write", Name: "write", Mul: -1, Div: 1000},
		},
	},
	{
		ID:    "sdc_%s_num_of_mapped_volumed",
		Title: "Mapped Volumes",
		Units: "volumes",
		Fam:   "sdc %s",
		Ctx:   "scaleio.sdc_%s_num_of_mapped_volumed",
		Dims: Dims{
			{ID: "sdc_%s_num_of_mapped_volumes", Name: "mapped"},
		},
	},
}

func newStoragePoolCharts(pool client.StoragePool) *Charts {
	charts := storagePoolCharts.Copy()
	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, pool.ID)
		chart.Fam = fmt.Sprintf(chart.Fam, pool.Name)
		chart.Ctx = fmt.Sprintf(chart.Ctx, pool.ID)
		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, pool.ID)
		}
	}
	return charts
}

func newSdcCharts(sdc client.Sdc) *Charts {
	charts := sdcCharts.Copy()
	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, sdc.ID)
		chart.Fam = fmt.Sprintf(chart.Fam, sdc.SdcIp)
		chart.Ctx = fmt.Sprintf(chart.Ctx, sdc.ID)
		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, sdc.ID)
		}
	}
	return charts
}

// TODO: remove stale charts?
func (s *ScaleIO) updateCharts() {
	s.updateStoragePoolCharts()
	s.updateSdcCharts()
}

func (s *ScaleIO) updateStoragePoolCharts() {
	for k, v := range s.discovered.pool {
		if s.charted.pool[k] {
			continue
		}
		s.charted.pool[k] = true
		s.addStoragePoolCharts(v)
	}
}

func (s *ScaleIO) updateSdcCharts() {
	for k, v := range s.discovered.sdc {
		if s.charted.sdc[k] {
			continue
		}
		s.charted.sdc[k] = true
		s.addSdcCharts(v)
	}
}

func (s *ScaleIO) addStoragePoolCharts(pool client.StoragePool) {
	charts := newStoragePoolCharts(pool)
	if err := s.Charts().Add(*charts...); err != nil {
		s.Warningf("couldn't add charts for storage pool '%s(%s)': %v", pool.ID, pool.Name, err)
	}
}

func (s *ScaleIO) addSdcCharts(sdc client.Sdc) {
	charts := newSdcCharts(sdc)
	if err := s.Charts().Add(*charts...); err != nil {
		s.Warningf("couldn't add charts for sdc '%s(%s)': %v", sdc.ID, sdc.SdcIp, err)
	}
}

package scaleio

import "github.com/netdata/go-orchestrator/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

var charts = Charts{
	// Capacity
	{
		ID:    "system_capacity_total",
		Title: "Total Capacity",
		Units: "KB",
		Fam:   "capacity",
		Ctx:   "scaleio.system_capacity_total",
		Dims: Dims{
			{ID: "system_capacity_max_capacity", Name: "total"},
		},
	},
	{
		ID:    "system_capacity",
		Title: "Capacity",
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
		ID:    "system_workload_primary_io_size_total",
		Title: "Primary Backend I/O Size Total (Read and Write)",
		Units: "KB",
		Fam:   "i/o workload",
		Ctx:   "scaleio.system_workload_primary_io_size_total",
		Dims: Dims{
			{ID: "system_backend_primary_io_size_read_write", Name: "io_size", Div: 1000},
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
	// Rebalance
	{
		ID:    "system_rebalance",
		Title: "Rebalance",
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
		Title: "Rebalance Pending Capacity",
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
		Title: "Rebuild Bandwidth Total (Forward, Backward and Normal)",
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
		Title: "Rebuild Pending Capacity Total (Forward, Backward and Normal)",
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
		Title: "Defined Components",
		Units: "number",
		Fam:   "components",
		Ctx:   "scaleio.system_defined_components",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "system_num_of_devices", Name: "devices"},
			{ID: "system_num_of_fault_sets", Name: "fault_sets"},
			{ID: "system_num_of_protection_domains", Name: "protection_domains"},
			{ID: "system_num_of_rfcache_devices", Name: "rfcache_devices"},
			{ID: "system_num_of_scsi_initiators", Name: "scsi_initiators"},
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

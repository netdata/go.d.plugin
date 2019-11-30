package scaleio

type metrics struct {
	SystemOverview struct {
		Capacity   systemCapacity   `stm:"capacity"`
		IOWorkload systemWorkload   `stm:""`
		Rebalance  systemRebalance  `stm:"rebalance"`
		Rebuild    systemRebuild    `stm:"rebuild"`
		Components systemComponents `stm:"num_of"`
	} `stm:"system"`
	Sdc         map[string]sdcStatistics         `stm:"sdc"`
	StoragePool map[string]storagePoolStatistics `stm:"storage_pool"`
}

type (
	systemCapacity struct {
		MaxCapacity int64 `stm:"max_capacity"`
		ThickInUse  int64 `stm:"thick_in_use"`
		ThinInUse   int64 `stm:"thin_in_use"`
		Snapshot    int64 `stm:"snapshot"`
		Spare       int64 `stm:"spare"`
		Decreased   int64 `stm:"decreased"` // not in statistics, should be calculated
		Unused      int64 `stm:"unused"`

		InUse                        int64   `stm:"in_use"`
		AvailableForVolumeAllocation int64   `stm:"available_for_volume_allocation"`
		Utilization                  float64 `stm:"utilization,100,1"`

		Protected         int64 `stm:"protected"`
		InMaintenance     int64 `stm:"in_maintenance"`
		Degraded          int64 `stm:"degraded"`
		Failed            int64 `stm:"failed"`
		UnreachableUnused int64 `stm:"unreachable_unused"`
	}
	systemComponents struct {
		Devices            int64 `stm:"devices"`
		FaultSets          int64 `stm:"fault_sets"`
		ProtectionDomains  int64 `stm:"protection_domains"`
		RfcacheDevices     int64 `stm:"rfcache_devices"`
		Sdc                int64 `stm:"sdc"`
		Sds                int64 `stm:"sds"`
		Snapshots          int64 `stm:"snapshots"`
		StoragePools       int64 `stm:"storage_pools"`
		MappedToAllVolumes int64 `stm:"mapped_to_all_volumes"`
		ThickBaseVolumes   int64 `stm:"thick_base_volumes"`
		ThinBaseVolumes    int64 `stm:"thin_base_volumes"`
		UnmappedVolumes    int64 `stm:"unmapped_volumes"`
		MappedVolumes      int64 `stm:"mapped_volumes"`
		Volumes            int64 `stm:"volumes"`
		VTrees             int64 `stm:"vtrees"`
	}
	systemWorkload struct {
		Total   bwIOPS `stm:"total"`
		Backend struct {
			Total     bwIOPS `stm:"total"`
			Primary   bwIOPS `stm:"primary"`
			Secondary bwIOPS `stm:"secondary"`
		} `stm:"backend"`
		Frontend bwIOPS `stm:"frontend_user_data"`
	}
	systemRebalance struct {
		TimeUntilFinish float64 `stm:"time_until_finish"`
		bwIOPSPending   `stm:""`
	}
	systemRebuild struct {
		Total    bwIOPSPending `stm:"total"`
		Forward  bwIOPSPending `stm:"forward"`
		Backward bwIOPSPending `stm:"backward"`
		Normal   bwIOPSPending `stm:"normal"`
	}
)

type (
	sdcStatistics struct {
		bwIOPS             `stm:""`
		MappedVolumes      int64 `stm:"num_of_mapped_volumes"`
		MDMConnectionState bool  `stm:"mdm_connection_state"`
	}
)

type (
	storagePoolStatistics struct {
		Capacity   storagePoolCapacity `stm:"capacity"`
		Components struct {
			Devices   int64 `stm:"devices"`
			Volumes   int64 `stm:"volumes"`
			Vtrees    int64 `stm:"vtrees"`
			Snapshots int64 `stm:"snapshots"`
		} `stm:"num_of"`
	}
	storagePoolCapacity = systemCapacity
)

type (
	readWrite struct {
		Read      float64 `stm:"read,1000,1"`
		Write     float64 `stm:"write,1000,1"`
		ReadWrite float64 `stm:"read_write,1000,1"`
	}
	bwIOPS struct {
		BW     readWrite `stm:"bandwidth"`
		IOPS   readWrite `stm:"iops"`
		IOSize readWrite `stm:"io_size"`
	}
	bwIOPSPending struct {
		bwIOPS  `stm:""`
		Pending int64 `stm:"pending_capacity_in_Kb"`
	}
)

func (rw *readWrite) set(r, w float64) {
	rw.Read = r
	rw.Write = w
	rw.ReadWrite = r + w
}

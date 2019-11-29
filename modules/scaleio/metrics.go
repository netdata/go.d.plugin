package scaleio

type metrics struct {
	SystemOverview struct {
		Capacity   systemCapacity   `stm:"capacity"`
		IOWorkload systemIOWorkload `stm:""`
		Rebalance  systemRebalance  `stm:"rebalance"`
		Rebuild    systemRebuild    `stm:"rebuild"`
		Components systemComponents `stm:"num_of"`
	} `stm:"system"`
	Sdc         map[string]sdcStatistics         `stm:"sdc"`
	StoragePool map[string]storagePoolStatistics `stm:"storage_pool"`
}

type (
	systemComponents struct {
		Devices            float64 `stm:"devices"`
		FaultSets          float64 `stm:"fault_sets"`
		ProtectionDomains  float64 `stm:"protection_domains"`
		RfcacheDevices     float64 `stm:"rfcache_devices"`
		ScsiInitiators     float64 `stm:"scsi_initiators"`
		Sdc                float64 `stm:"sdc"`
		Sds                float64 `stm:"sds"`
		Snapshots          float64 `stm:"snapshots"`
		StoragePools       float64 `stm:"storage_pools"`
		MappedToAllVolumes float64 `stm:"mapped_to_all_volumes"`
		ThickBaseVolumes   float64 `stm:"thick_base_volumes"`
		ThinBaseVolumes    float64 `stm:"thin_base_volumes"`
		UnmappedVolumes    float64 `stm:"unmapped_volumes"`
		MappedVolumes      float64 `stm:"mapped_volumes"`
		Volumes            float64 `stm:"volumes"`
		VolumesInDeletion  float64 `stm:"volumes_in_deletion"`
		VTrees             float64 `stm:"vtrees"`
	}

	systemCapacity struct {
		MaxCapacity                  float64 `stm:"max_capacity"`
		Protected                    float64 `stm:"protected"`
		InMaintenance                float64 `stm:"in_maintenance"`
		Degraded                     float64 `stm:"degraded"`
		Failed                       float64 `stm:"failed"`
		Spare                        float64 `stm:"spare"`
		UnreachableUnused            float64 `stm:"unreachable_unused"`
		Unused                       float64 `stm:"unused"`
		Decreased                    float64 `stm:"decreased"` // not in statistics, should be calculated
		AvailableForVolumeAllocation float64 `stm:"available_for_volume_allocation"`
		InUse                        float64 `stm:"in_use"`
		Limit                        float64 `stm:"limit"`
		SemiProtected                float64 `stm:"semi_protected"`
		SnapInUse                    float64 `stm:"snap_in_use"`
		SnapInUseOccupied            float64 `stm:"snap_in_use_occupied"`
		ThickInUse                   float64 `stm:"thick_in_use"`
		ThinAllocated                float64 `stm:"thin_allocated"`
		ThinInUse                    float64 `stm:"thin_in_use"`
		ThinFree                     float64 `stm:"thin_free"`
	}
	systemIOWorkload struct {
		Total   bwIOPS `stm:"total"`
		Backend struct {
			Total     bwIOPS `stm:"total"`
			Primary   bwIOPS `stm:"primary"`
			Secondary bwIOPS `stm:"secondary"`
		} `stm:"backend"`
		Frontend bwIOPS `stm:"frontend_user_data"`
	}
	systemRebalance bwIOPSPending
	systemRebuild   struct {
		Total    bwIOPSPending `stm:"total"`
		Forward  bwIOPSPending `stm:"forward"`
		Backward bwIOPSPending `stm:"backward"`
		Normal   bwIOPSPending `stm:"normal"`
	}
)

type (
	sdcStatistics struct {
		bwIOPS             `stm:""`
		MappedVolumes      float64 `stm:"num_of_mapped_volumes"`
		MDMConnectionState bool    `stm:"mdm_connection_state"`
	}
)

type (
	storagePoolStatistics struct {
		Capacity   storagePoolCapacity `stm:"capacity"`
		Components struct {
			Devices   float64 `stm:"devices"`
			Volumes   float64 `stm:"volumes"`
			Vtrees    float64 `stm:"vtrees"`
			Snapshots float64 `stm:"snapshots"`
		} `stm:"num_of"`
	}
	storagePoolCapacity struct {
		MaxCapacity                  float64 `stm:"max_capacity"`
		Protected                    float64 `stm:"protected"`
		InMaintenance                float64 `stm:"in_maintenance"`
		Degraded                     float64 `stm:"degraded"`
		Failed                       float64 `stm:"failed"`
		Spare                        float64 `stm:"spare"`
		UnreachableUnused            float64 `stm:"unreachable_unused"`
		Unused                       float64 `stm:"unused"`
		Decreased                    float64 `stm:"decreased"` // not in statistics, should be calculated
		AvailableForVolumeAllocation float64 `stm:"available_for_volume_allocation"`
	}
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
		Pending float64 `stm:"pending_capacity_in_Kb"`
	}
)

func (rw *readWrite) set(r, w float64) {
	rw.Read = r
	rw.Write = w
	rw.ReadWrite = r + w
}

package scaleio

import mtx "github.com/netdata/go.d.plugin/pkg/metrics"

type metrics struct {
	SystemOverview struct {
		Capacity   systemCapacity   `stm:"capacity"`
		IOWorkload systemIOWorkload `stm:""`
		Rebalance  systemRebalance  `stm:"rebalance"`
		Rebuild    systemRebuild    `stm:"rebuild"`
		Components systemComponents `stm:"num_of"`
	} `stm:"system"`
}

type systemComponents struct {
	Devices            mtx.Gauge `stm:"devices"`
	FaultSets          mtx.Gauge `stm:"fault_sets"`
	ProtectionDomains  mtx.Gauge `stm:"protection_domains"`
	RfcacheDevices     mtx.Gauge `stm:"rfcache_devices"`
	ScsiInitiators     mtx.Gauge `stm:"scsi_initiators"`
	Sdc                mtx.Gauge `stm:"sdc"`
	Sds                mtx.Gauge `stm:"sds"`
	Snapshots          mtx.Gauge `stm:"snapshots"`
	StoragePools       mtx.Gauge `stm:"storage_pools"`
	MappedToAllVolumes mtx.Gauge `stm:"mapped_to_all_volumes"`
	ThickBaseVolumes   mtx.Gauge `stm:"thick_base_volumes"`
	ThinBaseVolumes    mtx.Gauge `stm:"thin_base_volumes"`
	UnmappedVolumes    mtx.Gauge `stm:"unmapped_volumes"`
	MappedVolumes      mtx.Gauge `stm:"mapped_volumes"`
	Volumes            mtx.Gauge `stm:"volumes"`
	VolumesInDeletion  mtx.Gauge `stm:"volumes_in_deletion"`
	Vtrees             mtx.Gauge `stm:"vtrees"`
}

type systemCapacity struct {
	MaxCapacity                  mtx.Gauge `stm:"max_capacity"`
	Protected                    mtx.Gauge `stm:"protected"`
	InMaintenance                mtx.Gauge `stm:"in_maintenance"`
	Degraded                     mtx.Gauge `stm:"degraded"`
	Failed                       mtx.Gauge `stm:"failed"`
	Spare                        mtx.Gauge `stm:"spare"`
	UnreachableUnused            mtx.Gauge `stm:"unreachable_unused"`
	Unused                       mtx.Gauge `stm:"unused"`
	Decreased                    mtx.Gauge `stm:"decreased"` // not in statistics, should be calculated
	AvailableForVolumeAllocation mtx.Gauge `stm:"available_for_volume_allocation"`
	InUse                        mtx.Gauge `stm:"in_use"`
	Limit                        mtx.Gauge `stm:"limit"`
	SemiProtected                mtx.Gauge `stm:"semi_protected"`
	SnapInUse                    mtx.Gauge `stm:"snap_in_use"`
	SnapInUseOccupied            mtx.Gauge `stm:"snap_in_use_occupied"`
	ThickInUse                   mtx.Gauge `stm:"thick_in_use"`
	ThinAllocated                mtx.Gauge `stm:"thin_allocated"`
	ThinInUse                    mtx.Gauge `stm:"thin_in_use"`
	ThinFree                     mtx.Gauge `stm:"thin_free"`
}

type readWrite struct {
	Read      mtx.Gauge `stm:"read,1000,1"`
	Write     mtx.Gauge `stm:"write,1000,1"`
	ReadWrite mtx.Gauge `stm:"read_write,1000,1"`
}

func (rw *readWrite) set(r, w float64) {
	rw.Read.Set(r)
	rw.Write.Set(w)
	rw.ReadWrite.Set(r + w)
}

type bwIOPS struct {
	BW     readWrite `stm:"bandwidth"`
	IOPS   readWrite `stm:"iops"`
	IOSize readWrite `stm:"io_size"`
}

type bwIOPSPending struct {
	bwIOPS  `stm:""`
	Pending mtx.Gauge `stm:"pending_capacity_in_Kb"`
}

type systemIOWorkload struct {
	Total   bwIOPS `stm:"total"`
	Backend struct {
		Total     bwIOPS `stm:"total"`
		Primary   bwIOPS `stm:"primary"`
		Secondary bwIOPS `stm:"secondary"`
	} `stm:"backend"`
	Frontend bwIOPS `stm:"frontend_user_data"`
}

type systemRebalance bwIOPSPending

type systemRebuild struct {
	Total    bwIOPSPending `stm:"total"`
	Forward  bwIOPSPending `stm:"forward"`
	Backward bwIOPSPending `stm:"backward"`
	Normal   bwIOPSPending `stm:"normal"`
}

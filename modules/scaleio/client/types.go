package client

// https://github.com/dell/goscaleio/blob/master/types/v1/types.go

// For all 4xx and 5xx return codes, the body may contain an apiError instance
// with more specifics about the failure.
type apiError struct {
	Message        string
	HTTPStatusCode int
	ErrorCode      int
}

func (e apiError) Error() string {
	return e.Message
}

type Version struct {
	Major int64
	Minor int64
}

type Bwc struct {
	NumOccured      int64
	NumSeconds      int64
	TotalWeightInKb int64
}

type Link struct {
	Rel  string `json:"rel"`
	HREF string `json:"href"`
}

type Sdc struct {
	SystemID           string `json:"systemId"`
	SdcApproved        bool   `json:"sdcApproved"`
	SdcIp              string `json:"SdcIp"`
	SdcGuid            string `json:"sdcGuid"`
	MdmConnectionState string `json:"mdmConnectionState"`
	Name               string `json:"name"`
	ID                 string `json:"id"`
	Links              []Link `json:"links"`
}

type StoragePool struct {
	ProtectionDomainID string `json:"protectionDomainId"`
	Name               string `json:"name"`
	ID                 string `json:"id"`
	Links              []Link `json:"links"`
}

type Instances struct {
	StoragePool []StoragePool `json:"storagePoolList"`
	Sdc         []Sdc         `json:"sdcList"`
}

type (
	allIds bool

	SelectedStatisticsQuery struct {
		List []SelectedObject `json:"selectedStatisticsList"`
	}
	SelectedObject struct {
		Type       string   `json:"type"`
		IDs        []string `json:"ids,omitempty"`
		ALLIDs     allIds   `json:"allIds,omitempty"`
		Properties []string `json:"properties"`
	}
)

func (b allIds) MarshalJSON() ([]byte, error) {
	if b {
		return []byte("[]"), nil
	}
	return nil, nil
}

type SelectedStatistics struct {
	System      SystemStatistics
	Sdc         map[string]SdcStatistics
	StoragePool map[string]StoragePoolStatistics
}

type (
	SystemStatistics struct {
		CapacityAvailableForVolumeAllocationInKb float64
		MaxCapacityInKb                          float64
		CapacityLimitInKb                        float64
		ProtectedCapacityInKb                    float64
		DegradedFailedCapacityInKb               float64
		DegradedHealthyCapacityInKb              float64
		SpareCapacityInKb                        float64
		FailedCapacityInKb                       float64
		UnreachableUnusedCapacityInKb            float64
		InMaintenanceCapacityInKb                float64
		ThinCapacityAllocatedInKb                float64
		ThinCapacityInUseInKb                    float64
		ThickCapacityInUseInKb                   float64

		NumOfDevices            float64
		NumOfFaultSets          float64
		NumOfProtectionDomains  float64
		NumOfRfcacheDevices     float64
		NumOfSdc                float64
		NumOfSds                float64
		NumOfSnapshots          float64
		NumOfStoragePools       float64
		NumOfVolumes            float64
		NumOfVtrees             float64
		NumOfThickBaseVolumes   float64
		NumOfThinBaseVolumes    float64
		NumOfMappedToAllVolumes float64
		NumOfUnmappedVolumes    float64

		RebalanceReadBwc             Bwc
		RebalanceWriteBwc            Bwc
		PendingRebalanceCapacityInKb float64

		PendingNormRebuildCapacityInKb float64
		PendingBckRebuildCapacityInKb  float64
		PendingFwdRebuildCapacityInKb  float64
		NormRebuildReadBwc             Bwc // TODO: ???
		NormRebuildWriteBwc            Bwc // TODO: ???
		BckRebuildReadBwc              Bwc // failed node/disk is back alive
		BckRebuildWriteBwc             Bwc // failed node/disk is back alive
		FwdRebuildReadBwc              Bwc // node/disk fails
		FwdRebuildWriteBwc             Bwc // node/disk fails

		PrimaryReadBwc    Bwc // Backend (SDSs + Devices) Primary - Mater MDM
		PrimaryWriteBwc   Bwc // Backend (SDSs + Devices) Primary - Mater MDM
		SecondaryReadBwc  Bwc // Backend (SDSs + Devices, 2nd) Secondary - Slave MDM
		SecondaryWriteBwc Bwc // Backend (SDSs + Devices, 2nd) Secondary - Slave MDM
		UserDataReadBwc   Bwc // Frontend (Volumes + SDCs)
		UserDataWriteBwc  Bwc // Frontend (Volumes + SDCs)
		TotalReadBwc      Bwc // *ReadBwc
		TotalWriteBwc     Bwc // *WriteBwc

		BackgroundScanCompareCount     float64
		BackgroundScannedInMB          float64
		ActiveBckRebuildCapacityInKb   float64
		ActiveFwdRebuildCapacityInKb   float64
		ActiveMovingCapacityInKb       float64
		ActiveMovingInBckRebuildJobs   float64
		ActiveMovingInFwdRebuildJobs   float64
		ActiveMovingInNormRebuildJobs  float64
		ActiveMovingInRebalanceJobs    float64
		ActiveMovingOutBckRebuildJobs  float64
		ActiveMovingOutFwdRebuildJobs  float64
		ActiveMovingOutNormRebuildJobs float64
		ActiveMovingRebalanceJobs      float64
		ActiveNormRebuildCapacityInKb  float64
		ActiveRebalanceCapacityInKb    float64
		AtRestCapacityInKb             float64
		BckRebuildCapacityInKb         float64
		CapacityInUseInKb              float64
		DegradedFailedVacInKb          float64
		DegradedHealthyVacInKb         float64
		FailedVacInKb                  float64
		FixedReadErrorCount            float64
		FwdRebuildCapacityInKb         float64
		InMaintenanceVacInKb           float64
		InUseVacInKb                   float64
		MovingCapacityInKb             float64
		NormRebuildCapacityInKb        float64

		//NumOfScsiInitiators                             float64  // removed from version 3 of ScaleIO/VxFlex API
		//PendingMovingCapacityInKb                       float64
		//PendingMovingInBckRebuildJobs                   float64
		//PendingMovingInFwdRebuildJobs                   float64
		//PendingMovingInNormRebuildJobs                  float64
		//PendingMovingInRebalanceJobs                    float64
		//PendingMovingOutBckRebuildJobs                  float64
		//PendingMovingOutFwdRebuildJobs                  float64
		//PendingMovingOutNormrebuildJobs                 float64
		//PendingMovingRebalanceJobs                      float64
		//PrimaryReadFromDevBwc                           float64
		//PrimaryReadFromRmcacheBwc                       float64
		//PrimaryVacInKb                                  float64
		//ProtectedVacInKb                                float64
		//ProtectionDomainIds                             float64
		//RebalanceCapacityInKb                           float64
		//RebalancePerReceiveJobNetThrottlingInKbps       float64
		//RebalanceWaitSendQLength                        float64
		//RebuildPerReceiveJobNetThrottlingInKbps         float64
		//RebuildWaitSendQLength                          float64
		//RfacheReadHit                                   float64
		//RfacheWriteHit                                  float64
		//RfcacheAvgReadTime                              float64
		//RfcacheAvgWriteTime                             float64
		//RfcacheFdAvgReadTime                            float64
		//RfcacheFdAvgWriteTime                           float64
		//RfcacheFdCacheOverloaded                        float64
		//RfcacheFdInlightReads                           float64
		//RfcacheFdInlightWrites                          float64
		//RfcacheFdIoErrors                               float64
		//RfcacheFdMonitorErrorStuckIo                    float64
		//RfcacheFdReadTimeGreater1Min                    float64
		//RfcacheFdReadTimeGreater1Sec                    float64
		//RfcacheFdReadTimeGreater500Millis               float64
		//RfcacheFdReadTimeGreater5Sec                    float64
		//RfcacheFdReadsReceived                          float64
		//RfcacheFdWriteTimeGreater1Min                   float64
		//RfcacheFdWriteTimeGreater1Sec                   float64
		//RfcacheFdWriteTimeGreater500Millis              float64
		//RfcacheFdWriteTimeGreater5Sec                   float64
		//RfcacheFdWritesReceived                         float64
		//RfcacheIoErrors                                 float64
		//RfcacheIosOutstanding                           float64
		//RfcacheIosSkipped                               float64
		//RfcachePooIosOutstanding                        float64
		//RfcachePoolCachePages                           float64
		//RfcachePoolEvictions                            float64
		//RfcachePoolInLowMemoryCondition                 float64
		//RfcachePoolIoTimeGreater1Min                    float64
		//RfcachePoolLockTimeGreater1Sec                  float64
		//RfcachePoolLowResourcesInitiatedPassthroughMode float64
		//RfcachePoolNumCacheDevs                         float64
		//RfcachePoolNumSrcDevs                           float64
		//RfcachePoolPagesInuse                           float64
		//RfcachePoolReadHit                              float64
		//RfcachePoolReadMiss                             float64
		//RfcachePoolReadPendingG10Millis                 float64
		//RfcachePoolReadPendingG1Millis                  float64
		//RfcachePoolReadPendingG1Sec                     float64
		//RfcachePoolReadPendingG500Micro                 float64
		//RfcachePoolReadsPending                         float64
		//RfcachePoolSize                                 float64
		//RfcachePoolSourceIdMismatch                     float64
		//RfcachePoolSuspendedIos                         float64
		//RfcachePoolSuspendedPequestsRedundantSearchs    float64
		//RfcachePoolWriteHit                             float64
		//RfcachePoolWriteMiss                            float64
		//RfcachePoolWritePending                         float64
		//RfcachePoolWritePendingG10Millis                float64
		//RfcachePoolWritePendingG1Millis                 float64
		//RfcachePoolWritePendingG1Sec                    float64
		//RfcachePoolWritePendingG500Micro                float64
		//RfcacheReadMiss                                 float64
		//RfcacheReadsFromCache                           float64
		//RfcacheReadsPending                             float64
		//RfcacheReadsReceived                            float64
		//RfcacheReadsSkipped                             float64
		//RfcacheReadsSkippedAlignedSizeTooLarge          float64
		//RfcacheReadsSkippedHeavyLoad                    float64
		//RfcacheReadsSkippedInternalError                float64
		//RfcacheReadsSkippedLockIos                      float64
		//RfcacheReadsSkippedLowResources                 float64
		//RfcacheReadsSkippedMaxIoSize                    float64
		//RfcacheReadsSkippedStuckIo                      float64
		//RfcacheSkippedUnlinedWrite                      float64
		//RfcacheSourceDeviceReads                        float64
		//RfcacheSourceDeviceWrites                       float64
		//RfcacheWriteMiss                                float64
		//RfcacheWritePending                             float64
		//RfcacheWritesReceived                           float64
		//RfcacheWritesSkippedCacheMiss                   float64
		//RfcacheWritesSkippedHeavyLoad                   float64
		//RfcacheWritesSkippedInternalError               float64
		//RfcacheWritesSkippedLowResources                float64
		//RfcacheWritesSkippedMaxIoSize                   float64
		//RfcacheWritesSkippedStuckIo                     float64
		//RmPendingAllocatedInKb                          float64
		//Rmcache128kbEntryCount                          float64
		//Rmcache16kbEntryCount                           float64
		//Rmcache32kbEntryCount                           float64
		//Rmcache4kbEntryCount                            float64
		//Rmcache64kbEntryCount                           float64
		//Rmcache8kbEntryCount                            float64
		//RmcacheBigBlockEvictionCount                    float64
		//RmcacheBigBlockEvictionSizeCountInKb            float64
		//RmcacheCurrNumOf128kbEntries                    float64
		//RmcacheCurrNumOf16kbEntries                     float64
		//RmcacheCurrNumOf32kbEntries                     float64
		//RmcacheCurrNumOf4kbEntries                      float64
		//RmcacheCurrNumOf64kbEntries                     float64
		//RmcacheCurrNumOf8kbEntries                      float64
		//RmcacheEntryEvictionCount                       float64
		//RmcacheEntryEvictionSizeCountInKb               float64
		//RmcacheNoEvictionCount                          float64
		//RmcacheSizeInKb                                 float64
		//RmcacheSizeInUseInKb                            float64
		//RmcacheSkipCountCacheAllBusy                    float64
		//RmcacheSkipCountLargeIo                         float64
		//RmcacheSkipCountUnaligned4kbIo                  float64
		//ScsiInitiatorIds                                float64
		//SdcIds                                          float64
		//SecondaryReadFromDevBwc                         float64
		//SecondaryReadFromRmcacheBwc                     float64
		//SecondaryVacInKb                                float64
		//SemiProtectedCapacityInKb                       float64
		//SemiProtectedVacInKb                            float64
		//SnapCapacityInUseInKb                           float64
		//SnapCapacityInUseOccupiedInKb                   float64
		//UnusedCapacityInKb                              float64
	}

	SdcStatistics struct {
		NumOfMappedVolumes float64
		UserDataReadBwc    Bwc
		UserDataWriteBwc   Bwc
		//VolumeIds          float64
	}

	StoragePoolStatistics struct {
		CapacityAvailableForVolumeAllocationInKb float64
		MaxCapacityInKb                          float64
		CapacityLimitInKb                        float64
		DegradedFailedCapacityInKb               float64
		DegradedHealthyCapacityInKb              float64
		FailedCapacityInKb                       float64
		InMaintenanceCapacityInKb                float64
		ProtectedCapacityInKb                    float64
		SpareCapacityInKb                        float64
		UnreachableUnusedCapacityInKb            float64

		NumOfDevices   float64
		NumOfVolumes   float64
		NumOfVtrees    float64
		NumOfSnapshots float64

		//BackgroundScanCompareCount             float64
		//BackgroundScannedInMB                  float64
		//ActiveBckRebuildCapacityInKb           float64
		//ActiveFwdRebuildCapacityInKb           float64
		//ActiveMovingCapacityInKb               float64
		//ActiveMovingInBckRebuildJobs           float64
		//ActiveMovingInFwdRebuildJobs           float64
		//ActiveMovingInNormRebuildJobs          float64
		//ActiveMovingInRebalanceJobs            float64
		//ActiveMovingOutBckRebuildJobs          float64
		//ActiveMovingOutFwdRebuildJobs          float64
		//ActiveMovingOutNormRebuildJobs         float64
		//ActiveMovingRebalanceJobs              float64
		//ActiveNormRebuildCapacityInKb          float64
		//ActiveRebalanceCapacityInKb            float64
		//AtRestCapacityInKb                     float64
		//BckRebuildCapacityInKb                 float64
		//BckRebuildReadBwc                      float64
		//BckRebuildWriteBwc                     float64
		//CapacityInUseInKb                      float64
		//DegradedFailedVacInKb                  float64
		//DegradedHealthyVacInKb                 float64
		//DeviceIds                              float64
		//FailedVacInKb                          float64
		//FixedReadErrorCount                    float64
		//FwdRebuildCapacityInKb                 float64
		//FwdRebuildReadBwc                      float64
		//FwdRebuildWriteBwc                     float64
		//InMaintenanceVacInKb                   float64
		//InUseVacInKb                           float64
		//MovingCapacityInKb                     float64
		//NormRebuildCapacityInKb                float64
		//NormRebuildReadBwc                     float64
		//NormRebuildWriteBwc                    float64
		//NumOfMappedToAllVolumes                float64
		//NumOfThickBaseVolumes                  float64
		//NumOfThinBaseVolumes                   float64
		//NumOfUnmappedVolumes                   float64
		//NumOfVolumesInDeletion                 float64
		//PendingBckRebuildCapacityInKb          float64
		//PendingFwdRebuildCapacityInKb          float64
		//PendingMovingCapacityInKb              float64
		//PendingMovingInBckRebuildJobs          float64
		//PendingMovingInFwdRebuildJobs          float64
		//PendingMovingInNormRebuildJobs         float64
		//PendingMovingInRebalanceJobs           float64
		//PendingMovingOutBckRebuildJobs         float64
		//PendingMovingOutFwdRebuildJobs         float64
		//PendingMovingOutNormrebuildJobs        float64
		//PendingMovingRebalanceJobs             float64
		//PendingNormRebuildCapacityInKb         float64
		//PendingRebalanceCapacityInKb           float64
		//PrimaryReadBwc                         float64
		//PrimaryReadFromDevBwc                  float64
		//PrimaryReadFromRmcacheBwc              float64
		//PrimaryVacInKb                         float64
		//PrimaryWriteBwc                        float64
		//ProtectedVacInKb                       float64
		//RebalanceCapacityInKb                  float64
		//RebalanceReadBwc                       float64
		//RebalanceWriteBwc                      float64
		//RfacheReadHit                          float64
		//RfacheWriteHit                         float64
		//RfcacheAvgReadTime                     float64
		//RfcacheAvgWriteTime                    float64
		//RfcacheIoErrors                        float64
		//RfcacheIosOutstanding                  float64
		//RfcacheIosSkipped                      float64
		//RfcacheReadMiss                        float64
		//RfcacheReadsFromCache                  float64
		//RfcacheReadsPending                    float64
		//RfcacheReadsReceived                   float64
		//RfcacheReadsSkipped                    float64
		//RfcacheReadsSkippedAlignedSizeTooLarge float64
		//RfcacheReadsSkippedHeavyLoad           float64
		//RfcacheReadsSkippedInternalError       float64
		//RfcacheReadsSkippedLockIos             float64
		//RfcacheReadsSkippedLowResources        float64
		//RfcacheReadsSkippedMaxIoSize           float64
		//RfcacheReadsSkippedStuckIo             float64
		//RfcacheSkippedUnlinedWrite             float64
		//RfcacheSourceDeviceReads               float64
		//RfcacheSourceDeviceWrites              float64
		//RfcacheWriteMiss                       float64
		//RfcacheWritePending                    float64
		//RfcacheWritesReceived                  float64
		//RfcacheWritesSkippedCacheMiss          float64
		//RfcacheWritesSkippedHeavyLoad          float64
		//RfcacheWritesSkippedInternalError      float64
		//RfcacheWritesSkippedLowResources       float64
		//RfcacheWritesSkippedMaxIoSize          float64
		//RfcacheWritesSkippedStuckIo            float64
		//RmPendingAllocatedInKb                 float64
		//SecondaryReadBwc                       float64
		//SecondaryReadFromDevBwc                float64
		//SecondaryReadFromRmcacheBwc            float64
		//SecondaryVacInKb                       float64
		//SecondaryWriteBwc                      float64
		//SemiProtectedCapacityInKb              float64
		//SemiProtectedVacInKb                   float64
		//SnapCapacityInUseInKb                  float64
		//SnapCapacityInUseOccupiedInKb          float64
		//ThickCapacityInUseInKb                 float64
		//ThinCapacityAllocatedInKb              float64
		//ThinCapacityInUseInKb                  float64
		//TotalReadBwc                           float64
		//TotalWriteBwc                          float64
		//UnusedCapacityInKb                     float64
		//UserDataReadBwc                        float64
		//UserDataWriteBwc                       float64
		//VolumeIds                              float64
		//VtreeIds                               float64
	}
	DeviceStatistic struct {
		//	BackgroundScanCompareCount             float64
		//	BackgroundScannedInMB                  float64
		//	ActiveMovingInBckRebuildJobs           float64
		//	ActiveMovingInFwdRebuildJobs           float64
		//	ActiveMovingInNormRebuildJobs          float64
		//	ActiveMovingInRebalanceJobs            float64
		//	ActiveMovingOutBckRebuildJobs          float64
		//	ActiveMovingOutFwdRebuildJobs          float64
		//	ActiveMovingOutNormRebuildJobs         float64
		//	ActiveMovingRebalanceJobs              float64
		//	AvgReadLatencyInMicrosec               float64
		//	AvgReadSizeInBytes                     float64
		//	AvgWriteLatencyInMicrosec              float64
		//	AvgWriteSizeInBytes                    float64
		//	BckRebuildReadBwc                      float64
		//	BckRebuildWriteBwc                     float64
		//	CapacityInUseInKb                      float64
		//	CapacityLimitInKb                      float64
		//	DegradedFailedVacInKb                  float64
		//	DegradedHealthyVacInKb                 float64
		//	FailedVacInKb                          float64
		//	FixedReadErrorCount                    float64
		//	FwdRebuildReadBwc                      float64
		//	FwdRebuildWriteBwc                     float64
		//	InMaintenanceVacInKb                   float64
		//	InUseVacInKb                           float64
		//	MaxCapacityInKb                        float64
		//	NormRebuildReadBwc                     float64
		//	NormRebuildWriteBwc                    float64
		//	PendingMovingInBckRebuildJobs          float64
		//	PendingMovingInFwdRebuildJobs          float64
		//	PendingMovingInNormRebuildJobs         float64
		//	PendingMovingInRebalanceJobs           float64
		//	PendingMovingOutBckRebuildJobs         float64
		//	PendingMovingOutFwdRebuildJobs         float64
		//	PendingMovingOutNormrebuildJobs        float64
		//	PendingMovingRebalanceJobs             float64
		//	PrimaryReadBwc                         float64
		//	PrimaryReadFromDevBwc                  float64
		//	PrimaryReadFromRmcacheBwc              float64
		//	PrimaryVacInKb                         float64
		//	PrimaryWriteBwc                        float64
		//	ProtectedVacInKb                       float64
		//	RebalanceReadBwc                       float64
		//	RebalanceWriteBwc                      float64
		//	RfacheReadHit                          float64
		//	RfacheWriteHit                         float64
		//	RfcacheAvgReadTime                     float64
		//	RfcacheAvgWriteTime                    float64
		//	RfcacheIoErrors                        float64
		//	RfcacheIosOutstanding                  float64
		//	RfcacheIosSkipped                      float64
		//	RfcacheReadMiss                        float64
		//	RfcacheReadsFromCache                  float64
		//	RfcacheReadsPending                    float64
		//	RfcacheReadsReceived                   float64
		//	RfcacheReadsSkipped                    float64
		//	RfcacheReadsSkippedAlignedSizeTooLarge float64
		//	RfcacheReadsSkippedHeavyLoad           float64
		//	RfcacheReadsSkippedInternalError       float64
		//	RfcacheReadsSkippedLockIos             float64
		//	RfcacheReadsSkippedLowResources        float64
		//	RfcacheReadsSkippedMaxIoSize           float64
		//	RfcacheReadsSkippedStuckIo             float64
		//	RfcacheSkippedUnlinedWrite             float64
		//	RfcacheSourceDeviceReads               float64
		//	RfcacheSourceDeviceWrites              float64
		//	RfcacheWriteMiss                       float64
		//	RfcacheWritePending                    float64
		//	RfcacheWritesReceived                  float64
		//	RfcacheWritesSkippedCacheMiss          float64
		//	RfcacheWritesSkippedHeavyLoad          float64
		//	RfcacheWritesSkippedInternalError      float64
		//	RfcacheWritesSkippedLowResources       float64
		//	RfcacheWritesSkippedMaxIoSize          float64
		//	RfcacheWritesSkippedStuckIo            float64
		//	RmPendingAllocatedInKb                 float64
		//	SecondaryReadBwc                       float64
		//	SecondaryReadFromDevBwc                float64
		//	SecondaryReadFromRmcacheBwc            float64
		//	SecondaryVacInKb                       float64
		//	SecondaryWriteBwc                      float64
		//	SemiProtectedVacInKb                   float64
		//	SnapCapacityInUseInKb                  float64
		//	SnapCapacityInUseOccupiedInKb          float64
		//	ThickCapacityInUseInKb                 float64
		//	ThinCapacityAllocatedInKb              float64
		//	ThinCapacityInUseInKb                  float64
		//	TotalReadBwc                           float64
		//	TotalWriteBwc                          float64
		//	UnreachableUnusedCapacityInKb          float64
		//	UnusedCapacityInKb                     float64
	}
	FaultSetStatistics struct {
		//	BackgroundScanCompareCount                      float64
		//	BackgroundScannedInMB                           float64
		//	ActiveMovingInBckRebuildJobs                    float64
		//	ActiveMovingInFwdRebuildJobs                    float64
		//	ActiveMovingInNormRebuildJobs                   float64
		//	ActiveMovingInRebalanceJobs                     float64
		//	ActiveMovingOutBckRebuildJobs                   float64
		//	ActiveMovingOutFwdRebuildJobs                   float64
		//	ActiveMovingOutNormRebuildJobs                  float64
		//	ActiveMovingRebalanceJobs                       float64
		//	BckRebuildReadBwc                               float64
		//	BckRebuildWriteBwc                              float64
		//	CapacityInUseInKb                               float64
		//	CapacityLimitInKb                               float64
		//	DegradedFailedVacInKb                           float64
		//	DegradedHealthyVacInKb                          float64
		//	FailedVacInKb                                   float64
		//	FixedReadErrorCount                             float64
		//	FwdRebuildReadBwc                               float64
		//	FwdRebuildWriteBwc                              float64
		//	InMaintenanceVacInKb                            float64
		//	InUseVacInKb                                    float64
		//	MaxCapacityInKb                                 float64
		//	NormRebuildReadBwc                              float64
		//	NormRebuildWriteBwc                             float64
		//	NumOfSds                                        float64
		//	PendingMovingInBckRebuildJobs                   float64
		//	PendingMovingInFwdRebuildJobs                   float64
		//	PendingMovingInNormRebuildJobs                  float64
		//	PendingMovingInRebalanceJobs                    float64
		//	PendingMovingOutBckRebuildJobs                  float64
		//	PendingMovingOutFwdRebuildJobs                  float64
		//	PendingMovingOutNormrebuildJobs                 float64
		//	PendingMovingRebalanceJobs                      float64
		//	PrimaryReadBwc                                  float64
		//	PrimaryReadFromDevBwc                           float64
		//	PrimaryReadFromRmcacheBwc                       float64
		//	PrimaryVacInKb                                  float64
		//	PrimaryWriteBwc                                 float64
		//	ProtectedVacInKb                                float64
		//	RebalancePerReceiveJobNetThrottlingInKbps       float64
		//	RebalanceReadBwc                                float64
		//	RebalanceWaitSendQLength                        float64
		//	RebalanceWriteBwc                               float64
		//	RebuildPerReceiveJobNetThrottlingInKbps         float64
		//	RebuildWaitSendQLength                          float64
		//	RfacheReadHit                                   float64
		//	RfacheWriteHit                                  float64
		//	RfcacheAvgReadTime                              float64
		//	RfcacheAvgWriteTime                             float64
		//	RfcacheFdAvgReadTime                            float64
		//	RfcacheFdAvgWriteTime                           float64
		//	RfcacheFdCacheOverloaded                        float64
		//	RfcacheFdInlightReads                           float64
		//	RfcacheFdInlightWrites                          float64
		//	RfcacheFdIoErrors                               float64
		//	RfcacheFdMonitorErrorStuckIo                    float64
		//	RfcacheFdReadTimeGreater1Min                    float64
		//	RfcacheFdReadTimeGreater1Sec                    float64
		//	RfcacheFdReadTimeGreater500Millis               float64
		//	RfcacheFdReadTimeGreater5Sec                    float64
		//	RfcacheFdReadsReceived                          float64
		//	RfcacheFdWriteTimeGreater1Min                   float64
		//	RfcacheFdWriteTimeGreater1Sec                   float64
		//	RfcacheFdWriteTimeGreater500Millis              float64
		//	RfcacheFdWriteTimeGreater5Sec                   float64
		//	RfcacheFdWritesReceived                         float64
		//	RfcacheIoErrors                                 float64
		//	RfcacheIosOutstanding                           float64
		//	RfcacheIosSkipped                               float64
		//	RfcachePooIosOutstanding                        float64
		//	RfcachePoolCachePages                           float64
		//	RfcachePoolEvictions                            float64
		//	RfcachePoolInLowMemoryCondition                 float64
		//	RfcachePoolIoTimeGreater1Min                    float64
		//	RfcachePoolLockTimeGreater1Sec                  float64
		//	RfcachePoolLowResourcesInitiatedPassthroughMode float64
		//	RfcachePoolNumCacheDevs                         float64
		//	RfcachePoolNumSrcDevs                           float64
		//	RfcachePoolPagesInuse                           float64
		//	RfcachePoolReadHit                              float64
		//	RfcachePoolReadMiss                             float64
		//	RfcachePoolReadPendingG10Millis                 float64
		//	RfcachePoolReadPendingG1Millis                  float64
		//	RfcachePoolReadPendingG1Sec                     float64
		//	RfcachePoolReadPendingG500Micro                 float64
		//	RfcachePoolReadsPending                         float64
		//	RfcachePoolSize                                 float64
		//	RfcachePoolSourceIdMismatch                     float64
		//	RfcachePoolSuspendedIos                         float64
		//	RfcachePoolSuspendedPequestsRedundantSearchs    float64
		//	RfcachePoolWriteHit                             float64
		//	RfcachePoolWriteMiss                            float64
		//	RfcachePoolWritePending                         float64
		//	RfcachePoolWritePendingG10Millis                float64
		//	RfcachePoolWritePendingG1Millis                 float64
		//	RfcachePoolWritePendingG1Sec                    float64
		//	RfcachePoolWritePendingG500Micro                float64
		//	RfcacheReadMiss                                 float64
		//	RfcacheReadsFromCache                           float64
		//	RfcacheReadsPending                             float64
		//	RfcacheReadsReceived                            float64
		//	RfcacheReadsSkipped                             float64
		//	RfcacheReadsSkippedAlignedSizeTooLarge          float64
		//	RfcacheReadsSkippedHeavyLoad                    float64
		//	RfcacheReadsSkippedInternalError                float64
		//	RfcacheReadsSkippedLockIos                      float64
		//	RfcacheReadsSkippedLowResources                 float64
		//	RfcacheReadsSkippedMaxIoSize                    float64
		//	RfcacheReadsSkippedStuckIo                      float64
		//	RfcacheSkippedUnlinedWrite                      float64
		//	RfcacheSourceDeviceReads                        float64
		//	RfcacheSourceDeviceWrites                       float64
		//	RfcacheWriteMiss                                float64
		//	RfcacheWritePending                             float64
		//	RfcacheWritesReceived                           float64
		//	RfcacheWritesSkippedCacheMiss                   float64
		//	RfcacheWritesSkippedHeavyLoad                   float64
		//	RfcacheWritesSkippedInternalError               float64
		//	RfcacheWritesSkippedLowResources                float64
		//	RfcacheWritesSkippedMaxIoSize                   float64
		//	RfcacheWritesSkippedStuckIo                     float64
		//	RmPendingAllocatedInKb                          float64
		//	Rmcache128kbEntryCount                          float64
		//	Rmcache16kbEntryCount                           float64
		//	Rmcache32kbEntryCount                           float64
		//	Rmcache4kbEntryCount                            float64
		//	Rmcache64kbEntryCount                           float64
		//	Rmcache8kbEntryCount                            float64
		//	RmcacheBigBlockEvictionCount                    float64
		//	RmcacheBigBlockEvictionSizeCountInKb            float64
		//	RmcacheCurrNumOf128kbEntries                    float64
		//	RmcacheCurrNumOf16kbEntries                     float64
		//	RmcacheCurrNumOf32kbEntries                     float64
		//	RmcacheCurrNumOf4kbEntries                      float64
		//	RmcacheCurrNumOf64kbEntries                     float64
		//	RmcacheCurrNumOf8kbEntries                      float64
		//	RmcacheEntryEvictionCount                       float64
		//	RmcacheEntryEvictionSizeCountInKb               float64
		//	RmcacheNoEvictionCount                          float64
		//	RmcacheSizeInKb                                 float64
		//	RmcacheSizeInUseInKb                            float64
		//	RmcacheSkipCountCacheAllBusy                    float64
		//	RmcacheSkipCountLargeIo                         float64
		//	RmcacheSkipCountUnaligned4kbIo                  float64
		//	SdsIds                                          float64
		//	SecondaryReadBwc                                float64
		//	SecondaryReadFromDevBwc                         float64
		//	SecondaryReadFromRmcacheBwc                     float64
		//	SecondaryVacInKb                                float64
		//	SecondaryWriteBwc                               float64
		//	SemiProtectedVacInKb                            float64
		//	SnapCapacityInUseInKb                           float64
		//	SnapCapacityInUseOccupiedInKb                   float64
		//	ThickCapacityInUseInKb                          float64
		//	ThinCapacityAllocatedInKb                       float64
		//	ThinCapacityInUseInKb                           float64
		//	TotalReadBwc                                    float64
		//	TotalWriteBwc                                   float64
		//	UnreachableUnusedCapacityInKb                   float64
		//	UnusedCapacityInKb                              float64
	}
	ProtectionDomainStatistics struct {
		//	BackgroundScanCompareCount                      float64
		//	BackgroundScannedInMB                           float64
		//	ActiveBckRebuildCapacityInKb                    float64
		//	ActiveFwdRebuildCapacityInKb                    float64
		//	ActiveMovingCapacityInKb                        float64
		//	ActiveMovingInBckRebuildJobs                    float64
		//	ActiveMovingInFwdRebuildJobs                    float64
		//	ActiveMovingInNormRebuildJobs                   float64
		//	ActiveMovingInRebalanceJobs                     float64
		//	ActiveMovingOutBckRebuildJobs                   float64
		//	ActiveMovingOutFwdRebuildJobs                   float64
		//	ActiveMovingOutNormRebuildJobs                  float64
		//	ActiveMovingRebalanceJobs                       float64
		//	ActiveNormRebuildCapacityInKb                   float64
		//	ActiveRebalanceCapacityInKb                     float64
		//	AtRestCapacityInKb                              float64
		//	BckRebuildCapacityInKb                          float64
		//	BckRebuildReadBwc                               float64
		//	BckRebuildWriteBwc                              float64
		//	CapacityAvailableForVolumeAllocationInKb        float64
		//	CapacityInUseInKb                               float64
		//	CapacityLimitInKb                               float64
		//	DegradedFailedCapacityInKb                      float64
		//	DegradedFailedVacInKb                           float64
		//	DegradedHealthyCapacityInKb                     float64
		//	DegradedHealthyVacInKb                          float64
		//	FailedCapacityInKb                              float64
		//	FailedVacInKb                                   float64
		//	FaultSetIds                                     float64
		//	FixedReadErrorCount                             float64
		//	FwdRebuildCapacityInKb                          float64
		//	FwdRebuildReadBwc                               float64
		//	FwdRebuildWriteBwc                              float64
		//	InMaintenanceCapacityInKb                       float64
		//	InMaintenanceVacInKb                            float64
		//	InUseVacInKb                                    float64
		//	MaxCapacityInKb                                 float64
		//	MovingCapacityInKb                              float64
		//	NormRebuildCapacityInKb                         float64
		//	NormRebuildReadBwc                              float64
		//	NormRebuildWriteBwc                             float64
		//	NumOfFaultSets                                  float64
		//	NumOfMappedToAllVolumes                         float64
		//	NumOfSds                                        float64
		//	NumOfSnapshots                                  float64
		//	NumOfStoragePools                               float64
		//	NumOfThickBaseVolumes                           float64
		//	NumOfThinBaseVolumes                            float64
		//	NumOfUnmappedVolumes                            float64
		//	NumOfVolumesInDeletion                          float64
		//	PendingBckRebuildCapacityInKb                   float64
		//	PendingFwdRebuildCapacityInKb                   float64
		//	PendingMovingCapacityInKb                       float64
		//	PendingMovingInBckRebuildJobs                   float64
		//	PendingMovingInFwdRebuildJobs                   float64
		//	PendingMovingInNormRebuildJobs                  float64
		//	PendingMovingInRebalanceJobs                    float64
		//	PendingMovingOutBckRebuildJobs                  float64
		//	PendingMovingOutFwdRebuildJobs                  float64
		//	PendingMovingOutNormrebuildJobs                 float64
		//	PendingMovingRebalanceJobs                      float64
		//	PendingNormRebuildCapacityInKb                  float64
		//	PendingRebalanceCapacityInKb                    float64
		//	PrimaryReadBwc                                  float64
		//	PrimaryReadFromDevBwc                           float64
		//	PrimaryReadFromRmcacheBwc                       float64
		//	PrimaryVacInKb                                  float64
		//	PrimaryWriteBwc                                 float64
		//	ProtectedCapacityInKb                           float64
		//	ProtectedVacInKb                                float64
		//	RebalanceCapacityInKb                           float64
		//	RebalancePerReceiveJobNetThrottlingInKbps       float64
		//	RebalanceReadBwc                                float64
		//	RebalanceWaitSendQLength                        float64
		//	RebalanceWriteBwc                               float64
		//	RebuildPerReceiveJobNetThrottlingInKbps         float64
		//	RebuildWaitSendQLength                          float64
		//	RfacheReadHit                                   float64
		//	RfacheWriteHit                                  float64
		//	RfcacheAvgReadTime                              float64
		//	RfcacheAvgWriteTime                             float64
		//	RfcacheFdAvgReadTime                            float64
		//	RfcacheFdAvgWriteTime                           float64
		//	RfcacheFdCacheOverloaded                        float64
		//	RfcacheFdInlightReads                           float64
		//	RfcacheFdInlightWrites                          float64
		//	RfcacheFdIoErrors                               float64
		//	RfcacheFdMonitorErrorStuckIo                    float64
		//	RfcacheFdReadTimeGreater1Min                    float64
		//	RfcacheFdReadTimeGreater1Sec                    float64
		//	RfcacheFdReadTimeGreater500Millis               float64
		//	RfcacheFdReadTimeGreater5Sec                    float64
		//	RfcacheFdReadsReceived                          float64
		//	RfcacheFdWriteTimeGreater1Min                   float64
		//	RfcacheFdWriteTimeGreater1Sec                   float64
		//	RfcacheFdWriteTimeGreater500Millis              float64
		//	RfcacheFdWriteTimeGreater5Sec                   float64
		//	RfcacheFdWritesReceived                         float64
		//	RfcacheIoErrors                                 float64
		//	RfcacheIosOutstanding                           float64
		//	RfcacheIosSkipped                               float64
		//	RfcachePooIosOutstanding                        float64
		//	RfcachePoolCachePages                           float64
		//	RfcachePoolEvictions                            float64
		//	RfcachePoolInLowMemoryCondition                 float64
		//	RfcachePoolIoTimeGreater1Min                    float64
		//	RfcachePoolLockTimeGreater1Sec                  float64
		//	RfcachePoolLowResourcesInitiatedPassthroughMode float64
		//	RfcachePoolNumCacheDevs                         float64
		//	RfcachePoolNumSrcDevs                           float64
		//	RfcachePoolPagesInuse                           float64
		//	RfcachePoolReadHit                              float64
		//	RfcachePoolReadMiss                             float64
		//	RfcachePoolReadPendingG10Millis                 float64
		//	RfcachePoolReadPendingG1Millis                  float64
		//	RfcachePoolReadPendingG1Sec                     float64
		//	RfcachePoolReadPendingG500Micro                 float64
		//	RfcachePoolReadsPending                         float64
		//	RfcachePoolSize                                 float64
		//	RfcachePoolSourceIdMismatch                     float64
		//	RfcachePoolSuspendedIos                         float64
		//	RfcachePoolSuspendedPequestsRedundantSearchs    float64
		//	RfcachePoolWriteHit                             float64
		//	RfcachePoolWriteMiss                            float64
		//	RfcachePoolWritePending                         float64
		//	RfcachePoolWritePendingG10Millis                float64
		//	RfcachePoolWritePendingG1Millis                 float64
		//	RfcachePoolWritePendingG1Sec                    float64
		//	RfcachePoolWritePendingG500Micro                float64
		//	RfcacheReadMiss                                 float64
		//	RfcacheReadsFromCache                           float64
		//	RfcacheReadsPending                             float64
		//	RfcacheReadsReceived                            float64
		//	RfcacheReadsSkipped                             float64
		//	RfcacheReadsSkippedAlignedSizeTooLarge          float64
		//	RfcacheReadsSkippedHeavyLoad                    float64
		//	RfcacheReadsSkippedInternalError                float64
		//	RfcacheReadsSkippedLockIos                      float64
		//	RfcacheReadsSkippedLowResources                 float64
		//	RfcacheReadsSkippedMaxIoSize                    float64
		//	RfcacheReadsSkippedStuckIo                      float64
		//	RfcacheSkippedUnlinedWrite                      float64
		//	RfcacheSourceDeviceReads                        float64
		//	RfcacheSourceDeviceWrites                       float64
		//	RfcacheWriteMiss                                float64
		//	RfcacheWritePending                             float64
		//	RfcacheWritesReceived                           float64
		//	RfcacheWritesSkippedCacheMiss                   float64
		//	RfcacheWritesSkippedHeavyLoad                   float64
		//	RfcacheWritesSkippedInternalError               float64
		//	RfcacheWritesSkippedLowResources                float64
		//	RfcacheWritesSkippedMaxIoSize                   float64
		//	RfcacheWritesSkippedStuckIo                     float64
		//	RmPendingAllocatedInKb                          float64
		//	Rmcache128kbEntryCount                          float64
		//	Rmcache16kbEntryCount                           float64
		//	Rmcache32kbEntryCount                           float64
		//	Rmcache4kbEntryCount                            float64
		//	Rmcache64kbEntryCount                           float64
		//	Rmcache8kbEntryCount                            float64
		//	RmcacheBigBlockEvictionCount                    float64
		//	RmcacheBigBlockEvictionSizeCountInKb            float64
		//	RmcacheCurrNumOf128kbEntries                    float64
		//	RmcacheCurrNumOf16kbEntries                     float64
		//	RmcacheCurrNumOf32kbEntries                     float64
		//	RmcacheCurrNumOf4kbEntries                      float64
		//	RmcacheCurrNumOf64kbEntries                     float64
		//	RmcacheCurrNumOf8kbEntries                      float64
		//	RmcacheEntryEvictionCount                       float64
		//	RmcacheEntryEvictionSizeCountInKb               float64
		//	RmcacheNoEvictionCount                          float64
		//	RmcacheSizeInKb                                 float64
		//	RmcacheSizeInUseInKb                            float64
		//	RmcacheSkipCountCacheAllBusy                    float64
		//	RmcacheSkipCountLargeIo                         float64
		//	RmcacheSkipCountUnaligned4kbIo                  float64
		//	SdsIds                                          float64
		//	SecondaryReadBwc                                float64
		//	SecondaryReadFromDevBwc                         float64
		//	SecondaryReadFromRmcacheBwc                     float64
		//	SecondaryVacInKb                                float64
		//	SecondaryWriteBwc                               float64
		//	SemiProtectedCapacityInKb                       float64
		//	SemiProtectedVacInKb                            float64
		//	SnapCapacityInUseInKb                           float64
		//	SnapCapacityInUseOccupiedInKb                   float64
		//	SpareCapacityInKb                               float64
		//	StoragePoolIds                                  float64
		//	ThickCapacityInUseInKb                          float64
		//	ThinCapacityAllocatedInKb                       float64
		//	ThinCapacityInUseInKb                           float64
		//	TotalReadBwc                                    float64
		//	TotalWriteBwc                                   float64
		//	UnreachableUnusedCapacityInKb                   float64
		//	UnusedCapacityInKb                              float64
		//	UserDataReadBwc                                 float64
		//	UserDataWriteBwc                                float64
	}
	RFCacheDeviceStatistics struct {
		//	RfcacheFdAvgReadTime               float64
		//	RfcacheFdAvgWriteTime              float64
		//	RfcacheFdCacheOverloaded           float64
		//	RfcacheFdInlightReads              float64
		//	RfcacheFdInlightWrites             float64
		//	RfcacheFdIoErrors                  float64
		//	RfcacheFdMonitorErrorStuckIo       float64
		//	RfcacheFdReadTimeGreater1Min       float64
		//	RfcacheFdReadTimeGreater1Sec       float64
		//	RfcacheFdReadTimeGreater500Millis  float64
		//	RfcacheFdReadTimeGreater5Sec       float64
		//	RfcacheFdReadsReceived             float64
		//	RfcacheFdWriteTimeGreater1Min      float64
		//	RfcacheFdWriteTimeGreater1Sec      float64
		//	RfcacheFdWriteTimeGreater500Millis float64
		//	RfcacheFdWriteTimeGreater5Sec      float64
		//	RfcacheFdWritesReceived            float64
	}
	SdsStatistics struct {
		//	BackgroundScanCompareCount                      float64
		//	BackgroundScannedInMB                           float64
		//	ActiveMovingInBckRebuildJobs                    float64
		//	ActiveMovingInFwdRebuildJobs                    float64
		//	ActiveMovingInNormRebuildJobs                   float64
		//	ActiveMovingInRebalanceJobs                     float64
		//	ActiveMovingOutBckRebuildJobs                   float64
		//	ActiveMovingOutFwdRebuildJobs                   float64
		//	ActiveMovingOutNormRebuildJobs                  float64
		//	ActiveMovingRebalanceJobs                       float64
		//	BckRebuildReadBwc                               float64
		//	BckRebuildWriteBwc                              float64
		//	CapacityInUseInKb                               float64
		//	CapacityLimitInKb                               float64
		//	DegradedFailedVacInKb                           float64
		//	DegradedHealthyVacInKb                          float64
		//	DeviceIds                                       float64
		//	FailedVacInKb                                   float64
		//	FixedReadErrorCount                             float64
		//	FwdRebuildReadBwc                               float64
		//	FwdRebuildWriteBwc                              float64
		//	InMaintenanceVacInKb                            float64
		//	InUseVacInKb                                    float64
		//	MaxCapacityInKb                                 float64
		//	NormRebuildReadBwc                              float64
		//	NormRebuildWriteBwc                             float64
		//	NumOfDevices                                    float64
		//	NumOfRfcacheDevices                             float64
		//	PendingMovingInBckRebuildJobs                   float64
		//	PendingMovingInFwdRebuildJobs                   float64
		//	PendingMovingInNormRebuildJobs                  float64
		//	PendingMovingInRebalanceJobs                    float64
		//	PendingMovingOutBckRebuildJobs                  float64
		//	PendingMovingOutFwdRebuildJobs                  float64
		//	PendingMovingOutNormrebuildJobs                 float64
		//	PendingMovingRebalanceJobs                      float64
		//	PrimaryReadBwc                                  float64
		//	PrimaryReadFromDevBwc                           float64
		//	PrimaryReadFromRmcacheBwc                       float64
		//	PrimaryVacInKb                                  float64
		//	PrimaryWriteBwc                                 float64
		//	ProtectedVacInKb                                float64
		//	RebalancePerReceiveJobNetThrottlingInKbps       float64
		//	RebalanceReadBwc                                float64
		//	RebalanceWaitSendQLength                        float64
		//	RebalanceWriteBwc                               float64
		//	RebuildPerReceiveJobNetThrottlingInKbps         float64
		//	RebuildWaitSendQLength                          float64
		//	RfacheReadHit                                   float64
		//	RfacheWriteHit                                  float64
		//	RfcacheAvgReadTime                              float64
		//	RfcacheAvgWriteTime                             float64
		//	RfcacheDeviceIds                                float64
		//	RfcacheFdAvgReadTime                            float64
		//	RfcacheFdAvgWriteTime                           float64
		//	RfcacheFdCacheOverloaded                        float64
		//	RfcacheFdInlightReads                           float64
		//	RfcacheFdInlightWrites                          float64
		//	RfcacheFdIoErrors                               float64
		//	RfcacheFdMonitorErrorStuckIo                    float64
		//	RfcacheFdReadTimeGreater1Min                    float64
		//	RfcacheFdReadTimeGreater1Sec                    float64
		//	RfcacheFdReadTimeGreater500Millis               float64
		//	RfcacheFdReadTimeGreater5Sec                    float64
		//	RfcacheFdReadsReceived                          float64
		//	RfcacheFdWriteTimeGreater1Min                   float64
		//	RfcacheFdWriteTimeGreater1Sec                   float64
		//	RfcacheFdWriteTimeGreater500Millis              float64
		//	RfcacheFdWriteTimeGreater5Sec                   float64
		//	RfcacheFdWritesReceived                         float64
		//	RfcacheIoErrors                                 float64
		//	RfcacheIosOutstanding                           float64
		//	RfcacheIosSkipped                               float64
		//	RfcachePooIosOutstanding                        float64
		//	RfcachePoolCachePages                           float64
		//	RfcachePoolContinuosMem                         float64
		//	RfcachePoolEvictions                            float64
		//	RfcachePoolInLowMemoryCondition                 float64
		//	RfcachePoolIoTimeGreater1Min                    float64
		//	RfcachePoolLockTimeGreater1Sec                  float64
		//	RfcachePoolLowResourcesInitiatedPassthroughMode float64
		//	RfcachePoolMaxIoSize                            float64
		//	RfcachePoolNumCacheDevs                         float64
		//	RfcachePoolNumOfDriverTheads                    float64
		//	RfcachePoolNumSrcDevs                           float64
		//	RfcachePoolOpmode                               float64
		//	RfcachePoolPageSize                             float64
		//	RfcachePoolPagesInuse                           float64
		//	RfcachePoolReadHit                              float64
		//	RfcachePoolReadMiss                             float64
		//	RfcachePoolReadPendingG10Millis                 float64
		//	RfcachePoolReadPendingG1Millis                  float64
		//	RfcachePoolReadPendingG1Sec                     float64
		//	RfcachePoolReadPendingG500Micro                 float64
		//	RfcachePoolReadsPending                         float64
		//	RfcachePoolSize                                 float64
		//	RfcachePoolSourceIdMismatch                     float64
		//	RfcachePoolSuspendedIos                         float64
		//	RfcachePoolSuspendedIosMax                      float64
		//	RfcachePoolSuspendedPequestsRedundantSearchs    float64
		//	RfcachePoolWriteHit                             float64
		//	RfcachePoolWriteMiss                            float64
		//	RfcachePoolWritePending                         float64
		//	RfcachePoolWritePendingG10Millis                float64
		//	RfcachePoolWritePendingG1Millis                 float64
		//	RfcachePoolWritePendingG1Sec                    float64
		//	RfcachePoolWritePendingG500Micro                float64
		//	RfcacheReadMiss                                 float64
		//	RfcacheReadsFromCache                           float64
		//	RfcacheReadsPending                             float64
		//	RfcacheReadsReceived                            float64
		//	RfcacheReadsSkipped                             float64
		//	RfcacheReadsSkippedAlignedSizeTooLarge          float64
		//	RfcacheReadsSkippedHeavyLoad                    float64
		//	RfcacheReadsSkippedInternalError                float64
		//	RfcacheReadsSkippedLockIos                      float64
		//	RfcacheReadsSkippedLowResources                 float64
		//	RfcacheReadsSkippedMaxIoSize                    float64
		//	RfcacheReadsSkippedStuckIo                      float64
		//	RfcacheSkippedUnlinedWrite                      float64
		//	RfcacheSourceDeviceReads                        float64
		//	RfcacheSourceDeviceWrites                       float64
		//	RfcacheWriteMiss                                float64
		//	RfcacheWritePending                             float64
		//	RfcacheWritesReceived                           float64
		//	RfcacheWritesSkippedCacheMiss                   float64
		//	RfcacheWritesSkippedHeavyLoad                   float64
		//	RfcacheWritesSkippedInternalError               float64
		//	RfcacheWritesSkippedLowResources                float64
		//	RfcacheWritesSkippedMaxIoSize                   float64
		//	RfcacheWritesSkippedStuckIo                     float64
		//	RmPendingAllocatedInKb                          float64
		//	Rmcache128kbEntryCount                          float64
		//	Rmcache16kbEntryCount                           float64
		//	Rmcache32kbEntryCount                           float64
		//	Rmcache4kbEntryCount                            float64
		//	Rmcache64kbEntryCount                           float64
		//	Rmcache8kbEntryCount                            float64
		//	RmcacheBigBlockEvictionCount                    float64
		//	RmcacheBigBlockEvictionSizeCountInKb            float64
		//	RmcacheCurrNumOf128kbEntries                    float64
		//	RmcacheCurrNumOf16kbEntries                     float64
		//	RmcacheCurrNumOf32kbEntries                     float64
		//	RmcacheCurrNumOf4kbEntries                      float64
		//	RmcacheCurrNumOf64kbEntries                     float64
		//	RmcacheCurrNumOf8kbEntries                      float64
		//	RmcacheEntryEvictionCount                       float64
		//	RmcacheEntryEvictionSizeCountInKb               float64
		//	RmcacheNoEvictionCount                          float64
		//	RmcacheSizeInKb                                 float64
		//	RmcacheSizeInUseInKb                            float64
		//	RmcacheSkipCountCacheAllBusy                    float64
		//	RmcacheSkipCountLargeIo                         float64
		//	RmcacheSkipCountUnaligned4kbIo                  float64
		//	SecondaryReadBwc                                float64
		//	SecondaryReadFromDevBwc                         float64
		//	SecondaryReadFromRmcacheBwc                     float64
		//	SecondaryVacInKb                                float64
		//	SecondaryWriteBwc                               float64
		//	SemiProtectedVacInKb                            float64
		//	SnapCapacityInUseInKb                           float64
		//	SnapCapacityInUseOccupiedInKb                   float64
		//	ThickCapacityInUseInKb                          float64
		//	ThinCapacityAllocatedInKb                       float64
		//	ThinCapacityInUseInKb                           float64
		//	TotalReadBwc                                    float64
		//	TotalWriteBwc                                   float64
		//	UnreachableUnusedCapacityInKb                   float64
		//	UnusedCapacityInKb                              float64
	}
	VolumeStatistics struct {
		//	ChildVolumeIds            float64
		//	DescendantVolumeIds       float64
		//	MappedSdcIds              float64
		//	NumOfChildVolumes         float64
		//	NumOfDescendantVolumes    float64
		//	NumOfMappedScsiInitiators float64
		//	NumOfMappedSdcs           float64
		//	UserDataReadBwc           float64
		//	UserDataWriteBwc          float64
	}
	VTreeStatistics struct {
		//	BaseNetCapacityInUseInKb float64
		//	NetCapacityInUseInKb     float64
		//	NumOfVolumes             float64
		//	SnapNetCapacityInUseInKb float64
		//	TrimmedCapacityInKb      float64
		//	VolumeIds                float64
	}
)

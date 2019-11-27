package scaleio

import (
	"github.com/netdata/go.d.plugin/modules/scaleio/client"
	"github.com/netdata/go.d.plugin/pkg/stm"
)

/*
Starting from version 3 of ScaleIO/VxFlex API numOfScsiInitiators property is removed from the system selectedStatisticsQuery.
Reference: VxFlex OS v3.x REST API Reference Guide.pdf
*/
var (
	selectedStatisticsQuery = `
{
	"selectedStatisticsList": [{
			"type": "System",
			"properties": [
				"bckRebuildReadBwc",
				"bckRebuildWriteBwc",
				"fwdRebuildReadBwc",
				"fwdRebuildWriteBwc",
				"normRebuildReadBwc",
				"normRebuildWriteBwc",
				"rebalanceReadBwc",
				"rebalanceWriteBwc",
				"primaryReadBwc",
				"primaryWriteBwc",
				"secondaryReadBwc",
				"secondaryWriteBwc",
				"userDataReadBwc",
				"userDataWriteBwc",
				"totalReadBwc",
				"totalWriteBwc",
				"activeBckRebuildCapacityInKb",
				"activeFwdRebuildCapacityInKb",
				"activeMovingCapacityInKb",
				"activeNormRebuildCapacityInKb",
				"activeRebalanceCapacityInKb",
				"atRestCapacityInKb",
				"bckRebuildCapacityInKb",
				"capacityAvailableForVolumeAllocationInKb",
				"capacityInUseInKb",
				"capacityLimitInKb",
				"degradedFailedCapacityInKb",
				"degradedHealthyCapacityInKb",
				"failedCapacityInKb",
				"fwdRebuildCapacityInKb",
				"inMaintenanceCapacityInKb",
				"maxCapacityInKb",
				"movingCapacityInKb",
				"normRebuildCapacityInKb",
				"pendingBckRebuildCapacityInKb",
				"pendingFwdRebuildCapacityInKb",
				"pendingMovingCapacityInKb",
				"pendingNormRebuildCapacityInKb",
				"pendingRebalanceCapacityInKb",
				"protectedCapacityInKb",
				"rebalanceCapacityInKb",
				"semiProtectedCapacityInKb",
				"snapCapacityInUseInKb",
				"snapCapacityInUseOccupiedInKb",
				"spareCapacityInKb",
				"thickCapacityInUseInKb",
				"thinCapacityAllocatedInKb",
				"thinCapacityInUseInKb",
				"unreachableUnusedCapacityInKb",
				"unusedCapacityInKb",
				"numOfDevices",
				"numOfFaultSets",
				"numOfMappedToAllVolumes",
				"numOfProtectionDomains",
				"numOfRfcacheDevices",
				"numOfSdc",
				"numOfSds",
				"numOfSnapshots",
				"numOfStoragePools",
				"numOfThickBaseVolumes",
				"numOfThinBaseVolumes",
				"numOfUnmappedVolumes",
				"numOfVolumes",
				"numOfVolumesInDeletion",
				"numOfVtrees"
			]
		},
		{
			"type": "Sdc",
			"allIds": [],
			"properties": [
				"numOfMappedVolumes",
				"userDataReadBwc",
				"userDataWriteBwc",
				"volumeIds"
			]
		}
	]
}
`
)

type selectedStatistics struct {
	System struct {
		NumOfDevices                             float64
		NumOfFaultSets                           float64
		NumOfMappedToAllVolumes                  float64
		NumOfProtectionDomains                   float64
		NumOfRfcacheDevices                      float64
		NumOfScsiInitiators                      float64
		NumOfSdc                                 float64
		NumOfSds                                 float64
		NumOfSnapshots                           float64
		NumOfStoragePools                        float64
		NumOfThickBaseVolumes                    float64
		NumOfThinBaseVolumes                     float64
		NumOfUnmappedVolumes                     float64
		NumOfVolumes                             float64
		NumOfVolumesInDeletion                   float64
		NumOfVtrees                              float64
		ActiveBckRebuildCapacityInKb             float64
		ActiveFwdRebuildCapacityInKb             float64
		ActiveMovingCapacityInKb                 float64
		ActiveNormRebuildCapacityInKb            float64
		ActiveRebalanceCapacityInKb              float64
		AtRestCapacityInKb                       float64
		BckRebuildCapacityInKb                   float64
		CapacityAvailableForVolumeAllocationInKb float64
		CapacityInUseInKb                        float64
		CapacityLimitInKb                        float64
		DegradedFailedCapacityInKb               float64
		DegradedHealthyCapacityInKb              float64
		FailedCapacityInKb                       float64
		FwdRebuildCapacityInKb                   float64
		InMaintenanceCapacityInKb                float64
		MaxCapacityInKb                          float64
		MovingCapacityInKb                       float64
		NormRebuildCapacityInKb                  float64
		PendingBckRebuildCapacityInKb            float64
		PendingFwdRebuildCapacityInKb            float64
		PendingMovingCapacityInKb                float64
		PendingNormRebuildCapacityInKb           float64
		PendingRebalanceCapacityInKb             float64
		ProtectedCapacityInKb                    float64
		RebalanceCapacityInKb                    float64
		SemiProtectedCapacityInKb                float64
		SnapCapacityInUseInKb                    float64
		SnapCapacityInUseOccupiedInKb            float64
		SpareCapacityInKb                        float64
		ThickCapacityInUseInKb                   float64
		ThinCapacityAllocatedInKb                float64
		ThinCapacityInUseInKb                    float64
		UnreachableUnusedCapacityInKb            float64
		UnusedCapacityInKb                       float64
		NormRebuildReadBwc                       client.Bwc // TODO: ???
		NormRebuildWriteBwc                      client.Bwc // TODO: ???
		BckRebuildReadBwc                        client.Bwc // failed node/disk is back alive
		BckRebuildWriteBwc                       client.Bwc // failed node/disk is back alive
		FwdRebuildReadBwc                        client.Bwc // node/disk fails
		FwdRebuildWriteBwc                       client.Bwc // node/disk fails
		RebalanceReadBwc                         client.Bwc
		RebalanceWriteBwc                        client.Bwc
		PrimaryReadBwc                           client.Bwc // Backend (SDSs + Devices)
		PrimaryWriteBwc                          client.Bwc // Backend (SDSs + Devices)
		SecondaryReadBwc                         client.Bwc // Backend (SDSs + Devices, 2nd)
		SecondaryWriteBwc                        client.Bwc // Backend (SDSs + Devices, 2nd)
		UserDataReadBwc                          client.Bwc // Frontend (Volumes + SDCs)
		UserDataWriteBwc                         client.Bwc // Frontend (Volumes + SDCs)
		PrimaryReadFromDevBwc                    client.Bwc // TODO: ???
		PrimaryReadFromRmcacheBwc                client.Bwc // TODO: ???
		SecondaryReadFromDevBwc                  client.Bwc // TODO: ???
		SecondaryReadFromRmcacheBwc              client.Bwc // TODO: ???
		TotalReadBwc                             client.Bwc // *ReadBwc
		TotalWriteBwc                            client.Bwc // *WriteBwc
	}
	Sdc map[string]struct {
		NumOfMappedVolumes float64
		UserDataReadBwc    client.Bwc
		UserDataWriteBwc   client.Bwc
	}
}

func (s *ScaleIO) collect() (map[string]int64, error) {
	var stats selectedStatistics
	err := s.client.SelectedStatistics(&stats, selectedStatisticsQuery)
	if err != nil {
		return nil, err
	}

	var mx metrics
	s.collectSystemOverview(&mx, stats)
	s.collectSdcStats(&mx, stats)
	return stm.ToMap(mx), nil
}

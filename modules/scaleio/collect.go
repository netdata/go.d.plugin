package scaleio

import (
	"github.com/netdata/go.d.plugin/modules/scaleio/client"
	"github.com/netdata/go.d.plugin/pkg/stm"
)

/*
Starting from version 3 of ScaleIO/VxFlex API numOfScsiInitiators property is removed from the system selectedStatisticsQuery.
Reference: VxFlex OS v3.x REST API Reference Guide.pdf
*/

var query = client.SelectedStatisticsQuery{
	List: []client.SelectedObject{
		{
			Type: "System",
			Properties: []string{
				"maxCapacityInKb",
				"thickCapacityInUseInKb",
				"thinCapacityInUseInKb",
				"snapCapacityInUseOccupiedInKb",
				"spareCapacityInKb",
				"capacityLimitInKb",

				"protectedCapacityInKb",
				"degradedHealthyCapacityInKb",
				"degradedFailedCapacityInKb",
				"failedCapacityInKb",
				"unreachableUnusedCapacityInKb",
				"inMaintenanceCapacityInKb",

				"capacityInUseInKb",
				"capacityAvailableForVolumeAllocationInKb",

				"numOfDevices",
				"numOfFaultSets",
				"numOfProtectionDomains",
				"numOfRfcacheDevices",
				"numOfSdc",
				"numOfSds",
				"numOfSnapshots",
				"numOfStoragePools",
				"numOfVolumes",
				"numOfVtrees",
				"numOfThickBaseVolumes",
				"numOfThinBaseVolumes",
				"numOfMappedToAllVolumes",
				"numOfUnmappedVolumes",

				"rebalanceReadBwc",
				"rebalanceWriteBwc",
				"pendingRebalanceCapacityInKb",

				"pendingNormRebuildCapacityInKb",
				"pendingBckRebuildCapacityInKb",
				"pendingFwdRebuildCapacityInKb",
				"normRebuildReadBwc",
				"normRebuildWriteBwc",
				"bckRebuildReadBwc",
				"bckRebuildWriteBwc",
				"fwdRebuildReadBwc",
				"fwdRebuildWriteBwc",

				"primaryReadBwc",
				"primaryWriteBwc",
				"secondaryReadBwc",
				"secondaryWriteBwc",
				"userDataReadBwc",
				"userDataWriteBwc",
				"totalReadBwc",
				"totalWriteBwc",
			},
		},
		{
			Type:   "StoragePool",
			ALLIDs: true,
			Properties: []string{
				"maxCapacityInKb",
				"thickCapacityInUseInKb",
				"thinCapacityInUseInKb",
				"snapCapacityInUseOccupiedInKb",
				"spareCapacityInKb",
				"capacityLimitInKb",

				"protectedCapacityInKb",
				"degradedHealthyCapacityInKb",
				"degradedFailedCapacityInKb",
				"failedCapacityInKb",
				"unreachableUnusedCapacityInKb",
				"inMaintenanceCapacityInKb",

				"capacityInUseInKb",
				"capacityAvailableForVolumeAllocationInKb",

				"numOfDevices",
				"numOfVolumes",
				"numOfVtrees",
				"numOfSnapshots",
			},
		},
		{
			Type:   "Sdc",
			ALLIDs: true,
			Properties: []string{
				"userDataReadBwc",
				"userDataWriteBwc",

				"numOfMappedVolumes",
			},
		},
	},
}

const discoveryEvery = 5

func (s *ScaleIO) collect() (map[string]int64, error) {
	s.runs += 1
	if !s.lastDiscoveryOK || s.runs%discoveryEvery == 0 {
		if err := s.discovery(); err != nil {
			return nil, err
		}
	}

	stats, err := s.client.SelectedStatistics(query)
	if err != nil {
		return nil, err
	}

	var mx metrics
	s.collectSystemOverview(&mx, stats)
	s.collectSdc(&mx, stats)
	s.collectStoragePool(&mx, stats)
	s.updateCharts()
	return stm.ToMap(mx), nil
}

func (s *ScaleIO) discovery() error {
	ins, err := s.client.Instances()
	if err != nil {
		s.lastDiscoveryOK = false
		return err
	}

	s.discovered.pool = make(map[string]client.StoragePool, len(ins.StoragePool))
	for _, pool := range ins.StoragePool {
		s.discovered.pool[pool.ID] = pool
	}
	s.discovered.sdc = make(map[string]client.Sdc, len(ins.Sdc))
	for _, sdc := range ins.Sdc {
		s.discovered.sdc[sdc.ID] = sdc
	}
	s.lastDiscoveryOK = true
	return nil
}

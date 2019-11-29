package scaleio

import "github.com/netdata/go.d.plugin/modules/scaleio/client"

func (s *ScaleIO) collectStoragePool(mx *metrics, stats client.SelectedStatistics) {
	mx.StoragePool = make(map[string]storagePoolStatistics, len(stats.StoragePool))

	for k, v := range stats.StoragePool {
		if _, ok := s.discovered.pool[k]; !ok {
			continue
		}

		var m storagePoolStatistics
		m.Capacity.AvailableForVolumeAllocation = v.CapacityAvailableForVolumeAllocationInKb
		m.Capacity.MaxCapacity = v.MaxCapacityInKb

		{
			// General Capacity
			m.Capacity.Decreased = sum(v.MaxCapacityInKb, -v.CapacityLimitInKb)
			m.Capacity.Degraded = sum(v.DegradedFailedCapacityInKb, v.DegradedHealthyCapacityInKb)
			m.Capacity.Failed = v.FailedCapacityInKb
			m.Capacity.InMaintenance = v.InMaintenanceCapacityInKb
			m.Capacity.Protected = v.ProtectedCapacityInKb
			m.Capacity.Spare = v.SpareCapacityInKb
			m.Capacity.UnreachableUnused = v.UnreachableUnusedCapacityInKb
			// Note: can't use 'UnusedCapacityInKb' directly, dashboard shows calculated value
			used := sum(
				v.ProtectedCapacityInKb,
				v.InMaintenanceCapacityInKb,
				m.Capacity.Decreased,
				m.Capacity.Degraded,
				v.FailedCapacityInKb,
				v.SpareCapacityInKb,
				v.UnreachableUnusedCapacityInKb,
			)
			m.Capacity.Unused = sum(v.MaxCapacityInKb, -used)
		}

		m.Components.Devices = v.NumOfDevices
		m.Components.Snapshots = v.NumOfSnapshots
		m.Components.Volumes = v.NumOfVolumes
		m.Components.Vtrees = v.NumOfVtrees

		mx.StoragePool[k] = m
	}
}

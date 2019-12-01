package scaleio

import "github.com/netdata/go.d.plugin/modules/scaleio/client"

func (s *ScaleIO) collectStoragePool(mx *metrics, ss client.SelectedStatistics) {
	mx.StoragePool = make(map[string]storagePoolMetrics, len(ss.StoragePool))

	for id, stats := range ss.StoragePool {
		pool, ok := s.discovered.pool[id]
		if !ok {
			continue
		}
		var m storagePoolMetrics
		collectStoragePoolCapacity(&m, stats, pool)
		collectStoragePoolComponents(&m, stats)

		mx.StoragePool[id] = m
	}
}

func collectStoragePoolCapacity(pm *storagePoolMetrics, ps client.StoragePoolStatistics, pool client.StoragePool) {
	collectCapacity(&pm.Capacity.capacity, ps.CapacityStatistics)
	pm.Capacity.Utilization = calcCapacityUtilization(ps.CapacityInUseInKb, ps.MaxCapacityInKb, pool.SparePercentage)
}

func collectStoragePoolComponents(pm *storagePoolMetrics, ps client.StoragePoolStatistics) {
	pm.Components.Devices = ps.NumOfDevices
	pm.Components.Snapshots = ps.NumOfSnapshots
	pm.Components.Volumes = ps.NumOfVolumes
	pm.Components.Vtrees = ps.NumOfVtrees
}

func calcCapacityUtilization(inUse int64, max int64, sparePercent int64) float64 {
	spare := float64(max) / 100 * float64(sparePercent)
	return divFloat(float64(100*inUse), float64(max)-spare)
}

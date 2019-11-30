package scaleio

import "github.com/netdata/go.d.plugin/modules/scaleio/client"

func (s *ScaleIO) collectStoragePool(mx *metrics, stats client.SelectedStatistics) {
	mx.StoragePool = make(map[string]storagePoolStatistics, len(stats.StoragePool))

	for k, v := range stats.StoragePool {
		pool, ok := s.discovered.pool[k]
		if !ok {
			continue
		}
		var m storagePoolStatistics
		collectStoragePoolCapacity(&m, v, pool)
		collectStoragePoolComponents(&m, v)

		mx.StoragePool[k] = m
	}
}

func collectStoragePoolCapacity(m *storagePoolStatistics, s client.StoragePoolStatistics, pool client.StoragePool) {
	collectCapacity(&m.Capacity, s.CapacityStatistics)
	m.Capacity.Utilization = calcCapacityUtilization(s.CapacityInUseInKb, s.MaxCapacityInKb, pool.SparePercentage)
}

func collectStoragePoolComponents(m *storagePoolStatistics, s client.StoragePoolStatistics) {
	m.Components.Devices = s.NumOfDevices
	m.Components.Snapshots = s.NumOfSnapshots
	m.Components.Volumes = s.NumOfVolumes
	m.Components.Vtrees = s.NumOfVtrees
}

func calcCapacityUtilization(inUse int64, max int64, sparePercent int64) float64 {
	spare := float64(max) / 100 * float64(sparePercent)
	return divFloat(
		float64(100*inUse),
		float64(max)-spare,
	)
}

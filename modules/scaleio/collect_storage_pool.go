package scaleio

import "github.com/netdata/go.d.plugin/modules/scaleio/client"

func (s *ScaleIO) collectStoragePool(mx *metrics, stats client.SelectedStatistics) {
	mx.StoragePool = make(map[string]storagePoolStatistics, len(stats.StoragePool))

	for k, v := range stats.StoragePool {
		if _, ok := s.discovered.pool[k]; !ok {
			continue
		}
		var m storagePoolStatistics
		collectStoragePoolCapacity(&m, v)
		collectStoragePoolComponents(&m, v)

		mx.StoragePool[k] = m
	}
}

func collectStoragePoolCapacity(m *storagePoolStatistics, s client.StoragePoolStatistics) {
	collectCapacity(&m.Capacity, s.CapacityStatistics)
}

func collectStoragePoolComponents(m *storagePoolStatistics, s client.StoragePoolStatistics) {
	m.Components.Devices = s.NumOfDevices
	m.Components.Snapshots = s.NumOfSnapshots
	m.Components.Volumes = s.NumOfVolumes
	m.Components.Vtrees = s.NumOfVtrees
}

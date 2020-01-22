package cockroachdb

import "github.com/netdata/go.d.plugin/pkg/prometheus"

func collectStorage(pms prometheus.Metrics) storageMetrics {
	var sm storageMetrics
	collectStorageCapacity(&sm, pms)
	collectStorageLiveBytes(&sm, pms)
	collectStorageRockDB(&sm, pms)
	collectStorageFD(&sm, pms)
	collectStorageTimeSeries(&sm, pms)
	return sm
}

func collectStorageCapacity(sm *storageMetrics, pms prometheus.Metrics) {
	sm.Capacity.Total = pms.FindByName("capacity").Max()
	sm.Capacity.Available = pms.FindByName("capacity_available").Max()
	sm.Capacity.Used = pms.FindByName("capacity_used").Max()
	sm.Capacity.Reserved = pms.FindByName("capacity_reserved").Max()
	sm.Capacity.PercentageUsed = calcCapacityUsedPercentage(*sm)
}

func collectStorageLiveBytes(sm *storageMetrics, pms prometheus.Metrics) {
	sm.LiveBytes = pms.FindByName("livebytes").Max()
	sm.SysBytes = pms.FindByName("sysbytes").Max()
}

func collectStorageRockDB(sm *storageMetrics, pms prometheus.Metrics) {
	sm.RocksDB.ReadAmplifications = pms.FindByName("rocksdb_read_amplification").Max()
	sm.RocksDB.SSTables = pms.FindByName("rocksdb_num_sstables").Max()
	sm.RocksDB.Compactions = pms.FindByName("rocksdb_compactions").Max()
	sm.RocksDB.Flushes = pms.FindByName("rocksdb_flushes").Max()
	sm.RocksDB.BlockCache.Hits = pms.FindByName("rocksdb_block_cache_hits").Max()
	sm.RocksDB.BlockCache.Misses = pms.FindByName("rocksdb_block_cache_misses").Max()
	sm.RocksDB.BlockCache.HitRate = calcRocksDBCacheHitRate(*sm)
}

func collectStorageFD(sm *storageMetrics, pms prometheus.Metrics) {
	sm.FileDescriptors.Open = pms.FindByName("sys_fd_open").Max()
	sm.FileDescriptors.SoftLimit = pms.FindByName("sys_fd_softlimit").Max()

}

func collectStorageTimeSeries(sm *storageMetrics, pms prometheus.Metrics) {
	sm.TimeSeries.WriteSamples = pms.FindByName("timeseries_write_samples").Max()
	sm.TimeSeries.WriteErrors = pms.FindByName("timeseries_write_errors").Max()
	sm.TimeSeries.WriteBytes = pms.FindByName("timeseries_write_bytes").Max()
}

func calcCapacityUsedPercentage(sm storageMetrics) float64 {
	if sm.Capacity.Total == 0 {
		return 100
	}
	return (1 - sm.Capacity.Available/sm.Capacity.Total) * 100 * 1000
}

func calcRocksDBCacheHitRate(sm storageMetrics) float64 {
	if sm.RocksDB.BlockCache.Hits+sm.RocksDB.BlockCache.Misses == 0 {
		return 0
	}
	return sm.RocksDB.BlockCache.Hits / (sm.RocksDB.BlockCache.Hits + sm.RocksDB.BlockCache.Misses) * 100 * 1000
}

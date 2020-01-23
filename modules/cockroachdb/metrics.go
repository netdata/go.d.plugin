package cockroachdb

// https://github.com/cockroachdb/cockroach/blob/master/pkg/storage/metrics.go
// https://github.com/cockroachdb/cockroach/blob/master/pkg/ts/metrics.go
// https://github.com/cockroachdb/cockroach/blob/master/pkg/server/status/runtime.go

type metrics struct {
	Storage storageMetrics `stm:"storage"`
	Runtime runtimeMetrics `stm:"runtime"`
}

type (
	storageMetrics struct {
		Capacity struct {
			Total                float64 `stm:"total"`
			Used                 float64 `stm:"used"`
			Reserved             float64 `stm:"reserved"`
			Available            float64 `stm:"available"`
			Unusable             float64 `stm:"unusable"`
			Usable               float64 `stm:"usable"`
			TotalUsedPercentage  float64 `stm:"total_used_percentage"`
			UsableUsedPercentage float64 `stm:"usable_used_percentage"`
		} `stm:"capacity"`
		LiveBytes float64 `stm:"live_bytes"`
		SysBytes  float64 `stm:"sys_bytes"`
		RocksDB   struct {
			ReadAmplifications float64 `stm:"read_amplification"`
			SSTables           float64 `stm:"num_sstables"`
			BlockCache         struct {
				Bytes   float64 `stm:"bytes"`
				Hits    float64 `stm:"hits"`
				Misses  float64 `stm:"misses"`
				HitRate float64 `stm:"hit_rate"`
			} `stm:"block_cache"`
			Compactions float64 `stm:"compactions"`
			Flushes     float64 `stm:"flushes"`
		} `stm:"rocksdb"`
		FileDescriptors struct {
			Open      float64 `stm:"open"`
			SoftLimit float64 `stm:"soft_limit"`
		} `stm:"file_descriptors"`
		TimeSeries struct {
			WriteSamples float64 `stm:"write_samples"`
			WriteErrors  float64 `stm:"write_errors"`
			WriteBytes   float64 `stm:"write_bytes"`
		} `stm:"timeseries"`
	}
)

type (
	runtimeMetrics struct {
		LiveNodes float64 `stm:"live_nodes"`
		SysUptime float64 `stm:"uptime"`
		Memory    struct {
			RSS           float64 `stm:"rss"`
			GoAllocBytes  float64 `stm:"go_alloc_bytes"`
			GoTotalBytes  float64 `stm:"go_total_bytes"`
			CGoAllocBytes float64 `stm:"cgo_alloc_bytes"`
			CGoTotalBytes float64 `stm:"cgo_total_bytes"`
		} `stm:"memory"`
	}
)

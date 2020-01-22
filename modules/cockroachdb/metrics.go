package cockroachdb

type metrics struct {
	Storage storageMetrics `stm:"storage"`
}

type (
	storageMetrics struct {
		Capacity struct {
			Total          float64 `stm:"total"`
			Available      float64 `stm:"available"`
			Used           float64 `stm:"used"`
			PercentageUsed float64 `stm:"percentage_used"`
		} `stm:"capacity"`
		LiveBytes float64 `stm:"live_bytes"`
		SysBytes  float64 `stm:"sys_bytes"`
		RocksDB   struct {
			ReadAmplifications float64 `stm:"read_amplifications"`
			SSTables           float64 `stm:"num_sstables"`
			BlockCache         struct {
				Hits    float64 `stm:"hits"`
				Miss    float64 `stm:"hits"`
				HitRate float64 `stm:"hit_rate,1000,1"`
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

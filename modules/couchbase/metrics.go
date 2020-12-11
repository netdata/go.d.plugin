package couchbase

type cbMetrics struct {
	// https://developer.couchbase.com/resources/best-practice-guides/monitoring-guide.pdf
	BucketsStats []bucketsStats
}

func (m cbMetrics) empty() bool {
	switch {
	case m.hasBucketsStats():
		return false
	}
	return true
}
func (m cbMetrics) hasBucketsStats() bool { return len(m.BucketsStats) > 0 }

type bucketsStats struct {
	Name string `stm:"name" json:"name"`

	BasicStats struct {
		DataUsed               float64 `stm:"dataUsed" json:"dataUsed"`
		DiskFetches            float64 `stm:"diskFetches" json:"diskFetches"`
		ItemCount              float64 `stm:"itemCount" json:"itemCount"`
		DiskUsed               float64 `stm:"diskUsed" json:"diskUsed"`
		MemUsed                float64 `stm:"memUsed" json:"memUsed"`
		OpsPerSec              float64 `stm:"opsPerSec" json:"opsPerSec"`
		QuotaPercentUsed       float64 `stm:"quotaPercentUsed" json:"quotaPercentUsed"`
		VbActiveNumNonResident float64 `stm:"vbActiveNumNonResident" json:"vbActiveNumNonResident"`
	} `stm:"basicStats"`
}

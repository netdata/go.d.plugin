package dnsdist

// https://dnsdist.org/guides/webserver.html#get--jsonstat
// https://dnsdist.org/statistics.html

type statisticMetrics struct {
	AclDrops int64 `stm:"acl-drops" json:"acl-drops"` 
	CacheHits int64 `stm:"cache-hits" json:"cache-hits"`
	CacheMisses int64 `stm:"cache-misses" json:"cache-misses"`
	CPUSysMsec int64 `stm:"cpu-sys-msec" json:"cpu-sys-msec"` 
	CPUUserMsec int64 `stm:"cpu-user-msec" json:"cpu-user-msec"` 
	DownStreamSendErrors int64 `stm:"downstream-send-errors" json:"downstream-send-errors"` 
	DownStreamTimeout int64 `stm:"downstream-timeouts" json:"downstream-timeouts"` 
	DynBlocked int64 `stm:"dyn-blocked" json:"dyn-blocked"` 
	EmptyQueries int64 `stm:"empty-queries" json:"empty-queries"` 
	LatencyAvg100 int64  `stm:"empty-queries" json:"latency-avg100"`
	LatencyAvg1000 int64  `stm:"empty-queries" json:"latency-avg1000"`
	LatencyAvg10000 int64  `stm:"empty-queries" json:"latency-avg10000"`
	LatencyAvg1000000 int64  `stm:"empty-queries" json:"latency-avg1000000"`
	LatencySlow int64 `stm:"latency-slow" json:"latency-slow"` 
	Latency0 int64 `stm:"latency0-1" json:"latency0-1"` 
	Latency1 int64 `stm:"latency1-10" json:"latency1-10"` 
	Latency10 int64 `stm:"latency10-50" json:"latency10-50"` 
	Latency100 int64 `stm:"latency100-1000" json:"latency100-1000"` 
	Latency50 int64 `stm:"latency50-100" json:"latency50-100"` 
	NoPolicy int64 `stm:"no-policy" json:"no-policy"` 
	NonCompliantQueries int64 `stm:"noncompliant-queries" json:"noncompliant-queries"` 
	NonCompliantResponses int64 `stm:"noncompliant-responses" json:"noncompliant-responses"` 
	Queries int64 `stm:"queries" json:"queries"` 
	RdQueries int64 `stm:"rdqueries" json:"rdqueries"` 
	RealMemoryUsage int64 `stm:"real-memory-usage" json:"real-memory-usage"` 
	Responses int64 `stm:"responses" json:"responses"` 
	RuleDrop int64 `stm:"rule-drop" json:"rule-drop"` 
	RuleNxDomain int64 `stm:"rule-nxdomain" json:"rule-nxdomain"` 
	RuleRefused int64 `stm:"rule-refused" json:"rule-refused"` 
	SelfAnswered int64 `stm:"self-answered" json:"self-answered"` 
	ServFailResponses int64 `stm:"servfail-responses" json:"servfail-responses"` 
	TruncFailures int64 `stm:"trunc-failures" json:"trunc-failures"` 
}


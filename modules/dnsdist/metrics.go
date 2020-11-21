package dnsdist

// https://dnsdist.org/guides/webserver.html#get--jsonstat
// https://dnsdist.org/statistics.html

type statisticMetrics struct {
	AclDrops int64 `stm:"acl-drops" json:"acl-drops"` 
	CacheHits int64 `stm:"cache-hits" json:"cache-hits"`
	CacheMisses int64 `stm:"cache-misses" json:"cache-misses"`
	CPUIowait int64 `stm:"cpu-iowait" json:"cpu-iowait"` 
	CPUSteal int64 `stm:"cpu-steal" json:"cpu-steal"` 
	CPUSysMsec int64 `stm:"cpu-sys-msec" json:"cpu-sys-msec"` 
	CPUUserMsec int64 `stm:"cpu-user-msec" json:"cpu-user-msec"` 
	DohQueryPipeFull int64 `stm:"doh-query-pipe-full" json:"doh-query-pipe-full"` 
	DohResponsePipeFull int64 `stm:"doh-response-pipe-full" json:"doh-response-pipe-full"`
	DownStreamSendErrors int64 `stm:"downstream-send-errors" json:"downstream-send-errors"` 
	DownStreamTimeout int64 `stm:"downstream-timeouts" json:"downstream-timeouts"` 
	DynBlockNmgSize int64 `stm:"dyn-block-nmg-size" json:"dyn-block-nmg-size"` 
	DynBlocked int64 `stm:"dyn-blocked" json:"dyn-blocked"` 
	EmptyQueries int64 `stm:"empty-queries" json:"empty-queries"` 
	FdUsage int64 `stm:"fd-usage" json:"fd-usage"` 
	FrontendNoError int64 `stm:"frontend-noerror" json:"frontend-noerror"` 
	FrontEndNxDomain int64 `stm:"frontend-nxdomain" json:"frontend-nxdomain"` 
	FrontendServFail int64 `stm:"frontend-servfail" json:"frontend-servfail"` 
	LatencyAvg100 float64 `stm:"latency-avg100" json:"latency-avg100"` 
	LatencyAvg1000 float64 `stm:"latency-avg1000" json:"latency-avg1000"` 
	LatencyAvg10000 float64 `stm:"latency-avg10000" json:"latency-avg10000"` 
	LatencyAvg100000 float64 `stm:"latency-avg1000000" json:"latency-avg1000000"` 
	LatencyCount int64 `stm:"latency-count" json:"latency-count"` 
	LatencySlow int64 `stm:"latency-slow" json:"latency-slow"` 
	LatencySum int64 `stm:"latency-sum" json:"latency-sum"` 
	Latency0 int64 `stm:"latency0-1" json:"latency0-1"` 
	Latency1 int64 `stm:"latency1-10" json:"latency1-10"` 
	Latency10 int64 `stm:"latency10-50" json:"latency10-50"` 
	Latency100 int64 `stm:"latency100-1000" json:"latency100-1000"` 
	Latency50 int64 `stm:"latency50-100" json:"latency50-100"` 
	NoPolicy int64 `stm:"no-policy" json:"no-policy"` 
	NonCompliantQueries int64 `stm:"noncompliant-queries" json:"noncompliant-queries"` 
	NonCompliantResponses int64 `stm:"noncompliant-responses" json:"noncompliant-responses"` 
	CpacityDrops int64 `stm:"over-capacity-drops" json:"over-capacity-drops"` 
	PacketcacheHits int64 `stm:"packetcache-hits" json:"packetcache-hits"` 
	PacketCacheMisses int64 `stm:"packetcache-misses" json:"packetcache-misses"` 
	Queries int64 `stm:"queries" json:"queries"` 
	RdQueries int64 `stm:"rdqueries" json:"rdqueries"` 
	RealMemoryUsage int64 `stm:"real-memory-usage" json:"real-memory-usage"` 
	Responses int64 `stm:"responses" json:"responses"` 
	RuleDrop int64 `stm:"rule-drop" json:"rule-drop"` 
	RuleNxDomain int64 `stm:"rule-nxdomain" json:"rule-nxdomain"` 
	RuleRefused int64 `stm:"rule-refused" json:"rule-refused"` 
	RuleServFail int64 `stm:"rule-servfail" json:"rule-servfail"` 
	SecurityStatus int64 `stm:"security-status" json:"security-status"` 
	SelfAnswered int64 `stm:"self-answered" json:"self-answered"` 
	SereverPolicy int64 `stm:"server-policy" json:"server-policy"` 
	FirstAvailable int64 `stm:"firstAvailable" json:"firstAvailable"` 
	ServFailResponses int64 `stm:"servfail-responses" json:"servfail-responses"` 
	TooOldDrops int64 `stm:"too-old-drops" json:"too-old-drops"` 
	TruncFailures int64 `stm:"trunc-failures" json:"trunc-failures"` 
	UDPInErrors int64`stm:"udp-in-errors" json:"udp-in-errors"` 
	UDPNoReportErrors int64 `stm:"udp-noport-errors" json:"udp-noport-errors"` 
	UDPRecvBufErrors int64 `stm:"udp-recvbuf-errors" json:"udp-recvbuf-errors"` 
	UDPSndbufErrors int64 `stm:"udp-sndbuf-errors" json:"udp-sndbuf-errors"` 
	UpTime int64 `stm:"uptime" json:"uptime"`	
}
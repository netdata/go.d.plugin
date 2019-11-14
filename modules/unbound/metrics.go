package unbound

// https://github.com/NLnetLabs/unbound/blob/master/smallapp/unbound-control.c

type (
	cache struct {
		Hits float64 `stm:"hits"`
		Miss float64 `stm:"miss"`
	}
	dnsCrypt struct {
		Crypted   float64 `stm:"crypted"`
		Cert      float64 `stm:"cert"`
		ClearText float64 `stm:"clear_text"`
		Malformed float64 `stm:"malformed"`
	}
	requestList struct {
		Avg         float64 `stm:"avg"`
		Max         float64 `stm:"max"`
		Overwritten float64 `stm:"overwritten"`
		Exceeded    float64 `stm:"exceeded"`
		CurrentAll  float64 `stm:"current_all"`
		CurrentUser float64 `stm:"current_user"`
	}
	recursionTime struct {
		Avg    float64 `stm:"avg"`
		Median float64 `stm:"median"`
	}
	common struct {
		Queries          float64       `stm:"queries"`
		QueriesIPRL      float64       `stm:"queries_ip_ratelimited"`
		Cache            cache         `stm:"cache"`
		Prefetch         float64       `stm:"prefetch"`
		ZeroTTL          float64       `stm:"zero_ttl"`
		RecursiveReplies float64       `stm:"recursive_replies"`
		DNSCrypt         dnsCrypt      `stm:"dns_crypt"`
		RequestList      requestList   `stm:"request_list"`
		RecursionTime    recursionTime `stm:"recursion_time"`
		TCPUsage         float64       `stm:"tcp_usage"`
	}
	extended struct {
		QueryType   map[string]float64 `stm:"query_type"`
		QueryClass  map[string]float64 `stm:"query_class"`
		QueryOpCode map[string]float64 `stm:"query_opcode"`
		AnswerRCode map[string]float64 `stm:"answer_rcode"`
	}
)

type metricsData struct {
	common   `stm:"total"`
	extended `stm:""`
	Uptime   float64            `stm:"uptime"`
	Thread   map[string]*common `stm:"thread"`
}

func newMetricsData() *metricsData {
	return &metricsData{
		extended: extended{
			QueryType:   make(map[string]float64),
			QueryClass:  make(map[string]float64),
			QueryOpCode: make(map[string]float64),
			AnswerRCode: make(map[string]float64),
		},
		Thread: make(map[string]*common),
	}
}

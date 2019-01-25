package weblog

import (
	"sync"

	"github.com/netdata/go.d.plugin/pkg/metrics"
)

type metricsData struct {
	// mutex for all values
	Mux sync.RWMutex

	Requests metrics.Counter `stm:"requests"`

	ReqSuccessful metrics.Counter `stm:"req_successful"`
	ReqRedirect   metrics.Counter `stm:"req_redirect"`
	ReqBad        metrics.Counter `stm:"req_bad"`
	ReqError      metrics.Counter `stm:"req_error"`

	Req1xx       metrics.Counter `stm:"req_1xx"`
	Req2xx       metrics.Counter `stm:"req_2xx"`
	Req3xx       metrics.Counter `stm:"req_3xx"`
	Req4xx       metrics.Counter `stm:"req_4xx"`
	Req5xx       metrics.Counter `stm:"req_5xx"`
	ReqUnmatched metrics.Counter `stm:"req_unmatched"`

	ReqCode metrics.CounterVec `stm:"req_code"`

	ReqMethod metrics.CounterVec `stm:"req_method"`

	BytesSent     metrics.Counter `stm:"bytes_sent"`
	BytesReceived metrics.Counter `stm:"bytes_received"`

	RespTime             metrics.Summary   `stm:"resp_time"`
	RespTimeHist         metrics.Histogram `stm:"resp_time_hist"`
	RespTimeUpstream     metrics.Summary   `stm:"resp_time_upstream"`
	RespTimeUpstreamHist metrics.Histogram `stm:"resp_time_upstream_hist"`

	UniqueIPs metrics.UniqueCounter `stm:"uniq_ips"`

	CategorizedRespTime metrics.Summary `stm:"cat_resp_time"`
}

func newMetricsData(config Config) *metricsData {
	return &metricsData{
		ReqCode:              metrics.NewCounterVec(),
		ReqMethod:            metrics.NewCounterVec(),
		RespTime:             metrics.NewSummary(),
		RespTimeHist:         metrics.NewHistogram(config.Histogram),
		RespTimeUpstream:     metrics.NewSummary(),
		RespTimeUpstreamHist: metrics.NewHistogram(config.Histogram),
		UniqueIPs:            metrics.NewUniqueCounter(false),
	}
}

package weblog

import (
	"github.com/netdata/go.d.plugin/pkg/metrics"
)

type MetricsData struct {
	BytesSent     metrics.Counter `stm:"bytes_sent"`
	BytesReceived metrics.Counter `stm:"bytes_received"`

	Requests metrics.Counter `stm:"requests"`

	ReqUnmatched metrics.Counter `stm:"req_unmatched"`

	ReqMethod metrics.CounterVec `stm:"req_method"`

	ReqVersion metrics.CounterVec `stm:"req_version"`

	RespCode metrics.CounterVec `stm:"req_code"`

	RespSuccessful  metrics.Counter `stm:"resp_successful"`
	RespRedirect    metrics.Counter `stm:"resp_redirect"`
	RespClientError metrics.Counter `stm:"resp_client_error"`
	RespServerError metrics.Counter `stm:"resp_server_error"`

	Resp1xx metrics.Counter `stm:"resp_1xx"`
	Resp2xx metrics.Counter `stm:"resp_2xx"`
	Resp3xx metrics.Counter `stm:"resp_3xx"`
	Resp4xx metrics.Counter `stm:"resp_4xx"`
	Resp5xx metrics.Counter `stm:"resp_5xx"`

	RespTime             metrics.Summary   `stm:"resp_time,1000"`
	RespTimeHist         metrics.Histogram `stm:"resp_time_hist"`
	RespTimeUpstream     metrics.Summary   `stm:"resp_time_upstream,1000"`
	RespTimeUpstreamHist metrics.Histogram `stm:"resp_time_upstream_hist"`

	UniqueIPs metrics.UniqueCounter `stm:"uniq_ips"`

	CategorizedRequests metrics.CounterVec `stm:"cat_req"`
	CategorizedRespTime metrics.Summary    `stm:"cat_resp_time"`
}

func NewMetricsData(config Config) *MetricsData {
	return &MetricsData{
		RespCode:             metrics.NewCounterVec(),
		ReqMethod:            metrics.NewCounterVec(),
		ReqVersion:           metrics.NewCounterVec(),
		RespTime:             metrics.NewSummary(),
		RespTimeHist:         metrics.NewHistogram(config.Histogram),
		RespTimeUpstream:     metrics.NewSummary(),
		RespTimeUpstreamHist: metrics.NewHistogram(config.Histogram),
		UniqueIPs:            metrics.NewUniqueCounter(true),
		CategorizedRequests:  metrics.NewCounterVec(),
		CategorizedRespTime:  metrics.NewSummary(),
	}
}

func (m *MetricsData) Reset() {
	m.UniqueIPs.Reset()
}

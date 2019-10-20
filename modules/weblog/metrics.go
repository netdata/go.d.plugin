package weblog

import (
	"github.com/netdata/go.d.plugin/pkg/metrics"
)

type weblogSummary struct {
	metrics.Summary
}

// TODO: temporary workaround?
func (s weblogSummary) WriteTo(rv map[string]int64, key string, mul, div int) {
	s.Summary.WriteTo(rv, key, mul, div)
	if _, ok := rv[key+"_min"]; !ok {
		rv[key+"_min"] = 0
		rv[key+"_max"] = 0
		rv[key+"_avg"] = 0
	}
}

type MetricsData struct {
	Requests     metrics.Counter `stm:"requests"`
	ReqUnmatched metrics.Counter `stm:"req_unmatched"`

	ReqVhost metrics.CounterVec `stm:"req_vhost"`

	ReqIpv4    metrics.Counter       `stm:"req_ipv4"`
	ReqIpv6    metrics.Counter       `stm:"req_ipv6"`
	UniqueIPv4 metrics.UniqueCounter `stm:"unique_current_poll_ipv4"`
	UniqueIPv6 metrics.UniqueCounter `stm:"unique_current_poll_ipv6"`

	ReqMethod  metrics.CounterVec `stm:"req_method"`
	ReqURI     metrics.CounterVec `stm:"req_uri"`
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

	BytesSent     metrics.Counter `stm:"bytes_sent"`
	BytesReceived metrics.Counter `stm:"bytes_received"`

	RespTime             metrics.Summary   `stm:"resp_time,1000"`
	RespTimeHist         metrics.Histogram `stm:"resp_time_hist"`
	RespTimeUpstream     metrics.Summary   `stm:"resp_time_upstream,1000"`
	RespTimeUpstreamHist metrics.Histogram `stm:"resp_time_upstream_hist"`

	ReqCustom metrics.CounterVec `stm:"req_custom"`

	CategorizedRequests  metrics.CounterVec `stm:"cat_requests"`
	CategorizedBandwidth metrics.CounterVec `stm:"cat_bandwidth"`
	CategorizedRespTime  metrics.Summary    `stm:"cat_resp_time"`
}

func NewMetricsData(config Config) *MetricsData {
	return &MetricsData{
		RespCode:             metrics.NewCounterVec(),
		ReqMethod:            metrics.NewCounterVec(),
		ReqVersion:           metrics.NewCounterVec(),
		RespTime:             &weblogSummary{metrics.NewSummary()},
		RespTimeHist:         metrics.NewHistogram(config.Histogram),
		RespTimeUpstream:     &weblogSummary{metrics.NewSummary()},
		RespTimeUpstreamHist: metrics.NewHistogram(config.Histogram),
		UniqueIPv4:           metrics.NewUniqueCounter(true),
		UniqueIPv6:           metrics.NewUniqueCounter(true),
		CategorizedRequests:  metrics.NewCounterVec(),
		CategorizedRespTime:  metrics.NewSummary(),

		ReqURI:    newCounterVecFromCategories(config.URLCategories),
		ReqCustom: newCounterVecFromCategories(config.UserCategories),
		ReqVhost:  metrics.NewCounterVec(),
	}
}

func (m *MetricsData) Reset() {
	m.UniqueIPv4.Reset()
	m.UniqueIPv6.Reset()
	m.RespTime.Reset()
	m.RespTimeUpstream.Reset()
}

func newCounterVecFromCategories(cats []rawCategory) metrics.CounterVec {
	c := metrics.NewCounterVec()
	for _, v := range cats {
		_, _ = c.GetP(v.Name)
	}
	return c
}

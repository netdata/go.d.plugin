package squidlog

import "github.com/netdata/go.d.plugin/pkg/metrics"

func newSummary() metrics.Summary {
	return &summary{metrics.NewSummary()}
}

type summary struct {
	metrics.Summary
}

func (s summary) WriteTo(rv map[string]int64, key string, mul, div int) {
	s.Summary.WriteTo(rv, key, mul, div)
	if _, ok := rv[key+"_min"]; !ok {
		rv[key+"_min"] = 0
		rv[key+"_max"] = 0
		rv[key+"_avg"] = 0
	}
}

type metricsData struct {
	Requests     metrics.Counter `stm:"requests"`
	ReqUnmatched metrics.Counter `stm:"req_unmatched"`

	HTTPCode metrics.CounterVec `stm:"http_code"`
	HTTP1xx  metrics.Counter    `stm:"http_1xx"`
	HTTP2xx  metrics.Counter    `stm:"http_2xx"`
	HTTP3xx  metrics.Counter    `stm:"http_3xx"`
	HTTP4xx  metrics.Counter    `stm:"http_4xx"`
	HTTP5xx  metrics.Counter    `stm:"http_5xx"`

	ReqSuccess  metrics.Counter `stm:"req_type_success"`
	ReqRedirect metrics.Counter `stm:"req_type_redirect"`
	ReqBad      metrics.Counter `stm:"req_type_bad"`
	ReqError    metrics.Counter `stm:"req_type_error"`

	BytesSent metrics.Counter `stm:"bytes_sent"`
	RespTime  metrics.Summary `stm:"resp_time"`

	UniqueClients metrics.UniqueCounter `stm:"uniq_clients"`

	ReqMethod           metrics.CounterVec `stm:"req_method"`
	CacheCode           metrics.CounterVec `stm:"cache_code"`
	CacheCodeTransport  metrics.CounterVec `stm:"cache_code_transport"`
	CacheCodeHandling   metrics.CounterVec `stm:"cache_code_handling"`
	CacheCodeObject     metrics.CounterVec `stm:"cache_code_object"`
	CacheCodeLoadSource metrics.CounterVec `stm:"cache_code_load_source"`
	CacheCodeError      metrics.CounterVec `stm:"cache_code_error"`

	HierCode metrics.CounterVec `stm:"hier_code"`
	MimeType metrics.CounterVec `stm:"mime_type"`
}

func (m *metricsData) reset() {
	m.RespTime.Reset()
	m.UniqueClients.Reset()
}

func newMetricsData() *metricsData {
	return &metricsData{
		RespTime:            newSummary(),
		UniqueClients:       metrics.NewUniqueCounter(true),
		HTTPCode:            metrics.NewCounterVec(),
		ReqMethod:           metrics.NewCounterVec(),
		CacheCode:           metrics.NewCounterVec(),
		CacheCodeTransport:  metrics.NewCounterVec(),
		CacheCodeHandling:   metrics.NewCounterVec(),
		CacheCodeObject:     metrics.NewCounterVec(),
		CacheCodeLoadSource: metrics.NewCounterVec(),
		CacheCodeError:      metrics.NewCounterVec(),
		HierCode:            metrics.NewCounterVec(),
		MimeType:            metrics.NewCounterVec(),
	}
}

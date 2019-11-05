package weblog

import (
	"github.com/netdata/go.d.plugin/pkg/metrics"
)

func newWebLogSummary() metrics.Summary {
	return &weblogSummary{metrics.NewSummary()}
}

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

type (
	MetricsData struct {
		Requests     metrics.Counter `stm:"requests"`
		ReqUnmatched metrics.Counter `stm:"req_unmatched"`
		ReqFiltered  metrics.Counter `stm:"req_filtered"`

		RespStatusCode  metrics.CounterVec `stm:"resp_status_code"`
		RespSuccessful  metrics.Counter    `stm:"resp_successful"`
		RespRedirect    metrics.Counter    `stm:"resp_redirect"`
		RespClientError metrics.Counter    `stm:"resp_client_error"`
		RespServerError metrics.Counter    `stm:"resp_server_error"`
		Resp1xx         metrics.Counter    `stm:"resp_1xx"`
		Resp2xx         metrics.Counter    `stm:"resp_2xx"`
		Resp3xx         metrics.Counter    `stm:"resp_3xx"`
		Resp4xx         metrics.Counter    `stm:"resp_4xx"`
		Resp5xx         metrics.Counter    `stm:"resp_5xx"`

		UniqueIPv4      metrics.UniqueCounter `stm:"uniq_ipv4"`
		UniqueIPv6      metrics.UniqueCounter `stm:"uniq_ipv6"`
		BytesSent       metrics.Counter       `stm:"bytes_sent"`
		BytesReceived   metrics.Counter       `stm:"bytes_received"`
		ReqProcTime     metrics.Summary       `stm:"req_proc_time"`
		ReqProcTimeHist metrics.Histogram     `stm:"req_proc_time_hist"`
		UpsRespTime     metrics.Summary       `stm:"upstream_resp_time"`
		UpsRespTimeHist metrics.Histogram     `stm:"upstream_resp_time_hist"`

		ReqVhost          metrics.CounterVec `stm:"req_vhost"`
		ReqPort           metrics.CounterVec `stm:"req_port"`
		ReqHTTPScheme     metrics.Counter    `stm:"req_http_scheme"`
		ReqHTTPSScheme    metrics.Counter    `stm:"req_https_scheme"`
		ReqIPv4           metrics.Counter    `stm:"req_ipv4"`
		ReqIPv6           metrics.Counter    `stm:"req_ipv6"`
		ReqMethod         metrics.CounterVec `stm:"req_method"`
		ReqURLPattern     metrics.CounterVec `stm:"req_url_ptn"`
		ReqVersion        metrics.CounterVec `stm:"req_version"`
		ReqCustomPattern  metrics.CounterVec `stm:"req_custom_ptn"`
		ReqSSLProto       metrics.CounterVec `stm:"req_ssl_proto"`
		ReqSSLCipherSuite metrics.CounterVec `stm:"req_ssl_cipher_suite"`

		URLPatternStats patternStats `stm:"url_ptn"`
	}

	patternMetrics struct {
		RespStatusCode metrics.CounterVec `stm:"resp_status_code"`
		BytesSent      metrics.Counter    `stm:"bytes_sent"`
		BytesReceived  metrics.Counter    `stm:"bytes_received"`
		ReqProcTime    metrics.Summary    `stm:"req_proc_time"`
	}

	patternStats map[string]*patternMetrics
)

func NewMetricsData(config Config) *MetricsData {
	return &MetricsData{
		ReqVhost:          metrics.NewCounterVec(),
		ReqPort:           metrics.NewCounterVec(),
		RespStatusCode:    metrics.NewCounterVec(),
		ReqMethod:         metrics.NewCounterVec(),
		ReqVersion:        metrics.NewCounterVec(),
		ReqSSLProto:       metrics.NewCounterVec(),
		ReqSSLCipherSuite: metrics.NewCounterVec(),
		ReqProcTime:       newWebLogSummary(),
		ReqProcTimeHist:   metrics.NewHistogram(config.Histogram),
		UpsRespTime:       newWebLogSummary(),
		UpsRespTimeHist:   metrics.NewHistogram(config.Histogram),
		UniqueIPv4:        metrics.NewUniqueCounter(true),
		UniqueIPv6:        metrics.NewUniqueCounter(true),
		ReqURLPattern:     newCounterVecFromPatterns(config.URLPatterns),
		ReqCustomPattern:  newCounterVecFromPatterns(config.CustomPatterns),
		URLPatternStats:   newPatternStats(config.URLPatterns),
	}
}

func (m *MetricsData) Reset() {
	m.UniqueIPv4.Reset()
	m.UniqueIPv6.Reset()
	m.ReqProcTime.Reset()
	m.UpsRespTime.Reset()
	for _, v := range m.URLPatternStats {
		v.ReqProcTime.Reset()
	}
}

func newPatternStats(ps []userPattern) patternStats {
	stats := make(patternStats)
	for _, p := range ps {
		stats[p.Name] = &patternMetrics{
			RespStatusCode: metrics.NewCounterVec(),
			ReqProcTime:    newWebLogSummary(),
		}
	}
	return stats
}

func newCounterVecFromPatterns(ps []userPattern) metrics.CounterVec {
	c := metrics.NewCounterVec()
	for _, p := range ps {
		_, _ = c.GetP(p.Name)
	}
	return c
}

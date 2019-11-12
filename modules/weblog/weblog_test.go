package weblog

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/netdata/go.d.plugin/pkg/logs"
	"github.com/netdata/go.d.plugin/pkg/metrics"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testCommonLog, _ = ioutil.ReadFile("testdata/common.log")
	testFullLog, _   = ioutil.ReadFile("testdata/full.log")
	testCustomLog, _ = ioutil.ReadFile("testdata/custom.log")
)

func Test_readTestData(t *testing.T) {
	assert.NotNil(t, testFullLog)
	assert.NotNil(t, testCommonLog)
	assert.NotNil(t, testCustomLog)
}

func TestNew(t *testing.T) {
	// TODO:
}

func TestWebLog_Init(t *testing.T) {
	// TODO:
}

func TestWebLog_Check(t *testing.T) {
	// TODO:
}

func TestWebLog_Charts(t *testing.T) {
	// TODO:
}

func TestWebLog_Cleanup(t *testing.T) {
	// TODO:
}

func TestWebLog_Collect(t *testing.T) {
	weblog := prepareWebLogCollectFull(t)

	//m := weblog.Collect()
	//l := make([]string, 0)
	//for k := range m {
	//	l = append(l, k)
	//}
	//sort.Strings(l)
	//for _, value := range l {
	//	fmt.Println(fmt.Sprintf("\"%s\": %d,", value, m[value]))
	//}

	expected := map[string]int64{
		"bytes_received":                            1384573,
		"bytes_sent":                                1352758,
		"custom_field_drink_beer":                   223,
		"custom_field_drink_wine":                   232,
		"custom_field_side_dark":                    230,
		"custom_field_side_light":                   225,
		"req_bad":                                   54,
		"req_error":                                 0,
		"req_http_scheme":                           236,
		"req_https_scheme":                          219,
		"req_ipv4":                                  275,
		"req_ipv6":                                  180,
		"req_method_GET":                            162,
		"req_method_HEAD":                           138,
		"req_method_POST":                           155,
		"req_port_80":                               119,
		"req_port_81":                               101,
		"req_port_82":                               81,
		"req_port_83":                               74,
		"req_port_84":                               80,
		"req_proc_time_avg":                         246,
		"req_proc_time_count":                       455,
		"req_proc_time_hist_bucket_1":               0,
		"req_proc_time_hist_bucket_10":              3,
		"req_proc_time_hist_bucket_11":              6,
		"req_proc_time_hist_bucket_2":               0,
		"req_proc_time_hist_bucket_3":               0,
		"req_proc_time_hist_bucket_4":               0,
		"req_proc_time_hist_bucket_5":               0,
		"req_proc_time_hist_bucket_6":               0,
		"req_proc_time_hist_bucket_7":               0,
		"req_proc_time_hist_bucket_8":               0,
		"req_proc_time_hist_bucket_9":               1,
		"req_proc_time_hist_count":                  455,
		"req_proc_time_hist_sum":                    112018,
		"req_proc_time_max":                         498,
		"req_proc_time_min":                         2,
		"req_proc_time_sum":                         112018,
		"req_redirect":                              101,
		"req_ssl_cipher_suite_AES256-SHA":           127,
		"req_ssl_cipher_suite_DHE-RSA-AES256-SHA":   127,
		"req_ssl_cipher_suite_ECDHE-RSA-AES256-SHA": 105,
		"req_ssl_cipher_suite_PSK-RC4-SHA":          96,
		"req_ssl_proto_SSLv2":                       88,
		"req_ssl_proto_SSLv3":                       68,
		"req_ssl_proto_TLSv1":                       70,
		"req_ssl_proto_TLSv1.1":                     73,
		"req_ssl_proto_TLSv1.2":                     81,
		"req_ssl_proto_TLSv1.3":                     75,
		"req_success":                               300,
		"req_unmatched":                             45,
		"req_url_ptn_com":                           116,
		"req_url_ptn_net":                           107,
		"req_url_ptn_not_match":                     0,
		"req_url_ptn_org":                           109,
		"req_version_1.1":                           142,
		"req_version_2":                             159,
		"req_version_2.0":                           154,
		"req_vhost_198.51.100.1":                    96,
		"req_vhost_2001:db8:1ce::1":                 83,
		"req_vhost_localhost":                       86,
		"req_vhost_test.example.com":                97,
		"req_vhost_test.example.org":                93,
		"requests":                                  500,
		"resp_1xx":                                  123,
		"resp_2xx":                                  131,
		"resp_3xx":                                  101,
		"resp_4xx":                                  100,
		"resp_5xx":                                  0,
		"resp_code_100":                             68,
		"resp_code_101":                             55,
		"resp_code_200":                             66,
		"resp_code_201":                             65,
		"resp_code_300":                             51,
		"resp_code_301":                             50,
		"resp_code_400":                             54,
		"resp_code_401":                             46,
		"uniq_ipv4":                                 3,
		"uniq_ipv6":                                 2,
		"upstream_resp_time_avg":                    241,
		"upstream_resp_time_count":                  455,
		"upstream_resp_time_hist_bucket_1":          0,
		"upstream_resp_time_hist_bucket_10":         2,
		"upstream_resp_time_hist_bucket_11":         6,
		"upstream_resp_time_hist_bucket_2":          0,
		"upstream_resp_time_hist_bucket_3":          0,
		"upstream_resp_time_hist_bucket_4":          0,
		"upstream_resp_time_hist_bucket_5":          0,
		"upstream_resp_time_hist_bucket_6":          0,
		"upstream_resp_time_hist_bucket_7":          0,
		"upstream_resp_time_hist_bucket_8":          1,
		"upstream_resp_time_hist_bucket_9":          2,
		"upstream_resp_time_hist_count":             455,
		"upstream_resp_time_hist_sum":               109660,
		"upstream_resp_time_max":                    499,
		"upstream_resp_time_min":                    1,
		"upstream_resp_time_sum":                    109660,
		"url_ptn_com_bytes_received":                363260,
		"url_ptn_com_bytes_sent":                    349190,
		"url_ptn_com_req_proc_time_avg":             242,
		"url_ptn_com_req_proc_time_count":           116,
		"url_ptn_com_req_proc_time_max":             493,
		"url_ptn_com_req_proc_time_min":             8,
		"url_ptn_com_req_proc_time_sum":             28118,
		"url_ptn_com_resp_code_100":                 19,
		"url_ptn_com_resp_code_101":                 10,
		"url_ptn_com_resp_code_200":                 16,
		"url_ptn_com_resp_code_201":                 18,
		"url_ptn_com_resp_code_300":                 13,
		"url_ptn_com_resp_code_301":                 11,
		"url_ptn_com_resp_code_400":                 15,
		"url_ptn_com_resp_code_401":                 14,
		"url_ptn_net_bytes_received":                331553,
		"url_ptn_net_bytes_sent":                    316867,
		"url_ptn_net_req_proc_time_avg":             261,
		"url_ptn_net_req_proc_time_count":           107,
		"url_ptn_net_req_proc_time_max":             498,
		"url_ptn_net_req_proc_time_min":             11,
		"url_ptn_net_req_proc_time_sum":             28006,
		"url_ptn_net_resp_code_100":                 12,
		"url_ptn_net_resp_code_101":                 12,
		"url_ptn_net_resp_code_200":                 16,
		"url_ptn_net_resp_code_201":                 18,
		"url_ptn_net_resp_code_300":                 12,
		"url_ptn_net_resp_code_301":                 13,
		"url_ptn_net_resp_code_400":                 12,
		"url_ptn_net_resp_code_401":                 12,
		"url_ptn_not_match_bytes_received":          0,
		"url_ptn_not_match_bytes_sent":              0,
		"url_ptn_not_match_req_proc_time_avg":       0,
		"url_ptn_not_match_req_proc_time_count":     0,
		"url_ptn_not_match_req_proc_time_max":       0,
		"url_ptn_not_match_req_proc_time_min":       0,
		"url_ptn_not_match_req_proc_time_sum":       0,
		"url_ptn_org_bytes_received":                323159,
		"url_ptn_org_bytes_sent":                    322951,
		"url_ptn_org_req_proc_time_avg":             237,
		"url_ptn_org_req_proc_time_count":           109,
		"url_ptn_org_req_proc_time_max":             498,
		"url_ptn_org_req_proc_time_min":             4,
		"url_ptn_org_req_proc_time_sum":             25884,
		"url_ptn_org_resp_code_100":                 16,
		"url_ptn_org_resp_code_101":                 14,
		"url_ptn_org_resp_code_200":                 15,
		"url_ptn_org_resp_code_201":                 12,
		"url_ptn_org_resp_code_300":                 14,
		"url_ptn_org_resp_code_301":                 15,
		"url_ptn_org_resp_code_400":                 14,
		"url_ptn_org_resp_code_401":                 9,
	}

	assert.Equal(t, expected, weblog.Collect())
	testCharts(t, weblog)
}

func TestWebLog_Collect_CommonLogFormat(t *testing.T) {
	weblog := prepareWebLogCollectCommon(t)

	expected := map[string]int64{
		"bytes_received":                    1426307,
		"bytes_sent":                        0,
		"req_bad":                           66,
		"req_error":                         0,
		"req_http_scheme":                   0,
		"req_https_scheme":                  0,
		"req_ipv4":                          265,
		"req_ipv6":                          192,
		"req_method_GET":                    132,
		"req_method_HEAD":                   159,
		"req_method_POST":                   166,
		"req_proc_time_avg":                 0,
		"req_proc_time_count":               0,
		"req_proc_time_hist_bucket_1":       0,
		"req_proc_time_hist_bucket_10":      0,
		"req_proc_time_hist_bucket_11":      0,
		"req_proc_time_hist_bucket_2":       0,
		"req_proc_time_hist_bucket_3":       0,
		"req_proc_time_hist_bucket_4":       0,
		"req_proc_time_hist_bucket_5":       0,
		"req_proc_time_hist_bucket_6":       0,
		"req_proc_time_hist_bucket_7":       0,
		"req_proc_time_hist_bucket_8":       0,
		"req_proc_time_hist_bucket_9":       0,
		"req_proc_time_hist_count":          0,
		"req_proc_time_hist_sum":            0,
		"req_proc_time_max":                 0,
		"req_proc_time_min":                 0,
		"req_proc_time_sum":                 0,
		"req_redirect":                      97,
		"req_success":                       294,
		"req_unmatched":                     43,
		"req_version_1.1":                   157,
		"req_version_2":                     144,
		"req_version_2.0":                   156,
		"requests":                          500,
		"resp_1xx":                          130,
		"resp_2xx":                          103,
		"resp_3xx":                          97,
		"resp_4xx":                          127,
		"resp_5xx":                          0,
		"resp_code_100":                     66,
		"resp_code_101":                     64,
		"resp_code_200":                     50,
		"resp_code_201":                     53,
		"resp_code_300":                     47,
		"resp_code_301":                     50,
		"resp_code_400":                     66,
		"resp_code_401":                     61,
		"uniq_ipv4":                         3,
		"uniq_ipv6":                         2,
		"upstream_resp_time_avg":            0,
		"upstream_resp_time_count":          0,
		"upstream_resp_time_hist_bucket_1":  0,
		"upstream_resp_time_hist_bucket_10": 0,
		"upstream_resp_time_hist_bucket_11": 0,
		"upstream_resp_time_hist_bucket_2":  0,
		"upstream_resp_time_hist_bucket_3":  0,
		"upstream_resp_time_hist_bucket_4":  0,
		"upstream_resp_time_hist_bucket_5":  0,
		"upstream_resp_time_hist_bucket_6":  0,
		"upstream_resp_time_hist_bucket_7":  0,
		"upstream_resp_time_hist_bucket_8":  0,
		"upstream_resp_time_hist_bucket_9":  0,
		"upstream_resp_time_hist_count":     0,
		"upstream_resp_time_hist_sum":       0,
		"upstream_resp_time_max":            0,
		"upstream_resp_time_min":            0,
		"upstream_resp_time_sum":            0,
	}

	assert.Equal(t, expected, weblog.Collect())
	testCharts(t, weblog)
}

func TestWebLog_Collect_CustomLogs(t *testing.T) {
	weblog := prepareWebLogCollectCustom(t)

	expected := map[string]int64{
		"bytes_received":                    0,
		"bytes_sent":                        0,
		"custom_field_drink_beer":           52,
		"custom_field_drink_wine":           40,
		"custom_field_side_dark":            46,
		"custom_field_side_light":           46,
		"req_bad":                           0,
		"req_error":                         0,
		"req_http_scheme":                   0,
		"req_https_scheme":                  0,
		"req_ipv4":                          0,
		"req_ipv6":                          0,
		"req_proc_time_avg":                 0,
		"req_proc_time_count":               0,
		"req_proc_time_hist_bucket_1":       0,
		"req_proc_time_hist_bucket_10":      0,
		"req_proc_time_hist_bucket_11":      0,
		"req_proc_time_hist_bucket_2":       0,
		"req_proc_time_hist_bucket_3":       0,
		"req_proc_time_hist_bucket_4":       0,
		"req_proc_time_hist_bucket_5":       0,
		"req_proc_time_hist_bucket_6":       0,
		"req_proc_time_hist_bucket_7":       0,
		"req_proc_time_hist_bucket_8":       0,
		"req_proc_time_hist_bucket_9":       0,
		"req_proc_time_hist_count":          0,
		"req_proc_time_hist_sum":            0,
		"req_proc_time_max":                 0,
		"req_proc_time_min":                 0,
		"req_proc_time_sum":                 0,
		"req_redirect":                      0,
		"req_success":                       0,
		"req_unmatched":                     8,
		"requests":                          100,
		"resp_1xx":                          0,
		"resp_2xx":                          0,
		"resp_3xx":                          0,
		"resp_4xx":                          0,
		"resp_5xx":                          0,
		"uniq_ipv4":                         0,
		"uniq_ipv6":                         0,
		"upstream_resp_time_avg":            0,
		"upstream_resp_time_count":          0,
		"upstream_resp_time_hist_bucket_1":  0,
		"upstream_resp_time_hist_bucket_10": 0,
		"upstream_resp_time_hist_bucket_11": 0,
		"upstream_resp_time_hist_bucket_2":  0,
		"upstream_resp_time_hist_bucket_3":  0,
		"upstream_resp_time_hist_bucket_4":  0,
		"upstream_resp_time_hist_bucket_5":  0,
		"upstream_resp_time_hist_bucket_6":  0,
		"upstream_resp_time_hist_bucket_7":  0,
		"upstream_resp_time_hist_bucket_8":  0,
		"upstream_resp_time_hist_bucket_9":  0,
		"upstream_resp_time_hist_count":     0,
		"upstream_resp_time_hist_sum":       0,
		"upstream_resp_time_max":            0,
		"upstream_resp_time_min":            0,
		"upstream_resp_time_sum":            0,
	}

	assert.Equal(t, expected, weblog.Collect())
	testCharts(t, weblog)
}

func testCharts(t *testing.T, w *WebLog) {
	testRespCodeChart(t, w)
	testReqVhostChart(t, w)
	testReqPortChart(t, w)
	testReqSchemeChart(t, w)
	testReqHTTPMethodChart(t, w)
	testReqHTTPVersionChart(t, w)
	testReqClientCharts(t, w)
	testBandwidthChart(t, w)
	testReqURLPatternChart(t, w)
	testURLPatternStatsCharts(t, w)
	testReqProcTimeCharts(t, w)
	testUpsRespTimeCharts(t, w)
	testSSLProtoChart(t, w)
	testSSLCipherSuiteChart(t, w)
	testReqCustomFieldCharts(t, w)
}

func testReqProcTimeCharts(t *testing.T, w *WebLog) {
	if isEmptySummary(w.mx.ReqProcTime) {
		assert.Falsef(t, w.Charts().Has(reqProcTime.ID), "chart '%s' is created", reqProcTime.ID)
	} else {
		assert.Truef(t, w.Charts().Has(reqProcTime.ID), "chart '%s' is not created", reqProcTime.ID)
	}

	if isEmptyHistogram(w.mx.ReqProcTimeHist) {
		assert.Falsef(t, w.Charts().Has(reqProcTimeHist.ID), "chart '%s' is created", reqProcTimeHist.ID)
	} else {
		assert.Truef(t, w.Charts().Has(reqProcTimeHist.ID), "chart '%s' is not created", reqProcTimeHist.ID)
	}
}

func testUpsRespTimeCharts(t *testing.T, w *WebLog) {
	if isEmptySummary(w.mx.UpsRespTime) {
		assert.Falsef(t, w.Charts().Has(upsRespTime.ID), "chart '%s' is created", upsRespTime.ID)
	} else {
		assert.Truef(t, w.Charts().Has(upsRespTime.ID), "chart '%s' is not created", upsRespTime.ID)
	}

	if isEmptyHistogram(w.mx.UpsRespTimeHist) {
		assert.Falsef(t, w.Charts().Has(upsRespTimeHist.ID), "chart '%s' is created", upsRespTimeHist.ID)
	} else {
		assert.Truef(t, w.Charts().Has(upsRespTimeHist.ID), "chart '%s' is not created", upsRespTimeHist.ID)
	}
}

func testReqVhostChart(t *testing.T, w *WebLog) {
	if len(w.mx.ReqVhost) == 0 {
		assert.Falsef(t, w.Charts().Has(reqByVhost.ID), "chart '%s' is created", reqByVhost.ID)
		return
	}

	chart := w.Charts().Get(reqByVhost.ID)
	assert.NotNilf(t, chart, "chart '%s' is not created", reqByVhost.ID)
	if chart == nil {
		return
	}
	for v := range w.mx.ReqVhost {
		id := "req_vhost_" + v
		assert.Truef(t, chart.HasDim(id), "chart '%s' has no dim for '%s' vhost, expected '%s'", chart.ID, v, id)
	}
}

func testReqPortChart(t *testing.T, w *WebLog) {
	if len(w.mx.ReqPort) == 0 {
		assert.Falsef(t, w.Charts().Has(reqByPort.ID), "chart '%s' is created", reqByPort.ID)
		return
	}

	chart := w.Charts().Get(reqByPort.ID)
	assert.NotNilf(t, chart, "chart '%s' is not created", reqByPort.ID)
	if chart == nil {
		return
	}
	for v := range w.mx.ReqPort {
		id := "req_port_" + v
		assert.Truef(t, chart.HasDim(id), "chart '%s' has no dim for '%s' port, expected '%s'", chart.ID, v, id)
	}
}

func testReqHTTPMethodChart(t *testing.T, w *WebLog) {
	if len(w.mx.ReqMethod) == 0 {
		assert.Falsef(t, w.Charts().Has(reqByMethod.ID), "chart '%s' is created", reqByMethod.ID)
		return
	}

	chart := w.Charts().Get(reqByMethod.ID)
	assert.NotNilf(t, chart, "chart '%s' is not created", reqByMethod.ID)
	if chart == nil {
		return
	}
	for v := range w.mx.ReqMethod {
		id := "req_method_" + v
		assert.Truef(t, chart.HasDim(id), "chart '%s' has no dim for '%s' method, expected '%s'", chart.ID, v, id)
	}
}

func testReqHTTPVersionChart(t *testing.T, w *WebLog) {
	if len(w.mx.ReqVersion) == 0 {
		assert.Falsef(t, w.Charts().Has(reqByVersion.ID), "chart '%s' is created", reqByVersion.ID)
		return
	}

	chart := w.Charts().Get(reqByVersion.ID)
	assert.NotNilf(t, chart, "chart '%s' is not created", reqByVersion.ID)
	if chart == nil {
		return
	}
	for v := range w.mx.ReqVersion {
		id := "req_version_" + v
		assert.Truef(t, chart.HasDim(id), "chart '%s' has no dim for '%s' version, expected '%s'", chart.ID, v, id)
	}
}

func testReqSchemeChart(t *testing.T, w *WebLog) {
	if w.mx.ReqHTTPScheme.Value() == 0 && w.mx.ReqHTTPScheme.Value() == 0 {
		assert.Falsef(t, w.Charts().Has(reqByScheme.ID), "chart '%s' is created", reqByScheme.ID)
	} else {
		assert.Truef(t, w.Charts().Has(reqByScheme.ID), "chart '%s' is not created", reqByScheme.ID)
	}
}

func testReqClientCharts(t *testing.T, w *WebLog) {
	if w.mx.ReqIPv4.Value() == 0 && w.mx.ReqIPv6.Value() == 0 {
		assert.Falsef(t, w.Charts().Has(reqByIPProto.ID), "chart '%s' is created", reqByIPProto.ID)
	} else {
		assert.Truef(t, w.Charts().Has(reqByIPProto.ID), "chart '%s' is not created", reqByIPProto.ID)
	}

	if w.mx.UniqueIPv4.Value() == 0 && w.mx.UniqueIPv6.Value() == 0 {
		assert.Falsef(t, w.Charts().Has(uniqIPsCurPoll.ID), "chart '%s' is created", uniqIPsCurPoll.ID)
	} else {
		assert.Truef(t, w.Charts().Has(uniqIPsCurPoll.ID), "chart '%s' is not created", uniqIPsCurPoll.ID)
	}
}

func testBandwidthChart(t *testing.T, w *WebLog) {
	if w.mx.BytesSent.Value() == 0 && w.mx.BytesReceived.Value() == 0 {
		assert.Falsef(t, w.Charts().Has(bandwidth.ID), "chart '%s' is created", bandwidth.ID)
	} else {
		assert.Truef(t, w.Charts().Has(bandwidth.ID), "chart '%s' is not created", bandwidth.ID)
	}
}

func testReqURLPatternChart(t *testing.T, w *WebLog) {
	if isEmptyCounterVec(w.mx.ReqURLPattern) {
		assert.Falsef(t, w.Charts().Has(reqByURLPattern.ID), "chart '%s' is created", reqByURLPattern.ID)
		return
	}

	chart := w.Charts().Get(reqByURLPattern.ID)
	assert.NotNilf(t, chart, "chart '%s' is not created", reqByURLPattern.ID)
	if chart == nil {
		return
	}
	for v := range w.mx.ReqURLPattern {
		id := "req_url_ptn_" + v
		assert.True(t, chart.HasDim(id), "chart '%s' has no dim for '%s' pattern, expected '%s'", chart.ID, v, id)
	}
}

func testSSLProtoChart(t *testing.T, w *WebLog) {
	if len(w.mx.ReqSSLProto) == 0 {
		assert.Falsef(t, w.Charts().Has(reqBySSLProto.ID), "chart '%s' is created", reqBySSLProto.ID)
		return
	}

	chart := w.Charts().Get(reqBySSLProto.ID)
	assert.NotNilf(t, chart, "chart '%s' is not created", reqBySSLProto.ID)
	if chart == nil {
		return
	}
	for v := range w.mx.ReqSSLProto {
		id := "req_ssl_proto_" + v
		assert.Truef(t, chart.HasDim(id), "chart '%s' has no dim for '%s' ssl proto, expected '%s'", chart.ID, v, id)
	}
}

func testSSLCipherSuiteChart(t *testing.T, w *WebLog) {
	if len(w.mx.ReqSSLCipherSuite) == 0 {
		assert.Falsef(t, w.Charts().Has(reqBySSLCipherSuite.ID), "chart '%s' is created", reqBySSLCipherSuite.ID)
		return
	}

	chart := w.Charts().Get(reqBySSLCipherSuite.ID)
	assert.NotNilf(t, chart, "chart '%s' is not created", reqBySSLCipherSuite.ID)
	if chart == nil {
		return
	}
	for v := range w.mx.ReqSSLCipherSuite {
		id := "req_ssl_cipher_suite_" + v
		assert.Truef(t, chart.HasDim(id), "chart '%s' has no dim for '%s' ssl cipher suite, expected '%s'", chart.ID, v, id)
	}
}

func testReqCustomFieldCharts(t *testing.T, w *WebLog) {
	for _, cf := range w.CustomFields {
		var id string
		if w.customLog {
			id = fmt.Sprintf(matchesByCustomFieldPattern.ID, cf.Name)
		} else {
			id = fmt.Sprintf(reqByCustomFieldPattern.ID, cf.Name)
		}
		chart := w.Charts().Get(id)
		assert.NotNilf(t, chart, "chart '%s' is not created", id)
		if chart == nil {
			continue
		}

		for _, p := range cf.Patterns {
			id := fmt.Sprintf("custom_field_%s_%s", cf.Name, p.Name)
			assert.True(t, chart.HasDim(id), "chart '%s' has no dim for '%s' pattern, expected '%s'", chart.ID, p, id)
		}
	}
}

func testURLPatternStatsCharts(t *testing.T, w *WebLog) {
	for _, p := range w.URLPatterns {
		chartID := fmt.Sprintf(urlPatternRespCodes.ID, p.Name)
		chart := w.Charts().Get(chartID)
		assert.NotNilf(t, chart, "chart '%s' is not created", chartID)
		if chart == nil {
			continue
		}

		stats, ok := w.mx.URLPatternStats[p.Name]
		assert.Truef(t, ok, "url pattern '%s' has no metric in w.mx.URLPatternStats", p.Name)
		if !ok {
			continue
		}
		for v := range stats.RespCode {
			id := fmt.Sprintf("url_ptn_%s_resp_code_%s", p.Name, v)
			assert.Truef(t, chart.HasDim(id), "chart '%s' has no dim for '%s' code, expected '%s'", chartID, v, id)
		}
	}

	for _, p := range w.URLPatterns {
		id := fmt.Sprintf(urlPatternBandwidth.ID, p.Name)
		if w.mx.BytesSent.Value() == 0 && w.mx.BytesReceived.Value() == 0 {
			assert.Falsef(t, w.Charts().Has(id), "chart '%s' is created", id)
		} else {
			assert.Truef(t, w.Charts().Has(id), "chart '%s' is not created", id)
		}
	}

	for _, p := range w.URLPatterns {
		id := fmt.Sprintf(urlPatternReqProcTime.ID, p.Name)
		if isEmptySummary(w.mx.ReqProcTime) {
			assert.Falsef(t, w.Charts().Has(id), "chart '%s' is created", id)
		} else {
			assert.Truef(t, w.Charts().Has(id), "chart '%s' is not created", id)
		}
	}
}

func testRespCodeChart(t *testing.T, w *WebLog) {
	if isEmptyCounterVec(w.mx.RespCode) {
		if !w.GroupRespCodes {
			chart := w.Charts().Get(respCodes.ID)
			assert.Nilf(t, chart, "chart '%s' is not created", respCodes.ID)
			return
		}
	}

	if !w.GroupRespCodes {
		chart := w.Charts().Get(respCodes.ID)
		assert.NotNilf(t, chart, "chart '%s' is not created", respCodes.ID)
		if chart == nil {
			return
		}
		for v := range w.mx.RespCode {
			id := "resp_code_" + v
			assert.Truef(t, chart.HasDim(id), "chart '%s' has no dim for '%s' code, expected '%s'", chart.ID, v, id)
		}
		return
	}

	findCodes := func(class string) (codes []string) {
		for v := range w.mx.RespCode {
			if v[:1] == class {
				codes = append(codes, v)
			}
		}
		return codes
	}

	var n int
	ids := []string{
		respCodes1xx.ID,
		respCodes2xx.ID,
		respCodes3xx.ID,
		respCodes4xx.ID,
		respCodes5xx.ID,
	}
	for i, chartID := range ids {
		class := strconv.Itoa(i + 1)
		codes := findCodes(class)
		n += len(codes)
		chart := w.Charts().Get(chartID)
		assert.NotNilf(t, chart, "chart '%s' is not created", chartID)
		if chart == nil {
			return
		}
		for _, v := range codes {
			id := "resp_code_" + v
			assert.Truef(t, chart.HasDim(id), "chart '%s' has no dim for '%s' code, expected '%s'", chartID, v, id)
		}
	}
	assert.Equal(t, len(w.mx.RespCode), n)
}

var (
	emptySummary   = newWebLogSummary()
	emptyHistogram = metrics.NewHistogram(metrics.DefBuckets)
)

func isEmptySummary(s metrics.Summary) bool     { return reflect.DeepEqual(s, emptySummary) }
func isEmptyHistogram(h metrics.Histogram) bool { return reflect.DeepEqual(h, emptyHistogram) }

func isEmptyCounterVec(cv metrics.CounterVec) bool {
	for _, c := range cv {
		if c.Value() > 0 {
			return false
		}
	}
	return true
}

func prepareWebLogCollectFull(t *testing.T) *WebLog {
	t.Helper()
	format := strings.Join([]string{
		"$host:$server_port",
		"$scheme",
		"$remote_addr",
		"$ssl_protocol",
		"$ssl_cipher",
		`"$request"`,
		"$status",
		"$body_bytes_sent",
		"$request_length",
		"$request_time",
		"$upstream_response_time",
		"$side",
		"$drink",
	}, " ")

	cfg := Config{
		Parser: logs.ParserConfig{
			LogType: logs.TypeCSV,
			CSV: logs.CSVConfig{
				FieldsPerRecord:  -1,
				Delimiter:        ' ',
				TrimLeadingSpace: false,
				Format:           format,
				CheckField:       checkCSVFormatField,
			},
		},
		Path:        "testdata/full.log",
		ExcludePath: "",
		URLPatterns: []userPattern{
			{Name: "com", Match: "~ com$"},
			{Name: "org", Match: "~ org$"},
			{Name: "net", Match: "~ net$"},
			{Name: "not_match", Match: "* !*"},
		},
		CustomFields: []customField{
			{
				Name: "side",
				Patterns: []userPattern{
					{Name: "dark", Match: "= dark"},
					{Name: "light", Match: "= light"},
				},
			},
			{
				Name: "drink",
				Patterns: []userPattern{
					{Name: "beer", Match: "= beer"},
					{Name: "wine", Match: "= wine"},
				},
			},
		},
		Histogram:      metrics.DefBuckets,
		GroupRespCodes: true,
	}
	weblog := New()
	weblog.Config = cfg
	require.True(t, weblog.Init())
	require.True(t, weblog.Check())
	defer weblog.Cleanup()

	p, err := logs.NewCSVParser(weblog.Parser.CSV, bytes.NewReader(testFullLog))
	require.NoError(t, err)
	weblog.parser = p
	return weblog
}

func prepareWebLogCollectCommon(t *testing.T) *WebLog {
	t.Helper()
	format := strings.Join([]string{
		"$remote_addr",
		`"$request"`,
		"$status",
		"$body_bytes_sent",
	}, " ")

	cfg := Config{
		Parser: logs.ParserConfig{
			LogType: logs.TypeCSV,
			CSV: logs.CSVConfig{
				FieldsPerRecord:  -1,
				Delimiter:        ' ',
				TrimLeadingSpace: false,
				Format:           format,
				CheckField:       checkCSVFormatField,
			},
		},
		Path:           "testdata/common.log",
		ExcludePath:    "",
		URLPatterns:    nil,
		CustomFields:   nil,
		Histogram:      nil,
		GroupRespCodes: false,
	}

	weblog := New()
	weblog.Config = cfg
	require.True(t, weblog.Init())
	require.True(t, weblog.Check())
	defer weblog.Cleanup()

	p, err := logs.NewCSVParser(weblog.Parser.CSV, bytes.NewReader(testCommonLog))
	require.NoError(t, err)
	weblog.parser = p
	return weblog
}

func prepareWebLogCollectCustom(t *testing.T) *WebLog {
	t.Helper()
	format := strings.Join([]string{
		"$side",
		"$drink",
	}, " ")

	cfg := Config{
		Parser: logs.ParserConfig{
			LogType: logs.TypeCSV,
			CSV: logs.CSVConfig{
				FieldsPerRecord:  2,
				Delimiter:        ' ',
				TrimLeadingSpace: false,
				Format:           format,
				CheckField:       checkCSVFormatField,
			},
		},
		CustomFields: []customField{
			{
				Name: "side",
				Patterns: []userPattern{
					{Name: "dark", Match: "= dark"},
					{Name: "light", Match: "= light"},
				},
			},
			{
				Name: "drink",
				Patterns: []userPattern{
					{Name: "beer", Match: "= beer"},
					{Name: "wine", Match: "= wine"},
				},
			},
		},
		Path:           "testdata/custom.log",
		ExcludePath:    "",
		URLPatterns:    nil,
		Histogram:      nil,
		GroupRespCodes: false,
	}
	weblog := New()
	weblog.Config = cfg
	require.True(t, weblog.Init())
	require.True(t, weblog.Check())
	defer weblog.Cleanup()

	p, err := logs.NewCSVParser(weblog.Parser.CSV, bytes.NewReader(testCustomLog))
	require.NoError(t, err)
	weblog.parser = p
	return weblog
}

// generateLogs is used to populate 'testdata/full.log'
func generateLogs(w io.Writer, num int) error {
	var (
		vhost     = []string{"localhost", "test.example.com", "test.example.org", "198.51.100.1", "2001:db8:1ce::1"}
		scheme    = []string{"http", "https"}
		client    = []string{"localhost", "203.0.113.1", "203.0.113.2", "2001:db8:2ce:1", "2001:db8:2ce:2"}
		method    = []string{"GET", "HEAD", "POST"}
		url       = []string{"invalid.example", "example.com", "example.org", "example.net"}
		version   = []string{"1.1", "2", "2.0"}
		status    = []int{100, 101, 200, 201, 300, 301, 400, 401} // no 5xx on purpose
		sslProto  = []string{"TLSv1", "TLSv1.1", "TLSv1.2", "TLSv1.3", "SSLv2", "SSLv3"}
		sslCipher = []string{"ECDHE-RSA-AES256-SHA", "DHE-RSA-AES256-SHA", "AES256-SHA", "PSK-RC4-SHA"}

		customField1 = []string{"dark", "light"}
		customField2 = []string{"beer", "wine"}
	)

	var line string
	for i := 0; i < num; i++ {
		unmatched := randInt(1, 100) > 90
		if unmatched {
			line = "Unmatched! The rat the cat the dog chased killed ate the malt!\n"
		} else {
			// test.example.com:80 http 203.0.113.1 TLSv1 AES256-SHA "GET / HTTP/1.1" 200 1674 2674 3674 4674 dark beer
			line = fmt.Sprintf("%s:%d %s %s %s %s \"%s /%s HTTP/%s\" %d %d %d %d %d %s %s\n",
				randFromString(vhost),
				randInt(80, 85),
				randFromString(scheme),
				randFromString(client),
				randFromString(sslProto),
				randFromString(sslCipher),
				randFromString(method),
				randFromString(url),
				randFromString(version),
				randFromInt(status),
				randInt(1000, 5000),
				randInt(1000, 5000),
				randInt(1, 500),
				randInt(1, 500),
				randFromString(customField1),
				randFromString(customField2),
			)
		}
		_, err := fmt.Fprint(w, line)
		if err != nil {
			return err
		}
	}
	return nil
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func randFromString(s []string) string { return s[r.Intn(len(s))] }
func randFromInt(s []int) int          { return s[r.Intn(len(s))] }
func randInt(min, max int) int         { return r.Intn(max-min) + min }

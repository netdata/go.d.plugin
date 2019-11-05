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
	"github.com/netdata/go.d.plugin/pkg/matcher"
	"github.com/netdata/go.d.plugin/pkg/metrics"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testFormat = []string{
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
		"$custom",
	}
	testConfig = Config{
		Parser: logs.ParserConfig{
			LogType: logs.TypeCSV,
			CSV: logs.CSVConfig{
				Delimiter:        ' ',
				TrimLeadingSpace: false,
				Format:           strings.Join(testFormat, " "),
				CheckField:       checkCSVFormatField,
			},
		},
		Path:        "testdata/full.log",
		ExcludePath: "",
		Filter:      matcher.SimpleExpr{Excludes: []string{"~ ^/invalid"}},
		URLPatterns: []userPattern{
			{Name: "com", Match: "~ com$"},
			{Name: "org", Match: "~ org$"},
			{Name: "net", Match: "~ net$"},
			{Name: "not_match", Match: "= not_match"},
		},
		CustomPatterns: []userPattern{
			{Name: "dark", Match: "~ dark$"},
			{Name: "light", Match: "~ light$"},
			{Name: "not_match", Match: "= not_match"},
		},
		Histogram:      metrics.DefBuckets,
		GroupRespCodes: true,
	}
)

var (
	testFullLog, _ = ioutil.ReadFile("testdata/full.log")
)

func Test_readTestData(t *testing.T) {
	assert.NotNil(t, testFullLog)
}

func TestWebLog_Init(t *testing.T) {
}

func TestWebLog_Collect(t *testing.T) {
	weblog := New()
	weblog.Config = testConfig
	require.True(t, weblog.Init())

	p, err := logs.NewCSVParser(testConfig.Parser.CSV, bytes.NewReader(testFullLog))
	require.NoError(t, err)
	weblog.parser = p
	weblog.line = newEmptyLogLine()

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
		"bytes_received":                            1174596,
		"bytes_sent":                                1154933,
		"req_custom_ptn_dark":                       173,
		"req_custom_ptn_light":                      209,
		"req_custom_ptn_not_match":                  0,
		"req_filtered":                              118,
		"req_http_scheme":                           199,
		"req_https_scheme":                          183,
		"req_ipv4":                                  218,
		"req_ipv6":                                  164,
		"req_method_GET":                            125,
		"req_method_HEAD":                           118,
		"req_method_POST":                           139,
		"req_port_80":                               85,
		"req_port_81":                               74,
		"req_port_82":                               74,
		"req_port_83":                               80,
		"req_port_84":                               69,
		"req_proc_time_avg":                         261,
		"req_proc_time_count":                       382,
		"req_proc_time_hist_bucket_1":               0,
		"req_proc_time_hist_bucket_10":              4,
		"req_proc_time_hist_bucket_11":              5,
		"req_proc_time_hist_bucket_2":               0,
		"req_proc_time_hist_bucket_3":               0,
		"req_proc_time_hist_bucket_4":               0,
		"req_proc_time_hist_bucket_5":               0,
		"req_proc_time_hist_bucket_6":               0,
		"req_proc_time_hist_bucket_7":               0,
		"req_proc_time_hist_bucket_8":               1,
		"req_proc_time_hist_bucket_9":               2,
		"req_proc_time_hist_count":                  382,
		"req_proc_time_hist_sum":                    99779,
		"req_proc_time_max":                         498,
		"req_proc_time_min":                         1,
		"req_proc_time_sum":                         99779,
		"req_ssl_cipher_suite_AES256-SHA":           103,
		"req_ssl_cipher_suite_DHE-RSA-AES256-SHA":   97,
		"req_ssl_cipher_suite_ECDHE-RSA-AES256-SHA": 74,
		"req_ssl_cipher_suite_PSK-RC4-SHA":          108,
		"req_ssl_proto_SSLv2":                       66,
		"req_ssl_proto_SSLv3":                       59,
		"req_ssl_proto_TLSv1":                       58,
		"req_ssl_proto_TLSv1.1":                     63,
		"req_ssl_proto_TLSv1.2":                     59,
		"req_ssl_proto_TLSv1.3":                     77,
		"req_unmatched":                             50,
		"req_url_ptn_com":                           126,
		"req_url_ptn_net":                           137,
		"req_url_ptn_not_match":                     0,
		"req_url_ptn_org":                           119,
		"req_version_1.1":                           138,
		"req_version_2":                             122,
		"req_version_2.0":                           122,
		"req_vhost_198.51.100.1":                    67,
		"req_vhost_2001:db8:1ce::1":                 87,
		"req_vhost_localhost":                       66,
		"req_vhost_test.example.com":                77,
		"req_vhost_test.example.org":                85,
		"requests":                                  550,
		"resp_1xx":                                  85,
		"resp_2xx":                                  96,
		"resp_3xx":                                  107,
		"resp_4xx":                                  94,
		"resp_5xx":                                  0,
		"resp_client_error":                         94,
		"resp_redirect":                             107,
		"resp_server_error":                         0,
		"resp_status_code_100":                      39,
		"resp_status_code_101":                      46,
		"resp_status_code_200":                      57,
		"resp_status_code_201":                      39,
		"resp_status_code_300":                      57,
		"resp_status_code_301":                      50,
		"resp_status_code_400":                      52,
		"resp_status_code_401":                      42,
		"resp_successful":                           181,
		"uniq_ipv4":                                 3,
		"uniq_ipv6":                                 2,
		"upstream_resp_time_avg":                    247,
		"upstream_resp_time_count":                  382,
		"upstream_resp_time_hist_bucket_1":          0,
		"upstream_resp_time_hist_bucket_10":         2,
		"upstream_resp_time_hist_bucket_11":         6,
		"upstream_resp_time_hist_bucket_2":          0,
		"upstream_resp_time_hist_bucket_3":          0,
		"upstream_resp_time_hist_bucket_4":          0,
		"upstream_resp_time_hist_bucket_5":          0,
		"upstream_resp_time_hist_bucket_6":          0,
		"upstream_resp_time_hist_bucket_7":          0,
		"upstream_resp_time_hist_bucket_8":          0,
		"upstream_resp_time_hist_bucket_9":          0,
		"upstream_resp_time_hist_count":             382,
		"upstream_resp_time_hist_sum":               94414,
		"upstream_resp_time_max":                    498,
		"upstream_resp_time_min":                    3,
		"upstream_resp_time_sum":                    94414,
		"url_ptn_com_bytes_received":                382189,
		"url_ptn_com_bytes_sent":                    373539,
		"url_ptn_com_req_proc_time_avg":             255,
		"url_ptn_com_req_proc_time_count":           126,
		"url_ptn_com_req_proc_time_max":             495,
		"url_ptn_com_req_proc_time_min":             2,
		"url_ptn_com_req_proc_time_sum":             32164,
		"url_ptn_com_resp_status_code_100":          12,
		"url_ptn_com_resp_status_code_101":          15,
		"url_ptn_com_resp_status_code_200":          19,
		"url_ptn_com_resp_status_code_201":          16,
		"url_ptn_com_resp_status_code_300":          19,
		"url_ptn_com_resp_status_code_301":          19,
		"url_ptn_com_resp_status_code_400":          14,
		"url_ptn_com_resp_status_code_401":          12,
		"url_ptn_net_bytes_received":                431372,
		"url_ptn_net_bytes_sent":                    415012,
		"url_ptn_net_req_proc_time_avg":             258,
		"url_ptn_net_req_proc_time_count":           137,
		"url_ptn_net_req_proc_time_max":             498,
		"url_ptn_net_req_proc_time_min":             4,
		"url_ptn_net_req_proc_time_sum":             35414,
		"url_ptn_net_resp_status_code_100":          17,
		"url_ptn_net_resp_status_code_101":          21,
		"url_ptn_net_resp_status_code_200":          21,
		"url_ptn_net_resp_status_code_201":          9,
		"url_ptn_net_resp_status_code_300":          17,
		"url_ptn_net_resp_status_code_301":          20,
		"url_ptn_net_resp_status_code_400":          18,
		"url_ptn_net_resp_status_code_401":          14,
		"url_ptn_not_match_bytes_received":          0,
		"url_ptn_not_match_bytes_sent":              0,
		"url_ptn_not_match_req_proc_time_avg":       0,
		"url_ptn_not_match_req_proc_time_count":     0,
		"url_ptn_not_match_req_proc_time_max":       0,
		"url_ptn_not_match_req_proc_time_min":       0,
		"url_ptn_not_match_req_proc_time_sum":       0,
		"url_ptn_org_bytes_received":                361035,
		"url_ptn_org_bytes_sent":                    366382,
		"url_ptn_org_req_proc_time_avg":             270,
		"url_ptn_org_req_proc_time_count":           119,
		"url_ptn_org_req_proc_time_max":             498,
		"url_ptn_org_req_proc_time_min":             1,
		"url_ptn_org_req_proc_time_sum":             32201,
		"url_ptn_org_resp_status_code_100":          10,
		"url_ptn_org_resp_status_code_101":          10,
		"url_ptn_org_resp_status_code_200":          17,
		"url_ptn_org_resp_status_code_201":          14,
		"url_ptn_org_resp_status_code_300":          21,
		"url_ptn_org_resp_status_code_301":          11,
		"url_ptn_org_resp_status_code_400":          20,
		"url_ptn_org_resp_status_code_401":          16,
	}
	_ = expected

	assert.Equal(t, expected, weblog.Collect())
	testCharts(t, weblog)
}

func testCharts(t *testing.T, w *WebLog) {
	testRespStatusCodeChart(t, w)
	testReqVhostChart(t, w)
	testReqPortChart(t, w)
	testReqSchemeChart(t, w)
	testReqHTTPMethodChart(t, w)
	testReqHTTPVersionChart(t, w)
	testReqClientCharts(t, w)
	testBandwidthChart(t, w)
	testReqURLPatternChart(t, w)
	testReqCustomPatternChart(t, w)
	testURLPatternStatsCharts(t, w)
	testReqProcTimeCharts(t, w)
	testUpsRespTimeCharts(t, w)
	testSSLProtoChart(t, w)
	testSSLCipherSuiteChart(t, w)
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
		assert.Falsef(t, w.Charts().Has(reqPerVhost.ID), "chart '%s' is created", reqPerVhost.ID)
		return
	}

	chart := w.Charts().Get(reqPerVhost.ID)
	assert.NotNilf(t, chart, "chart '%s' is not created", reqPerVhost.ID)
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
		assert.Falsef(t, w.Charts().Has(reqPerPort.ID), "chart '%s' is created", reqPerPort.ID)
		return
	}

	chart := w.Charts().Get(reqPerPort.ID)
	assert.NotNilf(t, chart, "chart '%s' is not created", reqPerPort.ID)
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
		assert.Falsef(t, w.Charts().Has(reqPerMethod.ID), "chart '%s' is created", reqPerMethod.ID)
		return
	}

	chart := w.Charts().Get(reqPerMethod.ID)
	assert.NotNilf(t, chart, "chart '%s' is not created", reqPerMethod.ID)
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
		assert.Falsef(t, w.Charts().Has(reqPerVersion.ID), "chart '%s' is created", reqPerVersion.ID)
		return
	}

	chart := w.Charts().Get(reqPerVersion.ID)
	assert.NotNilf(t, chart, "chart '%s' is not created", reqPerVersion.ID)
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
		assert.Falsef(t, w.Charts().Has(reqPerScheme.ID), "chart '%s' is created", reqPerScheme.ID)
	} else {
		assert.Truef(t, w.Charts().Has(reqPerScheme.ID), "chart '%s' is not created", reqPerScheme.ID)
	}
}

func testReqClientCharts(t *testing.T, w *WebLog) {
	if w.mx.ReqIPv4.Value() == 0 && w.mx.ReqIPv6.Value() == 0 {
		assert.Falsef(t, w.Charts().Has(reqPerIPProto.ID), "chart '%s' is created", reqPerIPProto.ID)
	} else {
		assert.Truef(t, w.Charts().Has(reqPerIPProto.ID), "chart '%s' is not created", reqPerIPProto.ID)
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
	if len(w.mx.ReqURLPattern) == 0 || len(w.patURL) == 0 {
		assert.Falsef(t, w.Charts().Has(reqPerURLPattern.ID), "chart '%s' is created", reqPerURLPattern.ID)
		return
	}

	chart := w.Charts().Get(reqPerURLPattern.ID)
	assert.NotNilf(t, chart, "chart '%s' is not created", reqPerURLPattern.ID)
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
		assert.Falsef(t, w.Charts().Has(reqPerSSLProto.ID), "chart '%s' is created", reqPerSSLProto.ID)
		return
	}

	chart := w.Charts().Get(reqPerSSLProto.ID)
	assert.NotNilf(t, chart, "chart '%s' is not created", reqPerSSLProto.ID)
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
		assert.Falsef(t, w.Charts().Has(reqPerSSLCipherSuite.ID), "chart '%s' is created", reqPerSSLCipherSuite.ID)
		return
	}

	chart := w.Charts().Get(reqPerSSLCipherSuite.ID)
	assert.NotNilf(t, chart, "chart '%s' is not created", reqPerSSLCipherSuite.ID)
	if chart == nil {
		return
	}
	for v := range w.mx.ReqSSLCipherSuite {
		id := "req_ssl_cipher_suite_" + v
		assert.Truef(t, chart.HasDim(id), "chart '%s' has no dim for '%s' ssl cipher suite, expected '%s'", chart.ID, v, id)
	}
}

func testReqCustomPatternChart(t *testing.T, w *WebLog) {
	if len(w.mx.ReqCustomPattern) == 0 || len(w.patCustom) == 0 {
		assert.Falsef(t, w.Charts().Has(reqPerCustomPattern.ID), "chart '%s' is created", reqPerCustomPattern.ID)
		return
	}

	chart := w.Charts().Get(reqPerCustomPattern.ID)
	assert.NotNilf(t, chart, "chart '%s' is not created", reqPerCustomPattern.ID)
	if chart == nil {
		return
	}
	for v := range w.mx.ReqCustomPattern {
		id := "req_custom_ptn_" + v
		assert.True(t, chart.HasDim(id), "chart '%s' has no dim for '%s' pattern, expected '%s'", chart.ID, v, id)
	}
}

func testURLPatternStatsCharts(t *testing.T, w *WebLog) {
	for _, p := range w.patURL {
		chartID := fmt.Sprintf(perURLPatternRespStatusCode.ID, p.name)
		chart := w.Charts().Get(chartID)
		assert.NotNilf(t, chart, "chart '%s' is not created", chartID)
		if chart == nil {
			continue
		}

		stats, ok := w.mx.URLPatternStats[p.name]
		assert.Truef(t, ok, "url pattern '%s' has no metric in w.mx.URLPatternStats", p.name)
		if !ok {
			continue
		}
		for v := range stats.RespStatusCode {
			id := fmt.Sprintf("url_ptn_%s_resp_status_code_%s", p.name, v)
			assert.Truef(t, chart.HasDim(id), "chart '%s' has no dim for '%s' code, expected '%s'", chartID, v, id)
		}
	}

	for _, p := range w.patURL {
		id := fmt.Sprintf(perURLPatternBandwidth.ID, p.name)
		if w.mx.BytesSent.Value() == 0 && w.mx.BytesReceived.Value() == 0 {
			assert.Falsef(t, w.Charts().Has(id), "chart '%s' is created", id)
		} else {
			assert.Truef(t, w.Charts().Has(id), "chart '%s' is not created", id)
		}
	}

	for _, p := range w.patURL {
		id := fmt.Sprintf(perURLPatternReqProcTime.ID, p.name)
		if isEmptySummary(w.mx.ReqProcTime) {
			assert.Falsef(t, w.Charts().Has(id), "chart '%s' is created", id)
		} else {
			assert.Truef(t, w.Charts().Has(id), "chart '%s' is not created", id)
		}
	}
}

func testRespStatusCodeChart(t *testing.T, w *WebLog) {
	if !w.GroupRespCodes {
		chart := w.Charts().Get(respCodes.ID)
		assert.NotNilf(t, chart, "chart '%s' is not created", respCodes.ID)
		if chart == nil {
			return
		}
		for v := range w.mx.RespStatusCode {
			id := "resp_status_code_" + v
			assert.Truef(t, chart.HasDim(id), "chart '%s' has no dim for '%s' code, expected '%s'", chart.ID, v, id)
		}
		return
	}

	findCodes := func(class string) (codes []string) {
		for v := range w.mx.RespStatusCode {
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
		if len(codes) == 0 {
			assert.Nilf(t, chart, "chart '%s' is created", chartID)
			continue
		}
		assert.NotNilf(t, chart, "chart '%s' is not created", chartID)
		if chart == nil {
			return
		}
		for _, v := range codes {
			id := "resp_status_code_" + v
			assert.Truef(t, chart.HasDim(id), "chart '%s' has no dim for '%s' code, expected '%s'", chartID, v, id)
		}
	}
	assert.Equal(t, len(w.mx.RespStatusCode), n)
}

var (
	emptySummary   = newWebLogSummary()
	emptyHistogram = metrics.NewHistogram(metrics.DefBuckets)
)

func isEmptySummary(s metrics.Summary) bool     { return reflect.DeepEqual(s, emptySummary) }
func isEmptyHistogram(h metrics.Histogram) bool { return reflect.DeepEqual(h, emptyHistogram) }

// generateLogs is used to populate 'testdata/full.log'
func generateLogs(w io.Writer, matched, unmatched int) error {
	var (
		vhost     = []string{"localhost", "test.example.com", "test.example.org", "198.51.100.1", "2001:db8:1ce::1"}
		scheme    = []string{"http", "https"}
		client    = []string{"localhost", "203.0.113.1", "203.0.113.2", "2001:db8:2ce:1", "2001:db8:2ce:2"}
		method    = []string{"GET", "HEAD", "POST"}
		url       = []string{"invalid.example", "example.com", "example.org", "example.net"}
		version   = []string{"1.1", "2", "2.0"}
		status    = []int{100, 101, 200, 201, 300, 301, 400, 401} // not 5xx on purpose
		sslProto  = []string{"TLSv1", "TLSv1.1", "TLSv1.2", "TLSv1.3", "SSLv2", "SSLv3"}
		sslCipher = []string{"ECDHE-RSA-AES256-SHA", "DHE-RSA-AES256-SHA", "AES256-SHA", "PSK-RC4-SHA"}

		customs = []string{"dark", "light"}
	)
	// test.example.com:80 http 203.0.113.1 TLSv1 AES256-SHA "GET / HTTP/1.1" 200 1674 2674 3674 4674 custom_dark
	for i := 0; i < matched; i++ {
		_, err := fmt.Fprintf(w, "%s:%d %s %s %s %s \"%s /%s HTTP/%s\" %d %d %d %d %d custom_%s\n",
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
			randFromString(customs),
		)
		if err != nil {
			return err
		}
	}
	for i := 0; i < unmatched; i++ {
		_, err := fmt.Fprint(w, "Unmatched! The rat the cat the dog chased killed ate the malt!\n")
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

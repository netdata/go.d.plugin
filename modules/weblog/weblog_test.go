package weblog

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/netdata/go.d.plugin/pkg/logs"
	"github.com/netdata/go.d.plugin/pkg/matcher"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testFormat = []string{
		"$host:$server_port",
		"$scheme",
		"$remote_addr",
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
		},
		CustomPatterns: []userPattern{
			{Name: "dark", Match: "~ dark$"},
			{Name: "light", Match: "~ light$"},
		},
		Histogram:      []float64{10, 20, 30, 40, 100},
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
	//for _, v := range l {
	//	fmt.Println(fmt.Sprintf("\"%s\": %d,", v, m[v]))
	//}

	expected := map[string]int64{
		"bytes_received":                   1110682,
		"bytes_sent":                       1096866,
		"req_custom_ptn_dark":              199,
		"req_custom_ptn_light":             174,
		"req_filtered":                     127,
		"req_http_scheme":                  179,
		"req_https_scheme":                 194,
		"req_ipv4":                         235,
		"req_ipv6":                         138,
		"req_method_GET":                   130,
		"req_method_HEAD":                  124,
		"req_method_POST":                  119,
		"req_port_80":                      78,
		"req_port_81":                      70,
		"req_port_82":                      73,
		"req_port_83":                      85,
		"req_port_84":                      67,
		"req_proc_time_avg":                247,
		"req_proc_time_count":              373,
		"req_proc_time_hist_bucket_1":      12,
		"req_proc_time_hist_bucket_2":      20,
		"req_proc_time_hist_bucket_3":      23,
		"req_proc_time_hist_bucket_4":      29,
		"req_proc_time_hist_bucket_5":      74,
		"req_proc_time_hist_count":         373,
		"req_proc_time_hist_sum":           92141,
		"req_proc_time_max":                499,
		"req_proc_time_min":                1,
		"req_proc_time_sum":                92141,
		"req_unmatched":                    50,
		"req_url_ptn_com":                  116,
		"req_url_ptn_net":                  135,
		"req_url_ptn_org":                  122,
		"req_version_1.1":                  107,
		"req_version_2":                    133,
		"req_version_2.0":                  133,
		"req_vhost_198.51.100.1":           74,
		"req_vhost_2001:db8:1ce::1":        81,
		"req_vhost_localhost":              70,
		"req_vhost_test.example.com":       79,
		"req_vhost_test.example.org":       69,
		"requests":                         550,
		"resp_1xx":                         74,
		"resp_2xx":                         61,
		"resp_3xx":                         95,
		"resp_4xx":                         84,
		"resp_5xx":                         59,
		"resp_client_error":                84,
		"resp_redirect":                    95,
		"resp_server_error":                59,
		"resp_status_code_100":             42,
		"resp_status_code_101":             32,
		"resp_status_code_200":             28,
		"resp_status_code_201":             33,
		"resp_status_code_300":             51,
		"resp_status_code_301":             44,
		"resp_status_code_400":             44,
		"resp_status_code_401":             40,
		"resp_status_code_500":             29,
		"resp_status_code_501":             30,
		"resp_successful":                  135,
		"uniq_ipv4":                        3,
		"uniq_ipv6":                        2,
		"upstream_resp_time_avg":           244,
		"upstream_resp_time_count":         373,
		"upstream_resp_time_hist_bucket_1": 10,
		"upstream_resp_time_hist_bucket_2": 16,
		"upstream_resp_time_hist_bucket_3": 24,
		"upstream_resp_time_hist_bucket_4": 27,
		"upstream_resp_time_hist_bucket_5": 82,
		"upstream_resp_time_hist_count":    373,
		"upstream_resp_time_hist_sum":      91196,
		"upstream_resp_time_max":           499,
		"upstream_resp_time_min":           3,
		"upstream_resp_time_sum":           91196,
		"url_ptn_com_bytes_received":       339904,
		"url_ptn_com_bytes_sent":           342541,
		"url_ptn_com_req_proc_time_avg":    238,
		"url_ptn_com_req_proc_time_count":  116,
		"url_ptn_com_req_proc_time_max":    499,
		"url_ptn_com_req_proc_time_min":    3,
		"url_ptn_com_req_proc_time_sum":    27663,
		"url_ptn_com_resp_status_code_100": 13,
		"url_ptn_com_resp_status_code_101": 8,
		"url_ptn_com_resp_status_code_200": 7,
		"url_ptn_com_resp_status_code_201": 6,
		"url_ptn_com_resp_status_code_300": 20,
		"url_ptn_com_resp_status_code_301": 15,
		"url_ptn_com_resp_status_code_400": 12,
		"url_ptn_com_resp_status_code_401": 15,
		"url_ptn_com_resp_status_code_500": 9,
		"url_ptn_com_resp_status_code_501": 11,
		"url_ptn_net_bytes_received":       403033,
		"url_ptn_net_bytes_sent":           393199,
		"url_ptn_net_req_proc_time_avg":    247,
		"url_ptn_net_req_proc_time_count":  135,
		"url_ptn_net_req_proc_time_max":    497,
		"url_ptn_net_req_proc_time_min":    4,
		"url_ptn_net_req_proc_time_sum":    33347,
		"url_ptn_net_resp_status_code_100": 16,
		"url_ptn_net_resp_status_code_101": 10,
		"url_ptn_net_resp_status_code_200": 11,
		"url_ptn_net_resp_status_code_201": 16,
		"url_ptn_net_resp_status_code_300": 13,
		"url_ptn_net_resp_status_code_301": 17,
		"url_ptn_net_resp_status_code_400": 20,
		"url_ptn_net_resp_status_code_401": 11,
		"url_ptn_net_resp_status_code_500": 13,
		"url_ptn_net_resp_status_code_501": 8,
		"url_ptn_org_bytes_received":       367745,
		"url_ptn_org_bytes_sent":           361126,
		"url_ptn_org_req_proc_time_avg":    255,
		"url_ptn_org_req_proc_time_count":  122,
		"url_ptn_org_req_proc_time_max":    490,
		"url_ptn_org_req_proc_time_min":    1,
		"url_ptn_org_req_proc_time_sum":    31131,
		"url_ptn_org_resp_status_code_100": 13,
		"url_ptn_org_resp_status_code_101": 14,
		"url_ptn_org_resp_status_code_200": 10,
		"url_ptn_org_resp_status_code_201": 11,
		"url_ptn_org_resp_status_code_300": 18,
		"url_ptn_org_resp_status_code_301": 12,
		"url_ptn_org_resp_status_code_400": 12,
		"url_ptn_org_resp_status_code_401": 14,
		"url_ptn_org_resp_status_code_500": 7,
		"url_ptn_org_resp_status_code_501": 11,
	}
	_ = expected

	assert.Equal(t, expected, weblog.Collect())
	//testDynamicCharts(t, weblog)

	//sort.Slice(*weblog.charts, func(i, j int) bool {
	//	return (*weblog.charts)[i].Priority < (*weblog.charts)[j].Priority
	//})
	//
	//for _, v := range *weblog.charts {
	//	if v.Priority == 0 {
	//		t.Fatal(v.ID, "NO PRIORITY")
	//
	//	}
	//	fmt.Printf("%s '%s' %d\n", v.ID, v.Fam, v.Priority)
	//}
}

//func testDynamicCharts(t *testing.T, w *WebLog) {
//	if w.GroupRespCodes {
//		chart := w.Charts().Get(respCodes.ID)
//		require.NotNil(t, chart)
//		assert.Len(t, chart.Dims, len(w.mx.RespStatusCode))
//	} else {
//		for _, v := range []string{"_1xx", "_2xx", "_3xx", "_4xx", "_5xx", "_other"} {
//			chart := w.Charts().Get(respCodes.ID + v)
//			require.NotNil(t, chart)
//			assert.True(t, len(chart.Dims) > 0)
//		}
//	}
//	if w.col.vhost {
//		chart := w.Charts().Get(reqPerVhost.ID)
//		require.NotNil(t, chart)
//		assert.Len(t, chart.Dims, len(w.mx.ReqVhost))
//	}
//	if w.col.port {
//		chart := w.Charts().Get(reqPerPort.ID)
//		require.NotNil(t, chart)
//		assert.Len(t, chart.Dims, len(w.mx.ReqPort))
//	}
//	if w.col.scheme {
//		assert.NotNil(t, w.Charts().Get(reqPerScheme.ID))
//	}
//	if w.col.client {
//		assert.NotNil(t, w.Charts().Get(reqPerIPProto.ID))
//		assert.NotNil(t, w.Charts().Get(uniqIPsCurPoll.ID))
//	}
//	if w.col.method {
//		chart := w.Charts().Get(reqPerMethod.ID)
//		require.NotNil(t, chart)
//		assert.Len(t, chart.Dims, len(w.mx.ReqMethod))
//	}
//	if w.col.url && len(w.patURL) != 0 {
//		chart := w.Charts().Get(reqPerURLPattern.ID)
//		require.NotNil(t, chart)
//		assert.Len(t, chart.Dims, len(w.patURL))
//	}
//	if w.col.version {
//		chart := w.Charts().Get(reqPerVersion.ID)
//		require.NotNil(t, chart)
//		assert.Len(t, chart.Dims, len(w.mx.ReqVersion))
//	}
//	if w.col.reqSize || w.col.respSize {
//		assert.NotNil(t, w.Charts().Get(bandwidth.ID))
//	}
//	if w.col.custom && len(w.patCustom) != 0 {
//		chart := w.Charts().Get(reqPerCustomPattern.ID)
//		require.NotNil(t, chart)
//		assert.Len(t, chart.Dims, len(w.mx.ReqCustomPattern))
//	}
//	if w.col.reqProcTime {
//		assert.NotNil(t, w.Charts().Get(reqProcTime.ID))
//		if len(w.Histogram) != 0 {
//			assert.NotNil(t, w.Charts().Get(respTimeHist.ID))
//		}
//	}
//	if w.col.upRespTime {
//		assert.NotNil(t, w.Charts().Get(upsRespTime.ID))
//		if len(w.Histogram) != 0 {
//			assert.NotNil(t, w.Charts().Get(upsRespTimeHist.ID))
//		}
//	}
//
//	if w.col.url {
//		for _, cat := range w.patURL {
//			assert.NotNil(t, w.Charts().Get(respCodes.ID+"_"+cat.name))
//		}
//		if w.col.reqSize || w.col.respSize {
//			for _, cat := range w.patURL {
//				assert.NotNil(t, w.Charts().Get(bandwidth.ID+"_"+cat.name))
//			}
//		}
//		if w.col.reqProcTime {
//			for _, cat := range w.patURL {
//				assert.NotNil(t, w.Charts().Get(reqProcTime.ID+"_"+cat.name))
//			}
//		}
//	}
//}

// generateLogs is used to populate 'testdata/full.log'
func generateLogs(w io.Writer, matched, unmatched int) error {
	var (
		vhosts   = []string{"localhost", "test.example.com", "test.example.org", "198.51.100.1", "2001:db8:1ce::1"}
		schemes  = []string{"http", "https"}
		clients  = []string{"localhost", "203.0.113.1", "203.0.113.2", "2001:db8:2ce:1", "2001:db8:2ce:2"}
		methods  = []string{"GET", "HEAD", "POST"}
		urls     = []string{"invalid.example", "example.com", "example.org", "example.net"}
		versions = []string{"1.1", "2", "2.0"}
		statuses = []int{100, 101, 200, 201, 300, 301, 400, 401, 500, 501}
		customs  = []string{"dark", "light"}
	)
	// test.example.com:80 http 203.0.113.1 "GET / HTTP/1.1" 200 1674 2674 3674 4674 custom_dark
	for i := 0; i < matched; i++ {
		_, err := fmt.Fprintf(w, "%s:%d %s %s \"%s /%s HTTP/%s\" %d %d %d %d %d custom_%s\n",
			randFromString(vhosts),
			randInt(80, 85),
			randFromString(schemes),
			randFromString(clients),
			randFromString(methods),
			randFromString(urls),
			randFromString(versions),
			randFromInt(statuses),
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

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
		URLCategories: []rawCategory{
			{Name: "com", Match: "~ com$"},
			{Name: "org", Match: "~ org$"},
			{Name: "net", Match: "~ net$"},
		},
		UserCategories: []rawCategory{
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
		"bytes_received":                   1152899,
		"bytes_sent":                       1165589,
		"req_custom_dark":                  170,
		"req_custom_light":                 214,
		"req_filtered":                     116,
		"req_http_scheme":                  187,
		"req_https_scheme":                 197,
		"req_ipv4":                         210,
		"req_ipv6":                         174,
		"req_method_GET":                   127,
		"req_method_HEAD":                  134,
		"req_method_POST":                  123,
		"req_port_80":                      81,
		"req_port_81":                      76,
		"req_port_82":                      72,
		"req_port_83":                      79,
		"req_port_84":                      76,
		"req_proc_time_avg":                252,
		"req_proc_time_count":              384,
		"req_proc_time_hist_bucket_1":      6,
		"req_proc_time_hist_bucket_2":      14,
		"req_proc_time_hist_bucket_3":      23,
		"req_proc_time_hist_bucket_4":      29,
		"req_proc_time_hist_bucket_5":      73,
		"req_proc_time_hist_count":         384,
		"req_proc_time_hist_sum":           96992,
		"req_proc_time_max":                499,
		"req_proc_time_min":                4,
		"req_proc_time_sum":                96992,
		"req_unmatched":                    0,
		"req_url_com":                      118,
		"req_url_net":                      127,
		"req_url_org":                      139,
		"req_version_1.1":                  135,
		"req_version_2":                    120,
		"req_version_2.0":                  129,
		"req_vhost_198.51.100.1":           65,
		"req_vhost_2001:db8:1ce::1":        75,
		"req_vhost_localhost":              79,
		"req_vhost_test.example.com":       74,
		"req_vhost_test.example.org":       91,
		"requests":                         500,
		"resp_1xx":                         85,
		"resp_2xx":                         73,
		"resp_3xx":                         75,
		"resp_4xx":                         69,
		"resp_5xx":                         82,
		"resp_client_error":                69,
		"resp_code_100":                    41,
		"resp_code_101":                    44,
		"resp_code_200":                    44,
		"resp_code_201":                    29,
		"resp_code_300":                    35,
		"resp_code_301":                    40,
		"resp_code_400":                    40,
		"resp_code_401":                    29,
		"resp_code_500":                    42,
		"resp_code_501":                    40,
		"resp_redirect":                    75,
		"resp_server_error":                82,
		"resp_successful":                  158,
		"uniq_ipv4":                        3,
		"uniq_ipv6":                        2,
		"upstream_resp_time_avg":           242,
		"upstream_resp_time_count":         384,
		"upstream_resp_time_hist_bucket_1": 2,
		"upstream_resp_time_hist_bucket_2": 8,
		"upstream_resp_time_hist_bucket_3": 17,
		"upstream_resp_time_hist_bucket_4": 27,
		"upstream_resp_time_hist_bucket_5": 89,
		"upstream_resp_time_hist_count":    384,
		"upstream_resp_time_hist_sum":      93257,
		"upstream_resp_time_max":           499,
		"upstream_resp_time_min":           4,
		"upstream_resp_time_sum":           93257,
		"url_com_bytes_received":           371169,
		"url_com_bytes_sent":               366726,
		"url_com_req_code_100":             8,
		"url_com_req_code_101":             14,
		"url_com_req_code_200":             15,
		"url_com_req_code_201":             10,
		"url_com_req_code_300":             7,
		"url_com_req_code_301":             7,
		"url_com_req_code_400":             13,
		"url_com_req_code_401":             13,
		"url_com_req_code_500":             16,
		"url_com_req_code_501":             15,
		"url_com_req_proc_time_avg":        241,
		"url_com_req_proc_time_count":      118,
		"url_com_req_proc_time_max":        486,
		"url_com_req_proc_time_min":        15,
		"url_com_req_proc_time_sum":        28498,
		"url_net_bytes_received":           390107,
		"url_net_bytes_sent":               373219,
		"url_net_req_code_100":             17,
		"url_net_req_code_101":             16,
		"url_net_req_code_200":             15,
		"url_net_req_code_201":             9,
		"url_net_req_code_300":             11,
		"url_net_req_code_301":             13,
		"url_net_req_code_400":             14,
		"url_net_req_code_401":             6,
		"url_net_req_code_500":             13,
		"url_net_req_code_501":             13,
		"url_net_req_proc_time_avg":        243,
		"url_net_req_proc_time_count":      127,
		"url_net_req_proc_time_max":        499,
		"url_net_req_proc_time_min":        4,
		"url_net_req_proc_time_sum":        30895,
		"url_org_bytes_received":           391623,
		"url_org_bytes_sent":               425644,
		"url_org_req_code_100":             16,
		"url_org_req_code_101":             14,
		"url_org_req_code_200":             14,
		"url_org_req_code_201":             10,
		"url_org_req_code_300":             17,
		"url_org_req_code_301":             20,
		"url_org_req_code_400":             13,
		"url_org_req_code_401":             10,
		"url_org_req_code_500":             13,
		"url_org_req_code_501":             12,
		"url_org_req_proc_time_avg":        270,
		"url_org_req_proc_time_count":      139,
		"url_org_req_proc_time_max":        499,
		"url_org_req_proc_time_min":        7,
		"url_org_req_proc_time_sum":        37599,
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
//	if w.col.url && len(w.catURL) != 0 {
//		chart := w.Charts().Get(reqPerURL.ID)
//		require.NotNil(t, chart)
//		assert.Len(t, chart.Dims, len(w.catURL))
//	}
//	if w.col.version {
//		chart := w.Charts().Get(reqPerVersion.ID)
//		require.NotNil(t, chart)
//		assert.Len(t, chart.Dims, len(w.mx.ReqVersion))
//	}
//	if w.col.reqSize || w.col.respSize {
//		assert.NotNil(t, w.Charts().Get(bandwidth.ID))
//	}
//	if w.col.custom && len(w.catCustom) != 0 {
//		chart := w.Charts().Get(reqPerCustom.ID)
//		require.NotNil(t, chart)
//		assert.Len(t, chart.Dims, len(w.mx.ReqCustom))
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
//		for _, cat := range w.catURL {
//			assert.NotNil(t, w.Charts().Get(respCodes.ID+"_"+cat.name))
//		}
//		if w.col.reqSize || w.col.respSize {
//			for _, cat := range w.catURL {
//				assert.NotNil(t, w.Charts().Get(bandwidth.ID+"_"+cat.name))
//			}
//		}
//		if w.col.reqProcTime {
//			for _, cat := range w.catURL {
//				assert.NotNil(t, w.Charts().Get(reqProcTime.ID+"_"+cat.name))
//			}
//		}
//	}
//}

// generateLogs is used to populate 'testdata/full.log'
func generateLogs(w io.Writer, n int) error {
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
	const row = "%s:%d %s %s \"%s /%s HTTP/%s\" %d %d %d %d %d custom_%s\n"
	for i := 0; i < n; i++ {
		_, err := fmt.Fprintf(w, row,
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
	return nil
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func randFromString(s []string) string { return s[r.Intn(len(s))] }
func randFromInt(s []int) int          { return s[r.Intn(len(s))] }
func randInt(min, max int) int         { return r.Intn(max-min) + min }

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

//func init() {
//	logger.SetSeverity(logger.DEBUG)
//}

var (
	testFullLog, _ = ioutil.ReadFile("testdata/full.log")
)

func Test_readTestData(t *testing.T) {
	assert.NotNil(t, testFullLog)
}

func TestWebLog_Init(t *testing.T) {
	weblog := New()
	weblog.Config.Filter.Includes = []string{"~ .php$"}
	weblog.Config.URLCategories = []rawCategory{
		{"foo", "= foo"},
		{"bar", "= bar"},
	}
	weblog.Config.UserCategories = []rawCategory{
		{"baz", "= baz"},
		{"foobar", "= foobar"},
	}
	ok := weblog.Init()

	require.True(t, ok)
	assert.True(t, weblog.filter.MatchString("/abc.php"))
	assert.False(t, weblog.filter.MatchString("/abc.html"))

	assert.Len(t, weblog.urlCats, 2)
	assert.Equal(t, "foo", weblog.urlCats[0].name)
	assert.True(t, weblog.urlCats[0].Matcher.MatchString("foo"))
	assert.Equal(t, "bar", weblog.urlCats[1].name)
	assert.True(t, weblog.urlCats[1].Matcher.MatchString("bar"))

	assert.Len(t, weblog.userCats, 2)
	assert.Equal(t, "baz", weblog.userCats[0].name)
	assert.True(t, weblog.userCats[0].Matcher.MatchString("baz"))
	assert.Equal(t, "foobar", weblog.userCats[1].name)
	assert.True(t, weblog.userCats[1].Matcher.MatchString("foobar"))
}

func TestWebLog_Collect(t *testing.T) {
	format := []string{
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
	config := Config{
		Parser: logs.ParserConfig{
			LogType: logs.TypeCSV,
			CSV: logs.CSVConfig{
				Delimiter:        ' ',
				TrimLeadingSpace: false,
				Format:           strings.Join(format, " "),
				CheckField:       checkCSVFormatField,
			},
		},
		Path:        "testdata/full.log",
		ExcludePath: "",
		Filter:      matcher.SimpleExpr{Excludes: []string{"* /invalid*"}},
		URLCategories: []rawCategory{
			{Name: "com", Match: "* *com"},
			{Name: "org", Match: "* *org"},
			{Name: "net", Match: "* *net"},
		},
		UserCategories: []rawCategory{
			{Name: "dark", Match: "* *dark"},
			{Name: "light", Match: "* *light"},
		},
		Histogram:              nil,
		AggregateResponseCodes: true,
	}

	weblog := New()
	weblog.Config = config
	defer weblog.Cleanup()

	require.True(t, weblog.Init())
	p, err := logs.NewCSVParser(config.Parser.CSV, bytes.NewReader(testFullLog))
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
		"bytes_received":                    1095824,
		"bytes_sent":                        1128312,
		"req_code_200":                      52,
		"req_code_201":                      49,
		"req_code_300":                      45,
		"req_code_301":                      46,
		"req_code_400":                      51,
		"req_code_401":                      42,
		"req_code_500":                      46,
		"req_code_501":                      37,
		"req_custom_dark":                   192,
		"req_custom_light":                  176,
		"req_filtered":                      132,
		"req_http_scheme":                   203,
		"req_https_scheme":                  165,
		"req_ipv4":                          224,
		"req_ipv4_uniq":                     3,
		"req_ipv6":                          144,
		"req_ipv6_uniq":                     2,
		"req_method_GET":                    101,
		"req_method_HEAD":                   126,
		"req_method_POST":                   141,
		"req_port_80":                       87,
		"req_port_81":                       57,
		"req_port_82":                       81,
		"req_port_83":                       70,
		"req_port_84":                       73,
		"req_unmatched":                     0,
		"req_uri_com":                       127,
		"req_uri_net":                       112,
		"req_uri_org":                       129,
		"req_version_1.1":                   114,
		"req_version_2":                     132,
		"req_version_2.0":                   122,
		"req_vhost_198.51.100.1":            69,
		"req_vhost_2001:db8:1ce::1":         68,
		"req_vhost_localhost":               69,
		"req_vhost_test.example.com":        86,
		"req_vhost_test.example.org":        76,
		"requests":                          368,
		"resp_1xx":                          0,
		"resp_2xx":                          101,
		"resp_3xx":                          91,
		"resp_4xx":                          93,
		"resp_5xx":                          83,
		"resp_client_error":                 93,
		"resp_redirect":                     91,
		"resp_server_error":                 83,
		"resp_successful":                   101,
		"resp_time_avg":                     3002703,
		"resp_time_count":                   368,
		"resp_time_hist_bucket_1":           0,
		"resp_time_hist_bucket_10":          0,
		"resp_time_hist_bucket_11":          0,
		"resp_time_hist_bucket_2":           0,
		"resp_time_hist_bucket_3":           0,
		"resp_time_hist_bucket_4":           0,
		"resp_time_hist_bucket_5":           0,
		"resp_time_hist_bucket_6":           0,
		"resp_time_hist_bucket_7":           0,
		"resp_time_hist_bucket_8":           0,
		"resp_time_hist_bucket_9":           0,
		"resp_time_hist_count":              368,
		"resp_time_hist_sum":                1104995,
		"resp_time_max":                     4989000,
		"resp_time_min":                     1002000,
		"resp_time_sum":                     1104995000,
		"resp_time_upstream_avg":            2997956,
		"resp_time_upstream_count":          368,
		"resp_time_upstream_hist_bucket_1":  0,
		"resp_time_upstream_hist_bucket_10": 0,
		"resp_time_upstream_hist_bucket_11": 0,
		"resp_time_upstream_hist_bucket_2":  0,
		"resp_time_upstream_hist_bucket_3":  0,
		"resp_time_upstream_hist_bucket_4":  0,
		"resp_time_upstream_hist_bucket_5":  0,
		"resp_time_upstream_hist_bucket_6":  0,
		"resp_time_upstream_hist_bucket_7":  0,
		"resp_time_upstream_hist_bucket_8":  0,
		"resp_time_upstream_hist_bucket_9":  0,
		"resp_time_upstream_hist_count":     368,
		"resp_time_upstream_hist_sum":       1103248,
		"resp_time_upstream_max":            4990000,
		"resp_time_upstream_min":            1019000,
		"resp_time_upstream_sum":            1103248000,
		"uri_com_bytes_received":            404866,
		"uri_com_bytes_sent":                423446,
		"uri_com_req_code_200":              20,
		"uri_com_req_code_201":              17,
		"uri_com_req_code_300":              17,
		"uri_com_req_code_301":              19,
		"uri_com_req_code_400":              14,
		"uri_com_req_code_401":              12,
		"uri_com_req_code_500":              16,
		"uri_com_req_code_501":              12,
		"uri_com_resp_time_avg":             3080622,
		"uri_com_resp_time_count":           127,
		"uri_com_resp_time_max":             4903000,
		"uri_com_resp_time_min":             1002000,
		"uri_com_resp_time_sum":             391239000,
		"uri_net_bytes_received":            296230,
		"uri_net_bytes_sent":                321482,
		"uri_net_req_code_200":              17,
		"uri_net_req_code_201":              13,
		"uri_net_req_code_300":              15,
		"uri_net_req_code_301":              13,
		"uri_net_req_code_400":              17,
		"uri_net_req_code_401":              18,
		"uri_net_req_code_500":              13,
		"uri_net_req_code_501":              6,
		"uri_net_resp_time_avg":             2873098,
		"uri_net_resp_time_count":           112,
		"uri_net_resp_time_max":             4989000,
		"uri_net_resp_time_min":             1002000,
		"uri_net_resp_time_sum":             321787000,
		"uri_org_bytes_received":            394728,
		"uri_org_bytes_sent":                383384,
		"uri_org_req_code_200":              15,
		"uri_org_req_code_201":              19,
		"uri_org_req_code_300":              13,
		"uri_org_req_code_301":              14,
		"uri_org_req_code_400":              20,
		"uri_org_req_code_401":              12,
		"uri_org_req_code_500":              17,
		"uri_org_req_code_501":              19,
		"uri_org_resp_time_avg":             3038519,
		"uri_org_resp_time_count":           129,
		"uri_org_resp_time_max":             4914000,
		"uri_org_resp_time_min":             1022000,
		"uri_org_resp_time_sum":             391969000,
	}

	assert.Equal(t, expected, weblog.Collect())
}

var (
	vhosts   = []string{"localhost", "test.example.com", "test.example.org", "198.51.100.1", "2001:db8:1ce::1"}
	schemes  = []string{"http", "https"}
	clients  = []string{"localhost", "203.0.113.1", "203.0.113.2", "2001:db8:2ce:1", "2001:db8:2ce:2"}
	methods  = []string{"GET", "HEAD", "POST"}
	urls     = []string{"invalid.example", "example.com", "example.org", "example.net"}
	versions = []string{"1.1", "2", "2.0"}
	statuses = []int{200, 201, 300, 301, 400, 401, 500, 501}
	customs  = []string{"dark", "light"}
)

func generateLogs(w io.Writer, n int) error {
	// 	format := []string{
	//		"$host:$server_port",
	//		"$scheme",
	//		"$remote_addr",
	//		`"$request"`,
	//		"$status",
	//		"$body_bytes_sent",
	//		"$request_length",
	//		"$request_time",
	//		"$upstream_response_time",
	//		"$custom",
	//	}
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
			randInt(1000, 5000),
			randInt(1000, 5000),
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

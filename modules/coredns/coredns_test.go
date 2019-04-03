package coredns

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testNoLoad, _   = ioutil.ReadFile("testdata/no_load.txt")
	testSomeLoad, _ = ioutil.ReadFile("testdata/some_load.txt")
)

func TestNew(t *testing.T) {
	job := New()

	assert.IsType(t, (*CoreDNS)(nil), job)
	assert.Equal(t, defaultURL, job.URL)
	assert.Equal(t, defaultHTTPTimeout, job.Timeout.Duration)
}

func TestCoreDNS_Charts(t *testing.T) { assert.NotNil(t, New().Charts()) }

func TestCoreDNS_Cleanup(t *testing.T) { New().Cleanup() }

func TestCoreDNS_Init(t *testing.T) { assert.True(t, New().Init()) }

func TestCoreDNS_InitNG(t *testing.T) {
	job := New()
	job.URL = ""
	assert.False(t, job.Init())
}

func TestCoreDNS_Check(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testSomeLoad)
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL + "/metrics"
	require.True(t, job.Init())
	assert.True(t, job.Check())
}

func TestCoreDNS_CheckNG(t *testing.T) {
	job := New()
	job.URL = "http://127.0.0.1:38001/metrics"
	require.True(t, job.Init())
	assert.False(t, job.Check())
}

func TestCoreDNS_Collect(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testSomeLoad)
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL + "/metrics"
	require.True(t, job.Init())
	require.True(t, job.Check())

	//m := job.Collect()
	//l := make([]string, 0)
	//for k := range m {
	//	l = append(l, k)
	//}
	//sort.Strings(l)
	//for _, v := range l {
	//	fmt.Println(fmt.Sprintf("\"%s\": %d,", v, m[v]))
	//}

	expected := map[string]int64{
		"dns://:53_request_by_ip_family_v4":     119,
		"dns://:53_request_by_ip_family_v6":     62,
		"dns://:53_request_by_proto_tcp":        62,
		"dns://:53_request_by_proto_udp":        119,
		"dns://:53_request_by_status_dropped":   58,
		"dns://:53_request_by_status_processed": 123,
		"dns://:53_request_by_type_A":           89,
		"dns://:53_request_by_type_AAAA":        29,
		"dns://:53_request_by_type_ANY":         0,
		"dns://:53_request_by_type_CNAME":       0,
		"dns://:53_request_by_type_DNSKEY":      0,
		"dns://:53_request_by_type_DS":          0,
		"dns://:53_request_by_type_IXFR":        0,
		"dns://:53_request_by_type_MX":          1,
		"dns://:53_request_by_type_NS":          0,
		"dns://:53_request_by_type_NSEC":        0,
		"dns://:53_request_by_type_NSEC3":       0,
		"dns://:53_request_by_type_PTR":         0,
		"dns://:53_request_by_type_RRSIG":       0,
		"dns://:53_request_by_type_SOA":         0,
		"dns://:53_request_by_type_SRV":         0,
		"dns://:53_request_by_type_TXT":         0,
		"dns://:53_request_by_type_other":       0,
		"dns://:53_request_total":               0,
		"panic_total":                           99,
		"request_by_ip_family_v4":               119,
		"request_by_ip_family_v6":               62,
		"request_by_proto_tcp":                  0,
		"request_by_proto_udp":                  181,
		"request_by_status_dropped":             58,
		"request_by_status_processed":           123,
		"request_by_type_A":                     89,
		"request_by_type_AAAA":                  29,
		"request_by_type_ANY":                   0,
		"request_by_type_CNAME":                 0,
		"request_by_type_DNSKEY":                0,
		"request_by_type_DS":                    0,
		"request_by_type_IXFR":                  0,
		"request_by_type_MX":                    1,
		"request_by_type_NS":                    0,
		"request_by_type_NSEC":                  0,
		"request_by_type_NSEC3":                 0,
		"request_by_type_PTR":                   0,
		"request_by_type_RRSIG":                 0,
		"request_by_type_SOA":                   0,
		"request_by_type_SRV":                   0,
		"request_by_type_TXT":                   0,
		"request_by_type_other":                 0,
		"request_total":                         181,
	}

	assert.Equal(t, expected, job.Collect())
}

func TestCoreDNS_InvalidData(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("hello and goodbye"))
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL + "/metrics"
	require.True(t, job.Init())
	assert.False(t, job.Check())
}

func TestCoreDNS_404(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL + "/metrics"
	require.True(t, job.Init())
	assert.False(t, job.Check())
}

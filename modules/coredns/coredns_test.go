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
	job.PerServerStatsPermitFrom = "glob:*"
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
		"dns://:53_request_duration_seconds_bucket_+Inf":    0,
		"dns://:53_request_duration_seconds_bucket_0.00025": 108,
		"dns://:53_request_duration_seconds_bucket_0.0005":  5,
		"dns://:53_request_duration_seconds_bucket_0.001":   0,
		"dns://:53_request_duration_seconds_bucket_0.002":   0,
		"dns://:53_request_duration_seconds_bucket_0.004":   0,
		"dns://:53_request_duration_seconds_bucket_0.008":   2,
		"dns://:53_request_duration_seconds_bucket_0.016":   1,
		"dns://:53_request_duration_seconds_bucket_0.032":   1,
		"dns://:53_request_duration_seconds_bucket_0.064":   0,
		"dns://:53_request_duration_seconds_bucket_0.128":   0,
		"dns://:53_request_duration_seconds_bucket_0.256":   0,
		"dns://:53_request_duration_seconds_bucket_0.512":   0,
		"dns://:53_request_duration_seconds_bucket_1.024":   2,
		"dns://:53_request_duration_seconds_bucket_2.048":   0,
		"dns://:53_request_duration_seconds_bucket_4.096":   0,
		"dns://:53_request_duration_seconds_bucket_8.192":   0,
		"dns://:53_request_per_ip_family_v4":                119,
		"dns://:53_request_per_ip_family_v6":                62,
		"dns://:53_request_per_proto_tcp":                   62,
		"dns://:53_request_per_proto_udp":                   119,
		"dns://:53_request_per_status_dropped":              58,
		"dns://:53_request_per_status_processed":            65,
		"dns://:53_request_per_type_A":                      45,
		"dns://:53_request_per_type_AAAA":                   15,
		"dns://:53_request_per_type_ANY":                    0,
		"dns://:53_request_per_type_CNAME":                  0,
		"dns://:53_request_per_type_DNSKEY":                 0,
		"dns://:53_request_per_type_DS":                     0,
		"dns://:53_request_per_type_IXFR":                   0,
		"dns://:53_request_per_type_MX":                     1,
		"dns://:53_request_per_type_NS":                     0,
		"dns://:53_request_per_type_NSEC":                   0,
		"dns://:53_request_per_type_NSEC3":                  0,
		"dns://:53_request_per_type_PTR":                    0,
		"dns://:53_request_per_type_RRSIG":                  0,
		"dns://:53_request_per_type_SOA":                    0,
		"dns://:53_request_per_type_SRV":                    0,
		"dns://:53_request_per_type_TXT":                    0,
		"dns://:53_request_per_type_other":                  0,
		"dns://:53_request_total":                           181,
		"dns://:53_response_per_rcode_BADALG":               0,
		"dns://:53_response_per_rcode_BADCOOKIE":            0,
		"dns://:53_response_per_rcode_BADKEY":               0,
		"dns://:53_response_per_rcode_BADMODE":              0,
		"dns://:53_response_per_rcode_BADNAME":              0,
		"dns://:53_response_per_rcode_BADSIG":               0,
		"dns://:53_response_per_rcode_BADTIME":              0,
		"dns://:53_response_per_rcode_BADTRUNC":             0,
		"dns://:53_response_per_rcode_FORMERR":              0,
		"dns://:53_response_per_rcode_NOERROR":              3,
		"dns://:53_response_per_rcode_NOTAUTH":              0,
		"dns://:53_response_per_rcode_NOTIMP":               0,
		"dns://:53_response_per_rcode_NOTZONE":              0,
		"dns://:53_response_per_rcode_NXDOMAIN":             0,
		"dns://:53_response_per_rcode_NXRRSET":              0,
		"dns://:53_response_per_rcode_REFUSED":              0,
		"dns://:53_response_per_rcode_SERVFAIL":             58,
		"dns://:53_response_per_rcode_YXDOMAIN":             0,
		"dns://:53_response_per_rcode_YXRRSET":              0,
		"dns://:53_response_per_rcode_other":                0,
		"dns://:53_response_total":                          61,
		"no_matching_zone_dropped_total":                    111,
		"panic_total":                                       99,
		"request_duration_seconds_bucket_+Inf":              0,
		"request_duration_seconds_bucket_0.00025":           108,
		"request_duration_seconds_bucket_0.0005":            5,
		"request_duration_seconds_bucket_0.001":             0,
		"request_duration_seconds_bucket_0.002":             0,
		"request_duration_seconds_bucket_0.004":             0,
		"request_duration_seconds_bucket_0.008":             2,
		"request_duration_seconds_bucket_0.016":             1,
		"request_duration_seconds_bucket_0.032":             1,
		"request_duration_seconds_bucket_0.064":             0,
		"request_duration_seconds_bucket_0.128":             0,
		"request_duration_seconds_bucket_0.256":             0,
		"request_duration_seconds_bucket_0.512":             0,
		"request_duration_seconds_bucket_1.024":             2,
		"request_duration_seconds_bucket_2.048":             0,
		"request_duration_seconds_bucket_4.096":             0,
		"request_duration_seconds_bucket_8.192":             0,
		"request_per_ip_family_v4":                          119,
		"request_per_ip_family_v6":                          62,
		"request_per_proto_tcp":                             62,
		"request_per_proto_udp":                             119,
		"request_per_status_dropped":                        58,
		"request_per_status_processed":                      65,
		"request_per_type_A":                                45,
		"request_per_type_AAAA":                             15,
		"request_per_type_ANY":                              0,
		"request_per_type_CNAME":                            0,
		"request_per_type_DNSKEY":                           0,
		"request_per_type_DS":                               0,
		"request_per_type_IXFR":                             0,
		"request_per_type_MX":                               1,
		"request_per_type_NS":                               0,
		"request_per_type_NSEC":                             0,
		"request_per_type_NSEC3":                            0,
		"request_per_type_PTR":                              0,
		"request_per_type_RRSIG":                            0,
		"request_per_type_SOA":                              0,
		"request_per_type_SRV":                              0,
		"request_per_type_TXT":                              0,
		"request_per_type_other":                            0,
		"request_total":                                     181,
		"response_per_rcode_BADALG":                         0,
		"response_per_rcode_BADCOOKIE":                      0,
		"response_per_rcode_BADKEY":                         0,
		"response_per_rcode_BADMODE":                        0,
		"response_per_rcode_BADNAME":                        0,
		"response_per_rcode_BADSIG":                         0,
		"response_per_rcode_BADTIME":                        0,
		"response_per_rcode_BADTRUNC":                       0,
		"response_per_rcode_FORMERR":                        0,
		"response_per_rcode_NOERROR":                        3,
		"response_per_rcode_NOTAUTH":                        0,
		"response_per_rcode_NOTIMP":                         0,
		"response_per_rcode_NOTZONE":                        0,
		"response_per_rcode_NXDOMAIN":                       0,
		"response_per_rcode_NXRRSET":                        0,
		"response_per_rcode_REFUSED":                        0,
		"response_per_rcode_SERVFAIL":                       58,
		"response_per_rcode_YXDOMAIN":                       0,
		"response_per_rcode_YXRRSET":                        0,
		"response_per_rcode_other":                          0,
		"response_total":                                    61,
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

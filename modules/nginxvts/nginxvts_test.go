package nginxvts

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testStatusDataBasic, _ = ioutil.ReadFile("testdata/vts-basic.json")
	testStatusDataFull, _  = ioutil.ReadFile("testdata/vts-full.json")
)

func TestCleanup(t *testing.T) { New().Cleanup() }

func TestNew(t *testing.T) {
	job := New()

	assert.Implements(t, (*module.Module)(nil), job)
	assert.Equal(t, defaultURL, job.URL)
	assert.Equal(t, defaultHTTPTimeout, job.Timeout.Duration)
}

func TestInit(t *testing.T) {
	job0 := New()
	job1 := &NginxVts{
		Config: Config{
			HTTP: web.HTTP{
				Request: web.Request{
					URL: "",
				},
			},
		},
	}

	require.True(t, job0.Init())
	require.False(t, job1.Init())
}

func TestCheck(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testStatusDataBasic)
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL
	require.True(t, job.Init())
	assert.True(t, job.Check())
}

func TestCharts(t *testing.T) { assert.NotNil(t, New().Charts()) }

func TestCollect_basic(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testStatusDataBasic)
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL

	require.True(t, job.Init())
	require.True(t, job.Check())

	expected := map[string]int64{
		"loadmsec":             1606489796895,
		"nowmsec":              1606490116734,
		"connections_accepted": 12,
		"connections_active":   2,
		"connections_handled":  12,
		"connections_reading":  0,
		"connections_requests": 0,
		"connections_waiting":  1,
		"connections_writing":  1,
		"sharedzones_maxsize":  1048575,
		"sharedzones_usednode": 13,
		"sharedzones_usedsize": 45799,
	}

	assert.Equal(t, expected, job.Collect())
}

func TestCollect_full(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testStatusDataFull)
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL

	require.True(t, job.Init())
	require.True(t, job.Check())

	expected := map[string]int64{
		"cachezones_foo_cache_inbytes":                       234,
		"cachezones_foo_cache_maxsize":                       1073741824,
		"cachezones_foo_cache_outbytes":                      1070,
		"cachezones_foo_cache_responses_bypass":              0,
		"cachezones_foo_cache_responses_expired":             0,
		"cachezones_foo_cache_responses_hit":                 0,
		"cachezones_foo_cache_responses_miss":                3,
		"cachezones_foo_cache_responses_revalidated":         0,
		"cachezones_foo_cache_responses_scarce":              0,
		"cachezones_foo_cache_responses_stale":               0,
		"cachezones_foo_cache_responses_updating":            0,
		"cachezones_foo_cache_usedsize":                      0,
		"connections_accepted":                               12,
		"connections_active":                                 2,
		"connections_handled":                                12,
		"connections_reading":                                0,
		"connections_requests":                               0,
		"connections_waiting":                                1,
		"connections_writing":                                1,
		"filterzones_b.bar.com_/login_inbytes":               78,
		"filterzones_b.bar.com_/login_outbytes":              314,
		"filterzones_b.bar.com_/login_requestcounter":        1,
		"filterzones_b.bar.com_/login_responses_1xx":         0,
		"filterzones_b.bar.com_/login_responses_2xx":         0,
		"filterzones_b.bar.com_/login_responses_3xx":         0,
		"filterzones_b.bar.com_/login_responses_4xx":         0,
		"filterzones_b.bar.com_/login_responses_5xx":         1,
		"filterzones_b.bar.com_/login_responses_bypass":      0,
		"filterzones_b.bar.com_/login_responses_expired":     0,
		"filterzones_b.bar.com_/login_responses_hit":         0,
		"filterzones_b.bar.com_/login_responses_miss":        0,
		"filterzones_b.bar.com_/login_responses_revalidated": 0,
		"filterzones_b.bar.com_/login_responses_scarce":      0,
		"filterzones_b.bar.com_/login_responses_stale":       0,
		"filterzones_b.bar.com_/login_responses_updating":    0,
		"loadmsec":                                              1606489796895,
		"nowmsec":                                               1606490116734,
		"serverzones_a.foo.com_inbytes":                         156,
		"serverzones_a.foo.com_outbytes":                        692,
		"serverzones_a.foo.com_requestcounter":                  2,
		"serverzones_a.foo.com_responses_1xx":                   0,
		"serverzones_a.foo.com_responses_2xx":                   1,
		"serverzones_a.foo.com_responses_3xx":                   0,
		"serverzones_a.foo.com_responses_4xx":                   0,
		"serverzones_a.foo.com_responses_5xx":                   1,
		"serverzones_a.foo.com_responses_bypass":                0,
		"serverzones_a.foo.com_responses_expired":               0,
		"serverzones_a.foo.com_responses_hit":                   0,
		"serverzones_a.foo.com_responses_miss":                  2,
		"serverzones_a.foo.com_responses_revalidated":           0,
		"serverzones_a.foo.com_responses_scarce":                0,
		"serverzones_a.foo.com_responses_stale":                 0,
		"serverzones_a.foo.com_responses_updating":              0,
		"sharedzones_maxsize":                                   1048575,
		"sharedzones_usednode":                                  13,
		"sharedzones_usedsize":                                  45799,
		"upstreamzones_backend_10.0.0.110:16666_inbytes":        78,
		"upstreamzones_backend_10.0.0.110:16666_outbytes":       314,
		"upstreamzones_backend_10.0.0.110:16666_requestcounter": 1,
		"upstreamzones_backend_10.0.0.110:16666_responses_1xx":  0,
		"upstreamzones_backend_10.0.0.110:16666_responses_2xx":  0,
		"upstreamzones_backend_10.0.0.110:16666_responses_3xx":  0,
		"upstreamzones_backend_10.0.0.110:16666_responses_4xx":  0,
		"upstreamzones_backend_10.0.0.110:16666_responses_5xx":  1,
		"upstreamzones_backend_10.0.0.110:26666_inbytes":        156,
		"upstreamzones_backend_10.0.0.110:26666_outbytes":       756,
		"upstreamzones_backend_10.0.0.110:26666_requestcounter": 2,
		"upstreamzones_backend_10.0.0.110:26666_responses_1xx":  0,
		"upstreamzones_backend_10.0.0.110:26666_responses_2xx":  2,
		"upstreamzones_backend_10.0.0.110:26666_responses_3xx":  0,
		"upstreamzones_backend_10.0.0.110:26666_responses_4xx":  0,
		"upstreamzones_backend_10.0.0.110:26666_responses_5xx":  0,
	}

	assert.Equal(t, expected, job.Collect())
}

func TestInvalidData(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("Hello World!"))
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL

	require.True(t, job.Init())
	assert.False(t, job.Check())
}

func TestVts404(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL

	require.True(t, job.Init())
	assert.False(t, job.Check())

}

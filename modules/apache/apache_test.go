package apache

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testExtendedStatusData, _ = ioutil.ReadFile("testdata/extended-status.txt")
	testSimpleStatusData, _   = ioutil.ReadFile("testdata/simple-status.txt")
	testLighttpdStatusData, _ = ioutil.ReadFile("testdata/lighttpd-status.txt")
)

func TestApache_Cleanup(t *testing.T) { New().Cleanup() }

func TestNew(t *testing.T) {
	job := New()

	assert.Implements(t, (*module.Module)(nil), job)
	assert.Equal(t, defaultURL, job.URL)
	assert.Equal(t, defaultHTTPTimeout, job.Timeout.Duration)
	assert.NotNil(t, job.charts)
}

func TestApache_Init(t *testing.T) {
	job := New()

	require.True(t, job.Init())
	assert.NotNil(t, job.apiClient)
}

func TestApache_InitNG(t *testing.T) {
	job := New()

	job.URL = ""
	assert.False(t, job.Init())
}

func TestApache_Check(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testSimpleStatusData)
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL + "/server-status?auto"
	require.True(t, job.Init())
	assert.True(t, job.Check())
}

func TestApache_CheckNG(t *testing.T) {
	job := New()

	job.URL = "http://127.0.0.1:38001/server-status?auto"
	require.True(t, job.Init())
	assert.False(t, job.Check())
}

func TestApache_Charts(t *testing.T) { assert.NotNil(t, New().Charts()) }

func TestApache_Collect(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testSimpleStatusData)
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL + "/server-status?auto"
	require.True(t, job.Init())
	require.True(t, job.Check())

	expected := map[string]int64{
		"conns_async_closing":     0,
		"scoreboard_waiting":      74,
		"scoreboard_reading":      0,
		"busy_workers":            1,
		"conns_async_writing":     0,
		"conns_async_keep_alive":  0,
		"scoreboard_starting":     0,
		"idle_workers":            74,
		"conns_total":             0,
		"scoreboard_keepalive":    0,
		"scoreboard_idle_cleanup": 0,
		"scoreboard_sending":      1,
		"scoreboard_dns_lookup":   0,
		"scoreboard_closing":      0,
		"scoreboard_logging":      0,
		"scoreboard_finishing":    0,
		"scoreboard_open":         325,
	}

	assert.Equal(t, expected, job.Collect())
}

func TestApache_CollectExtended(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testExtendedStatusData)
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL + "/server-status?auto"
	require.True(t, job.Init())
	require.True(t, job.Check())

	expected := map[string]int64{
		"conns_async_closing":     0,
		"scoreboard_open":         300,
		"conns_async_keep_alive":  0,
		"uptime":                  256,
		"req_per_sec":             3515,
		"bytes_per_sec":           4800000,
		"scoreboard_waiting":      99,
		"scoreboard_reading":      0,
		"scoreboard_idle_cleanup": 0,
		"total_accesses":          9,
		"total_kBytes":            12,
		"scoreboard_starting":     0,
		"scoreboard_logging":      0,
		"scoreboard_finishing":    0,
		"conns_total":             0,
		"idle_workers":            99,
		"conns_async_writing":     0,
		"scoreboard_sending":      1,
		"scoreboard_keepalive":    0,
		"scoreboard_dns_lookup":   0,
		"scoreboard_closing":      0,
		"busy_workers":            1,
		"bytes_per_req":           136533000,
	}

	assert.Equal(t, expected, job.Collect())
}

func TestApache_InvalidData(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("hello and goodbye"))
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL + "/server-status?auto"
	require.True(t, job.Init())
	assert.False(t, job.Check())
}

func TestApache_LighttpdData(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testLighttpdStatusData)
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL + "/server-status?auto"
	require.True(t, job.Init())
	require.False(t, job.Check())
}

func TestApache_404(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL + "/server-status?auto"
	require.True(t, job.Init())
	assert.False(t, job.Check())
}

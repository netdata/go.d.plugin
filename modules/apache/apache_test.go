package apache

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	extendedStatus, _ = ioutil.ReadFile("testdata/status-extended.txt")
	simpleStatus, _   = ioutil.ReadFile("testdata/status-simple.txt")
)

func TestApache_Cleanup(t *testing.T) {
	New().Cleanup()
}

func TestNew(t *testing.T) {
	mod := New()

	assert.Implements(t, (*modules.Module)(nil), mod)
	assert.Equal(t, defaultURL, mod.URL)
	assert.Equal(t, defaultHTTPTimeout, mod.Timeout.Duration)
}

func TestApache_Init(t *testing.T) {
	mod := New()

	require.True(t, mod.Init())
	assert.NotNil(t, mod.apiClient)
}

func TestApache_InitNG(t *testing.T) {
	mod := New()

	mod.HTTP.Request = web.Request{URL: ""}
	assert.False(t, mod.Init())
}

func TestApache_Check(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/server-status" {
					_, _ = w.Write(simpleStatus)
					return
				}
			}))

	defer ts.Close()

	mod := New()

	mod.HTTP.Request = web.Request{URL: ts.URL + "/server-status?auto"}
	require.True(t, mod.Init())
	assert.True(t, mod.Check())
}

func TestApache_CheckNG(t *testing.T) {
	mod := New()

	mod.HTTP.Request = web.Request{URL: "http://127.0.0.1:38001/server-status?auto"}
	require.True(t, mod.Init())
	assert.False(t, mod.Check())
}

func TestApache_Charts(t *testing.T) {
	mod := New()

	assert.NotNil(t, mod.Charts())
}

func TestApache_Collect(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/server-status" {
					_, _ = w.Write(extendedStatus)
					return
				}
			}))
	defer ts.Close()

	mod := New()
	mod.HTTP.Request = web.Request{URL: ts.URL + "/server-status?auto"}

	require.True(t, mod.Init())
	require.True(t, mod.Check())

	metrics := mod.Collect()

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

	assert.Equal(t, expected, metrics)
}

func TestApache_InvalidData(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/server-status" {
					_, _ = w.Write([]byte("hello and goodbye"))
					return
				}
			}))
	defer ts.Close()

	mod := New()
	mod.HTTP.Request = web.Request{URL: ts.URL + "/server-status?auto"}

	require.True(t, mod.Init())
	assert.False(t, mod.Check())
}

func TestApache_404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer ts.Close()

	mod := New()
	mod.HTTP.Request = web.Request{URL: ts.URL + "/server-status?auto"}

	require.True(t, mod.Init())
	assert.False(t, mod.Check())
}

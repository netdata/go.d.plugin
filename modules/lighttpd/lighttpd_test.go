package lighttpd

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
	statusData, _        = ioutil.ReadFile("testdata/status.txt")
	statusInvalidData, _ = ioutil.ReadFile("testdata/status-invalid.txt")
)

func TestLighttpd_Cleanup(t *testing.T) {
	New().Cleanup()
}

func TestNew(t *testing.T) {
	assert.Implements(t, (*modules.Module)(nil), New())
}

func TestLighttpd_Init(t *testing.T) {
	mod := New()

	assert.True(t, mod.Init())
	assert.NotNil(t, mod.request)
	assert.NotNil(t, mod.client)
}

func TestApache_InitNG(t *testing.T) {
	mod := New()

	mod.HTTP.Request = web.Request{URL: mod.Request.URL[0 : len(mod.Request.URL)-1]}
	assert.False(t, mod.Init())
}

func TestLighttpd_Check(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/server-status" {
					_, _ = w.Write(statusData)
					return
				}
			}))

	defer ts.Close()

	mod := New()
	mod.HTTP.Request = web.Request{URL: ts.URL + "/server-status?auto"}

	require.True(t, mod.Init())
	assert.True(t, mod.Check())
}

func TestLighttpd_CheckNG(t *testing.T) {
	mod := New()

	mod.HTTP.Request = web.Request{URL: "http://127.0.0.1:38001/server-status?auto"}

	require.True(t, mod.Init())
	assert.False(t, mod.Check())
}

func TestLighttpd_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
	assert.NoError(t, modules.CheckCharts(*New().Charts()...))
}

func TestLighttpd_Collect(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/server-status" {
					_, _ = w.Write(statusData)
					return
				}
			}))
	defer ts.Close()

	mod := New()
	mod.HTTP.Request = web.Request{URL: ts.URL + "/server-status?auto"}

	require.True(t, mod.Init())
	require.True(t, mod.Check())
	require.True(t, len(mod.Collect()) > 0)

	expected := map[string]int64{
		assign(totalAccesses):       6384,
		assign(totalkBytes):         2245,
		assign(busyServers):         5,
		assign(idleServers):         123,
		assign(uptime):              2137,
		"scoreboard_waiting":        123,
		"scoreboard_open":           0,
		"scoreboard_close":          0,
		"scoreboard_hard_error":     0,
		"scoreboard_keepalive":      0,
		"scoreboard_read":           4,
		"scoreboard_read_post":      0,
		"scoreboard_write":          0,
		"scoreboard_handle_request": 1,
		"scoreboard_request_start":  0,
		"scoreboard_request_end":    0,
		"scoreboard_response_start": 0,
		"scoreboard_response_end":   0,
	}

	assert.Equal(t, expected, mod.metrics)
}

func TestLighttpd_InvalidData(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/server-status" {
					_, _ = w.Write(statusInvalidData)
					return
				}
			}))
	defer ts.Close()

	mod := New()
	mod.HTTP.Request = web.Request{URL: ts.URL + "/server-status?auto"}

	require.True(t, mod.Init())
	assert.False(t, mod.Check())
}

func TestLighttpd_404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer ts.Close()

	mod := New()
	mod.HTTP.Request = web.Request{URL: ts.URL + "/server-status?auto"}

	require.True(t, mod.Init())
	assert.False(t, mod.Check())
}

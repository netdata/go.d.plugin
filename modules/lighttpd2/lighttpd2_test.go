package lighttpd2

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
	serverStatus, _ = ioutil.ReadFile("testdata/status.txt")
)

func TestLighttpd2_Cleanup(t *testing.T) {
	New().Cleanup()
}

func TestNew(t *testing.T) {
	mod := New()

	assert.Implements(t, (*modules.Module)(nil), mod)
	assert.Equal(t, defaultURL, mod.URL)
	assert.Equal(t, defaultHTTPTimeout, mod.Timeout.Duration)
}

func TestLighttpd2_Init(t *testing.T) {
	mod := New()

	require.True(t, mod.Init())
	assert.NotNil(t, mod.apiClient)
}

func TestLighttpd2_InitNG(t *testing.T) {
	mod := New()

	mod.HTTP.Request = web.Request{URL: ""}
	assert.False(t, mod.Init())
}

func TestLighttpd2_Check(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/server-status" {
					_, _ = w.Write(serverStatus)
					return
				}
			}))

	defer ts.Close()

	mod := New()

	mod.HTTP.Request = web.Request{URL: ts.URL + "/server-status?format=plain"}
	require.True(t, mod.Init())
	assert.True(t, mod.Check())
}

func TestLighttpd_CheckNG(t *testing.T) {
	mod := New()

	mod.HTTP.Request = web.Request{URL: "http://127.0.0.1:38001/server-status?format=plain"}
	require.True(t, mod.Init())
	assert.False(t, mod.Check())
}

func TestLighttpd_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestLighttpd_Collect(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/server-status" {
					_, _ = w.Write(serverStatus)
					return
				}
			}))
	defer ts.Close()

	mod := New()
	mod.HTTP.Request = web.Request{URL: ts.URL + "/server-status?format=plain"}

	require.True(t, mod.Init())
	require.True(t, mod.Check())

	expected := map[string]int64{
		"traffic_out_abs":                 8811598713,
		"status_1xx":                      0,
		"requests_abs":                    640866,
		"traffic_out_avg_5sec":            45142,
		"status_5xx":                      572,
		"traffic_out_avg":                 25183,
		"status_2xx":                      515698,
		"status_4xx":                      52891,
		"connections_abs":                 8,
		"connection_state_start":          0,
		"connection_state_read_header":    0,
		"connection_state_handle_request": 1,
		"connection_state_keep_alive":     7,
		"connection_state_upgraded":       0,
		"traffic_in_avg_5sec":             2938,
		"requests_avg":                    1,
		"connections_avg":                 4,
		"connection_state_write_response": 0,
		"memory_usage":                    39006208,
		"traffic_in_abs":                  292744726,
		"requests_avg_5sec":               7,
		"connections_avg_5sec":            1,
		"status_3xx":                      70456,
		"uptime":                          349894,
		"traffic_in_avg":                  836,
	}

	assert.Equal(t, expected, mod.Collect())
}

func TestLighttpd_InvalidData(t *testing.T) {
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
	mod.HTTP.Request = web.Request{URL: ts.URL + "/server-status?format=plain"}

	require.True(t, mod.Init())
	assert.False(t, mod.Check())
}

func TestLighttpd_404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer ts.Close()

	mod := New()
	mod.HTTP.Request = web.Request{URL: ts.URL + "/server-status?format=plain"}

	require.True(t, mod.Init())
	assert.False(t, mod.Check())
}

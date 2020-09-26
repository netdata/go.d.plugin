package lighttpd2

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
	testStatusData, _ = ioutil.ReadFile("testdata/status.txt")
)

func TestLighttpd2_Cleanup(t *testing.T) { New().Cleanup() }

func TestNew(t *testing.T) {
	job := New()

	assert.Implements(t, (*module.Module)(nil), job)
	assert.Equal(t, defaultURL, job.URL)
	assert.Equal(t, defaultHTTPTimeout, job.Timeout.Duration)
}

func TestLighttpd2_Init(t *testing.T) {
	job := New()

	require.True(t, job.Init())
	assert.NotNil(t, job.apiClient)
}

func TestLighttpd2_InitNG(t *testing.T) {
	job := New()

	job.URL = ""
	assert.False(t, job.Init())
}

func TestLighttpd2_Check(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testStatusData)
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL + "/server-status?format=plain"
	require.True(t, job.Init())
	assert.True(t, job.Check())
}

func TestLighttpd2_CheckNG(t *testing.T) {
	job := New()

	job.URL = "http://127.0.0.1:38001/server-status?format=plain"
	require.True(t, job.Init())
	assert.False(t, job.Check())
}

func TestLighttpd2_Charts(t *testing.T) { assert.NotNil(t, New().Charts()) }

func TestLighttpd2_Collect(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testStatusData)
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL + "/server-status?format=plain"
	require.True(t, job.Init())
	require.True(t, job.Check())

	expected := map[string]int64{
		"traffic_out_abs":                 8811598713,
		"traffic_in_abs":                  292744726,
		"status_1xx":                      0,
		"status_2xx":                      515698,
		"status_3xx":                      70456,
		"status_4xx":                      52891,
		"status_5xx":                      572,
		"requests_abs":                    640866,
		"connection_abs":                  8,
		"connection_state_start":          0,
		"connection_state_read_header":    0,
		"connection_state_handle_request": 1,
		"connection_state_write_response": 0,
		"connection_state_keepalive":      7,
		"connection_state_upgraded":       0,
		"memory_usage":                    39006208,
		"uptime":                          349894,
	}

	assert.Equal(t, expected, job.Collect())
}

func TestLighttpd2_InvalidData(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("hello and goodbye"))
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL + "/server-status?format=plain"
	require.True(t, job.Init())
	assert.False(t, job.Check())
}

func TestLighttpd2_404(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL + "/server-status?format=plain"
	require.True(t, job.Init())
	assert.False(t, job.Check())
}

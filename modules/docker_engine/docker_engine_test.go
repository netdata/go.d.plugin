package docker_engine

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testMetrics, _ = ioutil.ReadFile("testdata/metrics.txt")

func TestNew(t *testing.T) {
	job := New()

	assert.IsType(t, (*DockerEngine)(nil), job)
	assert.Equal(t, defaultURL, job.URL)
	assert.Equal(t, defaultHTTPTimeout, job.Timeout.Duration)
}

func TestDockerEngine_Charts(t *testing.T) { assert.NotNil(t, New().Charts()) }

func TestDockerEngine_Cleanup(t *testing.T) { New().Cleanup() }

func TestDockerEngine_Init(t *testing.T) { assert.True(t, New().Init()) }

func TestDockerEngine_InitNG(t *testing.T) {
	job := New()
	job.URL = ""
	assert.False(t, job.Init())
}

func TestDockerEngine_Check(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testMetrics)
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL + "/metrics"
	require.True(t, job.Init())
	assert.True(t, job.Check())
}

func TestDockerEngine_CheckNG(t *testing.T) {
	job := New()
	job.URL = "http://127.0.0.1:38001/metrics"
	require.True(t, job.Init())
	assert.False(t, job.Check())
}

func TestExample_Collect(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testMetrics)
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL + "/metrics"
	require.True(t, job.Init())
	require.True(t, job.Check())

	expected := map[string]int64{
		"actions_changes": 1,
		"actions_commit":  1,
		"actions_create":  1,
		"actions_delete":  1,
		"actions_start":   1,
		"states_paused":   11,
		"states_running":  12,
		"states_stopped":  13,
	}

	assert.Equal(t, expected, job.Collect())
}

func TestDockerEngine_InvalidData(t *testing.T) {
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

func TestDockerEngine_404(t *testing.T) {
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

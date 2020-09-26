package consul

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	checks, _ = ioutil.ReadFile("testdata/checks.txt")
)

func TestNew(t *testing.T) {
	job := New()

	assert.Implements(t, (*module.Module)(nil), New())
	assert.NotNil(t, job.charts)
	assert.NotNil(t, job.activeChecks)
	assert.Equal(t, defaultURL, job.URL)
	assert.Equal(t, defaultHTTPTimeout, job.Timeout.Duration)
	assert.Equal(t, defaultMaxChecks, job.MaxChecks)
}

func TestConsul_Init(t *testing.T) {
	job := New()

	assert.True(t, job.Init())
	assert.NotNil(t, job.apiClient)
}

func TestConsul_Check(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/v1/agent/checks" {
					_, _ = w.Write(checks)
					return
				}
			}))
	defer ts.Close()

	job := New()
	job.HTTP.Request = web.Request{URL: ts.URL}

	require.True(t, job.Init())
	assert.True(t, job.Check())
}

func TestConsul_CheckNG(t *testing.T) {
	job := New()

	job.HTTP.Request = web.Request{URL: "http://127.0.0.1:38001"}

	require.True(t, job.Init())
	assert.False(t, job.Check())
}

func TestConsul_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestConsul_Cleanup(t *testing.T) {
	New().Cleanup()
}

func TestConsul_Collect(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/v1/agent/checks" {
					_, _ = w.Write(checks)
					return
				}
			}))
	defer ts.Close()

	job := New()
	job.HTTP.Request = web.Request{URL: ts.URL}

	assert.True(t, job.Init())

	metrics := job.Collect()
	assert.NotNil(t, metrics)

	expected := map[string]int64{
		"chk3":  2,
		"mysql": 2,
		"chk1":  0,
		"chk2":  2,
	}

	assert.Equal(t, expected, metrics)
	assert.Len(t, job.charts.Get("service_checks").Dims, 1)
	assert.Len(t, job.charts.Get("unbound_checks").Dims, 3)
}

func TestConsul_InvalidData(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/v1/agent/checks" {
					_, _ = w.Write([]byte("farewell, ashen one"))
					return
				}
			}))
	defer ts.Close()

	job := New()
	job.HTTP.Request = web.Request{URL: ts.URL}

	require.True(t, job.Init())
	assert.False(t, job.Check())
}

func TestConsul_404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer ts.Close()

	job := New()
	job.HTTP.Request = web.Request{URL: ts.URL}

	require.True(t, job.Init())
	assert.False(t, job.Check())
}

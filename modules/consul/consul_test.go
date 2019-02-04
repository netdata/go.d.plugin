package consul

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	checks, _ = ioutil.ReadFile("testdata/checks.txt")
)

func TestNew(t *testing.T) {
	mod := New()

	assert.Implements(t, (*module.Module)(nil), New())
	assert.NotNil(t, mod.charts)
	assert.NotNil(t, mod.activeChecks)
	assert.Equal(t, defaultMaxChecks, mod.MaxChecks)
}

func TestConsul_Init(t *testing.T) {
	mod := New()

	assert.True(t, mod.Init())
	assert.NotNil(t, mod.apiClient)
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

	mod := New()
	mod.HTTP.Request = web.Request{URL: ts.URL}

	require.True(t, mod.Init())
	assert.True(t, mod.Check())
}

func TestConsul_CheckNG(t *testing.T) {
	mod := New()

	mod.HTTP.Request = web.Request{URL: "http://127.0.0.1:38001"}

	require.True(t, mod.Init())
	assert.False(t, mod.Check())
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

	mod := New()
	mod.HTTP.Request = web.Request{URL: ts.URL}

	assert.True(t, mod.Init())

	metrics := mod.Collect()
	assert.NotNil(t, metrics)

	expected := map[string]int64{
		"chk3":  2,
		"mysql": 2,
		"chk1":  0,
		"chk2":  2,
	}

	assert.Equal(t, expected, metrics)
	assert.Len(t, mod.charts.Get("service_checks").Dims, 1)
	assert.Len(t, mod.charts.Get("unbound_checks").Dims, 3)
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

	mod := New()
	mod.HTTP.Request = web.Request{URL: ts.URL}

	require.True(t, mod.Init())
	assert.False(t, mod.Check())
}

func TestConsul_404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer ts.Close()

	mod := New()
	mod.HTTP.Request = web.Request{URL: ts.URL}

	require.True(t, mod.Init())
	assert.False(t, mod.Check())
}

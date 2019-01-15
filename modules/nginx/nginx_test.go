package nginx

import (
	"github.com/netdata/go.d.plugin/modules"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/stretchr/testify/assert"
)

var (
	status, _ = ioutil.ReadFile("testdata/status.txt")
)

func TestNginx_Cleanup(t *testing.T) {
	New().Cleanup()
}

func TestNew(t *testing.T) {
	mod := New()

	assert.Implements(t, (*modules.Module)(nil), mod)
	assert.Equal(t, defaultURL, mod.URL)
	assert.Equal(t, defaultHTTPTimeout, mod.Timeout.Duration)
}

func TestNginx_Init(t *testing.T) {
	mod := New()

	require.True(t, mod.Init())
	assert.NotNil(t, mod.apiClient)
}

func TestNginx_InitNG(t *testing.T) {
	mod := New()

	mod.HTTP.Request = web.Request{URL: ""}
	assert.False(t, mod.Init())
}

func TestNginx_Check(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/stub_status" {
					_, _ = w.Write(status)
					return
				}
			}))

	defer ts.Close()

	mod := New()

	mod.HTTP.Request = web.Request{URL: ts.URL + "/stub_status"}
	require.True(t, mod.Init())
	assert.True(t, mod.Check())
}

func TestNginx_CheckNG(t *testing.T) {
	mod := New()

	mod.HTTP.Request = web.Request{URL: "http://127.0.0.1:38001/stub_status"}
	require.True(t, mod.Init())
	assert.False(t, mod.Check())
}

func TestNginx_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
	assert.NoError(t, modules.CheckCharts(*New().Charts()...))
}

func TestNginx_Collect(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/stub_status" {
					_, _ = w.Write(status)
					return
				}
			}))
	defer ts.Close()

	mod := New()
	mod.HTTP.Request = web.Request{URL: ts.URL + "/stub_status"}

	assert.True(t, mod.Init())

	metrics := mod.Collect()
	assert.NotNil(t, metrics)

	expected := map[string]int64{
		"active":   1,
		"accepts":  36,
		"handled":  36,
		"requests": 126,
		"reading":  0,
		"writing":  1,
		"waiting":  0,
	}

	assert.Equal(t, expected, metrics)
}

func TestNginx_InvalidData(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/stub_status" {
					_, _ = w.Write([]byte("hello and goodbye"))
					return
				}
			}))
	defer ts.Close()

	mod := New()
	mod.HTTP.Request = web.Request{URL: ts.URL + "/stub_status"}

	require.True(t, mod.Init())
	assert.False(t, mod.Check())
}

func TestNginx_404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer ts.Close()

	mod := New()
	mod.HTTP.Request = web.Request{URL: ts.URL + "/stub_status"}

	require.True(t, mod.Init())
	assert.False(t, mod.Check())
}

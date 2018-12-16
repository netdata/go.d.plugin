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
	invalidStatus, _  = ioutil.ReadFile("testdata/status-invalid.txt")
)

func TestApache_Cleanup(t *testing.T) {
	New().Cleanup()
}

func TestNew(t *testing.T) {
	assert.Implements(t, (*modules.Module)(nil), New())
}

func TestApache_Init(t *testing.T) {
	mod := New()

	require.True(t, mod.Init())
	assert.NotNil(t, mod.request)
	assert.NotNil(t, mod.client)
}

func TestApache_InitNG(t *testing.T) {
	mod := New()

	mod.HTTP.Request = web.Request{URL: mod.Request.URL[0 : len(mod.Request.URL)-1]}
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
	assert.NotNil(t, New().Charts())

	mod := New()
	mod.extendedStats = true

	assert.True(t, len(*mod.Charts()) > len(*New().Charts()))
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
	require.NotNil(t, mod.Collect())

	expected := map[string]int64{
		assign(totalAccesses):       575,
		assign(totalkBytes):         433,
		assign(reqPerSec):           101590,
		assign(bytesPerSec):         78337800,
		assign(bytesPerReq):         77111700,
		assign(busyWorkers):         1,
		assign(idleWorkers):         49,
		assign(connsTotal):          2,
		assign(connsAsyncWriting):   0,
		assign(connsAsyncKeepAlive): 2,
		assign(connsAsyncClosing):   0,
		"scoreboard_waiting":        49,
		"scoreboard_starting":       0,
		"scoreboard_reading":        0,
		"scoreboard_sending":        1,
		"scoreboard_keepalive":      0,
		"scoreboard_dns_lookup":     0,
		"scoreboard_closing":        0,
		"scoreboard_logging":        0,
		"scoreboard_finishing":      0,
		"scoreboard_idle_cleanup":   0,
		"scoreboard_open":           100,
	}

	assert.Equal(t, expected, mod.metrics)
}

func TestApache_InvalidData(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/server-status" {
					_, _ = w.Write(invalidStatus)
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

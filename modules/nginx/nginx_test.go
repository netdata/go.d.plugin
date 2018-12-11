package nginx

import (
	"github.com/netdata/go.d.plugin/pkg/web"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testdata = []byte(`Active connections: 1
server accepts handled requests
36 36 126
Reading: 0 Writing: 1 Waiting: 0
`)

func TestNew(t *testing.T) {
	assert.IsType(t, (*Nginx)(nil), New())
}

func TestNginx_Init(t *testing.T) {
	mod := New()

	assert.True(t, mod.Init())
	assert.NotNil(t, mod.request)
	assert.NotNil(t, mod.client)
	assert.NotZero(t, mod.Timeout.Duration)
}

func TestNginx_Check(t *testing.T) {
	mod := New()

	mod.Init()
	assert.False(t, mod.Check())
}

func TestNginx_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestNginx_Cleanup(t *testing.T) {
	New().Cleanup()
}

func TestNginx_GatherMetrics(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/stub_status" {
					_, _ = w.Write(testdata)
					return
				}
			}))
	defer ts.Close()

	mod := New()
	mod.HTTP.RawRequest = web.RawRequest{URL: ts.URL + "/stub_status"}

	assert.True(t, mod.Init())

	metrics := mod.GatherMetrics()
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

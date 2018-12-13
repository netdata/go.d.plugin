package rabbitmq

import (
	"github.com/netdata/go.d.plugin/pkg/web"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/netdata/go.d.plugin/modules"
)

var (
	overviewData, _ = ioutil.ReadFile("testdata/overviewData.txt")
	nodesData, _    = ioutil.ReadFile("testdata/nodesData.txt")
	badData, _      = ioutil.ReadFile("testdata/badData.txt")
)

func TestNew(t *testing.T) {
	assert.Implements(t, (*modules.Module)(nil), New())

}

func TestRabbitmq_Cleanup(t *testing.T) {
	New().Cleanup()
}

func TestRabbitmq_Init(t *testing.T) {

}

func TestRabbitmq_Check(t *testing.T) {

}

func TestRabbitmq_Charts(t *testing.T) {

}

func TestRabbitmq_GatherMetrics(t *testing.T) {

}

func TestRabbitmq_BadData(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(badData)
			}))
	defer ts.Close()

	mod := New()
	mod.HTTP.RawRequest = web.RawRequest{URL: ts.URL}

	assert.True(t, mod.Init())
	assert.False(t, mod.Check())
}

func TestRabbitmq_404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer ts.Close()

	mod := New()
	mod.HTTP.RawRequest = web.RawRequest{URL: ts.URL}

	require.True(t, mod.Init())
	assert.False(t, mod.Check())
}

package rabbitmq

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
	overviewData, _ = ioutil.ReadFile("testdata/overview.txt")
	nodesData, _    = ioutil.ReadFile("testdata/nodes.txt")
	badData, _      = ioutil.ReadFile("testdata/bad.txt")
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
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/api/overview":
					_, _ = w.Write(overviewData)
				case "/api/nodes":
					_, _ = w.Write(nodesData)
				}
			}))
	defer ts.Close()

	mod := New()
	mod.HTTP.RawRequest = web.RawRequest{URL: ts.URL}

	require.True(t, mod.Init())
	require.True(t, mod.Check())
	require.NotZero(t, mod.GatherMetrics())

	expected := map[string]int64{
		"message_stats_ack":                    1,
		"message_stats_publish":                2,
		"message_stats_publish_in":             3,
		"message_stats_publish_out":            4,
		"message_stats_confirm":                5,
		"message_stats_deliver":                6,
		"message_stats_deliver_no_ack":         7,
		"message_stats_get":                    8,
		"message_stats_get_no_ack":             9,
		"message_stats_deliver_get":            10,
		"message_stats_redeliver":              11,
		"message_stats_return_unroutable":      666,
		"object_totals_channels":               44,
		"object_totals_connections":            44,
		"object_totals_consumers":              65,
		"object_totals_exchanges":              43,
		"object_totals_queues":                 62,
		"queue_totals_messages_ready":          150,
		"queue_totals_messages_unacknowledged": 99,
		"fd_used":                              74,
		"sockets_used":                         40,
		"mem_used":                             95463672,
		"disk_free":                            79654567936,
		"proc_used":                            621,
		"run_queue":                            0,
	}

	assert.Equal(t, expected, mod.metrics)
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

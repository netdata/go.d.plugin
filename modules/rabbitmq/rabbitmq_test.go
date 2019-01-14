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
	nodeData, _     = ioutil.ReadFile("testdata/node.txt")
)

func TestRabbitmq_Cleanup(t *testing.T) {
	New().Cleanup()
}

func TestNew(t *testing.T) {
	assert.Implements(t, (*modules.Module)(nil), New())

}

func TestRabbitmq_Init(t *testing.T) {
	mod := New()

	assert.Implements(t, (*modules.Module)(nil), mod)
	assert.Equal(t, defURL, mod.URL)
	assert.Equal(t, defUsername, mod.Username)
	assert.Equal(t, defPassword, mod.Password)
	assert.Equal(t, defHTTPTimeout, mod.Timeout.Duration)
}

func TestRabbitmq_Check(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/api/overview":
					_, _ = w.Write(overviewData)
				case "/api/node/rabbit@rbt0":
					_, _ = w.Write(nodeData)
				}
			}))
	defer ts.Close()

	mod := New()
	mod.HTTP.Request = web.Request{URL: ts.URL}

	require.True(t, mod.Init())
	assert.True(t, mod.Check())
}

func TestApache_CheckNG(t *testing.T) {
	mod := New()

	require.True(t, mod.Init())
	assert.False(t, mod.Check())
}

func TestRabbitmq_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestRabbitmq_Collect(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/api/overview":
					_, _ = w.Write(overviewData)
				case "/api/node/rabbit@rbt0":
					_, _ = w.Write(nodeData)
				}
			}))
	defer ts.Close()

	mod := New()
	mod.HTTP.Request = web.Request{URL: ts.URL}

	require.True(t, mod.Init())
	require.True(t, mod.Check())

	expected := map[string]int64{
		"message_stats_deliver_no_ack":         7,
		"run_queue":                            0,
		"mem_used":                             75022616,
		"message_stats_return_unroutable":      666,
		"message_stats_get":                    8,
		"object_totals_consumers":              65,
		"object_totals_queues":                 62,
		"message_stats_deliver":                6,
		"message_stats_deliver_get":            10,
		"message_stats_confirm":                5,
		"message_stats_publish_in":             3,
		"message_stats_publish":                2,
		"message_stats_get_no_ack":             9,
		"message_stats_redeliver":              11,
		"sockets_used":                         40,
		"fd_used":                              75,
		"object_totals_channels":               44,
		"message_stats_ack":                    1,
		"object_totals_exchanges":              43,
		"proc_used":                            622,
		"disk_free":                            79493152768,
		"queue_totals_messages_unacknowledged": 99,
		"message_stats_publish_out":            4,
		"object_totals_connections":            44,
		"queue_totals_messages_ready":          150,
	}

	assert.Equal(t, expected, mod.Collect())
}

func TestRabbitmq_InvalidData(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("hello and goodbye"))
			}))
	defer ts.Close()

	mod := New()
	mod.HTTP.Request = web.Request{URL: ts.URL}

	assert.True(t, mod.Init())
	assert.False(t, mod.Check())
}

func TestRabbitmq_404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer ts.Close()

	mod := New()
	mod.HTTP.Request = web.Request{URL: ts.URL}

	require.True(t, mod.Init())
	assert.False(t, mod.Check())
}

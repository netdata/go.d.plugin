// SPDX-License-Identifier: GPL-3.0-or-later

package rabbitmq

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testOverviewData, _ = os.ReadFile("testdata/overview.json")
	testNodeData, _     = os.ReadFile("testdata/node.json")
	testVhostsData, _   = os.ReadFile("testdata/vhosts.json")
)

func newTestRabbitMQHTTPServer() *httptest.Server {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				default:
					w.WriteHeader(404)
				case "/api/overview":
					_, _ = w.Write(testOverviewData)
				case "/api/nodes/rabbit@rbt0":
					_, _ = w.Write(testNodeData)
				case "/api/vhosts":
					_, _ = w.Write(testVhostsData)
				}
			}))
	return ts
}

func Test_readTestData(t *testing.T) {
	assert.NotNil(t, testOverviewData)
	assert.NotNil(t, testNodeData)
	assert.NotNil(t, testVhostsData)
}

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}

func TestRabbitMQ_Cleanup(t *testing.T) {
	New().Cleanup()
}

func TestRabbitMQ_Init(t *testing.T) {
	job := New()

	assert.True(t, job.Init())
}

func TestRabbitMQ_InitErrorOnCreatingClientWrongTLSCA(t *testing.T) {
	job := New()
	job.Client.TLSConfig.TLSCA = "testdata/tls"

	assert.False(t, job.Init())
}

func TestRabbitMQ_Check(t *testing.T) {
	ts := newTestRabbitMQHTTPServer()
	defer ts.Close()

	job := New()
	job.URL = ts.URL
	require.True(t, job.Init())

	assert.True(t, job.Check())
}

func TestHDFS_CheckError404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL
	require.True(t, job.Init())

	assert.False(t, job.Check())
}

func TestRabbitMQ_CheckNoResponse(t *testing.T) {
	job := New()
	job.URL = "http://127.0.0.1:38001"
	require.True(t, job.Init())

	assert.False(t, job.Check())
}

func TestRabbitMQ_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestRabbitMQ_Collect(t *testing.T) {
	ts := newTestRabbitMQHTTPServer()
	defer ts.Close()
	job := New()
	job.URL = ts.URL
	require.True(t, job.Init())
	require.True(t, job.Check())

	expected := map[string]int64{
		"disk_free":                                         79493152768,
		"fd_used":                                           75,
		"mem_used":                                          75022616,
		"message_stats_ack":                                 1,
		"message_stats_confirm":                             5,
		"message_stats_deliver":                             6,
		"message_stats_deliver_get":                         10,
		"message_stats_deliver_no_ack":                      7,
		"message_stats_get":                                 8,
		"message_stats_get_no_ack":                          9,
		"message_stats_publish":                             2,
		"message_stats_publish_in":                          3,
		"message_stats_publish_out":                         4,
		"message_stats_redeliver":                           11,
		"message_stats_return_unroutable":                   666,
		"object_totals_channels":                            44,
		"object_totals_connections":                         44,
		"object_totals_consumers":                           65,
		"object_totals_exchanges":                           43,
		"object_totals_queues":                              62,
		"proc_used":                                         622,
		"queue_totals_messages_ready":                       150,
		"queue_totals_messages_unacknowledged":              99,
		"run_queue":                                         0,
		"sockets_used":                                      40,
		"vhost_/check_api_message_stats_ack":                208961440,
		"vhost_/check_api_message_stats_confirm":            210205428,
		"vhost_/check_api_message_stats_deliver":            209220446,
		"vhost_/check_api_message_stats_deliver_get":        209220446,
		"vhost_/check_api_message_stats_deliver_no_ack":     0,
		"vhost_/check_api_message_stats_get":                0,
		"vhost_/check_api_message_stats_get_no_ack":         0,
		"vhost_/check_api_message_stats_publish":            209597605,
		"vhost_/check_api_message_stats_publish_in":         0,
		"vhost_/check_api_message_stats_publish_out":        0,
		"vhost_/check_api_message_stats_redeliver":          210205428,
		"vhost_/check_api_message_stats_return_unroutable":  210205428,
		"vhost_/search_api_message_stats_ack":               210205368,
		"vhost_/search_api_message_stats_confirm":           174130170,
		"vhost_/search_api_message_stats_deliver":           210205428,
		"vhost_/search_api_message_stats_deliver_get":       210205428,
		"vhost_/search_api_message_stats_deliver_no_ack":    210205428,
		"vhost_/search_api_message_stats_get":               210205428,
		"vhost_/search_api_message_stats_get_no_ack":        210205428,
		"vhost_/search_api_message_stats_publish":           210127507,
		"vhost_/search_api_message_stats_publish_in":        0,
		"vhost_/search_api_message_stats_publish_out":       0,
		"vhost_/search_api_message_stats_redeliver":         60,
		"vhost_/search_api_message_stats_return_unroutable": 210205428,
	}

	assert.Equal(t, expected, job.Collect())
}

func TestRabbitMQ_AddVhostsChartsAfterCollect(t *testing.T) {
	ts := newTestRabbitMQHTTPServer()
	defer ts.Close()
	job := New()
	job.URL = ts.URL
	require.True(t, job.Init())
	require.True(t, job.Check())
	require.NotNil(t, job.Collect())

	assert.True(t, job.charts.Has("vhost_/search_api_message_stats"))
	assert.True(t, job.charts.Has("vhost_/check_api_message_stats"))
}

func TestRabbitMQ_CollectReceiveNoResponse(t *testing.T) {
	job := New()
	job.URL = "http://127.0.0.1:38001/jmx"
	require.True(t, job.Init())

	assert.Nil(t, job.Collect())
}

func TestRabbitMQ_CollectReceiveUnexpectedJSONResponse(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(`{"ByteSlice":"AAAAAQID","SingleByte":10,"IntSlice":[0,0,0,1,2,3]}`))
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL
	assert.True(t, job.Init())

	assert.Nil(t, job.Collect())
}

func TestRabbitMQ_CollectReceiveNotJSONResponse(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("hello and goodbye"))
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL
	assert.True(t, job.Init())

	assert.Nil(t, job.Collect())
}

func TestRabbitMQ_CollectReceive404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL
	require.True(t, job.Init())

	assert.Nil(t, job.Collect())
}

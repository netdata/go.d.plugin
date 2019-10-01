package rabbitmq

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/netdata/go-orchestrator/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testOverviewData, _ = ioutil.ReadFile("testdata/overview.json")
	testNodeData, _     = ioutil.ReadFile("testdata/node.json")
	testVhostsData, _   = ioutil.ReadFile("testdata/vhosts.json")
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
				case "/api/node/rabbit@rbt0":
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

func TestRabbitMQ_Init(t *testing.T) {
	job := New()

	assert.True(t, job.Init())
}

func TestRabbitMQ_InitErrorOnCreatingClientWrongTLSCA(t *testing.T) {
	job := New()
	job.ClientTLSConfig.TLSCA = "testdata/tls"

	assert.False(t, job.Init())
}

func TestHDFS_Check(t *testing.T) {
	ts := newTestRabbitMQHTTPServer()
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL
	require.True(t, job.Init())

	assert.True(t, job.Check())
}

func TestRabbitMQ_Cleanup(t *testing.T) {
	New().Cleanup()
}

func TestRabbitMQ_InvalidData(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("hello and goodbye"))
			}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL

	assert.True(t, job.Init())
	assert.False(t, job.Check())
}

func TestRabbitMQ_404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL

	require.True(t, job.Init())
	assert.False(t, job.Check())
}

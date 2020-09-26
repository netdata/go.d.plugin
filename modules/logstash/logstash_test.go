package logstash

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	jvmStatusData, _ = ioutil.ReadFile("testdata/stats.json")
)

func Test_readTestData(t *testing.T) {
	assert.NotNil(t, jvmStatusData)
}

func TestLogstash_Cleanup(t *testing.T) {
	New().Cleanup()
}

func TestNew(t *testing.T) {
	logstash := New()

	assert.Implements(t, (*module.Module)(nil), logstash)
}

func TestLogstash_Init(t *testing.T) {
	logstash := New()

	assert.True(t, logstash.Init())
}

func TestLogstash_Init_ErrorOnValidatingConfigURLNotSet(t *testing.T) {
	logstash := New()
	logstash.URL = ""

	assert.False(t, logstash.Init())
}

func TestWMI_Init_ErrorOnCreatingClientWrongTLSCA(t *testing.T) {
	logstash := New()
	logstash.Client.TLSConfig.TLSCA = "testdata/tls"

	assert.False(t, logstash.Init())
}

func TestLogstash_Check(t *testing.T) {
	logstash, ts := prepareClientServerValidResponse(t)
	defer ts.Close()

	assert.True(t, logstash.Check())
}

func TestWMI_Check_ErrorOnCollectConnectionRefused(t *testing.T) {
	logstash := New()
	logstash.URL = "http://127.0.0.1:38001/metrics"
	require.True(t, logstash.Init())

	assert.False(t, logstash.Check())
}

func TestLogstash_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestLogstash_Collect(t *testing.T) {
	logstash, ts := prepareClientServerValidResponse(t)
	defer ts.Close()

	expected := map[string]int64{
		"event_duration_in_millis":                                 0,
		"event_filtered":                                           0,
		"event_in":                                                 0,
		"event_out":                                                0,
		"event_queue_push_duration_in_millis":                      0,
		"jvm_gc_collectors_eden_collection_count":                  5796,
		"jvm_gc_collectors_eden_collection_time_in_millis":         45008,
		"jvm_gc_collectors_old_collection_count":                   7,
		"jvm_gc_collectors_old_collection_time_in_millis":          3263,
		"jvm_mem_heap_committed_in_bytes":                          528154624,
		"jvm_mem_heap_used_in_bytes":                               189973480,
		"jvm_mem_heap_used_percent":                                35,
		"jvm_mem_pools_eden_committed_in_bytes":                    69795840,
		"jvm_mem_pools_eden_used_in_bytes":                         2600120,
		"jvm_mem_pools_old_committed_in_bytes":                     449642496,
		"jvm_mem_pools_old_used_in_bytes":                          185944824,
		"jvm_mem_pools_survivor_committed_in_bytes":                8716288,
		"jvm_mem_pools_survivor_used_in_bytes":                     1428536,
		"jvm_threads_count":                                        28,
		"jvm_uptime_in_millis":                                     699809475,
		"pipelines_pipeline-1_event_duration_in_millis":            5027018,
		"pipelines_pipeline-1_event_filtered":                      567639,
		"pipelines_pipeline-1_event_in":                            567639,
		"pipelines_pipeline-1_event_out":                           567639,
		"pipelines_pipeline-1_event_queue_push_duration_in_millis": 84241,
		"process_open_file_descriptors":                            101,
	}

	collected := logstash.Collect()
	assert.Equal(t, expected, collected)
	testCharts(t, logstash, collected)
}

func TestLogstash_Collect_ReturnsNothingWhenBadData(t *testing.T) {
	logstash, ts := prepareClientServerBadData(t)
	defer ts.Close()

	assert.Nil(t, logstash.Collect())
}

func TestLogstash_Collect_ReturnsNothingWhen404(t *testing.T) {
	logstash, ts := prepareClientServerResponse404(t)
	defer ts.Close()

	assert.Nil(t, logstash.Collect())
}

func testCharts(t *testing.T, logstash *Logstash, collected map[string]int64) {
	ensurePipelinesChartsCreated(t, logstash)
	ensureCollectedHasAllChartsDimsVarsIDs(t, logstash, collected)
}

func ensurePipelinesChartsCreated(t *testing.T, logstash *Logstash) {
	for id := range logstash.collectedPipelines {
		for _, chart := range *pipelineCharts(id) {
			assert.Truef(t, logstash.Charts().Has(chart.ID), "chart '%' is not created", chart.ID)
		}
	}
}

func ensureCollectedHasAllChartsDimsVarsIDs(t *testing.T, logstash *Logstash, collected map[string]int64) {
	for _, chart := range *logstash.Charts() {
		for _, dim := range chart.Dims {
			_, ok := collected[dim.ID]
			assert.Truef(t, ok, "collected metrics has no data for dim '%s' chart '%s'", dim.ID, chart.ID)
		}
		for _, v := range chart.Vars {
			_, ok := collected[v.ID]
			assert.Truef(t, ok, "collected metrics has no data for var '%s' chart '%s'", v.ID, chart.ID)
		}
	}
}

func prepareClientServerValidResponse(t *testing.T) (*Logstash, *httptest.Server) {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(jvmStatusData)
		}))

	logstash := New()
	logstash.URL = ts.URL
	require.True(t, logstash.Init())
	return logstash, ts
}

func prepareClientServerBadData(t *testing.T) (*Logstash, *httptest.Server) {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("hello and\n goodbye"))
		}))

	logstash := New()
	logstash.URL = ts.URL
	require.True(t, logstash.Init())
	return logstash, ts
}

func prepareClientServerResponse404(t *testing.T) (*Logstash, *httptest.Server) {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

	logstash := New()
	logstash.URL = ts.URL
	require.True(t, logstash.Init())
	return logstash, ts
}

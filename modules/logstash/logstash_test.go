package logstash

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
	jvmStatusData, _ = ioutil.ReadFile("testdata/stats.json")
)

func TestLogstash_Cleanup(t *testing.T) {
	New().Cleanup()
}

func TestNew(t *testing.T) {
	job := New()

	assert.Implements(t, (*module.Module)(nil), job)
}

func TestLogstash_Init(t *testing.T) {
	job := New()

	require.True(t, job.Init())
	assert.NotNil(t, job.apiClient)
}

func TestLogstash_InitNG(t *testing.T) {
	job := New()

	job.HTTP.Request = web.Request{UserURL: ""}
	assert.False(t, job.Init())
}

func TestLogstash_Check(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == jvmStatusPath {
					_, _ = w.Write(jvmStatusData)
					return
				}
			}))

	defer ts.Close()

	job := New()

	job.UserURL = ts.URL
	require.True(t, job.Init())
	assert.True(t, job.Check())
}

func TestLogstash_CheckNG(t *testing.T) {
	job := New()

	job.UserURL = "http://127.0.0.1:38001"
	require.True(t, job.Init())
	assert.False(t, job.Check())
}

func TestLogstash_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestLogstash_Collect(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == jvmStatusPath {
					_, _ = w.Write(jvmStatusData)
					return
				}
			}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL

	assert.True(t, job.Init())

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

	assert.Equal(t, expected, job.Collect())
}

func TestLogstash_InvalidData(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == jvmStatusPath {
					_, _ = w.Write([]byte("hello and goodbye"))
					return
				}
			}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL

	require.True(t, job.Init())
	assert.False(t, job.Check())
}

func TestLogstash_404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL

	require.True(t, job.Init())
	assert.False(t, job.Check())
}

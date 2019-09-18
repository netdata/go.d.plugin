package hdfs

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/netdata/go-orchestrator/module"
)

var (
	testDataNodeData, _ = ioutil.ReadFile("testdata/datanode.json")
)

func Test_readTestData(t *testing.T) {
	assert.NotNil(t, testDataNodeData)
}

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}

func TestHDFS_Init(t *testing.T) {
	job := New()

	assert.True(t, job.Init())
}

func TestHDFS_InitErrorOnCreatingClientWrongTLSCA(t *testing.T) {
	job := New()
	job.ClientTLSConfig.TLSCA = "testdata/tls"

	assert.False(t, job.Init())
}

func TestHDFS_Check(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testDataNodeData)
			}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL
	require.True(t, job.Init())

	assert.True(t, job.Check())
}

func TestHDFS_CheckNoResponse(t *testing.T) {
	job := New()
	job.UserURL = "http://127.0.0.1:38001/jmx"
	require.True(t, job.Init())

	assert.False(t, job.Check())
}

func TestHDFS_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestHDFS_Cleanup(t *testing.T) {
	New().Cleanup()
}

func TestHDFS_Collect(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testDataNodeData)
			}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL
	require.True(t, job.Init())
	require.True(t, job.Check())

	expected := map[string]int64{
		"jvm_gc_count":                       155,
		"jvm_gc_num_info_threshold_exceeded": 0,
		"jvm_gc_num_warn_threshold_exceeded": 0,
		"jvm_gc_time_millis":                 672,
		"jvm_gc_total_extra_sleep_time":      8783,
		"jvm_log_error":                      1,
		"jvm_log_fatal":                      0,
		"jvm_log_info":                       257,
		"jvm_log_warn":                       2,
		"jvm_mem_heap_committed":             60,
		"jvm_mem_heap_max":                   843,
		"jvm_mem_heap_used":                  18,
		"jvm_mem_max":                        843,
		"jvm_mem_non_heap_committed":         54,
		"jvm_mem_non_heap_max":               -1,
		"jvm_mem_non_heap_used":              53,
		"jvm_threads_blocked":                0,
		"jvm_threads_new":                    0,
		"jvm_threads_runnable":               11,
		"jvm_threads_terminated":             0,
		"jvm_threads_timed_waiting":          25,
		"jvm_threads_waiting":                11,
	}

	assert.Equal(t, expected, job.Collect())
}

func TestHDFS_CollectNoResponse(t *testing.T) {
	job := New()
	job.UserURL = "http://127.0.0.1:38001/jmx"
	require.True(t, job.Init())

	assert.Nil(t, job.Collect())
}

func TestHDFS_CollectReceiveInvalidResponse(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("hello and\ngoodbye!\n"))
			}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL
	require.True(t, job.Init())

	assert.Nil(t, job.Collect())
}

func TestHDFS_CollectReceive404(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL
	require.True(t, job.Init())

	assert.Nil(t, job.Collect())
}

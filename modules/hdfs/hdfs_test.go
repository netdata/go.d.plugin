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
	testJvmData, _ = ioutil.ReadFile("testdata/jvm.json")
)

func Test_readTestData(t *testing.T) {
	assert.NotNil(t, testJvmData)
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
				_, _ = w.Write(testJvmData)
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
				_, _ = w.Write(testJvmData)
			}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL
	require.True(t, job.Init())
	require.True(t, job.Check())

	expected := map[string]int64{
		"jvm_gc_count":                       1699,
		"jvm_gc_num_info_threshold_exceeded": 0,
		"jvm_gc_num_warn_threshold_exceeded": 0,
		"jvm_gc_time_millis":                 3483,
		"jvm_gc_total_extra_sleep_time":      1944,
		"jvm_log_error":                      0,
		"jvm_log_fatal":                      0,
		"jvm_log_info":                       3382077,
		"jvm_log_warn":                       3378983,
		"jvm_mem_heap_committed":             67,
		"jvm_mem_heap_max":                   843,
		"jvm_mem_heap_used":                  26,
		"jvm_mem_max":                        843,
		"jvm_mem_non_heap_committed":         67,
		"jvm_mem_non_heap_max":               -1,
		"jvm_mem_non_heap_used":              66,
		"jvm_threads_blocked":                0,
		"jvm_threads_new":                    0,
		"jvm_threads_runnable":               7,
		"jvm_threads_terminated":             0,
		"jvm_threads_timed_waiting":          34,
		"jvm_threads_waiting":                6,
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

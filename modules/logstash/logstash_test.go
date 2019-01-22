package logstash

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
	jvmStatusData, _ = ioutil.ReadFile("testdata/jvm-stats.txt")
)

func TestLogstash_Cleanup(t *testing.T) {
	New().Cleanup()
}

func TestNew(t *testing.T) {
	job := New()

	assert.Implements(t, (*modules.Module)(nil), job)
	assert.Equal(t, defaultURL, job.URL)
	assert.Equal(t, defaultHTTPTimeout, job.Timeout.Duration)
}

func TestLogstash_Init(t *testing.T) {
	job := New()

	require.True(t, job.Init())
	assert.NotNil(t, job.apiClient)
}

func TestLogstash_InitNG(t *testing.T) {
	job := New()

	job.HTTP.Request = web.Request{URL: ""}
	assert.False(t, job.Init())
}

func TestLogstash_Check(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == jvmStatusURI {
					_, _ = w.Write(jvmStatusData)
					return
				}
			}))

	defer ts.Close()

	job := New()

	job.HTTP.Request = web.Request{URL: ts.URL}
	require.True(t, job.Init())
	assert.True(t, job.Check())
}

func TestLogstash_CheckNG(t *testing.T) {
	job := New()

	job.HTTP.Request = web.Request{URL: "http://127.0.0.1:38001"}
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
				if r.URL.Path == jvmStatusURI {
					_, _ = w.Write(jvmStatusData)
					return
				}
			}))
	defer ts.Close()

	job := New()
	job.HTTP.Request = web.Request{URL: ts.URL}

	assert.True(t, job.Init())

	expected := map[string]int64{
		"jvm_mem_heap_used_percent":                         14,
		"jvm_mem_pools_survivor_used_in_bytes":              288776,
		"jvm_mem_pools_old_used_in_bytes":                   148656848,
		"jvm_mem_pools_old_committed_in_bytes":              229322752,
		"jvm_mem_pools_young_committed_in_bytes":            71630848,
		"jvm_gc_collectors_old_collection_count":            12,
		"jvm_gc_collectors_young_collection_count":          1033,
		"jvm_mem_heap_committed_in_bytes":                   309866496,
		"jvm_mem_heap_used_in_bytes":                        151686096,
		"jvm_mem_pools_survivor_committed_in_bytes":         8912896,
		"jvm_mem_pools_young_used_in_bytes":                 2740472,
		"jvm_gc_collectors_old_collection_time_in_millis":   607,
		"jvm_gc_collectors_young_collection_time_in_millis": 4904,
	}

	assert.Equal(t, expected, job.Collect())
}

func TestLogstash_InvalidData(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == jvmStatusURI {
					_, _ = w.Write([]byte("hello and goodbye"))
					return
				}
			}))
	defer ts.Close()

	job := New()
	job.HTTP.Request = web.Request{URL: ts.URL}

	require.True(t, job.Init())
	assert.False(t, job.Check())
}

func TestLogstash_404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer ts.Close()

	job := New()
	job.HTTP.Request = web.Request{URL: ts.URL}

	require.True(t, job.Init())
	assert.False(t, job.Check())
}

package springboot2

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testdata, _ = ioutil.ReadFile("tests/testdata.txt")

func TestSpringboot2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/actuator/prometheus" {
			_, _ = w.Write(testdata)
			return
		}
	}))
	defer ts.Close()
	job := New()
	job.HTTP.Request.URL = ts.URL + "/actuator/prometheus"

	assert.True(t, job.Init())

	assert.True(t, job.Check())

	data := job.Collect()

	assert.EqualValues(
		t,
		map[string]int64{
			"threads":                 23,
			"threads_daemon":          21,
			"resp_1xx":                1,
			"resp_2xx":                19,
			"resp_3xx":                1,
			"resp_4xx":                4,
			"resp_5xx":                1,
			"heap_used_eden":          129649936,
			"heap_used_survivor":      8900136,
			"heap_used_old":           17827920,
			"heap_committed_eden":     153616384,
			"heap_committed_survivor": 8912896,
			"heap_committed_old":      40894464,
			"mem_free":                47045752,
			"uptime":                  191730,
		},
		data,
	)
}

func TestSpringboot2_404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer ts.Close()
	job := New()
	job.HTTP.Request.URL = ts.URL + "/actuator/prometheus"

	job.Init()

	assert.False(t, job.Check())

	job.Cleanup()
}

func TestSpringBoot2_Charts(t *testing.T) {
	job := New()
	charts := job.Charts()

	assert.True(t, charts.Has("response_codes"))
	assert.True(t, charts.Has("uptime"))
}

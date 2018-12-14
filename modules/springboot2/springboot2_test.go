package springboot2

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/netdata/go.d.plugin/pkg/web"

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
	module := &SpringBoot2{
		HTTP: web.HTTP{
			RawRequest: web.RawRequest{
				URL: ts.URL + "/actuator/prometheus",
			},
		},
	}

	assert.True(t, module.Init())

	assert.True(t, module.Check())

	data := module.Collect()

	assert.EqualValues(
		t,
		map[string]int64{
			"threads":                 23,
			"threads_daemon":          21,
			"resp_1xx":                0,
			"resp_2xx":                19,
			"resp_3xx":                0,
			"resp_4xx":                4,
			"resp_5xx":                0,
			"heap_used_eden":          129649936,
			"heap_used_survivor":      8900136,
			"heap_used_old":           17827920,
			"heap_committed_eden":     153616384,
			"heap_committed_survivor": 8912896,
			"heap_committed_old":      40894464,
			"mem_free":                47045752,
		},
		data,
	)
}

func TestSpringboot2_404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer ts.Close()
	module := &SpringBoot2{
		HTTP: web.HTTP{
			RawRequest: web.RawRequest{
				URL: ts.URL + "/actuator/prometheus",
			},
		},
	}

	module.Init()

	assert.False(t, module.Check())
}

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

	data := module.GatherMetrics()

	assert.EqualValues(
		t,
		map[string]int64{"threads": 24, "threads_daemon": 20},
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

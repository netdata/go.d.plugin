package prometheus

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/netdata/go.d.plugin/pkg/web"
	"github.com/prometheus/prometheus/pkg/labels"

	"github.com/stretchr/testify/assert"
)

var testdata, _ = ioutil.ReadFile("tests/testdata.txt")
var testdataNometa, _ = ioutil.ReadFile("tests/testdata.nometa.txt")

func TestPrometheus404(t *testing.T) {
	tsMux := http.NewServeMux()
	tsMux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})
	ts := httptest.NewServer(tsMux)
	defer ts.Close()

	req := web.Request{UserURL: ts.URL + "/metrics"}
	prom := New(http.DefaultClient, req)
	res, err := prom.Scrape()

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestPrometheusPlain(t *testing.T) {
	tsMux := http.NewServeMux()
	tsMux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(testdata)
	})
	ts := httptest.NewServer(tsMux)
	defer ts.Close()

	req := web.Request{UserURL: ts.URL + "/metrics"}
	prom := New(http.DefaultClient, req)
	res, err := prom.Scrape()

	assert.NoError(t, err)
	verifyTestData(t, res)
}

func TestPrometheusGzip(t *testing.T) {
	counter := 0
	rawTestData := [][]byte{testdata, testdataNometa}
	tsMux := http.NewServeMux()
	tsMux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(200)
		gz := new(bytes.Buffer)
		ww := gzip.NewWriter(gz)
		_, _ = ww.Write(rawTestData[counter])
		_ = ww.Close()
		_, _ = gz.WriteTo(w)
		counter++
	})
	ts := httptest.NewServer(tsMux)
	defer ts.Close()

	req := web.Request{UserURL: ts.URL + "/metrics"}
	prom := New(http.DefaultClient, req)

	for i := 0; i < 2; i++ {
		res, err := prom.Scrape()
		assert.NoError(t, err)
		verifyTestData(t, res)
	}
}

func TestParse(t *testing.T) {
	res := Metrics{}
	err := parse(testdata, &res)
	assert.NoError(t, err)

	verifyTestData(t, res)
}

func verifyTestData(t *testing.T, m Metrics) {
	assert.Equal(t, 410, len(m))
	assert.Equal(t, "go_gc_duration_seconds", m[0].Labels.Get("__name__"))
	assert.Equal(t, "0", m[0].Labels.Get("quantile"))
	assert.InDelta(t, 4.9351e-05, m[0].Value, 0.0001)

	notExistYet := m.FindByName("not_exist_yet")
	assert.NotNil(t, notExistYet)
	assert.Len(t, notExistYet, 0)

	targetInterval := m.FindByName("prometheus_target_interval_length_seconds")
	assert.Len(t, targetInterval, 5)

	matcher, _ := labels.NewMatcher(labels.MatchEqual, "quantile", "0.9")
	intervalQ90 := targetInterval.Match(matcher)
	assert.Len(t, intervalQ90, 1)
	assert.InDelta(t, 0.052614556, intervalQ90[0].Value, 0.000001)
}

// SPDX-License-Identifier: GPL-3.0-or-later

package prometheus

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/netdata/go.d.plugin/pkg/prometheus/selector"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testdata, _ = os.ReadFile("tests/testdata.txt")
var testdataNometa, _ = os.ReadFile("tests/testdata.nometa.txt")

func TestPrometheus404(t *testing.T) {
	tsMux := http.NewServeMux()
	tsMux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})
	ts := httptest.NewServer(tsMux)
	defer ts.Close()

	req := web.Request{URL: ts.URL + "/metrics"}
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

	req := web.Request{URL: ts.URL + "/metrics"}
	prom := New(http.DefaultClient, req)
	res, err := prom.Scrape()

	assert.NoError(t, err)
	verifyTestData(t, res)
}

func TestPrometheusPlainWithSelector(t *testing.T) {
	tsMux := http.NewServeMux()
	tsMux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(testdata)
	})
	ts := httptest.NewServer(tsMux)
	defer ts.Close()

	req := web.Request{URL: ts.URL + "/metrics"}
	sr, err := selector.Parse("go_gc*")
	require.NoError(t, err)
	prom := NewWithSelector(http.DefaultClient, req, sr)

	res, err := prom.Scrape()
	require.NoError(t, err)

	for _, v := range res {
		assert.Truef(t, strings.HasPrefix(v.Name(), "go_gc"), v.Name())
	}
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

	req := web.Request{URL: ts.URL + "/metrics"}
	prom := New(http.DefaultClient, req)

	for i := 0; i < 2; i++ {
		res, err := prom.Scrape()
		assert.NoError(t, err)
		verifyTestData(t, res)
	}
}

func TestParse(t *testing.T) {
	res := Metrics{}
	prom := prometheus{}
	err := prom.parse(testdata, &res, Metadata{})
	assert.NoError(t, err)

	res.Sort()
	verifyTestData(t, res)
}

func verifyTestData(t *testing.T, ms Metrics) {
	assert.Equal(t, 410, len(ms))
	assert.Equal(t, "go_gc_duration_seconds", ms[0].Labels.Get("__name__"))
	assert.Equal(t, "0.25", ms[0].Labels.Get("quantile"))
	assert.InDelta(t, 4.9351e-05, ms[0].Value, 0.0001)

	notExistYet := ms.FindByName("not_exist_yet")
	assert.NotNil(t, notExistYet)
	assert.Len(t, notExistYet, 0)

	targetInterval := ms.FindByName("prometheus_target_interval_length_seconds")
	assert.Len(t, targetInterval, 5)

	m, _ := labels.NewMatcher(labels.MatchEqual, "quantile", "0.9")
	intervalQ90 := targetInterval.Match(m)
	assert.Len(t, intervalQ90, 1)
	assert.InDelta(t, 0.052614556, intervalQ90[0].Value, 0.000001)
}

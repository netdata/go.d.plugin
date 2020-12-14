package couchbase

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	basicstat, _ = ioutil.ReadFile("testdata/basicstat.json")
)

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}

func Test_testDataIsCorrectlyReadAndValid(t *testing.T) {
	for name, data := range map[string][]byte{
		"basicstat": basicstat,
	} {
		require.NotNilf(t, data, name)
	}
}

func TestCouchbase_Init(t *testing.T) {
	tests := map[string]struct {
		config   Config
		wantFail bool
	}{
		"success on default config": {
			config: New().Config,
		},
		"fails on unset 'URL'": {
			wantFail: true,
			config: Config{
				HTTP: web.HTTP{
					Request: web.Request{
						URL: "",
					},
				},
			},
		},
		"fails on invalid URL": {
			wantFail: true,
			config: Config{
				HTTP: web.HTTP{
					Request: web.Request{
						URL: "127.0.0.1:9090",
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cb := New()
			cb.Config = test.config

			if test.wantFail {
				assert.False(t, cb.Init())
			} else {
				assert.True(t, cb.Init())
			}
		})
	}
}

func TestCouchbase_Collect(t *testing.T) {
	tests := map[string]struct {
		prepare       func() *Couchbase
		wantCollected map[string]int64
	}{
		"all stats": {
			prepare: func() *Couchbase {
				cb := New()
				return cb
			},
			wantCollected: map[string]int64{
				"basic_stats_beer-sample_data":                13660906,
				"basic_stats_beer-sample_disk":                90936691,
				"basic_stats_beer-sample_fetches":             0,
				"basic_stats_beer-sample_item_count":          107303,
				"basic_stats_beer-sample_mem":                 42653560,
				"basic_stats_beer-sample_num_non_resident":    0,
				"basic_stats_beer-sample_ops":                 0,
				"basic_stats_beer-sample_quota_used":          40,
				"basic_stats_gamesim-sample_data":             18177,
				"basic_stats_gamesim-sample_disk":             17091065,
				"basic_stats_gamesim-sample_fetches":          0,
				"basic_stats_gamesim-sample_item_count":       586,
				"basic_stats_gamesim-sample_mem":              19567320,
				"basic_stats_gamesim-sample_num_non_resident": 0,
				"basic_stats_gamesim-sample_ops":              0,
				"basic_stats_gamesim-sample_quota_used":       18,
				"basic_stats_netdata_data":                    0,
				"basic_stats_netdata_disk":                    40896724,
				"basic_stats_netdata_fetches":                 0,
				"basic_stats_netdata_item_count":              100000,
				"basic_stats_netdata_mem":                     36873944,
				"basic_stats_netdata_num_non_resident":        0,
				"basic_stats_netdata_ops":                     0,
				"basic_stats_netdata_quota_used":              0,
				"basic_stats_travel-sample_data":              0,
				"basic_stats_travel-sample_disk":              49636018,
				"basic_stats_travel-sample_fetches":           0,
				"basic_stats_travel-sample_item_count":        31591,
				"basic_stats_travel-sample_mem":               44275464,
				"basic_stats_travel-sample_num_non_resident":  0,
				"basic_stats_travel-sample_ops":               0,
				"basic_stats_travel-sample_quota_used":        42,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cb, cleanup := prepareCouchbase(t, test.prepare)
			defer cleanup()

			var collected map[string]int64
			for i := 0; i < 10; i++ {
				collected = cb.Collect()
			}

			assert.Equal(t, test.wantCollected, collected)
			ensureCollectedHasAllChartsDimsVarsIDs(t, cb, collected)
		})
	}
}

func TestCouchbase_Check(t *testing.T) {
	tests := map[string]struct {
		prepare  func(*testing.T) (cb *Couchbase, cleanup func())
		wantFail bool
	}{
		"valid data":         {prepare: prepareCouchbaseValidData},
		"invalid data":       {prepare: prepareCouchbaseInvalidData, wantFail: true},
		"404":                {prepare: prepareCouchbase404, wantFail: true},
		"connection refused": {prepare: prepareCouchbaseConnectionRefused, wantFail: true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cb, cleanup := test.prepare(t)
			defer cleanup()

			if test.wantFail {
				assert.False(t, cb.Check())
			} else {
				assert.True(t, cb.Check())
			}
		})
	}
}

func prepareCouchbase(t *testing.T, createCB func() *Couchbase) (cb *Couchbase, cleanup func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(basicstat)
		}))

	cb = createCB()
	cb.URL = srv.URL
	require.True(t, cb.Init())

	return cb, srv.Close
}

func prepareCouchbaseInvalidData(t *testing.T) (*Couchbase, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("hello and\n goodbye"))
		}))
	cb := New()
	cb.URL = srv.URL
	require.True(t, cb.Init())

	return cb, srv.Close
}

func prepareCouchbase404(t *testing.T) (*Couchbase, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
	cb := New()
	cb.URL = srv.URL
	require.True(t, cb.Init())

	return cb, srv.Close
}

func prepareCouchbaseConnectionRefused(t *testing.T) (*Couchbase, func()) {
	t.Helper()
	cb := New()
	cb.URL = "http://127.0.0.1:9090"
	require.True(t, cb.Init())

	return cb, func() {}
}

func prepareCouchbaseValidData(t *testing.T) (es *Couchbase, cleanup func()) {
	return prepareCouchbase(t, New)
}

func ensureCollectedHasAllChartsDimsVarsIDs(t *testing.T, cb *Couchbase, collected map[string]int64) {
	for _, chart := range *cb.Charts() {
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

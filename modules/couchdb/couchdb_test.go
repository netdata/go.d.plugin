package couchdb

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/tlscfg"
	"github.com/netdata/go.d.plugin/pkg/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	responseRoot, _        = ioutil.ReadFile("testdata/v311_single_root.json")
	responseActiveTasks, _ = ioutil.ReadFile("testdata/v311_single_active_tasks.json")
	responseNodeStats, _   = ioutil.ReadFile("testdata/v311_node_stats.json")
	responseNodeSystem, _  = ioutil.ReadFile("testdata/v311_node_system.json")
	responseDatabase, _    = ioutil.ReadFile("testdata/v311_database.json")
)

func Test_testDataIsCorrectlyReadAndValid(t *testing.T) {
	for name, data := range map[string][]byte{
		"responseRoot":        responseRoot,
		"responseActiveTasks": responseActiveTasks,
		"responseNodeStats":   responseNodeStats,
		"responseNodeSystem":  responseNodeSystem,
		"responseDatabase":    responseDatabase,
	} {
		require.NotNilf(t, data, name)
	}
}

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}

func TestCouchDB_Init(t *testing.T) {
	tests := map[string]struct {
		config          Config
		wantNumOfCharts int
		wantFail        bool
	}{
		"default": {
			wantNumOfCharts: numOfCharts(
				dbActivityCharts,
				httpTrafficBreakdownCharts,
				serverOperationsCharts,
				erlangStatisticsCharts,
			),
			config: New().Config,
		},
		"URL not set": {
			wantFail: true,
			config: Config{
				HTTP: web.HTTP{
					Request: web.Request{URL: ""},
				}},
		},
		"invalid TLSCA": {
			wantFail: true,
			config: Config{
				HTTP: web.HTTP{
					Client: web.Client{
						TLSConfig: tlscfg.TLSConfig{TLSCA: "testdata/tls"},
					},
				}},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			es := New()
			es.Config = test.config

			if test.wantFail {
				assert.False(t, es.Init())
			} else {
				assert.True(t, es.Init())
				assert.Equal(t, test.wantNumOfCharts, len(*es.Charts()))
			}
		})
	}
}

func TestCouchDB_Check(t *testing.T) {
	tests := map[string]struct {
		prepare  func(*testing.T) (cdb *CouchDB, cleanup func())
		wantFail bool
	}{
		"valid data":         {prepare: prepareCouchDBValidData},
		"invalid data":       {prepare: prepareCouchDBInvalidData, wantFail: true},
		"404":                {prepare: prepareCouchDB404, wantFail: true},
		"connection refused": {prepare: prepareCouchDBConnectionRefused, wantFail: true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cdb, cleanup := test.prepare(t)
			defer cleanup()

			if test.wantFail {
				assert.False(t, cdb.Check())
			} else {
				assert.True(t, cdb.Check())
			}
		})
	}
}

func TestCouchDB_Charts(t *testing.T) {
	assert.Nil(t, New().Charts())
}

func TestCouchDB_Cleanup(t *testing.T) {
	assert.NotPanics(t, New().Cleanup)
}

func TestCouchDB_Collect(t *testing.T) {
	tests := map[string]struct {
		prepare       func() *CouchDB
		wantCollected map[string]int64
	}{
		"all stats": {
			prepare: func() *CouchDB {
				cdb := New()
				return cdb
			},
			wantCollected: map[string]int64{

				// node stats
				"couch_replicator_jobs_crashed":         1,
				"couch_replicator_jobs_penging":         1,
				"couch_replicator_jobs_running":         1,
				"couchdb_database_reads":                1,
				"couchdb_database_writes":               14,
				"couchdb_httpd_request_methods_COPY":    1,
				"couchdb_httpd_request_methods_DELETE":  1,
				"couchdb_httpd_request_methods_GET":     75544,
				"couchdb_httpd_request_methods_HEAD":    1,
				"couchdb_httpd_request_methods_OPTIONS": 1,
				"couchdb_httpd_request_methods_POST":    15,
				"couchdb_httpd_request_methods_PUT":     3,
				"couchdb_httpd_status_codes_200":        75294,
				"couchdb_httpd_status_codes_201":        15,
				"couchdb_httpd_status_codes_202":        1,
				"couchdb_httpd_status_codes_2xx":        2,
				"couchdb_httpd_status_codes_3xx":        3,
				"couchdb_httpd_status_codes_4xx":        258,
				"couchdb_httpd_status_codes_5xx":        2,
				"couchdb_httpd_view_reads":              1,
				"couchdb_open_os_files":                 1,

				// node system
				"context_switches":          22614499,
				"ets_table_count":           116,
				"internal_replication_jobs": 1,
				"io_input":                  49674812,
				"io_output":                 686400800,
				"memory_atom_used":          488328,
				"memory_atom":               504433,
				"memory_binary":             297696,
				"memory_code":               11252688,
				"memory_ets":                1579120,
				"memory_other":              20427855,
				"memory_processes":          9161448,
				"os_proc_count":             1,
				"process_count":             296,
				"reductions":                43211228312,
				"run_queue":                 1,

				// active tasks
				"active_tasks_database_compaction": 1,
				"active_tasks_indexer":             2,
				"active_tasks_replication":         1,
				"active_tasks_view_compaction":     1,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cdb, cleanup := prepareCouchDB(t, test.prepare)
			defer cleanup()

			var collected map[string]int64
			for i := 0; i < 10; i++ {
				collected = cdb.Collect()
			}

			assert.Equal(t, test.wantCollected, collected)
			ensureCollectedHasAllChartsDimsVarsIDs(t, cdb, collected)
		})
	}
}

func ensureCollectedHasAllChartsDimsVarsIDs(t *testing.T, cdb *CouchDB, collected map[string]int64) {
	for _, chart := range *cdb.Charts() {
		if chart.Obsolete {
			continue
		}
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

func prepareCouchDB(t *testing.T, createCDB func() *CouchDB) (cdb *CouchDB, cleanup func()) {
	t.Helper()
	srv := prepareCouchDBEndpoint()

	cdb = createCDB()
	cdb.URL = srv.URL
	require.True(t, cdb.Init())

	return cdb, srv.Close
}

func prepareCouchDBValidData(t *testing.T) (cdb *CouchDB, cleanup func()) {
	return prepareCouchDB(t, New)
}

func prepareCouchDBInvalidData(t *testing.T) (*CouchDB, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("hello and\n goodbye"))
		}))
	cdb := New()
	cdb.URL = srv.URL
	require.True(t, cdb.Init())

	return cdb, srv.Close
}

func prepareCouchDB404(t *testing.T) (*CouchDB, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
	cdb := New()
	cdb.URL = srv.URL
	require.True(t, cdb.Init())

	return cdb, srv.Close
}

func prepareCouchDBConnectionRefused(t *testing.T) (*CouchDB, func()) {
	t.Helper()
	cdb := New()
	cdb.URL = "http://127.0.0.1:38001"
	require.True(t, cdb.Init())

	return cdb, func() {}
}

func prepareCouchDBEndpoint() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case urlPathOverviewStats:
				_, _ = w.Write(responseNodeStats)
			case urlPathSystemStats:
				_, _ = w.Write(responseNodeSystem)
			case urlPathActiveTasks:
				_, _ = w.Write(responseActiveTasks)
			case "/":
				_, _ = w.Write(responseRoot)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}))
}

func numOfCharts(charts ...Charts) (num int) {
	for _, v := range charts {
		num += len(v)
	}
	return num
}

package elasticsearch

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	v790SingleNodesLocalStats, _ = ioutil.ReadFile("testdata/v790_single_nodes_local_stats.json")
	v790SingleClusterHealth, _   = ioutil.ReadFile("testdata/v790_single_cluster_health.json")
	v790SingleClusterStats, _    = ioutil.ReadFile("testdata/v790_single_cluster_stats.json")
)

func Test_testDataIsCorrectlyReadAndValid(t *testing.T) {
	for name, data := range map[string][]byte{
		"v790SingleNodesLocalStats": v790SingleNodesLocalStats,
		"v790SingleClusterHealth":   v790SingleClusterHealth,
		"v790SingleClusterStats":    v790SingleClusterStats,
	} {
		require.NotNilf(t, data, name)
	}
}

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}

func TestElasticsearch_Init(t *testing.T) {
	tests := map[string]struct {
		config   Config
		wantFail bool
	}{
		"default": {
			config: New().Config,
		},
		"URL not set": {
			wantFail: true,
			config: Config{web.HTTP{
				Request: web.Request{UserURL: ""},
			}},
		},
		"invalid TLSCA": {
			wantFail: true,
			config: Config{web.HTTP{
				Client: web.Client{ClientTLSConfig: web.ClientTLSConfig{TLSCA: "testdata/tls"}},
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
			}
		})
	}
}

func TestElasticsearch_Check(t *testing.T) {

}

func TestElasticsearch_Charts(t *testing.T) {

}

func TestElasticsearch_Cleanup(t *testing.T) {
	assert.NotPanics(t, New().Cleanup)
}

func TestElasticsearch_Collect(t *testing.T) {
	tests := map[string]struct {
		prepare       func(t *testing.T) (es *Elasticsearch, cleanup func())
		wantCollected map[string]int64
	}{
		"v790": {
			prepare: prepareElasticsearch,
			wantCollected: map[string]int64{
				"cluster_health_active_shards":                   0,
				"cluster_health_active_shards_percent_as_number": 100,
				"cluster_health_delayed_unassigned_shards":       0,
				"cluster_health_initializing_shards":             0,
				"cluster_health_number_of_data_nodes":            1,
				"cluster_health_number_of_in_flight_fetch":       0,
				"cluster_health_number_of_nodes":                 1,
				"cluster_health_number_of_pending_tasks":         0,
				"cluster_health_relocating_shards":               0,
				"cluster_health_unassigned_shards":               0,
				"cluster_stats_indices_count":                    0,
				"cluster_stats_indices_docs_count":               0,
				"cluster_stats_indices_query_cache_hit_count":    0,
				"cluster_stats_indices_query_cache_miss_count":   0,
				"cluster_stats_indices_shards_total":             0,
				"cluster_stats_indices_store_size_in_bytes":      0,
				"cluster_stats_nodes_count_coordinating_only":    0,
				"cluster_stats_nodes_count_data":                 1,
				"cluster_stats_nodes_count_ingest":               1,
				"cluster_stats_nodes_count_master":               1,
				"cluster_stats_nodes_count_total":                1,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			es, cleanup := test.prepare(t)
			defer cleanup()

			collected := es.Collect()

			assert.Equal(t, test.wantCollected, collected)
		})
	}
}

func prepareElasticsearch(t *testing.T) (es *Elasticsearch, cleanup func()) {
	t.Helper()
	srv := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case urlPathNodesLocalStats:
					_, _ = w.Write(v790SingleNodesLocalStats)
				case urlPathClusterHealth:
					_, _ = w.Write(v790SingleClusterHealth)
				case urlPathClusterStats:
					_, _ = w.Write(v790SingleClusterStats)
				default:
					w.WriteHeader(http.StatusNotFound)
				}
			}))

	es = New()
	es.UserURL = srv.URL
	require.True(t, es.Init())

	return es, srv.Close
}

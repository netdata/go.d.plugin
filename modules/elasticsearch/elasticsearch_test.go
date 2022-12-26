// SPDX-License-Identifier: GPL-3.0-or-later

package elasticsearch

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/netdata/go.d.plugin/pkg/tlscfg"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	v790NodesLocalStats, _ = os.ReadFile("testdata/v7.9.0/nodes_local_stats.json")
	v790ClusterHealth, _   = os.ReadFile("testdata/v7.9.0/cluster_health.json")
	v790ClusterStats, _    = os.ReadFile("testdata/v7.9.0/cluster_stats.json")
	v790CatIndicesStats, _ = os.ReadFile("testdata/v7.9.0/cat_indices_stats.json")
	v790Info, _            = os.ReadFile("testdata/v7.9.0/info.json")
)

func TestNew2(t *testing.T) {
	fmt.Println(indexDimID("open", ""))
}

func Test_testDataIsCorrectlyReadAndValid(t *testing.T) {
	for name, data := range map[string][]byte{
		"v790NodesLocalStats": v790NodesLocalStats,
		"v790ClusterHealth":   v790ClusterHealth,
		"v790ClusterStats":    v790ClusterStats,
		"v790CatIndicesStats": v790CatIndicesStats,
		"v790Info":            v790Info,
	} {
		require.NotNilf(t, data, name)
	}
}

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}

func TestElasticsearch_Init(t *testing.T) {
	tests := map[string]struct {
		config          Config
		wantNumOfCharts int
		wantFail        bool
	}{
		"default": {
			wantNumOfCharts: numOfCharts(
				nodeCharts,
				clusterHealthCharts,
				clusterStatsCharts,
			),
			config: New().Config,
		},
		"all stats": {
			wantNumOfCharts: numOfCharts(
				nodeCharts,
				clusterHealthCharts,
				clusterStatsCharts,
			),
			config: Config{
				HTTP: web.HTTP{
					Request: web.Request{URL: "http://127.0.0.1:38001"},
				},
				DoNodeStats:     true,
				DoClusterHealth: true,
				DoClusterStats:  true,
				DoIndicesStats:  true,
			},
		},
		"only node_stats": {
			wantNumOfCharts: len(nodeCharts),
			config: Config{
				HTTP: web.HTTP{
					Request: web.Request{URL: "http://127.0.0.1:38001"},
				},
				DoNodeStats:     true,
				DoClusterHealth: false,
				DoClusterStats:  false,
				DoIndicesStats:  false,
			},
		},
		"only cluster_health": {
			wantNumOfCharts: len(clusterHealthCharts),
			config: Config{
				HTTP: web.HTTP{
					Request: web.Request{URL: "http://127.0.0.1:38001"},
				},
				DoNodeStats:     false,
				DoClusterHealth: true,
				DoClusterStats:  false,
				DoIndicesStats:  false,
			},
		},
		"only cluster_stats": {
			wantNumOfCharts: len(clusterStatsCharts),
			config: Config{
				HTTP: web.HTTP{
					Request: web.Request{URL: "http://127.0.0.1:38001"},
				},
				DoNodeStats:     false,
				DoClusterHealth: false,
				DoClusterStats:  true,
				DoIndicesStats:  false,
			},
		},
		"only indices_stats": {
			wantNumOfCharts: 0,
			config: Config{
				HTTP: web.HTTP{
					Request: web.Request{URL: "http://127.0.0.1:38001"},
				},
				DoNodeStats:     false,
				DoClusterHealth: false,
				DoClusterStats:  false,
				DoIndicesStats:  true,
			},
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
		"all API calls are disabled": {
			wantFail: true,
			config: Config{
				HTTP: web.HTTP{
					Request: web.Request{URL: "http://127.0.0.1:38001"},
				},
				DoNodeStats:     false,
				DoClusterHealth: false,
				DoClusterStats:  false,
				DoIndicesStats:  false,
			},
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

func TestElasticsearch_Check(t *testing.T) {
	tests := map[string]struct {
		prepare  func(*testing.T) (es *Elasticsearch, cleanup func())
		wantFail bool
	}{
		"valid data":         {prepare: prepareElasticsearchValidData},
		"invalid data":       {prepare: prepareElasticsearchInvalidData, wantFail: true},
		"404":                {prepare: prepareElasticsearch404, wantFail: true},
		"connection refused": {prepare: prepareElasticsearchConnectionRefused, wantFail: true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			es, cleanup := test.prepare(t)
			defer cleanup()

			if test.wantFail {
				assert.False(t, es.Check())
			} else {
				assert.True(t, es.Check())
			}
		})
	}
}

func TestElasticsearch_Charts(t *testing.T) {
	assert.Nil(t, New().Charts())
}

func TestElasticsearch_Cleanup(t *testing.T) {
	assert.NotPanics(t, New().Cleanup)
}

func TestElasticsearch_Collect(t *testing.T) {
	tests := map[string]struct {
		prepare       func() *Elasticsearch
		wantCollected map[string]int64
	}{
		"v790: all stats": {
			prepare: func() *Elasticsearch {
				es := New()
				es.DoNodeStats = true
				es.DoClusterHealth = true
				es.DoClusterStats = true
				es.DoIndicesStats = true
				return es
			},
			wantCollected: map[string]int64{
				"cluster_active_primary_shards":                          1,
				"cluster_active_shards":                                  1,
				"cluster_active_shards_percent_as_number":                100,
				"cluster_delayed_unassigned_shards":                      1,
				"cluster_indices_count":                                  3,
				"cluster_indices_docs_count":                             1,
				"cluster_indices_query_cache_hit_count":                  1,
				"cluster_indices_query_cache_miss_count":                 1,
				"cluster_indices_shards_primaries":                       3,
				"cluster_indices_shards_replication":                     1,
				"cluster_indices_shards_total":                           3,
				"cluster_indices_store_size_in_bytes":                    624,
				"cluster_initializing_shards":                            1,
				"cluster_nodes_count_coordinating_only":                  1,
				"cluster_nodes_count_data":                               1,
				"cluster_nodes_count_ingest":                             1,
				"cluster_nodes_count_master":                             1,
				"cluster_nodes_count_ml":                                 1,
				"cluster_nodes_count_remote_cluster_client":              1,
				"cluster_nodes_count_total":                              1,
				"cluster_nodes_count_transform":                          1,
				"cluster_nodes_count_voting_only":                        1,
				"cluster_number_of_data_nodes":                           1,
				"cluster_number_of_in_flight_fetch":                      1,
				"cluster_number_of_nodes":                                1,
				"cluster_number_of_pending_tasks":                        1,
				"cluster_relocating_shards":                              1,
				"cluster_status_green":                                   1,
				"cluster_status_red":                                     0,
				"cluster_status_yellow":                                  0,
				"cluster_unassigned_shards":                              1,
				"node_breakers_accounting_tripped":                       1,
				"node_breakers_fielddata_tripped":                        1,
				"node_breakers_in_flight_requests_tripped":               1,
				"node_breakers_model_inference_tripped":                  1,
				"node_breakers_parent_tripped":                           1,
				"node_breakers_request_tripped":                          1,
				"node_http_current_open":                                 3,
				"node_index_my-index-000001_stats_docs_count":            1,
				"node_index_my-index-000001_stats_health_green":          0,
				"node_index_my-index-000001_stats_health_red":            0,
				"node_index_my-index-000001_stats_health_yellow":         1,
				"node_index_my-index-000001_stats_shards_count":          1,
				"node_index_my-index-000001_stats_store_size_in_bytes":   208,
				"node_index_my-index-000002_stats_docs_count":            1,
				"node_index_my-index-000002_stats_health_green":          0,
				"node_index_my-index-000002_stats_health_red":            0,
				"node_index_my-index-000002_stats_health_yellow":         1,
				"node_index_my-index-000002_stats_shards_count":          1,
				"node_index_my-index-000002_stats_store_size_in_bytes":   208,
				"node_index_my-index-000003_stats_docs_count":            1,
				"node_index_my-index-000003_stats_health_green":          0,
				"node_index_my-index-000003_stats_health_red":            0,
				"node_index_my-index-000003_stats_health_yellow":         1,
				"node_index_my-index-000003_stats_shards_count":          1,
				"node_index_my-index-000003_stats_store_size_in_bytes":   208,
				"node_indices_fielddata_evictions":                       1,
				"node_indices_fielddata_memory_size_in_bytes":            1,
				"node_indices_flush_total":                               1,
				"node_indices_flush_total_time_in_millis":                1,
				"node_indices_indexing_index_current":                    1,
				"node_indices_indexing_index_time_in_millis":             1,
				"node_indices_indexing_index_total":                      1,
				"node_indices_refresh_total":                             1,
				"node_indices_refresh_total_time_in_millis":              1,
				"node_indices_search_fetch_current":                      1,
				"node_indices_search_fetch_time_in_millis":               1,
				"node_indices_search_fetch_total":                        1,
				"node_indices_search_query_current":                      1,
				"node_indices_search_query_time_in_millis":               1,
				"node_indices_search_query_total":                        1,
				"node_indices_segments_count":                            1,
				"node_indices_segments_doc_values_memory_in_bytes":       1,
				"node_indices_segments_fixed_bit_set_memory_in_bytes":    1,
				"node_indices_segments_index_writer_memory_in_bytes":     1,
				"node_indices_segments_memory_in_bytes":                  1,
				"node_indices_segments_norms_memory_in_bytes":            1,
				"node_indices_segments_points_memory_in_bytes":           1,
				"node_indices_segments_stored_fields_memory_in_bytes":    1,
				"node_indices_segments_term_vectors_memory_in_bytes":     1,
				"node_indices_segments_terms_memory_in_bytes":            1,
				"node_indices_segments_version_map_memory_in_bytes":      1,
				"node_indices_translog_operations":                       1,
				"node_indices_translog_size_in_bytes":                    1,
				"node_indices_translog_uncommitted_operations":           1,
				"node_indices_translog_uncommitted_size_in_bytes":        1,
				"node_jvm_buffer_pools_direct_count":                     15,
				"node_jvm_buffer_pools_direct_total_capacity_in_bytes":   6321124,
				"node_jvm_buffer_pools_direct_used_in_bytes":             6321125,
				"node_jvm_buffer_pools_mapped_count":                     1,
				"node_jvm_buffer_pools_mapped_total_capacity_in_bytes":   1,
				"node_jvm_buffer_pools_mapped_used_in_bytes":             1,
				"node_jvm_gc_collectors_old_collection_count":            1,
				"node_jvm_gc_collectors_old_collection_time_in_millis":   1,
				"node_jvm_gc_collectors_young_collection_count":          16,
				"node_jvm_gc_collectors_young_collection_time_in_millis": 184,
				"node_jvm_mem_heap_committed_in_bytes":                   1073741824,
				"node_jvm_mem_heap_used_in_bytes":                        363166720,
				"node_jvm_mem_heap_used_percent":                         33,
				"node_process_max_file_descriptors":                      1048576,
				"node_process_open_file_descriptors":                     258,
				"node_thread_pool_analyze_queue":                         1,
				"node_thread_pool_analyze_rejected":                      1,
				"node_thread_pool_fetch_shard_started_queue":             1,
				"node_thread_pool_fetch_shard_started_rejected":          1,
				"node_thread_pool_fetch_shard_store_queue":               1,
				"node_thread_pool_fetch_shard_store_rejected":            1,
				"node_thread_pool_flush_queue":                           1,
				"node_thread_pool_flush_rejected":                        1,
				"node_thread_pool_force_merge_queue":                     1,
				"node_thread_pool_force_merge_rejected":                  1,
				"node_thread_pool_generic_queue":                         1,
				"node_thread_pool_generic_rejected":                      1,
				"node_thread_pool_get_queue":                             1,
				"node_thread_pool_get_rejected":                          1,
				"node_thread_pool_listener_queue":                        1,
				"node_thread_pool_listener_rejected":                     1,
				"node_thread_pool_management_queue":                      1,
				"node_thread_pool_management_rejected":                   1,
				"node_thread_pool_refresh_queue":                         1,
				"node_thread_pool_refresh_rejected":                      1,
				"node_thread_pool_search_queue":                          1,
				"node_thread_pool_search_rejected":                       1,
				"node_thread_pool_search_throttled_queue":                1,
				"node_thread_pool_search_throttled_rejected":             1,
				"node_thread_pool_snapshot_queue":                        1,
				"node_thread_pool_snapshot_rejected":                     1,
				"node_thread_pool_warmer_queue":                          1,
				"node_thread_pool_warmer_rejected":                       1,
				"node_thread_pool_write_queue":                           1,
				"node_thread_pool_write_rejected":                        1,
				"node_transport_rx_count":                                1,
				"node_transport_rx_size_in_bytes":                        1,
				"node_transport_tx_count":                                1,
				"node_transport_tx_size_in_bytes":                        1,
			},
		},
		"v790: only node_stats": {
			prepare: func() *Elasticsearch {
				es := New()
				es.DoNodeStats = true
				es.DoClusterHealth = false
				es.DoClusterStats = false
				es.DoIndicesStats = false
				return es
			},
			wantCollected: map[string]int64{
				"node_breakers_accounting_tripped":                       1,
				"node_breakers_fielddata_tripped":                        1,
				"node_breakers_in_flight_requests_tripped":               1,
				"node_breakers_model_inference_tripped":                  1,
				"node_breakers_parent_tripped":                           1,
				"node_breakers_request_tripped":                          1,
				"node_http_current_open":                                 3,
				"node_indices_fielddata_evictions":                       1,
				"node_indices_fielddata_memory_size_in_bytes":            1,
				"node_indices_flush_total":                               1,
				"node_indices_flush_total_time_in_millis":                1,
				"node_indices_indexing_index_current":                    1,
				"node_indices_indexing_index_time_in_millis":             1,
				"node_indices_indexing_index_total":                      1,
				"node_indices_refresh_total":                             1,
				"node_indices_refresh_total_time_in_millis":              1,
				"node_indices_search_fetch_current":                      1,
				"node_indices_search_fetch_time_in_millis":               1,
				"node_indices_search_fetch_total":                        1,
				"node_indices_search_query_current":                      1,
				"node_indices_search_query_time_in_millis":               1,
				"node_indices_search_query_total":                        1,
				"node_indices_segments_count":                            1,
				"node_indices_segments_doc_values_memory_in_bytes":       1,
				"node_indices_segments_fixed_bit_set_memory_in_bytes":    1,
				"node_indices_segments_index_writer_memory_in_bytes":     1,
				"node_indices_segments_memory_in_bytes":                  1,
				"node_indices_segments_norms_memory_in_bytes":            1,
				"node_indices_segments_points_memory_in_bytes":           1,
				"node_indices_segments_stored_fields_memory_in_bytes":    1,
				"node_indices_segments_term_vectors_memory_in_bytes":     1,
				"node_indices_segments_terms_memory_in_bytes":            1,
				"node_indices_segments_version_map_memory_in_bytes":      1,
				"node_indices_translog_operations":                       1,
				"node_indices_translog_size_in_bytes":                    1,
				"node_indices_translog_uncommitted_operations":           1,
				"node_indices_translog_uncommitted_size_in_bytes":        1,
				"node_jvm_buffer_pools_direct_count":                     15,
				"node_jvm_buffer_pools_direct_total_capacity_in_bytes":   6321124,
				"node_jvm_buffer_pools_direct_used_in_bytes":             6321125,
				"node_jvm_buffer_pools_mapped_count":                     1,
				"node_jvm_buffer_pools_mapped_total_capacity_in_bytes":   1,
				"node_jvm_buffer_pools_mapped_used_in_bytes":             1,
				"node_jvm_gc_collectors_old_collection_count":            1,
				"node_jvm_gc_collectors_old_collection_time_in_millis":   1,
				"node_jvm_gc_collectors_young_collection_count":          16,
				"node_jvm_gc_collectors_young_collection_time_in_millis": 184,
				"node_jvm_mem_heap_committed_in_bytes":                   1073741824,
				"node_jvm_mem_heap_used_in_bytes":                        363166720,
				"node_jvm_mem_heap_used_percent":                         33,
				"node_process_max_file_descriptors":                      1048576,
				"node_process_open_file_descriptors":                     258,
				"node_thread_pool_analyze_queue":                         1,
				"node_thread_pool_analyze_rejected":                      1,
				"node_thread_pool_fetch_shard_started_queue":             1,
				"node_thread_pool_fetch_shard_started_rejected":          1,
				"node_thread_pool_fetch_shard_store_queue":               1,
				"node_thread_pool_fetch_shard_store_rejected":            1,
				"node_thread_pool_flush_queue":                           1,
				"node_thread_pool_flush_rejected":                        1,
				"node_thread_pool_force_merge_queue":                     1,
				"node_thread_pool_force_merge_rejected":                  1,
				"node_thread_pool_generic_queue":                         1,
				"node_thread_pool_generic_rejected":                      1,
				"node_thread_pool_get_queue":                             1,
				"node_thread_pool_get_rejected":                          1,
				"node_thread_pool_listener_queue":                        1,
				"node_thread_pool_listener_rejected":                     1,
				"node_thread_pool_management_queue":                      1,
				"node_thread_pool_management_rejected":                   1,
				"node_thread_pool_refresh_queue":                         1,
				"node_thread_pool_refresh_rejected":                      1,
				"node_thread_pool_search_queue":                          1,
				"node_thread_pool_search_rejected":                       1,
				"node_thread_pool_search_throttled_queue":                1,
				"node_thread_pool_search_throttled_rejected":             1,
				"node_thread_pool_snapshot_queue":                        1,
				"node_thread_pool_snapshot_rejected":                     1,
				"node_thread_pool_warmer_queue":                          1,
				"node_thread_pool_warmer_rejected":                       1,
				"node_thread_pool_write_queue":                           1,
				"node_thread_pool_write_rejected":                        1,
				"node_transport_rx_count":                                1,
				"node_transport_rx_size_in_bytes":                        1,
				"node_transport_tx_count":                                1,
				"node_transport_tx_size_in_bytes":                        1,
			},
		},
		"v790: only cluster_health": {
			prepare: func() *Elasticsearch {
				es := New()
				es.DoNodeStats = false
				es.DoClusterHealth = true
				es.DoClusterStats = false
				es.DoIndicesStats = false
				return es
			},
			wantCollected: map[string]int64{
				"cluster_active_primary_shards":           1,
				"cluster_active_shards":                   1,
				"cluster_active_shards_percent_as_number": 100,
				"cluster_delayed_unassigned_shards":       1,
				"cluster_initializing_shards":             1,
				"cluster_number_of_data_nodes":            1,
				"cluster_number_of_in_flight_fetch":       1,
				"cluster_number_of_nodes":                 1,
				"cluster_number_of_pending_tasks":         1,
				"cluster_relocating_shards":               1,
				"cluster_status_green":                    1,
				"cluster_status_red":                      0,
				"cluster_status_yellow":                   0,
				"cluster_unassigned_shards":               1,
			},
		},
		"v790: only cluster_stats": {
			prepare: func() *Elasticsearch {
				es := New()
				es.DoNodeStats = false
				es.DoClusterHealth = false
				es.DoClusterStats = true
				es.DoIndicesStats = false
				return es
			},
			wantCollected: map[string]int64{
				"cluster_indices_count":                     3,
				"cluster_indices_docs_count":                1,
				"cluster_indices_query_cache_hit_count":     1,
				"cluster_indices_query_cache_miss_count":    1,
				"cluster_indices_shards_primaries":          3,
				"cluster_indices_shards_replication":        1,
				"cluster_indices_shards_total":              3,
				"cluster_indices_store_size_in_bytes":       624,
				"cluster_nodes_count_coordinating_only":     1,
				"cluster_nodes_count_data":                  1,
				"cluster_nodes_count_ingest":                1,
				"cluster_nodes_count_master":                1,
				"cluster_nodes_count_ml":                    1,
				"cluster_nodes_count_remote_cluster_client": 1,
				"cluster_nodes_count_total":                 1,
				"cluster_nodes_count_transform":             1,
				"cluster_nodes_count_voting_only":           1,
			},
		},
		"v790: only indices_stats": {
			prepare: func() *Elasticsearch {
				es := New()
				es.DoNodeStats = false
				es.DoClusterHealth = false
				es.DoClusterStats = false
				es.DoIndicesStats = true
				return es
			},
			wantCollected: map[string]int64{
				"node_index_my-index-000001_stats_docs_count":          1,
				"node_index_my-index-000001_stats_health_green":        0,
				"node_index_my-index-000001_stats_health_red":          0,
				"node_index_my-index-000001_stats_health_yellow":       1,
				"node_index_my-index-000001_stats_shards_count":        1,
				"node_index_my-index-000001_stats_store_size_in_bytes": 208,
				"node_index_my-index-000002_stats_docs_count":          1,
				"node_index_my-index-000002_stats_health_green":        0,
				"node_index_my-index-000002_stats_health_red":          0,
				"node_index_my-index-000002_stats_health_yellow":       1,
				"node_index_my-index-000002_stats_shards_count":        1,
				"node_index_my-index-000002_stats_store_size_in_bytes": 208,
				"node_index_my-index-000003_stats_docs_count":          1,
				"node_index_my-index-000003_stats_health_green":        0,
				"node_index_my-index-000003_stats_health_red":          0,
				"node_index_my-index-000003_stats_health_yellow":       1,
				"node_index_my-index-000003_stats_shards_count":        1,
				"node_index_my-index-000003_stats_store_size_in_bytes": 208,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			es, cleanup := prepareElasticsearch(t, test.prepare)
			defer cleanup()

			var mx map[string]int64
			for i := 0; i < 10; i++ {
				mx = es.Collect()
			}

			assert.Equal(t, test.wantCollected, mx)
			ensureCollectedHasAllChartsDimsVarsIDs(t, es, mx)
		})
	}
}

func ensureCollectedHasAllChartsDimsVarsIDs(t *testing.T, es *Elasticsearch, collected map[string]int64) {
	for _, chart := range *es.Charts() {
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

func prepareElasticsearch(t *testing.T, createES func() *Elasticsearch) (es *Elasticsearch, cleanup func()) {
	t.Helper()
	srv := prepareElasticsearchEndpoint()

	es = createES()
	es.URL = srv.URL
	require.True(t, es.Init())

	return es, srv.Close
}

func prepareElasticsearchValidData(t *testing.T) (es *Elasticsearch, cleanup func()) {
	return prepareElasticsearch(t, New)
}

func prepareElasticsearchInvalidData(t *testing.T) (*Elasticsearch, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("hello and\n goodbye"))
		}))
	es := New()
	es.URL = srv.URL
	require.True(t, es.Init())

	return es, srv.Close
}

func prepareElasticsearch404(t *testing.T) (*Elasticsearch, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
	es := New()
	es.URL = srv.URL
	require.True(t, es.Init())

	return es, srv.Close
}

func prepareElasticsearchConnectionRefused(t *testing.T) (*Elasticsearch, func()) {
	t.Helper()
	es := New()
	es.URL = "http://127.0.0.1:38001"
	require.True(t, es.Init())

	return es, func() {}
}

func prepareElasticsearchEndpoint() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case urlPathLocalNodeStats:
				_, _ = w.Write(v790NodesLocalStats)
			case urlPathClusterHealth:
				_, _ = w.Write(v790ClusterHealth)
			case urlPathClusterStats:
				_, _ = w.Write(v790ClusterStats)
			case urlPathIndicesStats:
				_, _ = w.Write(v790CatIndicesStats)
			case "/":
				_, _ = w.Write(v790Info)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}))
}

func numOfCharts(charts ...module.Charts) (num int) {
	for _, v := range charts {
		num += len(v)
	}
	return num
}

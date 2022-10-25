// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	dataMetrics, _ = os.ReadFile("testdata/metrics.txt")
)

func Test_TestData(t *testing.T) {
	for name, data := range map[string][]byte{
		"dataMetrics": dataMetrics,
	} {
		assert.NotNilf(t, data, name)
	}
}

func TestNew(t *testing.T) {
	assert.IsType(t, (*Cassandra)(nil), New())
}

func TestCassandra_Init(t *testing.T) {
	tests := map[string]struct {
		config   Config
		wantFail bool
	}{
		"success if 'url' is set": {
			config: Config{
				HTTP: web.HTTP{Request: web.Request{URL: "http://127.0.0.1:7072"}}},
		},
		"success on default config": {
			wantFail: false,
			config:   New().Config,
		},
		"fails if 'url' is unset": {
			wantFail: true,
			config:   Config{HTTP: web.HTTP{Request: web.Request{URL: ""}}},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c := New()
			c.Config = test.config

			if test.wantFail {
				assert.False(t, c.Init())
			} else {
				assert.True(t, c.Init())
			}
		})
	}
}

func TestCassandra_Check(t *testing.T) {
	tests := map[string]struct {
		prepare  func() (c *Cassandra, cleanup func())
		wantFail bool
	}{
		"success on valid response": {
			prepare: prepareCassandra,
		},
		"fails if endpoint returns invalid data": {
			wantFail: true,
			prepare:  prepareCassandraInvalidData,
		},
		"fails on connection refused": {
			wantFail: true,
			prepare:  prepareCassandraConnectionRefused,
		},
		"fails on 404 response": {
			wantFail: true,
			prepare:  prepareCassandraResponse404,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c, cleanup := test.prepare()
			defer cleanup()

			require.True(t, c.Init())

			if test.wantFail {
				assert.False(t, c.Check())
			} else {
				assert.True(t, c.Check())
			}
		})
	}
}

func TestCassandra_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestCassandra_Collect(t *testing.T) {
	tests := map[string]struct {
		prepare       func() (c *Cassandra, cleanup func())
		wantCollected map[string]int64
	}{
		"success on valid response": {
			prepare: prepareCassandra,
			wantCollected: map[string]int64{
				"client_request_failures_reads":                                0,
				"client_request_failures_writes":                               0,
				"client_request_latency_reads":                                 2279017,
				"client_request_latency_writes":                                2285110,
				"client_request_timeouts_reads":                                0,
				"client_request_timeouts_writes":                               0,
				"client_request_total_latency_reads":                           155385053,
				"client_request_total_latency_writes":                          100869491,
				"client_request_unavailables_reads":                            0,
				"client_request_unavailables_writes":                           0,
				"compaction_bytes_compacted":                                   840350211,
				"compaction_completed_tasks":                                   185,
				"compaction_pending_tasks":                                     0,
				"dropped_messages_one_minute":                                  0,
				"jvm_gc_cms_count":                                             4,
				"jvm_gc_cms_time":                                              208,
				"jvm_gc_parnew_count":                                          324,
				"jvm_gc_parnew_time":                                           5020,
				"jvm_memory_heap_used":                                         1262549184,
				"jvm_memory_nonheap_used":                                      95144624,
				"key_cache_hit_ratio":                                          76552,
				"key_cache_hits":                                               5323032,
				"key_cache_misses":                                             1630411,
				"key_cache_size":                                               131861248,
				"key_cache_utilization":                                        13972,
				"row_cache_hit_ratio":                                          0,
				"row_cache_hits":                                               0,
				"row_cache_misses":                                             0,
				"row_cache_size":                                               0,
				"row_cache_utilization":                                        0,
				"storage_exceptions":                                           0,
				"storage_load":                                                 630321359,
				"thread_pool_CacheCleanupExecutor_active_tasks":                0,
				"thread_pool_CacheCleanupExecutor_blocked_tasks":               0,
				"thread_pool_CacheCleanupExecutor_pending_tasks":               0,
				"thread_pool_CacheCleanupExecutor_total_blocked_tasks":         0,
				"thread_pool_CompactionExecutor_active_tasks":                  0,
				"thread_pool_CompactionExecutor_blocked_tasks":                 0,
				"thread_pool_CompactionExecutor_pending_tasks":                 0,
				"thread_pool_CompactionExecutor_total_blocked_tasks":           0,
				"thread_pool_GossipStage_active_tasks":                         0,
				"thread_pool_GossipStage_blocked_tasks":                        0,
				"thread_pool_GossipStage_pending_tasks":                        0,
				"thread_pool_GossipStage_total_blocked_tasks":                  0,
				"thread_pool_HintsDispatcher_active_tasks":                     0,
				"thread_pool_HintsDispatcher_blocked_tasks":                    0,
				"thread_pool_HintsDispatcher_pending_tasks":                    0,
				"thread_pool_HintsDispatcher_total_blocked_tasks":              0,
				"thread_pool_MemtableFlushWriter_active_tasks":                 0,
				"thread_pool_MemtableFlushWriter_blocked_tasks":                0,
				"thread_pool_MemtableFlushWriter_pending_tasks":                0,
				"thread_pool_MemtableFlushWriter_total_blocked_tasks":          0,
				"thread_pool_MemtablePostFlush_active_tasks":                   0,
				"thread_pool_MemtablePostFlush_blocked_tasks":                  0,
				"thread_pool_MemtablePostFlush_pending_tasks":                  0,
				"thread_pool_MemtablePostFlush_total_blocked_tasks":            0,
				"thread_pool_MemtableReclaimMemory_active_tasks":               0,
				"thread_pool_MemtableReclaimMemory_blocked_tasks":              0,
				"thread_pool_MemtableReclaimMemory_pending_tasks":              0,
				"thread_pool_MemtableReclaimMemory_total_blocked_tasks":        0,
				"thread_pool_MutationStage_active_tasks":                       0,
				"thread_pool_MutationStage_blocked_tasks":                      0,
				"thread_pool_MutationStage_pending_tasks":                      0,
				"thread_pool_MutationStage_total_blocked_tasks":                0,
				"thread_pool_Native-Transport-Requests_active_tasks":           0,
				"thread_pool_Native-Transport-Requests_blocked_tasks":          0,
				"thread_pool_Native-Transport-Requests_pending_tasks":          0,
				"thread_pool_Native-Transport-Requests_total_blocked_tasks":    0,
				"thread_pool_PendingRangeCalculator_active_tasks":              0,
				"thread_pool_PendingRangeCalculator_blocked_tasks":             0,
				"thread_pool_PendingRangeCalculator_pending_tasks":             0,
				"thread_pool_PendingRangeCalculator_total_blocked_tasks":       0,
				"thread_pool_PerDiskMemtableFlushWriter_0_active_tasks":        0,
				"thread_pool_PerDiskMemtableFlushWriter_0_blocked_tasks":       0,
				"thread_pool_PerDiskMemtableFlushWriter_0_pending_tasks":       0,
				"thread_pool_PerDiskMemtableFlushWriter_0_total_blocked_tasks": 0,
				"thread_pool_ReadStage_active_tasks":                           0,
				"thread_pool_ReadStage_blocked_tasks":                          0,
				"thread_pool_ReadStage_pending_tasks":                          0,
				"thread_pool_ReadStage_total_blocked_tasks":                    0,
				"thread_pool_Sampler_active_tasks":                             0,
				"thread_pool_Sampler_blocked_tasks":                            0,
				"thread_pool_Sampler_pending_tasks":                            0,
				"thread_pool_Sampler_total_blocked_tasks":                      0,
				"thread_pool_SecondaryIndexManagement_active_tasks":            0,
				"thread_pool_SecondaryIndexManagement_blocked_tasks":           0,
				"thread_pool_SecondaryIndexManagement_pending_tasks":           0,
				"thread_pool_SecondaryIndexManagement_total_blocked_tasks":     0,
				"thread_pool_ValidationExecutor_active_tasks":                  0,
				"thread_pool_ValidationExecutor_blocked_tasks":                 0,
				"thread_pool_ValidationExecutor_pending_tasks":                 0,
				"thread_pool_ValidationExecutor_total_blocked_tasks":           0,
				"thread_pool_ViewBuildExecutor_active_tasks":                   0,
				"thread_pool_ViewBuildExecutor_blocked_tasks":                  0,
				"thread_pool_ViewBuildExecutor_pending_tasks":                  0,
				"thread_pool_ViewBuildExecutor_total_blocked_tasks":            0,
			},
		},
		"fails if endpoint returns invalid data": {
			prepare: prepareCassandraInvalidData,
		},
		"fails on connection refused": {
			prepare: prepareCassandraConnectionRefused,
		},
		"fails on 404 response": {
			prepare: prepareCassandraResponse404,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c, cleanup := test.prepare()
			defer cleanup()

			require.True(t, c.Init())

			mx := c.Collect()

			assert.Equal(t, test.wantCollected, mx)
		})
	}
}

func prepareCassandra() (c *Cassandra, cleanup func()) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(dataMetrics)
		}))

	c = New()
	c.URL = ts.URL
	return c, ts.Close
}

func prepareCassandraInvalidData() (c *Cassandra, cleanup func()) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("hello and\n goodbye"))
		}))

	c = New()
	c.URL = ts.URL
	return c, ts.Close
}

func prepareCassandraConnectionRefused() (c *Cassandra, cleanup func()) {
	c = New()
	c.URL = "http://127.0.0.1:38001"
	return c, func() {}
}

func prepareCassandraResponse404() (c *Cassandra, cleanup func()) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

	c = New()
	c.URL = ts.URL
	return c, ts.Close
}

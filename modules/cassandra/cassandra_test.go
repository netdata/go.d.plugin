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
	vMetrics, _ = os.ReadFile("testdata/metrics.txt")
)

func Test_TestData(t *testing.T) {
	for name, data := range map[string][]byte{
		"vMetrics": vMetrics,
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
				"cache_hit_ratio":                      84671,
				"cache_hits":                           116000,
				"cache_misses":                         21000,
				"cache_size":                           1560,
				"client_request_failures_reads":        0,
				"client_request_failures_writes":       0,
				"client_request_latency_reads":         1000,
				"client_request_latency_writes":        0,
				"client_request_timeouts_reads":        0,
				"client_request_timeouts_writes":       0,
				"client_request_total_latency_reads":   13002000,
				"client_request_total_latency_writes":  0,
				"client_request_unavailables_reads":    0,
				"client_request_unavailables_writes":   0,
				"compaction_bytes_compacted":           128749,
				"compaction_completed_tasks":           45,
				"compaction_pending_tasks":             0,
				"dropped_messages_one_minute":          0,
				"jvm_gc_cms_count":                     1000,
				"jvm_gc_cms_time":                      59,
				"jvm_gc_parnew_count":                  365000,
				"jvm_gc_parnew_time":                   1101,
				"storage_exceptions":                   0,
				"storage_load":                         145789366,
				"thread_pools_currently_blocked_tasks": 0,
				"thread_pools_total_blocked_tasks":     0,
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
			_, _ = w.Write(vMetrics)
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

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
		"fails on default config": {
			wantFail: true,
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

func TestCachestat_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestCachestat_Cleanup(t *testing.T) {
	assert.NotPanics(t, New().Cleanup)
}

func TestCassandra_Collect(t *testing.T) {
	tests := map[string]struct {
		prepare       func() (c *Cassandra, cleanup func())
		wantCollected map[string]int64
	}{
		"success on valid response": {
			prepare: prepareCassandra,
			wantCollected: map[string]int64{
				"org_apache_cassandra_metrics_clientrequest_count_Read":           1,
				"org_apache_cassandra_metrics_clientrequest_count_Write":          0,
				"org_apache_cassandra_metrics_table_count_CompactionBytesWritten": 7186,
				"org_apache_cassandra_metrics_table_count_HitRate":                87,
				"org_apache_cassandra_metrics_table_count_LiveDiskSpaceUsed":      124808,
				"org_apache_cassandra_metrics_table_count_PendingCompactions":     0,
				"org_apache_cassandra_metrics_table_count_ReadLatency":            100,
				"org_apache_cassandra_metrics_table_count_TotalDiskSpaceUsed":     124808,
				"org_apache_cassandra_metrics_table_count_WriteLatency":           46,
				"system_up_time": 0,
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

			collected := c.Collect()

			if collected != nil && test.wantCollected != nil {
				collected["system_up_time"] = test.wantCollected["system_up_time"]
			}

			assert.Equal(t, test.wantCollected, collected)
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

// SPDX-License-Identifier: GPL-3.0-or-later

package consul

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
	dataAgentChecks, _  = os.ReadFile("testdata/checks.txt")
	dataAgentMetrics, _ = os.ReadFile("testdata/metrics.json")
)

func Test_testDataIsValid(t *testing.T) {
	for name, data := range map[string][]byte{
		"dataAgentChecks":  dataAgentChecks,
		"dataAgentMetrics": dataAgentMetrics,
	} {
		require.NotNilf(t, data, name)
	}
}

func TestConsul_Init(t *testing.T) {
	tests := map[string]struct {
		wantFail bool
		config   Config
	}{
		"success with default": {
			wantFail: false,
			config:   New().Config,
		},
		"fail when URL not set": {
			wantFail: true,
			config: Config{
				HTTP: web.HTTP{
					Request: web.Request{URL: ""},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			consul := New()
			consul.Config = test.config

			if test.wantFail {
				assert.False(t, consul.Init())
			} else {
				assert.True(t, consul.Init())
			}
		})
	}
}

func TestConsul_Check(t *testing.T) {
	tests := map[string]struct {
		wantFail bool
		prepare  func(t *testing.T) (consul *Consul, cleanup func())
	}{
		"success on response from Consul": {
			wantFail: false,
			prepare:  caseConsulResponse,
		},
		"fail on invalid data response": {
			wantFail: true,
			prepare:  caseInvalidDataResponse,
		},
		"fail on connection refused": {
			wantFail: true,
			prepare:  caseConnectionRefused,
		},
		"fail on 404 response": {
			wantFail: true,
			prepare:  case404,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			consul, cleanup := test.prepare(t)
			defer cleanup()

			if test.wantFail {
				assert.False(t, consul.Check())
			} else {
				assert.True(t, consul.Check())
			}
		})
	}
}

func TestConsul_Collect(t *testing.T) {
	tests := map[string]struct {
		prepare         func(t *testing.T) (consul *Consul, cleanup func())
		wantNumOfCharts int
		wantMetrics     map[string]int64
	}{
		"success on response from Consul": {
			prepare:         caseConsulResponse,
			wantNumOfCharts: 11,
			wantMetrics: map[string]int64{
				"consul.autopilot.failure_tolerance":    0,
				"consul.autopilot.healthy.no":           0,
				"consul.autopilot.healthy.yes":          1,
				"consul.client.rpc":                     1,
				"consul.runtime.alloc_bytes":            14164144,
				"consul.runtime.sys_bytes":              34685960,
				"consul.runtime.total_gc_pause_ns":      9367254,
				"consul.server.isLeader.no":             0,
				"consul.server.isLeader.yes":            1,
				"health_check_chk1_critical_status":     0,
				"health_check_chk1_maintenance_status":  0,
				"health_check_chk1_passing_status":      1,
				"health_check_chk1_warning_status":      0,
				"health_check_chk2_critical_status":     1,
				"health_check_chk2_maintenance_status":  0,
				"health_check_chk2_passing_status":      0,
				"health_check_chk2_warning_status":      0,
				"health_check_chk3_critical_status":     1,
				"health_check_chk3_maintenance_status":  0,
				"health_check_chk3_passing_status":      0,
				"health_check_chk3_warning_status":      0,
				"health_check_mysql_critical_status":    1,
				"health_check_mysql_maintenance_status": 0,
				"health_check_mysql_passing_status":     0,
				"health_check_mysql_warning_status":     0,
			},
		},
		"success on response from Consul with filtered checks": {
			prepare:         caseConsulResponseWithFilteredChecks,
			wantNumOfCharts: 8,
			wantMetrics: map[string]int64{
				"consul.autopilot.failure_tolerance":    0,
				"consul.autopilot.healthy.no":           0,
				"consul.autopilot.healthy.yes":          1,
				"consul.client.rpc":                     1,
				"consul.runtime.alloc_bytes":            14164144,
				"consul.runtime.sys_bytes":              34685960,
				"consul.runtime.total_gc_pause_ns":      9367254,
				"consul.server.isLeader.no":             0,
				"consul.server.isLeader.yes":            1,
				"health_check_mysql_critical_status":    1,
				"health_check_mysql_maintenance_status": 0,
				"health_check_mysql_passing_status":     0,
				"health_check_mysql_warning_status":     0,
			},
		},
		"fail on invalid data response": {
			prepare:         caseInvalidDataResponse,
			wantNumOfCharts: 0,
			wantMetrics:     nil,
		},
		"fail on connection refused": {
			prepare:         caseConnectionRefused,
			wantNumOfCharts: 0,
			wantMetrics:     nil,
		},
		"fail on 404 response": {
			prepare:         case404,
			wantNumOfCharts: 0,
			wantMetrics:     nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			consul, cleanup := test.prepare(t)
			defer cleanup()

			mx := consul.Collect()

			require.Equal(t, test.wantMetrics, mx)
			if len(test.wantMetrics) > 0 {
				assert.Equal(t, test.wantNumOfCharts, len(*consul.Charts()))
			}
		})
	}
}

func caseConsulResponse(t *testing.T) (*Consul, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case urlPathAgentChecks:
				_, _ = w.Write(dataAgentChecks)
			case urlPathAgentMetrics:
				_, _ = w.Write(dataAgentMetrics)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}))

	consul := New()
	consul.URL = srv.URL

	require.True(t, consul.Init())

	return consul, srv.Close
}

func caseConsulResponseWithFilteredChecks(t *testing.T) (*Consul, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case urlPathAgentChecks:
				_, _ = w.Write(dataAgentChecks)
			case urlPathAgentMetrics:
				_, _ = w.Write(dataAgentMetrics)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}))

	consul := New()
	consul.URL = srv.URL
	consul.ChecksSelector = "!chk* *"

	require.True(t, consul.Init())

	return consul, srv.Close
}

func caseInvalidDataResponse(t *testing.T) (*Consul, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("hello and\n goodbye"))
		}))

	consul := New()
	consul.URL = srv.URL

	require.True(t, consul.Init())

	return consul, srv.Close
}

func caseConnectionRefused(t *testing.T) (*Consul, func()) {
	t.Helper()
	consul := New()
	require.True(t, consul.Init())

	return consul, func() {}
}

func case404(t *testing.T) (*Consul, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

	consul := New()
	require.True(t, consul.Init())

	return consul, srv.Close
}

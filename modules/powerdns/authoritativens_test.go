package powerdns

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/netdata/go.d.plugin/pkg/tlscfg"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	v430statistics, _     = ioutil.ReadFile("testdata/v4.3.0/statistics.json")
	recursorStatistics, _ = ioutil.ReadFile("testdata/recursor/statistics.json")
)

func Test_testDataIsCorrectlyReadAndValid(t *testing.T) {
	for name, data := range map[string][]byte{
		"v430statistics":     v430statistics,
		"recursorStatistics": recursorStatistics,
	} {
		require.NotNilf(t, data, name)
	}
}

func TestNew(t *testing.T) {
	assert.IsType(t, (*AuthoritativeNS)(nil), New())
}

func TestRecursor_Init(t *testing.T) {
	tests := map[string]struct {
		config   Config
		wantFail bool
	}{
		"success on default config": {
			config: New().Config,
		},
		"fails on unset URL": {
			wantFail: true,
			config: Config{
				HTTP: web.HTTP{
					Request: web.Request{URL: ""},
				},
			},
		},
		"fails on invalid TLSCA": {
			wantFail: true,
			config: Config{
				HTTP: web.HTTP{
					Request: web.Request{
						URL: "http://127.0.0.1:38001",
					},
					Client: web.Client{
						TLSConfig: tlscfg.TLSConfig{TLSCA: "testdata/tls"},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ns := New()
			ns.Config = test.config

			if test.wantFail {
				assert.False(t, ns.Init())
			} else {
				assert.True(t, ns.Init())
			}
		})
	}
}

func TestRecursor_Check(t *testing.T) {
	tests := map[string]struct {
		prepare  func() (p *AuthoritativeNS, cleanup func())
		wantFail bool
	}{
		"success on valid response v4.3.0": {
			prepare: preparePowerDNSAuthoritativeNSV430,
		},
		"fails on response from PowerDNS Recursor": {
			wantFail: true,
			prepare:  preparePowerDNSAuthoritativeNSRecursorData,
		},
		"fails on 404 response": {
			wantFail: true,
			prepare:  preparePowerDNSAuthoritativeNS404,
		},
		"fails on connection refused": {
			wantFail: true,
			prepare:  preparePowerDNSAuthoritativeNSConnectionRefused,
		},
		"fails on response with invalid data": {
			wantFail: true,
			prepare:  preparePowerDNSAuthoritativeNSInvalidData,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			recursor, cleanup := test.prepare()
			defer cleanup()
			require.True(t, recursor.Init())

			if test.wantFail {
				assert.False(t, recursor.Check())
			} else {
				assert.True(t, recursor.Check())
			}
		})
	}
}

func TestRecursor_Charts(t *testing.T) {
	recursor := New()
	require.True(t, recursor.Init())
	assert.NotNil(t, recursor.Charts())
}

func TestRecursor_Cleanup(t *testing.T) {
	assert.NotPanics(t, New().Cleanup)
}

func TestRecursor_Collect(t *testing.T) {
	tests := map[string]struct {
		prepare       func() (p *AuthoritativeNS, cleanup func())
		wantCollected map[string]int64
	}{
		"success on valid response v4.3.0": {
			prepare: preparePowerDNSAuthoritativeNSV430,
			wantCollected: map[string]int64{
				"all-outqueries":                41,
				"answers-slow":                  1,
				"answers0-1":                    1,
				"answers1-10":                   1,
				"answers10-100":                 1,
				"answers100-1000":               1,
				"auth-zone-queries":             1,
				"auth4-answers-slow":            1,
				"auth4-answers0-1":              1,
				"auth4-answers1-10":             5,
				"auth4-answers10-100":           35,
				"auth4-answers100-1000":         1,
				"auth6-answers-slow":            1,
				"auth6-answers0-1":              1,
				"auth6-answers1-10":             1,
				"auth6-answers10-100":           1,
				"auth6-answers100-1000":         1,
				"cache-entries":                 171,
				"cache-hits":                    1,
				"cache-misses":                  1,
				"case-mismatches":               1,
				"chain-resends":                 1,
				"client-parse-errors":           1,
				"concurrent-queries":            1,
				"cpu-msec-thread-0":             439,
				"cpu-msec-thread-1":             445,
				"cpu-msec-thread-2":             466,
				"dlg-only-drops":                1,
				"dnssec-authentic-data-queries": 1,
				"dnssec-check-disabled-queries": 1,
				"dnssec-queries":                1,
				"dnssec-result-bogus":           1,
				"dnssec-result-indeterminate":   1,
				"dnssec-result-insecure":        1,
				"dnssec-result-nta":             1,
				"dnssec-result-secure":          5,
				"dnssec-validations":            5,
				"dont-outqueries":               1,
				"ecs-queries":                   1,
				"ecs-responses":                 1,
				"edns-ping-matches":             1,
				"edns-ping-mismatches":          1,
				"empty-queries":                 1,
				"failed-host-entries":           1,
				"fd-usage":                      32,
				"ignored-packets":               1,
				"ipv6-outqueries":               1,
				"ipv6-questions":                1,
				"malloc-bytes":                  1,
				"max-cache-entries":             1000000,
				"max-mthread-stack":             1,
				"max-packetcache-entries":       500000,
				"negcache-entries":              1,
				"no-packet-error":               1,
				"noedns-outqueries":             1,
				"noerror-answers":               1,
				"noping-outqueries":             1,
				"nsset-invalidations":           1,
				"nsspeeds-entries":              78,
				"nxdomain-answers":              1,
				"outgoing-timeouts":             1,
				"outgoing4-timeouts":            1,
				"outgoing6-timeouts":            1,
				"over-capacity-drops":           1,
				"packetcache-entries":           1,
				"packetcache-hits":              1,
				"packetcache-misses":            1,
				"policy-drops":                  1,
				"policy-result-custom":          1,
				"policy-result-drop":            1,
				"policy-result-noaction":        1,
				"policy-result-nodata":          1,
				"policy-result-nxdomain":        1,
				"policy-result-truncate":        1,
				"qa-latency":                    1,
				"qname-min-fallback-success":    1,
				"query-pipe-full-drops":         1,
				"questions":                     1,
				"real-memory-usage":             44773376,
				"rebalanced-queries":            1,
				"resource-limits":               1,
				"security-status":               3,
				"server-parse-errors":           1,
				"servfail-answers":              1,
				"spoof-prevents":                1,
				"sys-msec":                      1520,
				"tcp-client-overflow":           1,
				"tcp-clients":                   1,
				"tcp-outqueries":                1,
				"tcp-questions":                 1,
				"throttle-entries":              1,
				"throttled-out":                 1,
				"throttled-outqueries":          1,
				"too-old-drops":                 1,
				"truncated-drops":               1,
				"udp-in-errors":                 1,
				"udp-noport-errors":             1,
				"udp-recvbuf-errors":            1,
				"udp-sndbuf-errors":             1,
				"unauthorized-tcp":              1,
				"unauthorized-udp":              1,
				"unexpected-packets":            1,
				"unreachables":                  1,
				"uptime":                        1624,
				"user-msec":                     465,
				"variable-responses":            1,
				"x-our-latency":                 1,
				"x-ourtime-slow":                1,
				"x-ourtime0-1":                  1,
				"x-ourtime1-2":                  1,
				"x-ourtime16-32":                1,
				"x-ourtime2-4":                  1,
				"x-ourtime4-8":                  1,
				"x-ourtime8-16":                 1,
			},
		},
		"fails on response from PowerDNS Recursor": {
			prepare: preparePowerDNSAuthoritativeNSRecursorData,
		},
		"fails on 404 response": {
			prepare: preparePowerDNSAuthoritativeNS404,
		},
		"fails on connection refused": {
			prepare: preparePowerDNSAuthoritativeNSConnectionRefused,
		},
		"fails on response with invalid data": {
			prepare: preparePowerDNSAuthoritativeNSInvalidData,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ns, cleanup := test.prepare()
			defer cleanup()
			require.True(t, ns.Init())

			collected := ns.Collect()

			l := make([]string, 0)
			for k := range collected {
				l = append(l, k)
			}
			sort.Strings(l)
			for _, value := range l {
				fmt.Println(fmt.Sprintf("\"%s\": %d,", value, collected[value]))
			}

			assert.Equal(t, test.wantCollected, collected)
			if len(test.wantCollected) > 0 {
				ensureCollectedHasAllChartsDimsVarsIDs(t, ns, collected)
			}
		})
	}
}

func ensureCollectedHasAllChartsDimsVarsIDs(t *testing.T, ns *AuthoritativeNS, collected map[string]int64) {
	for _, chart := range *ns.Charts() {
		if chart.Obsolete {
			continue
		}
		for _, dim := range chart.Dims {
			_, ok := collected[dim.ID]
			assert.Truef(t, ok, "chart '%s' dim '%s': no dim in collected", dim.ID, chart.ID)
		}
		for _, v := range chart.Vars {
			_, ok := collected[v.ID]
			assert.Truef(t, ok, "chart '%s' dim '%s': no dim in collected", v.ID, chart.ID)
		}
	}
}

func preparePowerDNSAuthoritativeNSV430() (*AuthoritativeNS, func()) {
	srv := preparePowerDNSAuthoritativeNSEndpoint()
	ns := New()
	ns.URL = srv.URL

	return ns, srv.Close
}

func preparePowerDNSAuthoritativeNSRecursorData() (*AuthoritativeNS, func()) {
	srv := preparePowerDNSRecursorEndpoint()
	ns := New()
	ns.URL = srv.URL

	return ns, srv.Close
}

func preparePowerDNSAuthoritativeNSInvalidData() (*AuthoritativeNS, func()) {
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("hello and\n goodbye"))
		}))
	ns := New()
	ns.URL = srv.URL

	return ns, srv.Close
}

func preparePowerDNSAuthoritativeNS404() (*AuthoritativeNS, func()) {
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
	ns := New()
	ns.URL = srv.URL

	return ns, srv.Close
}

func preparePowerDNSAuthoritativeNSConnectionRefused() (*AuthoritativeNS, func()) {
	ns := New()
	ns.URL = "http://127.0.0.1:38001"

	return ns, func() {}
}

func preparePowerDNSAuthoritativeNSEndpoint() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case urlPathLocalStatistics:
				_, _ = w.Write(v430statistics)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}))
}

func preparePowerDNSRecursorEndpoint() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case urlPathLocalStatistics:
				_, _ = w.Write(recursorStatistics)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}))
}

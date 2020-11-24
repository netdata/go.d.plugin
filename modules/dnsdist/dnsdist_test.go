package dnsdist

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/netdata/go.d.plugin/pkg/tlscfg"
	"github.com/netdata/go.d.plugin/pkg/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	dnsdistStatisicsV151, _ = ioutil.ReadFile("testdata/v1.5.1/statistics.json")
)

func Test_testDataIsCorrectlyReadAndValid(t *testing.T) {
	for name, data := range map[string][]byte{
		"dnsdistStatistics": dnsdistStatisicsV151,
	} {
		require.NotNilf(t, data, name)
	}
}

func TestNew(t *testing.T) {
	assert.IsType(t, (*DNSdist)(nil), New())
}

func Test_Init(t *testing.T) {
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

func Test_Charts(t *testing.T) {
	dist := New()
	require.True(t, dist.Init())
	assert.NotNil(t, dist.Charts())
}

func Test_Cleanup(t *testing.T) {
	assert.NotPanics(t, New().Cleanup)
}

func Test_Check(t *testing.T) {
	tests := map[string]struct {
		prepare func() (p *DNSdist, cleanup func())
		wantFail bool
	}{
		"success" : {
			prepare: preparePowerDNSdistV151,
			wantFail: false,
		},
		"fail on 404 response" : {
			prepare: preparePowerDNSdist404,
			wantFail: true,
		},
		"fail on connection refused" : {
			prepare: preparePowerDNSdistConnectionRefused,
			wantFail: true,
		},
		"fail with invalid data" : {
			prepare: preparePowerDNSdistInvalidData,
			wantFail: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			dist, cleanup := test.prepare()
			defer cleanup()

			require.True(t, dist.Init())
			if test.wantFail {
				assert.False(t, dist.Check())
			} else {
				assert.True(t, dist.Check())
			}
		})
	}
}

func Test_Collect(t *testing.T) {
	tests := map[string]struct {
		prepare  func()  (p *DNSdist, cleanup func())
		wantCollected map[string]int64
	}{
	/*	"success" : {
			prepare: preparePowerDNSdistV151,
			wantCollected: map[string]int64 {
				"acl-drops": 0,
				"cache-hits": 0,
				"cache-misses": 0,
				"cpu-iowait": 39284,
				"cpu-steal": 0,
				"cpu-sys-msec": 411,
				"cpu-user-msec": 939,
				"doh-query-pipe-full": 0,
				"doh-response-pipe-full": 0,
				"downstream-send-errors": 0,
				"downstream-timeouts": 0,
				"dyn-block-nmg-size": 0,
				"dyn-blocked": 0,
				"empty-queries": 0,
				"fd-usage": 22,
				"frontend-noerror": 1003,
				"frontend-nxdomain": 0,
				"frontend-servfail": 0,
				"latency-avg100": 14237,
				"latency-avg1000": 9728,
				"latency-avg10000": 1514,
				"latency-avg1000000": 15,
				"latency-count": 1003,
				"latency-slow": 0,
				"latency-sum": 15474,
				"latency0-1": 0,
				"latency1-10": 3,
				"latency10-50": 996,
				"latency100-1000": 4,
				"latency50-100": 0,
				"no-policy": 0,
				"noncompliant-queries": 0,
				"noncompliant-responses": 0,
				"over-capacity-drops": 0,
				"packetcache-hits": 0,
				"packetcache-misses": 0,
				"queries": 1003,
				"rdqueries": 1003,
				"real-memory-usage": 202125312,
				"responses": 1003,
				"rule-drop": 0,
				"rule-nxdomain": 0,
				"rule-refused": 0,
				"rule-servfail": 0,
				"security-status": 0,
				"self-answered": 0,
				"servfail-responses": 0,
				"too-old-drops": 0,
				"trunc-failures": 0,
				"udp-in-errors": 38,
				"udp-noport-errors": 1102,
				"udp-recvbuf-errors": 0,
				"udp-sndbuf-errors": 179,
				"uptime": 394,
			},
		}, */
		"fail on 404 response" : {
			prepare: preparePowerDNSdist404,
		},
		"fail on connection refused" : {
			prepare: preparePowerDNSdistConnectionRefused,
		},
		"fail with invalid data" : {
			prepare: preparePowerDNSdistInvalidData,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			dist, cleanup := test.prepare()
			defer cleanup()

			require.True(t, dist.Init())

			collected := dist.Collect()

			assert.Equal(t, test.wantCollected, collected)
			/*
			if len(test.wantCollected) > 0 {
				ensureCollectedHasAllChartsDimsVarsIDs(t, dist, collected)
			}
			*/

			dist.Cleanup()
		})
	}
}

/*
func ensureCollectedHasAllChartsDimsVarsIDs(t *testing.T, ns *DNSdist, collected map[string]int64) {
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
*/

func preparePowerDNSdistV151() (*DNSdist, func()) {
	srv := preparePowerDNSDistEndpoint()
	ns := New()
	ns.URL = srv.URL

	return ns, srv.Close
}

func preparePowerDNSDistEndpoint() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		// take a look on decoded and clean request
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(r.URL.Path)
			switch r.URL.Path {
			case urlPathLocalStatistics:
				_, _ = w.Write(dnsdistStatisicsV151)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}))
}

func preparePowerDNSdist404() (*DNSdist, func()) {
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
	ns := New()
	ns.URL = srv.URL

	return ns, srv.Close
}

func preparePowerDNSdistConnectionRefused() (*DNSdist, func()) {
	ns := New()
	ns.URL = "http://127.0.0.1:38001"

	return ns, func() {}
}

func preparePowerDNSdistInvalidData() (*DNSdist, func()) {
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("hello and\n goodbye"))
		}))
	ns := New()
	ns.URL = srv.URL

	return ns, srv.Close
}

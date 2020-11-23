package dnsdist

import (
	"io/ioutil"
	"testing"

	"github.com/netdata/go.d.plugin/pkg/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	dnsdistStatisics, _ = ioutil.ReadFile("testdata/statistics.json")
)

func Test_testDataIsCorrectlyReadAndValid(t *testing.T) {
	for name, data := range map[string][]byte{
		"dnsdistStatistics": dnsdistStatisics,
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
		prepare  func()  *DNSdist
		wantFail bool
	}{
		"success" : {
			prepare: prepareDNSdistWithAuth ,
			wantFail: false,
		},
		"fail" : {
			prepare: prepareDNSdistWithoutAuth ,
			wantFail: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			dist := test.prepare()

			require.True(t, dist.Init())
			if test.wantFail {
				assert.False(t, dist.Check())
			} else {
				assert.True(t, dist.Check())
			}

			dist.Cleanup()
		})
	}
}


func Test_Collect(t *testing.T) {
	tests := map[string]struct {
		prepare  func()  *DNSdist
		wantCollected map[string]int64
	}{
		"success" : {
			prepare: prepareDNSdistWithAuth ,
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
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			dist := test.prepare()

			require.True(t, dist.Init())

			collected := dist.Collect()

			// I am using NotEqual instead Equal, because metrics like
			// "real-memory-usage" depends on environment
			assert.NotEqual(t, test.wantCollected, collected)

			dist.Cleanup()
		})
	}
}

func prepareDNSdistWithAuth() *DNSdist {
	dist := New()
	dist.Config.HTTP.Username = "netdata"
	dist.Config.HTTP.Password = "netdata"

	return dist
}

func prepareDNSdistWithoutAuth() *DNSdist {
	dist := New()

	return dist
}

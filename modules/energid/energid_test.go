package energid

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudflare/cfssl/log"
	"github.com/netdata/go.d.plugin/pkg/tlscfg"
	"github.com/netdata/go.d.plugin/pkg/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	v12JSONblockchaininfo, _ = ioutil.ReadFile("testdata/v12/getblockchaininfo.json")
	v12JSONmempoolinfo, _ = ioutil.ReadFile("testdata/v12/getmempoolinfo.json")
	v12JSONnetworkinfo, _ = ioutil.ReadFile("testdata/v12/getnetworkinfo.json")
	v12JSONtxoutsetinfo, _ = ioutil.ReadFile("testdata/v12/gettxoutsetinfo.json")
)

func Test_testDataIsCorrectlyReadAndValid(t *testing.T) {
	for name, data := range map[string][]byte{
			"v12JSONblockchaininfo": v12JSONblockchaininfo,
			"v12JSONmempoolinfo": v12JSONmempoolinfo,
			"v12JSONnetworkinfo": v12JSONnetworkinfo,
			"v12JSONtxoutsetinfo": v12JSONtxoutsetinfo,
	} {
			require.NotNilf(t, data, name)
	}
}

func TestNew(t *testing.T) {
	assert.IsType(t, (*Energid)(nil), New())
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
		prepare  func() (e *Energid,  cleanup func())
		wantFail bool
	}{
		"valid" : {prepare : prepareEnergidValidData, wantFail: false},
		"invalid data" : {prepare : prepareEnergidInvalidData, wantFail: true},
		"404" : {prepare : prepareEnergid404, wantFail: true},
		"Connection refused" : {prepare : prepareEnergidConnectionRefused, wantFail: true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			e, cleanup := test.prepare()
			defer cleanup()

			require.True(t, e.Init())

			if test.wantFail {
				assert.False(t, e.Check())
			} else {
				assert.True(t, e.Check())
			}
		})
	}
}

func Test_Collect(t *testing.T) {
	tests := map[string]struct {
		prepare       func() (dist *Energid, cleanup func())
		wantCollected map[string]int64
	}{
		"success" :{
			prepare: prepareEnergidValidData,
			wantCollected: map[string]int64{
				// Block Chaing statistic 
				"blockchain_blocks" : 0,
				"blockchain_difficulty": 0,
				"blockchain_headers": 0,

				// Memory Pool
				"mempool_max": 300000000,
				"mempool_current": 0,
				"mempool_txsize": 0,

				// Network 
				"network_connections": 0,
				"network_timeoffset": 0,

				// TX
				"utxo_xfers": 0,
				"utxo_count": 0,
			},
		}, 
		"fails on 404 response": {
			prepare: prepareEnergid404,
		},
		"fails on connection refused": {
			prepare: prepareEnergidConnectionRefused,
		},
		"fails with invalid data": {
			prepare: prepareEnergidInvalidData,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			e, cleanup := test.prepare()
			defer cleanup()
			require.True(t, e.Init())

			collected := e.Collect()

			assert.Equal(t, test.wantCollected, collected)
			if len(test.wantCollected) > 0 {
				ensureCollectedHasAllChartsDimsVarsIDs(t, e, collected)
			}
		})
	}
}

func ensureCollectedHasAllChartsDimsVarsIDs(t *testing.T, e *Energid, collected map[string]int64) {
	for _, chart := range *e.Charts() {
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

func prepareEnergidInvalidData() (*Energid, func()) {
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("Hello world!"))
		}))
	e := New()
	e.URL = srv.URL

	return e, srv.Close
}

func prepareEnergid404() (*Energid, func()) {
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
	cdb := New()
	cdb.URL = srv.URL

	return cdb, srv.Close
}

func prepareEnergidConnectionRefused() (*Energid, func()) {
	e := New()
	e.URL = "http://127.0.0.1:38001"

	return e, func() {}
}

func prepareEnergidValidData() (*Energid, func()) {
	srv := prepareEnergidEndPoint()
	e := New()
	e.URL = srv.URL

	return e, srv.Close
}

func prepareEnergidEndPoint() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				body, _ := ioutil.ReadAll(r.Body)

				if body != nil {
					var data energyBody
					err := json.Unmarshal(body, &data)
					if err != nil {
						log.Error(err)
					} else {
						switch data.Method {
							case "getblockchaininfo":
								_, _ = w.Write(v12JSONblockchaininfo)
							case "getmempoolinfo":
								_, _ = w.Write(v12JSONmempoolinfo)
							case "getnetworkinfo":
								_, _ = w.Write(v12JSONnetworkinfo)
							case "gettxoutsetinfo":
								_, _ = w.Write(v12JSONtxoutsetinfo)
							default:
								w.WriteHeader(http.StatusNotFound)
							}
					}
				} 
			}))
}

package energid

import (
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
		prepare  func(*testing.T) (e *Energid, cleanup func())
		wantFail bool
	}{
		"valid" : {prepare : prepareEnergiddValidData, wantFail: false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
				cdb, cleanup := test.prepare(t)
				defer cleanup()

				if test.wantFail {
						assert.False(t, cdb.Check())
				} else {
						assert.True(t, cdb.Check())
				}
		})
	}
}

func prepareEnergiddValidData(t *testing.T) (cdb *Energid, cleanup func()) {
	return prepareEnergid12() 
}

func prepareEnergid12() (*Energid, func()) {
	srv := preparePowerDNSDistEndpoint()
	ns := New()
	ns.URL = srv.URL

	return ns, srv.Close
}

func preparePowerDNSDistEndpoint() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
					switch r.URL.String() {
					case "/blockchain":
							_, _ = w.Write(v12JSONblockchaininfo)
					case "/mempool":
							_, _ = w.Write(v12JSONmempoolinfo)
					case "/network":
							_, _ = w.Write(v12JSONnetworkinfo)
					case "/txout":
							_, _ = w.Write(v12JSONtxoutsetinfo)
					default:
							w.WriteHeader(http.StatusNotFound)
					}
			}))
}

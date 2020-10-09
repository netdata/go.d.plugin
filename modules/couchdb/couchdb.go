package couchdb

import (
	"fmt"
	"net/http"
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	module.Register("couchdb", module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 10,
		},
		Create: func() module.Module { return New() },
	})
}

func New() *CouchDB {
	return &CouchDB{
		Config: Config{
			HTTP: web.HTTP{
				Request: web.Request{
					URL: "http://127.0.0.1:5984",
				},
				Client: web.Client{
					Timeout: web.Duration{Duration: time.Second * 5},
				},
			},
			Node: "nonode@nohost",
		},

		urlPathOverviewStats: "/_node/%s/_stats",
		urlPathSystemStats:   "/_node/%s/_system",
	}
}

type (
	Config struct {
		web.HTTP  `yaml:",inline"`
		Node      string `yaml:"node"`
		Databases string `yaml:"databases"`
	}

	CouchDB struct {
		module.Base
		Config `yaml:",inline"`

		httpClient *http.Client
		charts     *module.Charts

		urlPathOverviewStats string
		urlPathSystemStats   string
	}
)

func (cdb *CouchDB) Cleanup() {
	if cdb.httpClient == nil {
		return
	}
	cdb.httpClient.CloseIdleConnections()
}

func (cdb *CouchDB) Init() bool {
	err := cdb.validateConfig()
	if err != nil {
		cdb.Errorf("check configuration: %v", err)
		return false
	}

	cdb.urlPathOverviewStats = fmt.Sprintf(cdb.urlPathOverviewStats, cdb.Config.Node)
	cdb.urlPathSystemStats = fmt.Sprintf(cdb.urlPathSystemStats, cdb.Config.Node)

	httpClient, err := cdb.initHTTPClient()
	if err != nil {
		cdb.Errorf("init HTTP client: %v", err)
		return false
	}
	cdb.httpClient = httpClient

	charts, err := cdb.initCharts()
	if err != nil {
		cdb.Errorf("init charts: %v", err)
		return false
	}
	cdb.charts = charts

	return true
}

func (cdb *CouchDB) Check() bool {
	if err := cdb.pingCouchDB(); err != nil {
		cdb.Error(err)
		return false
	}
	return true //TODO: len(cdb.Collect()) > 0
}

func (cdb *CouchDB) Charts() *Charts {
	return cdb.charts
}

func (cdb *CouchDB) Collect() map[string]int64 {
	mx, err := cdb.collect()
	if err != nil {
		cdb.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

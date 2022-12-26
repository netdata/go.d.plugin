// SPDX-License-Identifier: GPL-3.0-or-later

package elasticsearch

import (
	"net/http"
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go.d.plugin/agent/module"
)

func init() {
	module.Register("elasticsearch", module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 5,
		},
		Create: func() module.Module { return New() },
	})
}

func New() *Elasticsearch {
	return &Elasticsearch{
		Config: Config{
			HTTP: web.HTTP{
				Request: web.Request{
					URL: "http://127.0.0.1:9200",
				},
				Client: web.Client{
					Timeout: web.Duration{Duration: time.Second * 5},
				},
			},
			DoNodeStats:     true,
			DoClusterStats:  true,
			DoClusterHealth: true,
			DoIndicesStats:  false,
		},
		indices: make(map[string]bool),
	}
}

type Config struct {
	web.HTTP        `yaml:",inline"`
	DoNodeStats     bool `yaml:"collect_node_stats"`
	DoClusterHealth bool `yaml:"collect_cluster_health"`
	DoClusterStats  bool `yaml:"collect_cluster_stats"`
	DoIndicesStats  bool `yaml:"collect_indices_stats"`
}

type Elasticsearch struct {
	module.Base
	Config `yaml:",inline"`

	httpClient *http.Client
	charts     *module.Charts

	indices map[string]bool
}

func (es *Elasticsearch) Init() bool {
	err := es.validateConfig()
	if err != nil {
		es.Errorf("check configuration: %v", err)
		return false
	}

	httpClient, err := es.initHTTPClient()
	if err != nil {
		es.Errorf("init HTTP client: %v", err)
		return false
	}
	es.httpClient = httpClient

	charts, err := es.initCharts()
	if err != nil {
		es.Errorf("init charts: %v", err)
		return false
	}
	es.charts = charts

	return true
}

func (es *Elasticsearch) Check() bool {
	if err := es.pingElasticsearch(); err != nil {
		es.Error(err)
		return false
	}
	return len(es.Collect()) > 0
}

func (es *Elasticsearch) Charts() *module.Charts {
	return es.charts
}

func (es *Elasticsearch) Collect() map[string]int64 {
	mx, err := es.collect()
	if err != nil {
		es.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (es *Elasticsearch) Cleanup() {
	if es.httpClient != nil {
		es.httpClient.CloseIdleConnections()
	}
}

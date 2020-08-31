package elasticsearch

import (
	"net/http"
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
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
					UserURL: "http://127.0.0.1:9200",
					//UserURL: "http://192.168.88.250:9200/",
				},
				Client: web.Client{
					Timeout: web.Duration{Duration: time.Second * 5},
				},
			},
		},
		collectedIndices: make(map[string]bool),
	}
}

type (
	Config struct {
		web.HTTP `yaml:",inline"`
	}
	Elasticsearch struct {
		module.Base
		Config `yaml:",inline"`

		httpClient       *http.Client
		charts           *module.Charts
		collectedIndices map[string]bool
	}
)

func (es *Elasticsearch) Cleanup() {
	if es.httpClient == nil {
		return
	}
	es.httpClient.CloseIdleConnections()
}

func (es *Elasticsearch) Init() bool {
	if es.UserURL == "" {
		es.Error("URL not set")
		return false
	}

	client, err := web.NewHTTPClient(es.Client)
	if err != nil {
		es.Errorf("init HTTP client: %v", err)
		return false
	}
	es.httpClient = client

	es.charts = newLocalNodeCharts()
	if err := es.charts.Add(*newClusterHealthCharts()...); err != nil {
		es.Errorf("init charts: add cluster health charts: $v", err)
		return false
	}
	if err := es.charts.Add(*newClusterStatsCharts()...); err != nil {
		es.Errorf("init charts: add cluster stats charts: $v", err)
		return false
	}

	return true
}

func (es *Elasticsearch) Check() bool {
	return len(es.Collect()) > 0
}

func (es *Elasticsearch) Charts() *module.Charts {
	return es.charts
}

func (es Elasticsearch) Collect() map[string]int64 {
	mx, err := es.collect()
	if err != nil {
		es.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

package elasticsearch

import (
	"net/http"
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 5,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("elasticsearch", creator)
}

func New() *Elasticsearch {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{UserURL: "http://127.0.0.1:9200"},
			Client:  web.Client{Timeout: web.Duration{Duration: time.Second * 5}},
		},
	}
	return &Elasticsearch{
		Config: config,
	}
}

type (
	Config struct {
		web.HTTP `yaml:",inline"`
	}
	Elasticsearch struct {
		module.Base
		Config `yaml:",inline"`

		httpClient *http.Client
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

	return true
}

func (es *Elasticsearch) Check() bool {
	return len(es.Collect()) > 0
}

func (es Elasticsearch) Charts() *module.Charts {
	return nil
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

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

func (Elasticsearch) Cleanup() {}

func (e *Elasticsearch) Init() bool {
	if e.UserURL == "" {
		e.Error("URL not set")
		return false
	}

	client, err := web.NewHTTPClient(e.Client)
	if err != nil {
		e.Errorf("init HTTP client: %v", err)
		return false
	}
	e.httpClient = client

	return true
}

func (e *Elasticsearch) Check() bool {
	return len(e.Collect()) > 0
}

func (e Elasticsearch) Charts() *module.Charts {
	return nil
}

func (e Elasticsearch) Collect() map[string]int64 {
	mx, err := e.collect()
	if err != nil {
		e.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

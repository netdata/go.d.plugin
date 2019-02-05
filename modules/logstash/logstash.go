package logstash

import (
	"time"

	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("logstash", creator)
}

const (
	defaultURL         = "http://localhost:9600"
	defaultHTTPTimeout = time.Second
)

// New creates Logstash with default values.
func New() *Logstash {

	return &Logstash{
		HTTP: web.HTTP{
			Request: web.Request{URL: defaultURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
	}
}

// Logstash logstash module.
type Logstash struct {
	module.Base
	web.HTTP  `yaml:",inline"`
	apiClient *apiClient
}

// Cleanup makes cleanup.
func (Logstash) Cleanup() {}

// Init makes initialization.
func (l *Logstash) Init() bool {
	if l.URL == "" {
		l.Error("URL is not set")
		return false
	}

	client, err := web.NewHTTPClient(l.Client)

	if err != nil {
		l.Error(err)
		return false
	}

	l.apiClient = &apiClient{
		req:        l.Request,
		httpClient: client,
	}

	l.Debugf("using URL %s", l.URL)
	l.Debugf("using timeout: %s", l.Timeout.Duration)

	return true
}

// Check makes check.
func (l *Logstash) Check() bool {
	return len(l.Collect()) > 0
}

// Charts creates Charts.
func (Logstash) Charts() *Charts {
	return charts.Copy()
}

// Collect collects metrics.
func (l *Logstash) Collect() map[string]int64 {
	jvmStats, err := l.apiClient.jvmStats()

	if err != nil {
		l.Error(err)
		return nil
	}

	return stm.ToMap(jvmStats)
}

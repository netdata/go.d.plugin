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
	defaultHTTPTimeout = time.Second * 2
)

// New creates Logstash with default values.
func New() *Logstash {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{UserURL: defaultURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
	}

	return &Logstash{Config: config}
}

// Config is the Logstash module configuration.
type Config struct {
	web.HTTP `yaml:",inline"`
}

// Logstash Logstash module.
type Logstash struct {
	module.Base
	Config    `yaml:",inline"`
	apiClient *apiClient
	charts    *Charts
}

// Cleanup makes cleanup.
func (Logstash) Cleanup() {}

// Init makes initialization.
func (l *Logstash) Init() bool {
	if err := l.ParseUserURL(); err != nil {
		l.Errorf("error on parsing url '%s' : %v", l.Request.UserURL, err)
		return false
	}

	if l.URL.Host == "" {
		l.Error("URL is not set")
		return false
	}

	client, err := web.NewHTTPClient(l.Client)

	if err != nil {
		l.Error(err)
		return false
	}

	l.apiClient = newAPIClient(client, l.Request)

	l.Debugf("using URL %s", l.URL)
	l.Debugf("using timeout: %s", l.Timeout.Duration)

	return true
}

// Check makes check.
func (l *Logstash) Check() bool { return len(l.Collect()) > 0 }

// Charts creates Charts.
func (l *Logstash) Charts() *Charts {
	if l.charts == nil {
		l.charts = charts.Copy()
	}
	return l.charts
}

// Collect collects metrics.
func (l *Logstash) Collect() map[string]int64 {
	jvmStats, err := l.apiClient.jvmStats()

	if err != nil {
		l.Error(err)
		return nil
	}

	for id := range jvmStats.Pipelines {
		chartID := "pipeline_" + id + "_event"
		if !l.Charts().Has(chartID) {
			l.Charts().Add(createPipelineChart(id)...)
		}
	}

	return stm.ToMap(jvmStats)
}

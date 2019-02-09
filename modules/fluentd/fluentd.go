package fluentd

import (
	"time"

	"github.com/netdata/go.d.plugin/pkg/matcher"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("fluentd", creator)
}

const (
	defaultURL         = "http://127.0.0.1:24220"
	defaultHTTPTimeout = time.Second * 2
)

type Config struct {
	web.HTTP         `yaml:",inline"`
	PermitPluginType string `yaml:"permit_plugin_type"`
}

// New creates Fluentd with default values.
func New() *Fluentd {
	return &Fluentd{
		Config: Config{
			HTTP: web.HTTP{
				Request: web.Request{URL: defaultURL},
				Client:  web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
			}},
		charts: &Charts{},
	}
}

// Fluentd Fluentd module.
type Fluentd struct {
	module.Base
	Config `yaml:",inline"`

	apiClient        *apiClient
	permitPluginType matcher.Matcher
	charts           *Charts
}

// Cleanup makes cleanup.
func (Fluentd) Cleanup() {}

// Init makes initialization.
func (f *Fluentd) Init() bool {
	if f.URL == "" {
		f.Error("URL is not set")
		return false
	}

	if f.PermitPluginType != "" {
		m, err := matcher.NewSimplePatternsMatcher(f.PermitPluginType)
		if err != nil {
			f.Errorf("error on creating permit_plugin matcher : %v", err)
			return false
		}
		f.permitPluginType = matcher.WithCache(m)
	}

	client, err := web.NewHTTPClient(f.Client)

	if err != nil {
		f.Errorf("error on creating client : %v", err)
		return false
	}

	f.apiClient = newAPIClient(client, f.Request)

	f.Debugf("using URL %s", f.URL)
	f.Debugf("using timeout: %s", f.Timeout.Duration)

	return true
}

// Check makes check.
func (f Fluentd) Check() bool { return len(f.Collect()) > 0 }

// Charts creates Charts.
func (f Fluentd) Charts() *Charts { return f.charts }

// Collect collects metrics.
func (f *Fluentd) Collect() map[string]int64 {
	info, err := f.apiClient.getPluginsInfo()

	if err != nil {
		f.Error(err)
		return nil
	}

	metrics := make(map[string]int64)

	for _, p := range info.Payload {
		if p.RetryCount == nil && p.BufferQueueLength == nil && p.BufferTotalQueuedSize == nil {
			continue
		}

		if !f.permitPluginType.MatchString(p.Type) {
			continue
		}

		f.collectPlugin(metrics, p)
	}

	return metrics
}

func (f *Fluentd) collectPlugin(metrics map[string]int64, plugin pluginData) {

}

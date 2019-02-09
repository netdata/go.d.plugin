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
	web.HTTP     `yaml:",inline"`
	PermitPlugin string `yaml:"permit_plugin"`
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

	apiClient    *apiClient
	permitPlugin matcher.Matcher
	charts       *Charts
}

// Cleanup makes cleanup.
func (Fluentd) Cleanup() {}

// Init makes initialization.
func (f *Fluentd) Init() bool {
	if f.URL == "" {
		f.Error("URL is not set")
		return false
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
	return nil
}

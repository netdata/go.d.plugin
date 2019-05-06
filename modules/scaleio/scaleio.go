package scaleio

import (
	"time"

	"github.com/netdata/go.d.plugin/modules/scaleio/api"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("scaleio", creator)
}

const (
	// defaultURL         = "https://127.0.0.1"
	defaultURL         = "https://100.127.0.201"
	defaultHTTPTimeout = time.Second * 2
)

// New creates ScaleIO with default values.
func New() *ScaleIO {
	config := Config{
		HTTP: web.HTTP{
			// Request: web.Request{UserURL: defaultURL},
			Request: web.Request{UserURL: defaultURL, Username: "admin", Password: "Admin1!"},
			Client: web.Client{
				Timeout:         web.Duration{Duration: defaultHTTPTimeout},
				ClientTLSConfig: web.ClientTLSConfig{InsecureSkipVerify: true}, // TODO: remove
			},
		},
	}
	return &ScaleIO{Config: config}
}

// Config is the ScaleIO module configuration.
type Config struct {
	web.HTTP `yaml:",inline"`
}

type scaleIOAPIClient interface {
	Login() error
	Logout() error
	IsLoggedIn() bool
	GetSelectedStatistics(dst interface{}, query string) error
}

// ScaleIO ScaleIO module.
type ScaleIO struct {
	module.Base
	Config    `yaml:",inline"`
	apiClient scaleIOAPIClient
}

// Cleanup makes cleanup.
func (s *ScaleIO) Cleanup() {
	if s.apiClient == nil {
		return
	}
	_ = s.apiClient.Logout()
}

// Init makes initialization.
func (s *ScaleIO) Init() bool {
	if err := s.ParseUserURL(); err != nil {
		s.Errorf("error on parsing URL '%s' : %v", s.UserURL, err)
		return false
	}

	if s.URL.Host == "" {
		s.Error("URL is not set")
		return false
	}

	client, err := web.NewHTTPClient(s.Client)

	if err != nil {
		s.Errorf("error on creating http client : %v", err)
		return false
	}

	s.apiClient = api.NewClient(client, s.Request)

	s.Debugf("using URL %s", s.URL)
	s.Debugf("using timeout: %s", s.Timeout.Duration)

	return true
}

// Check makes check.
func (s *ScaleIO) Check() bool {
	if err := s.apiClient.Login(); err != nil {
		s.Error(err)
		return false
	}
	return len(s.Collect()) > 0
}

// Charts returns Charts.
func (s ScaleIO) Charts() *module.Charts { return charts.Copy() }

// Collect collects metrics.
func (s *ScaleIO) Collect() map[string]int64 {
	mx, err := s.collect()

	if err != nil {
		s.Error(err)
		return nil
	}

	return mx
}

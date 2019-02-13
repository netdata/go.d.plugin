package kubernetes

import (
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("kubernetes", creator)
}

const (
	defaultHTTPTimeout = time.Second * 2
	defaultURL         = "http://127.0.0.1:10255"
)

type Config struct {
	web.HTTP `yaml:",inline"`
}

// New creates Kubernetes with default values.
func New() *Kubernetes {
	return &Kubernetes{
		Config: Config{
			HTTP: web.HTTP{
				Request: web.Request{URL: defaultURL},
				Client:  web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
			},
		},
		activePods: make(map[string]bool),
	}
}

// Kubernetes Kubernetes module.
type Kubernetes struct {
	module.Base
	Config `yaml:",inline"`

	apiClient apiClient
	// TODO: likely wrong
	activePods map[string]bool
}

// Cleanup makes cleanup.
func (Kubernetes) Cleanup() {}

// Init makes initialization.
func (Kubernetes) Init() bool { return false }

// Check makes check.
func (Kubernetes) Check() bool { return true }

// Charts creates Charts.
func (Kubernetes) Charts() *Charts { return nil }

// Collect collects metrics.
func (Kubernetes) Collect() map[string]int64 { return nil }

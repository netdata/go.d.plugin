package kubernetes

import (
	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("kubernetes", creator)
}

// New creates Kubernetes with default values.
func New() *Kubernetes { return &Kubernetes{} }

// Kubernetes Kubernetes module.
type Kubernetes struct{ module.Base }

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

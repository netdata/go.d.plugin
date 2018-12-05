package example

import (
	"math/rand"

	"github.com/netdata/go.d.plugin/modules"
)

// New creates Example with default values
func New() *Example {
	return &Example{
		metrics: make(map[string]int64),
	}
}

// Example example module
type Example struct {
	modules.Base // should be embedded by every module

	metrics map[string]int64
}

// Cleanup makes cleanup
func (Example) Cleanup() {}

// Init makes initialization
func (Example) Init() bool {
	return true
}

// Check makes check
func (Example) Check() bool {
	return true
}

// Charts creates Charts
func (Example) Charts() *Charts {
	return charts.Copy()
}

// GatherMetrics gathers metrics
func (e *Example) GatherMetrics() map[string]int64 {
	e.metrics["random0"] = rand.Int63n(100)
	e.metrics["random1"] = rand.Int63n(100)

	return e.metrics
}

func init() {
	creator := modules.Creator{
		DisabledByDefault: true,
		Create:            func() modules.Module { return New() },
	}

	modules.Register("example", creator)
}

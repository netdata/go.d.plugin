package example

import (
	"math/rand"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			Disabled: true,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("example", creator)
}

// New creates Example with default values
func New() *Example {
	return &Example{
		metrics: make(map[string]int64),
	}
}

// Example example module
type Example struct {
	module.Base // should be embedded by every module

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

// Collect collects metrics
func (e *Example) Collect() map[string]int64 {
	e.metrics["random0"] = rand.Int63n(100)
	e.metrics["random1"] = rand.Int63n(100)

	return e.metrics
}

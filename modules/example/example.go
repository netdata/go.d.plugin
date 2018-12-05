package example

import (
	"math/rand"

	"github.com/netdata/go.d.plugin/modules"
)

// Example module
type Example struct {
	modules.Base // should be embedded by every module

	data map[string]int64
}

// New returns Example with default values
func New() *Example {
	return &Example{
		data: make(map[string]int64),
	}
}

// Cleanup makes cleanup
func (Example) Cleanup() {}

// Init makes initialization
func (Example) Init() bool {
	return true
}

// Check makes check
func (e *Example) Check() bool {
	return true
}

// Charts creates Charts
func (Example) Charts() *Charts {
	return charts.Copy()
}

// GatherMetrics gathers metrics
func (e *Example) GatherMetrics() map[string]int64 {
	e.data["random0"] = rand.Int63n(100)
	e.data["random1"] = rand.Int63n(100)

	return e.data
}

func init() {
	creator := modules.Creator{
		DisabledByDefault: true,
		Create:            func() modules.Module { return New() },
	}

	modules.Register("example", creator)
}

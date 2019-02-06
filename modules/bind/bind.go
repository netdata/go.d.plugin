package bind

import (
	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		// DisabledByDefault: true,
		Create: func() module.Module { return New() },
	}

	module.Register("example", creator)
}

// New creates Bind with default values.
func New() *Bind {
	return &Bind{}
}

// Bind bind module.
type Bind struct {
	module.Base
}

// Cleanup makes cleanup.
func (Bind) Cleanup() {}

// Init makes initialization.
func (Bind) Init() bool {
	return true
}

// Check makes check.
func (Bind) Check() bool {
	return true
}

// Charts creates Charts.
func (Bind) Charts() *Charts {
	return nil
}

// Collect collects metrics.
func (Bind) Collect() map[string]int64 {
	return nil
}

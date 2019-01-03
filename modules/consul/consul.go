package consul

import (
	"github.com/netdata/go.d.plugin/modules"
)

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("consul", creator)
}

// New creates Consul with default values
func New() *Consul {
	return &Consul{}
}

// Consul consul module
type Consul struct {
	modules.Base
}

// Cleanup makes cleanup
func (Consul) Cleanup() {}

// Init makes initialization
func (Consul) Init() bool {
	return false
}

// Check makes check
func (Consul) Check() bool {
	return false
}

// Charts creates Charts
func (Consul) Charts() *Charts {
	return nil
}

// Collect collects metrics
func (c *Consul) Collect() map[string]int64 {
	return nil
}

package freeradius

import (
	"github.com/netdata/go.d.plugin/modules"
)

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("freeradius", creator)
}

// New creates Example with default values
func New() *Freeradius {
	return &Freeradius{}
}

// Freeradius freeradius module
type Freeradius struct {
	modules.Base
}

// Cleanup makes cleanup
func (Freeradius) Cleanup() {}

// Init makes initialization
func (Freeradius) Init() bool {
	return false
}

// Check makes check
func (Freeradius) Check() bool {
	return false
}

// Charts creates Charts
func (Freeradius) Charts() *Charts {
	return charts.Copy()
}

// Collect collects metrics
func (Freeradius) Collect() map[string]int64 {
	return nil
}

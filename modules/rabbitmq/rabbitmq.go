package rabbitmq

import (
	"github.com/netdata/go.d.plugin/modules"
)

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("rabbitmq", creator)
}

// New creates Rabbitmq with default values
func New() *Rabbitmq {
	return &Rabbitmq{}
}

// Rabbitmq rabbitmq module
type Rabbitmq struct {
	modules.Base // should be embedded by every module

}

// Cleanup makes cleanup
func (Rabbitmq) Cleanup() {}

// Init makes initialization
func (Rabbitmq) Init() bool {
	return false
}

// Check makes check
func (Rabbitmq) Check() bool {
	return false
}

// Charts creates Charts
func (Rabbitmq) Charts() *Charts {
	return charts.Copy()
}

// GatherMetrics gathers metrics
func (Rabbitmq) GatherMetrics() map[string]int64 {
	return nil
}

package openvpn

import (
	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("openvpn", creator)
}

// New creates OpenVPN with default values.
func New() *OpenVPN {
	return &OpenVPN{}
}

// OpenVPN OpenVPN module.
type OpenVPN struct {
	module.Base // should be embedded by every module

	apiClient apiClient
}

// Cleanup makes cleanup.
func (OpenVPN) Cleanup() {}

// Init makes initialization.
func (OpenVPN) Init() bool { return false }

// Check makes check.
func (OpenVPN) Check() bool { return false }

// Charts creates Charts.
func (OpenVPN) Charts() *Charts { return charts.Copy() }

// Collect collects metrics.
func (o *OpenVPN) Collect() map[string]int64 {
	mx, err := o.collect()

	if err != nil {
		o.Error(err)
		return nil
	}

	return mx
}

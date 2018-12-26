package freeradius

import (
	"context"
	"time"

	"github.com/netdata/go.d.plugin/modules"

	"layeh.com/radius"
)

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("freeradius", creator)
}

// New creates Freeradius with default values
func New() *Freeradius {
	return &Freeradius{
		Address: "127.0.0.1",
		Port:    18121,
		Secret:  "adminsecret",

		exchanger: &radius.Client{
			Retry:           time.Second,
			MaxPacketErrors: 10,
		},
	}
}

type exchanger interface {
	Exchange(ctx context.Context, packet *radius.Packet, address string) (*radius.Packet, error)
}

// Freeradius freeradius module
type Freeradius struct {
	modules.Base

	Address string
	Port    int
	Secret  string

	exchanger exchanger
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

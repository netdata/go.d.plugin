package freeradius

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"time"

	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/web"

	"layeh.com/radius"
	"layeh.com/radius/rfc2869"
)

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("freeradius", creator)
}

const (
	defAddress = "127.0.0.1"
	defPort    = 18121
	defSecret  = "adminsecret"
)

// New creates Freeradius with default values
func New() *Freeradius {
	return &Freeradius{
		Address: defAddress,
		Port:    defPort,
		Secret:  defSecret,
		Timeout: web.Duration{Duration: time.Second},

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
	Timeout web.Duration

	exchanger exchanger
}

// Cleanup makes cleanup
func (Freeradius) Cleanup() {}

// Init makes initialization
func (f *Freeradius) Init() bool {
	if f.Address == "" {
		f.Address = defAddress
	}
	if f.Port == 0 {
		f.Port = defPort
	}
	if f.Secret == "" {
		f.Secret = defSecret
	}

	return true
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

func newStatusServerPacket(secret string) (*radius.Packet, error) {
	packet := radius.New(radius.CodeStatusServer, []byte(secret))

	if err := PacketType_Set(packet, PacketType_Value_AccessAccept); err != nil {
		return nil, err
	}
	if err := FreeRADIUSStatisticsType_Set(packet, FreeRADIUSStatisticsType_Value_All); err != nil {
		return nil, err
	}

	if err := rfc2869.MessageAuthenticator_Set(packet, make([]byte, 16)); err != nil {
		return nil, err
	}

	hash := hmac.New(md5.New, packet.Secret)

	encode, err := packet.Encode()

	if err != nil {
		return nil, err
	}

	hash.Write(encode)

	if err := rfc2869.MessageAuthenticator_Set(packet, hash.Sum(nil)); err != nil {
		return nil, err
	}

	return packet, nil
}

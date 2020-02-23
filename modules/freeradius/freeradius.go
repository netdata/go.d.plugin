package freeradius

import (
	"context"
	"errors"
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
	"layeh.com/radius"
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			Disabled: true,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("freeradius", creator)
}

func New() *FreeRADIUS {
	cfg := Config{
		Address: "127.0.0.1",
		Port:    18121,
		Secret:  "adminsecret",
		Timeout: web.Duration{Duration: time.Second},
	}
	return &FreeRADIUS{
		Config: cfg,
	}
}

type (
	client interface {
		Exchange(ctx context.Context, packet *radius.Packet, address string) (*radius.Packet, error)
	}
	Config struct {
		Address string
		Port    int
		Secret  string
		Timeout web.Duration
	}
	FreeRADIUS struct {
		module.Base
		Config `yaml:",inline"`
		client
	}
)

func (f FreeRADIUS) validateConfig() error {
	if f.Address == "" {
		return errors.New("address not set")
	}
	if f.Port == 0 {
		return errors.New("port not set")
	}
	if f.Secret == "" {
		return errors.New("secret not set")
	}
	return nil
}

func (f *FreeRADIUS) initClient() {
	f.client = &radius.Client{
		Retry:           time.Second,
		MaxPacketErrors: 10,
	}
}

func (f *FreeRADIUS) Init() bool {
	err := f.validateConfig()
	if err != nil {
		f.Errorf("error on validating config: %v", err)
		return false
	}

	f.initClient()
	return true
}

func (f FreeRADIUS) Check() bool {
	return len(f.Collect()) > 0
}

func (FreeRADIUS) Charts() *Charts {
	return charts.Copy()
}

func (f *FreeRADIUS) Collect() map[string]int64 {
	mx, err := f.collect()
	if err != nil {
		f.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (FreeRADIUS) Cleanup() {}

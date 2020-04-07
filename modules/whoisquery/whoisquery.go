package whoisquery

import (
	"errors"
	"time"

	"github.com/netdata/go-orchestrator/module"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 60,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("whoisquery", creator)
}

func New() *WhoisQuery {
	return &WhoisQuery{
		Config: Config{
			Timeout:       web.Duration{Duration: time.Second * 2},
			DaysUntilWarn: 14,
			DaysUntilCrit: 7,
		},
	}
}

type Config struct {
	web.ClientTLSConfig `yaml:",inline"`
	Timeout             web.Duration
	Source              string
	DaysUntilWarn       int64 `yaml:"days_until_expiration_warning"`
	DaysUntilCrit       int64 `yaml:"days_until_expiration_critical"`
}

type WhoisQuery struct {
	module.Base
	Config `yaml:",inline"`
	prov   provider
}

func (x WhoisQuery) validateConfig() error {
	if x.Source == "" {
		return errors.New("source is not set")
	}
	return nil
}

func (x *WhoisQuery) initProvider() error {
	p, err := newProvider(x.Config)
	if err != nil {
		return err
	}

	x.prov = p
	return nil
}

func (x *WhoisQuery) Init() bool {
	if err := x.validateConfig(); err != nil {
		x.Errorf("error on validating config: %v", err)
		return false
	}

	if err := x.initProvider(); err != nil {
		x.Errorf("error on initializing whois provider: %v", err)
		return false
	}
	return true
}

func (x *WhoisQuery) Check() bool {
	return len(x.Collect()) > 0
}

func (x WhoisQuery) Charts() *Charts {
	return charts.Copy()
}

func (x *WhoisQuery) Collect() map[string]int64 {
	mx, err := x.collect()
	if err != nil {
		x.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (WhoisQuery) Cleanup() {}

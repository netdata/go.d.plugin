package whoisquery

import (
	"errors"

	"github.com/netdata/go.d.plugin/agent/module"
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
			DaysUntilWarn: 90,
			DaysUntilCrit: 30,
		},
	}
}

type Config struct {
	Source        string
	DaysUntilWarn int64 `yaml:"days_until_expiration_warning"`
	DaysUntilCrit int64 `yaml:"days_until_expiration_critical"`
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

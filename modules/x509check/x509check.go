package x509check

import (
	"errors"
	"time"

	"github.com/netdata/go.d.plugin/pkg/tlscfg"
	"github.com/netdata/go.d.plugin/pkg/web"

	cfssllog "github.com/cloudflare/cfssl/log"
	"github.com/netdata/go.d.plugin/agent/module"
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 60,
		},
		Create: func() module.Module { return New() },
	}

	cfssllog.Level = cfssllog.LevelFatal
	module.Register("x509check", creator)
}

func New() *X509Check {
	return &X509Check{
		Config: Config{
			Timeout:           web.Duration{Duration: time.Second * 2},
			DaysUntilWarn:     14,
			DaysUntilCritical: 7,
		},
	}
}

type Config struct {
	Source            string
	Timeout           web.Duration
	tlscfg.TLSConfig  `yaml:",inline"`
	DaysUntilWarn     int64 `yaml:"days_until_expiration_warning"`
	DaysUntilCritical int64 `yaml:"days_until_expiration_critical"`
	CheckRevocation   bool  `yaml:"check_revocation_status"`
}

type X509Check struct {
	module.Base
	Config `yaml:",inline"`
	prov   provider
}

func (x X509Check) validateConfig() error {
	if x.Source == "" {
		return errors.New("source is not set")
	}
	return nil
}

func (x *X509Check) initProvider() error {
	p, err := newProvider(x.Config)
	if err != nil {
		return err
	}

	x.prov = p
	return nil
}

func (x *X509Check) Init() bool {
	if err := x.validateConfig(); err != nil {
		x.Errorf("error on validating config: %v", err)
		return false
	}

	if err := x.initProvider(); err != nil {
		x.Errorf("error on initializing certificate provider: %v", err)
		return false
	}
	return true
}

func (x *X509Check) Check() bool {
	return len(x.Collect()) > 0
}

func (x X509Check) Charts() *Charts {
	if x.CheckRevocation {
		return withRevocationCharts.Copy()
	}
	return charts.Copy()
}

func (x *X509Check) Collect() map[string]int64 {
	mx, err := x.collect()
	if err != nil {
		x.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (X509Check) Cleanup() {}

package x509check

import (
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

const (
	defaultConnTimeout   = time.Second * 2
	defaultDaysUntilWarn = 14
	defaultDaysUntilCrit = 7
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("x509check", creator)
}

// New creates X509Check with default values
func New() *X509Check {
	return &X509Check{
		Config: Config{
			Timeout:       web.Duration{Duration: defaultConnTimeout},
			DaysUntilWarn: defaultDaysUntilWarn,
			DaysUntilCrit: defaultDaysUntilCrit,
		},
	}
}

// Config is the x509Check module configuration.
type Config struct {
	web.ClientTLSConfig `yaml:",inline"`
	Timeout             web.Duration
	Source              string
	DaysUntilWarn       int `yaml:"days_until_expiration_warning"`
	DaysUntilCrit       int `yaml:"days_until_expiration_critical"`
}

// X509Check X509Check module.
type X509Check struct {
	module.Base
	Config `yaml:",inline"`
	certGetter
}

// Cleanup makes cleanup.
func (X509Check) Cleanup() {}

// Init makes initialization.
func (x *X509Check) Init() bool {
	getter, err := newCertGetter(x.Config)

	if err != nil {
		x.Error(err)
		return false
	}

	x.certGetter = getter

	return true
}

// Check makes check.
func (x *X509Check) Check() bool {
	return len(x.Collect()) > 0
}

// Charts creates Charts.
func (X509Check) Charts() *Charts {
	return charts.Copy()
}

// Collect collects metrics.
func (x *X509Check) Collect() map[string]int64 {
	certs, err := x.getCert()

	if err != nil {
		x.Error(err)
		return nil
	}

	if len(certs) == 0 {
		x.Error("no certificate was provided by '%s'", x.Config.Source)
		return nil
	}

	now := time.Now()
	notAfter := certs[0].NotAfter

	return map[string]int64{
		"time":                           int64(notAfter.Sub(now).Seconds()),
		"days_until_expiration_warning":  int64(x.DaysUntilWarn),
		"days_until_expiration_critical": int64(x.DaysUntilCrit),
	}
}

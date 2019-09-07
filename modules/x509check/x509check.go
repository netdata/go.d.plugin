package x509check

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/netdata/go.d.plugin/modules/x509check/cert"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 60,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("x509check", creator)
}

// New creates X509Check with default values
func New() *X509Check {
	return &X509Check{
		Config: Config{
			Timeout:       web.Duration{Duration: time.Second * 2},
			DaysUntilWarn: 14,
			DaysUntilCrit: 7,
		},
	}
}

type gatherer interface {
	Gather() ([]*x509.Certificate, error)
}

// Config is the x509Check module configuration.
type Config struct {
	web.ClientTLSConfig `yaml:",inline"`
	Timeout             web.Duration
	Source              string
	DaysUntilWarn       int64 `yaml:"days_until_expiration_warning"`
	DaysUntilCrit       int64 `yaml:"days_until_expiration_critical"`
}

// X509Check X509Check module.
type X509Check struct {
	module.Base
	Config `yaml:",inline"`
	gatherer
}

// Cleanup makes cleanup.
func (X509Check) Cleanup() {}

func (x X509Check) createGatherer() (gatherer, error) {
	if x.Source == "" {
		return nil, errors.New("'source' parameter is mandatory, but it's not set")
	}

	u, err := url.Parse(x.Source)
	if err != nil {
		return nil, fmt.Errorf("error on parsing source : %v", err)
	}

	tlsCfg, err := web.NewTLSConfig(x.ClientTLSConfig)
	if err != nil {
		return nil, fmt.Errorf("error on creating tls config : %v", err)
	}
	if tlsCfg == nil {
		tlsCfg = &tls.Config{}
	}
	tlsCfg.ServerName = u.Hostname()

	switch u.Scheme {
	case "file":
		return cert.NewFile(u.Path), nil
	case "https", "udp", "udp4", "udp6", "tcp", "tcp4", "tcp6":
		if u.Scheme == "https" {
			u.Scheme = "tcp"
		}
		return cert.NewNet(u, tlsCfg, x.Timeout.Duration), nil
	case "smtp":
		u.Scheme = "tcp"
		return cert.NewSMTP(u, tlsCfg, x.Timeout.Duration), nil

	}
	return nil, fmt.Errorf("unsupported scheme in '%s'", u)
}

// Init makes initialization.
func (x *X509Check) Init() bool {
	g, err := x.createGatherer()
	if err != nil {
		x.Error(err)
		return false
	}

	x.gatherer = g
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
	mx, err := x.collect()
	if err != nil {
		x.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

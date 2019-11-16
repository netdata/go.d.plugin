package unbound

import (
	"crypto/tls"
	"strings"
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("unbound", creator)
}

func New() *Unbound {
	config := Config{
		// "/etc/unbound/unbound.conf"
		Address:  "192.168.88.223:8953",
		ConfPath: "/Users/ilyam/Projects/goland/go.d.plugin/modules/unbound/testdata/unbound.conf",
		Timeout:  web.Duration{Duration: time.Second * 2},
	}

	return &Unbound{
		Config:   config,
		charts:   charts.Copy(),
		curCache: newCollectCache(),
		cache:    newCollectCache(),
	}
}

type unboundClient interface {
	send(command string) ([]string, error)
}

type (
	Config struct {
		Address             string       `yaml:"address"`
		ConfPath            string       `yaml:"conf_path"`
		Timeout             web.Duration `yaml:"timeout"`
		DisableTLS          bool         `yaml:"disable_tls"`
		Cumulative          bool         `yaml:"cumulative_stats"`
		web.ClientTLSConfig `yaml:",inline"`
	}
	Unbound struct {
		module.Base
		Config `yaml:",inline"`

		client   unboundClient
		cache    collectCache
		curCache collectCache

		prevTotQueries float64
		hasExtCharts   bool

		charts *module.Charts
	}
)

func (Unbound) Cleanup() {}

func (u *Unbound) Init() bool {
	cfg, err := readConfig(u.ConfPath)
	if err != nil {
		u.Warning(err)
	}
	if cfg != nil {
		if cfg.isRemoteControlDisabled() {
			u.Info("remote control is disabled in the configuration file")
			return false
		}
		u.applyConfig(cfg)
	}

	cl, err := u.createClient()
	if err != nil {
		u.Errorf("creating client: %v", err)
		return false
	}
	u.client = cl
	return true
}

func (u *Unbound) Check() bool {
	return len(u.Collect()) > 0
}

func (u Unbound) Charts() *module.Charts {
	return u.charts
}

func (u *Unbound) Collect() map[string]int64 {
	mx, err := u.collect()
	if err != nil {
		u.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}

	return mx
}

func (u *Unbound) applyConfig(cfg *unboundConfig) {
	if cfg.hasServerSection() {
		if cfg.Server.StatisticsCumulative.isSet() {
			u.Cumulative = cfg.Server.StatisticsCumulative.bool()
		}
	}
	if !cfg.hasRemoteControlSection() {
		return
	}
	if cfg.RemoteControl.ControlUseCert.isSet() {
		u.DisableTLS = cfg.RemoteControl.ControlUseCert.bool()
	}
	if cfg.RemoteControl.ControlKeyFile.isSet() {
		u.TLSKey = string(cfg.RemoteControl.ControlKeyFile)
	}
	if cfg.RemoteControl.ControlCertFile.isSet() {
		u.TLSCert = string(cfg.RemoteControl.ControlCertFile)
	}
}

func (u *Unbound) createClient() (uClient unboundClient, err error) {
	useTLS := !isUnixSocket(u.Address) && !u.DisableTLS
	var tlsCfg *tls.Config
	if useTLS {
		if tlsCfg, err = web.NewTLSConfig(u.ClientTLSConfig); err != nil {
			return nil, err
		}
	}
	uClient = newClient(clientConfig{
		address: u.Address,
		timeout: u.Timeout.Duration,
		useTLS:  useTLS,
		tlsConf: tlsCfg,
	})
	return uClient, err
}

func isUnixSocket(address string) bool { return strings.HasPrefix(address, "/") }

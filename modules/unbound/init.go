package unbound

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/web"

	"gopkg.in/yaml.v2"
)

type (
	str     string
	strBool string
)

func (s str) isSet() bool     { return s != "" }
func (s strBool) isSet() bool { return s != "" }
func (s strBool) bool() bool  { return s == "yes" }

type unboundConfig struct {
	Server *struct {
		StatisticsCumulative strBool `yaml:"statistics-cumulative"`
	} `yaml:"server"`
	RemoteControl *struct {
		ControlEnable    strBool `yaml:"control-enable"`
		ControlInterface str     `yaml:"control-interface"`
		ControlPort      str     `yaml:"control-port"`
		ControlUseCert   strBool `yaml:"control-use-cert"`
		ControlKeyFile   str     `yaml:"control-key-file"`
		ControlCertFile  str     `yaml:"control-cert-file"`
	} `yaml:"remote-control"`
}

func (c unboundConfig) hasServerSection() bool {
	return c.Server != nil
}
func (c unboundConfig) hasRemoteControlSection() bool {
	return c.RemoteControl != nil
}
func (c unboundConfig) isRemoteControlDisabled() bool {
	return c.hasRemoteControlSection() && c.RemoteControl.ControlEnable.isSet() && !c.RemoteControl.ControlEnable.bool()
}

func (u *Unbound) initConfig() bool {
	if u.ConfPath == "" {
		return true
	}
	cfg, err := readConfig(u.ConfPath)
	if err != nil {
		u.Warning(err)
	}
	if cfg == nil {
		return true
	}
	if cfg.isRemoteControlDisabled() {
		u.Info("remote control is disabled in the configuration file")
		return false
	}
	u.applyConfig(cfg)
	return true
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
	if cfg.RemoteControl.ControlInterface.isSet() {
		u.Address = string(cfg.RemoteControl.ControlInterface)
	}
	if cfg.RemoteControl.ControlPort.isSet() && !isUnixSocket(u.Address) {
		host, _, _ := net.SplitHostPort(u.Address)
		u.Address = net.JoinHostPort(host, string(cfg.RemoteControl.ControlPort))
	}
}

func (u *Unbound) initClient() (err error) {
	var tlsCfg *tls.Config

	useTLS := !isUnixSocket(u.Address) && !u.DisableTLS
	if useTLS {
		if tlsCfg, err = web.NewTLSConfig(u.ClientTLSConfig); err != nil {
			return err
		}
	}

	u.client = newClient(clientConfig{
		address: u.Address,
		timeout: u.Timeout.Duration,
		useTLS:  useTLS,
		tlsConf: tlsCfg,
	})
	return nil
}

func isUnixSocket(address string) bool {
	return strings.HasPrefix(address, "/")
}

func readConfig(config string) (*unboundConfig, error) {
	b, err := ioutil.ReadFile(config)
	if err != nil {
		return nil, fmt.Errorf("error on reading config: %v", err)
	}

	b = adjustUnboundConfig(b)

	var cfg unboundConfig
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("error on parsing config: %v", err)
	}
	return &cfg, nil
}

func adjustUnboundConfig(cfg []byte) []byte {
	return bytes.ReplaceAll(cfg, []byte("\t"), []byte(" "))
}

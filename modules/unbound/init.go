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
	Srv *struct {
		Cumulative strBool `yaml:"statistics-cumulative"`
	} `yaml:"server"`
	RC *struct {
		Enable    strBool `yaml:"control-enable"`
		Interface str     `yaml:"control-interface"`
		Port      str     `yaml:"control-port"`
		UseCert   strBool `yaml:"control-use-cert"`
		KeyFile   str     `yaml:"control-key-file"`
		CertFile  str     `yaml:"control-cert-file"`
	} `yaml:"remote-control"`
}

func (c unboundConfig) hasServer() bool {
	return c.Srv != nil
}
func (c unboundConfig) hasRemoteControl() bool {
	return c.RC != nil
}
func (c unboundConfig) isRemoteControlDisabled() bool {
	return c.hasRemoteControl() && c.RC.Enable.isSet() && !c.RC.Enable.bool()
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
	if cfg.hasServer() && cfg.Srv.Cumulative.isSet() {
		u.Cumulative = cfg.Srv.Cumulative.bool()
	}
	if !cfg.hasRemoteControl() {
		return
	}
	if cfg.RC.UseCert.isSet() {
		u.DisableTLS = cfg.RC.UseCert.bool()
	}
	if cfg.RC.KeyFile.isSet() {
		u.TLSKey = string(cfg.RC.KeyFile)
	}
	if cfg.RC.CertFile.isSet() {
		u.TLSCert = string(cfg.RC.CertFile)
	}
	if cfg.RC.Interface.isSet() {
		u.Address = string(cfg.RC.Interface)
	}
	if cfg.RC.Port.isSet() && !isUnixSocket(u.Address) {
		host, _, _ := net.SplitHostPort(u.Address)
		port := string(cfg.RC.Port)
		u.Address = net.JoinHostPort(host, port)
	}
}

func (u *Unbound) initClient() (err error) {
	var tlsCfg *tls.Config

	useTLS := !isUnixSocket(u.Address) && !u.DisableTLS
	//if useTLS && (u.TLSCert == "" || u.TLSKey == "") {
	//	return errors.New("")
	//}
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

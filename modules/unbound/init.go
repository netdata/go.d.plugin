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
		Cumulative strBool `yaml:"statistics-cumulative,omitempty"`
	} `yaml:"server"`
	RC *struct {
		Enable    strBool `yaml:"control-enable,omitempty"`
		Interface str     `yaml:"control-interface,omitempty"`
		Port      str     `yaml:"control-port,omitempty"`
		UseCert   strBool `yaml:"control-use-cert,omitempty"`
		KeyFile   str     `yaml:"control-key-file,omitempty"`
		CertFile  str     `yaml:"control-cert-file,omitempty"`
	} `yaml:"remote-control"`
}

func (c unboundConfig) String() string { b, _ := yaml.Marshal(c); return string(b) }

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
	// TODO: config parameters auto detection by reading the config file feature is questionable
	// unbound config file is not in yaml format, it looks like yaml but it is not, for example it allows such config
	// remote-control:
	//   control-interface: 0.0.0.0
	//   control-interface: /var/run/unbound.sock
	// Module will try to get stats from /var/run/unbound.sock and fail. Unbound doesnt allow to query stats from
	// unix socket when control-interface enabled on ip interface
	if u.ConfPath == "" {
		u.Info("'conf_path' not set, skipping parameters auto detection")
		return true
	}
	u.Info("reading '%s'", u.ConfPath)
	cfg, err := readConfig(u.ConfPath)
	if err != nil {
		u.Warningf("%v, skipping parameters auto detection", err)
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
	u.Debugf("applying configuration:\n%s", cfg)
	if cfg.hasServer() && cfg.Srv.Cumulative.isSet() {
		u.Debugf("found 'statistics-cumulative', applying 'cumulative_stats': %v", cfg.Srv.Cumulative.bool())
		u.Cumulative = cfg.Srv.Cumulative.bool()
	}
	if !cfg.hasRemoteControl() {
		return
	}
	if cfg.RC.UseCert.isSet() {
		u.Debugf("found 'control-use-cert', applying 'disable_tls': %v", !cfg.RC.UseCert.bool())
		u.DisableTLS = !cfg.RC.UseCert.bool()
	}
	if cfg.RC.KeyFile.isSet() {
		u.Debugf("found 'control-key-file', applying 'tls_key': %s", cfg.RC.KeyFile)
		u.TLSKey = string(cfg.RC.KeyFile)
	}
	if cfg.RC.CertFile.isSet() {
		u.Debugf("found 'control-cert-file', applying 'tls_cert': %s", cfg.RC.CertFile)
		u.TLSCert = string(cfg.RC.CertFile)
	}
	if cfg.RC.Interface.isSet() {
		u.Debugf("found 'control-interface', applying 'address': %s", cfg.RC.CertFile)
		u.Address = string(cfg.RC.Interface)
	}
	if cfg.RC.Port.isSet() && !isUnixSocket(u.Address) {
		host, _, _ := net.SplitHostPort(u.Address)
		port := string(cfg.RC.Port)
		address := net.JoinHostPort(host, port)
		u.Debugf("found 'control-port', applying 'address': %s", cfg.RC.CertFile)
		u.Address = address
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
	// unbound config is not yaml syntax file, but the fix makes it readable at least
	return bytes.ReplaceAll(cfg, []byte("\t"), []byte(" "))
}

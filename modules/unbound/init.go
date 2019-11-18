package unbound

import (
	"bytes"
	"crypto/tls"
	"errors"
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
func (s str) value() string   { return string(s) }
func (s strBool) isSet() bool { return s != "" }
func (s strBool) value() bool { return s == "yes" }

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

func (c unboundConfig) String() string         { b, _ := yaml.Marshal(c); return string(b) }
func (c unboundConfig) hasServer() bool        { return c.Srv != nil }
func (c unboundConfig) hasRemoteControl() bool { return c.RC != nil }

func (c unboundConfig) isRemoteControlDisabled() bool {
	return c.hasRemoteControl() && c.RC.Enable.isSet() && !c.RC.Enable.value()
}

func (u *Unbound) initConfig() (enabled bool) {
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
	u.Infof("reading '%s'", u.ConfPath)
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
		if cfg.Srv.Cumulative.value() != u.Cumulative {
			u.Debugf("changing 'cumulative_stats': %v => %v", u.Cumulative, cfg.Srv.Cumulative.value())
			u.Cumulative = cfg.Srv.Cumulative.value()
		}
	}

	if !cfg.hasRemoteControl() {
		return
	}
	if cfg.RC.UseCert.isSet() {
		if cfg.RC.UseCert.value() != u.DisableTLS {
			u.Debugf("changing 'disable_tls': %v => %v", u.DisableTLS, !cfg.RC.UseCert.value())
			u.DisableTLS = !cfg.RC.UseCert.value()
		}
	}

	if cfg.RC.KeyFile.isSet() {
		if cfg.RC.KeyFile.value() != u.TLSKey {
			u.Debugf("changing 'tls_key': '%s' => '%s'", u.TLSKey, cfg.RC.KeyFile)
			u.TLSKey = cfg.RC.KeyFile.value()
		}
	}

	if cfg.RC.CertFile.isSet() {
		if cfg.RC.CertFile.value() != u.TLSCert {
			u.Debugf("changing 'tls_cert': '%s' => '%s'", u.TLSCert, cfg.RC.CertFile)
			u.TLSCert = cfg.RC.CertFile.value()
		}
	}

	if cfg.RC.Interface.isSet() {
		if v := adjustControlInterface(cfg.RC.Interface.value()); v != u.Address {
			u.Debugf("changing 'address': '%s' => '%s'", u.Address, v)
			u.Address = v
		}
	}

	if cfg.RC.Port.isSet() && !isUnixSocket(u.Address) {
		if host, port, err := net.SplitHostPort(u.Address); err == nil && port != cfg.RC.Port.value() {
			address := net.JoinHostPort(host, cfg.RC.Port.value())
			u.Debugf("changing 'address': '%s' => '%s'", u.Address, address)
			u.Address = address
		}
	}
}

func (u *Unbound) initClient() (err error) {
	var tlsCfg *tls.Config
	useTLS := !isUnixSocket(u.Address) && !u.DisableTLS
	if useTLS && (u.TLSCert == "" || u.TLSKey == "") {
		return errors.New("'tls_cert' or 'tls_key' is missing")
	}
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

func adjustControlInterface(value string) string {
	if isUnixSocket(value) {
		return value
	}
	if value == "0.0.0.0" {
		value = "127.0.0.1"
	}
	return net.JoinHostPort(value, "8953")
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
	// unbound config format is not yaml, but the fix makes it readable at least
	return bytes.ReplaceAll(cfg, []byte("\t"), []byte(" "))
}

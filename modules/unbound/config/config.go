package config

import "strings"

// UnboundConfig represents Unbound configuration file.
type UnboundConfig struct {
	Cumulative string // statistics-cumulative
	Enable     string // control-enable
	Interface  string // control-interface
	Port       string // control-port
	UseCert    string // control-use-cert
	KeyFile    string // control-key-file
	CertFile   string // control-cert-file
}

func fromOptions(options []option) *UnboundConfig {
	cfg := &UnboundConfig{}
	for _, opt := range options {
		switch opt.name {
		default:
		case "control-interface":
			applyControlInterface(cfg, opt.value)
		case "statistics-cumulative":
			cfg.Cumulative = opt.value
		case "control-enable":
			cfg.Enable = opt.value
		case "control-port":
			cfg.Port = opt.value
		case "control-use-key-file":
			cfg.KeyFile = opt.value
		case "control-use-cert-file":
			cfg.CertFile = opt.value
		}
	}
	return cfg
}

func applyControlInterface(cfg *UnboundConfig, value string) {
	if cfg.Interface == "" {
		cfg.Interface = value
		return
	}
	if !isUnixSocket(value) || isUnixSocket(cfg.Interface) {
		cfg.Interface = value
	}
}

func isUnixSocket(address string) bool {
	return strings.HasPrefix(address, "/")
}

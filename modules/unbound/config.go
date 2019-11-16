package unbound

import (
	"bytes"
	"fmt"
	"io/ioutil"

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
	return c.hasRemoteControlSection() && c.RemoteControl.ControlEnable.isSet() && c.RemoteControl.ControlEnable.bool()
}

func adjustUnboundConfig(cfg []byte) []byte {
	return bytes.ReplaceAll(cfg, []byte("\t"), []byte(" "))
}

func readConfig(config string) (*unboundConfig, error) {
	if config == "" {
		return nil, nil
	}
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

package snmp

import (
	"fmt"
	"time"

	gosnmp "github.com/gosnmp/gosnmp"
	"github.com/netdata/go.d.plugin/agent/module"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("snmp", creator)
}

//Everything is initialized at Init()
func New() *SNMP {
	return &SNMP{}
}

type (
	Config struct {
		SNMPClient  gosnmp.GoSNMP
		Name        string         `default:"127.0.0.1" yaml:"hostname"`
		UpdateEvery int            `default:3 yaml:"update_every"`
		Community   *string        `yaml:"community,omitempty"`
		User        *User          `yaml:"user,omitempty"`
		Options     *Options       `yaml:"options,omitempty"`
		Settings    []ChartsConfig `yaml:"charts,omitempty"`
	}
	User struct {
		Name      string `yaml:"name"`
		Level     int    `default:1 yaml:"level"`
		AuthProto int    `default:1 yaml:"auth_proto"`
		AuthKey   string `yaml:"auth_key"`
		PrivProto int    `default:1 yaml:"priv_proto"`
		PrivKey   string `yaml:"priv_key"`
	}
	Options struct {
		Port    int `default:161 yaml:"port"`
		Retries int `default:1 yaml:"retries"`
		Timeout int `default:2 yaml:"timeout"`
		Version int `default:1 yaml:"version"`
	}
	ChartsConfig struct {
		Title         string      `yaml:"title"`
		Priority      int         `default:7000 yaml:"priority"`
		Units         *string     `yaml:"units,omitempty"`
		Type          *string     `yaml:"type,omitempty"`
		Family        *string     `yaml:"family,omitempty"`
		MultiplyRange [2]int      `yaml:"multiply_range,omitempty"`
		Dimensions    []Dimension `yaml:"dimensions,omitempty"`
	}
	Dimension struct {
		Name       string `yaml:"name"`
		OID        string `yaml:"oid"`
		Algorithm  string `yaml:"algorithm"`
		Multiplier int    `yaml:"multiplier"`
		Divisor    int    `yaml:"divisor"`
	}
)

type SNMP struct {
	module.Base
	Config `yaml:",inline"`
	charts *module.Charts
}

func (s *SNMP) Init() bool {
	err := s.validateConfig()
	if err != nil {
		s.Errorf("config validation: %v", err)
		return false
	}

	//Default SNMP connection params
	params := &gosnmp.GoSNMP{
		Target:    s.Name,
		Port:      uint16(s.Options.Port),
		Community: *s.Community,
		Version:   gosnmp.Version2c,
		Timeout:   time.Duration(2) * time.Second,
		Logger:    gosnmp.NewLogger(s.Logger),
	}
	switch s.Options.Version {
	case 1:
		version := gosnmp.Version1
		params.Version = version
		params.Timeout = time.Duration(s.Options.Timeout) * time.Second

	case 2:
		version := gosnmp.Version2c
		params.Version = version
		params.Timeout = time.Duration(s.Options.Timeout) * time.Second

	case 3:
		version := gosnmp.Version3
		params.Version = version
		params.Timeout = time.Duration(s.Options.Timeout) * time.Second
		params.SecurityModel = gosnmp.SnmpV3SecurityModel(s.User.Level)
		params.MsgFlags = gosnmp.SnmpV3MsgFlags(s.User.AuthProto)
		params.SecurityParameters = &gosnmp.UsmSecurityParameters{
			AuthenticationProtocol:   gosnmp.SnmpV3AuthProtocol(s.User.AuthProto),
			AuthenticationPassphrase: s.User.AuthKey,
			PrivacyProtocol:          gosnmp.SnmpV3PrivProtocol(s.User.PrivProto),
			PrivacyPassphrase:        s.User.PrivKey,
		}

	default:
		s.Errorf("invalid SNMP version: %d", s.Options.Version)
		return false
	}

	err = params.Connect()
	if err != nil {
		s.Errorf("SNMP Connect fail: %v", err)
		return false
	}
	s.Config.SNMPClient = *params

	if len(s.Settings) > 0 {
		s.charts = newChart(s.Settings)
	} else {
		c := snmp_chart_template.Copy()
		c.ID = fmt.Sprintf(c.ID, 1)
		c.Title = fmt.Sprint(c.Title, "default")
		c.AddDim(default_dims[0])
		c.AddDim(default_dims[1])
		s.charts = &module.Charts{c}
	}
	return true
}

func (s *SNMP) Check() bool {
	return len(s.Collect()) > 0
}

func (s *SNMP) Charts() *module.Charts {
	return s.charts
}

func (s *SNMP) Collect() map[string]int64 {
	mx, err := s.collect()
	if err != nil {
		s.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (s SNMP) Cleanup() {
	params := s.Config.SNMPClient
	params.Conn.Close()
}

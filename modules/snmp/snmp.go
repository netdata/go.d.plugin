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
	comm := "public"
	return &SNMP{
		Config: Config{
			Name:        "127.0.0.1",
			MaxOIDs:     60,
			Community:   &comm,
			UpdateEvery: 3,
			Options: &Options{
				Port:    161,
				Retries: 1,
				Timeout: 2,
				Version: 2,
			},
		},
	}
}

type (
	Config struct {
		SNMPClient  *gosnmp.GoSNMP
		Name        string         `yaml:"hostname"`
		MaxOIDs     int            `yaml:"max_request_size"`
		UpdateEvery int            `yaml:"update_every"`
		Community   *string        `yaml:"community,omitempty"`
		User        *User          `yaml:"user,omitempty"`
		Options     *Options       `yaml:"options,omitempty"`
		ChartInput  []ChartsConfig `yaml:"charts,omitempty"`
	}
	User struct {
		Name      string `yaml:"name"`
		Level     int    `yaml:"level"`
		AuthProto int    `yaml:"auth_proto"`
		AuthKey   string `yaml:"auth_key"`
		PrivProto int    `yaml:"priv_proto"`
		PrivKey   string `yaml:"priv_key"`
	}
	Options struct {
		Port    int `yaml:"port"`
		Retries int `yaml:"retries"`
		Timeout int `yaml:"timeout"`
		Version int `yaml:"version"`
	}
	ChartsConfig struct {
		Title         string      `yaml:"title"`
		Priority      int         `yaml:"priority"`
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
		MaxOids:   s.MaxOIDs,
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
	s.Config.SNMPClient = params
	if len(s.ChartInput) > 0 {
		s.charts = newChart(s.ChartInput)
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
	if s.Config.SNMPClient != nil {
		params := s.Config.SNMPClient
		params.Conn.Close()
	}
}

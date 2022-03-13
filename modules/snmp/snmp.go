package snmp

import (
	"fmt"
	"time"

	gosnmp "github.com/gosnmp/gosnmp"
	"github.com/netdata/go.d.plugin/agent/module"
)

var snmpHandler = gosnmp.NewHandler

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("snmp", creator)
}

func New() *SNMP {
	comm := "public"
	return &SNMP{
		Config: Config{
			Name:        "127.0.0.1",
			Community:   &comm,
			UpdateEvery: 3,
			Options: &Options{
				Port:    161,
				Retries: 1,
				Timeout: 2,
				Version: 2,
				MaxOIDs: 60,
			},
		},
	}
}

type (
	Config struct {
		SNMPClient  gosnmp.Handler
		Name        string         `yaml:"hostname"`
		UpdateEvery int            `yaml:"update_every"`
		Community   *string        `yaml:"community"`
		User        *User          `yaml:"user"`
		Options     *Options       `yaml:"options"`
		ChartInput  []ChartsConfig `yaml:"charts"`
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
		MaxOIDs int `yaml:"max_request_size"`
	}
	ChartsConfig struct {
		Title         string      `yaml:"title"`
		Priority      int         `yaml:"priority"`
		Units         *string     `yaml:"units"`
		Type          *string     `yaml:"type"`
		Family        *string     `yaml:"family"`
		MultiplyRange []int       `yaml:"multiply_range"`
		Dimensions    []Dimension `yaml:"dimensions"`
	}
	Dimension struct {
		Name       string  `yaml:"name"`
		OID        string  `yaml:"oid"`
		Algorithm  *string `yaml:"algorithm"`
		Multiplier *int    `yaml:"multiplier"`
		Divisor    *int    `yaml:"divisor"`
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

	snmpClient := snmpHandler()

	//Default SNMP connection params
	snmpClient.SetTarget(s.Name)
	snmpClient.SetPort(uint16(s.Options.Port))
	snmpClient.SetMaxOids(s.Options.MaxOIDs)
	snmpClient.SetLogger(gosnmp.NewLogger(s.Logger))
	snmpClient.SetTimeout(time.Duration(s.Options.Timeout) * time.Second)

	switch s.Options.Version {
	case 1:
		snmpClient.SetCommunity(*s.Community)
		snmpClient.SetVersion(gosnmp.Version1)

	case 2:
		snmpClient.SetCommunity(*s.Community)
		snmpClient.SetVersion(gosnmp.Version2c)

	case 3:
		snmpClient.SetVersion(gosnmp.Version3)
		snmpClient.SetSecurityModel(gosnmp.UserSecurityModel)
		snmpClient.SetMsgFlags(gosnmp.SnmpV3MsgFlags(s.User.Level))
		snmpClient.SetSecurityParameters(&gosnmp.UsmSecurityParameters{
			UserName:                 s.User.Name,
			AuthenticationProtocol:   gosnmp.SnmpV3AuthProtocol(s.User.AuthProto),
			AuthenticationPassphrase: s.User.AuthKey,
			PrivacyProtocol:          gosnmp.SnmpV3PrivProtocol(s.User.PrivProto),
			PrivacyPassphrase:        s.User.PrivKey,
		})

	default:
		s.Errorf("invalid SNMP version: %d", s.Options.Version)
		return false
	}

	err = snmpClient.Connect()
	if err != nil {
		s.Errorf("SNMP Connect fail: %v", err)
		return false
	}
	s.SNMPClient = snmpClient

	if len(s.ChartInput) > 0 {
		s.charts, err = allCharts(s.ChartInput)
		if err != nil {
			s.Errorf("Population of charts failed: %v", err)
			return false
		}
	} else {
		c := defaultSNMPchart.Copy()
		c.ID = fmt.Sprintf(c.ID, 1)
		c.Title = fmt.Sprint(c.Title, "default")
		for _, d := range defaultDims {
			if err = c.AddDim(d); err != nil {
				s.Errorf("Population of charts failed: %v", err)
				return false
			}
		}
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

func (s *SNMP) Cleanup() {
	if s.SNMPClient != nil {
		s.SNMPClient.Close()
	}
}

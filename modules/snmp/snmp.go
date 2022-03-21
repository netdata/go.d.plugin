package snmp

import (
	"github.com/gosnmp/gosnmp"
	"github.com/netdata/go.d.plugin/agent/module"
)

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
			Hostname:    "127.0.0.1",
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
		Hostname    string         `yaml:"hostname"`
		UpdateEvery int            `yaml:"update_every"`
		Community   *string        `yaml:"community"`
		User        *User          `yaml:"user"`
		Options     *Options       `yaml:"options"`
		ChartInput  []ChartsConfig `yaml:"charts"`
	}
	User struct {
		Name      string `yaml:"name"`
		Level     string `yaml:"level"`
		AuthProto string `yaml:"auth_proto"`
		AuthKey   string `yaml:"auth_key"`
		PrivProto string `yaml:"priv_proto"`
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
		ID            string      `yaml:"id"`
		Title         string      `yaml:"title"`
		Units         *string     `yaml:"units"`
		Type          *string     `yaml:"type"`
		Family        *string     `yaml:"family"`
		Priority      int         `yaml:"priority"`
		MultiplyRange []int       `yaml:"multiply_range"`
		Dimensions    []Dimension `yaml:"dimensions"`
	}
	Dimension struct {
		OID        string  `yaml:"oid"`
		Name       string  `yaml:"name"`
		Algorithm  *string `yaml:"algorithm"`
		Multiplier *int    `yaml:"multiplier"`
		Divisor    *int    `yaml:"divisor"`
	}
)

type SNMP struct {
	module.Base
	snmpClient gosnmp.Handler
	Config     `yaml:",inline"`
	charts     *module.Charts
}

func (s *SNMP) Init() bool {
	err := s.validateConfig()
	if err != nil {
		s.Errorf("config validation: %v", err)
		return false
	}

	snmpClient, err := s.initSNMPClient()
	if err != nil {
		s.Errorf("SNMP client initialization: %v", err)
		return false
	}

	err = snmpClient.Connect()
	if err != nil {
		s.Errorf("SNMP client connect: %v", err)
		return false
	}
	s.snmpClient = snmpClient

	s.charts, err = newCharts(s.ChartInput)
	if err != nil {
		s.Errorf("Population of charts failed: %v", err)
		return false
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
	if s.snmpClient != nil {
		s.snmpClient.Close()
	}
}

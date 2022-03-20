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
	snmpHandler gosnmp.Handler
	Config      `yaml:",inline"`
	charts      *module.Charts
}

func (s *SNMP) Init() bool {
	err := s.validateConfig()
	if err != nil {
		s.Errorf("config validation: %v", err)
		return false
	}

	s.snmpHandler, err = s.initSNMPClient()
	if err != nil {
		s.Errorf("SNMP Connect fail: %v", err)
		return false
	}

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
	if s.snmpHandler != nil {
		s.snmpHandler.Close()
	}
}

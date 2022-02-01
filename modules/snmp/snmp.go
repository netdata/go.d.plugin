package snmp

import (
	"math/rand"
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

const (
	defaultTarget    = "127.0.0.1"
	defaultPort      = 161
	defaultCommunity = "public"
)

func New() *SNMP {
	c := snmp_chart_template.Copy()
	sCharts := &module.Charts{}
	sCharts.Add(c)

	return &SNMP{
		randInt: func() int64 { return rand.Int63n(100) },
		charts:  sCharts,
	}
}

type (
	Config struct {
		SNMPClient gosnmp.GoSNMP
		Name       string         `yaml:"hostname"`
		Port       int            `yaml:"port"`
		Community  string         `yaml:"community"`
		Settings   []ChartsConfig `yaml:"charts"`
	}
	ChartsConfig struct {
		Title    string `yaml:"title"`
		Priority int    `yaml:"priority"`
	}
)

type SNMP struct {
	module.Base
	Config  `yaml:",inline"`
	randInt func() int64
	charts  *module.Charts
}

func (s *SNMP) Init() bool {
	err := s.validateConfig()
	if err != nil {
		s.Errorf("config validation: %v", err)
		return false
	}
	if len(s.Settings) > 0 {
		s.charts = newChart(s.Settings)
	}

	params := &gosnmp.GoSNMP{
		Target:    s.Name,
		Port:      uint16(s.Port),
		Community: s.Community,
		Version:   gosnmp.Version2c,
		Timeout:   time.Duration(2) * time.Second,
		Logger:    gosnmp.NewLogger(s.Logger),
	}

	err = params.Connect()
	if err != nil {
		s.Errorf("SNMP Connect fail: %v", err)
		return false
	}
	s.Config.SNMPClient = *params
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

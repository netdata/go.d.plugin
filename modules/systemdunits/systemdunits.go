// +build linux

package systemdunits

import (
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/matcher"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	module.Register("systemdunits", module.Creator{
		Defaults: module.Defaults{
			Disabled:    true,
			UpdateEvery: 10, // gathering systemd units is CPU-intensive op
		},
		Create: func() module.Module { return New() },
	})
}

func New() *SystemdUnits {
	return &SystemdUnits{
		Config: Config{
			Include: []string{
				"*.service",
			},
			Timeout: web.Duration{Duration: time.Second * 2},
		},

		client:         newSystemdDBusClient(),
		collectedUnits: make(map[string]bool),
	}
}

type Config struct {
	Include []string     `yaml:"include"`
	Timeout web.Duration `yaml:"timeout"`
}

type SystemdUnits struct {
	module.Base
	Config `yaml:",inline"`

	client systemdClient
	conn   systemdConnection

	systemdVersion int
	collectedUnits map[string]bool
	sr             matcher.Matcher

	charts *module.Charts
}

func (s *SystemdUnits) Init() bool {
	err := s.validateConfig()
	if err != nil {
		s.Errorf("config validation: %v", err)
		return false
	}

	sr, err := s.initSelector()
	if err != nil {
		s.Errorf("init selector: %v", err)
		return false
	}
	s.sr = sr

	cs, err := s.initCharts()
	if err != nil {
		s.Errorf("init charts: %v", err)
		return false
	}
	s.charts = cs

	s.Debugf("unit names patterns: %v", s.Include)
	s.Debugf("timeout: %s", s.Timeout)
	return true
}

func (s *SystemdUnits) Check() bool {
	return len(s.Collect()) > 0
}

func (s *SystemdUnits) Charts() *module.Charts {
	return s.charts
}

func (s *SystemdUnits) Collect() map[string]int64 {
	ms, err := s.collect()
	if err != nil {
		s.Error(err)
	}

	if len(ms) == 0 {
		return nil
	}
	return ms
}

func (s *SystemdUnits) Cleanup() {
	s.closeConnection()
}

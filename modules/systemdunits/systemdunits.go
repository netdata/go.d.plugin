// +build linux

package systemdunits

import (
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
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
		charts:         charts.Copy(),
	}
}

type Config struct {
	Include []string     `yaml:"include"`
	Timeout web.Duration `yaml:"timeout"`
}

type (
	SystemdUnits struct {
		module.Base
		Config `yaml:",inline"`

		client systemdClient
		conn   systemdConnection

		collectedUnits map[string]bool
		charts         *module.Charts
	}
)

func (s *SystemdUnits) Init() bool {
	if len(s.Include) == 0 {
		s.Error("'include' option not set")
		return false
	}
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
	mx, err := s.collect()
	if err != nil {
		s.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (s *SystemdUnits) Cleanup() {
	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
	}
}

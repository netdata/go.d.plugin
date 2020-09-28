package systemdunits

import (
	"github.com/netdata/go.d.plugin/agent/module"

	"github.com/coreos/go-systemd/v22/dbus"
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
			UnitsPatterns: []string{
				"*.service",
			},
		},

		client:         newSystemdDBusClient(),
		collectedUnits: make(map[string]bool),
		charts:         charts.Copy(),
	}
}

type Config struct {
	UnitsPatterns []string `yaml:"units_patterns"`
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
	systemdClient interface {
		connect() (systemdConnection, error)
	}
	systemdConnection interface {
		Close()
		ListUnitsByPatterns(states []string, patterns []string) ([]dbus.UnitStatus, error)
	}
)

func (s *SystemdUnits) Init() bool {
	if len(s.UnitsPatterns) == 0 {
		s.Error("'unit_patterns' not set")
		return false
	}
	return true
}

func (s *SystemdUnits) Check() bool {
	return len(s.Collect()) > 0
}

func (s *SystemdUnits) Charts() *Charts {
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

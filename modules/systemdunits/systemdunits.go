package systemdunits

import (
	"github.com/netdata/go-orchestrator/module"
	"github.com/netdata/go.d.plugin/pkg/matcher"
)

type Config struct {
	Selector matcher.SimpleExpr `yaml:"selector"`
}

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			Disabled:    true,
			UpdateEvery: 1,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("systemdunits", creator)
}

// New creates SystemdUnits with default values
func New() *SystemdUnits {
	return &SystemdUnits{
		charts: charts.Copy(),
	}
}

// SystemdUnits systemdunits module
type SystemdUnits struct {
	module.Base // should be embedded by every module
	Config      `yaml:",inline"`
	charts      *module.Charts
	selector    matcher.Matcher
}

// Cleanup makes cleanup
func (SystemdUnits) Cleanup() {}

// Init makes initialization
func (s *SystemdUnits) Init() bool {
	if !s.Selector.Empty() {
		m, err := s.Selector.Parse()
		if err != nil {
			s.Errorf("error on creating per user stats matcher : %v", err)
		}
		s.selector = matcher.WithCache(m)
	}

	return true
}

// Check makes check
func (s SystemdUnits) Check() bool {
	return len(s.Collect()) > 0
}

// Charts creates Charts
func (s *SystemdUnits) Charts() *Charts {
	return s.charts
}

// Collect collects metrics
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

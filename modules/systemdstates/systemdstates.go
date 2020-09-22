package systemdstates

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
			Disabled:    false,
			UpdateEvery: 1,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("systemdstates", creator)
}

// New creates SystemdStates with default values
func New() *SystemdStates {
	return &SystemdStates{
		charts: charts.Copy(),
	}
}

// SystemdStates SystemdStates module
type SystemdStates struct {
	module.Base // should be embedded by every module
	Config      `yaml:",inline"`
	charts      *module.Charts
	selector    matcher.Matcher
}

// Cleanup makes cleanup
func (SystemdStates) Cleanup() {}

// Init makes initialization
func (s *SystemdStates) Init() bool {

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
func (s SystemdStates) Check() bool {
	return len(s.Collect()) > 0
}

// Charts creates Charts
func (s *SystemdStates) Charts() *Charts {
	return s.charts
}

// Collect collects metrics
func (s *SystemdStates) Collect() map[string]int64 {
	mx, err := s.collect()
	if err != nil {
		s.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

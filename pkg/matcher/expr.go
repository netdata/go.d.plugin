package matcher

import "fmt"

type (
	// SimpleExpr is a simple expression to describe the condition:
	//     (includes[0].Match(v) || includes[1].Match(v) || ...) && !(excludes[0].Match(v) || excludes[1].Match(v) || ...)
	SimpleExpr struct {
		Includes []string `yaml:"includes" json:"includes"`
		Excludes []string `yaml:"excludes" json:"excludes"`

		includes []Matcher
		excludes []Matcher
	}
)

// Parse parse the given matchers in Includes and Excludes
func (s *SimpleExpr) Parse() error {
	if len(s.Includes) > 0 {
		s.includes = make([]Matcher, 0, len(s.Includes))
		for _, item := range s.Includes {
			m, err := Parse(item)
			if err != nil {
				return fmt.Errorf("parse matcher %q error: %v", item, err)
			}
			s.includes = append(s.includes, m)
		}
	}
	if len(s.Excludes) > 0 {
		s.excludes = make([]Matcher, 0, len(s.Excludes))
		for _, item := range s.Excludes {
			m, err := Parse(item)
			if err != nil {
				return fmt.Errorf("parse matcher %q error: %v", item, err)
			}
			s.excludes = append(s.excludes, m)
		}
	}
	return nil
}

// Match match against []byte
func (s *SimpleExpr) Match(b []byte) bool {
	if len(s.includes) > 0 {
		for _, m := range s.includes {
			if !m.Match(b) {
				return false
			}
		}
	}
	if len(s.excludes) > 0 {
		for _, m := range s.excludes {
			if m.Match(b) {
				return false
			}
		}
	}
	return true
}

// MatchString match against string
func (s *SimpleExpr) MatchString(str string) bool {
	if len(s.includes) > 0 {
		for _, m := range s.includes {
			if !m.MatchString(str) {
				return false
			}
		}
	}
	if len(s.excludes) > 0 {
		for _, m := range s.excludes {
			if m.MatchString(str) {
				return false
			}
		}
	}
	return true
}

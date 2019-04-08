package matcher

import "fmt"

type (
	// SimpleExpr is a simple expression to describe the condition:
	//     (includes[0].Match(v) || includes[1].Match(v) || ...) && !(excludes[0].Match(v) || excludes[1].Match(v) || ...)
	SimpleExpr struct {
		Includes []string `yaml:"includes" json:"includes"`
		Excludes []string `yaml:"excludes" json:"excludes"`

		matcher Matcher
	}
)

// Parse parse the given matchers in Includes and Excludes
func (s *SimpleExpr) Parse() error {
	var (
		includes = FALSE()
		excludes = FALSE()
	)
	for _, item := range s.Includes {
		m, err := Parse(item)
		if err != nil {
			return fmt.Errorf("parse matcher %q error: %v", item, err)
		}
		includes = Or(includes, m)
	}
	for _, item := range s.Excludes {
		m, err := Parse(item)
		if err != nil {
			return fmt.Errorf("parse matcher %q error: %v", item, err)
		}
		excludes = Or(excludes, m)
	}

	if len(s.Includes) == 0 {
		includes = TRUE()
	}
	s.matcher = And(includes, Not(excludes))
	return nil
}

// Match match against []byte
func (s *SimpleExpr) Match(b []byte) bool {
	return s.matcher.Match(b)
}

// MatchString match against string
func (s *SimpleExpr) MatchString(str string) bool {
	return s.matcher.MatchString(str)
}

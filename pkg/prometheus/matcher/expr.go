package matcher

import "fmt"

type Expr struct {
	Includes []string `yaml:"includes"`
	Excludes []string `yaml:"excludes"`
}

func (e Expr) Empty() bool {
	return len(e.Includes) == 0 && len(e.Excludes) == 0

}

func (e Expr) Parse() (Matcher, error) {
	if e.Empty() {
		return nil, nil
	}

	var matchers []Matcher
	var includes Matcher
	var excludes Matcher

	for _, item := range e.Includes {
		m, err := Parse(item)
		if err != nil {
			return nil, fmt.Errorf("parse matcher '%s': %v", item, err)
		}
		matchers = append(matchers, m)
	}

	switch len(matchers) {
	case 0:
		includes = trueMatcher{}
	case 1:
		includes = matchers[0]
	default:
		includes = or(matchers[0], matchers[1], matchers[2:]...)
	}

	matchers = matchers[:0]
	for _, item := range e.Excludes {
		m, err := Parse(item)
		if err != nil {
			return nil, fmt.Errorf("parse matcher '%s': %v", item, err)
		}
		matchers = append(matchers, m)
	}

	switch len(matchers) {
	case 0:
		excludes = falseMatcher{}
	case 1:
		excludes = matchers[0]
	default:
		excludes = or(matchers[0], matchers[1], matchers[2:]...)
	}

	return and(includes, not(excludes)), nil
}

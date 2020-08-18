package matcher

import "fmt"

type Expr struct {
	Allow []string `yaml:"allow"`
	Deny  []string `yaml:"deny"`
}

func (e Expr) Empty() bool {
	return len(e.Allow) == 0 && len(e.Deny) == 0

}

func (e Expr) Parse() (Matcher, error) {
	if e.Empty() {
		return nil, nil
	}

	var matchers []Matcher
	var allow Matcher
	var deny Matcher

	for _, item := range e.Allow {
		m, err := Parse(item)
		if err != nil {
			return nil, fmt.Errorf("parse matcher '%s': %v", item, err)
		}
		matchers = append(matchers, m)
	}

	switch len(matchers) {
	case 0:
		allow = trueMatcher{}
	case 1:
		allow = matchers[0]
	default:
		allow = Or(matchers[0], matchers[1], matchers[2:]...)
	}

	matchers = matchers[:0]
	for _, item := range e.Deny {
		m, err := Parse(item)
		if err != nil {
			return nil, fmt.Errorf("parse matcher '%s': %v", item, err)
		}
		matchers = append(matchers, m)
	}

	switch len(matchers) {
	case 0:
		deny = falseMatcher{}
	case 1:
		deny = matchers[0]
	default:
		deny = Or(matchers[0], matchers[1], matchers[2:]...)
	}

	return And(allow, Not(deny)), nil
}

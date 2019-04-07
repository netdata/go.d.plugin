package coredns

import (
	"fmt"

	"github.com/netdata/go.d.plugin/pkg/matcher"
)

type filter struct {
	Include []string `yaml:"include"`
	Exclude []string `yaml:"exclude"`
}

func (f filter) isEmpty() bool {
	return len(f.Include) == 0 && len(f.Exclude) == 0
}

func (f filter) createMatcher() (matcher.Matcher, error) {
	var (
		includes []matcher.Matcher
		excludes []matcher.Matcher
		include  = matcher.TRUE()
		exclude  = matcher.FALSE()
	)

	for _, line := range f.Include {
		m, err := matcher.Parse(line)
		if err != nil {
			return nil, fmt.Errorf("error on parsing line '%s' : %v", line, err)
		}
		includes = append(includes, m)
	}
	for _, line := range f.Exclude {
		m, err := matcher.Parse(line)
		if err != nil {
			return nil, fmt.Errorf("error on parsing line '%s' : %v", line, err)
		}
		excludes = append(excludes, m)
	}

	switch len(includes) {
	default:
		include = matcher.Or(includes[0], includes[1], includes[2:]...)
	case 0:
	case 1:
		include = includes[0]
	}

	switch len(excludes) {
	default:
		exclude = matcher.Or(excludes[0], excludes[1], excludes[2:]...)
	case 0:
	case 1:
		exclude = excludes[0]
	}

	return matcher.And(include, matcher.Not(exclude)), nil
}

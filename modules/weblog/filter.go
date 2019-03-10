package weblog

import (
	"fmt"

	"github.com/netdata/go.d.plugin/pkg/matcher"
)

type rawFilter struct {
	Include string `yaml:"include"`
	Exclude string `yaml:"exclude"`
}

func (r rawFilter) String() string {
	return fmt.Sprintf(`{"include": %q, "exclude": %q}`, r.Include, r.Exclude)
}

func NewFilter(raw rawFilter) (matcher.Matcher, error) {
	var (
		include matcher.Matcher
		exclude matcher.Matcher
		err     error
	)

	if raw.Include == "" && raw.Exclude == "" {
		return matcher.TRUE(), nil
	}

	if raw.Include != "" {
		if include, err = matcher.Parse(raw.Include); err != nil {
			return nil, err
		}
	}

	if raw.Exclude != "" {
		if exclude, err = matcher.Parse(raw.Exclude); err != nil {
			return nil, err
		}
		exclude = matcher.Not(exclude)
	}

	if include == nil {
		return exclude, nil
	}
	if exclude == nil {
		return include, nil
	}

	return matcher.And(include, exclude), nil
}

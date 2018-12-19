package weblog

import "fmt"

type rawFilter struct {
	Include string
	Exclude string
}

func (r rawFilter) String() string {
	return fmt.Sprintf("{include: %s, exclude: %s}", r.Include, r.Exclude)
}

type filter struct {
	include matcher
	exclude matcher
}

func (f *filter) match(s string) bool {
	includeOK := true
	excludeOK := false

	if f.include != nil {
		includeOK = f.include.match(s)
	}

	if f.exclude != nil {
		excludeOK = f.exclude.match(s)
	}

	return includeOK && !excludeOK
}

func newFilter(raw rawFilter) (matcher, error) {
	var f filter
	if raw.Include == "" && raw.Exclude == "" {
		return &f, nil
	}

	var err error

	if raw.Include != "" {
		if f.include, err = newMatcher(raw.Include); err != nil {
			return nil, err
		}
	}

	if raw.Exclude != "" {
		if f.exclude, err = newMatcher(raw.Exclude); err != nil {
			return nil, err
		}
	}

	return &f, nil
}

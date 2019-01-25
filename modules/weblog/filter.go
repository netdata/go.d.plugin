package weblog

import (
	"fmt"

	"github.com/netdata/go.d.plugin/pkg/matcher"
)

type rawfilter struct {
	Include string
	Exclude string
}

func (r rawfilter) String() string {
	return fmt.Sprintf("{include: %s, exclude: %s}", r.Include, r.Exclude)
}

type filter struct {
	include matcher.Matcher
	exclude matcher.Matcher
}

func (f *filter) Match(b []byte) bool {
	includeOK := true
	excludeOK := false

	if f.include != nil {
		includeOK = f.include.Match(b)
	}

	if f.exclude != nil {
		excludeOK = f.exclude.Match(b)
	}

	return includeOK && !excludeOK
}

func (f *filter) MatchString(s string) bool {
	includeOK := true
	excludeOK := false

	if f.include != nil {
		includeOK = f.include.MatchString(s)
	}

	if f.exclude != nil {
		excludeOK = f.exclude.MatchString(s)
	}

	return includeOK && !excludeOK
}

func newFilter(raw rawfilter) (matcher.Matcher, error) {
	var f filter
	if raw.Include == "" && raw.Exclude == "" {
		return &f, nil
	}

	var err error

	if raw.Include != "" {
		if f.include, err = matcher.Parse(raw.Include); err != nil {
			return nil, err
		}
	}

	if raw.Exclude != "" {
		if f.exclude, err = matcher.Parse(raw.Exclude); err != nil {
			return nil, err
		}
	}

	return &f, nil
}

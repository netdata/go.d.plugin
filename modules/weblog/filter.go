package weblog

import "github.com/netdata/go.d.plugin/modules/weblog/matcher"

type RawFilter struct {
	Include string
	Exclude string
}

type Filter struct {
	include matcher.Matcher
	exclude matcher.Matcher
}

func (f *Filter) Match(s string) bool {
	includeOK := true
	excludeOK := true

	if f.include != nil {
		includeOK = f.include.Match(s)
	}

	if f.exclude != nil {
		excludeOK = !f.exclude.Match(s)
	}

	return includeOK && excludeOK
}

func createFilter(raw RawFilter) (matcher.Matcher, error) {
	var f Filter

	if raw.Include == "" && raw.Exclude == "" {
		return &f, nil
	}

	if raw.Include != "" {
		m, err := matcher.New(raw.Include)
		if err != nil {
			return nil, err
		}
		f.include = m

	}

	if raw.Exclude != "" {
		m, err := matcher.New(raw.Exclude)
		if err != nil {
			return nil, err
		}
		f.exclude = m
	}

	return &f, nil

}

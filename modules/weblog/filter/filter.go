package filter

import "github.com/netdata/go.d.plugin/modules/weblog/matcher"

type Filter interface {
	Filter(line string) bool
}

type Raw struct {
	Include string
	Exclude string
}

type filter struct {
	include matcher.Matcher
	exclude matcher.Matcher
}

func (f *filter) Filter(s string) bool {
	includeOK := true
	excludeOK := true

	if f.include != nil {
		includeOK = f.include.Match(s)
	}

	if f.exclude != nil {
		excludeOK = f.exclude.Match(s)
	}

	return includeOK && !excludeOK
}

func New(raw Raw) (Filter, error) {
	var f filter

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

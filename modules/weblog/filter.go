package weblog

type rawFilter struct {
	Include string
	Exclude string

	include matcher
	exclude matcher
}

type filter struct {
	include matcher
	exclude matcher
}

func (f *filter) match(s string) bool {
	includeOK := true
	excludeOK := true

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

	if raw.Include != "" {
		m, err := newMatcher(raw.Include)
		if err != nil {
			return nil, err
		}
		f.include = m

	}

	if raw.Exclude != "" {
		m, err := newMatcher(raw.Exclude)
		if err != nil {
			return nil, err
		}
		f.exclude = m
	}

	return &f, nil
}

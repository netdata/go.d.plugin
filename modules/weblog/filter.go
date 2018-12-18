package weblog

type rawFilter struct {
	Include string
	Exclude string
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
	if raw.Include == "" && raw.Exclude == "" {
		return nil, nil
	}

	var (
		fil filter
		err error
	)

	if raw.Include != "" {
		if fil.include, err = newMatcher(raw.Include); err != nil {
			return nil, err
		}
	}

	if raw.Exclude != "" {
		if fil.exclude, err = newMatcher(raw.Exclude); err != nil {
			return nil, err
		}
	}

	return &fil, nil
}

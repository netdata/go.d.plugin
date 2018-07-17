package web_log

type filter struct {
	include matcher
	exclude matcher
}

func (f *filter) exist() bool {
	return f.include != nil || f.exclude != nil
}

func (f *filter) filter(s string) bool {
	i, e := true, true
	if f.include != nil {
		i = f.include.match(s)
	}
	if f.exclude != nil {
		e = !f.exclude.match(s)
	}
	return i && e
}

type rawFilter struct {
	Include string `yaml:"include"`
	Exclude string `yaml:"exclude"`
}

func getFilter(r rawFilter) (filter, error) {
	var f filter
	if r.Include == "" && r.Exclude == "" {
		return f, nil
	}

	if r.Include != "" {
		m, err := getMatcher(r.Include)
		if err != nil {
			return f, err
		}
		f.include = m
	}

	if r.Exclude != "" {
		m, err := getMatcher(r.Exclude)
		if err != nil {
			return f, err
		}
		f.exclude = m
	}

	return f, nil
}

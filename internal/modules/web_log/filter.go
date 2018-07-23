package web_log

import "github.com/l2isbad/go.d.plugin/internal/modules/web_log/matcher"

type filter struct {
	include matcher.Matcher
	exclude matcher.Matcher
}

func (f *filter) exist() bool {
	return f.include != nil || f.exclude != nil
}

func (f *filter) filter(s string) bool {
	i, e := true, true
	if f.include != nil {
		i = f.include.Match(s)
	}
	if f.exclude != nil {
		e = !f.exclude.Match(s)
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
		m, err := matcher.New(r.Include)
		if err != nil {
			return f, err
		}
		f.include = m
	}

	if r.Exclude != "" {
		m, err := matcher.New(r.Exclude)
		if err != nil {
			return f, err
		}
		f.exclude = m
	}

	return f, nil
}

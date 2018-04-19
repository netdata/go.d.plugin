package web_log

import (
	"strings"
	"regexp"
)

type filter interface {
	match(string) bool
}

type rawFilter struct {
	Include  string `toml:"include"`
	Exclude  string `toml:"exclude"`
	UseRegex bool   `toml:"use_regex"`
}

type strFilter struct {
	include string
	exclude string
}

func (f *strFilter) match(s string) bool {
	if f.include != "" {
		return strings.Contains(s, f.include)
	}
	if f.exclude != "" {
		return !strings.Contains(s, f.exclude)
	}
	return true
}

type regexFilter struct {
	include *regexp.Regexp
	exclude *regexp.Regexp
}

func (f *regexFilter) match(s string) bool {
	if f.include != nil {
		return f.include.MatchString(s)
	}
	if f.exclude != nil {
		return !f.exclude.MatchString(s)
	}
	return true
}

func getFilter(f rawFilter) (filter, error) {
	if f.Include == "" && f.Exclude == "" {
		return nil, nil
	}
	if !f.UseRegex {
		return &strFilter{f.Include, f.Exclude}, nil
	}

	rf := &regexFilter{}
	if f.Include != "" {
		if r, err := regexp.Compile(f.Include); err != nil {
			return nil, err
		} else {
			rf.include = r
		}
	}
	if f.Exclude != "" {
		if r, err := regexp.Compile(f.Exclude); err != nil {
			return nil, err
		} else {
			rf.exclude = r
		}
	}
	return rf, nil
}

package web_log

import (
	"regexp"
	"strings"
)

type filter interface {
	match(string) bool
}

type rawFilter struct {
	Include  string `yaml:"include"`
	Exclude  string `yaml:"exclude"`
	UseRegex bool   `yaml:"use_regex"`
}

type strFilter struct {
	include string
	exclude string
}

func (f *strFilter) match(s string) bool {
	if f.include != "" && !strings.Contains(s, f.include) {
		return false
	}
	if f.exclude != "" && strings.Contains(s, f.exclude) {
		return false
	}
	return true
}

type regexFilter struct {
	include *regexp.Regexp
	exclude *regexp.Regexp
}

func (f *regexFilter) match(s string) bool {
	if f.include != nil && !f.include.MatchString(s) {
		return false
	}
	if f.exclude != nil && f.exclude.MatchString(s) {
		return false
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

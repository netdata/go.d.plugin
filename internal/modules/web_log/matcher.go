package web_log

import (
	"regexp"
	"strings"
)

type matcher interface {
	compile() error
	match(string) bool
}

func newMatcher(f matchRaw) matcher {
	if f.UseRegex {
		return &matchRe{
			include: f.Include,
			exclude: f.Exclude,
		}
	}
	return &matchStr{
		include: f.Include,
		exclude: f.Exclude,
	}
}

type matchRaw struct {
	Include  string `yaml:"include"`
	Exclude  string `yaml:"exclude"`
	UseRegex bool   `yaml:"use_regex"`
}

func (m matchRaw) exist() bool {
	return m.Include != "" || m.Exclude != ""
}

type matchStr struct {
	include string
	exclude string
}

func (m *matchStr) compile() error {
	return nil
}

func (m *matchStr) match(s string) bool {
	i, e := true, true
	if m.include != "" {
		i = strings.Contains(s, m.include)
	}
	if m.exclude != "" {
		e = !strings.Contains(s, m.exclude)
	}
	return i && e
}

type matchRe struct {
	include string
	exclude string
	inc     *regexp.Regexp
	exc     *regexp.Regexp
}

func (m *matchRe) compile() error {
	if m.include != "" {
		r, err := regexp.Compile(m.include)
		if err != nil {
			return err
		}
		m.inc = r
	}

	if m.exclude != "" {
		r, err := regexp.Compile(m.exclude)
		if err != nil {
			return err
		}
		m.exc = r
	}
	return nil
}

func (m *matchRe) match(s string) bool {
	i, e := true, true
	if m.inc != nil {
		i = m.inc.MatchString(s)
	}
	if m.exc != nil {
		e = !m.exc.MatchString(s)
	}
	return i && e
}

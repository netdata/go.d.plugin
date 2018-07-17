package web_log

import (
	"regexp"
	"strings"
)

type matcher interface {
	match(string) bool
}

type stringMatch struct {
	v string
}

func (m stringMatch) match(s string) bool {
	return strings.Contains(s, m.v)
}

type regexMatch struct {
	v *regexp.Regexp
}

func (m regexMatch) match(s string) bool {
	return m.v.MatchString(s)
}

// TODO: super simple and super questionable, simple > questionable?
func isStringRegex(s string) bool {
	return strings.HasPrefix(s, "s:")
}

func getMatcher(s string) (matcher, error) {
	if isStringRegex(s) {
		return stringMatch{s[2:]}, nil
	}

	re, err := regexp.Compile(s)

	if err != nil {
		return nil, err
	}

	return regexMatch{re}, nil
}

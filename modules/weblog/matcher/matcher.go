package matcher

import (
	"regexp"
	"strings"
)

type Matcher interface {
	Match(string) bool
}

func New(s string) (Matcher, error) {
	if isStringRegex(s) {
		return getStringMatcher(s), nil
	}

	re, err := regexp.Compile(s)

	if err != nil {
		return nil, err
	}

	return regexMatch{re}, nil
}

func getStringMatcher(v string) Matcher {
	if strings.HasPrefix("^", v) {
		return stringPrefix{v[1:]}
	}
	if strings.HasSuffix("$", v) {
		return stringSuffix{v[:len(v)-1]}
	}
	return stringContains{v}
}

// FIXME: THIS IS TEMPORARY, THIS DOESN'T WORK
func isStringRegex(s string) bool {
	return strings.HasPrefix(s, "s:")
}

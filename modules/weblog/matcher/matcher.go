package matcher

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	methodString = "string"
	methodRegexp = "regexp"
)

// Matcher Matcher
type Matcher interface {
	Match(string) bool
}

// New creates new Matcher.
// It returns Matcher and error if any.
// Examples:
// 'string=GET'
// 'string=^GOT'
// 'regexp=G[QWERTY]T'
func New(rawExpr string) (Matcher, error) {
	v := strings.SplitN(rawExpr, "=", 2)

	if len(v) == 2 && v[1] == "" || len(v) != 2 {
		return nil, fmt.Errorf("unsupported match syntax: %s", rawExpr)
	}

	method, expr := v[0], v[1]

	switch method {
	case methodString:
		return stringMatcherFactory(expr), nil
	case methodRegexp:
		return regexpMatcherFactory(expr)
	}

	return nil, fmt.Errorf("unsupported match method: %s", method)
}

func stringMatcherFactory(expr string) Matcher {
	if strings.HasPrefix(expr, "^") {
		return &stringPrefix{expr[1:]}
	}
	if strings.HasSuffix(expr, "$") {
		return &stringSuffix{expr[:len(expr)-1]}
	}

	return &stringContains{expr}
}

func regexpMatcherFactory(expr string) (Matcher, error) {
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	return &regexMatch{re}, nil
}

type regexMatch struct {
	v *regexp.Regexp
}

func (m regexMatch) Match(s string) bool {
	return m.v.MatchString(s)
}

type stringContains struct {
	v string
}

func (m stringContains) Match(s string) bool {
	return strings.Contains(s, m.v)
}

type stringPrefix struct {
	v string
}

func (m stringPrefix) Match(s string) bool {
	return strings.HasPrefix(s, m.v)
}

type stringSuffix struct {
	v string
}

func (m stringSuffix) Match(s string) bool {
	return strings.HasSuffix(s, m.v)
}

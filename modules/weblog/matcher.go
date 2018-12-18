package weblog

import (
	"fmt"
	"regexp"
	"strings"
)

type matcher interface {
	match(string) bool
}

// Valid options:
// 'string=GET'
// 'string=^GOT'
// 'regexp=G[QWERTY]T'
func newMatcher(rawExpr string) (matcher, error) {
	v := strings.SplitN(rawExpr, "=", 2)

	if len(v) == 2 && v[1] == "" || len(v) != 2 {
		return nil, fmt.Errorf("unsupported match syntax: %s", rawExpr)
	}

	method, expr := v[0], v[1]

	switch method {
	case "string":
		return newStringMatcher(expr), nil
	case "regexp":
		return newRegexpMatcher(expr)
	}

	return nil, fmt.Errorf("unsupported Match method: %s", method)
}

func newStringMatcher(expr string) matcher {
	if strings.HasPrefix(expr, "^") {
		return &stringPrefixMatcher{expr[1:]}
	}
	if strings.HasSuffix(expr, "$") {
		return &stringSuffixMatcher{expr[:len(expr)-1]}
	}

	return &stringContainsMatcher{expr}
}

func newRegexpMatcher(expr string) (matcher, error) {
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	return &regexpMatcher{re}, nil
}

type stringContainsMatcher struct{ v string }

func (m stringContainsMatcher) match(s string) bool { return strings.Contains(s, m.v) }

type stringPrefixMatcher struct{ v string }

func (m stringPrefixMatcher) match(s string) bool { return strings.HasPrefix(s, m.v) }

type stringSuffixMatcher struct{ v string }

func (m stringSuffixMatcher) match(s string) bool { return strings.HasSuffix(s, m.v) }

type regexpMatcher struct{ v *regexp.Regexp }

func (m regexpMatcher) match(s string) bool { return m.v.MatchString(s) }

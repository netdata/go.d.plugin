package matcher

import "strings"

func stringMatcherFactory(expr string) Matcher {
	if strings.HasPrefix(expr, "^") {
		return &StringPrefixMatcher{expr[1:]}
	}
	if strings.HasSuffix(expr, "$") {
		return &StringSuffixMatcher{expr[:len(expr)-1]}
	}

	return &StringContainsMatcher{expr}
}

type StringContainsMatcher struct {
	v string
}

func (m StringContainsMatcher) Match(s string) bool {
	return strings.Contains(s, m.v)
}

type StringPrefixMatcher struct {
	v string
}

func (m StringPrefixMatcher) Match(s string) bool {
	return strings.HasPrefix(s, m.v)
}

type StringSuffixMatcher struct {
	v string
}

func (m StringSuffixMatcher) Match(s string) bool {
	return strings.HasSuffix(s, m.v)
}

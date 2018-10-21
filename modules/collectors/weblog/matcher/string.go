package matcher

import "strings"

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

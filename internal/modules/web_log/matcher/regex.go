package matcher

import "regexp"

type regexMatch struct {
	v *regexp.Regexp
}

func (m regexMatch) Match(s string) bool {
	return m.v.MatchString(s)
}

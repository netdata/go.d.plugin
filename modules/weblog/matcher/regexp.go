package matcher

import "regexp"

func regexpMatcherFactory(expr string) (Matcher, error) {
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	return &RegexpMatcher{re}, nil
}

type RegexpMatcher struct {
	v *regexp.Regexp
}

func (m RegexpMatcher) Match(s string) bool {
	return m.v.MatchString(s)
}

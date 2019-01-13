package matcher

type (
	trueMatcher  struct{}
	falseMatcher struct{}
	andMatcher   struct{ lhs, rhs Matcher }
	orMatcher    struct{ lhs, rhs Matcher }
	negMatcher   struct{ Matcher }
)

var (
	matcherT trueMatcher
	matcherF falseMatcher
)

// TRUE returns a matcher which always returns true
func TRUE() Matcher {
	return matcherT
}

// FALSE returns a matcher which always returns false
func FALSE() Matcher {
	return matcherF
}

// Not returns a matcher which negative the sub-matcher's result
func Not(m Matcher) Matcher {
	return negMatcher{m}
}

// And returns a matcher which returns true only if all of it's sub-matcher return true
func And(lhs, rhs Matcher) Matcher {
	return andMatcher{lhs, rhs}
}

// Or returns a matcher which returns true if any of it's sub-matcher return true
func Or(lhs, rhs Matcher) Matcher {
	return orMatcher{lhs, rhs}
}

func (trueMatcher) Match(b []byte) bool       { return true }
func (trueMatcher) MatchString(s string) bool { return true }

func (falseMatcher) Match(b []byte) bool       { return false }
func (falseMatcher) MatchString(s string) bool { return false }

func (m andMatcher) Match(b []byte) bool       { return m.lhs.Match(b) && m.rhs.Match(b) }
func (m andMatcher) MatchString(s string) bool { return m.lhs.MatchString(s) && m.rhs.MatchString(s) }

func (m orMatcher) Match(b []byte) bool       { return m.lhs.Match(b) || m.rhs.Match(b) }
func (m orMatcher) MatchString(s string) bool { return m.lhs.MatchString(s) || m.rhs.MatchString(s) }

func (m negMatcher) Match(b []byte) bool       { return !m.Matcher.Match(b) }
func (m negMatcher) MatchString(s string) bool { return !m.Matcher.MatchString(s) }

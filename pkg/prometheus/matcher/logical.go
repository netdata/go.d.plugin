package matcher

import (
	"github.com/prometheus/prometheus/pkg/labels"
)

type (
	trueMatcher  struct{}
	falseMatcher struct{}
	negMatcher   struct{ m Matcher }
	andMatcher   struct{ lhs, rhs Matcher }
	orMatcher    struct{ lhs, rhs Matcher }
)

func (trueMatcher) Matches(_ labels.Labels) bool    { return true }
func (falseMatcher) Matches(_ labels.Labels) bool   { return false }
func (m negMatcher) Matches(lbs labels.Labels) bool { return !m.m.Matches(lbs) }
func (m andMatcher) Matches(lbs labels.Labels) bool { return m.lhs.Matches(lbs) && m.rhs.Matches(lbs) }
func (m orMatcher) Matches(lbs labels.Labels) bool  { return m.lhs.Matches(lbs) || m.rhs.Matches(lbs) }

// And returns a matcher which returns true only if all of it's sub-matcher return true
func And(lhs, rhs Matcher, others ...Matcher) Matcher {
	m := andMatcher{lhs: lhs, rhs: rhs}
	if len(others) == 0 {
		return m
	}
	return And(m, others[0], others[1:]...)
}

// Or returns a matcher which returns true if any of it's sub-matcher return true
func Or(lhs, rhs Matcher, others ...Matcher) Matcher {
	m := orMatcher{lhs: lhs, rhs: rhs}
	if len(others) == 0 {
		return m
	}
	return Or(m, others[0], others[1:]...)
}

// Not returns a matcher which opposites the sub-matcher's result
func Not(m Matcher) Matcher {
	return negMatcher{m}
}
